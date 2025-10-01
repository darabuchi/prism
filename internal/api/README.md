# API 层

## 说明

HTTP API 实现，包括路由、中间件、处理器等。

## 文件

- `router.go`: 路由配置
- `request.go`: 请求结构定义
- `response.go`: 响应封装
- `validator.go`: 参数验证

## 子模块

### middleware/ - 中间件
- `auth.go`: 认证
- `cors.go`: 跨域
- `logger.go`: 日志
- `metrics.go`: 监控
- `recovery.go`: 恢复

### handler/ - 处理器
- `subscription.go`: 订阅接口
- `node.go`: 节点接口
- `route.go`: 路由接口
- `task.go`: 任务接口
- `health.go`: 健康检查

## API 设计规范

遵循 RESTful 风格：
- GET: 查询
- POST: 创建
- PUT: 更新
- DELETE: 删除

## 参考文档

- [API 设计规范](../../docs/编码规范/API设计规范.md)
