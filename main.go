// Package logstash_exporter
// Time    : 2021/7/22 11:08 上午
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package main

import (
	"github.com/Achillesxu/logstash_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
	"runtime"
)

var (
	/*
		BuildVersion, BuildDate, BuildCommitSha are filled in by the build script
	*/
	BuildVersion   = "<<< filled in by build >>>"
	BuildDate      = "<<< filled in by build >>>"
	BuildCommitSha = "<<< filled in by build >>>"
	NameSpace      = "logstash"
	MetricsPath    = "/metrics"

	logstashEndpoint    string
	exporterBindAddress string
	logstashUsage       string
	isDebug             bool
)

func init() {
	flag.StringVarP(&logstashEndpoint, "logstash_endpoint", "l", "http://localhost:9600", "logstash metric endpoint")
	flag.StringVarP(&exporterBindAddress, "web_listen_address", "w", ":9198", "http server for /metric and more")
	flag.StringVarP(&logstashUsage, "logstash_usage", "u", "logstash", "logstash_usage, for instance: sms, to cope with sms message")
	flag.BoolVar(&isDebug, "debug", false, "Output verbose debug information")
}

func main() {
	flag.Parse()
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		DisableQuote:  true,
		FullTimestamp: true,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "time",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "msg",
		},
	})
	log.SetReportCaller(true)

	switch isDebug {
	case true:
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabling debug output")
	default:
		log.SetLevel(log.InfoLevel)
	}

	log.Infof("Logstash Metrics Exporter %s    build date: %s    sha1: %s    Go: %s    GOOS: %s    GOARCH: %s",
		BuildVersion, BuildDate, BuildCommitSha,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
	instanceHostName, err := os.Hostname()
	if err != nil {
		log.Fatalf("get hostname failed: %#v", err)
	}

	registry := prometheus.NewRegistry()

	exp, err := exporter.NewLogstashExporter(exporter.Options{
		Namespace:     NameSpace,
		LogstashUsage: logstashUsage,
		Hostname:      instanceHostName,
		EndPoint:      logstashEndpoint,
		MetricsPath:   MetricsPath,
		Registry:      registry,
		BuildInfo: exporter.BuildInfo{
			Version:   BuildVersion,
			CommitSha: BuildCommitSha,
			Date:      BuildDate,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Providing metrics at %s%s", exporterBindAddress, MetricsPath)
	log.Infof("logstash_endpoint addr: %s", logstashEndpoint)
	log.Fatal(http.ListenAndServe(exporterBindAddress, exp))
}
