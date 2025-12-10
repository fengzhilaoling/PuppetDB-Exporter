package puppetdb

import (
	"encoding/json"
	"fmt"
)

// MetricsClient 扩展的PuppetDB客户端，专门用于获取指标
type MetricsClient struct {
	*PuppetDB
}

// NewMetricsClient 创建指标客户端
func NewMetricsClient(pdb *PuppetDB) *MetricsClient {
	return &MetricsClient{PuppetDB: pdb}
}

// GetPopulationMetrics 获取人口统计指标
func (mc *MetricsClient) GetPopulationMetrics() (map[string]float64, error) {
	metrics := make(map[string]float64)

	// 获取节点数量
	nodeCount, err := mc.getMBeanValue("puppetlabs.puppetdb.population:name=num-nodes")
	if err == nil {
		metrics["nodes"] = nodeCount
	}

	// 获取资源数量
	resourceCount, err := mc.getMBeanValue("puppetlabs.puppetdb.population:name=num-resources")
	if err == nil {
		metrics["resources"] = resourceCount
	}

	// 获取平均资源数
	avgResources, err := mc.getMBeanValue("puppetlabs.puppetdb.population:name=avg-resources-per-node")
	if err == nil {
		metrics["avg_resources_per_node"] = avgResources
	}

	// 获取重复资源百分比
	duplicatePct, err := mc.getMBeanValue("puppetlabs.puppetdb.population:name=pct-resource-dupes")
	if err == nil {
		metrics["resource_duplicates_pct"] = duplicatePct
	}

	return metrics, nil
}

// GetStorageMetrics 获取存储层指标
func (mc *MetricsClient) GetStorageMetrics() (map[string]float64, error) {
	metrics := make(map[string]float64)

	// 获取重复编录百分比
	duplicatePct, err := mc.getMBeanValue("puppetlabs.puppetdb.storage:name=duplicate-pct")
	if err == nil {
		metrics["duplicate_pct"] = duplicatePct
	}

	// 获取GC时间
	gcTime, err := mc.getMBeanValue("puppetlabs.puppetdb.storage:name=gc-time")
	if err == nil {
		metrics["gc_time"] = gcTime
	}

	// 获取替换事实时间
	replaceFactsTime, err := mc.getMBeanValue("puppetlabs.puppetdb.storage:name=replace-facts-time")
	if err == nil {
		metrics["replace_facts_time"] = replaceFactsTime
	}

	// 获取替换编录时间
	replaceCatalogTime, err := mc.getMBeanValue("puppetlabs.puppetdb.storage:name=replace-catalog-time")
	if err == nil {
		metrics["replace_catalog_time"] = replaceCatalogTime
	}

	return metrics, nil
}

// GetCommandMetrics 获取命令处理指标
func (mc *MetricsClient) GetCommandMetrics() (map[string]map[string]float64, error) {
	metrics := make(map[string]map[string]float64)

	// 获取全局命令指标
	globalMetrics := []string{"seen", "processed", "fatal", "retried", "awaiting-retry", "depth"}
	for _, metric := range globalMetrics {
		value, err := mc.getMBeanValue(fmt.Sprintf("puppetlabs.puppetdb.mq:name=global.%s", metric))
		if err == nil {
			if metrics["global"] == nil {
				metrics["global"] = make(map[string]float64)
			}
			metrics["global"][metric] = value
		}
	}

	return metrics, nil
}

// GetDBMetrics 获取数据库连接池指标
func (mc *MetricsClient) GetDBMetrics() (map[string]map[string]float64, error) {
	metrics := make(map[string]map[string]float64)

	// HikariCP指标映射
	hikariMetrics := []string{"ActiveConnections", "IdleConnections", "TotalConnections", "WaitTime"}
	pools := []string{"PDBReadPool", "PDBWritePool"}

	for _, pool := range pools {
		metrics[pool] = make(map[string]float64)
		for _, metric := range hikariMetrics {
			value, err := mc.getMBeanValue(fmt.Sprintf("puppetlabs.puppetdb.database:%s.%s", pool, metric))
			if err == nil {
				metrics[pool][metric] = value
			}
		}
	}

	return metrics, nil
}

// GetJVMMetrics 获取JVM指标
func (mc *MetricsClient) GetJVMMetrics() (map[string]float64, error) {
	metrics := make(map[string]float64)

	// 内存指标
	memoryTypes := []string{"HeapMemoryUsage", "NonHeapMemoryUsage"}
	for _, memoryType := range memoryTypes {
		value, err := mc.getMBeanValue(fmt.Sprintf("java.lang:type=Memory.%s.used", memoryType))
		if err == nil {
			metrics[fmt.Sprintf("memory_%s_used", memoryType)] = value
		}

		maxValue, err := mc.getMBeanValue(fmt.Sprintf("java.lang:type=Memory.%s.max", memoryType))
		if err == nil {
			metrics[fmt.Sprintf("memory_%s_max", memoryType)] = maxValue
		}
	}

	// 线程指标
	threadCount, err := mc.getMBeanValue("java.lang:type=Threading.ThreadCount")
	if err == nil {
		metrics["threads_active"] = threadCount
	}

	return metrics, nil
}

// GetHTTPMetrics 获取HTTP服务指标
func (mc *MetricsClient) GetHTTPMetrics() (map[string]map[string]float64, error) {
	metrics := make(map[string]map[string]float64)

	// 主要HTTP端点
	endpoints := []string{"/pdb/query/v4/nodes", "/pdb/query/v4/resources", "/metrics/v1/mbeans", "/metrics/v2"}

	for _, endpoint := range endpoints {
		metrics[endpoint] = make(map[string]float64)

		// 获取服务时间
		serviceTime, err := mc.getMBeanValue(fmt.Sprintf("puppetlabs.puppetdb.http:name=%s.service-time.mean", endpoint))
		if err == nil {
			metrics[endpoint]["service_time_mean"] = serviceTime
		}

		// 获取请求计数（200状态码）
		requestCount, err := mc.getMBeanValue(fmt.Sprintf("puppetlabs.puppetdb.http:name=%s.200.count", endpoint))
		if err == nil {
			metrics[endpoint]["requests_200"] = requestCount
		}
	}

	return metrics, nil
}

// getMBeanValue 获取单个MBean指标值
func (mc *MetricsClient) getMBeanValue(mbeanName string) (float64, error) {
	endpoint := fmt.Sprintf("/metrics/v1/mbeans/%s", mbeanName)

	var result map[string]interface{}
	err := mc.get(endpoint, "", &result)
	if err != nil {
		return 0, err
	}

	// 根据不同的MBean类型提取数值
	return mc.extractValueFromMBean(result, mbeanName)
} // GetMetricsBulk 批量获取多个MBean指标
func (mc *MetricsClient) GetMetricsBulk(mbeanNames []string) (map[string]float64, error) {
	endpoint := "/metrics/v1/mbeans"

	// 准备请求体
	requestBody, err := json.Marshal(mbeanNames)
	if err != nil {
		return nil, err
	}

	resp, err := mc.Post(endpoint, "application/json", requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get bulk metrics: status %d", resp.StatusCode)
	}

	var results []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	// 解析结果
	metrics := make(map[string]float64)
	for i, result := range results {
		if i < len(mbeanNames) {
			value, err := mc.extractValueFromMBean(result, mbeanNames[i])
			if err == nil {
				metrics[mbeanNames[i]] = value
			}
		}
	}

	return metrics, nil
}

// extractValueFromMBean 从MBean结果中提取数值
func (mc *MetricsClient) extractValueFromMBean(result map[string]interface{}, mbeanName string) (float64, error) {
	// 根据MBean类型提取相应的值
	if val, ok := result["Value"]; ok {
		if floatVal, ok := val.(float64); ok {
			return floatVal, nil
		}
		if intVal, ok := val.(int); ok {
			return float64(intVal), nil
		}
	}

	// 对于某些指标，可能需要提取特定的字段
	if val, ok := result["Count"]; ok {
		if floatVal, ok := val.(float64); ok {
			return floatVal, nil
		}
	}

	if val, ok := result["Mean"]; ok {
		if floatVal, ok := val.(float64); ok {
			return floatVal, nil
		}
	}

	return 0, fmt.Errorf("no suitable value found in MBean result")
}

// GetAvailableMBeans 获取所有可用的MBean列表
func (mc *MetricsClient) GetAvailableMBeans() ([]string, error) {
	var result map[string]string
	err := mc.get("/metrics/v1/mbeans", "", &result)
	if err != nil {
		return nil, err
	}

	mbeans := make([]string, 0, len(result))
	for name := range result {
		mbeans = append(mbeans, name)
	}

	return mbeans, nil
}
