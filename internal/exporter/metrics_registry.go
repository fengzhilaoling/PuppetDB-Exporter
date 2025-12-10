package exporter

import "github.com/prometheus/client_golang/prometheus"

// MetricsRegistry 指标注册表，统一管理所有指标
type MetricsRegistry struct {
	nodeMetrics        *NodeMetrics
	serviceMetrics     *ServiceMetrics
	systemMetrics      *SystemMetrics
	metricsV2          *MetricsV2
	performanceMetrics *PerformanceMetrics
	puppetDBMetrics    *PuppetDBMetrics
}

// NewMetricsRegistry 创建指标注册表
func NewMetricsRegistry(namespace string, categories map[string]struct{}) *MetricsRegistry {
	return &MetricsRegistry{
		nodeMetrics:        NewNodeMetrics(namespace, categories),
		serviceMetrics:     NewServiceMetrics(namespace),
		systemMetrics:      NewSystemMetrics(namespace),
		metricsV2:          NewMetricsV2(namespace),
		performanceMetrics: NewPerformanceMetrics(namespace),
		puppetDBMetrics:    NewPuppetDBMetrics(namespace),
	}
}

// RegisterAll 注册所有指标
func (mr *MetricsRegistry) RegisterAll() {
	mr.nodeMetrics.Register()
	mr.serviceMetrics.Register()
	mr.systemMetrics.Register()
	mr.metricsV2.Register()
	mr.performanceMetrics.Register()
	mr.puppetDBMetrics.Register()
}

// GetNodeMetrics 获取节点指标
func (mr *MetricsRegistry) GetNodeMetrics() *NodeMetrics {
	return mr.nodeMetrics
}

// GetServiceMetrics 获取服务指标
func (mr *MetricsRegistry) GetServiceMetrics() *ServiceMetrics {
	return mr.serviceMetrics
}

// GetSystemMetrics 获取系统指标
func (mr *MetricsRegistry) GetSystemMetrics() *SystemMetrics {
	return mr.systemMetrics
}

// GetMetricsV2 获取 metrics v2 指标
func (mr *MetricsRegistry) GetMetricsV2() *MetricsV2 {
	return mr.metricsV2
}

// GetPerformanceMetrics 获取性能指标
func (mr *MetricsRegistry) GetPerformanceMetrics() *PerformanceMetrics {
	return mr.performanceMetrics
}

// GetPuppetDBMetrics 获取PuppetDB核心指标
func (mr *MetricsRegistry) GetPuppetDBMetrics() *PuppetDBMetrics {
	return mr.puppetDBMetrics
}

// Describe 输出所有指标描述
func (mr *MetricsRegistry) Describe(ch chan<- *prometheus.Desc) {
	// 这里可以遍历所有指标并调用它们的 Describe 方法
	// 为了简化，暂时留空，由主 exporter 处理
}

// Collect 收集所有指标
func (mr *MetricsRegistry) Collect(ch chan<- prometheus.Metric) {
	// 这里可以遍历所有指标并调用它们的 Collect 方法
	// 为了简化，暂时留空，由主 exporter 处理
}
