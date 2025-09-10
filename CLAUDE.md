# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Prism is a multi-platform proxy client based on mihomo/clash core, designed for node pool users. It consists of four main components:

1. **Core** (Go): Proxy service with RESTful API, based on mihomo/clash
2. **Desktop** (Tauri + React): Cross-platform desktop application
3. **CLI** (Go + Cobra): Command-line management tool
4. **Android** (Kotlin + Compose): Native mobile application

## Essential Commands

### Environment Setup
```bash
make doctor          # Check development environment and get fix suggestions
make setup-dev       # Install all dependencies and setup development environment
```

### Development
```bash
make dev            # Start full development environment (Core + Desktop)
make dev-core       # Start only Core service
make dev-desktop    # Start only desktop client
```

### Building
```bash
make build          # Build all components
make build-core     # Build Core service only
make build-desktop  # Build desktop application only
make build-all      # Cross-platform build for all architectures
```

### Testing
```bash
make test           # Run all tests
make test-core      # Run Core tests with coverage
make test-desktop   # Run desktop client tests
```

### Code Quality
```bash
make lint           # Run linters for all components
make fmt            # Format all code
```

### Docker
```bash
make docker-build   # Build Docker image
make docker-run     # Run containerized version
docker-compose up   # Start full stack with monitoring
```

### Component-Specific Commands

#### Core Service
```bash
cd core
go run cmd/prism-core/main.go --config config/config.yaml
go test -v ./...
go build -o ../dist/prism-core cmd/prism-core/main.go
```

#### Desktop Client
```bash
cd desktop
npm install
npm run tauri dev    # Development with hot reload
npm run tauri build  # Production build
```

#### CLI Tool
```bash
cd cli
go run main.go status            # Check service status
go run main.go --server http://localhost:9090 status
```

## Architecture Overview

### Core Service Architecture
- **Entry Point**: `core/cmd/prism-core/main.go` - Application bootstrap with graceful shutdown
- **API Layer**: `core/internal/api/` - Gin-based REST API with middleware
- **Proxy Core**: `core/internal/core/proxy.go` - Abstraction layer for mihomo integration
- **Configuration**: `core/internal/config/` - Viper-based config management with defaults
- **Storage**: Uses SQLite via GORM for persistence

### Desktop Client Architecture
- **Tauri Backend**: `desktop/src-tauri/src/main.rs` - System integration, tray, window management
- **React Frontend**: `desktop/src/` - Modern React app with TypeScript
- **State Management**: Zustand store with persistence (`desktop/src/store/useAppStore.ts`)
- **UI Components**: Ant Design with custom styling and dark mode support
- **API Communication**: Axios client for REST API calls to Core service

### CLI Tool Architecture  
- **Cobra Framework**: `cli/cmd/` - Command structure with root, status, and other commands
- **API Client**: `cli/internal/client/client.go` - HTTP client for Core service communication
- **Rich Output**: Colored terminal output with tables and formatted responses

### Cross-Component Communication
- **Core ↔ Desktop**: HTTP REST API on port 9090, WebSocket for real-time updates
- **Core ↔ CLI**: HTTP REST API with JSON responses
- **Configuration**: YAML files with environment variable overrides
- **Data Flow**: Desktop/CLI → Core API → mihomo proxy engine

### Key Design Patterns
- **Modular Architecture**: Each component is independent with clear interfaces
- **Configuration Management**: Centralized config with sensible defaults
- **Error Handling**: Structured error responses with proper HTTP status codes
- **State Management**: Reactive state with persistence for UI components
- **Graceful Shutdown**: Proper cleanup handling in Core service

### Development Workflow
- Use `make doctor` first to ensure environment is properly configured
- Start Core service first (`make dev-core`) before desktop client
- Desktop client expects Core service on `http://localhost:9090`
- All components use the same API contract defined in `docs/api-specification.md`
- Hot reload available for both desktop frontend and Core service during development

### Important File Locations
- **Core Config**: `core/config/config.yaml` - Service configuration
- **Desktop Config**: `desktop/src-tauri/tauri.conf.json` - Tauri app configuration  
- **API Routes**: `core/internal/api/server.go` - REST endpoint definitions
- **Desktop Store**: `desktop/src/store/useAppStore.ts` - Application state management
- **Build Scripts**: `Makefile` - All build and development commands

### Testing Strategy
- Go tests use testify framework with coverage reporting
- Desktop tests use standard React testing patterns
- Integration tests verify API contracts between components
- Docker-based testing for full system verification