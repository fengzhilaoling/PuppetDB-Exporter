package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// SystemMetrics 定义系统健康评分相关的指标
type SystemMetrics struct {
	healthScore   *prometheus.GaugeVec
	failureRate   *prometheus.GaugeVec
	degradedNodes *prometheus.GaugeVec
}

// NewSystemMetrics 创建系统指标实例
func NewSystemMetrics(namespace string) *SystemMetrics {
	sm := &SystemMetrics{}

	sm.healthScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "system_health_score",
			Help:      "PuppetDB system health score (0-100), calculated based on node status",
		}, []string{})

	sm.failureRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "node_failure_rate",
			Help:      "Node failure rate percentage",
		}, []string{})

	sm.degradedNodes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "degraded_nodes",
			Help:      "Number of degraded nodes (failed + unreported)",
		}, []string{})

	return sm
}

// Register 注册所有系统指标
func (sm *SystemMetrics) Register() {
	prometheus.MustRegister(sm.healthScore)
	prometheus.MustRegister(sm.failureRate)
	prometheus.MustRegister(sm.degradedNodes)
}

// UpdateSystemMetrics 更新系统健康评分指标
func (sm *SystemMetrics) UpdateSystemMetrics(statuses map[string]int) {
	var totalNodes, healthyNodes, warningNodes, criticalNodes int
	for _, status := range statuses {
		totalNodes++
		switch {
		case status == 0:
			healthyNodes++
		case status > 0 && status <= 2:
			warningNodes++
		default:
			criticalNodes++
		}
	}

	// 健康评分 = (健康节点数 / 总节点数) * 100
	if totalNodes > 0 {
		healthScore := (float64(healthyNodes) / float64(totalNodes)) * 100
		sm.healthScore.WithLabelValues().Set(healthScore)
	}

	// 节点失败率 = (失败节点数 / 总节点数) * 100
	if totalNodes > 0 {
		failureRate := (float64(statuses["failed"]) / float64(totalNodes)) * 100
		sm.failureRate.WithLabelValues().Set(failureRate)
	}

	// 降级节点数 = 失败节点数 + 未报告节点数
	degradedNodes := statuses["failed"] + statuses["unreported"]
	sm.degradedNodes.WithLabelValues().Set(float64(degradedNodes))
}
