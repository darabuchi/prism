# Rule Collector - Prism 规则收集工具

自动从多个知名规则源收集、解析和导出路由规则，用于 Prism 代理系统的智能路由功能。

## 功能特性

- **多源聚合**：支持 BlackMatrix7、Loyalsoldier、ACL4SSR 等主流规则源
- **智能缓存**：自动缓存下载的规则，避免重复请求
- **规则去重**：自动去除重复规则，优化规则集大小
- **分类导出**：按动作（PROXY、DIRECT、REJECT）分类导出规则
- **格式转换**：自动将各种格式转换为 Prism 规则格式

## 数据源

### BlackMatrix7
- 来源：https://github.com/blackmatrix7/ios_rule_script
- 特点：规则分类详细，覆盖面广
- 更新频率：每日更新

### Loyalsoldier
- 来源：https://github.com/Loyalsoldier/clash-rules
- 特点：规则精简，维护活跃
- 更新频率：每日更新

### ACL4SSR
- 来源：https://github.com/ACL4SSR/ACL4SSR
- 特点：规则全面，社区维护
- 更新频率：定期更新

## 使用方法

### 运行收集器

```bash
cd scripts/tools/rule-collector
make run
```

### 清理缓存

```bash
make clean
```

### 构建二进制

```bash
make build
```

## 输出文件

规则文件将导出到 `output/` 目录：

```
output/
├── PROXY.txt         # 代理规则
├── DIRECT.txt        # 直连规则
├── REJECT.txt        # 拒绝规则
└── all_rules.txt     # 所有规则汇总
```

## 规则格式

导出的规则使用 Prism 标准格式：

```
TYPE,PAYLOAD,ACTION
```

示例：
```
DOMAIN,google.com,PROXY
DOMAIN-SUFFIX,example.com,DIRECT
IP-CIDR,192.168.0.0/16,DIRECT
GEOIP,CN,DIRECT
```

## 规则分类

### 代理规则 (PROXY)

- **流媒体**：YouTube、Netflix、Disney+、Spotify 等
- **AI 服务**：OpenAI、Claude、Gemini、Copilot 等
- **开发工具**：GitHub、GitLab、Docker Hub、NPM 等
- **社交平台**：Telegram、Twitter、Facebook、Instagram 等
- **科技公司**：Google、Apple、Microsoft、Amazon 等

### 直连规则 (DIRECT)

- **局域网**：局域网地址段
- **国内服务**：国内常用网站和应用
- **国内流媒体**：哔哩哔哩、爱奇艺、腾讯视频等
- **国内公司**：阿里巴巴、腾讯、百度、字节跳动等

### 拒绝规则 (REJECT)

- **广告拦截**：广告域名和追踪域名
- **隐私保护**：隐私追踪和数据收集
- **反劫持**：防止 DNS 劫持和 HTTP 劫持

## 自定义规则

### 添加自定义数据源

1. 在 `collector/` 目录创建新的处理器文件
2. 实现 `Handler` 接口：
   ```go
   type Handler interface {
       Download(path string) ([]byte, error)
       Parse(body []byte, action string) ([]rules.Rule, error)
       NeedUpdate(info os.FileInfo) bool
   }
   ```
3. 在 `main.go` 中注册处理器

### 修改规则配置

编辑 `main.go` 中的 `parseList` 数组，添加或修改规则源配置：

```go
parseList := []ParseConfig{
    {BLACKMATRIX7, "YouTube/YouTube.yaml", "PROXY"},
    // 添加更多配置...
}
```

## 缓存机制

- 缓存目录：`tmp/cache/`
- 缓存文件以 SHA256 哈希命名
- 缓存有效期：24 小时
- 清理缓存：`make clean`

## 性能优化

- 使用进度条显示处理进度
- 智能缓存减少网络请求
- 规则去重优化内存使用
- 并发处理提高速度

## 故障排除

### 下载失败

- 检查网络连接
- 确认数据源可访问
- 查看错误日志

### 解析错误

- 检查规则格式是否正确
- 查看不支持的规则类型
- 跳过无效规则继续处理

### 缓存问题

- 清理缓存重新下载：`make clean && make run`
- 检查缓存目录权限

## 扩展开发

### 添加新的规则类型

在 `collector/parser.go` 中扩展解析逻辑：

```go
func parseRuleLine(line, action string) (rules.Rule, error) {
    // 添加新的规则类型处理
}
```

### 实现规则覆盖检查

在 `collector.go` 中完善 `ruleCovers` 方法：

```go
func (c *Collector) ruleCovers(existing, new rules.Rule) bool {
    // 实现智能覆盖检查
}
```

## 依赖项

- `github.com/darabuchi/prism/pkg/rules` - Prism 规则包
- `github.com/pterm/pterm` - 终端美化输出

## 许可证

本工具遵循 Prism 项目的许可证。

## 参考资料

- [BlackMatrix7 规则文档](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash)
- [Loyalsoldier 规则文档](https://github.com/Loyalsoldier/clash-rules)
- [ACL4SSR 规则文档](https://github.com/ACL4SSR/ACL4SSR)
- [Prism 规则引擎文档](../../../pkg/rules/README.md)
