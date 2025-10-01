# I18n Package 国际化包

实现 `github.com/lazygophers/lrpc/middleware/xerror.I18n` 接口的多语言本地化系统。

## 功能特性

- 📁 **YAML 文件加载** - 从目录批量加载多语言 YAML 文件
- 🔄 **自动回退** - 支持语言回退机制（如 zh-CN -> zh -> en）
- 🗺️ **错误码映射** - 将 i18n 键（如 "common.success"）映射到错误码
- 🔒 **并发安全** - 使用读写锁保护共享数据
- 🌍 **多语言支持** - 支持任意数量的语言

## 架构设计

```
┌─────────────────┐
│  xerror.I18n   │ ◄── 接口定义
└────────┬────────┘
         │ implements
         │
┌────────▼────────┐
│   Localizer    │ ◄── 本地化器实现
│                 │
│ - messages map  │ ◄── locale -> code -> message
│ - fallback      │ ◄── 默认语言（en）
│ - mutex         │ ◄── 并发保护
└────────┬────────┘
         │ loads
         │
┌────────▼────────┐
│   YAML Files   │ ◄── 资源文件
│                 │
│ - en.yaml       │
│ - zh-CN.yaml    │
│ - fr.yaml       │
│ - ...           │
└─────────────────┘
```

## 核心类型

### Localizer

```go
type Localizer struct {
    messages map[string]map[int32]string  // locale -> code -> message
    mu       sync.RWMutex
    fallback string                       // 默认语言
}
```

### 错误码映射

```go
// codes.go
var keyToCodeMap = map[string]int32{
    "common.success": 0,
    "common.internal_error": 1000,
    // ...
}
```

## API 说明

### 创建本地化器

```go
// 创建新的本地化器实例
localizer := i18n.NewLocalizer()

// 或使用全局默认实例
i18n.LoadFromDir("./resource/localize")
```

### 加载资源文件

```go
// 从目录加载所有 YAML 文件
err := localizer.LoadFromDir("./resource/localize")

// 或加载单个文件
err := localizer.LoadFromFile("zh-CN", "./zh-CN.yaml")
```

### 本地化消息

```go
// 实现 xerror.I18n 接口
msg, ok := localizer.Localize(1000, "zh-CN")
// msg = "内部错误", ok = true

// 支持多个语言，按顺序尝试
msg, ok := localizer.Localize(1000, "zh-TW", "zh-CN", "en")
// 尝试：zh-TW -> zh -> zh-CN -> zh -> en

// 使用全局函数
msg, ok := i18n.Localize(1000, "fr")
```

### 设置回退语言

```go
// 默认回退语言是 "en"
i18n.SetFallback("zh-CN")
```

## YAML 文件格式

资源文件使用嵌套的 YAML 结构：

```yaml
# zh-CN.yaml
common:
  success: "操作成功"
  internal_error: "内部错误"

auth:
  unauthorized: "未授权"
  forbidden: "禁止访问"
```

文件会被扁平化为：
- `common.success` -> code 0 -> "操作成功"
- `common.internal_error` -> code 1000 -> "内部错误"
- `auth.unauthorized` -> code 2000 -> "未授权"

## 语言回退机制

本地化器支持多层回退：

```
用户请求: zh-TW
    ↓
尝试 zh-TW
    ↓ (未找到)
尝试 zh (去掉区域代码)
    ↓ (未找到)
尝试 en (默认语言)
    ↓
返回消息或 not found
```

示例：

```go
// messages 中只有 en 和 zh-CN
localizer.Localize(1000, "zh-TW")
// 尝试顺序：zh-TW -> zh -> en
// 如果 zh 存在，使用 zh
// 否则回退到 en
```

## 与 xerror 集成

### 1. 设置 I18n 实现

```go
import (
    "github.com/darabuchi/prism/pkg/i18n"
    "github.com/lazygophers/lrpc/middleware/xerror"
)

func init() {
    // 加载资源文件
    i18n.LoadFromDir("./resource/localize")

    // 设置 xerror 的 i18n 实现
    xerror.SetI18n(i18n.NewLocalizer())
}
```

### 2. 使用 xerror 创建本地化错误

```go
// xerror 会自动调用 Localize 方法
err := xerror.NewError(1000, "zh-CN")
// err.Code = 1000
// err.Msg = "内部错误"
```

## 添加新语言

### 1. 创建语言文件

创建 `resource/localize/{locale}.yaml`，例如 `ja.yaml`（日语）：

```yaml
common:
  success: "操作が成功しました"
  internal_error: "内部エラー"
```

### 2. 加载文件

```go
// 如果使用 LoadFromDir，新文件会自动加载
i18n.LoadFromDir("./resource/localize")
```

### 3. 使用新语言

```go
err := xerror.NewError(1000, "ja")
// err.Msg = "内部エラー"
```

## 并发安全

`Localizer` 使用读写锁保护内部数据：

- `LoadFromDir` / `LoadFromFile` - 写锁
- `Localize` - 读锁

可以安全地在多个 goroutine 中使用：

```go
// goroutine 1
msg1, _ := i18n.Localize(1000, "zh-CN")

// goroutine 2
msg2, _ := i18n.Localize(2000, "en")

// goroutine 3 (需要小心)
i18n.LoadFromDir("./new-locales")  // 可能会短暂阻塞读操作
```

## 性能考虑

1. **内存占用** - 所有翻译都加载到内存中
   - 对于大型应用，考虑按需加载

2. **查找性能** - O(1) map 查找
   - 非常高效，适合高并发场景

3. **初始化时间** - 启动时一次性加载
   - 推荐在 `init()` 或 `main()` 开始时加载

## 错误处理

### 加载失败

```go
if err := i18n.LoadFromDir("./locales"); err != nil {
    // 文件不存在、格式错误等
    log.Fatal(err)
}
```

### 翻译不存在

```go
msg, ok := i18n.Localize(9999, "zh-CN")
if !ok {
    // 错误码未定义或翻译缺失
    // 回退到默认语言或使用错误码本身
}
```

## 测试

### 单元测试示例

```go
func TestLocalizer(t *testing.T) {
    localizer := i18n.NewLocalizer()

    // 加载测试资源
    err := localizer.LoadFromFile("en", "./testdata/en.yaml")
    if err != nil {
        t.Fatal(err)
    }

    // 测试本地化
    msg, ok := localizer.Localize(1000, "en")
    if !ok {
        t.Error("localization failed")
    }

    if msg != "Internal error" {
        t.Errorf("expected 'Internal error', got '%s'", msg)
    }
}
```

## 扩展性

### 支持新的数据源

可以扩展 `Localizer` 支持其他数据源：

```go
// 从数据库加载
func (l *Localizer) LoadFromDB(db *sql.DB) error {
    // 实现数据库加载逻辑
}

// 从远程 API 加载
func (l *Localizer) LoadFromAPI(url string) error {
    // 实现 API 加载逻辑
}
```

### 动态重载

```go
// 监听文件变化，动态重载
func (l *Localizer) Watch(dirPath string) {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add(dirPath)

    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                l.LoadFromDir(dirPath)
            }
        }
    }
}
```

## 相关文档

- [错误码定义](../errcode/README.md)
- [多语言资源文件规范](../../resource/localize/README.md)
- [xerror 接口文档](https://github.com/lazygophers/lrpc/tree/main/middleware/xerror)
