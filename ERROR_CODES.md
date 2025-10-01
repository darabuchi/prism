# Error Codes 错误码说明

本项目使用 `github.com/lazygophers/lrpc/middleware/xerror` 作为统一的错误处理机制，配合多语言 i18n 系统提供国际化的错误消息。

## 错误码注册表

所有错误码定义在 [`error_code.json`](error_code.json) 文件中，这是一个机器可读的错误码注册表，可用于：
- 文档生成
- API 规范
- 客户端错误码同步
- 代码生成工具

## 错误码范围

| 范围 | 分类 | 说明 |
|------|------|------|
| 0 | Success | 操作成功 |
| 1000-1999 | Common | 通用错误（内部错误、请求无效等） |
| 2000-2999 | Auth | 认证和授权错误 |
| 3000-3999 | Resource | 资源相关错误（订阅、节点、路由） |
| 4000-4999 | Validation | 输入验证错误 |
| 5000-5999 | Service | 服务错误（数据库、缓存） |
| 6000-6999 | Queue | 消息队列错误 |
| 7000-7999 | Config | 配置错误 |
| 8000-8999 | Network | 网络错误（代理、DNS） |

## 支持的语言

系统内置 **7 种语言**支持：

| 代码 | 语言 | 示例 |
|------|------|------|
| `en` | English | "Internal error" |
| `zh-CN` | 简体中文 | "内部错误" |
| `zh-TW` | 繁體中文 | "內部錯誤" |
| `fr` | Français | "Erreur interne" |
| `ru` | Русский | "Внутренняя ошибка" |
| `es` | Español | "Error interno" |
| `ar` | العربية | "خطأ داخلي" |

## 使用方法

### 1. 初始化系统

```go
import (
    "github.com/darabuchi/prism/pkg/i18n"
    "github.com/lazygophers/lrpc/middleware/xerror"
)

func main() {
    // 加载多语言资源
    if err := i18n.LoadFromDir("./resource/localize"); err != nil {
        log.Fatal(err)
    }

    // 设置 xerror 的 i18n 实现
    xerror.SetI18n(i18n.NewLocalizer())
}
```

### 2. 创建错误

```go
// 创建带本地化消息的错误
err := xerror.NewError(1000, "zh-CN")
// err.Code = 1000
// err.Msg = "内部错误"

// 创建带自定义消息的错误
err := xerror.NewErrorWithMsg(5000, "数据库连接超时")
```

### 3. 检查错误

```go
// 获取错误码
code := xerror.GetCode(err)  // 1000

// 获取错误消息
msg := xerror.GetMsg(err)    // "内部错误"

// 检查特定错误码
if xerror.CheckCode(err, 3000) {
    // 处理订阅不存在的情况
}
```

### 4. HTTP API 集成

```go
func HandleRequest(c *gin.Context) {
    // 从请求头获取语言偏好
    lang := c.GetHeader("Accept-Language")
    if lang == "" {
        lang = "en"  // 默认英语
    }

    result, err := service.DoSomething()
    if err != nil {
        code := xerror.GetCode(err)
        if code == -1 {
            // 不是 xerror，返回通用错误
            c.JSON(500, gin.H{
                "code": 1000,
                "msg":  xerror.NewError(1000, lang).Msg,
            })
            return
        }

        // 返回本地化的错误
        localErr := xerror.NewError(code, lang)
        c.JSON(getHTTPStatus(code), gin.H{
            "code": localErr.Code,
            "msg":  localErr.Msg,
        })
        return
    }

    c.JSON(200, result)
}

func getHTTPStatus(code int32) int {
    // 根据错误码返回适当的 HTTP 状态码
    switch {
    case code >= 2000 && code < 3000:
        return 401  // 认证错误
    case code >= 4000 && code < 5000:
        return 400  // 验证错误
    case code == 1002:
        return 404  // 未找到
    default:
        return 500  // 内部错误
    }
}
```

## 添加新错误码

### 1. 更新 error_code.json

```json
{
  "code": 9000,
  "key": "custom.error",
  "category": "custom",
  "description": "自定义错误"
}
```

### 2. 更新 pkg/i18n/codes.go

```go
var keyToCodeMap = map[string]int32{
    // ... 现有映射
    "custom.error": 9000,
}
```

### 3. 更新所有语言文件

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

对其他 5 种语言重复此操作。

### 4. 使用新错误码

```go
err := xerror.NewError(9000, "zh-CN")
// err.Msg = "自定义错误消息"
```

## 最佳实践

1. **一致性**: 同一类错误使用统一的错误码
2. **不暴露细节**: 对外 API 返回标准错误码，日志记录详细错误
3. **语言检测**: 从 HTTP Header `Accept-Language` 获取用户语言
4. **错误包装**: 使用 `NewErrorWithMsg` 保留底层错误详情
5. **文档同步**: 添加新错误码时同步更新 `error_code.json`
6. **回退机制**: i18n 系统会自动回退到英语（默认语言）

## 相关文档

- [i18n 包文档](pkg/i18n/README.md) - 多语言实现细节
- [多语言资源规范](resource/localize/README.md) - YAML 文件格式
- [xerror 文档](https://github.com/lazygophers/lrpc/tree/main/middleware/xerror) - 错误处理库
