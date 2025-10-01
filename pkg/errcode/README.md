# Error Code 错误码管理

基于 `github.com/lazygophers/lrpc/middleware/xerror` 的统一错误码管理系统，支持多语言国际化。

## 特性

✨ **统一错误码** - 全局唯一的错误码标识
🌍 **多语言支持** - 支持 7 种语言（英语、简体中文、繁体中文、法语、俄语、西班牙语、阿拉伯语）
📊 **自动本地化** - 通过 i18n 系统自动返回对应语言的错误消息
🎯 **类型安全** - 强类型错误码定义
🔗 **错误链** - 支持 `errors.Is` 错误判断

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

## 支持的语言

- **en** - English (default)
- **zh-CN** - 简体中文
- **zh-TW** - 繁體中文
- **fr** - Français (法语)
- **ru** - Русский (俄语)
- **es** - Español (西班牙语)
- **ar** - العربية (阿拉伯语)

## 快速开始

### 1. 初始化错误码系统

```go
import "github.com/darabuchi/prism/pkg/errcode"

func main() {
    // 加载多语言资源文件
    if err := errcode.Init("./resource/localize"); err != nil {
        log.Fatal(err)
    }
}
```

### 2. 创建错误

```go
// 方式 1：创建带本地化消息的错误
err := errcode.New(errcode.ErrSubscriptionNotFound, "zh-CN")

// 方式 2：创建带自定义消息的错误
err := errcode.NewWithMsg(errcode.ErrDatabaseError, "连接超时")

// 方式 3：直接使用 xerror
import "github.com/lazygophers/lrpc/middleware/xerror"
err := xerror.NewError(errcode.ErrInvalidRequest, "zh-CN")
```

### 3. 错误检查

```go
import (
    "errors"
    "github.com/lazygophers/lrpc/middleware/xerror"
    "github.com/darabuchi/prism/pkg/errcode"
)

// 检查错误码
if xerror.CheckCode(err, errcode.ErrSubscriptionNotFound) {
    // 处理订阅不存在的情况
}

// 获取错误码
code := xerror.GetCode(err)

// 获取错误消息
msg := xerror.GetMsg(err)
```

### 4. HTTP API 响应

```go
import "github.com/lazygophers/lrpc/middleware/xerror"

func GetSubscription(c *gin.Context) {
    sub, err := service.GetSubscription(id)
    if err != nil {
        // 获取客户端语言
        lang := c.GetHeader("Accept-Language") // 如 "zh-CN"

        // 创建本地化错误
        if xerror.CheckCode(err, errcode.ErrSubscriptionNotFound) {
            apiErr := errcode.New(errcode.ErrSubscriptionNotFound, lang)
            c.JSON(404, gin.H{
                "code": apiErr.Code,
                "msg":  apiErr.Msg,
            })
            return
        }

        // 内部错误不暴露详情
        c.JSON(500, gin.H{
            "code": errcode.ErrInternal,
            "msg":  errcode.New(errcode.ErrInternal, lang).Msg,
        })
        return
    }

    c.JSON(200, sub)
}
```

## 使用场景

### 业务逻辑层

```go
import "github.com/lazygophers/lrpc/middleware/xerror"

func (s *SubscriptionService) CreateSubscription(req *CreateRequest) error {
    // 验证 URL
    if req.URL == "" {
        return xerror.NewError(errcode.ErrURLRequired, "zh-CN")
    }

    // 检查是否已存在
    exists, err := s.repo.Exists(req.URL)
    if err != nil {
        // 包装底层错误
        return xerror.NewErrorWithMsg(errcode.ErrDatabaseError, err.Error())
    }

    if exists {
        return xerror.NewError(errcode.ErrSubscriptionExists, "zh-CN")
    }

    // 创建订阅
    if err := s.repo.Create(req); err != nil {
        return xerror.NewErrorWithMsg(errcode.ErrDatabaseError, err.Error())
    }

    return nil
}
```

### 多语言支持

```go
// 根据用户偏好返回不同语言的错误消息
func HandleError(err error, userLang string) string {
    code := xerror.GetCode(err)
    if code == -1 {
        return err.Error()
    }

    // 创建本地化错误
    localizedErr := errcode.New(code, userLang)
    return localizedErr.Msg
}

// 使用示例
msg1 := HandleError(err, "en")      // "Subscription not found"
msg2 := HandleError(err, "zh-CN")   // "订阅不存在"
msg3 := HandleError(err, "fr")      // "Abonnement introuvable"
```

## 添加新错误码

### 1. 在 `errcode.go` 中定义错误码

```go
const (
    // 自定义错误 (9000-9999)
    ErrCustomError int32 = 9000
)
```

### 2. 在 `pkg/i18n/codes.go` 中添加映射

```go
var keyToCodeMap = map[string]int32{
    // ... 现有映射
    "custom.error": 9000,
}
```

### 3. 在资源文件中添加翻译

**resource/localize/zh-CN.yaml:**
```yaml
custom:
  error: "自定义错误消息"
```

**resource/localize/en.yaml:**
```yaml
custom:
  error: "Custom error message"
```

对其他语言文件重复此操作。

## API 参考

### 初始化

- `Init(dirPath string) error` - 初始化错误码系统，加载多语言资源

### 创建错误

- `New(code int32, langs ...string) *xerror.Error` - 创建带本地化消息的错误
- `NewWithMsg(code int32, msg string) *xerror.Error` - 创建带自定义消息的错误

### 错误检查（来自 xerror）

- `xerror.GetCode(err error) int32` - 获取错误码
- `xerror.GetMsg(err error) string` - 获取错误消息
- `xerror.CheckCode(err error, code int32) bool` - 检查错误码
- `xerror.Is(err1, err2 error) bool` - 判断错误类型

## 最佳实践

1. **统一初始化** - 在应用启动时调用 `errcode.Init()` 加载资源
2. **语言检测** - 从 HTTP Header `Accept-Language` 获取用户语言偏好
3. **错误码一致性** - 同一类错误使用统一的错误码
4. **不暴露内部错误** - 对外 API 返回标准错误码，日志记录详细错误
5. **回退机制** - i18n 系统会自动回退到英语（默认语言）
6. **错误包装** - 使用 `NewErrorWithMsg` 包装底层错误，保留详细信息用于日志

## 相关文档

- [i18n 资源文件格式](../../resource/localize/README.md)
- [xerror 文档](https://github.com/lazygophers/lrpc/tree/main/middleware/xerror)
- [多语言支持规范](../../resource/localize/README.md)
