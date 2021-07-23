// Package exporter
// Time    : 2021/7/22 2:24 下午
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package exporter

import "github.com/prometheus/client_golang/prometheus"

func newMetricDesc(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "", metricName), docString, labels, nil)
}

func (e *LogstashExporter) registerConstMetricGauge(ch chan<- prometheus.Metric, metric string, val float64, labels ...string) {
	e.registerConstMetric(ch, metric, val, prometheus.GaugeValue, labels...)
}

func (e *LogstashExporter) registerConstMetric(ch chan<- prometheus.Metric, metric string, val float64, valType prometheus.ValueType, labelValues ...string) {
	desc := e.metricDescriptions[metric]
	if desc == nil {
		desc = newMetricDesc(e.options.Namespace, metric, metric+" metric", labelValues)
	}

	if m, err := prometheus.NewConstMetric(desc, valType, val, labelValues...); err == nil {
		ch <- m
	}
}
