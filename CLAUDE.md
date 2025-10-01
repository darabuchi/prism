# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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
   - RESTful API and WebSocket event system
   - Task queue for background operations
   - Database persistence (SQLite/MySQL + BBolt cache)

### Key Components

**Proxy System:**
- Central `DefaultTunnel` coordinates all proxy connections
- Manages active node pool (default: top 10 nodes by score)
- Implements circuit breaker pattern for node health tracking
- Maintains session persistence and blocklist for failed connections
- Mixed port (HTTP + SOCKS5, default: 7899) and optional TUN interface

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
- WebSocket events for live updates

See 产品功能文档.md sections 9.1-9.8 for detailed API documentation.

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
- Product Documentation: See 产品功能文档.md (Chinese) for comprehensive system behavior documentation
