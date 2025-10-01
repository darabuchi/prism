# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important: Memory and Rules System

**When the user says "记住：[content]" or "Remember: [content]"**, you MUST:
1. Immediately add the content to this CLAUDE.md file under the appropriate section
2. If the content is a general rule/guideline, add it to the "Project Rules and Guidelines" section below
3. If the content is technical/architectural, add it to the relevant section
4. Commit the changes to git with a clear message
5. Confirm to the user that the information has been recorded

## 编码规范强制要求

**在进行任何代码开发、修改或审查之前，必须严格遵守以下编码规范文档：**

1. **[数据库设计规范](docs/编码规范/数据库设计规范.md)** - 必须遵守
   - 所有时间字段使用 TIMESTAMP 类型
   - 所有索引包含 deleted_at 字段
   - 禁止使用基本类型指针
   - 表结构、字段命名、索引设计等规范

2. **[API 设计规范](docs/编码规范/API设计规范.md)** - 必须遵守
   - 所有接口使用 POST 方法
   - 禁止路径变量，所有参数通过 Request Body 传递
   - 统一响应格式、错误处理、认证方式等规范

3. **[后端开发规范](docs/编码规范/后端开发规范.md)** - 必须遵守
   - 必须使用指定的依赖包（LRPC、lazygophers）
   - 错误处理、日志记录、数据验证等规范
   - 代码组织、测试标准等规范

**违规检查：**
- 在提交代码前，必须对照各规范文档的"违规检查清单"进行自查
- 任何违反规范的代码都不应被提交
- 如发现已有代码违反规范，应立即修正

## Project Rules and Guidelines

<!-- User-defined rules and memories will be added here -->
<!-- Format: Add new rules as bullet points with date stamp -->

- [2025-01-01] **工作流程规则**：及时提交变更 - 完成任何有意义的修改后应立即commit，保持Git历史清晰
- [2025-01-01] **文档维护规则**：文档相关的内容以及代码核心的内容需要及时更新到CLAUDE.md的相关索引中，保持文档与代码同步
- [2025-10-01] **编码规范文档化**：所有编码规范已移至 `docs/编码规范/` 目录下的独立文档，必须严格遵守（详见上方"编码规范强制要求"）

---

## Project Overview

This is **Prism**, an intelligent proxy client system built on top of the Mihomo/Meta kernel. The project provides automated proxy node management, comprehensive performance testing, intelligent routing, and subscription management. It continuously monitors proxy node performance and automatically selects optimal nodes based on multiple metrics.

## Architecture

### High-Level Structure

The project consists of two main layers:

1. **Mihomo Core** (`pkg/mihomo/`) - The underlying proxy kernel forked from MetaCubeX/mihomo
   - Handles actual proxy protocols (Shadowsocks, VMess, VLESS, Trojan, Hysteria, TUIC, etc.)
   - Implements low-level networking (TUN, TPROXY, HTTP/SOCKS5 listeners)
   - Provides DNS resolution, GeoIP routing, and rule matching
   - Contains the core tunnel and connection management logic

2. **Prism Application Layer** - Custom business logic on top of Mihomo
   - Subscription management and automated updates
   - Multi-dimensional node testing (latency, speed, unlock detection, geo-location, route quality)
   - Intelligent node scoring and selection algorithms
   - RESTful API, WebSocket and SSE (Server-Sent Events) for real-time updates
   - Task queue for background operations
   - Database persistence (SQLite/MySQL + BBolt cache)

### Key Components

**Proxy System:**
- Central `DefaultTunnel` coordinates all proxy connections
- Manages active node pool (default: top 10 nodes by score)
- Implements circuit breaker pattern for node health tracking
- Maintains session persistence and blocklist for failed connections
- Mixed port (HTTP + SOCKS5, default: 7899) and optional TUN interface
- **MITM Support**: HTTP/HTTPS traffic interception with dynamic certificate generation, request/response modification, and JavaScript scripting API (see `docs/产品设计文档.md` Section 5.2)

**Routing Engine:**
- Rule types: Domain, DomainSuffix, DomainKeyword, DomainRegex, GEOIP, IPCIDR, Process, ProcessPath
- Target types: Direct, Reject, Proxy, plus service-specific (OpenAI, Netflix, YouTube, Disney, etc.)
- Priority: Custom rules → Domain Trie → IP ranges → Built-in rules → Default
- Uses Trie data structure for O(m) domain lookups

**Node Selection:**
- Pluggable selection pipeline: Pre-checks → Selection strategy → Post-checks
- Strategies: Sticky (session persistence), RandomRange, WeakRandomRange, RandomOne, Range (score-based)
- Filters: BlockSession, BeforeRequest (circuit breaker), Unlock, TotalRequest, Country
- Fallback chain for graceful degradation

**Testing System:**
- **Delay Test**: 10 samples to google.com/generate_204, default 10min interval
- **Speed Test**: Download/upload via Ookla/Cloudflare/LibreSpeed, monthly interval, daytime only
- **Unlock Detection**: Tests OpenAI, Netflix, YouTube, Disney+ access, monthly interval
- **Geo-Location**: Determines entry/exit country via ipinfo.io
- **IP Risk**: Assesses IP reputation and threat level
- **Route Quality**: Traceroute analysis for network path optimization

**Scoring Algorithm:**
- Base factors: (1 - delay/10000) × 10 + download/50MB
- Protocol bonuses: SS (+3), Trojan (+2), VMess (+1)
- Unlock bonus: +0.025 per service (max +0.1)
- Geo bonus: +0.005 (different entry/exit) or +0.0025 (same)
- Route quality: +3 for high-quality paths
- Penalties: -10 (CN exit), +3 (CN entry), -100 (dead nodes)

**Subscription Management:**
- Auto-updates at configurable intervals (default: 6 hours)
- Supports Clash YAML, V2Ray links, mixed formats
- Tracks traffic stats, expiration, node counts
- Auto-disable after max retries (default: 5)
- Optional sharing endpoints for external clients

**Task Queue:**
- Priority levels: High (user-initiated), Normal (scheduled), Low (maintenance)
- Task types: UpdateSubscribe, CheckNodeDelay, CheckNodeDownload, CheckNodeUnlock, CheckNodeRegion, CheckNodeRoutes
- Concurrency control: Download tests mutex-protected (one at a time)
- Automatic retry with backoff for transient errors

**Scheduled Jobs:**
- Wake subscriptions (every 5 min): Queue up to 3 due updates
- Wake node checks (every 2 min): Queue up to 100 due tests
- Reload nodes (every 5 min): Refresh active node pool
- Daily cleanup (12:05-12:25): Remove disabled subscriptions/nodes, expired data, orphaned records
- Database optimization (hourly :55): VACUUM/optimize tables

## Commands

### Building

Build the Mihomo kernel for current platform:
```bash
cd pkg/mihomo
make darwin-arm64      # For macOS ARM
make linux-amd64-v3    # For Linux AMD64
make windows-amd64-v3  # For Windows AMD64
```

Build all common platforms:
```bash
cd pkg/mihomo
make all
```

Build with gvisor TUN stack (recommended):
```bash
cd pkg/mihomo
go build -tags with_gvisor -trimpath
```

### Testing

Run tests in Mihomo:
```bash
cd pkg/mihomo
go test ./...
```

Run linting:
```bash
cd pkg/mihomo
golangci-lint run ./...
```

Test configuration and exit:
```bash
./mihomo -t -f config.yaml
```

### Running

Run with custom config:
```bash
./mihomo -f config.yaml -d /path/to/home
```

Run with environment variables:
```bash
export CLASH_HOME_DIR=/path/to/home
export CLASH_CONFIG_FILE=/path/to/config.yaml
./mihomo
```

Convert ruleset (utility command):
```bash
./mihomo convert-ruleset [options]
```

Generate configurations:
```bash
./mihomo generate [options]
```

## Development Guidelines

### Code Organization

- **Mihomo kernel code** lives in `pkg/mihomo/` and should be treated as upstream dependency
- **Prism application code** should be kept separate from mihomo core
- The project uses Go 1.25.1 with go modules

### Backend Dependencies

**LRPC Middleware**:
- `github.com/lazygophers/lrpc/middleware/storage/db` - Database middleware
- `github.com/lazygophers/lrpc/middleware/xerror` - Error handling middleware
- `github.com/lazygophers/lrpc/middleware/i18n` - Internationalization middleware
- `github.com/lazygophers/lrpc/middleware/core` - Core middleware

**Utility Packages**:
- `github.com/lazygophers/log` - Logging
- `github.com/lazygophers/utils/candy` - Utility functions
- `github.com/lazygophers/utils/json` - JSON handling
- `github.com/lazygophers/utils/xtime` - Time utilities
- `github.com/lazygophers/utils/validator` - Data validation
- `github.com/lazygophers/utils/cryptox` - Cryptography utilities
- `github.com/lazygophers/utils/app` - Application utilities
- `github.com/lazygophers/utils/wait` - Wait/synchronization utilities
- `github.com/lazygophers/utils/runtime` - Runtime utilities
- `github.com/lazygophers/utils/config` - Configuration management

### Important Patterns

**Circuit Breaker Implementation:**
- Each node maintains state: Open (healthy) → Closed (unhealthy) → Half-Open (probing)
- Thresholds depend on request count (60%-90% success rate required)
- 5-minute rolling window, tracks last 200 requests
- 10% probe chance when closed

**Node Lifecycle:**
- Created via subscription update → Tested → Scored → Activated → Used → Degraded → Disabled/Deleted
- Soft delete with 1-week retention before hard delete
- Dead nodes use exponential backoff for testing (interval × death_count)

**Configuration Changes:**
- Some take immediate effect (rules, plugins)
- Some require proxy reload (ports, TUN)
- Some require restart (database path, cache)
- Always validate before persistence

**Private Network Handling:**
- Private IPs (10.0.0.0/8, 192.168.0.0/16, etc.) always bypass proxy
- Known DNS servers (223.5.5.5, 119.29.29.29) forced to direct connection
- Automatic detection prevents proxy loops

### Testing Considerations

- Speed tests only run during daytime (8am-6pm) unless forced
- Speed tests skip if network busy (>200KB/s)
- Download tests are mutex-protected (serial execution)
- All tests respect circuit breaker state but don't get blocked by it

### Performance Notes

- Active node pool limited to configurable count (default: 10) for memory efficiency
- Domain rules use Trie for O(m) lookups instead of linear search
- Session cache uses LRU with 300s TTL, max 128 entries
- Blocklist cache uses LRU with 60s TTL, max 32 entries
- Database vacuuming runs hourly to prevent bloat

### Protocol Support

Supported proxy protocols in Mihomo:
- Shadowsocks (SS)
- ShadowsocksR (SSR)
- VMess
- VLESS
- Trojan
- Hysteria 1 & 2
- WireGuard
- TUIC
- Mieru
- AnyTLS

### API Integration

The system exposes comprehensive RESTful APIs for:
- Subscription CRUD operations
- Node management and testing
- Configuration updates
- Real-time connection monitoring
- Rule matching simulation
- WebSocket and SSE (Server-Sent Events) for real-time updates

See 产品设计文档.md sections 9.1-9.8 for detailed API documentation.

## Common Gotchas

- **Do not use `net.DefaultResolver`** - The main.go intentionally traps this to prevent DNS leaks
- **TUN interface requires elevated privileges** - Must run as root/admin
- **Mihomo expects specific config structure** - See `pkg/mihomo/docs/config.yaml` for reference
- **Go proxy may be needed in China** - Use `go env -w GOPROXY=https://goproxy.io,direct`
- **License restriction** - GPL-3.0 with naming restriction: downstream projects cannot use "mihomo" in their name
- **Platform-specific features** - TUN and system proxy integration depend on OS support

## Resources

- Mihomo Documentation: https://wiki.metacubex.one/
- Mihomo Dashboard: https://github.com/MetaCubeX/metacubexd
- Product Documentation: See 产品设计文档.md (Chinese) for comprehensive platform support, deployment options and system behavior documentation
