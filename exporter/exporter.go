// Package exporter
// Time    : 2021/7/22 1:25 下午
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type BuildInfo struct {
	Version   string
	CommitSha string
	Date      string
}

type Options struct {
	Namespace    string
	EndPoint     string
	LogstashName string
	Hostname     string
	MetricsPath  string
	Registry     *prometheus.Registry
	BuildInfo    BuildInfo
}

// LogstashExporter implements the prometheus.Exporter interface, and exports Logstash metrics.
type LogstashExporter struct {
	sync.Mutex

	namespace string
	endpoint  string

	totalScrapes   prometheus.Counter
	scrapeDuration prometheus.Summary

	metricDescriptions map[string]*prometheus.Desc

	options   Options
	mux       *http.ServeMux
	buildInfo BuildInfo
}

func NewLogstashExporter(opts Options) (*LogstashExporter, error) {
	log.Debugf("LogstashExporter options: %#v", opts)
	e := &LogstashExporter{
		namespace: opts.Namespace,
		endpoint:  opts.EndPoint,
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Name:      "exporter_scrapes_total",
			Help:      "Current total logstash scrapes.",
			ConstLabels: map[string]string{
				"hostname":      opts.Hostname,
				"logstash_name": opts.LogstashName,
			},
		}),
		scrapeDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: opts.Namespace,
			Name:      "exporter_scrape_duration_seconds",
			Help:      "Durations of scrapes by the exporter",
			ConstLabels: map[string]string{
				"hostname":      opts.Hostname,
				"logstash_name": opts.LogstashName,
			},
		}),

		options:   opts,
		buildInfo: opts.BuildInfo,
	}
	e.metricDescriptions = map[string]*prometheus.Desc{}
	for k, desc := range map[string]struct {
		txt    string
		labels []string
	}{
		"up": {txt: "Information about the Redis instance", labels: []string{"hostname", "logstash_name"}},
	} {
		e.metricDescriptions[k] = newMetricDesc(opts.Namespace, k, desc.txt, desc.labels)
	}

	e.mux = http.NewServeMux()

	if e.options.Registry != nil {
		e.options.Registry.MustRegister(e)
		e.mux.Handle(e.options.MetricsPath, promhttp.HandlerFor(
			e.options.Registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError},
		))
	}

	e.mux.HandleFunc("/", e.indexHandler)
	e.mux.HandleFunc("/health", e.healthHandler)

	return e, nil
}

// Describe outputs Redis metric descriptions.
func (e *LogstashExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range e.metricDescriptions {
		ch <- desc
	}
	//
	// for _, v := range e.metricMapGauges {
	// 	ch <- newMetricDescr(e.options.Namespace, v, v+" metric", nil)
	// }
	//
	// for _, v := range e.metricMapCounters {
	// 	ch <- newMetricDescr(e.options.Namespace, v, v+" metric", nil)
	// }

	ch <- e.totalScrapes.Desc()
	ch <- e.scrapeDuration.Desc()
}

// Collect fetches new metrics from the RedisHost and updates the appropriate metrics.
func (e *LogstashExporter) Collect(ch chan<- prometheus.Metric) {
	e.Lock()
	defer e.Unlock()
	e.totalScrapes.Inc()

	startTime := time.Now()
	var up float64
	up = 1
	e.registerConstMetricGauge(ch, "up", up, e.options.Hostname, e.options.LogstashName)

	took := time.Since(startTime).Seconds()
	e.scrapeDuration.Observe(took)

	ch <- e.totalScrapes
	ch <- e.scrapeDuration
}
