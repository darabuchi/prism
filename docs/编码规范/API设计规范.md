# API 设计规范

## 1. HTTP 方法规范

### 1.1 强制使用 POST
- **所有 API 接口必须使用 POST 方法**
- 禁止使用 GET、PUT、DELETE、PATCH 等其他 HTTP 方法
- 统一使用 POST 简化 API 设计和客户端实现

#### 原因
- 简化接口设计，降低学习成本
- 避免 GET 请求的 URL 长度限制
- 提高参数传递的安全性（通过 Body 而非 URL）
- 便于统一的请求/响应处理

## 2. 路径设计规范

### 2.1 禁止路径变量
- **禁止在路径中使用变量**（如 `{id}`、`{uuid}` 等）
- 所有变量必须通过 Request Body 传递

#### 错误示例 ❌
```
POST /api/v1/users/{id}
POST /api/v1/subscriptions/{id}/update
POST /api/v1/nodes/{id}/test
```

#### 正确示例 ✅
```
POST /api/v1/users/detail
POST /api/v1/subscriptions/update
POST /api/v1/nodes/test
```

### 2.2 路径命名规范
路径应该体现资源和操作：`/api/v{version}/{resource}/{action}`

#### 资源名称（Resource）
- 使用复数形式
- 使用小写字母
- 多个单词使用连字符 `-` 分隔

示例：
- `subscriptions`
- `nodes`
- `connection-logs`
- `mitm-scripts`

#### 操作名称（Action）
标准操作名称：

| 操作 | 路径后缀 | 说明 | Request Body |
|------|---------|------|--------------|
| 创建 | `/create` 或直接 `/resource` | 创建资源 | 资源数据 |
| 列表查询 | `/list` | 获取列表 | 分页、过滤参数 |
| 详情查询 | `/detail` 或 `/get` | 获取详情 | `id` |
| 更新 | `/update` | 更新资源 | `id` + 更新数据 |
| 删除 | `/delete` | 删除资源 | `id` |
| 批量操作 | `/batch` | 批量处理 | `ids[]` + 操作数据 |
| 自定义操作 | `/custom-action` | 特定操作 | 操作相关数据 |

### 2.3 路径示例

```
# 订阅管理
POST /api/v1/subscriptions              # 创建订阅
POST /api/v1/subscriptions/list         # 查询列表
POST /api/v1/subscriptions/detail       # 查询详情
POST /api/v1/subscriptions/update       # 更新订阅
POST /api/v1/subscriptions/delete       # 删除订阅
POST /api/v1/subscriptions/force-update # 强制更新

# 节点管理
POST /api/v1/nodes/list                 # 查询列表
POST /api/v1/nodes/detail               # 查询详情
POST /api/v1/nodes/test                 # 触发测试
POST /api/v1/nodes/batch                # 批量操作
POST /api/v1/nodes/active               # 获取活动节点

# 配置管理
POST /api/v1/config/get                 # 获取配置
POST /api/v1/config/update              # 更新配置
POST /api/v1/config/rules/get           # 获取规则
POST /api/v1/config/rules/update        # 更新规则
```

## 3. 请求格式规范

### 3.1 Content-Type
- 必须使用 `application/json`
- 所有参数通过 JSON Body 传递

### 3.2 Request Body 结构
```json
{
  "id": "resource_id",           // 资源 ID（如需要）
  "page": 1,                     // 分页参数
  "page_size": 20,               // 每页数量
  "filter_field": "value",       // 过滤条件
  "data": {                      // 业务数据
    "field1": "value1",
    "field2": "value2"
  }
}
```

### 3.3 参数命名规范
- 使用 snake_case 命名（全小写，下划线分隔）
- 禁止使用 camelCase 或 PascalCase

#### 正确示例 ✅
```json
{
  "user_id": 123,
  "page_size": 20,
  "created_at": "2025-10-01T00:00:00Z",
  "is_enabled": true
}
```

#### 错误示例 ❌
```json
{
  "userId": 123,          // ❌ camelCase
  "PageSize": 20,         // ❌ PascalCase
  "createdAt": "...",     // ❌ camelCase
  "IsEnabled": true       // ❌ PascalCase
}
```

## 4. 响应格式规范

### 4.1 统一响应结构
所有 API 必须返回统一的响应格式：

```json
{
  "code": 0,              // 状态码：0=成功，非0=失败
  "message": "success",   // 提示信息
  "data": {}             // 响应数据
}
```

### 4.2 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "abc123",
    "name": "Example"
  }
}
```

### 4.3 列表响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": "item1",
        "name": "Item 1"
      }
    ]
  }
}
```

### 4.4 错误响应
```json
{
  "code": 1001,
  "message": "参数错误：id 不能为空",
  "data": null
}
```

### 4.5 状态码规范
| 范围 | 含义 | 示例 |
|------|------|------|
| 0 | 成功 | 0 |
| 1000-1999 | 客户端错误 | 1001=参数错误, 1002=未授权 |
| 2000-2999 | 服务端错误 | 2001=数据库错误, 2002=外部服务错误 |
| 3000-3999 | 业务错误 | 3001=订阅已存在, 3002=节点不可用 |

## 5. 分页规范

### 5.1 请求参数
```json
{
  "page": 1,              // 页码，从 1 开始
  "page_size": 20         // 每页数量，默认 20，最大 100
}
```

### 5.2 响应结构
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,         // 总记录数
    "page": 1,            // 当前页码
    "page_size": 20,      // 每页数量
    "items": []           // 数据列表
  }
}
```

## 6. 过滤和排序

### 6.1 过滤参数
直接在 Request Body 中传递过滤条件：

```json
{
  "page": 1,
  "page_size": 20,
  "enabled": true,        // 过滤：是否启用
  "type": "ss",          // 过滤：类型
  "alive": true          // 过滤：是否存活
}
```

### 6.2 排序参数
```json
{
  "page": 1,
  "page_size": 20,
  "order_by": "created_at",    // 排序字段
  "order_direction": "desc"     // 排序方向：asc/desc
}
```

## 7. 认证和授权

### 7.1 认证方式
- 使用 JWT Token 认证
- Token 通过 HTTP Header 传递：`Authorization: Bearer <token>`

### 7.2 请求头
```
POST /api/v1/subscriptions/list
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

## 8. 错误处理

### 8.1 参数验证错误
```json
{
  "code": 1001,
  "message": "参数验证失败",
  "data": {
    "errors": [
      {
        "field": "email",
        "message": "邮箱格式不正确"
      },
      {
        "field": "password",
        "message": "密码长度不能少于 8 位"
      }
    ]
  }
}
```

### 8.2 业务错误
```json
{
  "code": 3001,
  "message": "订阅已存在",
  "data": {
    "existing_id": "sub123"
  }
}
```

### 8.3 系统错误
```json
{
  "code": 2001,
  "message": "数据库连接失败",
  "data": null
}
```

## 9. 版本控制

### 9.1 版本号
- 在 URL 路径中包含版本号：`/api/v1/...`
- 主版本号变更时，向后不兼容的修改
- 次版本号在响应 Header 中返回（可选）

### 9.2 版本升级
- v1 → v2 时，保持 v1 接口可用一段时间（如 6 个月）
- 提前通知客户端接口废弃计划
- 通过响应头告知客户端使用的是旧版本

## 10. 性能优化

### 10.1 响应压缩
- 启用 gzip/brotli 压缩
- 减少响应体积，提高传输速度

### 10.2 缓存策略
- 对于只读接口，设置合理的缓存时间
- 使用 ETag 或 Last-Modified 实现条件请求

### 10.3 字段过滤
支持客户端指定需要的字段（可选）：

```json
{
  "id": "node123",
  "fields": ["id", "name", "score"]  // 只返回这些字段
}
```

## 11. 违规检查清单

在提交代码前，必须检查：
- [ ] 是否所有接口都使用 POST 方法
- [ ] 路径中是否包含变量（如 `{id}`）
- [ ] 是否所有参数都通过 Request Body 传递
- [ ] 参数命名是否使用 snake_case
- [ ] 响应格式是否符合统一结构
- [ ] 错误码是否在规定范围内
- [ ] 是否添加了适当的认证授权
- [ ] 是否考虑了分页和性能优化

## 12. API 文档示例

### 12.1 创建订阅
```
POST /api/v1/subscriptions
Content-Type: application/json
Authorization: Bearer <token>

Request Body:
{
  "url": "https://example.com/subscription",
  "title": "My Subscription",
  "proxy_type": "auto"
}

Response:
{
  "code": 0,
  "message": "订阅创建成功",
  "data": {
    "id": "sub_abc123",
    "url": "https://example.com/subscription",
    "title": "My Subscription",
    "created_at": "2025-10-01T12:00:00Z"
  }
}
```

### 12.2 查询订阅列表
```
POST /api/v1/subscriptions/list
Content-Type: application/json
Authorization: Bearer <token>

Request Body:
{
  "page": 1,
  "page_size": 20,
  "enabled": true
}

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": "sub_abc123",
        "title": "My Subscription",
        "enabled": true,
        "created_at": "2025-10-01T12:00:00Z"
      }
    ]
  }
}
```

### 12.3 更新订阅
```
POST /api/v1/subscriptions/update
Content-Type: application/json
Authorization: Bearer <token>

Request Body:
{
  "id": "sub_abc123",
  "title": "Updated Title",
  "enabled": false
}

Response:
{
  "code": 0,
  "message": "订阅更新成功",
  "data": {
    "id": "sub_abc123",
    "title": "Updated Title",
    "enabled": false,
    "updated_at": "2025-10-01T13:00:00Z"
  }
}
```

### 12.4 删除订阅
```
POST /api/v1/subscriptions/delete
Content-Type: application/json
Authorization: Bearer <token>

Request Body:
{
  "id": "sub_abc123"
}

Response:
{
  "code": 0,
  "message": "订阅删除成功",
  "data": null
}
```
