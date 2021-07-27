// Package exporter
// Time    : 2021/7/22 3:29 下午
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package exporter

import (
	"fmt"
	"net/http"
)

func (e *LogstashExporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.mux.ServeHTTP(w, r)
}

func (e *LogstashExporter) healthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(fmt.Sprintf("%s-%s exporter is running", e.options.Namespace, e.options.LogstashUsage)))
}

func (e *LogstashExporter) indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`<html>
<head><title>Logstash Exporter ` + e.buildInfo.Version + `</title></head>
<body>
<h1>Logstash Exporter ` + e.buildInfo.Version + `</h1>
<p><a href='` + e.options.MetricsPath + `'>Metrics</a></p>
</body>
</html>
`))
}
