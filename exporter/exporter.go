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

const RootPath = "/"

type BuildInfo struct {
	Version   string
	CommitSha string
	Date      string
}

type Options struct {
	Namespace     string
	EndPoint      string
	LogstashUsage string
	Hostname      string
	MetricsPath   string
	Registry      *prometheus.Registry
	BuildInfo     BuildInfo
}

type Collector interface {
	Collect(ch chan<- prometheus.Metric)
}

// LogstashExporter implements the prometheus.Exporter interface, and exports Logstash metrics.
type LogstashExporter struct {
	sync.Mutex

	namespace string
	endpoint  string

	totalScrapes   *prometheus.CounterVec
	scrapeDuration *prometheus.SummaryVec
	logstashUp     *prometheus.GaugeVec

	reqClient  *ReqClient
	collectors []Collector

	options   Options
	mux       *http.ServeMux
	buildInfo BuildInfo
}

func NewLogstashExporter(opts Options) (*LogstashExporter, error) {
	log.Debugf("LogstashExporter options: %#v", opts)
	e := &LogstashExporter{
		namespace: opts.Namespace,
		endpoint:  opts.EndPoint,
		totalScrapes: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Name:      "exporter_scrapes_total",
			Help:      "Current total logstash scrapes.",
		}, []string{"hostname", "logstash_usage"}),
		scrapeDuration: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: opts.Namespace,
			Name:      "exporter_scrape_duration_seconds",
			Help:      "Durations of scrapes by the exporter",
		}, []string{"hostname", "logstash_usage"}),
		logstashUp: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: opts.Namespace,
			Name:      "up",
			Help:      "Information about the Logstash instance",
		}, []string{"hostname", "logstash_usage"}),

		options:   opts,
		buildInfo: opts.BuildInfo,
	}

	e.reqClient = NewReqClient(opts.EndPoint)

	nodeStatCollector, _ := NewNodeStatsCollector(e)

	e.collectors = append(e.collectors, nodeStatCollector)

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
func (e *LogstashExporter) Describe(_ chan<- *prometheus.Desc) {
}

// Collect fetches new metrics from the RedisHost and updates the appropriate metrics.
func (e *LogstashExporter) Collect(ch chan<- prometheus.Metric) {
	e.Lock()
	defer e.Unlock()
	e.totalScrapes.WithLabelValues(e.options.Hostname, e.options.LogstashUsage).Inc()
	startTime := time.Now()
	nrc := NewReqClient(e.endpoint)

	rootInfo, err := GetLogstashRootInfo(nrc, RootPath)
	if err != nil {
		e.logstashUp.WithLabelValues(e.options.Hostname, e.options.LogstashUsage).Set(0)
		log.Errorf("request %s%s", e.endpoint, RootPath)
	} else {
		e.logstashUp.WithLabelValues(rootInfo.Host, e.options.LogstashUsage).Set(1)
		wg := sync.WaitGroup{}
		wg.Add(len(e.collectors))
		for _, c := range e.collectors {
			go func(c Collector) {
				c.Collect(ch)
				wg.Done()
			}(c)
		}
		wg.Wait()
		took := time.Since(startTime).Seconds()
		e.scrapeDuration.WithLabelValues(e.options.Hostname, e.options.LogstashUsage).Observe(took)
	}
	e.logstashUp.Collect(ch)
	e.totalScrapes.Collect(ch)
	e.scrapeDuration.Collect(ch)
}
