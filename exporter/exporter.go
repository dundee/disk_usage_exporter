package exporter

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dundee/gdu/v5/pkg/analyze"
	"github.com/dundee/gdu/v5/pkg/fs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_disk_usage_bytes",
			Help: "Disk usage of the directory/file",
		},
		[]string{"path"},
	)
	diskUsageLevel1 = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_disk_usage_level_1_bytes",
			Help: "Disk usage of the directory/file level 1",
		},
		[]string{"path"},
	)
	collectors = []prometheus.Collector{diskUsage, diskUsageLevel1}
)

func init() {
	prometheus.MustRegister(collectors...)
}

// Exporter is the type to be used to start HTTP server and run the analysis
type Exporter struct {
	paths          map[string]int
	ignoreDirPaths map[string]struct{}
	followSymlinks bool
}

// NewExporter creates new Exporter
func NewExporter(paths map[string]int, followSymlinks bool) *Exporter {
	return &Exporter{
		followSymlinks: followSymlinks,
		paths:          paths,
	}
}

func (e *Exporter) runAnalysis() {
	defer debug.FreeOSMemory()

	for path, level := range e.paths {
		analyzer := analyze.CreateAnalyzer()
		analyzer.SetFollowSymlinks(e.followSymlinks)
		dir := analyzer.AnalyzeDir(path, e.shouldDirBeIgnored, false)
		dir.UpdateStats(fs.HardLinkedItems{})
		e.reportItem(dir, 0, level)
	}

	log.Info("Analysis done")
}

// SetIgnoreDirPaths sets paths to ignore
func (e *Exporter) SetIgnoreDirPaths(paths []string) {
	e.ignoreDirPaths = make(map[string]struct{}, len(paths))
	for _, path := range paths {
		e.ignoreDirPaths[path] = struct{}{}
	}
}

func (e *Exporter) shouldDirBeIgnored(name, path string) bool {
	_, ok := e.ignoreDirPaths[path]
	return ok
}

func (e *Exporter) reportItem(item fs.Item, level, maxLevel int) {
	if level == maxLevel {
		diskUsage.WithLabelValues(item.GetPath()).Set(float64(item.GetUsage()))
	} else if level == 1 {
		diskUsageLevel1.WithLabelValues(item.GetPath()).Set(float64(item.GetUsage()))
	}

	if item.IsDir() && level+1 <= maxLevel {
		for _, entry := range item.GetFiles() {
			e.reportItem(entry, level+1, maxLevel)
		}
	}
}

func (e *Exporter) WriteToTextfile(name string) {
	e.runAnalysis()

	// Use a custom registry to drop go stats
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors...)

	err := prometheus.WriteToTextfile(name, registry)
	if err != nil {
		log.Fatalf("WriteToTextfile error: %v", err)
	}
	log.Printf("Stored stats in file %s", name)
}

func (e *Exporter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.runAnalysis()
	promhttp.Handler().ServeHTTP(w, req)
}

// RunServer starts HTTP server loop
func (e *Exporter) RunServer(addr string) {
	http.Handle("/", http.HandlerFunc(ServeIndex))
	http.Handle("/metrics", e)

	log.Printf("Providing metrics at http://%s/metrics", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// ServeIndex serves index page
func ServeIndex(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	res := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
	<meta name="viewport" content="width=device-width">
	<title>Disk Usage Prometheus Exporter</title>
</head>
<body>
<h1>Disk Usage Prometheus Exporter</h1>
<p>
	<a href="/metrics">Metrics</a>
</p>
<p>
	<a href="https://github.com/dundee/disk_usage_exporter">Homepage</a>
</p>
</body>
</html>
`
	fmt.Fprint(w, res)
}
