package puppetdb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

		// 获取待处理连接数
		pendingValue, err := mc.getMBeanValue(fmt.Sprintf("puppetlabs.puppetdb.database:name=%s.pool.PendingConnections", pool))
		if err == nil {
			metrics[pool]["PendingConnections"] = pendingValue
		}
	}

	return metrics, nil
}

// GetDBPoolUsageMetrics 获取数据库连接池使用统计指标
func (mc *MetricsClient) GetDBPoolUsageMetrics() (map[string]map[string]float64, error) {
	metrics := make(map[string]map[string]float64)
	pools := []string{"PDBReadPool", "PDBWritePool"}

	for _, pool := range pools {
		metrics[pool] = make(map[string]float64)

		// 获取Usage统计
		usageResult, err := mc.getMBeanFullData(fmt.Sprintf("puppetlabs.puppetdb.database:name=%s.pool.Usage", pool))
		if err == nil && usageResult != nil {
			// 提取关键百分位数和统计信息
			if val, ok := usageResult["Mean"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["UsageMean"] = floatVal
				}
			}
			if val, ok := usageResult["75thPercentile"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["Usage75thPercentile"] = floatVal
				}
			}
			if val, ok := usageResult["95thPercentile"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["Usage95thPercentile"] = floatVal
				}
			}
			if val, ok := usageResult["99thPercentile"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["Usage99thPercentile"] = floatVal
				}
			}
			if val, ok := usageResult["Max"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["UsageMax"] = floatVal
				}
			}
		}

		// 获取Wait统计
		waitResult, err := mc.getMBeanFullData(fmt.Sprintf("puppetlabs.puppetdb.database:name=%s.pool.Wait", pool))
		if err == nil && waitResult != nil {
			// 提取关键百分位数和统计信息
			if val, ok := waitResult["Mean"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["WaitMean"] = floatVal / 1000.0 // 转换为秒
				}
			}
			if val, ok := waitResult["75thPercentile"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["Wait75thPercentile"] = floatVal / 1000.0 // 转换为秒
				}
			}
			if val, ok := waitResult["95thPercentile"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["Wait95thPercentile"] = floatVal / 1000.0 // 转换为秒
				}
			}
			if val, ok := waitResult["99thPercentile"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["Wait99thPercentile"] = floatVal / 1000.0 // 转换为秒
				}
			}
			if val, ok := waitResult["Max"]; ok {
				if floatVal, err := mc.parseMBeanValue(val, pool); err == nil {
					metrics[pool]["WaitMax"] = floatVal / 1000.0 // 转换为秒
				}
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
	endpoints := []string{"/pdb/query/v4/nodes", "/pdb/query/v4/resources", "/metrics/v2/read", "/metrics/v2"}

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
	// 使用 /metrics/v2/read 接口替代 /metrics/v1/mbeans
	endpoint := "/metrics/v2/read"

	// 准备请求体，使用新的格式
	requestBody, err := json.Marshal(map[string]interface{}{
		"mbean": mbeanName,
	})
	if err != nil {
		return 0, err
	}

	resp, err := mc.Post(endpoint, "application/json", requestBody)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("failed to get mbean %s: status %d", mbeanName, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// 根据不同的MBean类型提取数值
	return mc.extractValueFromMBean(result, mbeanName)
}

// getMBeanFullData 获取MBean的完整数据（不提取具体值）
func (mc *MetricsClient) getMBeanFullData(mbeanName string) (map[string]interface{}, error) {
	// 使用 /metrics/v2/read 接口替代 /metrics/v1/mbeans
	endpoint := "/metrics/v2/read"

	// 准备请求体，使用新的格式
	requestBody, err := json.Marshal(map[string]interface{}{
		"mbean": mbeanName,
	})
	if err != nil {
		return nil, err
	}

	resp, err := mc.Post(endpoint, "application/json", requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get mbean %s: status %d", mbeanName, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 返回value字段中的数据
	if val, ok := result["value"]; ok {
		if valueMap, ok := val.(map[string]interface{}); ok {
			return valueMap, nil
		}
	}

	return nil, fmt.Errorf("no value data found in MBean result for %s", mbeanName)
}

// GetMetricsBulk 批量获取多个MBean指标
func (mc *MetricsClient) GetMetricsBulk(mbeanNames []string) (map[string]float64, error) {
	endpoint := "/metrics/v2/read"

	// 准备请求体，使用新的格式
	request := make([]map[string]interface{}, len(mbeanNames))
	for i, mbeanName := range mbeanNames {
		request[i] = map[string]interface{}{
			"mbean": mbeanName,
		}
	}

	requestBody, err := json.Marshal(request)
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
	// 首先尝试提取 value 字段中的 Value 子字段（新的 /metrics/v2/read 格式）
	if val, ok := result["value"]; ok {
		if valueMap, ok := val.(map[string]interface{}); ok {
			// 检查 value.Value 字段
			if innerVal, ok := valueMap["Value"]; ok {
				return mc.parseMBeanValue(innerVal, mbeanName)
			}
			// 检查其他可能的字段
			if innerVal, ok := valueMap["value"]; ok {
				return mc.parseMBeanValue(innerVal, mbeanName)
			}
		} else {
			// 如果 value 不是 map，直接解析
			return mc.parseMBeanValue(val, mbeanName)
		}
	}

	// 兼容旧格式：直接提取 Value 字段
	if val, ok := result["Value"]; ok {
		return mc.parseMBeanValue(val, mbeanName)
	}

	// 对于某些指标，可能需要提取特定的字段
	if val, ok := result["Count"]; ok {
		return mc.parseMBeanValue(val, mbeanName)
	}

	if val, ok := result["Mean"]; ok {
		return mc.parseMBeanValue(val, mbeanName)
	}

	// 对于某些MBean，可能需要检查特定的属性
	// 例如，内存相关的MBean可能有 used, max 等字段
	if strings.Contains(mbeanName, "Memory") {
		if val, ok := result["used"]; ok {
			return mc.parseMBeanValue(val, mbeanName)
		}
		if val, ok := result["max"]; ok {
			return mc.parseMBeanValue(val, mbeanName)
		}
	}

	// 对于线程相关的MBean
	if strings.Contains(mbeanName, "Threading") {
		if val, ok := result["ThreadCount"]; ok {
			return mc.parseMBeanValue(val, mbeanName)
		}
	}

	// 对于HTTP相关的MBean，检查 mean 或 count 字段
	if strings.Contains(mbeanName, "http") {
		if val, ok := result["mean"]; ok {
			return mc.parseMBeanValue(val, mbeanName)
		}
		if val, ok := result["count"]; ok {
			return mc.parseMBeanValue(val, mbeanName)
		}
	}

	return 0, fmt.Errorf("no suitable value found in MBean result for %s", mbeanName)
}

// parseMBeanValue 解析MBean值，支持多种格式
func (mc *MetricsClient) parseMBeanValue(val interface{}, mbeanName string) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		// 处理字符串格式的数值，包括分数格式如 "4/83"
		return mc.parseStringValue(v, mbeanName)
	default:
		return 0, fmt.Errorf("unsupported value type %T for MBean %s", val, mbeanName)
	}
}

// parseStringValue 解析字符串格式的数值
func (mc *MetricsClient) parseStringValue(strVal string, mbeanName string) (float64, error) {
	// 首先尝试直接解析为浮点数
	if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
		return floatVal, nil
	}

	// 处理分数格式，如 "4/83"
	if strings.Contains(strVal, "/") {
		parts := strings.Split(strVal, "/")
		if len(parts) == 2 {
			numerator, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			denominator, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil && denominator != 0 {
				return numerator / denominator, nil
			}
		}
	}

	return 0, fmt.Errorf("cannot parse string value '%s' for MBean %s", strVal, mbeanName)
}

// GetAvailableMBeans 获取所有可用的MBean列表
func (mc *MetricsClient) GetAvailableMBeans() ([]string, error) {
	// 使用 /metrics/v2/list 接口替代 /metrics/v1/mbeans
	var result map[string]interface{}
	err := mc.get("/metrics/v2/list", "", &result)
	if err != nil {
		return nil, err
	}

	// 解析嵌套的JSON结构，提取MBean名称
	var mbeans []string

	// 检查value字段是否存在
	if value, ok := result["value"]; ok {
		if valueMap, ok := value.(map[string]interface{}); ok {
			// 遍历value中的每个MBean及其属性
			for mbeanName, attributes := range valueMap {
				if attrsMap, ok := attributes.(map[string]interface{}); ok {
					// 遍历所有属性，构建完整的MBean名称
					for attrKey := range attrsMap {
						if attrKey == "" {
							// 空属性，只使用主名称
							mbeans = append(mbeans, mbeanName)
						} else {
							// 非空属性，组合成完整的MBean名称
							fullMBeanName := fmt.Sprintf("%s:%s", mbeanName, attrKey)
							mbeans = append(mbeans, fullMBeanName)
						}
					}
				} else {
					// 如果属性不是map类型，只添加主MBean名称
					mbeans = append(mbeans, mbeanName)
				}
			}
		}
	} else {
		// 如果没有value字段，直接遍历顶层对象（兼容旧格式）
		for mbeanName := range result {
			// 跳过request字段，只处理MBean名称
			if mbeanName != "request" {
				mbeans = append(mbeans, mbeanName)
			}
		}
	}

	return mbeans, nil
}
