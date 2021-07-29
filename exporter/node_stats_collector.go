// Package exporter
// Time    : 2021/7/26 2:10 下午
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// NodeStatsCollector type
type NodeStatsCollector struct {
	export  *LogstashExporter
	ReqPath string

	LogstashInfo *prometheus.Desc

	JvmThreadsCount     *prometheus.Desc
	JvmThreadsPeakCount *prometheus.Desc

	MemHeapUsedInBytes         *prometheus.Desc
	MemHeapCommittedInBytes    *prometheus.Desc
	MemHeapMaxInBytes          *prometheus.Desc
	MemHeapUsedPercent         *prometheus.Desc
	MemNonHeapUsedInBytes      *prometheus.Desc
	MemNonHeapCommittedInBytes *prometheus.Desc

	MemPoolPeakUsedInBytes  *prometheus.Desc
	MemPoolUsedInBytes      *prometheus.Desc
	MemPoolPeakMaxInBytes   *prometheus.Desc
	MemPoolMaxInBytes       *prometheus.Desc
	MemPoolCommittedInBytes *prometheus.Desc

	GCCollectionTimeInMillis *prometheus.Desc
	GCCollectionCount        *prometheus.Desc

	ProcessOpenFileDescriptors     *prometheus.Desc
	ProcessPeakOpenFileDescriptors *prometheus.Desc
	ProcessMaxFileDescriptors      *prometheus.Desc
	ProcessMemTotalVirtualInBytes  *prometheus.Desc
	ProcessCPUTotalInMillis        *prometheus.Desc
	ProcessCPUPercent              *prometheus.Desc

	PipelineDuration       *prometheus.Desc
	PipelineEventsIn       *prometheus.Desc
	PipelineEventsFiltered *prometheus.Desc
	PipelineEventsOut      *prometheus.Desc
}

func NewNodeStatsCollector(e *LogstashExporter) (*NodeStatsCollector, error) {
	const subsystem = "node_stats"
	return &NodeStatsCollector{
		export:  e,
		ReqPath: "/_node/stats",

		LogstashInfo: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, "", "instance_info"),
			"instance_info",
			[]string{"version", "http_address"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		JvmThreadsCount: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "jvm_threads_count"),
			"jvm_threads_count",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		JvmThreadsPeakCount: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "jvm_threads_peak_count"),
			"jvm_threads_peak_count",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemHeapUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_heap_used_bytes"),
			"mem_heap_used_bytes",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemHeapUsedPercent: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "heap_used_percent"),
			"heap_used_percent",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemHeapCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_heap_committed_bytes"),
			"mem_heap_committed_bytes",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemHeapMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_heap_max_bytes"),
			"mem_heap_max_bytes",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemNonHeapUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_nonheap_used_bytes"),
			"mem_nonheap_used_bytes",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemNonHeapCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_nonheap_committed_bytes"),
			"mem_nonheap_committed_bytes",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemPoolUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_pool_used_bytes"),
			"mem_pool_used_bytes",
			[]string{"pool"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemPoolPeakUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_pool_peak_used_bytes"),
			"mem_pool_peak_used_bytes",
			[]string{"pool"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemPoolPeakMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_pool_peak_max_bytes"),
			"mem_pool_peak_max_bytes",
			[]string{"pool"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemPoolMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_pool_max_bytes"),
			"mem_pool_max_bytes",
			[]string{"pool"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		MemPoolCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "mem_pool_committed_bytes"),
			"mem_pool_committed_bytes",
			[]string{"pool"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		GCCollectionTimeInMillis: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "gc_collection_duration_seconds_total"),
			"gc_collection_duration_seconds_total",
			[]string{"collector"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		GCCollectionCount: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "gc_collection_total"),
			"gc_collection_total",
			[]string{"collector"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		ProcessOpenFileDescriptors: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "process_open_file_descriptors"),
			"process_open_file_descriptors",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		ProcessPeakOpenFileDescriptors: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "process_peak_open_file_descriptors"),
			"process_peak_open_file_descriptors",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		ProcessMaxFileDescriptors: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "process_max_file_descriptors"),
			"process_max_file_descriptors",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		ProcessMemTotalVirtualInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "process_mem_total_virtual_bytes"),
			"process_mem_total_virtual_bytes",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		ProcessCPUTotalInMillis: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "process_cpu_total_seconds_total"),
			"process_cpu_total_seconds_total",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		ProcessCPUPercent: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "process_cpu_percent"),
			"process_cpu_percent",
			nil,
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		PipelineDuration: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "pipeline_duration_seconds_total"),
			"pipeline_duration_seconds_total",
			[]string{"pipeline"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		PipelineEventsIn: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "pipeline_events_in_total"),
			"pipeline_events_in_total",
			[]string{"pipeline"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		PipelineEventsFiltered: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "pipeline_events_filtered_total"),
			"pipeline_events_filtered_total",
			[]string{"pipeline"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),

		PipelineEventsOut: prometheus.NewDesc(
			prometheus.BuildFQName(e.namespace, subsystem, "pipeline_events_out_total"),
			"pipeline_events_out_total",
			[]string{"pipeline"},
			map[string]string{"hostname": e.options.Hostname, "logstash_usage": e.options.LogstashUsage},
		),
	}, nil
}

func (c *NodeStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := GetLogstashNodeStats(c.export.reqClient, c.ReqPath, c.export.options.ScrapeTimeoutMillisecond)
	if err != nil {
		log.Fatalf("GetLogstashNodeStats <%s> error: <%#v>", c.ReqPath, err)
	} else {

		ch <- prometheus.MustNewConstMetric(
			c.LogstashInfo,
			prometheus.GaugeValue,
			float64(1),
			stats.Version, stats.HTTPAddress)

		ch <- prometheus.MustNewConstMetric(
			c.JvmThreadsCount,
			prometheus.GaugeValue,
			float64(stats.Jvm.Threads.Count),
		)

		ch <- prometheus.MustNewConstMetric(
			c.JvmThreadsPeakCount,
			prometheus.GaugeValue,
			float64(stats.Jvm.Threads.PeakCount),
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemHeapUsedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.HeapUsedInBytes),
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemHeapUsedPercent,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.HeapUsedPercent),
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemHeapMaxInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.HeapMaxInBytes),
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemHeapCommittedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.HeapCommittedInBytes),
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemNonHeapUsedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.NonHeapUsedInBytes),
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemNonHeapCommittedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.NonHeapCommittedInBytes),
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolPeakUsedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
			"old",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolUsedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.UsedInBytes),
			"old",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolPeakMaxInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
			"old",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolMaxInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.MaxInBytes),
			"old",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolCommittedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.CommittedInBytes),
			"old",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolPeakUsedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
			"young",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolUsedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Young.UsedInBytes),
			"young",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolPeakMaxInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
			"young",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolMaxInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Young.MaxInBytes),
			"young",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolCommittedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Young.CommittedInBytes),
			"young",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolPeakUsedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
			"survivor",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolUsedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Survivor.UsedInBytes),
			"survivor",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolPeakMaxInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
			"survivor",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolMaxInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Survivor.MaxInBytes),
			"survivor",
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemPoolCommittedInBytes,
			prometheus.GaugeValue,
			float64(stats.Jvm.Mem.Pools.Survivor.CommittedInBytes),
			"survivor",
		)

		ch <- prometheus.MustNewConstMetric(
			c.GCCollectionTimeInMillis,
			prometheus.CounterValue,
			float64(stats.Jvm.Gc.Collectors.Old.CollectionTimeInMillis),
			"old",
		)

		ch <- prometheus.MustNewConstMetric(
			c.GCCollectionCount,
			prometheus.GaugeValue,
			float64(stats.Jvm.Gc.Collectors.Old.CollectionCount),
			"old",
		)

		ch <- prometheus.MustNewConstMetric(
			c.GCCollectionTimeInMillis,
			prometheus.CounterValue,
			float64(stats.Jvm.Gc.Collectors.Young.CollectionTimeInMillis),
			"young",
		)

		ch <- prometheus.MustNewConstMetric(
			c.GCCollectionCount,
			prometheus.GaugeValue,
			float64(stats.Jvm.Gc.Collectors.Young.CollectionCount),
			"young",
		)

		ch <- prometheus.MustNewConstMetric(
			c.ProcessOpenFileDescriptors,
			prometheus.GaugeValue,
			float64(stats.Process.OpenFileDescriptors),
		)

		ch <- prometheus.MustNewConstMetric(
			c.ProcessPeakOpenFileDescriptors,
			prometheus.GaugeValue,
			float64(stats.Process.PeakOpenFileDescriptors),
		)

		ch <- prometheus.MustNewConstMetric(
			c.ProcessMaxFileDescriptors,
			prometheus.GaugeValue,
			float64(stats.Process.MaxFileDescriptors),
		)

		ch <- prometheus.MustNewConstMetric(
			c.ProcessMemTotalVirtualInBytes,
			prometheus.GaugeValue,
			float64(stats.Process.Mem.TotalVirtualInBytes),
		)

		ch <- prometheus.MustNewConstMetric(
			c.ProcessCPUTotalInMillis,
			prometheus.CounterValue,
			float64(stats.Process.CPU.TotalInMillis/1000),
		)

		ch <- prometheus.MustNewConstMetric(
			c.ProcessCPUPercent,
			prometheus.GaugeValue,
			float64(stats.Process.CPU.Percent),
		)

		// For backwards compatibility with Logstash 5
		pipelines := make(map[string]Pipeline)
		if len(stats.Pipelines) == 0 {
			pipelines["main"] = stats.Pipeline
		} else {
			pipelines = stats.Pipelines
		}

		for pipelineID, pipeline := range pipelines {
			ch <- prometheus.MustNewConstMetric(
				c.PipelineDuration,
				prometheus.CounterValue,
				float64(pipeline.Events.DurationInMillis/1000),
				pipelineID,
			)

			ch <- prometheus.MustNewConstMetric(
				c.PipelineEventsIn,
				prometheus.CounterValue,
				float64(pipeline.Events.In),
				pipelineID,
			)

			ch <- prometheus.MustNewConstMetric(
				c.PipelineEventsFiltered,
				prometheus.CounterValue,
				float64(pipeline.Events.Filtered),
				pipelineID,
			)

			ch <- prometheus.MustNewConstMetric(
				c.PipelineEventsOut,
				prometheus.CounterValue,
				float64(pipeline.Events.Out),
				pipelineID,
			)
		}
	}
}
