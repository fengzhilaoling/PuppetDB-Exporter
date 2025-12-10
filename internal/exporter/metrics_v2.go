package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsV2 定义 /metrics/v2 相关的指标
type MetricsV2 struct {
	status    *prometheus.GaugeVec
	timestamp *prometheus.GaugeVec
	info      *prometheus.GaugeVec
	config    *prometheus.GaugeVec
}

// NewMetricsV2 创建 metrics v2 指标实例
func NewMetricsV2(namespace string) *MetricsV2 {
	m2 := &MetricsV2{}

	m2.status = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "metrics_v2_status",
		Help:      "Status code returned by /metrics/v2 endpoint (HTTP-style status).",
	}, []string{})

	m2.timestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "metrics_v2_timestamp",
		Help:      "Response timestamp from /metrics/v2 endpoint (UNIX epoch).",
	}, []string{})

	m2.info = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "metrics_v2_info",
		Help:      "Information fields from /metrics/v2 exported as labels (always 1).",
	}, []string{"product", "vendor", "version", "agent", "protocol"})

	m2.config = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "metrics_v2_config",
		Help:      "Boolean/numeric representation of config items from /metrics/v2 (label key represents config item name).",
	}, []string{"key"})

	return m2
}

// Register 注册所有 metrics v2 指标
func (m2 *MetricsV2) Register() {
	prometheus.MustRegister(m2.status)
	prometheus.MustRegister(m2.timestamp)
	prometheus.MustRegister(m2.info)
	prometheus.MustRegister(m2.config)
}

// UpdateMetricsV2 更新 metrics v2 相关指标
func (m2 *MetricsV2) UpdateMetricsV2(metricsV2 MetricsV2Data) {
	// status and timestamp (no labels)
	if m2.status != nil {
		m2.status.WithLabelValues().Set(float64(metricsV2.Status))
	}
	if m2.timestamp != nil {
		m2.timestamp.WithLabelValues().Set(float64(metricsV2.Timestamp))
	}

	// value map may contain agent/info/config
	if val := metricsV2.Value; val != nil {
		// info fields -> expose as labels on info metric
		// Try common keys at top-level value
		product := toString(val["product"])
		vendor := toString(val["vendor"])
		version := toString(val["version"])
		agent := toString(val["agent"])
		protocol := toString(val["protocol"])
		if product == "" {
			// fallback: value["info"] may be a map
			if info, ok := val["info"].(map[string]interface{}); ok {
				product = toString(info["product"])
				vendor = toString(info["vendor"])
				version = toString(info["version"])
				agent = toString(info["agent"])
				protocol = toString(info["protocol"])
			}
		}
		if m2.info != nil {
			m2.info.With(prometheus.Labels{"product": product, "vendor": vendor, "version": version, "agent": agent, "protocol": protocol}).Set(1)
		}

		// config map: many keys likely present under value.config
		if cfg, ok := val["config"].(map[string]interface{}); ok {
			for k, v := range cfg {
				if m2.config != nil {
					// try boolean/number conversion
					m2.config.With(prometheus.Labels{"key": k}).Set(boolOrNumberToFloat(v))
				}
			}
		}
	}
}

// MetricsV2Data metrics v2 数据结构
type MetricsV2Data struct {
	Status    int
	Timestamp int64
	Value     map[string]interface{}
}
