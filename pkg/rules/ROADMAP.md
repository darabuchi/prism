# Rules 包完善方案

基于对 Mihomo/Clash 规则系统的研究，本文档提出 Prism Rules 包的完善方案。

## 当前实现状态

### ✅ 已实现功能

1. **基础规则类型**
   - DOMAIN (域名精确匹配)
   - DOMAIN-SUFFIX (域名后缀匹配)
   - DOMAIN-KEYWORD (域名关键字匹配)
   - DOMAIN-REGEX (域名正则匹配)
   - IP-CIDR (IP CIDR 匹配)
   - IP-CIDR6 (IPv6 CIDR 匹配)
   - GEOIP (GeoIP 国家代码匹配)
   - IP-ASN (ASN 号码匹配)
   - DST-PORT (目标端口匹配)
   - SRC-PORT (源端口匹配)
   - PROCESS-NAME (进程名称匹配)
   - PROCESS-PATH (进程路径匹配)
   - MATCH (匹配所有流量)

2. **核心功能**
   - 规则解析器（ParseRule, LoadRulesFromFile）
   - 规则匹配（Match）
   - GeoIP 提供者支持
   - 完整的单元测试

3. **参数支持**
   - no-resolve (IP 规则跳过域名解析)

## 待完善功能

参考 Mihomo/Clash 的实现，以下是建议添加的高级功能：

### 1. 逻辑规则（Logic Rules）⭐⭐⭐

支持规则的逻辑组合，实现复杂的匹配策略。

#### NOT (取反规则)
```text
NOT,((DOMAIN,google.com)),DIRECT
# 匹配所有非 google.com 的域名
```

#### OR (或规则)
```text
OR,((DOMAIN,google.com),(DOMAIN,youtube.com)),PROXY
# 匹配 google.com 或 youtube.com
```

#### AND (与规则)
```text
AND,((DOMAIN-SUFFIX,google.com),(DST-PORT,443)),PROXY
# 同时匹配 google.com 后缀和 443 端口
```

**实现建议**:
```go
// logic.go
type Logic struct {
    base
    ruleType RuleType // NOT, OR, AND
    rules    []Rule
}

func NewLogic(ruleType RuleType, rules []Rule, action ActionType) *Logic
func (l *Logic) Match(metadata *Metadata) bool {
    switch l.ruleType {
    case TypeNOT:
        return !l.rules[0].Match(metadata)
    case TypeOR:
        for _, rule := range l.rules {
            if rule.Match(metadata) {
                return true
            }
        }
        return false
    case TypeAND:
        for _, rule := range l.rules {
            if !rule.Match(metadata) {
                return false
            }
        }
        return true
    }
}
```

### 2. 规则集（Rule-Set）⭐⭐⭐

支持从外部文件加载规则集，便于规则管理和复用。

#### 规则集提供者（Provider）
```yaml
rule-providers:
  reject:
    type: http
    url: "https://example.com/reject.yaml"
    interval: 86400
    path: ./ruleset/reject.yaml

  direct:
    type: file
    path: ./ruleset/direct.yaml
```

#### 使用规则集
```text
RULE-SET,reject,REJECT
RULE-SET,direct,DIRECT
```

**实现建议**:
```go
// provider/provider.go
type Provider interface {
    Name() string
    Type() ProviderType
    Rules() []Rule
    Update() error
}

// provider/http_provider.go
type HTTPProvider struct {
    url      string
    interval time.Duration
    path     string
    rules    []Rule
}

// provider/file_provider.go
type FileProvider struct {
    path  string
    rules []Rule
}

// ruleset.go
type RuleSet struct {
    base
    provider Provider
}

func (r *RuleSet) Match(metadata *Metadata) bool {
    for _, rule := range r.provider.Rules() {
        if rule.Match(metadata) {
            return true
        }
    }
    return false
}
```

### 3. 扩展规则类型⭐⭐

#### GEOSITE (域名地理位置规则)
```text
GEOSITE,cn,DIRECT
GEOSITE,google,PROXY
```

基于 V2Ray 的 geosite.dat 数据库。

#### IN-TYPE (入站类型规则)
```text
IN-TYPE,HTTP,PROXY
IN-TYPE,SOCKS5,DIRECT
```

匹配入站连接类型。

#### IN-NAME (入站名称规则)
```text
IN-NAME,socks-in,PROXY
IN-NAME,http-in,DIRECT
```

匹配入站端口的名称。

#### IN-USER (入站用户规则)
```text
IN-USER,admin,PROXY
IN-USER,guest,REJECT
```

匹配认证用户名。

#### NETWORK-TYPE (网络类型规则)
```text
NETWORK,TCP,DIRECT
NETWORK,UDP,PROXY
```

匹配 TCP 或 UDP。

#### UID (用户 ID 规则)
```text
UID,1000,DIRECT
```

匹配 Linux/Android 的 UID。

#### DSCP (DSCP 规则)
```text
DSCP,46,PROXY
```

匹配 DSCP 值。

#### IP-SUFFIX (IP 后缀规则)
```text
IP-SUFFIX,8.8.8.8/24,1,DIRECT
```

匹配 IP 的特定后缀位。

**实现建议**:
```go
// geosite.go
type GeoSite struct {
    base
    category string
    matcher  *DomainMatcher // 域名匹配器
}

// in_type.go
type InType struct {
    base
    inType string // HTTP, SOCKS5, etc.
}

// network_type.go
type NetworkType struct {
    base
    network string // TCP, UDP
}
```

### 4. 正则表达式扩展⭐

支持更多正则匹配类型。

#### PROCESS-NAME-REGEX
```text
PROCESS-NAME-REGEX,^chrome.*$,PROXY
```

#### PROCESS-PATH-REGEX
```text
PROCESS-PATH-REGEX,^/usr/bin/.*$,DIRECT
```

**实现建议**:
```go
// process_regex.go
type ProcessRegex struct {
    base
    pattern *regexp.Regexp
    isPath  bool
}
```

### 5. SUB-RULE (子规则)⭐⭐

支持规则分组和嵌套。

```text
# 定义子规则
SUB-RULE,(DOMAIN-SUFFIX,google.com),google-rules
# ... google-rules 中的规则

# 使用子规则
google-rules:
  - DOMAIN,mail.google.com,DIRECT
  - DOMAIN,drive.google.com,PROXY
```

**实现建议**:
```go
// subrule.go
type SubRule struct {
    base
    subRuleName string
    subRules    map[string][]Rule
}
```

## 实现优先级

### 高优先级 (1-2 个月)
1. **逻辑规则** (NOT, OR, AND) - 提供复杂匹配能力
2. **规则集提供者** (Rule-Set) - 提高规则管理效率
3. **文档完善** - 补充使用示例和最佳实践

### 中优先级 (2-3 个月)
4. **扩展规则类型** (IN-TYPE, IN-NAME, NETWORK-TYPE 等)
5. **GEOSITE 支持** - 基于 V2Ray geosite.dat
6. **正则扩展** (PROCESS-NAME-REGEX, PROCESS-PATH-REGEX)

### 低优先级 (3-6 个月)
7. **SUB-RULE 支持**
8. **性能优化** (规则索引、缓存优化)
9. **规则可视化** (Web UI 支持)

## 兼容性考虑

### 与 Clash/Mihomo 兼容
- 支持 Clash 的规则格式
- 兼容 Clash Premium 的高级特性
- 保持 Mihomo 的规则语法

### 扩展性设计
- 清晰的接口定义
- 插件化规则类型
- 自定义规则支持

## 性能优化建议

### 1. 规则索引
```go
type RuleIndex struct {
    domainMap    map[string][]Rule      // 精确域名索引
    suffixTree   *SuffixTree             // 后缀树索引
    ipTrie       *IPTrie                 // IP 前缀树
    portRanges   *IntervalTree           // 端口范围索引
}
```

### 2. 缓存优化
```go
type RuleCache struct {
    domainCache *lru.Cache[string, ActionType]
    ipCache     *lru.Cache[string, ActionType]
    maxSize     int
}
```

### 3. 并发优化
- 使用 sync.Map 存储规则
- 规则匹配支持并发
- 提供者更新异步化

## 测试策略

### 单元测试
- 每种规则类型的完整测试
- 边界情况测试
- 性能基准测试

### 集成测试
- 规则链匹配测试
- 规则集加载测试
- GeoIP/GeoSite 集成测试

### 性能测试
```bash
# 规则匹配性能测试
go test -bench=BenchmarkRuleMatch -benchtime=10s

# 规则解析性能测试
go test -bench=BenchmarkRuleParse -benchtime=10s

# 内存使用测试
go test -bench=BenchmarkRuleMemory -benchmem
```

## 文档完善

### 用户文档
- [ ] 规则语法完整说明
- [ ] 使用示例和最佳实践
- [ ] 常见问题和故障排查
- [ ] 性能调优指南

### 开发者文档
- [ ] API 参考文档
- [ ] 扩展开发指南
- [ ] 架构设计文档
- [ ] 贡献指南

## 里程碑

### v1.0 (当前)
- ✅ 基础规则类型
- ✅ 规则解析器
- ✅ GeoIP 支持
- ✅ 基础文档

### v1.1 (1 个月)
- [ ] 逻辑规则 (NOT, OR, AND)
- [ ] 参数支持增强
- [ ] 性能优化
- [ ] 文档完善

### v1.2 (2 个月)
- [ ] 规则集提供者 (HTTP, File)
- [ ] 规则集支持
- [ ] 规则热重载
- [ ] Web UI 基础

### v2.0 (3 个月)
- [ ] GEOSITE 支持
- [ ] 扩展规则类型
- [ ] 正则扩展
- [ ] SUB-RULE 支持
- [ ] 完整的 Web UI

### v2.1 (6 个月)
- [ ] 规则可视化编辑
- [ ] 规则调试工具
- [ ] 性能监控
- [ ] 规则分析报告

## 贡献指南

欢迎贡献！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

### 代码规范
- 遵循 Go 代码规范
- 添加完整的注释
- 包含单元测试
- 更新相关文档

### 提交规范
```
<type>(<scope>): <subject>

<body>

<footer>
```

类型：
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式
- refactor: 重构
- perf: 性能优化
- test: 测试
- chore: 构建/工具

## 参考资料

- [Clash 文档](https://github.com/Dreamacro/clash/wiki)
- [Mihomo 文档](https://wiki.metacubex.one/)
- [V2Ray 规则](https://www.v2ray.com/)
- [Surge 规则](https://manual.nssurge.com/rule/ruleset.html)

## 许可证

MIT License

---

**更新日期**: 2025-01-10
**版本**: 1.0
**维护者**: Prism Team
