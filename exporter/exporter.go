package exporter

import (
	"fmt"
	"net/http"
	"path/filepath"

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
)

func init() {
	prometheus.MustRegister(diskUsage)
	prometheus.MustRegister(diskUsageLevel1)
}

// Exporter is the type to be used to start HTTP server and run the analysis
type Exporter struct {
	ignoreDirPaths map[string]struct{}
	maxLevel       int
	path           string
}

// NewExporter creates new Exporter
func NewExporter(maxLevel int, path string) *Exporter {
	return &Exporter{
		maxLevel: maxLevel,
		path:     filepath.Clean(path),
	}
}

func (e *Exporter) runAnalysis() {
	analyzer := analyze.CreateAnalyzer()
	dir := analyzer.AnalyzeDir(e.path, e.shouldDirBeIgnored, false)
	dir.UpdateStats(fs.HardLinkedItems{})
	e.reportItem(dir, 0)
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

func (e *Exporter) reportItem(item fs.Item, level int) {
	if level == e.maxLevel {
		diskUsage.WithLabelValues(item.GetPath()).Set(float64(item.GetUsage()))
	} else if level == 1 {
		diskUsageLevel1.WithLabelValues(item.GetPath()).Set(float64(item.GetUsage()))
	}

	if item.IsDir() && level+1 <= e.maxLevel {
		for _, entry := range item.GetFiles() {
			e.reportItem(entry, level+1)
		}
	}
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
