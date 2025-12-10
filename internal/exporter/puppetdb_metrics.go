package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// PuppetDBMetrics 定义PuppetDB核心性能指标
type PuppetDBMetrics struct {
	// 命令处理指标
	commandsProcessed          *prometheus.CounterVec
	commandsProcessingDuration *prometheus.HistogramVec
	commandQueueDepth          prometheus.Gauge

	// 存储层指标
	storageDuplicatePct           prometheus.Gauge
	storageGCDuration             prometheus.Gauge
	storageReplaceFactsDuration   prometheus.Gauge
	storageReplaceCatalogDuration prometheus.Gauge

	// 人口统计指标
	populationNodes               prometheus.Gauge
	populationResources           prometheus.Gauge
	populationAvgResourcesPerNode prometheus.Gauge

	// HTTP 服务指标
	httpRequestsTotal     *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	httpActiveConnections prometheus.Gauge

	// 数据库连接池指标
	dbConnectionsActive  *prometheus.GaugeVec
	dbConnectionsIdle    *prometheus.GaugeVec
	dbConnectionsTotal   *prometheus.GaugeVec
	dbConnectionsPending *prometheus.GaugeVec
	dbConnectionWaitTime *prometheus.HistogramVec

	// 数据库连接池使用统计指标
	dbPoolUsageMean           *prometheus.GaugeVec
	dbPoolUsage75thPercentile *prometheus.GaugeVec
	dbPoolUsage95thPercentile *prometheus.GaugeVec
	dbPoolUsage99thPercentile *prometheus.GaugeVec
	dbPoolUsageMax            *prometheus.GaugeVec

	// 数据库连接池等待时间统计指标
	dbPoolWaitMean           *prometheus.GaugeVec
	dbPoolWait75thPercentile *prometheus.GaugeVec
	dbPoolWait95thPercentile *prometheus.GaugeVec
	dbPoolWait99thPercentile *prometheus.GaugeVec
	dbPoolWaitMax            *prometheus.GaugeVec

	// JVM 指标
	jvmMemoryUsed    *prometheus.GaugeVec
	jvmMemoryMax     *prometheus.GaugeVec
	jvmThreadsActive prometheus.Gauge
	jvmGCDuration    *prometheus.HistogramVec
}

// NewPuppetDBMetrics 创建PuppetDB指标实例
func NewPuppetDBMetrics(namespace string) *PuppetDBMetrics {
	pm := &PuppetDBMetrics{}

	// 命令处理指标
	pm.commandsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "commands_processed_total",
			Help:      "Total number of commands processed by PuppetDB",
		},
		[]string{"command", "version", "status"},
	)

	pm.commandsProcessingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "commands_processing_duration_seconds",
			Help:      "Command processing duration in seconds",
		},
		[]string{"command", "version"},
	)

	pm.commandQueueDepth = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "command_queue_depth",
			Help:      "Current depth of the command queue",
		},
	)

	// 存储层指标
	pm.storageDuplicatePct = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "storage_duplicate_percentage",
			Help:      "Percentage of catalogs that are duplicates",
		},
	)

	pm.storageGCDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "storage_gc_duration_seconds",
			Help:      "Storage garbage collection duration in seconds",
		},
	)

	pm.storageReplaceFactsDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "storage_replace_facts_duration_seconds",
			Help:      "Time taken to replace facts in seconds",
		},
	)

	pm.storageReplaceCatalogDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "storage_replace_catalog_duration_seconds",
			Help:      "Time taken to replace catalogs in seconds",
		},
	)

	// 人口统计指标
	pm.populationNodes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "population_nodes_total",
			Help:      "Total number of nodes in PuppetDB",
		},
	)

	pm.populationResources = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "population_resources_total",
			Help:      "Total number of resources in PuppetDB",
		},
	)

	pm.populationAvgResourcesPerNode = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "population_avg_resources_per_node",
			Help:      "Average number of resources per node",
		},
	)

	// HTTP 服务指标
	pm.httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"endpoint", "method", "status"},
	)

	pm.httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request duration in seconds",
		},
		[]string{"endpoint", "method"},
	)

	pm.httpActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "http_active_connections",
			Help:      "Number of active HTTP connections",
		},
	)

	// 数据库连接池指标
	pm.dbConnectionsActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_connections_active",
			Help:      "Number of active database connections",
		},
		[]string{"pool"},
	)

	pm.dbConnectionsIdle = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_connections_idle",
			Help:      "Number of idle database connections",
		},
		[]string{"pool"},
	)

	pm.dbConnectionsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_connections_total",
			Help:      "Total number of database connections",
		},
		[]string{"pool"},
	)

	pm.dbConnectionsPending = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_connections_pending",
			Help:      "Number of pending database connections",
		},
		[]string{"pool"},
	)

	pm.dbConnectionWaitTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "db_connection_wait_time_seconds",
			Help:      "Database connection wait time in seconds",
		},
		[]string{"pool"},
	)

	// 数据库连接池使用统计指标
	pm.dbPoolUsageMean = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_usage_mean",
			Help:      "Mean value of database pool usage",
		},
		[]string{"pool"},
	)

	pm.dbPoolUsage75thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_usage_75th_percentile",
			Help:      "75th percentile of database pool usage",
		},
		[]string{"pool"},
	)

	pm.dbPoolUsage95thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_usage_95th_percentile",
			Help:      "95th percentile of database pool usage",
		},
		[]string{"pool"},
	)

	pm.dbPoolUsage99thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_usage_99th_percentile",
			Help:      "99th percentile of database pool usage",
		},
		[]string{"pool"},
	)

	pm.dbPoolUsageMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_usage_max",
			Help:      "Maximum value of database pool usage",
		},
		[]string{"pool"},
	)

	// 数据库连接池等待时间统计指标
	pm.dbPoolWaitMean = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_wait_mean_seconds",
			Help:      "Mean wait time for database pool connections in seconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolWait75thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_wait_75th_percentile_seconds",
			Help:      "75th percentile of wait time for database pool connections in seconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolWait95thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_wait_95th_percentile_seconds",
			Help:      "95th percentile of wait time for database pool connections in seconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolWait99thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_wait_99th_percentile_seconds",
			Help:      "99th percentile of wait time for database pool connections in seconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolWaitMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_wait_max_seconds",
			Help:      "Maximum wait time for database pool connections in seconds",
		},
		[]string{"pool"},
	)

	// JVM 指标
	pm.jvmMemoryUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_memory_used_bytes",
			Help:      "JVM memory used in bytes",
		},
		[]string{"type"},
	)

	pm.jvmMemoryMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_memory_max_bytes",
			Help:      "JVM memory max in bytes",
		},
		[]string{"type"},
	)

	pm.jvmThreadsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_threads_active",
			Help:      "Number of active JVM threads",
		},
	)

	pm.jvmGCDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "jvm_gc_duration_seconds",
			Help:      "JVM garbage collection duration in seconds",
		},
		[]string{"gc"},
	)

	return pm
}

// Register 注册所有PuppetDB指标
func (pm *PuppetDBMetrics) Register() {
	// 命令处理指标
	prometheus.MustRegister(pm.commandsProcessed)
	prometheus.MustRegister(pm.commandsProcessingDuration)
	prometheus.MustRegister(pm.commandQueueDepth)

	// 存储层指标
	prometheus.MustRegister(pm.storageDuplicatePct)
	prometheus.MustRegister(pm.storageGCDuration)
	prometheus.MustRegister(pm.storageReplaceFactsDuration)
	prometheus.MustRegister(pm.storageReplaceCatalogDuration)

	// 人口统计指标
	prometheus.MustRegister(pm.populationNodes)
	prometheus.MustRegister(pm.populationResources)
	prometheus.MustRegister(pm.populationAvgResourcesPerNode)

	// HTTP 服务指标
	prometheus.MustRegister(pm.httpRequestsTotal)
	prometheus.MustRegister(pm.httpRequestDuration)
	prometheus.MustRegister(pm.httpActiveConnections)

	// 数据库连接池指标
	prometheus.MustRegister(pm.dbConnectionsActive)
	prometheus.MustRegister(pm.dbConnectionsIdle)
	prometheus.MustRegister(pm.dbConnectionsTotal)
	prometheus.MustRegister(pm.dbConnectionsPending)
	prometheus.MustRegister(pm.dbConnectionWaitTime)

	// 数据库连接池使用统计指标
	prometheus.MustRegister(pm.dbPoolUsageMean)
	prometheus.MustRegister(pm.dbPoolUsage75thPercentile)
	prometheus.MustRegister(pm.dbPoolUsage95thPercentile)
	prometheus.MustRegister(pm.dbPoolUsage99thPercentile)
	prometheus.MustRegister(pm.dbPoolUsageMax)

	// 数据库连接池等待时间统计指标
	prometheus.MustRegister(pm.dbPoolWaitMean)
	prometheus.MustRegister(pm.dbPoolWait75thPercentile)
	prometheus.MustRegister(pm.dbPoolWait95thPercentile)
	prometheus.MustRegister(pm.dbPoolWait99thPercentile)
	prometheus.MustRegister(pm.dbPoolWaitMax)

	// JVM 指标
	prometheus.MustRegister(pm.jvmMemoryUsed)
	prometheus.MustRegister(pm.jvmMemoryMax)
	prometheus.MustRegister(pm.jvmThreadsActive)
	prometheus.MustRegister(pm.jvmGCDuration)
}

// UpdateCommandMetrics 更新命令处理指标
func (pm *PuppetDBMetrics) UpdateCommandMetrics(command string, version string, status string, duration float64) {
	pm.commandsProcessed.WithLabelValues(command, version, status).Inc()
	if duration > 0 {
		pm.commandsProcessingDuration.WithLabelValues(command, version).Observe(duration)
	}
}

// UpdateStorageMetrics 更新存储层指标
func (pm *PuppetDBMetrics) UpdateStorageMetrics(duplicatePct float64, gcDuration float64, replaceFactsDuration float64, replaceCatalogDuration float64) {
	if duplicatePct >= 0 {
		pm.storageDuplicatePct.Set(duplicatePct)
	}
	if gcDuration >= 0 {
		pm.storageGCDuration.Set(gcDuration)
	}
	if replaceFactsDuration >= 0 {
		pm.storageReplaceFactsDuration.Set(replaceFactsDuration)
	}
	if replaceCatalogDuration >= 0 {
		pm.storageReplaceCatalogDuration.Set(replaceCatalogDuration)
	}
}

// UpdatePopulationMetrics 更新人口统计指标
func (pm *PuppetDBMetrics) UpdatePopulationMetrics(nodes float64, resources float64, avgResourcesPerNode float64) {
	if nodes >= 0 {
		pm.populationNodes.Set(nodes)
	}
	if resources >= 0 {
		pm.populationResources.Set(resources)
	}
	if avgResourcesPerNode >= 0 {
		pm.populationAvgResourcesPerNode.Set(avgResourcesPerNode)
	}
}

// UpdateHTTPMetrics 更新HTTP服务指标
func (pm *PuppetDBMetrics) UpdateHTTPMetrics(endpoint string, method string, status string, duration float64) {
	pm.httpRequestsTotal.WithLabelValues(endpoint, method, status).Inc()
	if duration > 0 {
		pm.httpRequestDuration.WithLabelValues(endpoint, method).Observe(duration)
	}
}

// UpdateDBMetrics 更新数据库连接池指标
func (pm *PuppetDBMetrics) UpdateDBMetrics(pool string, active float64, idle float64, total float64, waitTime float64) {
	if active >= 0 {
		pm.dbConnectionsActive.WithLabelValues(pool).Set(active)
	}
	if idle >= 0 {
		pm.dbConnectionsIdle.WithLabelValues(pool).Set(idle)
	}
	if total >= 0 {
		pm.dbConnectionsTotal.WithLabelValues(pool).Set(total)
	}
	if waitTime > 0 {
		pm.dbConnectionWaitTime.WithLabelValues(pool).Observe(waitTime)
	}
}

// UpdateDBPoolPendingConnections 更新数据库连接池待处理连接数
func (pm *PuppetDBMetrics) UpdateDBPoolPendingConnections(pool string, pending float64) {
	if pending >= 0 {
		pm.dbConnectionsPending.WithLabelValues(pool).Set(pending)
	}
}

// UpdateDBPoolUsageStats 更新数据库连接池使用统计
func (pm *PuppetDBMetrics) UpdateDBPoolUsageStats(pool string, mean float64, p75 float64, p95 float64, p99 float64, max float64) {
	if mean >= 0 {
		pm.dbPoolUsageMean.WithLabelValues(pool).Set(mean)
	}
	if p75 >= 0 {
		pm.dbPoolUsage75thPercentile.WithLabelValues(pool).Set(p75)
	}
	if p95 >= 0 {
		pm.dbPoolUsage95thPercentile.WithLabelValues(pool).Set(p95)
	}
	if p99 >= 0 {
		pm.dbPoolUsage99thPercentile.WithLabelValues(pool).Set(p99)
	}
	if max >= 0 {
		pm.dbPoolUsageMax.WithLabelValues(pool).Set(max)
	}
}

// UpdateDBPoolWaitStats 更新数据库连接池等待时间统计
func (pm *PuppetDBMetrics) UpdateDBPoolWaitStats(pool string, mean float64, p75 float64, p95 float64, p99 float64, max float64) {
	if mean >= 0 {
		pm.dbPoolWaitMean.WithLabelValues(pool).Set(mean)
	}
	if p75 >= 0 {
		pm.dbPoolWait75thPercentile.WithLabelValues(pool).Set(p75)
	}
	if p95 >= 0 {
		pm.dbPoolWait95thPercentile.WithLabelValues(pool).Set(p95)
	}
	if p99 >= 0 {
		pm.dbPoolWait99thPercentile.WithLabelValues(pool).Set(p99)
	}
	if max >= 0 {
		pm.dbPoolWaitMax.WithLabelValues(pool).Set(max)
	}
}

// UpdateCommandQueueDepth 更新命令队列深度
func (pm *PuppetDBMetrics) UpdateCommandQueueDepth(depth float64) {
	pm.commandQueueDepth.Set(depth)
}

// UpdateJVMMetrics 更新JVM指标
func (pm *PuppetDBMetrics) UpdateJVMMetrics(memoryType string, used float64, max float64, threads float64, gcType string, gcDuration float64) {
	if used >= 0 {
		pm.jvmMemoryUsed.WithLabelValues(memoryType).Set(used)
	}
	if max >= 0 {
		pm.jvmMemoryMax.WithLabelValues(memoryType).Set(max)
	}
	if threads >= 0 {
		pm.jvmThreadsActive.Set(threads)
	}
	if gcDuration > 0 {
		pm.jvmGCDuration.WithLabelValues(gcType).Observe(gcDuration)
	}
}
