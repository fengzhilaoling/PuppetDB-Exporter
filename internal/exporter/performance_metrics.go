package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// PerformanceMetrics 定义性能相关的指标
type PerformanceMetrics struct {
	scrapeDuration  *prometheus.HistogramVec
	scrapeErrors    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestsTotal   *prometheus.CounterVec
}

// NewPerformanceMetrics 创建性能指标实例
func NewPerformanceMetrics(namespace string) *PerformanceMetrics {
	pm := &PerformanceMetrics{}

	pm.scrapeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "puppetdb_exporter_scrape_duration_seconds",
			Help: "PuppetDB exporter scrape duration in seconds",
		},
		[]string{"endpoint"},
	)

	pm.scrapeErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "puppetdb_exporter_scrape_errors_total",
			Help: "Total number of scrape errors",
		},
		[]string{"endpoint", "error_type"},
	)

	pm.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "puppetdb_exporter_request_duration_seconds",
			Help: "PuppetDB API request duration in seconds",
		},
		[]string{"endpoint", "method"},
	)

	pm.requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "puppetdb_exporter_requests_total",
			Help: "Total number of PuppetDB API requests",
		},
		[]string{"endpoint", "status"},
	)

	return pm
}

// Register 注册所有性能指标
func (pm *PerformanceMetrics) Register() {
	prometheus.MustRegister(pm.scrapeDuration)
	prometheus.MustRegister(pm.scrapeErrors)
	prometheus.MustRegister(pm.requestDuration)
	prometheus.MustRegister(pm.requestsTotal)
}

// RecordScrapeDuration 记录抓取耗时
func (pm *PerformanceMetrics) RecordScrapeDuration(endpoint string, duration float64) {
	pm.scrapeDuration.With(prometheus.Labels{"endpoint": endpoint}).Observe(duration)
}

// RecordScrapeError 记录抓取错误
func (pm *PerformanceMetrics) RecordScrapeError(endpoint, errorType string) {
	pm.scrapeErrors.With(prometheus.Labels{"endpoint": endpoint, "error_type": errorType}).Inc()
}

// RecordRequestDuration 记录请求耗时
func (pm *PerformanceMetrics) RecordRequestDuration(endpoint, method string, duration float64) {
	pm.requestDuration.With(prometheus.Labels{"endpoint": endpoint, "method": method}).Observe(duration)
}

// RecordRequestTotal 记录请求总数
func (pm *PerformanceMetrics) RecordRequestTotal(endpoint, status string) {
	pm.requestsTotal.With(prometheus.Labels{"endpoint": endpoint, "status": status}).Inc()
}
