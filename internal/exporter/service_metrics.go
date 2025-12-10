package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// ServiceMetrics 定义服务相关的指标
type ServiceMetrics struct {
	up         *prometheus.GaugeVec
	info       *prometheus.GaugeVec
	queueDepth *prometheus.GaugeVec
}

// NewServiceMetrics 创建服务指标实例
func NewServiceMetrics(namespace string) *ServiceMetrics {
	sm := &ServiceMetrics{}

	sm.up = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "service_up",
		Help:      "Whether service is running (1=running, 0=not running).",
	}, []string{"service", "version"})

	sm.info = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "service_info",
		Help:      "Service information metric (for recording service version and status, always 1).",
	}, []string{"service", "version", "state"})

	sm.queueDepth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "service_queue_depth",
		Help:      "Service queue depth (indicating number of unprocessed tasks).",
	}, []string{"service"})

	return sm
}

// Register 注册所有服务指标
func (sm *ServiceMetrics) Register() {
	prometheus.MustRegister(sm.up)
	prometheus.MustRegister(sm.info)
	prometheus.MustRegister(sm.queueDepth)
}

// Reset 重置所有服务指标
func (sm *ServiceMetrics) Reset() {
	sm.up.Reset()
	sm.queueDepth.Reset()
	sm.info.Reset()
}

// UpdateServiceMetrics 更新服务指标
func (sm *ServiceMetrics) UpdateServiceMetrics(services []ServiceInfo) {
	for _, svc := range services {
		svcName := svc.Name
		if svcName == "" {
			svcName = "puppetdb"
		}

		// 服务状态
		if svc.Up {
			sm.up.With(prometheus.Labels{"service": svcName, "version": svc.Version}).Set(1)
		} else {
			sm.up.With(prometheus.Labels{"service": svcName, "version": svc.Version}).Set(0)
		}

		// 服务信息
		sm.info.With(prometheus.Labels{"service": svcName, "version": svc.Version, "state": svc.State}).Set(1)

		// 队列深度
		sm.queueDepth.With(prometheus.Labels{"service": svcName}).Set(float64(svc.QueueDepth))
	}
}

// ServiceInfo 服务信息结构
// 用于从PuppetDB API获取服务状态信息
type ServiceInfo struct {
	Name       string
	Version    string
	State      string
	Up         bool
	QueueDepth int
}
