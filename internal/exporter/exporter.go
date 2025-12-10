package exporter

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/camptocamp/prometheus-puppetdb-exporter/internal/puppetdb"
)

// Exporter implements the prometheus.Exporter interface, and exports PuppetDB metrics
type Exporter struct {
	client          *puppetdb.PuppetDB
	metricsClient   *puppetdb.MetricsClient
	namespace       string
	metricsRegistry *MetricsRegistry
}

var (
	// reserved for future mapping
	_ = 0
)

// convertReportMetrics 将 puppetdb 的报告指标转换为内部格式
func convertReportMetrics(reportMetrics []puppetdb.ReportMetric) []ReportMetric {
	result := make([]ReportMetric, len(reportMetrics))
	for i, rm := range reportMetrics {
		result[i] = ReportMetric{
			Category: rm.Category,
			Name:     rm.Name,
			Value:    rm.Value,
		}
	}
	return result
}

// NewPuppetDBExporter returns a new exporter of PuppetDB metrics.
func NewPuppetDBExporter(url, certPath, caPath, keyPath string, sslSkipVerify bool, categories map[string]struct{}) (e *Exporter, err error) {
	e = &Exporter{
		namespace: "puppetdb",
	}

	// 创建指标注册表
	e.metricsRegistry = NewMetricsRegistry(e.namespace, categories)
	e.metricsRegistry.RegisterAll()

	opts := &puppetdb.Options{
		URL:        url,
		CertPath:   certPath,
		CACertPath: caPath,
		KeyPath:    keyPath,
		SSLVerify:  sslSkipVerify,
	}

	e.client, err = puppetdb.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create PuppetDB client: %v", err)
	}

	// 创建MetricsClient - 必须在e.client初始化之后
	e.metricsClient = puppetdb.NewMetricsClient(e.client)
	if err != nil {
		log.Fatalf("failed to create new client: %s", err)
		return
	}

	return
}

// Describe outputs PuppetDB metric descriptions
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	// 指标注册表会处理所有指标描述
	e.metricsRegistry.Describe(ch)
}

// Collect fetches new metrics from the PuppetDB and updates the appropriate metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// 指标注册表会处理所有指标收集
	e.metricsRegistry.Collect(ch)
}

// Scrape scrapes PuppetDB and update metrics
func (e *Exporter) Scrape(interval time.Duration, unreportedNode string, categories map[string]struct{}) {
	var statuses map[string]int

	unreportedDuration, err := time.ParseDuration(unreportedNode)
	if err != nil {
		log.Errorf("failed to parse unreported duration: %s", err)
		return
	}

	for {
		statuses = make(map[string]int)

		// 记录节点抓取开始时间
		scrapeStart := time.Now()
		nodes, err := e.client.Nodes()
		if err != nil {
			log.Errorf("failed to get nodes: %s", err)
			e.metricsRegistry.GetPerformanceMetrics().RecordScrapeError("nodes", "connection_error")
		}
		// 记录节点抓取耗时
		e.metricsRegistry.GetPerformanceMetrics().RecordScrapeDuration("nodes", time.Since(scrapeStart).Seconds())

		// 重置指标
		e.metricsRegistry.GetNodeMetrics().Reset()
		e.metricsRegistry.GetServiceMetrics().Reset()

		for _, node := range nodes {
			var deactivated string
			if node.Deactivated == "" {
				deactivated = "false"
			} else {
				deactivated = "true"
			}

			if node.ReportTimestamp == "" {
				if deactivated == "false" {
					statuses["unreported"]++
				}
				continue
			}
			latestReport, err := time.Parse(time.RFC3339, node.ReportTimestamp)
			if err != nil {
				if deactivated == "false" {
					statuses["unreported"]++
				}
				log.Errorf("failed to parse report timestamp: %s", err)
				continue
			}

			// 创建节点信息结构体
			nodeInfo := NodeInfo{
				Certname:                node.Certname,
				ReportEnvironment:       node.ReportEnvironment,
				ReportTimestamp:         node.ReportTimestamp,
				Deactivated:             node.Deactivated,
				LatestReportHash:        node.LatestReportHash,
				LatestReportNoop:        node.LatestReportNoop,
				LatestReportNoopPending: node.LatestReportNoopPending,
				LatestReportStatus:      node.LatestReportStatus,
				CachedCatalogStatus:     node.CachedCatalogStatus,
				CatalogTimestamp:        node.CatalogTimestamp,
				FactsTimestamp:          node.FactsTimestamp,
			}

			// 更新节点指标
			e.metricsRegistry.GetNodeMetrics().UpdateNodeMetrics(nodeInfo, unreportedDuration, time.Now())

			if deactivated == "false" {
				if latestReport.Add(unreportedDuration).Before(time.Now()) {
					statuses["unreported"]++
				} else if node.LatestReportStatus == "" {
					statuses["unreported"]++
				} else {
					statuses[node.LatestReportStatus]++
				}
			}

			if node.LatestReportHash != "" {
				reportMetrics, _ := e.client.ReportMetrics(node.LatestReportHash)
				e.metricsRegistry.GetNodeMetrics().UpdateReportMetrics(nodeInfo, convertReportMetrics(reportMetrics))
			}
		}

		// Scrape service status endpoints and expose metrics
		serviceScrapeStart := time.Now()
		services, serr := e.client.Services()
		if serr != nil {
			log.Errorf("failed to get services: %s", serr)
			e.metricsRegistry.GetPerformanceMetrics().RecordScrapeError("services", "connection_error")
		}
		e.metricsRegistry.GetPerformanceMetrics().RecordScrapeDuration("services", time.Since(serviceScrapeStart).Seconds())

		if serr == nil {
			// 转换服务信息格式
			serviceInfos := make([]ServiceInfo, 0)
			for svcName, info := range services {
				serviceInfo := ServiceInfo{
					Name:    svcName,
					Version: info.ServiceVersion,
					State:   info.State,
					Up:      true, // 简化处理，假设服务正常运行
				}
				serviceInfos = append(serviceInfos, serviceInfo)
			}
			e.metricsRegistry.GetServiceMetrics().UpdateServiceMetrics(serviceInfos)
		}

		// Scrape /metrics/v2 and expose useful values
		metricsV2ScrapeStart := time.Now()
		metricsV2, merr := e.client.MetricsV2()
		if merr != nil {
			log.Errorf("failed to get metrics v2: %s", merr)
			e.metricsRegistry.GetPerformanceMetrics().RecordScrapeError("metrics_v2", "connection_error")
		}
		e.metricsRegistry.GetPerformanceMetrics().RecordScrapeDuration("metrics_v2", time.Since(metricsV2ScrapeStart).Seconds())

		if merr == nil {
			// 转换 metrics v2 数据格式
			metricsV2Data := MetricsV2Data{
				Status:    metricsV2.Status,
				Timestamp: metricsV2.Timestamp,
				Value:     metricsV2.Value,
			}
			e.metricsRegistry.GetMetricsV2().UpdateMetricsV2(metricsV2Data)
		}

		// 更新节点状态计数
		e.metricsRegistry.GetNodeMetrics().UpdateStatusCount(statuses)

		// 更新系统健康评分
		e.metricsRegistry.GetSystemMetrics().UpdateSystemMetrics(statuses)

		// 收集PuppetDB核心指标
		if e.metricsClient != nil {
			// 收集人口统计指标
			populationMetrics, err := e.metricsClient.GetPopulationMetrics()
			if err == nil {
				e.metricsRegistry.GetPuppetDBMetrics().UpdatePopulationMetrics(
					populationMetrics["nodes"],
					populationMetrics["resources"],
					populationMetrics["avg_resources_per_node"],
				)
			}

			// 收集存储层指标
			storageMetrics, err := e.metricsClient.GetStorageMetrics()
			if err == nil {
				e.metricsRegistry.GetPuppetDBMetrics().UpdateStorageMetrics(
					storageMetrics["duplicate_pct"],
					storageMetrics["gc_time"],
					storageMetrics["replace_facts_time"],
					storageMetrics["replace_catalog_time"],
				)
			}

			// 收集命令处理指标
			commandMetrics, err := e.metricsClient.GetCommandMetrics()
			if err == nil && commandMetrics["global"] != nil {
				global := commandMetrics["global"]
				if depth, ok := global["depth"]; ok {
					e.metricsRegistry.GetPuppetDBMetrics().UpdateCommandMetrics("global", "", "", 0)
					// 设置队列深度
					e.metricsRegistry.GetPuppetDBMetrics().UpdateCommandQueueDepth(depth)
				}
			}

			// 收集数据库指标
			dbMetrics, err := e.metricsClient.GetDBMetrics()
			if err == nil {
				for pool, metrics := range dbMetrics {
					e.metricsRegistry.GetPuppetDBMetrics().UpdateDBMetrics(
						pool,
						metrics["ActiveConnections"],
						metrics["IdleConnections"],
						metrics["TotalConnections"],
						metrics["WaitTime"],
					)
				}
			}

			// 收集JVM指标
			jvmMetrics, err := e.metricsClient.GetJVMMetrics()
			if err == nil {
				// 更新内存指标
				if used, ok := jvmMetrics["memory_HeapMemoryUsage_used"]; ok {
					e.metricsRegistry.GetPuppetDBMetrics().UpdateJVMMetrics("heap", used, -1, -1, "", 0)
				}
				if max, ok := jvmMetrics["memory_HeapMemoryUsage_max"]; ok {
					e.metricsRegistry.GetPuppetDBMetrics().UpdateJVMMetrics("heap", -1, max, -1, "", 0)
				}
				// 更新线程指标
				if threads, ok := jvmMetrics["threads_active"]; ok {
					e.metricsRegistry.GetPuppetDBMetrics().UpdateJVMMetrics("", -1, -1, threads, "", 0)
				}
			}
		}

		time.Sleep(interval)
	}
}
