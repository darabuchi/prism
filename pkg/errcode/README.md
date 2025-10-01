# Error Code 错误码管理

统一的错误码管理系统，支持多语言国际化。

## 特性

✨ **统一错误码** - 全局唯一的错误码标识
🌍 **多语言支持** - 基于 i18n 资源文件的本地化
📊 **HTTP 状态码** - 每个错误码对应标准 HTTP 状态码
🔗 **错误链** - 基于 xerrors 的错误包装和原因追踪
📝 **详细信息** - 支持添加错误详情和上下文
🎯 **类型安全** - 强类型错误码定义
🔍 **错误格式化** - 支持 %+v 格式化输出完整错误链

## 错误码范围

| 范围 | 模块 | 说明 |
|------|------|------|
| 0 | 成功 | 操作成功 |
| 1000-1999 | 通用错误 | 基础错误类型 |
| 2000-2999 | 认证授权 | 认证和授权相关错误 |
| 3000-3999 | 资源错误 | 订阅、节点、路由等资源错误 |
| 4000-4999 | 验证错误 | 输入验证相关错误 |
| 5000-5999 | 服务错误 | 数据库、缓存等服务错误 |
| 6000-6999 | 队列错误 | 消息队列相关错误 |
| 7000-7999 | 配置错误 | 配置文件和设置错误 |
| 8000-8999 | 网络错误 | 网络、代理、DNS 等错误 |

## 快速开始

### 1. 初始化本地化器

```go
import "github.com/darabuchi/prism/pkg/errcode"

func main() {
    // 加载本地化文件
    err := errcode.LoadLocalizations("./resource/localize")
    if err != nil {
        log.Fatal(err)
    }
}
```

### 2. 创建错误

```go
// 方式 1：直接创建本地化错误
err := errcode.NewLocalizedError("zh-CN", errcode.ErrSubscriptionNotFound)

// 方式 2：创建带格式化的错误
err := errcode.NewLocalizedErrorf("zh-CN", errcode.ErrFieldRequired, "username")

// 方式 3：包装现有错误
err := errcode.WrapLocalizedError("zh-CN", errcode.ErrDatabaseError, originalErr)

// 方式 4：使用错误码直接创建
err := errcode.ErrSubscriptionNotFound.New("订阅不存在")
```

### 3. 添加错误详情

```go
err := errcode.NewLocalizedError("zh-CN", errcode.ErrSubscriptionNotFound)
err.WithDetail("subscription_id", "12345")
err.WithDetail("user_id", "67890")

// 或批量添加
err.WithDetails(map[string]interface{}{
    "subscription_id": "12345",
    "user_id": "67890",
})
```

### 4. 错误检查

```go
// 检查错误码（支持错误链）
if errcode.Is(err, errcode.ErrSubscriptionNotFound) {
    // 处理订阅不存在的情况
}

// 获取错误码（从错误链中提取）
code := errcode.GetCode(err)

// 获取 HTTP 状态码（从错误链中提取）
httpStatus := errcode.GetHTTPStatus(err)
```

### 5. 错误链和格式化

```go
// 创建错误链
err1 := errors.New("database connection error")
err2 := errcode.WrapLocalizedError("zh-CN", errcode.ErrDatabaseError, err1)
err3 := errcode.WrapLocalizedError("zh-CN", errcode.ErrInternal, err2)

// 简单输出
fmt.Printf("%s\n", err3)
// 输出: [1000] 内部错误: [5000] 数据库错误: database connection error

// 详细输出（包含完整错误链）
fmt.Printf("%+v\n", err3)
// 输出:
// [1000] 内部错误
// Caused by: [5000] 数据库错误
// Caused by: database connection error

// 检查错误链中的任意错误码
if errcode.Is(err3, errcode.ErrDatabaseError) {
    // 即使 err3 外层是 ErrInternal，也能识别内层的 ErrDatabaseError
    fmt.Println("检测到数据库错误")
}
```

## 使用场景

### HTTP API 响应

```go
func GetSubscription(c *gin.Context) {
    sub, err := service.GetSubscription(id)
    if err != nil {
        if errcode.Is(err, errcode.ErrSubscriptionNotFound) {
            c.JSON(err.HTTPStatus(), gin.H{
                "code": err.Code(),
                "message": err.Message,
                "details": err.Details,
            })
            return
        }

        // 处理其他错误
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": errcode.ErrInternal.Code,
            "message": "Internal server error",
        })
        return
    }

    c.JSON(http.StatusOK, sub)
}
```

### 业务逻辑层

```go
func (s *SubscriptionService) CreateSubscription(req *CreateRequest) error {
    // 验证 URL
    if req.URL == "" {
        return errcode.NewLocalizedError("zh-CN", errcode.ErrURLRequired)
    }

    // 检查是否已存在
    exists, err := s.repo.Exists(req.URL)
    if err != nil {
        return errcode.WrapLocalizedError("zh-CN", errcode.ErrDatabaseError, err)
    }

    if exists {
        return errcode.NewLocalizedError("zh-CN", errcode.ErrSubscriptionExists)
    }

    // 创建订阅
    if err := s.repo.Create(req); err != nil {
        return errcode.WrapLocalizedError("zh-CN", errcode.ErrDatabaseError, err).
            WithDetail("url", req.URL)
    }

    return nil
}
```

### 队列消息处理

```go
func ProcessMessage(msg *queue.Message) error {
    // 处理消息
    if err := process(msg.Payload); err != nil {
        // 根据错误类型决定是否重试
        if errcode.Is(err, errcode.ErrTimeout) {
            // 临时错误，可以重试
            return queue.Retry()
        }

        if errcode.Is(err, errcode.ErrValidationFailed) {
            // 永久错误，不重试
            return err
        }

        // 其他错误，使用默认策略
        return err
    }

    return nil
}
```

## 自定义错误码

如果需要添加新的错误码，按以下步骤：

### 1. 在 `codes.go` 中定义错误码

```go
// ErrCustomError 自定义错误
ErrCustomError = &ErrCode{
    Code:     9000,
    Key:      "custom.error",
    HTTPCode: http.StatusBadRequest,
}
```

### 2. 在资源文件中添加翻译

**zh-CN.yaml:**
```yaml
custom:
  error: "自定义错误消息"
```

**en-US.yaml:**
```yaml
custom:
  error: "Custom error message"
```

## API 参考

### 核心类型

#### ErrCode
```go
type ErrCode struct {
    Code     int    // 错误码
    Key      string // i18n 键名
    HTTPCode int    // HTTP 状态码
}
```

#### CodedError
```go
type CodedError struct {
    ErrCode *ErrCode               // 错误码
    Message string                 // 本地化消息
    Cause   error                  // 原始错误
    Details map[string]interface{} // 错误详情
}
```

### 主要方法

#### 创建错误
- `NewLocalizedError(locale string, code *ErrCode) *CodedError`
- `NewLocalizedErrorf(locale string, code *ErrCode, args ...interface{}) *CodedError`
- `WrapLocalizedError(locale string, code *ErrCode, err error) *CodedError`
- `WrapLocalizedErrorf(locale string, code *ErrCode, err error, args ...interface{}) *CodedError`

#### 错误方法
- `WithCause(err error) *CodedError` - 添加原始错误
- `WithDetail(key string, value interface{}) *CodedError` - 添加详情
- `WithDetails(details map[string]interface{}) *CodedError` - 批量添加详情
- `Code() int` - 获取错误码
- `HTTPStatus() int` - 获取 HTTP 状态码

#### 错误检查
- `Is(err error, code *ErrCode) bool` - 判断错误类型
- `GetCode(err error) int` - 提取错误码
- `GetHTTPStatus(err error) int` - 提取 HTTP 状态码

### 本地化方法

- `LoadLocalizations(dirPath string) error` - 加载本地化文件
- `T(locale, key string) string` - 翻译
- `Tf(locale, key string, args ...interface{}) string` - 格式化翻译

## 最佳实践

1. **统一错误处理** - 在 API 层统一处理错误响应
2. **使用错误链** - 使用 `Wrap` 保留原始错误信息，支持错误追踪
3. **添加上下文** - 使用 `WithDetail` 添加调试信息
4. **区分错误类型** - 使用 `Is` 判断具体错误类型（支持错误链查找）
5. **本地化消息** - 始终使用本地化函数创建面向用户的错误
6. **详细日志** - 使用 `%+v` 格式化记录完整的错误链和详情
7. **错误传播** - 在跨层传递错误时，使用 `Wrap` 添加上下文而非创建新错误
8. **避免丢失信息** - 不要只返回错误消息字符串，保留完整的错误对象

## 示例项目

查看 `examples/` 目录获取完整的使用示例。

## 相关文档

- [i18n 资源文件格式](../../resource/localize/README.md)
- [API 错误处理规范](../../docs/api-error-handling.md)
