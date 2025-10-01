# 错误码生成工具

根据 `error_code.json` 自动生成 `error_code.gen.go` 文件。

## 功能

- 读取 `error_code.json` 错误码定义
- 验证错误码唯一性和格式
- 生成 Go 常量定义文件
- 生成错误码查询函数

## 使用方法

### 快速开始

```bash
# 构建并运行（生成错误码文件）
make

# 或者分步执行
make build   # 构建工具
make run     # 运行生成
```

### 其他命令

```bash
make install   # 安装到 GOPATH/bin
make validate  # 验证 JSON 格式
make clean     # 清理构建产物
make help      # 显示帮助
```

### 手动运行

```bash
# 构建工具
go build -o bin/gen-error-code .

# 运行工具
./bin/gen-error-code \
  --input ../../../error_code.json \
  --output ../../../error_code.gen.go
```

## 生成的文件

生成的 `error_code.gen.go` 包含：

1. **错误码常量**：所有错误码的 Go 常量定义
2. **错误码信息结构**：包含完整的错误信息
3. **查询函数**：
   - `GetErrorCodeInfo(code int32)` - 获取错误码详细信息
   - `IsValidErrorCode(code int32)` - 检查错误码是否有效

## 示例

### error_code.json

```json
{
  "version": "1.0.0",
  "description": "Prism 项目错误码注册表",
  "categories": [
    {
      "name": "common",
      "range": "0, 1000-1999",
      "description": "通用错误"
    }
  ],
  "errors": [
    {
      "code": 1000,
      "key": "common.internal_error",
      "category": "common",
      "description": "内部错误"
    }
  ]
}
```

### 生成的 error_code.gen.go

```go
package prism

const (
    // ErrInternalError 内部错误 (1000)
    ErrInternalError int32 = 1000
)

// GetErrorCodeInfo 获取错误码信息
func GetErrorCodeInfo(code int32) (ErrorCodeInfo, bool) {
    // ...
}
```

## 错误码规范

### 命名规则

- JSON key 格式: `category.name`（例如 `common.internal_error`）
- Go 常量名: `Err` + 大驼峰命名（例如 `ErrInternalError`）

### 分类范围

参考 `error_code.json` 中的 categories 定义：

- 0, 1000-1999: common (通用错误)
- 2000-2999: auth (认证授权)
- 3000-3999: resource (资源相关)
- 4000-4999: validation (输入验证)
- 5000-5999: service (服务错误)
- 6000-6999: queue (消息队列)
- 7000-7999: config (配置相关)
- 8000-8999: network (网络相关)

## 开发流程

1. 在 `error_code.json` 中添加新的错误码定义
2. 运行 `make` 或 `make run` 重新生成代码
3. 在业务代码中使用生成的常量：

```go
import "github.com/lazygophers/lrpc/middleware/xerror"

// 使用错误码
return xerror.New(ErrInternalError, "something went wrong")
return xerror.WrapError(err, ErrDatabaseError, "database operation failed")
```

## 注意事项

- ⚠️ 不要手动编辑 `error_code.gen.go`，所有修改应在 `error_code.json` 中进行
- ✅ 添加新错误码前，确保错误码和 key 不重复
- ✅ 错误码应在对应的范围内
- ✅ 运行 `make validate` 验证 JSON 格式
