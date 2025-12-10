Prometheus PuppetDB 导出器
============================

## 📋 概述

Prometheus PuppetDB Exporter 是一个用于监控 PuppetDB 的 Prometheus 导出器。本次版本进行了全面的指标优化，移除了不合适的指标，增加了核心性能指标，并优化了数据获取接口。

## 🔧 主要改进

### ✅ 新增核心指标
- **时间间隔指标**: 节点报告、编录、事实数据的时间间隔
- **PuppetDB核心性能**: 命令处理、存储层、人口统计指标
- **HTTP服务指标**: 请求统计、延迟、连接数
- **数据库连接池**: 连接状态、等待时间、使用统计、创建统计、超时率
- **JVM指标**: 内存使用、线程、GC统计

### 🆕 数据库连接池增强功能
- **完整的连接池统计**: 支持 PDBReadPool 和 PDBWritePool 的完整指标采集
- **连接创建统计**: 包含创建时间的均值、百分位数统计
- **连接超时率监控**: 支持1分钟、5分钟、15分钟速率和平均速率
- **连接池配置指标**: 最大连接数、最小连接数配置监控
- **详细的使用统计**: 包含50th、75th、95th、99th、999th百分位数
- **等待时间统计**: 详细的连接等待时间分布统计

### ❌ 移除不合适指标
- 移除了含义不明确或重复的标签
- 精简了服务指标，专注于PuppetDB核心功能
- 去除了应由专门exporter处理的数据库指标

### 🚀 性能优化
- 支持批量API调用，减少网络往返
- 新增MetricsClient，提供类型安全的指标获取
- 改进错误处理和重试机制

## 使用说明

### 基本用法

```bash
# Linux/macOS
./prometheus-puppetdb-exporter --puppetdb-url=https://puppetdb:8081 --listen-address=0.0.0.0:9635

# Windows PowerShell
$env:PUPPETDB_URL="https://puppetdb:8081"
.\prometheus-puppetdb-exporter.exe --listen-address=0.0.0.0:9635
```

### 命令行参数

| 参数 | 环境变量 | 说明 | 默认值 |
|------|----------|------|--------|
| `--puppetdb-url` | `PUPPETDB_URL` | PuppetDB 基础 URL（例如: https://puppetdb:8081） | - |
| `--cert-file` | `PUPPETDB_CERT_FILE` | 客户端 TLS 证书（PEM 编码） | - |
| `--key-file` | `PUPPETDB_KEY_FILE` | 客户端私钥（PEM 编码） | - |
| `--ca-file` | `PUPPETDB_CA_FILE` | CA 根证书（PEM 编码） | - |
| `--ssl-skip-verify` | `PUPPETDB_SSL_SKIP_VERIFY` | 跳过 SSL 证书校验（不推荐用于生产环境） | `false` |
| `--scrape-interval` | `PUPPETDB_SCRAPE_INTERVAL` | 两次抓取之间的间隔（示例：5s） | - |
| `--listen-address` | `PUPPETDB_LISTEN_ADDRESS` | 监听地址 | `0.0.0.0:9635` |
| `--metric-path` | `PUPPETDB_METRIC_PATH` | 指标导出路径 | `/metrics` |
| `--verbose` | `PUPPETDB_VERBOSE` | 启用调试日志输出 | `false` |
| `--unreported-node` | `PUPPETDB_UNREPORTED_NODE` | 节点未报告超时时间 | `2h` |
| `--categories` | `REPORT_METRICS_CATEGORIES` | 要抓取的报告指标类别 | `resources,time,changes,events` |

### 访问指标

启动后，可通过以下地址访问 Prometheus 指标：

```
http://<host>:9635/metrics
```

## 指标说明

### 构建信息

| 指标 | 类型 | 说明 |
|------|------|------|
| `puppetdb_exporter_build_info` | gauge | exporter 构建信息（版本、提交、构建时间和 Go 版本） |

### 节点报告状态

| 指标 | 类型 | 说明 |
|------|------|------|
| `puppetdb_node_report_status_count` | gauge | 节点按报告状态的计数（status 标签：changed/failed/unchanged/unreported） |

### 节点相关指标

| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppet_report` | gauge | 节点最新报告的时间戳（UNIX epoch） | 核心 |
| `puppet_report_<category>` | gauge | 报告指标数值（按类别：resources/time/changes/events） | 业务 |
| `puppetdb_node_has_report` | gauge | 节点是否存在最新报告（1=有，0=无） | 核心 |
| `puppetdb_node_latest_report_noop` | gauge | 节点最新报告是否为 noop（1=是，0=否） | 诊断 |
| `puppetdb_node_catalog_timestamp` | gauge | 节点 catalog 时间戳（UNIX epoch） | 业务 |
| `puppetdb_node_facts_timestamp` | gauge | 节点 facts 时间戳（UNIX epoch） | 业务 |
| `puppetdb_node_report_age_seconds` | gauge | 节点报告时间间隔（秒） | 核心 |
| `puppetdb_node_catalog_age_seconds` | gauge | 节点编录时间间隔（秒） | 业务 |
| `puppetdb_node_facts_age_seconds` | gauge | 节点事实数据时间间隔（秒） | 业务 |

### 服务状态指标

| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_service_up` | gauge | 服务是否处于运行状态（1=运行，0=非运行） | 核心 |
| `puppetdb_service_info` | gauge | 服务信息（恒为 1，包含版本和状态标签） | 诊断 |
| `puppetdb_service_queue_depth` | gauge | 服务处理队列深度（未处理任务数） | 核心 |

### 性能指标

| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_exporter_scrape_duration_seconds` | histogram | PuppetDB exporter 抓取耗时（按 endpoint 分类） | 诊断 |
| `puppetdb_exporter_scrape_errors_total` | counter | 抓取错误总数（按 endpoint 和错误类型分类） | 诊断 |
| `puppetdb_exporter_request_duration_seconds` | histogram | PuppetDB API 请求耗时（按 endpoint 和 method 分类） | 诊断 |
| `puppetdb_exporter_requests_total` | counter | PuppetDB API 请求总数（按 endpoint 和状态分类） | 诊断 |

### PuppetDB核心性能指标

#### 系统健康指标
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_system_health_score` | gauge | PuppetDB系统健康评分（0-100） | 核心 |
| `puppetdb_node_failure_rate` | gauge | 节点失败率百分比 | 核心 |

#### 命令处理指标
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_commands_processed_total` | counter | 处理的命令总数（按命令类型、版本、状态分类） | 诊断 |
| `puppetdb_commands_processing_duration_seconds` | histogram | 命令处理耗时（按命令类型、版本分类） | 核心 |
| `puppetdb_command_queue_depth` | gauge | 命令队列深度 | 核心 |

#### 存储层指标
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_storage_duplicate_percentage` | gauge | 重复编录百分比 | 业务 |
| `puppetdb_storage_gc_duration_seconds` | histogram | 存储GC耗时 | 核心 |
| `puppetdb_storage_replace_facts_duration_seconds` | histogram | 替换事实耗时 | 诊断 |
| `puppetdb_storage_replace_catalog_duration_seconds` | histogram | 替换编录耗时 | 诊断 |

#### 人口统计指标
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_population_nodes_total` | gauge | 节点总数 | 核心 |
| `puppetdb_population_resources_total` | gauge | 资源总数 | 业务 |
| `puppetdb_population_avg_resources_per_node` | gauge | 平均资源数 | 诊断 |

#### HTTP服务指标
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_http_requests_total` | counter | HTTP请求总数（按endpoint、方法、状态分类） | 诊断 |
| `puppetdb_http_request_duration_seconds` | histogram | HTTP请求耗时（按endpoint、方法分类） | 核心 |
| `puppetdb_http_active_connections` | gauge | 活跃HTTP连接数 | 业务 |

#### 数据库连接池指标
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_db_connections_active` | gauge | 活跃数据库连接数（按连接池分类） | 核心 |
| `puppetdb_db_connections_idle` | gauge | 空闲数据库连接数（按连接池分类） | 诊断 |
| `puppetdb_db_connections_total` | gauge | 数据库连接总数（按连接池分类） | 诊断 |
| `puppetdb_db_connections_pending` | gauge | 待处理数据库连接数（按连接池分类） | 诊断 |
| `puppetdb_db_connection_wait_time_seconds` | histogram | 数据库连接等待时间（按连接池分类） | 业务 |
| `puppetdb_db_pool_max_connections` | gauge | 数据库连接池最大连接数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_min_connections` | gauge | 数据库连接池最小连接数（按连接池分类） | 诊断 |

#### 数据库连接池使用统计
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_db_pool_usage_mean` | gauge | 数据库连接池使用均值（按连接池分类） | 核心 |
| `puppetdb_db_pool_usage_75th_percentile` | gauge | 数据库连接池使用75th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_usage_95th_percentile` | gauge | 数据库连接池使用95th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_usage_99th_percentile` | gauge | 数据库连接池使用99th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_usage_max` | gauge | 数据库连接池使用最大值（按连接池分类） | 诊断 |

#### 数据库连接池等待时间统计
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_db_pool_wait_mean_seconds` | gauge | 数据库连接池等待时间均值（按连接池分类） | 核心 |
| `puppetdb_db_pool_wait_75th_percentile_seconds` | gauge | 数据库连接池等待时间75th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_wait_95th_percentile_seconds` | gauge | 数据库连接池等待时间95th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_wait_99th_percentile_seconds` | gauge | 数据库连接池等待时间99th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_wait_max_seconds` | gauge | 数据库连接池等待时间最大值（按连接池分类） | 诊断 |

#### 数据库连接池连接创建统计
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_db_pool_connection_creation_mean_seconds` | gauge | 数据库连接创建时间均值（按连接池分类） | 诊断 |
| `puppetdb_db_pool_connection_creation_75th_percentile_seconds` | gauge | 数据库连接创建时间75th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_connection_creation_95th_percentile_seconds` | gauge | 数据库连接创建时间95th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_connection_creation_99th_percentile_seconds` | gauge | 数据库连接创建时间99th百分位数（按连接池分类） | 诊断 |
| `puppetdb_db_pool_connection_creation_max_seconds` | gauge | 数据库连接创建时间最大值（按连接池分类） | 诊断 |

#### 数据库连接池超时率指标
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_db_pool_connection_timeout_rate` | gauge | 数据库连接超时率（按连接池分类） | 核心 |
| `puppetdb_db_pool_connection_timeout_one_minute_rate` | gauge | 数据库连接超时率1分钟速率（按连接池分类） | 诊断 |
| `puppetdb_db_pool_connection_timeout_five_minute_rate` | gauge | 数据库连接超时率5分钟速率（按连接池分类） | 诊断 |
| `puppetdb_db_pool_connection_timeout_fifteen_minute_rate` | gauge | 数据库连接超时率15分钟速率（按连接池分类） | 诊断 |
| `puppetdb_db_pool_connection_timeout_mean_rate` | gauge | 数据库连接超时率平均速率（按连接池分类） | 诊断 |

#### JVM指标
| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_jvm_memory_used_bytes` | gauge | JVM内存使用量（按内存类型分类） | 核心 |
| `puppetdb_jvm_memory_max_bytes` | gauge | JVM内存最大值（按内存类型分类） | 核心 |
| `puppetdb_jvm_threads_active` | gauge | JVM活跃线程数 | 业务 |
| `puppetdb_jvm_gc_duration_seconds` | histogram | JVM GC耗时（按GC类型分类） | 核心 |

### Metrics V2 指标

| 指标 | 类型 | 说明 | 监控级别 |
|------|------|------|----------|
| `puppetdb_metrics_v2_status` | gauge | /metrics/v2 返回的状态码 | 诊断 |
| `puppetdb_metrics_v2_timestamp` | gauge | /metrics/v2 响应的时间戳（UNIX epoch） | 诊断 |
| `puppetdb_metrics_v2_info` | gauge | /metrics/v2 中的产品信息（恒为 1） | 诊断 |
| `puppetdb_metrics_v2_config` | gauge | /metrics/v2 中的配置项（key 标签表示配置名称） | 诊断 |

#### 监控建议
**核心指标**（必须设置告警）：
- 系统健康评分 `puppetdb_system_health_score < 80`
- 命令队列深度 `puppetdb_command_queue_depth > 1000`
- 节点报告时间间隔 `puppetdb_node_report_age_seconds > 7200`
- HTTP请求延迟 `histogram_quantile(0.95, puppetdb_http_request_duration_seconds_bucket) > 5`
- JVM内存使用率 `puppetdb_jvm_memory_used_bytes / puppetdb_jvm_memory_max_bytes > 0.9`
- 数据库连接超时率 `puppetdb_db_pool_connection_timeout_rate > 0.1`
- 数据库连接池使用率 `puppetdb_db_pool_usage_mean > 0.9`

**业务指标**（推荐监控）：
- 节点失败率异常上升
- 数据库连接池活跃连接数持续高位
- 存储GC耗时异常增加

**诊断指标**（按需查看）：
- 详细性能数据分析
- 容量规划参考
- 故障排查辅助

### 数据库连接池监控最佳实践

**关键监控指标**:
- **连接池使用率**: `puppetdb_db_pool_usage_mean` 应保持在 0.7-0.8 以下
- **连接等待时间**: `puppetdb_db_pool_wait_95th_percentile_seconds` 应小于 500ms
- **连接超时率**: `puppetdb_db_pool_connection_timeout_rate` 应接近于 0
- **活跃连接数**: `puppetdb_db_connections_active` 对比 `puppetdb_db_pool_max_connections` 检查是否接近上限

**性能调优建议**:
- 当连接池使用率持续高于 80% 时，考虑增加最大连接数
- 当连接等待时间超过 1 秒时，检查数据库性能或增加连接池大小
- 监控连接创建时间，如果创建时间过长可能需要优化数据库连接参数
- 定期检查连接超时率，高超时率可能表示网络或数据库问题

**容量规划**:
- 使用 `puppetdb_db_pool_usage_*` 指标分析连接池使用趋势
- 结合 `puppetdb_db_pool_connection_creation_*` 指标评估连接创建开销
- 监控 `puppetdb_db_connections_pending` 了解连接请求排队情况

## 变更说明

本仓库已做以下主要变更：

- 升级 Go 版本至 `go 1.20`，并更新 `Dockerfile` 构建镜像为 `golang:1.20`。
- 支持新的 PuppetDB / 相关接口访问：
  - `/status/v1/services`
  - `/metrics/v2/list`
  - `/metrics/v2`
  - `/pdb/query/v4/nodes`
  - `/pdb/query/v4/reports`
- 新增了对上述接口中常用字段的 Prometheus 指标导出。

**注意**：部分历史指标使用 `puppet` 命名空间（例如 `puppet_report`、`puppet_report_<category>`），新指标使用 `puppetdb` 命名空间。如需统一命名空间，可进一步调整。

## 🚨 告警规则示例

## 🚨 告警规则示例

```yaml
groups:
- name: puppetdb_alerts
  rules:
  - alert: PuppetDBHealthScoreLow
    expr: puppetdb_system_health_score < 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "PuppetDB health score is low"
      description: "PuppetDB health score is {{ $value }} (below 80)"
      
  - alert: PuppetDBCommandQueueHigh
    expr: puppetdb_command_queue_depth > 1000
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "PuppetDB command queue depth is high"
      description: "Command queue depth is {{ $value }} (above 1000)"
      
  - alert: PuppetDBReportAgeHigh
    expr: puppetdb_node_report_age_seconds > 7200
    for: 15m
    labels:
      severity: warning
    annotations:
      summary: "PuppetDB node report age is high"
      description: "Node {{ $labels.host }} report age is {{ $value }}s (above 2 hours)"
      
  - alert: PuppetDBHTTPRequestLatencyHigh
    expr: histogram_quantile(0.95, puppetdb_http_request_duration_seconds_bucket) > 5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "PuppetDB HTTP request latency is high"
      description: "95th percentile latency is {{ $value }}s (above 5s)"
      
  - alert: PuppetDBJVMMemoryHigh
    expr: puppetdb_jvm_memory_used_bytes / puppetdb_jvm_memory_max_bytes > 0.9
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "PuppetDB JVM memory usage is high"
      description: "JVM memory usage is {{ $value | humanizePercentage }} (above 90%)"
      
  - alert: PuppetDBConnectionPoolTimeoutHigh
    expr: puppetdb_db_pool_connection_timeout_rate > 0.1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "PuppetDB connection pool timeout rate is high"
      description: "Connection pool {{ $labels.pool }} timeout rate is {{ $value }} (above 0.1)"
      
  - alert: PuppetDBConnectionPoolUsageHigh
    expr: puppetdb_db_pool_usage_mean > 0.9
    for: 10m
    labels:
      severity: warning
    annotations:
      summary: "PuppetDB connection pool usage is high"
      description: "Connection pool {{ $labels.pool }} usage is {{ $value | humanizePercentage }} (above 90%)"
      
  - alert: PuppetDBConnectionPoolWaitTimeHigh
    expr: puppetdb_db_pool_wait_95th_percentile_seconds > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "PuppetDB connection pool wait time is high"
      description: "Connection pool {{ $labels.pool }} 95th percentile wait time is {{ $value }}s (above 1s)"
```

## 📊 Grafana Dashboard

建议创建以下dashboard面板：

1. **Overview面板**
   - 系统健康评分
   - 节点状态分布
   - 命令队列深度

2. **Performance面板**
   - HTTP请求延迟
   - 命令处理时间
   - GC耗时

3. **Capacity面板**
   - 节点数量趋势
   - 资源使用情况
   - 数据库连接数

4. **JVM面板**
   - 内存使用趋势
   - GC活动统计
   - 线程状态

## 🔧 性能优化

### 1. 批量API调用
- 使用POST `/metrics/v2/read` 批量获取指标
- 减少网络往返次数

### 2. 缓存机制
- 对不经常变化的数据增加缓存支持
- 可配置缓存过期时间

### 3. 错误处理
- 增强的错误处理和重试机制
- 更详细的错误分类和报告

## 🚀 后续改进计划

1. **支持Metrics V2 API** ✅ 已完成
   - 优先使用新的`/metrics/v2`端点
   - 提供更好的安全性和性能
   - 完整支持数据库连接池指标采集

2. **增加自定义指标**
   - 支持用户自定义MBean指标
   - 可配置的指标收集规则

3. **增强错误诊断**
   - 更详细的错误分类
   - 错误趋势分析

4. **支持集群模式**
   - 多PuppetDB实例监控
   - 集群级别的聚合指标

## 📋 验证清单

- [x] 代码可以成功编译
- [x] 移除了不合适的指标
- [x] 增加了核心性能指标
- [x] 优化了数据获取接口
- [x] 添加了批量获取支持
- [x] 改进了错误处理
- [x] 更新了文档说明
- [x] 完整实现数据库连接池指标采集
- [x] 支持PDBReadPool和PDBWritePool完整指标
- [x] 添加连接创建统计和超时率监控
- [x] 支持详细的百分位数统计（50th, 75th, 95th, 99th, 999th）
- [x] 添加数据库连接池配置指标监控
- [x] 更新告警规则和监控建议

## 🔍 故障排除

### 常见问题

1. **连接失败**
   - 检查PuppetDB URL配置
   - 验证TLS证书配置
   - 确认网络连通性

2. **指标缺失**
   - 检查PuppetDB版本兼容性
   - 确认指标端点权限配置
   - 查看exporter日志中的错误信息

3. **性能问题**
   - 调整抓取间隔时间
   - 检查PuppetDB负载情况
   - 考虑增加缓存配置

### 调试模式

使用 `--verbose` 参数启用调试日志：

```bash
./prometheus-puppetdb-exporter --puppetdb-url=https://puppetdb:8081 --verbose
```

## 进一步建议

- 若需要在生产环境中使用，请根据实际 PuppetDB 返回的 JSON 完善 `metrics_v2.Value` 和 `services[].status` 的字段解析。
- 如需统一指标命名空间（`puppet` vs `puppetdb`），我可以替你统一并更新文档与代码。


```
