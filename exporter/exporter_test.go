package exporter_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/dundee/disk_usage_exporter/exporter"
	"github.com/stretchr/testify/assert"
)

func TestRunAnalysis(t *testing.T) {
	fin := createTestDir()
	defer fin()

	w := &mockedResponseWriter{}

	e := exporter.NewExporter(2, "test_dir")
	e.SetIgnoreDirPaths([]string{"/proc"})
	e.ServeHTTP(w, &http.Request{
		Header: make(http.Header),
	})

	assert.NotContains(t, string(w.buff), "test_dir/test_dir")
	assert.Contains(t, string(w.buff), "node_disk_usage_bytes")
}

func TestIndex(t *testing.T) {
	w := &mockedResponseWriter{}
	exporter.ServeIndex(w, &http.Request{
		Header: make(http.Header),
	})
	assert.Contains(t, string(w.buff), "<h1>Disk Usage Prometheus Exporter</h1>")
}

type mockedResponseWriter struct {
	buff []byte
}

func (w *mockedResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (w *mockedResponseWriter) Write(b []byte) (int, error) {
	w.buff = append(w.buff, b...)
	return len(b), nil
}

func (w *mockedResponseWriter) WriteHeader(statusCode int) {

}

func createTestDir() func() {
	os.MkdirAll("test_dir/nested/subnested", os.ModePerm)
	os.WriteFile("test_dir/nested/subnested/file", []byte("hello"), 0644)
	os.WriteFile("test_dir/nested/file2", []byte("go"), 0644)
	return func() {
		err := os.RemoveAll("test_dir")
		if err != nil {
			panic(err)
		}
	}
}
