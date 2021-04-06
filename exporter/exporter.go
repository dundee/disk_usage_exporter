package exporter

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dundee/gdu/v4/analyze"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_disk_usage_bytes",
			Help: "Disk usage of the directory/file",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(diskUsage)
}

type Exporter struct {
	analyzer       analyze.Analyzer
	ignoreDirPaths map[string]struct{}
	maxLevel       int
}

func NewExporter(maxLevel int) *Exporter {
	return &Exporter{
		maxLevel: maxLevel,
		analyzer: analyze.CreateAnalyzer(),
	}
}

func (e *Exporter) Run(path string) {
	dir := e.analyzer.AnalyzeDir(path, e.ShouldDirBeIgnored)
	e.ReportItem(dir, 0)
}

// SetIgnoreDirPaths sets paths to ignore
func (e *Exporter) SetIgnoreDirPaths(paths []string) {
	e.ignoreDirPaths = make(map[string]struct{}, len(paths))
	for _, path := range paths {
		e.ignoreDirPaths[path] = struct{}{}
	}
}

// ShouldDirBeIgnored returns true if given path should be ignored
func (e *Exporter) ShouldDirBeIgnored(path string) bool {
	_, ok := e.ignoreDirPaths[path]
	return ok
}

func (e *Exporter) ReportItem(item analyze.Item, level int) {
	if level == e.maxLevel {
		diskUsage.WithLabelValues(item.GetPath()).Set(float64(item.GetUsage()))
	}

	if item.IsDir() && level+1 <= e.maxLevel {
		for _, entry := range item.(*analyze.Dir).Files {
			e.ReportItem(entry, level+1)
		}
	}
}

func RunAnalysis(path string, ignoreDirs []string, level int) {
	exporter := NewExporter(level)
	exporter.SetIgnoreDirPaths(ignoreDirs)

	for {
		exporter.Run(path)
		time.Sleep(time.Minute * 1)
	}
}

func RunServer(addr string) {
	http.Handle("/", http.HandlerFunc(serveIndex))
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func serveIndex(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	res := `<!DOCTYPE>
<html>
<body>
<h1>go DiskUsage("Prometheus Exporter")</h1>
<a href="/metrics">Metrics</a>
</html>
`
	fmt.Fprint(w, res)
}
