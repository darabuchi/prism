# Prism - API 规格说明文档

## API 概述

Prism Core 提供 RESTful API 接口，供各平台客户端调用。API 基于 HTTP/HTTPS 协议，使用 JSON 格式进行数据交换。

### Base URL
```
http://localhost:9090/api/v1
```

### 认证方式
使用 Bearer Token 认证：
```
Authorization: Bearer <token>
```

## 核心 API 接口

### 1. 系统管理

#### 1.1 获取系统状态
```http
GET /system/status
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version": "1.0.0",
    "uptime": 3600,
    "memory_usage": 128.5,
    "cpu_usage": 15.2,
    "proxy_status": "running",
    "connected_clients": 5,
    "total_upload": 1048576,
    "total_download": 2097152
  }
}
```

#### 1.2 获取系统配置
```http
GET /system/config
```

#### 1.3 更新系统配置
```http
PUT /system/config
Content-Type: application/json

{
  "log_level": "info",
  "api_port": 9090,
  "proxy_port": 7890,
  "allow_lan": false,
  "bind_address": "127.0.0.1"
}
```

### 2. 订阅管理

#### 2.1 获取订阅列表
```http
GET /subscriptions
```

**查询参数:**
- `page`: 页码，默认 1
- `size`: 每页大小，默认 20
- `status`: 按状态过滤 (active/inactive/error)
- `auto_update`: 按自动更新状态过滤

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 10,
    "page": 1,
    "size": 20,
    "subscriptions": [
      {
        "id": 1,
        "name": "高速订阅",
        "url": "https://example.com/subscribe",
        "user_agent": "clash",
        "auto_update": true,
        "update_interval": 3600,
        "total_nodes": 50,
        "active_nodes": 45,
        "unique_new_nodes": 5,
        "status": "active",
        "last_update": "2024-01-15T10:00:00Z",
        "last_success": "2024-01-15T10:00:00Z",
        "error_count": 0,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
      }
    ]
  }
}
```

#### 2.2 创建订阅
```http
POST /subscriptions
Content-Type: application/json

{
  "name": "新订阅",
  "url": "https://example.com/subscribe",
  "user_agent": "clash",
  "auto_update": true,
  "update_interval": 3600,
  "node_pool_ids": [1, 2]
}
```

#### 2.3 更新订阅
```http
PUT /subscriptions/{subscription_id}
Content-Type: application/json

{
  "name": "更新的订阅名称",
  "auto_update": false,
  "update_interval": 7200
}
```

#### 2.4 手动更新订阅
```http
POST /subscriptions/{subscription_id}/update
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_fetched": 60,
    "valid_nodes": 55,
    "new_nodes": 8,
    "global_new_nodes": 3,
    "updated_nodes": 2,
    "removed_nodes": 1,
    "duration": 2350
  }
}
```

#### 2.5 获取订阅统计
```http
GET /subscriptions/{subscription_id}/stats
```

#### 2.6 获取订阅更新日志
```http
GET /subscriptions/{subscription_id}/logs
```

### 3. 节点池管理

#### 3.1 获取节点池列表
```http
GET /nodepools
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success", 
  "data": [
    {
      "id": 1,
      "name": "高速节点池",
      "description": "优质高速节点",
      "total_subscriptions": 3,
      "total_nodes": 150,
      "active_nodes": 128,
      "survival_rate": 85.33,
      "enabled": true,
      "priority": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-10T00:00:00Z"
    }
  ]
}
```

#### 3.2 创建节点池
```http
POST /nodepools
Content-Type: application/json

{
  "name": "新节点池",
  "description": "节点池描述",
  "enabled": true,
  "priority": 1
}
```

#### 3.3 获取节点池详情
```http
GET /nodepools/{pool_id}
```

#### 3.4 更新节点池
```http
PUT /nodepools/{pool_id}
Content-Type: application/json

{
  "name": "更新的节点池",
  "description": "新描述",
  "enabled": true,
  "priority": 2
}
```

#### 3.5 删除节点池
```http
DELETE /nodepools/{pool_id}
```

#### 3.6 关联订阅到节点池
```http
POST /nodepools/{pool_id}/subscriptions
Content-Type: application/json

{
  "subscription_ids": [1, 2, 3],
  "enabled": true,
  "priority": 1
}
```

### 4. 节点管理

#### 4.1 获取节点列表
```http
GET /nodes
```

**查询参数:**
- `page`: 页码，默认 1
- `size`: 每页大小，默认 20
- `subscription_id`: 按订阅过滤
- `node_pool_id`: 按节点池过滤
- `country`: 按国家过滤
- `protocol`: 按协议过滤
- `status`: 按状态过滤 (online/offline/testing/unknown)
- `sort`: 排序字段 (delay/upload_speed/download_speed/last_test)
- `order`: 排序方向 (asc/desc)

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "size": 20,
    "nodes": [
      {
        "id": 1,
        "subscription_id": 1,
        "node_pool_id": 1,
        "name": "🇭🇰 香港 HK01",
        "hash": "1a2b3c4d5e6f...",
        "server": "hk01.example.com",
        "port": 443,
        "protocol": "vmess",
        "country": "HK",
        "country_name": "Hong Kong",
        "city": "Hong Kong",
        "isp": "HKT",
        "delay": 50,
        "upload_speed": 104857600,
        "download_speed": 157286400,
        "loss_rate": 0.5,
        "status": "online",
        "last_test": "2024-01-15T10:00:00Z",
        "last_online": "2024-01-15T09:55:00Z",
        "continuous_failures": 0,
        "streaming_unlock": {
          "netflix": {"available": true, "region": "HK"},
          "youtube": {"available": true, "region": "HK"},
          "chatgpt": {"available": true}
        },
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
      }
    ]
  }
}
```

#### 4.2 获取节点详情
```http
GET /nodes/{node_id}
```

#### 4.3 测试单个节点
```http
POST /nodes/{node_id}/test
Content-Type: application/json

{
  "test_types": ["delay", "speed", "streaming"],
  "test_config": {
    "delay_url": "http://www.gstatic.com/generate_204",
    "timeout": 5000,
    "streaming_services": ["netflix", "youtube", "chatgpt"]
  }
}
```

#### 4.4 批量测试节点
```http
POST /nodes/batch-test
Content-Type: application/json

{
  "node_ids": [1, 2, 3],
  "test_types": ["delay", "speed"],
  "test_config": {
    "delay_url": "http://www.gstatic.com/generate_204",
    "timeout": 5000,
    "concurrent": 10
  }
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "test-task-12345",
    "total_nodes": 3,
    "status": "running"
  }
}
```

#### 4.5 获取测试任务状态
```http
GET /nodes/test-tasks/{task_id}
```

#### 4.6 获取节点测试历史
```http
GET /nodes/{node_id}/test-history
```

**查询参数:**
- `test_type`: 测试类型过滤
- `start_time`: 开始时间
- `end_time`: 结束时间
- `limit`: 返回数量限制

#### 4.7 智能节点选择
```http
GET /nodes/best-selection
```

**查询参数:**
- `node_pool_id`: 节点池限制
- `country`: 国家偏好
- `protocol`: 协议偏好
- `streaming`: 流媒体解锁需求
- `count`: 返回节点数量，默认1

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "node": {
        "id": 1,
        "name": "🇭🇰 香港 HK01",
        "delay": 45,
        "score": 95.5
      },
      "selection_reason": {
        "delay_score": 90,
        "speed_score": 95,
        "stability_score": 98,
        "streaming_score": 100
      }
    }
  ]
}
```

### 5. 统计和分析

#### 5.1 获取整体统计
```http
GET /stats/overview
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_subscriptions": 10,
    "active_subscriptions": 8,
    "total_node_pools": 5,
    "total_nodes": 500,
    "active_nodes": 450,
    "overall_survival_rate": 90.0,
    "total_tests_today": 1200,
    "successful_tests_today": 1080
  }
}
```

#### 5.2 获取地区分布统计
```http
GET /stats/geo-distribution
```

#### 5.3 获取协议分布统计
```http
GET /stats/protocol-distribution
```

#### 5.4 获取性能趋势
```http
GET /stats/performance-trend
```

**查询参数:**
- `period`: 时间周期 (hour/day/week/month)
- `node_pool_id`: 节点池过滤
- `country`: 国家过滤

### 6. 自动化任务

#### 6.1 获取自动更新任务状态
```http
GET /tasks/auto-update
```

#### 6.2 立即执行自动更新
```http
POST /tasks/auto-update/trigger
```

#### 6.3 获取定时测试任务状态
```http
GET /tasks/scheduled-test
```

### 7. 规则管理

#### 7.1 获取规则列表
```http
GET /rules
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "rule-1",
      "name": "直连规则",
      "type": "DOMAIN",
      "payload": "baidu.com",
      "proxy": "DIRECT",
      "enabled": true,
      "priority": 1
    },
    {
      "id": "rule-2", 
      "name": "代理规则",
      "type": "DOMAIN-SUFFIX",
      "payload": "google.com",
      "proxy": "PROXY",
      "enabled": true,
      "priority": 2
    }
  ]
}
```

#### 5.2 创建规则
```http
POST /rules
Content-Type: application/json

{
  "name": "新规则",
  "type": "DOMAIN-SUFFIX",
  "payload": "example.com", 
  "proxy": "PROXY",
  "enabled": true,
  "priority": 10
}
```

#### 5.3 更新规则
```http
PUT /rules/{rule_id}
```

#### 5.4 删除规则
```http
DELETE /rules/{rule_id}
```

### 6. 代理控制

#### 6.1 获取代理状态
```http
GET /proxy/status
```

#### 6.2 切换代理模式
```http
PUT /proxy/mode
Content-Type: application/json

{
  "mode": "rule"
}
```

**支持的模式:**
- `direct`: 直连模式
- `global`: 全局代理模式  
- `rule`: 规则模式

#### 6.3 选择代理节点
```http
PUT /proxy/select
Content-Type: application/json

{
  "node_id": "node-1"
}
```

#### 6.4 获取当前代理节点
```http
GET /proxy/current
```

### 7. 流量统计

#### 7.1 获取实时流量
```http
GET /traffic/realtime
```

**响应示例:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "upload_speed": 1024,
    "download_speed": 2048,
    "total_upload": 10485760,
    "total_download": 20971520,
    "connections": 15
  }
}
```

#### 7.2 获取历史流量统计
```http
GET /traffic/history
```

**查询参数:**
- `period`: 时间周期 (hour/day/week/month)
- `start_time`: 开始时间
- `end_time`: 结束时间

### 8. 连接管理

#### 8.1 获取活跃连接列表
```http
GET /connections
```

#### 8.2 关闭指定连接
```http
DELETE /connections/{connection_id}
```

#### 8.3 关闭所有连接
```http
DELETE /connections
```

### 9. 日志管理

#### 9.1 获取日志列表
```http
GET /logs
```

**查询参数:**
- `level`: 日志级别过滤
- `limit`: 返回数量限制
- `since`: 起始时间

### 10. 配置管理

#### 10.1 导出配置
```http
GET /config/export
```

#### 10.2 导入配置  
```http
POST /config/import
Content-Type: multipart/form-data

file: <config.yaml>
```

#### 10.3 重载配置
```http
POST /config/reload
```

## WebSocket 接口

### 实时数据推送
```
ws://localhost:9090/api/v1/ws
```

**订阅消息格式:**
```json
{
  "type": "subscribe",
  "topics": ["traffic", "connections", "logs"]
}
```

**推送消息格式:**
```json
{
  "type": "traffic",
  "timestamp": 1640995200,
  "data": {
    "upload_speed": 1024,
    "download_speed": 2048
  }
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0      | 成功 |
| 1001   | 参数错误 |
| 1002   | 认证失败 |
| 1003   | 权限不足 |
| 2001   | 节点池不存在 |
| 2002   | 节点不存在 |
| 2003   | 订阅链接无效 |
| 3001   | 配置错误 |
| 3002   | 代理启动失败 |
| 5000   | 服务器内部错误 |

## 接口限制

- API 请求频率限制：100 次/分钟
- 单个请求最大数据量：10MB
- WebSocket 连接数限制：10 个
- 批量操作最大数量：100 个

## SDK 支持

### JavaScript/TypeScript
```bash
npm install @prism/api-client
```

### Go
```bash
go get github.com/prism/go-client
```

### Java/Kotlin
```gradle
implementation 'com.prism:api-client:1.0.0'
```

这些 API 接口为客户端提供了完整的代理管理功能，确保各平台客户端能够统一、高效地与核心服务进行交互。