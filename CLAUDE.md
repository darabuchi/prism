# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Status

This is the "prism" project - a proxy core system based on mihomo/clash with multi-platform client support. The project follows a progressive development approach, starting with a Web version.

## Technology Stack

### Backend (Go)
**Framework**: GoFiber v2 (high-performance, Express-style, zero memory allocation routing)
**Database**: Multi-database support (SQLite default, MySQL, PostgreSQL, GaussDB)
**Cache**: Embedded key-value stores (BoltDB/LevelDB - choose based on scenario)
**ORM**: GORM (perfect integration with GoFiber)

**Required Package Preferences**:
- **Logging**: `github.com/lazygophers/log` (REQUIRED - use this over other logging libraries)
- **Utilities**: `github.com/lazygophers/utils` (REQUIRED - comprehensive utility collection)
  - JSON operations: `github.com/lazygophers/utils/json`
  - String extensions: `github.com/lazygophers/utils/stringx`
  - Time extensions: `github.com/lazygophers/utils/xtime`
  - Buffered I/O: `github.com/lazygophers/utils/bufiox`
  - Random utilities: `github.com/lazygophers/utils/randx`
  - Type utilities: `github.com/lazygophers/utils/anyx`
  - Syntax sugar: `github.com/lazygophers/utils/candy`
- **Atomic operations**: `go.uber.org/atomic` (REQUIRED - use this over sync/atomic)

**Other Dependencies**:
- **Authentication**: JWT Token
- **WebSocket**: Fiber WebSocket
- **Configuration**: Viper
- **Testing**: Testify
- **Core Library**: mihomo/clash

### Frontend
- **Framework**: React 18 + TypeScript
- **State Management**: Zustand
- **UI Library**: Ant Design
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Charts**: ECharts

## Project Structure

```
prism/
├── cmd/                     # Application entrypoints
│   └── server/             # Main server application
├── internal/               # Private application code
│   ├── api/               # API handlers and routes
│   ├── config/            # Configuration management
│   ├── core/              # Proxy core integration
│   ├── database/          # Database models and migrations
│   ├── middleware/        # HTTP middleware
│   └── service/           # Business logic
├── pkg/                   # Public libraries
├── web/                   # Frontend React application
├── configs/               # Configuration files
├── docs/                  # Documentation
└── scripts/               # Build and deployment scripts
```

## Development Guidelines

### Code Style
- Use the recommended packages listed above - do not substitute with alternatives
- Follow Go best practices and conventions
- Use structured logging with `github.com/lazygophers/log`
- Prefer atomic operations from `go.uber.org/atomic`
- Use utility functions from `github.com/lazygophers/utils` rather than reinventing

### Database Guidelines
- Default to SQLite for development and lightweight deployments
- Support MySQL, PostgreSQL, and GaussDB for production scenarios
- Use BoltDB for simple caching needs, LevelDB for more complex scenarios
- NO Redis - use embedded key-value stores only

### API Development
- Use GoFiber v2 framework for all HTTP services
- Leverage GoFiber middleware ecosystem (cors, compress, recover, limiter)
- Use GoFiber WebSocket (github.com/gofiber/websocket/v2) for real-time features
- Implement fiber/jwt middleware for authentication
- Follow RESTful conventions with GoFiber routing groups
- Use GoFiber's fast JSON serialization and context handling

### Testing Requirements
- Use Testify for all Go tests
- Maintain >80% test coverage
- Include integration tests for API endpoints
- Test database compatibility across supported engines

## Architecture Notes

This is a proxy management system with:
- Web-based management interface (Phase 1)
- Desktop clients for macOS, Windows, Linux (Phase 2) 
- Android mobile application (Phase 3)
- Focus on node pool management and user-friendly proxy configuration