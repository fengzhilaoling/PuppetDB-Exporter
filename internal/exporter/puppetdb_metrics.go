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

	// HTTP 端点详细指标
	httpRequestRate            *prometheus.GaugeVec
	httpServiceTimePercentiles *prometheus.GaugeVec
	httpServiceTimeStats       *prometheus.GaugeVec

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

	// 数据库连接池配置指标
	dbPoolMaxConnections *prometheus.GaugeVec
	dbPoolMinConnections *prometheus.GaugeVec

	// 数据库连接池连接创建统计指标
	dbPoolConnectionCreationMean           *prometheus.GaugeVec
	dbPoolConnectionCreation75thPercentile *prometheus.GaugeVec
	dbPoolConnectionCreation95thPercentile *prometheus.GaugeVec
	dbPoolConnectionCreation99thPercentile *prometheus.GaugeVec
	dbPoolConnectionCreationMax            *prometheus.GaugeVec
	dbPoolConnectionCreationCount          *prometheus.GaugeVec

	// 数据库连接池连接超时率指标
	dbPoolConnectionTimeoutRateOneMinute     *prometheus.GaugeVec
	dbPoolConnectionTimeoutRateFiveMinute    *prometheus.GaugeVec
	dbPoolConnectionTimeoutRateFifteenMinute *prometheus.GaugeVec
	dbPoolConnectionTimeoutRateMean          *prometheus.GaugeVec
	dbPoolConnectionTimeoutRateCount         *prometheus.GaugeVec

	// JVM 指标
	jvmMemoryUsed    *prometheus.GaugeVec
	jvmMemoryMax     *prometheus.GaugeVec
	jvmThreadsActive prometheus.Gauge
	jvmGCDuration    *prometheus.HistogramVec

	// JVM 内存池指标
	jvmMemoryPoolUsed      *prometheus.GaugeVec
	jvmMemoryPoolCommitted *prometheus.GaugeVec
	jvmMemoryPoolMax       *prometheus.GaugeVec
	jvmMemoryPoolPeakUsed  *prometheus.GaugeVec

	// JVM 垃圾收集器指标
	jvmGCCollectionCount    *prometheus.CounterVec
	jvmGCCollectionTime     *prometheus.CounterVec
	jvmGCLastGcInfoDuration *prometheus.GaugeVec

	// JVM 运行时系统指标
	jvmClassLoadingLoadedClassCount      prometheus.Gauge
	jvmClassLoadingUnloadedClassCount    prometheus.Counter
	jvmClassLoadingTotalLoadedClassCount prometheus.Counter

	jvmCompilationTotalTime prometheus.Counter

	jvmOperatingSystemOpenFileDescriptors    prometheus.Gauge
	jvmOperatingSystemCommittedVirtualMemory prometheus.Gauge
	jvmOperatingSystemFreePhysicalMemory     prometheus.Gauge
	jvmOperatingSystemSystemLoadAverage      prometheus.Gauge
	jvmOperatingSystemProcessCpuLoad         prometheus.Gauge
	jvmOperatingSystemFreeSwapSpace          prometheus.Gauge
	jvmOperatingSystemTotalPhysicalMemory    prometheus.Gauge
	jvmOperatingSystemTotalSwapSpace         prometheus.Gauge
	jvmOperatingSystemProcessCpuTime         prometheus.Counter
	jvmOperatingSystemMaxFileDescriptors     prometheus.Gauge
	jvmOperatingSystemSystemCpuLoad          prometheus.Gauge
	jvmOperatingSystemAvailableProcessors    prometheus.Gauge
	jvmOperatingSystemCpuLoad                prometheus.Gauge
	jvmOperatingSystemFreeMemory             prometheus.Gauge

	jvmRuntimeUptime    prometheus.Counter
	jvmRuntimeStartTime prometheus.Gauge

	jvmThreadingTotalStartedThreads          prometheus.Counter
	jvmThreadingPeakThreadCount              prometheus.Gauge
	jvmThreadingDaemonThreadCount            prometheus.Gauge
	jvmThreadingCurrentThreadAllocatedBytes  prometheus.Gauge
	jvmThreadingThreadAllocatedMemoryEnabled prometheus.Gauge
	jvmThreadingThreadCpuTimeEnabled         prometheus.Gauge
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

	// HTTP 端点详细指标
	pm.httpRequestRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "http_request_rate",
			Help:      "HTTP request rate in requests per second",
		},
		[]string{"endpoint", "rate_type"},
	)

	pm.httpServiceTimePercentiles = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "http_service_time_percentile_seconds",
			Help:      "HTTP service time percentiles in seconds",
		},
		[]string{"endpoint", "percentile"},
	)

	pm.httpServiceTimeStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "http_service_time_stats_seconds",
			Help:      "HTTP service time statistics in seconds",
		},
		[]string{"endpoint", "stat"},
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

	// 数据库连接池配置指标
	pm.dbPoolMaxConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_max_connections",
			Help:      "Maximum number of connections in the database pool",
		},
		[]string{"pool"},
	)

	pm.dbPoolMinConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_min_connections",
			Help:      "Minimum number of connections in the database pool",
		},
		[]string{"pool"},
	)

	// 数据库连接池连接创建统计指标
	pm.dbPoolConnectionCreationMean = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_creation_mean_ms",
			Help:      "Mean time for database connection creation in milliseconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionCreation75thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_creation_75th_percentile_ms",
			Help:      "75th percentile of database connection creation time in milliseconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionCreation95thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_creation_95th_percentile_ms",
			Help:      "95th percentile of database connection creation time in milliseconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionCreation99thPercentile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_creation_99th_percentile_ms",
			Help:      "99th percentile of database connection creation time in milliseconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionCreationMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_creation_max_ms",
			Help:      "Maximum time for database connection creation in milliseconds",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionCreationCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_creation_count",
			Help:      "Total count of database connection creation operations",
		},
		[]string{"pool"},
	)

	// 数据库连接池连接超时率指标
	pm.dbPoolConnectionTimeoutRateOneMinute = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_timeout_rate_one_minute",
			Help:      "One minute rate of database connection timeouts per second",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionTimeoutRateFiveMinute = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_timeout_rate_five_minute",
			Help:      "Five minute rate of database connection timeouts per second",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionTimeoutRateFifteenMinute = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_timeout_rate_fifteen_minute",
			Help:      "Fifteen minute rate of database connection timeouts per second",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionTimeoutRateMean = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_timeout_rate_mean",
			Help:      "Mean rate of database connection timeouts per second",
		},
		[]string{"pool"},
	)

	pm.dbPoolConnectionTimeoutRateCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_pool_connection_timeout_rate_count",
			Help:      "Total count of database connection timeouts",
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

	// JVM 内存池指标
	pm.jvmMemoryPoolUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_memory_pool_used_bytes",
			Help:      "JVM memory pool used in bytes",
		},
		[]string{"pool"},
	)

	pm.jvmMemoryPoolCommitted = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_memory_pool_committed_bytes",
			Help:      "JVM memory pool committed in bytes",
		},
		[]string{"pool"},
	)

	pm.jvmMemoryPoolMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_memory_pool_max_bytes",
			Help:      "JVM memory pool max in bytes",
		},
		[]string{"pool"},
	)

	pm.jvmMemoryPoolPeakUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_memory_pool_peak_used_bytes",
			Help:      "JVM memory pool peak used in bytes",
		},
		[]string{"pool"},
	)

	// JVM 垃圾收集器指标
	pm.jvmGCCollectionCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "jvm_gc_collection_count",
			Help:      "Total number of garbage collection events",
		},
		[]string{"gc"},
	)

	pm.jvmGCCollectionTime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "jvm_gc_collection_time_seconds",
			Help:      "Total time spent in garbage collection in seconds",
		},
		[]string{"gc"},
	)

	pm.jvmGCLastGcInfoDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_gc_last_gc_info_duration_seconds",
			Help:      "Duration of the last garbage collection in seconds",
		},
		[]string{"gc"},
	)

	// JVM 运行时系统指标
	pm.jvmClassLoadingLoadedClassCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_class_loading_loaded_class_count",
			Help:      "Number of classes currently loaded in the JVM",
		},
	)

	pm.jvmClassLoadingUnloadedClassCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "jvm_class_loading_unloaded_class_count",
			Help:      "Total number of classes unloaded by the JVM",
		},
	)

	pm.jvmClassLoadingTotalLoadedClassCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "jvm_class_loading_total_loaded_class_count",
			Help:      "Total number of classes loaded by the JVM",
		},
	)

	pm.jvmCompilationTotalTime = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "jvm_compilation_total_time_seconds",
			Help:      "Total time spent in JVM compilation in seconds",
		},
	)

	pm.jvmOperatingSystemOpenFileDescriptors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_open_file_descriptors",
			Help:      "Number of open file descriptors",
		},
	)

	pm.jvmOperatingSystemCommittedVirtualMemory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_committed_virtual_memory_bytes",
			Help:      "Amount of committed virtual memory in bytes",
		},
	)

	pm.jvmOperatingSystemFreePhysicalMemory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_free_physical_memory_bytes",
			Help:      "Amount of free physical memory in bytes",
		},
	)

	pm.jvmOperatingSystemSystemLoadAverage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_system_load_average",
			Help:      "System load average",
		},
	)

	pm.jvmOperatingSystemProcessCpuLoad = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_process_cpu_load",
			Help:      "Process CPU load",
		},
	)

	pm.jvmOperatingSystemFreeSwapSpace = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_free_swap_space_bytes",
			Help:      "Amount of free swap space in bytes",
		},
	)

	pm.jvmOperatingSystemTotalPhysicalMemory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_total_physical_memory_bytes",
			Help:      "Total amount of physical memory in bytes",
		},
	)

	pm.jvmOperatingSystemTotalSwapSpace = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_total_swap_space_bytes",
			Help:      "Total amount of swap space in bytes",
		},
	)

	pm.jvmOperatingSystemProcessCpuTime = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_process_cpu_time_seconds",
			Help:      "Total CPU time used by the process in seconds",
		},
	)

	pm.jvmOperatingSystemMaxFileDescriptors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_max_file_descriptors",
			Help:      "Maximum number of file descriptors",
		},
	)

	pm.jvmOperatingSystemSystemCpuLoad = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_system_cpu_load",
			Help:      "System CPU load",
		},
	)

	pm.jvmOperatingSystemAvailableProcessors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_available_processors",
			Help:      "Number of available processors",
		},
	)

	pm.jvmOperatingSystemCpuLoad = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_cpu_load",
			Help:      "CPU load",
		},
	)

	pm.jvmOperatingSystemFreeMemory = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_operating_system_free_memory_bytes",
			Help:      "Amount of free memory in bytes",
		},
	)

	pm.jvmRuntimeUptime = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "jvm_runtime_uptime_seconds",
			Help:      "JVM runtime uptime in seconds",
		},
	)

	pm.jvmRuntimeStartTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_runtime_start_time_seconds",
			Help:      "JVM runtime start time in seconds since epoch",
		},
	)

	pm.jvmThreadingTotalStartedThreads = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "jvm_threading_total_started_threads",
			Help:      "Total number of threads started",
		},
	)

	pm.jvmThreadingPeakThreadCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_threading_peak_thread_count",
			Help:      "Peak number of threads",
		},
	)

	pm.jvmThreadingDaemonThreadCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_threading_daemon_thread_count",
			Help:      "Number of daemon threads",
		},
	)

	pm.jvmThreadingCurrentThreadAllocatedBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_threading_current_thread_allocated_bytes",
			Help:      "Bytes allocated for the current thread",
		},
	)

	pm.jvmThreadingThreadAllocatedMemoryEnabled = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_threading_thread_allocated_memory_enabled",
			Help:      "Whether thread allocated memory tracking is enabled",
		},
	)

	pm.jvmThreadingThreadCpuTimeEnabled = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "jvm_threading_thread_cpu_time_enabled",
			Help:      "Whether thread CPU time tracking is enabled",
		},
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
	prometheus.MustRegister(pm.httpRequestRate)
	prometheus.MustRegister(pm.httpServiceTimePercentiles)
	prometheus.MustRegister(pm.httpServiceTimeStats)

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

	// 数据库连接池配置指标
	prometheus.MustRegister(pm.dbPoolMaxConnections)
	prometheus.MustRegister(pm.dbPoolMinConnections)

	// 数据库连接池连接创建统计指标
	prometheus.MustRegister(pm.dbPoolConnectionCreationMean)
	prometheus.MustRegister(pm.dbPoolConnectionCreation75thPercentile)
	prometheus.MustRegister(pm.dbPoolConnectionCreation95thPercentile)
	prometheus.MustRegister(pm.dbPoolConnectionCreation99thPercentile)
	prometheus.MustRegister(pm.dbPoolConnectionCreationMax)
	prometheus.MustRegister(pm.dbPoolConnectionCreationCount)

	// 数据库连接池连接超时率指标
	prometheus.MustRegister(pm.dbPoolConnectionTimeoutRateOneMinute)
	prometheus.MustRegister(pm.dbPoolConnectionTimeoutRateFiveMinute)
	prometheus.MustRegister(pm.dbPoolConnectionTimeoutRateFifteenMinute)
	prometheus.MustRegister(pm.dbPoolConnectionTimeoutRateMean)
	prometheus.MustRegister(pm.dbPoolConnectionTimeoutRateCount)

	// JVM 指标
	prometheus.MustRegister(pm.jvmMemoryUsed)
	prometheus.MustRegister(pm.jvmMemoryMax)
	prometheus.MustRegister(pm.jvmThreadsActive)
	prometheus.MustRegister(pm.jvmGCDuration)

	// JVM 内存池指标
	prometheus.MustRegister(pm.jvmMemoryPoolUsed)
	prometheus.MustRegister(pm.jvmMemoryPoolCommitted)
	prometheus.MustRegister(pm.jvmMemoryPoolMax)
	prometheus.MustRegister(pm.jvmMemoryPoolPeakUsed)

	// JVM 垃圾收集器指标
	prometheus.MustRegister(pm.jvmGCCollectionCount)
	prometheus.MustRegister(pm.jvmGCCollectionTime)
	prometheus.MustRegister(pm.jvmGCLastGcInfoDuration)

	// JVM 运行时系统指标
	prometheus.MustRegister(pm.jvmClassLoadingLoadedClassCount)
	prometheus.MustRegister(pm.jvmClassLoadingUnloadedClassCount)
	prometheus.MustRegister(pm.jvmClassLoadingTotalLoadedClassCount)

	prometheus.MustRegister(pm.jvmCompilationTotalTime)

	prometheus.MustRegister(pm.jvmOperatingSystemOpenFileDescriptors)
	prometheus.MustRegister(pm.jvmOperatingSystemCommittedVirtualMemory)
	prometheus.MustRegister(pm.jvmOperatingSystemFreePhysicalMemory)
	prometheus.MustRegister(pm.jvmOperatingSystemSystemLoadAverage)
	prometheus.MustRegister(pm.jvmOperatingSystemProcessCpuLoad)
	prometheus.MustRegister(pm.jvmOperatingSystemFreeSwapSpace)
	prometheus.MustRegister(pm.jvmOperatingSystemTotalPhysicalMemory)
	prometheus.MustRegister(pm.jvmOperatingSystemTotalSwapSpace)
	prometheus.MustRegister(pm.jvmOperatingSystemProcessCpuTime)
	prometheus.MustRegister(pm.jvmOperatingSystemMaxFileDescriptors)
	prometheus.MustRegister(pm.jvmOperatingSystemSystemCpuLoad)
	prometheus.MustRegister(pm.jvmOperatingSystemAvailableProcessors)
	prometheus.MustRegister(pm.jvmOperatingSystemCpuLoad)
	prometheus.MustRegister(pm.jvmOperatingSystemFreeMemory)

	prometheus.MustRegister(pm.jvmRuntimeUptime)
	prometheus.MustRegister(pm.jvmRuntimeStartTime)

	prometheus.MustRegister(pm.jvmThreadingTotalStartedThreads)
	prometheus.MustRegister(pm.jvmThreadingPeakThreadCount)
	prometheus.MustRegister(pm.jvmThreadingDaemonThreadCount)
	prometheus.MustRegister(pm.jvmThreadingCurrentThreadAllocatedBytes)
	prometheus.MustRegister(pm.jvmThreadingThreadAllocatedMemoryEnabled)
	prometheus.MustRegister(pm.jvmThreadingThreadCpuTimeEnabled)
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

// UpdateDBPoolConfig 更新数据库连接池配置指标
func (pm *PuppetDBMetrics) UpdateDBPoolConfig(pool string, maxConnections float64, minConnections float64) {
	if maxConnections >= 0 {
		pm.dbPoolMaxConnections.WithLabelValues(pool).Set(maxConnections)
	}
	if minConnections >= 0 {
		pm.dbPoolMinConnections.WithLabelValues(pool).Set(minConnections)
	}
}

// UpdateDBPoolConnectionCreationStats 更新数据库连接池连接创建统计
func (pm *PuppetDBMetrics) UpdateDBPoolConnectionCreationStats(pool string, mean float64, p75 float64, p95 float64, p99 float64, max float64, count float64) {
	if mean >= 0 {
		pm.dbPoolConnectionCreationMean.WithLabelValues(pool).Set(mean)
	}
	if p75 >= 0 {
		pm.dbPoolConnectionCreation75thPercentile.WithLabelValues(pool).Set(p75)
	}
	if p95 >= 0 {
		pm.dbPoolConnectionCreation95thPercentile.WithLabelValues(pool).Set(p95)
	}
	if p99 >= 0 {
		pm.dbPoolConnectionCreation99thPercentile.WithLabelValues(pool).Set(p99)
	}
	if max >= 0 {
		pm.dbPoolConnectionCreationMax.WithLabelValues(pool).Set(max)
	}
	if count >= 0 {
		pm.dbPoolConnectionCreationCount.WithLabelValues(pool).Set(count)
	}
}

// UpdateDBPoolConnectionTimeoutRateStats 更新数据库连接池连接超时率统计
func (pm *PuppetDBMetrics) UpdateDBPoolConnectionTimeoutRateStats(pool string, oneMinute float64, fiveMinute float64, fifteenMinute float64, mean float64, count float64) {
	if oneMinute >= 0 {
		pm.dbPoolConnectionTimeoutRateOneMinute.WithLabelValues(pool).Set(oneMinute)
	}
	if fiveMinute >= 0 {
		pm.dbPoolConnectionTimeoutRateFiveMinute.WithLabelValues(pool).Set(fiveMinute)
	}
	if fifteenMinute >= 0 {
		pm.dbPoolConnectionTimeoutRateFifteenMinute.WithLabelValues(pool).Set(fifteenMinute)
	}
	if mean >= 0 {
		pm.dbPoolConnectionTimeoutRateMean.WithLabelValues(pool).Set(mean)
	}
	if count >= 0 {
		pm.dbPoolConnectionTimeoutRateCount.WithLabelValues(pool).Set(count)
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

// UpdateJVMHeapMemoryPoolMetrics 更新JVM堆内存池指标
func (pm *PuppetDBMetrics) UpdateJVMHeapMemoryPoolMetrics(pool string, used float64, committed float64, max float64, peakUsed float64) {
	if used >= 0 {
		pm.jvmMemoryPoolUsed.WithLabelValues(pool).Set(used)
	}
	if committed >= 0 {
		pm.jvmMemoryPoolCommitted.WithLabelValues(pool).Set(committed)
	}
	if max >= 0 {
		pm.jvmMemoryPoolMax.WithLabelValues(pool).Set(max)
	}
	if peakUsed >= 0 {
		pm.jvmMemoryPoolPeakUsed.WithLabelValues(pool).Set(peakUsed)
	}
}

// UpdateJVMGarbageCollectorMetrics 更新JVM垃圾收集器指标
func (pm *PuppetDBMetrics) UpdateJVMGarbageCollectorMetrics(gc string, collectionCount float64, collectionTime float64, lastGcInfoDuration float64) {
	if collectionCount >= 0 {
		pm.jvmGCCollectionCount.WithLabelValues(gc).Set(collectionCount)
	}
	if collectionTime >= 0 {
		pm.jvmGCCollectionTime.WithLabelValues(gc).Set(collectionTime)
	}
	if lastGcInfoDuration >= 0 {
		pm.jvmGCLastGcInfoDuration.WithLabelValues(gc).Set(lastGcInfoDuration)
	}
}

// UpdateJVMClassLoadingMetrics 更新JVM类加载指标
func (pm *PuppetDBMetrics) UpdateJVMClassLoadingMetrics(loadedClassCount float64, unloadedClassCount float64, totalLoadedClassCount float64) {
	if loadedClassCount >= 0 {
		pm.jvmClassLoadingLoadedClassCount.Set(loadedClassCount)
	}
	if unloadedClassCount >= 0 {
		pm.jvmClassLoadingUnloadedClassCount.Add(unloadedClassCount)
	}
	if totalLoadedClassCount >= 0 {
		pm.jvmClassLoadingTotalLoadedClassCount.Add(totalLoadedClassCount)
	}
}

// UpdateJVMCompilationMetrics 更新JVM编译指标
func (pm *PuppetDBMetrics) UpdateJVMCompilationMetrics(totalTime float64) {
	if totalTime >= 0 {
		pm.jvmCompilationTotalTime.Add(totalTime)
	}
}

// UpdateJVMSysMetrics 更新JVM系统指标
func (pm *PuppetDBMetrics) UpdateJVMSysMetrics(openFileDescriptors float64, committedVirtualMemory float64, freePhysicalMemory float64, systemLoadAverage float64, processCpuLoad float64, freeSwapSpace float64, totalPhysicalMemory float64, totalSwapSpace float64, processCpuTime float64, maxFileDescriptors float64, systemCpuLoad float64, availableProcessors float64, cpuLoad float64, freeMemory float64) {
	if openFileDescriptors >= 0 {
		pm.jvmOperatingSystemOpenFileDescriptors.Set(openFileDescriptors)
	}
	if committedVirtualMemory >= 0 {
		pm.jvmOperatingSystemCommittedVirtualMemory.Set(committedVirtualMemory)
	}
	if freePhysicalMemory >= 0 {
		pm.jvmOperatingSystemFreePhysicalMemory.Set(freePhysicalMemory)
	}
	if systemLoadAverage >= 0 {
		pm.jvmOperatingSystemSystemLoadAverage.Set(systemLoadAverage)
	}
	if processCpuLoad >= 0 {
		pm.jvmOperatingSystemProcessCpuLoad.Set(processCpuLoad)
	}
	if freeSwapSpace >= 0 {
		pm.jvmOperatingSystemFreeSwapSpace.Set(freeSwapSpace)
	}
	if totalPhysicalMemory >= 0 {
		pm.jvmOperatingSystemTotalPhysicalMemory.Set(totalPhysicalMemory)
	}
	if totalSwapSpace >= 0 {
		pm.jvmOperatingSystemTotalSwapSpace.Set(totalSwapSpace)
	}
	if processCpuTime >= 0 {
		pm.jvmOperatingSystemProcessCpuTime.Add(processCpuTime)
	}
	if maxFileDescriptors >= 0 {
		pm.jvmOperatingSystemMaxFileDescriptors.Set(maxFileDescriptors)
	}
	if systemCpuLoad >= 0 {
		pm.jvmOperatingSystemSystemCpuLoad.Set(systemCpuLoad)
	}
	if availableProcessors >= 0 {
		pm.jvmOperatingSystemAvailableProcessors.Set(availableProcessors)
	}
	if cpuLoad >= 0 {
		pm.jvmOperatingSystemCpuLoad.Set(cpuLoad)
	}
	if freeMemory >= 0 {
		pm.jvmOperatingSystemFreeMemory.Set(freeMemory)
	}
}

// UpdateJVMTimingMetrics 更新JVM计时指标
func (pm *PuppetDBMetrics) UpdateJVMTimingMetrics(uptime float64, startTime float64) {
	if uptime >= 0 {
		pm.jvmRuntimeUptime.Add(uptime)
	}
	if startTime >= 0 {
		pm.jvmRuntimeStartTime.Set(startTime)
	}
}

// UpdateJVMThreadingMetrics 更新JVM线程指标
func (pm *PuppetDBMetrics) UpdateJVMThreadingMetrics(totalStartedThreads float64, peakThreadCount float64, daemonThreadCount float64, currentThreadAllocatedBytes float64, threadAllocatedMemoryEnabled float64, threadCpuTimeEnabled float64) {
	if totalStartedThreads >= 0 {
		pm.jvmThreadingTotalStartedThreads.Add(totalStartedThreads)
	}
	if peakThreadCount >= 0 {
		pm.jvmThreadingPeakThreadCount.Set(peakThreadCount)
	}
	if daemonThreadCount >= 0 {
		pm.jvmThreadingDaemonThreadCount.Set(daemonThreadCount)
	}
	if currentThreadAllocatedBytes >= 0 {
		pm.jvmThreadingCurrentThreadAllocatedBytes.Set(currentThreadAllocatedBytes)
	}
	if threadAllocatedMemoryEnabled >= 0 {
		pm.jvmThreadingThreadAllocatedMemoryEnabled.Set(threadAllocatedMemoryEnabled)
	}
	if threadCpuTimeEnabled >= 0 {
		pm.jvmThreadingThreadCpuTimeEnabled.Set(threadCpuTimeEnabled)
	}
}

// UpdateHTTPDetailedMetrics 更新HTTP端点详细指标
func (pm *PuppetDBMetrics) UpdateHTTPDetailedMetrics(endpoint string, metrics map[string]float64) {
	// 更新请求速率指标
	if oneMinuteRate, ok := metrics["requests_200_one_minute_rate"]; ok && oneMinuteRate >= 0 {
		pm.httpRequestRate.WithLabelValues(endpoint, "one_minute").Set(oneMinuteRate)
	}
	if fiveMinuteRate, ok := metrics["requests_200_five_minute_rate"]; ok && fiveMinuteRate >= 0 {
		pm.httpRequestRate.WithLabelValues(endpoint, "five_minute").Set(fiveMinuteRate)
	}
	if fifteenMinuteRate, ok := metrics["requests_200_fifteen_minute_rate"]; ok && fifteenMinuteRate >= 0 {
		pm.httpRequestRate.WithLabelValues(endpoint, "fifteen_minute").Set(fifteenMinuteRate)
	}
	if meanRate, ok := metrics["requests_200_mean_rate"]; ok && meanRate >= 0 {
		pm.httpRequestRate.WithLabelValues(endpoint, "mean").Set(meanRate)
	}

	// 更新服务时间分位数指标
	if p50, ok := metrics["service_time_50th_percentile"]; ok && p50 >= 0 {
		pm.httpServiceTimePercentiles.WithLabelValues(endpoint, "50").Set(p50)
	}
	if p75, ok := metrics["service_time_75th_percentile"]; ok && p75 >= 0 {
		pm.httpServiceTimePercentiles.WithLabelValues(endpoint, "75").Set(p75)
	}
	if p95, ok := metrics["service_time_95th_percentile"]; ok && p95 >= 0 {
		pm.httpServiceTimePercentiles.WithLabelValues(endpoint, "95").Set(p95)
	}
	if p98, ok := metrics["service_time_98th_percentile"]; ok && p98 >= 0 {
		pm.httpServiceTimePercentiles.WithLabelValues(endpoint, "98").Set(p98)
	}
	if p99, ok := metrics["service_time_99th_percentile"]; ok && p99 >= 0 {
		pm.httpServiceTimePercentiles.WithLabelValues(endpoint, "99").Set(p99)
	}
	if p999, ok := metrics["service_time_999th_percentile"]; ok && p999 >= 0 {
		pm.httpServiceTimePercentiles.WithLabelValues(endpoint, "999").Set(p999)
	}

	// 更新服务时间统计指标
	if mean, ok := metrics["service_time_mean"]; ok && mean >= 0 {
		pm.httpServiceTimeStats.WithLabelValues(endpoint, "mean").Set(mean)
	}
	if stddev, ok := metrics["service_time_stddev"]; ok && stddev >= 0 {
		pm.httpServiceTimeStats.WithLabelValues(endpoint, "stddev").Set(stddev)
	}
	if min, ok := metrics["service_time_min"]; ok && min >= 0 {
		pm.httpServiceTimeStats.WithLabelValues(endpoint, "min").Set(min)
	}
	if max, ok := metrics["service_time_max"]; ok && max >= 0 {
		pm.httpServiceTimeStats.WithLabelValues(endpoint, "max").Set(max)
	}
}
