# Prism - 部署和分发指南

## 部署概述

Prism 采用模块化架构，支持多种部署方式。本指南涵盖开发、测试、生产环境的部署策略，以及各平台客户端的分发方案。

## Core 服务部署

### Docker 部署 (推荐)

#### Dockerfile
```dockerfile
# core/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prism-core cmd/prism-core/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/prism-core .
COPY --from=builder /app/config ./config

EXPOSE 9090 7890

CMD ["./prism-core"]
```

#### docker-compose.yml
```yaml
version: '3.8'

services:
  prism-core:
    build: ./core
    ports:
      - "9090:9090"  # API 端口
      - "7890:7890"  # 代理端口
    volumes:
      - ./data:/data
      - ./config:/config
    environment:
      - PRISM_CONFIG_PATH=/config/config.yaml
      - PRISM_DATA_PATH=/data
    restart: unless-stopped
    
  # 可选：数据库
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  redis_data:
```

#### 启动部署
```bash
# 构建和启动
docker-compose up -d

# 查看日志
docker-compose logs -f prism-core

# 停止服务
docker-compose down
```

### 二进制部署

#### 系统服务配置
```ini
# /etc/systemd/system/prism-core.service
[Unit]
Description=Prism Core Service
After=network.target

[Service]
Type=simple
User=prism
Group=prism
WorkingDirectory=/opt/prism
ExecStart=/opt/prism/prism-core
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

#### 安装脚本
```bash
#!/bin/bash
# scripts/install.sh

set -e

INSTALL_DIR="/opt/prism"
SERVICE_USER="prism"
SERVICE_FILE="/etc/systemd/system/prism-core.service"

echo "安装 Prism Core..."

# 创建用户
if ! id "$SERVICE_USER" &>/dev/null; then
    useradd --system --shell /bin/false --home-dir /nonexistent --no-create-home $SERVICE_USER
fi

# 创建目录
mkdir -p $INSTALL_DIR
mkdir -p $INSTALL_DIR/config
mkdir -p $INSTALL_DIR/data

# 复制文件
cp prism-core $INSTALL_DIR/
cp -r config/* $INSTALL_DIR/config/
chmod +x $INSTALL_DIR/prism-core

# 设置权限
chown -R $SERVICE_USER:$SERVICE_USER $INSTALL_DIR

# 安装服务
cp prism-core.service $SERVICE_FILE
systemctl daemon-reload
systemctl enable prism-core

echo "安装完成！使用以下命令管理服务："
echo "启动: sudo systemctl start prism-core"
echo "停止: sudo systemctl stop prism-core"
echo "状态: sudo systemctl status prism-core"
```

### Kubernetes 部署

#### Deployment
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prism-core
  labels:
    app: prism-core
spec:
  replicas: 2
  selector:
    matchLabels:
      app: prism-core
  template:
    metadata:
      labels:
        app: prism-core
    spec:
      containers:
      - name: prism-core
        image: prism/core:latest
        ports:
        - containerPort: 9090
        - containerPort: 7890
        env:
        - name: PRISM_CONFIG_PATH
          value: "/config/config.yaml"
        volumeMounts:
        - name: config
          mountPath: /config
        - name: data
          mountPath: /data
        livenessProbe:
          httpGet:
            path: /api/v1/system/status
            port: 9090
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/system/status
            port: 9090
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: config
        configMap:
          name: prism-config
      - name: data
        persistentVolumeClaim:
          claimName: prism-data
```

#### Service
```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: prism-core-service
spec:
  selector:
    app: prism-core
  ports:
  - name: api
    port: 9090
    targetPort: 9090
  - name: proxy
    port: 7890
    targetPort: 7890
  type: LoadBalancer
```

## 桌面客户端分发

### 代码签名

#### macOS 签名
```bash
# 开发者证书签名
codesign --force --deep --sign "Developer ID Application: Your Name" --options runtime "target/release/bundle/macos/Prism.app"

# 公证 (Notarization)
xcrun notarytool submit "target/release/bundle/dmg/Prism_1.0.0_x64.dmg" --keychain-profile "AC_PASSWORD"

# 验证签名
codesign --verify --verbose "target/release/bundle/macos/Prism.app"
spctl --assess --verbose "target/release/bundle/macos/Prism.app"
```

#### Windows 签名
```powershell
# 使用 signtool 签名
signtool sign /f certificate.p12 /p password /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 "Prism_1.0.0_x64.msi"

# 验证签名
signtool verify /pa "Prism_1.0.0_x64.msi"
```

### 自动更新

#### Tauri 更新配置
```json
// src-tauri/tauri.conf.json
{
  "tauri": {
    "updater": {
      "active": true,
      "endpoints": [
        "https://releases.prism.app/{{target}}/{{arch}}/{{current_version}}"
      ],
      "dialog": true,
      "pubkey": "your-public-key-here"
    }
  }
}
```

#### 更新服务器
```go
// 简单的更新服务器实现
package main

import (
    "encoding/json"
    "net/http"
    "github.com/gin-gonic/gin"
)

type UpdateInfo struct {
    Version   string `json:"version"`
    Notes     string `json:"notes"`
    Pub_date  string `json:"pub_date"`
    Platforms map[string]Platform `json:"platforms"`
}

type Platform struct {
    Signature string `json:"signature"`
    URL       string `json:"url"`
}

func checkUpdate(c *gin.Context) {
    target := c.Param("target")
    arch := c.Param("arch")
    currentVersion := c.Param("current_version")
    
    // 检查是否有新版本
    latestVersion := getLatestVersion()
    if currentVersion >= latestVersion {
        c.Status(204) // No Content
        return
    }
    
    updateInfo := UpdateInfo{
        Version:  latestVersion,
        Notes:    "Bug fixes and improvements",
        Pub_date: "2024-01-15T12:00:00Z",
        Platforms: map[string]Platform{
            fmt.Sprintf("%s-%s", target, arch): {
                Signature: getSignature(target, arch, latestVersion),
                URL:       fmt.Sprintf("https://releases.prism.app/download/%s/%s/%s", target, arch, latestVersion),
            },
        },
    }
    
    c.JSON(200, updateInfo)
}
```

## Android 应用分发

### Google Play 发布

#### 应用签名
```bash
# 生成签名密钥
keytool -genkey -v -keystore prism-release-key.keystore -alias prism -keyalg RSA -keysize 2048 -validity 10000

# 配置 build.gradle
android {
    signingConfigs {
        release {
            storeFile file('../prism-release-key.keystore')
            storePassword 'your-store-password'
            keyAlias 'prism'
            keyPassword 'your-key-password'
        }
    }
    
    buildTypes {
        release {
            signingConfig signingConfigs.release
            minifyEnabled true
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }
}
```

#### 构建发布版本
```bash
# 构建 AAB (推荐)
./gradlew bundleRelease

# 构建 APK
./gradlew assembleRelease
```

### F-Droid 发布

#### Metadata 文件
```yaml
# metadata/com.prism.android.yml
Categories:
  - Internet
License: GPL-3.0-or-later
AuthorName: Prism Team
AuthorWebSite: https://prism.app
WebSite: https://prism.app
SourceCode: https://github.com/prism/prism-android
IssueTracker: https://github.com/prism/prism-android/issues

AutoName: Prism
Description: |-
    Prism is a proxy client based on mihomo/clash core.
    
    Features:
    * Support for multiple proxy protocols
    * Node pool management
    * Rule-based routing
    * Traffic statistics

RequiresRoot: false

Builds:
  - versionName: '1.0.0'
    versionCode: 1
    commit: v1.0.0
    subdir: android
    gradle:
      - yes
```

### 直接分发 (APK)

#### 下载页面
```html
<!DOCTYPE html>
<html>
<head>
    <title>Download Prism</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <div class="container">
        <h1>Download Prism</h1>
        
        <div class="platform">
            <h2>Android</h2>
            <a href="/download/prism-android-v1.0.0.apk" class="download-btn">
                Download APK (v1.0.0)
            </a>
            <p>Size: 25MB | Updated: 2024-01-15</p>
        </div>
        
        <div class="platform">
            <h2>Desktop</h2>
            <div class="desktop-downloads">
                <a href="/download/prism-desktop-macos-v1.0.0.dmg">macOS (Intel)</a>
                <a href="/download/prism-desktop-macos-arm64-v1.0.0.dmg">macOS (Apple Silicon)</a>
                <a href="/download/prism-desktop-windows-v1.0.0.msi">Windows</a>
                <a href="/download/prism-desktop-linux-v1.0.0.appimage">Linux (AppImage)</a>
            </div>
        </div>
    </div>
</body>
</html>
```

## CI/CD 流水线

### GitHub Actions

#### 构建流水线
```yaml
# .github/workflows/build.yml
name: Build and Release

on:
  push:
    tags:
      - 'v*'
  pull_request:
    branches: [ main ]

jobs:
  build-core:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.21'
    
    - name: Build Core
      run: |
        cd core
        make build-all
    
    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: prism-core
        path: dist/

  build-desktop:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
    
    - name: Install Rust
      uses: actions-rs/toolchain@v1
      with:
        toolchain: stable
    
    - name: Build Desktop App
      run: |
        cd desktop
        npm install
        npm run tauri build
    
    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: prism-desktop-${{ matrix.os }}
        path: desktop/src-tauri/target/release/bundle/

  build-android:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up JDK
      uses: actions/setup-java@v3
      with:
        java-version: '17'
        distribution: 'temurin'
    
    - name: Setup Android SDK
      uses: android-actions/setup-android@v2
    
    - name: Build Android APK
      run: |
        cd android
        ./gradlew assembleRelease
    
    - name: Upload APK
      uses: actions/upload-artifact@v3
      with:
        name: prism-android
        path: android/app/build/outputs/apk/release/

  release:
    needs: [build-core, build-desktop, build-android]
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    steps:
    - name: Download all artifacts
      uses: actions/download-artifact@v3
    
    - name: Create Release
      uses: softprops/action-gh-release@v1
      with:
        files: |
          prism-core/*
          prism-desktop-*/*
          prism-android/*
        generate_release_notes: true
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Docker 镜像发布
```yaml
# .github/workflows/docker.yml
name: Docker Build and Push

on:
  push:
    tags:
      - 'v*'

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout
      uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
    
    - name: Login to Docker Hub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKERHUB_USERNAME }}
        password: ${{ secrets.DOCKERHUB_TOKEN }}
    
    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: prism/core
        tags: |
          type=ref,event=tag
          type=raw,value=latest
    
    - name: Build and push
      uses: docker/build-push-action@v4
      with:
        context: ./core
        platforms: linux/amd64,linux/arm64
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
```

## 监控和维护

### 健康检查
```go
// internal/api/health.go
package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type HealthResponse struct {
    Status    string            `json:"status"`
    Version   string            `json:"version"`
    Uptime    int64            `json:"uptime"`
    Checks    map[string]string `json:"checks"`
}

func (s *Server) healthCheck(c *gin.Context) {
    checks := make(map[string]string)
    
    // 检查代理服务状态
    if s.proxyCore.IsRunning() {
        checks["proxy"] = "healthy"
    } else {
        checks["proxy"] = "unhealthy"
    }
    
    // 检查数据库连接
    if s.db.Ping() == nil {
        checks["database"] = "healthy"
    } else {
        checks["database"] = "unhealthy"
    }
    
    status := "healthy"
    for _, check := range checks {
        if check == "unhealthy" {
            status = "unhealthy"
            break
        }
    }
    
    c.JSON(http.StatusOK, HealthResponse{
        Status:  status,
        Version: s.version,
        Uptime:  s.uptime,
        Checks:  checks,
    })
}
```

### 日志配置
```yaml
# config/logging.yaml
logging:
  level: info
  format: json
  output: 
    - type: file
      path: /var/log/prism/app.log
      max_size: 100MB
      max_files: 10
    - type: stdout
      format: console
  
  loggers:
    proxy: debug
    api: info
    database: warn
```

### 性能监控
```go
// 集成 Prometheus 指标
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "prism_active_connections",
        Help: "Number of active proxy connections",
    })
    
    TotalTraffic = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "prism_traffic_bytes_total",
        Help: "Total traffic in bytes",
    }, []string{"direction"})
    
    RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name: "prism_request_duration_seconds",
        Help: "API request duration",
    }, []string{"method", "endpoint"})
)
```

## 备份和恢复

### 配置备份
```bash
#!/bin/bash
# scripts/backup.sh

BACKUP_DIR="/backup/prism"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# 备份配置文件
tar -czf $BACKUP_DIR/config_$DATE.tar.gz -C /opt/prism config/

# 备份数据库
sqlite3 /opt/prism/data/prism.db ".backup $BACKUP_DIR/database_$DATE.db"

# 清理旧备份 (保留30天)
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete
find $BACKUP_DIR -name "*.db" -mtime +30 -delete

echo "备份完成: $BACKUP_DIR"
```

这份部署指南提供了 Prism 项目从开发到生产的完整部署方案，涵盖了各种部署场景和平台分发策略，确保项目能够稳定可靠地运行。