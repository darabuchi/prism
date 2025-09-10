# Prism - 多平台客户端开发指南

## 开发环境准备

### 通用依赖
- **Node.js**: 18.0+ (桌面客户端)
- **Go**: 1.21+ (Core 和 CLI)
- **Git**: 版本控制
- **Docker**: 容器化部署

### 平台特定依赖

#### macOS 开发环境
```bash
# 安装 Xcode Command Line Tools
xcode-select --install

# 安装 Rust (Tauri 依赖)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 安装 Tauri CLI
cargo install tauri-cli
```

#### Windows 开发环境
```powershell
# 安装 Visual Studio Build Tools
# 安装 Rust
winget install Rustlang.Rustup

# 安装 Tauri CLI
cargo install tauri-cli
```

#### Linux 开发环境
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y curl build-essential libssl-dev pkg-config

# 安装 Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 安装系统依赖
sudo apt install -y libwebkit2gtk-4.0-dev libgtk-3-dev libayatana-appindicator3-dev librsvg2-dev
```

#### Android 开发环境
```bash
# 安装 Android Studio 和 SDK
# 配置环境变量
export ANDROID_HOME=$HOME/Android/Sdk
export PATH=$PATH:$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools
```

## Core 开发 (prism-core)

### 项目结构
```
core/
├── cmd/
│   └── prism-core/
│       └── main.go              # 入口文件
├── internal/
│   ├── api/                     # API 服务器
│   ├── config/                  # 配置管理
│   ├── core/                    # 代理核心
│   ├── nodepool/               # 节点池管理
│   ├── subscription/           # 订阅管理
│   ├── rule/                   # 规则引擎
│   └── storage/                # 数据存储
├── pkg/
│   ├── logger/                 # 日志系统
│   └── utils/                  # 工具函数
├── config/
│   └── config.yaml             # 默认配置
├── go.mod
└── go.sum
```

### 开发步骤

#### 1. 初始化项目
```bash
cd core
go mod init github.com/yourusername/prism-core
go mod tidy
```

#### 2. 集成 mihomo
```go
// internal/core/proxy.go
package core

import (
    "github.com/metacubex/mihomo/config"
    "github.com/metacubex/mihomo/hub/executor"
)

type ProxyCore struct {
    config *config.Config
}

func NewProxyCore() *ProxyCore {
    return &ProxyCore{}
}

func (p *ProxyCore) Start() error {
    cfg, err := config.Parse([]byte(configContent))
    if err != nil {
        return err
    }
    
    executor.ApplyConfig(cfg, true)
    return nil
}
```

#### 3. API 服务器实现
```go
// internal/api/server.go
package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Server struct {
    engine *gin.Engine
    port   string
}

func NewServer(port string) *Server {
    r := gin.Default()
    s := &Server{
        engine: r,
        port:   port,
    }
    s.setupRoutes()
    return s
}

func (s *Server) setupRoutes() {
    v1 := s.engine.Group("/api/v1")
    {
        v1.GET("/system/status", s.getSystemStatus)
        v1.GET("/nodepools", s.getNodePools)
        v1.POST("/nodepools", s.createNodePool)
        // 更多路由...
    }
}
```

### 构建和测试
```bash
# 构建
go build -o bin/prism-core cmd/prism-core/main.go

# 运行测试
go test ./...

# 生成代码覆盖率报告
go test -cover ./...
```

## 桌面客户端开发 (prism-desktop)

### 技术栈
- **前端**: React 18 + TypeScript + Vite
- **UI**: Ant Design 5.0
- **状态管理**: Zustand
- **路由**: React Router 6
- **HTTP**: Axios
- **图表**: Recharts
- **后端**: Tauri 2.0 (Rust)

### 项目初始化
```bash
# 创建 Tauri 项目
npm create tauri-app@latest prism-desktop
cd prism-desktop

# 安装依赖
npm install
npm install antd @ant-design/icons
npm install zustand react-router-dom axios recharts
npm install -D @types/node
```

### Tauri 配置
```json
// src-tauri/tauri.conf.json
{
  "package": {
    "productName": "Prism",
    "version": "1.0.0"
  },
  "tauri": {
    "allowlist": {
      "all": false,
      "http": {
        "all": true,
        "request": true
      },
      "notification": {
        "all": true
      },
      "systemTray": {
        "all": true
      },
      "window": {
        "all": false,
        "minimize": true,
        "maximize": true,
        "close": true,
        "hide": true,
        "show": true
      }
    },
    "systemTray": {
      "iconPath": "icons/icon.png",
      "iconAsTemplate": true
    },
    "windows": [
      {
        "title": "Prism",
        "width": 1200,
        "height": 800,
        "minWidth": 800,
        "minHeight": 600,
        "center": true,
        "visible": false
      }
    ]
  }
}
```

### 状态管理
```typescript
// src/store/useAppStore.ts
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface SystemStatus {
  version: string
  uptime: number
  proxy_status: 'running' | 'stopped'
  memory_usage: number
  cpu_usage: number
}

interface AppState {
  systemStatus: SystemStatus | null
  currentNode: any | null
  nodePools: any[]
  
  // Actions
  setSystemStatus: (status: SystemStatus) => void
  setCurrentNode: (node: any) => void
  setNodePools: (pools: any[]) => void
}

export const useAppStore = create<AppState>()(
  persist(
    (set) => ({
      systemStatus: null,
      currentNode: null,
      nodePools: [],
      
      setSystemStatus: (status) => set({ systemStatus: status }),
      setCurrentNode: (node) => set({ currentNode: node }),
      setNodePools: (pools) => set({ nodePools: pools }),
    }),
    {
      name: 'prism-store',
    }
  )
)
```

### API 客户端
```typescript
// src/services/api.ts
import axios from 'axios'

const BASE_URL = 'http://localhost:9090/api/v1'

export const apiClient = axios.create({
  baseURL: BASE_URL,
  timeout: 10000,
})

export const systemApi = {
  getStatus: () => apiClient.get('/system/status'),
  getConfig: () => apiClient.get('/system/config'),
  updateConfig: (config: any) => apiClient.put('/system/config', config),
}

export const nodePoolApi = {
  getList: () => apiClient.get('/nodepools'),
  create: (data: any) => apiClient.post('/nodepools', data),
  update: (id: string, data: any) => apiClient.put(`/nodepools/${id}`, data),
  delete: (id: string) => apiClient.delete(`/nodepools/${id}`),
}
```

### 主要组件
```tsx
// src/components/Dashboard/index.tsx
import React, { useEffect } from 'react'
import { Card, Row, Col, Statistic, Button } from 'antd'
import { PlayCircleOutlined, PauseCircleOutlined } from '@ant-design/icons'
import { useAppStore } from '../../store/useAppStore'
import { systemApi } from '../../services/api'
import TrafficChart from './TrafficChart'

const Dashboard: React.FC = () => {
  const { systemStatus, setSystemStatus } = useAppStore()

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const response = await systemApi.getStatus()
        setSystemStatus(response.data.data)
      } catch (error) {
        console.error('Failed to fetch system status:', error)
      }
    }

    fetchStatus()
    const interval = setInterval(fetchStatus, 5000) // 每5秒更新
    return () => clearInterval(interval)
  }, [setSystemStatus])

  const toggleProxy = async () => {
    // 切换代理状态的逻辑
  }

  return (
    <div className="dashboard">
      <Row gutter={[16, 16]}>
        <Col span={6}>
          <Card>
            <Statistic
              title="代理状态"
              value={systemStatus?.proxy_status === 'running' ? '运行中' : '已停止'}
              prefix={systemStatus?.proxy_status === 'running' ? 
                <PlayCircleOutlined style={{ color: '#52c41a' }} /> : 
                <PauseCircleOutlined style={{ color: '#ff4d4f' }} />
              }
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic 
              title="内存使用" 
              value={systemStatus?.memory_usage || 0} 
              suffix="MB"
              precision={1}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic 
              title="CPU使用率" 
              value={systemStatus?.cpu_usage || 0} 
              suffix="%"
              precision={1}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="运行时间" value={systemStatus?.uptime || 0} suffix="秒" />
          </Card>
        </Col>
      </Row>
      
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col span={24}>
          <Card title="流量统计" extra={
            <Button 
              type="primary" 
              onClick={toggleProxy}
              icon={systemStatus?.proxy_status === 'running' ? 
                <PauseCircleOutlined /> : 
                <PlayCircleOutlined />
              }
            >
              {systemStatus?.proxy_status === 'running' ? '停止代理' : '启动代理'}
            </Button>
          }>
            <TrafficChart />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
```

### 构建和开发
```bash
# 开发模式
npm run tauri dev

# 构建生产版本
npm run tauri build

# 只构建前端
npm run build
```

## Android 客户端开发 (prism-android)

### 项目结构
```
android/
├── app/
│   ├── src/main/
│   │   ├── java/com/prism/android/
│   │   │   ├── ui/                    # UI 组件
│   │   │   ├── data/                  # 数据层
│   │   │   ├── domain/                # 业务逻辑层
│   │   │   ├── di/                    # 依赖注入
│   │   │   └── service/               # 后台服务
│   │   ├── res/                       # 资源文件
│   │   └── AndroidManifest.xml
│   └── build.gradle
├── build.gradle
└── gradle.properties
```

### 依赖配置
```gradle
// app/build.gradle
android {
    compileSdk 34
    
    defaultConfig {
        applicationId "com.prism.android"
        minSdk 21
        targetSdk 34
        versionCode 1
        versionName "1.0"
    }
    
    compileOptions {
        sourceCompatibility JavaVersion.VERSION_17
        targetCompatibility JavaVersion.VERSION_17
    }
    
    kotlinOptions {
        jvmTarget = '17'
    }
    
    buildFeatures {
        compose = true
    }
    
    composeOptions {
        kotlinCompilerExtensionVersion = '1.5.4'
    }
}

dependencies {
    // Jetpack Compose
    implementation 'androidx.compose.ui:ui:1.5.4'
    implementation 'androidx.compose.ui:ui-tooling-preview:1.5.4'
    implementation 'androidx.compose.material3:material3:1.1.2'
    implementation 'androidx.activity:activity-compose:1.8.0'
    
    // Architecture Components
    implementation 'androidx.lifecycle:lifecycle-viewmodel-compose:2.6.2'
    implementation 'androidx.navigation:navigation-compose:2.7.4'
    
    // Networking
    implementation 'com.squareup.retrofit2:retrofit:2.9.0'
    implementation 'com.squareup.retrofit2:converter-gson:2.9.0'
    implementation 'com.squareup.okhttp3:logging-interceptor:4.11.0'
    
    // Database
    implementation 'androidx.room:room-runtime:2.6.0'
    implementation 'androidx.room:room-ktx:2.6.0'
    kapt 'androidx.room:room-compiler:2.6.0'
    
    // Dependency Injection
    implementation 'com.google.dagger:hilt-android:2.48'
    implementation 'androidx.hilt:hilt-navigation-compose:1.1.0'
    kapt 'com.google.dagger:hilt-compiler:2.48'
}
```

### VPN 服务实现
```kotlin
// service/VpnService.kt
package com.prism.android.service

import android.net.VpnService
import android.os.ParcelFileDescriptor
import kotlinx.coroutines.*

class PrismVpnService : VpnService() {
    private var vpnInterface: ParcelFileDescriptor? = null
    private val serviceJob = SupervisorJob()
    private val serviceScope = CoroutineScope(Dispatchers.IO + serviceJob)
    
    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        return when (intent?.action) {
            ACTION_START_VPN -> {
                startVpn()
                START_STICKY
            }
            ACTION_STOP_VPN -> {
                stopVpn()
                START_NOT_STICKY
            }
            else -> START_NOT_STICKY
        }
    }
    
    private fun startVpn() {
        val builder = Builder()
        builder.setSession("Prism VPN")
        builder.setMtu(1500)
        builder.addAddress("10.0.0.2", 24)
        builder.addRoute("0.0.0.0", 0)
        builder.addDnsServer("8.8.8.8")
        
        // 设置允许的应用
        val allowedApps = getSelectedApps()
        allowedApps.forEach { packageName ->
            try {
                builder.addAllowedApplication(packageName)
            } catch (e: PackageManager.NameNotFoundException) {
                // 应用不存在
            }
        }
        
        vpnInterface = builder.establish()
        
        serviceScope.launch {
            runVpnLoop()
        }
    }
    
    private fun stopVpn() {
        vpnInterface?.close()
        vpnInterface = null
        serviceJob.cancelChildren()
    }
    
    private suspend fun runVpnLoop() {
        // VPN 数据包处理逻辑
        // 与 Core API 通信，转发流量
    }
    
    companion object {
        const val ACTION_START_VPN = "START_VPN"
        const val ACTION_STOP_VPN = "STOP_VPN"
    }
}
```

### Jetpack Compose UI
```kotlin
// ui/dashboard/DashboardScreen.kt
package com.prism.android.ui.dashboard

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DashboardScreen(
    viewModel: DashboardViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp)
    ) {
        // 连接状态卡片
        Card(
            modifier = Modifier.fillMaxWidth(),
            elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
        ) {
            Column(
                modifier = Modifier.padding(16.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = if (uiState.isConnected) "已连接" else "未连接",
                    style = MaterialTheme.typography.headlineMedium
                )
                
                Spacer(modifier = Modifier.height(16.dp))
                
                Button(
                    onClick = { viewModel.toggleConnection() },
                    colors = ButtonDefaults.buttonColors(
                        containerColor = if (uiState.isConnected) 
                            MaterialTheme.colorScheme.error else 
                            MaterialTheme.colorScheme.primary
                    )
                ) {
                    Text(if (uiState.isConnected) "断开连接" else "连接")
                }
            }
        }
        
        Spacer(modifier = Modifier.height(16.dp))
        
        // 流量统计
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            TrafficCard(
                title = "上传",
                value = uiState.uploadTraffic,
                modifier = Modifier.weight(1f)
            )
            
            Spacer(modifier = Modifier.width(8.dp))
            
            TrafficCard(
                title = "下载", 
                value = uiState.downloadTraffic,
                modifier = Modifier.weight(1f)
            )
        }
    }
}

@Composable
fun TrafficCard(
    title: String,
    value: String,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier,
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(
                text = title,
                style = MaterialTheme.typography.bodyMedium
            )
            Text(
                text = value,
                style = MaterialTheme.typography.headlineSmall
            )
        }
    }
}
```

### 构建配置
```bash
# Debug 构建
./gradlew assembleDebug

# Release 构建
./gradlew assembleRelease

# 安装到设备
./gradlew installDebug
```

## CLI 工具开发 (prism-cli)

### 项目结构
```
cli/
├── cmd/
│   ├── root.go                 # 根命令
│   ├── start.go               # 启动命令
│   ├── stop.go                # 停止命令
│   ├── status.go              # 状态命令
│   └── config.go              # 配置命令
├── internal/
│   ├── client/                # API 客户端
│   └── config/                # 配置管理
├── go.mod
└── main.go
```

### Cobra CLI 实现
```go
// cmd/root.go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "prism",
    Short: "Prism 代理客户端命令行工具",
    Long:  `Prism CLI 是用于管理 Prism 代理服务的命令行工具。`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件 (默认 $HOME/.prism.yaml)")
    rootCmd.PersistentFlags().StringP("server", "s", "http://localhost:9090", "Prism Core 服务器地址")
    viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)
        
        viper.AddConfigPath(home)
        viper.SetConfigType("yaml")
        viper.SetConfigName(".prism")
    }
    
    viper.AutomaticEnv()
    viper.ReadInConfig()
}
```

## 跨平台构建

### 自动化构建脚本
```bash
#!/bin/bash
# scripts/build.sh

set -e

VERSION=${VERSION:-"v1.0.0"}
PLATFORMS="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64"

echo "构建 Prism $VERSION"

# 构建 Core
echo "构建 Core..."
cd core
for platform in $PLATFORMS; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    echo "构建 $GOOS/$GOARCH"
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=$VERSION" -o ../dist/prism-core-$GOOS-$GOARCH cmd/prism-core/main.go
    if [ $GOOS = "windows" ]; then
        mv ../dist/prism-core-$GOOS-$GOARCH ../dist/prism-core-$GOOS-$GOARCH.exe
    fi
done
cd ..

# 构建桌面客户端
echo "构建桌面客户端..."
cd desktop
npm install
npm run tauri build
cd ..

# 构建 Android 应用
echo "构建 Android 应用..."
cd android
./gradlew assembleRelease
cd ..

# 构建 CLI
echo "构建 CLI..."
cd cli
for platform in $PLATFORMS; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    echo "构建 CLI $GOOS/$GOARCH"
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=$VERSION" -o ../dist/prism-cli-$GOOS-$GOARCH main.go
    if [ $GOOS = "windows" ]; then
        mv ../dist/prism-cli-$GOOS-$GOARCH ../dist/prism-cli-$GOOS-$GOARCH.exe
    fi
done
cd ..

echo "构建完成！文件位于 dist/ 目录"
```

## 测试策略

### 单元测试
```bash
# Go 测试
go test -v ./...

# TypeScript 测试
npm test

# Android 测试
./gradlew test
```

### 集成测试
```bash
# 启动测试环境
docker-compose -f docker-compose.test.yml up -d

# 运行集成测试
npm run test:integration
```

这份开发指南涵盖了 Prism 项目各个组件的开发细节，为团队成员提供了完整的开发参考。