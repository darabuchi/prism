# 多语言国际化 (i18n) 规范

本文档定义项目的多语言国际化规范和最佳实践。

## 目录结构

```
resource/localize/
├── README.md         # 本规范文档
├── zh-CN.yaml        # 简体中文（默认语言）
├── en-US.yaml        # 美式英语
├── zh-TW.yaml        # 繁体中文（可选）
├── ja.yaml           # 日语（可选）
├── ko.yaml           # 韩语（可选）
├── pt-BR.yaml        # 巴西葡萄牙语（可选）
└── ...               # 其他语言
```

## 语言代码规范

本项目采用 **ISO 639-1**（语言代码）和 **ISO 3166-1**（国家/地区代码）组合的标准。

参考规范：[Lingoes 语言代码表](http://www.lingoes.net/zh/translator/langcode.htm)

### 语言代码格式

支持两种格式：

1. **纯语言代码**（ISO 639-1）：`{language}`
   - 适用于不需要区分地区变体的语言
   - 例如：`ja`（日语）、`ko`（韩语）

2. **语言-地区代码**（ISO 639-1 + ISO 3166-1）：`{language}-{REGION}`
   - 适用于需要区分地区变体的语言
   - 例如：`zh-CN`（简体中文）、`en-US`（美式英语）

### 常用语言代码表

#### 中文变体

| 语言代码 | 语言名称 | ISO 639-1 | ISO 3166-1 | 说明 |
|---------|---------|-----------|------------|------|
| `zh` | 中文 | zh | - | 通用中文（不推荐，建议明确指定地区） |
| `zh-CN` | 简体中文 | zh | CN | 中国大陆（**推荐作为默认**） |
| `zh-TW` | 繁体中文 | zh | TW | 中国台湾 |
| `zh-HK` | 繁体中文 | zh | HK | 中国香港 |
| `zh-SG` | 简体中文 | zh | SG | 新加坡 |

#### 英语变体

| 语言代码 | 语言名称 | ISO 639-1 | ISO 3166-1 | 说明 |
|---------|---------|-----------|------------|------|
| `en` | 英语 | en | - | 通用英语 |
| `en-US` | 美国英语 | en | US | 美国（**推荐**） |
| `en-GB` | 英国英语 | en | GB | 英国 |
| `en-AU` | 澳大利亚英语 | en | AU | 澳大利亚 |
| `en-CA` | 加拿大英语 | en | CA | 加拿大 |

#### 其他常用语言

| 语言代码 | 语言名称 | ISO 639-1 | ISO 3166-1 | 说明 |
|---------|---------|-----------|------------|------|
| `ja` | 日语 | ja | - | 日本（单一语言，无需地区代码） |
| `ko` | 韩语 | ko | - | 韩国（单一语言，无需地区代码） |
| `fr` | 法语 | fr | - | 法国（也可用 `fr-FR`） |
| `de` | 德语 | de | - | 德国（也可用 `de-DE`） |
| `es` | 西班牙语 | es | - | 西班牙（也可用 `es-ES`） |
| `pt` | 葡萄牙语 | pt | - | 葡萄牙 |
| `pt-BR` | 巴西葡萄牙语 | pt | BR | 巴西 |
| `ru` | 俄语 | ru | - | 俄罗斯 |
| `ar` | 阿拉伯语 | ar | - | 阿拉伯语 |
| `hi` | 印地语 | hi | - | 印度 |
| `th` | 泰语 | th | - | 泰国 |
| `vi` | 越南语 | vi | - | 越南 |
| `id` | 印尼语 | id | - | 印度尼西亚 |
| `ms` | 马来语 | ms | - | 马来西亚 |
| `tr` | 土耳其语 | tr | - | 土耳其 |
| `pl` | 波兰语 | pl | - | 波兰 |
| `nl` | 荷兰语 | nl | - | 荷兰 |
| `it` | 意大利语 | it | - | 意大利 |
| `sv` | 瑞典语 | sv | - | 瑞典 |
| `da` | 丹麦语 | da | - | 丹麦 |
| `no` | 挪威语 | no | - | 挪威 |
| `fi` | 芬兰语 | fi | - | 芬兰 |
| `cs` | 捷克语 | cs | - | 捷克 |
| `hu` | 匈牙利语 | hu | - | 匈牙利 |
| `ro` | 罗马尼亚语 | ro | - | 罗马尼亚 |
| `uk` | 乌克兰语 | uk | - | 乌克兰 |
| `el` | 希腊语 | el | - | 希腊 |
| `he` | 希伯来语 | he | - | 以色列 |

### 文件命名规范

**基本规则**：

- **必须**使用 `{language}.yaml` 或 `{language-REGION}.yaml` 格式
- **语言代码**（ISO 639-1）使用**小写字母**（2个字符）
- **地区代码**（ISO 3166-1）使用**大写字母**（2个字符）
- 使用**连字符** `-` 分隔语言和地区
- **不允许**使用下划线 `_` 或其他分隔符

✅ **正确示例**：
```
zh-CN.yaml          # 简体中文（推荐）
zh-TW.yaml          # 繁体中文
en-US.yaml          # 美式英语（推荐）
en-GB.yaml          # 英式英语
ja.yaml             # 日语（无地区代码）
ko.yaml             # 韩语（无地区代码）
pt-BR.yaml          # 巴西葡萄牙语
```

❌ **错误示例**：
```
zh_CN.yaml          # ❌ 使用了下划线
zh-cn.yaml          # ❌ 地区代码应大写
ZH-CN.yaml          # ❌ 语言代码应小写
zhcn.yaml           # ❌ 缺少分隔符
chinese.yaml        # ❌ 应使用 ISO 639-1 代码
zh_Hans.yaml        # ❌ 不符合标准格式
```

### 语言代码选择建议

1. **中文**：明确使用 `zh-CN` 或 `zh-TW`，不要使用 `zh`
2. **英语**：推荐使用 `en-US`，如需区分英式英语可使用 `en-GB`
3. **单一语言**：如果语言没有明显的地区差异（如日语、韩语），可以只使用语言代码（`ja`, `ko`）
4. **有地区差异的语言**：必须使用完整的语言-地区代码（如 `pt-BR` 巴西葡萄牙语）

### 默认语言设置

- **项目默认语言**：`zh-CN`（简体中文）
- **备用语言**：`en-US`（美式英语）
- **Fallback 顺序**：请求的语言 → 默认语言 → 键名本身

## YAML 文件结构

### 基本格式

```yaml
# 注释说明
category:
  subcategory:
    key: "翻译文本"
    another_key: "另一个翻译"
```

### 层级组织规范

使用 **点分隔符** 表示层级关系，最终生成的 key 为 `category.subcategory.key`

#### 推荐的分类结构

```yaml
# ==================== 通用消息 ====================
common:
  success: "操作成功"
  error: "操作失败"

# ==================== 认证和授权 ====================
auth:
  login_success: "登录成功"
  logout_success: "退出登录"

# ==================== 业务模块 ====================
subscription:
  not_found: "订阅不存在"
  created: "订阅创建成功"

node:
  not_found: "节点不存在"
  test_failed: "节点测试失败"

# ==================== 验证错误 ====================
validation:
  required: "字段不能为空"
  invalid_format: "格式无效"

# ==================== 业务消息 ====================
messages:
  server_started: "服务器已启动"
  operation_completed: "操作已完成"
```

### 分类命名规范

| 分类 | 用途 | 示例 |
|------|------|------|
| `common` | 通用消息和错误 | `common.success`, `common.error` |
| `auth` | 认证和授权 | `auth.login_required`, `auth.forbidden` |
| `validation` | 输入验证 | `validation.email_invalid`, `validation.required` |
| `subscription` | 订阅相关 | `subscription.not_found`, `subscription.created` |
| `node` | 节点相关 | `node.unavailable`, `node.test_failed` |
| `route` | 路由相关 | `route.not_found`, `route.invalid` |
| `queue` | 队列相关 | `queue.full`, `queue.timeout` |
| `config` | 配置相关 | `config.invalid`, `config.parse_failed` |
| `proxy` | 代理相关 | `proxy.connection_failed`, `proxy.timeout` |
| `network` | 网络相关 | `network.unreachable`, `network.dns_failed` |
| `service` | 服务错误 | `service.database_error`, `service.cache_error` |
| `messages` | 业务消息 | `messages.server_started`, `messages.task_completed` |

## 占位符规范

### 使用 Printf 风格占位符

```yaml
# 单个占位符
validation:
  field_required: "字段 %s 不能为空"
  field_too_long: "字段 %s 过长，最大长度为 %d"

messages:
  server_started: "服务器已启动，监听端口 %d"
  items_found: "找到 %d 个项目"
```

### 占位符类型

| 占位符 | 类型 | 示例 |
|--------|------|------|
| `%s` | 字符串 | `"用户 %s 已登录"` |
| `%d` | 整数 | `"找到 %d 条记录"` |
| `%f` | 浮点数 | `"进度 %.2f%%"` |
| `%v` | 任意类型 | `"值为 %v"` |

### 占位符顺序

对于多个占位符，**必须**保持顺序一致：

```yaml
# zh-CN.yaml
validation:
  range_error: "字段 %s 的值必须在 %d 到 %d 之间"

# en-US.yaml
validation:
  range_error: "Field %s must be between %d and %d"
```

使用时：
```go
errcode.Tf("zh-CN", "validation.range_error", "age", 18, 65)
// 输出: 字段 age 的值必须在 18 到 65 之间
```

### 命名占位符（不推荐）

如果确实需要命名占位符，使用注释说明：

```yaml
# %[1]s: 用户名, %[2]d: 年龄
messages:
  user_info: "用户 %[1]s 的年龄是 %[2]d 岁"
```

## 翻译键名规范

### 命名约定

1. **使用小写字母和下划线**
   - ✅ `subscription_not_found`
   - ❌ `subscriptionNotFound`（驼峰命名）
   - ❌ `Subscription-Not-Found`（连字符）

2. **使用描述性名称**
   - ✅ `validation.email_invalid`
   - ❌ `validation.err1`

3. **保持简洁但明确**
   - ✅ `auth.unauthorized`
   - ❌ `auth.user_is_not_authorized_to_access_this_resource`

4. **使用动词-名词结构（操作消息）**
   - ✅ `messages.subscription_created`
   - ✅ `messages.node_test_started`

5. **使用名词-形容词结构（状态消息）**
   - ✅ `node.unavailable`
   - ✅ `queue.full`

## 翻译内容规范

### 文本规范

1. **保持简洁明了**
   ```yaml
   # ✅ 好
   auth:
     unauthorized: "未授权，请先登录"

   # ❌ 不好（过于冗长）
   auth:
     unauthorized: "您当前未登录系统，无法访问此资源，请先登录后再试"
   ```

2. **使用正式语气**
   ```yaml
   # ✅ 好
   common:
     error: "操作失败"

   # ❌ 不好（过于口语化）
   common:
     error: "哎呀，出错了"
   ```

3. **避免技术术语**
   ```yaml
   # ✅ 好（面向用户）
   service:
     database_error: "数据库错误"

   # ❌ 不好（过于技术化）
   service:
     database_error: "SQL query execution failed with errno 1064"
   ```

4. **保持一致的术语**
   - 在整个项目中统一使用相同的术语翻译
   - 维护术语表（见下文）

### 标点符号规范

| 语言 | 规范 |
|------|------|
| 中文 | 使用中文标点（。，！？：；） |
| 英文 | 使用英文标点（., !, ?, :, ;） |
| 日文 | 使用日文标点（。、！？：；） |

```yaml
# zh-CN.yaml
common:
  success: "操作成功"  # 无标点或使用中文句号

# en-US.yaml
common:
  success: "Operation succeeded"  # 无标点或使用英文句号
```

## 术语表

维护统一的术语翻译表：

| 英文 | 简体中文 | 繁体中文 | 日语 |
|------|---------|---------|------|
| Subscription | 订阅 | 訂閱 | サブスクリプション |
| Node | 节点 | 節點 | ノード |
| Route | 路由 | 路由 | ルート |
| Proxy | 代理 | 代理 | プロキシ |
| Queue | 队列 | 佇列 | キュー |
| Configuration | 配置 | 設定 | 設定 |
| Authentication | 认证 | 認證 | 認証 |
| Authorization | 授权 | 授權 | 認可 |
| Validation | 验证 | 驗證 | 検証 |

## 添加新语言

### 步骤

1. **创建新的 YAML 文件**
   ```bash
   # 根据需要选择合适的语言代码格式
   # 示例1: 日语（无需地区代码）
   cp resource/localize/en-US.yaml resource/localize/ja.yaml

   # 示例2: 巴西葡萄牙语（需要地区代码）
   cp resource/localize/en-US.yaml resource/localize/pt-BR.yaml
   ```

2. **翻译所有键值**
   - 保持键名不变
   - 只翻译值（value）部分
   - 保持占位符位置和顺序

3. **验证文件格式**
   ```bash
   # 使用 YAML 验证工具
   yamllint resource/localize/ja.yaml
   ```

4. **测试翻译**
   ```go
   // 在代码中测试
   msg := errcode.T("ja", "common.success")
   fmt.Println(msg) // 输出日语翻译

   // 测试带地区代码的语言
   msg = errcode.T("pt-BR", "common.success")
   fmt.Println(msg) // 输出巴西葡萄牙语翻译
   ```

5. **更新文档**
   - 在本 README 中添加新语言的说明
   - 更新支持的语言列表
   - 如果是新的语言变体，更新常用语言代码表

## 使用示例

### 在代码中使用

```go
package main

import (
    "github.com/darabuchi/prism/pkg/errcode"
)

func main() {
    // 1. 加载本地化文件
    errcode.LoadLocalizations("./resource/localize")

    // 2. 获取翻译
    msg := errcode.T("zh-CN", "common.success")
    fmt.Println(msg) // 输出: 操作成功

    // 3. 格式化翻译
    msg = errcode.Tf("zh-CN", "messages.server_started", 8080)
    fmt.Println(msg) // 输出: 服务器已启动，监听端口 8080

    // 4. 创建本地化错误
    err := errcode.NewLocalizedError("zh-CN", errcode.ErrSubscriptionNotFound)
    fmt.Println(err) // 输出: [3000] 订阅不存在
}
```

### 根据用户语言偏好

```go
func GetUserLocale(r *http.Request) string {
    // 1. 从请求头获取
    acceptLang := r.Header.Get("Accept-Language")

    // 2. 从查询参数获取
    if lang := r.URL.Query().Get("lang"); lang != "" {
        return lang
    }

    // 3. 从 Cookie 获取
    if cookie, err := r.Cookie("locale"); err == nil {
        return cookie.Value
    }

    // 4. 解析 Accept-Language 头
    // "zh-CN,zh;q=0.9,en;q=0.8" -> "zh-CN"

    // 5. 使用默认语言
    return "zh-CN"
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    locale := GetUserLocale(r)

    err := DoSomething()
    if err != nil {
        codedErr := errcode.WrapLocalizedError(locale, errcode.ErrInternal, err)
        // 返回本地化的错误消息
    }
}
```

## 最佳实践

### 1. 完整性检查

确保所有语言文件包含相同的键：

```bash
# 提取所有键
yq eval 'keys | .[]' zh-CN.yaml | sort > keys-zh.txt
yq eval 'keys | .[]' en-US.yaml | sort > keys-en.txt

# 比较差异
diff keys-zh.txt keys-en.txt
```

### 2. 避免硬编码文本

❌ 不好：
```go
return errors.New("订阅不存在")
```

✅ 好：
```go
return errcode.NewLocalizedError(locale, errcode.ErrSubscriptionNotFound)
```

### 3. 统一管理翻译

- 所有面向用户的文本都应该在 YAML 文件中定义
- 不要在代码中直接写入翻译文本
- 使用错误码系统统一管理错误消息

### 4. 定期审查

- 定期检查未翻译的键
- 审查翻译质量
- 更新过时的翻译

### 5. 版本控制

- 翻译文件纳入版本控制
- 翻译变更需要 Code Review
- 使用 Git blame 追踪翻译来源

### 6. 注释说明

对于复杂的翻译，添加注释说明上下文：

```yaml
# 用于显示订阅更新的成功消息，%s 为订阅名称
messages:
  subscription_updated: "订阅 %s 已更新"

# 用于显示批量操作结果，%d 为成功数量，第二个 %d 为总数
messages:
  batch_result: "成功处理 %d/%d 项"
```

## 工具推荐

### YAML 编辑器

- **VS Code** + `YAML` 扩展
- **IntelliJ IDEA** 内置 YAML 支持
- **在线编辑器**: [YAML Lint](http://www.yamllint.com/)

### 翻译工具

- **机器翻译**: Google Translate, DeepL
- **术语管理**: Terminology 数据库
- **本地化平台**: Crowdin, Lokalise

### 验证工具

```bash
# 安装 yamllint
pip install yamllint

# 验证 YAML 文件
yamllint resource/localize/*.yaml
```

## 参考资料

### 语言编码标准

- **[Lingoes 语言代码表](http://www.lingoes.net/zh/translator/langcode.htm)** - 本项目采用的主要参考规范
- **[ISO 639-1](https://en.wikipedia.org/wiki/ISO_639-1)** - 两字母语言代码标准
- **[ISO 639-2](https://www.loc.gov/standards/iso639-2/)** - 三字母语言代码标准
- **[ISO 3166-1](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)** - 两字母国家/地区代码标准
- **[IANA Language Subtag Registry](https://www.iana.org/assignments/language-subtag-registry/)** - 语言标签注册表

### 技术规范

- **[YAML 1.2 规范](https://yaml.org/spec/1.2/spec.html)** - YAML 格式标准
- **[Go i18n 最佳实践](https://phrase.com/blog/posts/internationalization-i18n-go/)** - Go 语言国际化指南
- **[Unicode CLDR](http://cldr.unicode.org/)** - Unicode 通用语言数据仓库

### 工具和资源

- **[ISO 639 Language Codes](https://www.iso.org/iso-639-language-code)** - ISO 官方语言代码页面
- **[Language Code Table](http://www.lingoes.net/en/translator/langcode.htm)** - Lingoes 英文版语言代码表

## 常见问题

### Q: 如何处理复数形式？

A: 使用不同的键来处理复数：

```yaml
messages:
  item_count_one: "找到 1 个项目"
  item_count_other: "找到 %d 个项目"
```

### Q: 如何处理性别差异？

A: 在键名中包含性别信息：

```yaml
messages:
  welcome_male: "欢迎，先生"
  welcome_female: "欢迎，女士"
  welcome_neutral: "欢迎"
```

### Q: 如何处理长文本？

A: 使用 YAML 多行字符串：

```yaml
messages:
  terms_of_service: |
    第一段条款内容...

    第二段条款内容...
```

### Q: 如何更新已有翻译？

A: 遵循以下流程：

1. 修改源语言文件（通常是 `zh-CN.yaml`）
2. 同步更新其他语言文件
3. 提交 Pull Request
4. Code Review 确认
5. 合并到主分支

## 贡献指南

欢迎贡献新的语言翻译！请遵循以上规范，并确保：

1. ✅ 翻译准确、自然
2. ✅ 保持键名不变
3. ✅ 占位符位置正确
4. ✅ YAML 格式正确
5. ✅ 通过 Code Review

---

**最后更新**: 2025-10-01
**维护者**: Prism Team
