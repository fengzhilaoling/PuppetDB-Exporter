package exporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// NodeMetrics 定义节点相关的指标
type NodeMetrics struct {
	reportStatusCount *prometheus.GaugeVec
	hasReport         *prometheus.GaugeVec
	latestReportNoop  *prometheus.GaugeVec
	catalogTimestamp  *prometheus.GaugeVec
	factsTimestamp    *prometheus.GaugeVec
	report            *prometheus.GaugeVec
	reportAge         *prometheus.GaugeVec
	catalogAge        *prometheus.GaugeVec
	factsAge          *prometheus.GaugeVec
	reportMetrics     map[string]*prometheus.GaugeVec
}

// NewNodeMetrics 创建节点指标实例
func NewNodeMetrics(namespace string, categories map[string]struct{}) *NodeMetrics {
	nm := &NodeMetrics{
		reportMetrics: make(map[string]*prometheus.GaugeVec),
	}

	nm.reportStatusCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "node_report_status_count",
		Help:      "Number of nodes by latest report status (e.g., changed/failed/unchanged/unreported).",
	}, []string{"status"})

	nm.hasReport = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "node_has_report",
		Help:      "Whether node has latest report (1=yes, 0=no).",
	}, []string{"environment", "host"})

	nm.latestReportNoop = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "node_latest_report_noop",
		Help:      "Whether node's latest report is noop (1=yes, 0=no).",
	}, []string{"environment", "host"})

	// 新增：节点报告时间间隔指标
	nm.reportAge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "node_report_age_seconds",
		Help:      "Age of node's latest report in seconds.",
	}, []string{"environment", "host"})

	// 新增：编录时间间隔指标
	nm.catalogAge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "node_catalog_age_seconds",
		Help:      "Age of node's catalog in seconds.",
	}, []string{"environment", "host"})

	// 新增：事实数据时间间隔指标
	nm.factsAge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "node_facts_age_seconds",
		Help:      "Age of node's facts in seconds.",
	}, []string{"environment", "host"})

	nm.catalogTimestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "node_catalog_timestamp",
		Help:      "Node catalog timestamp (UNIX epoch).",
	}, []string{"environment", "host"})

	nm.factsTimestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "node_facts_timestamp",
		Help:      "Node facts timestamp (UNIX epoch).",
	}, []string{"environment", "host"})

	nm.report = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "puppet",
		Name:      "report",
		Help:      "Timestamp of node's latest report (UNIX epoch).",
	}, []string{"environment", "host", "deactivated"})

	// 为每个分类创建报告指标
	for category := range categories {
		metricName := fmt.Sprintf("report_%s", category)
		nm.reportMetrics[category] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "puppet",
			Name:      metricName,
			Help:      fmt.Sprintf("Number of %s by status (divided by name/environment/host)", category),
		}, []string{"name", "environment", "host"})
	}

	return nm
}

// Register 注册所有节点指标
func (nm *NodeMetrics) Register() {
	prometheus.MustRegister(nm.reportStatusCount)
	prometheus.MustRegister(nm.hasReport)
	prometheus.MustRegister(nm.latestReportNoop)
	prometheus.MustRegister(nm.catalogTimestamp)
	prometheus.MustRegister(nm.factsTimestamp)
	prometheus.MustRegister(nm.report)
	prometheus.MustRegister(nm.reportAge)
	prometheus.MustRegister(nm.catalogAge)
	prometheus.MustRegister(nm.factsAge)

	for _, metric := range nm.reportMetrics {
		prometheus.MustRegister(metric)
	}
}

// Reset 重置所有节点指标
func (nm *NodeMetrics) Reset() {
	nm.report.Reset()
	nm.reportStatusCount.Reset()

	for _, metric := range nm.reportMetrics {
		metric.Reset()
	}
}

// UpdateNodeMetrics 更新节点相关指标
func (nm *NodeMetrics) UpdateNodeMetrics(node NodeInfo, unreportedDuration time.Duration, now time.Time) {
	var deactivated string
	if node.Deactivated == "" {
		deactivated = "false"
	} else {
		deactivated = "true"
	}

	if node.ReportTimestamp == "" {
		return
	}

	latestReport, err := time.Parse(time.RFC3339, node.ReportTimestamp)
	if err != nil {
		log.Errorf("failed to parse report timestamp: %s", err)
		return
	}

	nm.report.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname, "deactivated": deactivated}).Set(float64(latestReport.Unix()))

	// 节点层面的指标
	if node.LatestReportHash != "" {
		nm.hasReport.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(1)
	} else {
		nm.hasReport.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(0)
	}

	if node.LatestReportNoop {
		nm.latestReportNoop.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(1)
	} else {
		nm.latestReportNoop.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(0)
	}

	// 计算并设置时间间隔指标
	reportAge := now.Sub(latestReport).Seconds()
	nm.reportAge.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(reportAge)

	// parse catalog_timestamp and facts_timestamp if present
	if node.CatalogTimestamp != "" {
		if t, err := time.Parse(time.RFC3339, node.CatalogTimestamp); err == nil {
			nm.catalogTimestamp.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(float64(t.Unix()))
			// 计算编录时间间隔
			catalogAge := now.Sub(t).Seconds()
			nm.catalogAge.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(catalogAge)
		}
	}
	if node.FactsTimestamp != "" {
		if t, err := time.Parse(time.RFC3339, node.FactsTimestamp); err == nil {
			nm.factsTimestamp.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(float64(t.Unix()))
			// 计算事实数据时间间隔
			factsAge := now.Sub(t).Seconds()
			nm.factsAge.With(prometheus.Labels{"environment": node.ReportEnvironment, "host": node.Certname}).Set(factsAge)
		}
	}
}

// UpdateReportMetrics 更新报告指标
func (nm *NodeMetrics) UpdateReportMetrics(node NodeInfo, reportMetrics []ReportMetric) {
	if node.LatestReportHash == "" {
		return
	}

	for _, reportMetric := range reportMetrics {
		if metric, ok := nm.reportMetrics[reportMetric.Category]; ok {
			displayName := strings.ReplaceAll(reportMetric.Name, "_", " ")
			metric.With(prometheus.Labels{"name": displayName, "environment": node.ReportEnvironment, "host": node.Certname}).Set(reportMetric.Value)
		}
	}
}

// UpdateStatusCount 更新状态计数
func (nm *NodeMetrics) UpdateStatusCount(statuses map[string]int) {
	for statusName, statusValue := range statuses {
		nm.reportStatusCount.With(prometheus.Labels{"status": statusName}).Set(float64(statusValue))
	}
}

// NodeInfo 节点信息结构体
type NodeInfo struct {
	Certname                string
	ReportEnvironment       string
	ReportTimestamp         string
	Deactivated             string
	LatestReportHash        string
	LatestReportNoop        bool
	LatestReportNoopPending bool
	LatestReportStatus      string
	CachedCatalogStatus     string
	CatalogTimestamp        string
	FactsTimestamp          string
}

// ReportMetric 报告指标结构体
type ReportMetric struct {
	Category string
	Name     string
	Value    float64
}
