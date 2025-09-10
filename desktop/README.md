# Prism Desktop

基于 Tauri 的现代化代理管理工具桌面端应用。

## 技术栈

- **前端**: React 18 + TypeScript + Vite
- **UI 库**: Ant Design 5.x
- **状态管理**: Zustand
- **图表**: Recharts
- **桌面框架**: Tauri 1.5
- **构建工具**: Vite + Tauri CLI
- **开发工具**: ESLint + TypeScript + 热重载

## 功能特性

### 核心功能
- 🚀 **订阅管理** - 创建、编辑、删除、启用/禁用订阅
- 📊 **仪表盘** - 实时统计和性能监控
- 🌐 **节点管理** - 节点测试、筛选、详情查看
- ⚙️ **系统设置** - 配置管理和系统信息

### 桌面特性
- 🖥️ **系统托盘** - 后台运行，托盘控制
- 📱 **原生体验** - 跨平台原生应用
- 🔔 **系统通知** - 重要事件通知
- 🔄 **自动更新** - 应用自动更新

## 开发环境

### 系统要求
- Node.js 18+ 
- Rust 1.60+
- 操作系统: Windows 10+, macOS 10.15+, Linux (Ubuntu 18.04+)

### 安装依赖

```bash
# 安装前端依赖
npm install

# 安装 Tauri CLI
npm install --save-dev @tauri-apps/cli
```

### 开发模式

#### 1. Web 端开发（推荐用于样式调试）
```bash
# 启动纯 Web 开发服务器
npm run dev
```
- 访问 `http://localhost:1420` 查看 Web 版本
- 支持热重载，样式修改实时生效
- 可使用浏览器开发者工具调试
- 适合样式开发和响应式布局测试

#### 2. Tauri 桌面开发
```bash
# 启动 Tauri 开发模式
npm run tauri:dev
```
这将同时启动：
- Vite 开发服务器 (http://localhost:1420)
- Tauri 桌面应用窗口
- 可以对比 Web 版本和桌面版本的显示效果

#### 3. 预览构建版本
```bash
# 构建并预览生产版本
npm run build
npm run preview
```
- 查看生产环境构建后的效果
- 验证构建优化后的样式

### 样式调试指南

#### 实时样式开发
1. 运行 `npm run dev` 启动 Web 开发服务器
2. 在浏览器中访问 `http://localhost:1420`
3. 按 F12 打开开发者工具进行样式调试
4. 修改组件样式文件，页面会自动热重载

#### 桌面版本验证
1. 运行 `npm run tauri:dev` 启动桌面应用
2. 在桌面应用中右键点击 → "检查元素" 进行调试
3. 对比 Web 版本 (`localhost:1420`) 和桌面版本的显示差异

#### 响应式测试
- 在浏览器开发者工具中切换设备模拟器
- 测试不同屏幕尺寸下的布局效果
- 验证 Ant Design 组件的响应式行为

### 构建应用

```bash
# 构建前端资源
npm run build

# 构建桌面应用
npm run tauri:build
```

构建完成后，安装包将位于 `src-tauri/target/release/bundle/` 目录中。

## 项目结构

```
desktop/
├── src/                    # React 前端源码
│   ├── components/         # 通用组件
│   ├── pages/             # 页面组件
│   ├── store/             # 状态管理
│   ├── types/             # TypeScript 类型定义
│   ├── utils/             # 工具函数
│   ├── App.tsx            # 主应用组件
│   └── main.tsx           # 应用入口
├── src-tauri/             # Tauri 后端源码
│   ├── src/               # Rust 源码
│   ├── icons/             # 应用图标
│   ├── Cargo.toml         # Rust 依赖配置
│   └── tauri.conf.json    # Tauri 配置
├── public/                # 静态资源
├── dist/                  # 构建输出
├── package.json           # 项目配置
├── vite.config.ts         # Vite 配置
└── tsconfig.json          # TypeScript 配置
```

## 主要页面

### 仪表盘 (Dashboard)
- 系统概览统计
- 性能趋势图表
- 协议分布饼图
- 最近活动列表

### 订阅管理 (Subscriptions)
- 订阅列表展示
- CRUD 操作
- 批量导入导出
- 手动更新
- 统计信息和日志

### 节点管理 (Nodes)
- 节点列表和筛选
- 批量测试功能
- 节点详情抽屉
- 性能指标显示

### 系统设置 (Settings)
- 服务器配置
- 代理设置
- 日志配置
- 系统信息

## API 集成

应用通过 HTTP API 与 Prism Core 后端服务通信：

- **基础 URL**: `http://localhost:9090/api/v1`
- **请求拦截器**: 自动处理错误和认证
- **响应拦截器**: 统一错误提示和数据处理

主要 API 模块：
- `subscriptionAPI` - 订阅管理
- `nodeAPI` - 节点管理  
- `statsAPI` - 统计分析
- `systemAPI` - 系统管理

## 状态管理

使用 Zustand 进行状态管理，主要状态包括：

- **subscriptions** - 订阅数据和操作
- **nodes** - 节点数据和操作
- **nodePools** - 节点池数据和操作
- **stats** - 统计数据

每个状态模块都包含：
- 数据存储
- 加载状态
- 错误处理
- CRUD 操作方法

## 系统托盘功能

- **左键点击** - 显示/隐藏主窗口
- **右键菜单**:
  - 显示主窗口
  - 隐藏主窗口
  - 退出应用

## 开发指南

### 开发工作流程

#### 样式开发推荐流程
1. **启动开发环境**
   ```bash
   npm run dev  # 启动 Web 开发服务器
   ```

2. **样式调试**
   - 浏览器访问 `http://localhost:1420`
   - 使用浏览器开发者工具 (F12) 实时调试样式
   - 修改 CSS/组件文件，享受热重载体验

3. **桌面效果验证**
   ```bash
   npm run tauri:dev  # 启动桌面应用验证效果
   ```

4. **多端对比测试**
   - 同时打开 Web 版本和桌面版本
   - 对比不同平台下的显示效果
   - 确保样式在桌面 WebView 中正常显示

### 开发调试技巧

#### 浏览器开发者工具
- **Elements 面板**: 实时修改 CSS，查看样式层级
- **Console 面板**: 查看 JavaScript 错误和日志
- **Network 面板**: 监控 API 请求和响应
- **Application 面板**: 查看 LocalStorage 和状态

#### Tauri 桌面调试
- **右键菜单**: 在桌面应用中右键 → "检查元素"
- **WebView 调试**: 桌面应用内置 Chrome DevTools
- **Rust 日志**: 在终端中查看 Tauri 后端日志
- **系统托盘**: 测试托盘功能和窗口状态

#### 代码热重载
- **前端热重载**: 修改 React/TypeScript 文件自动刷新
- **样式热重载**: CSS 修改立即生效，无需刷新页面
- **Tauri 重载**: 修改 Rust 代码需要重启 `npm run tauri:dev`

### 添加新页面

1. 在 `src/pages/` 创建页面组件
2. 在 `src/App.tsx` 添加路由
3. 在 `src/components/Sidebar.tsx` 添加导航菜单
4. 使用 `npm run dev` 在浏览器中预览新页面

### 添加新 API

1. 在 `src/utils/api.ts` 添加 API 方法
2. 在 `src/types/index.ts` 定义数据类型
3. 在对应的 store 中添加状态管理
4. 在浏览器 Network 面板中测试 API 调用

### 添加 Tauri 命令

1. 在 `src-tauri/src/main.rs` 添加命令函数
2. 在 `invoke_handler` 中注册命令
3. 在前端通过 `invoke` 调用
4. 重启 `npm run tauri:dev` 使 Rust 代码生效

### 样式开发最佳实践

#### Ant Design 主题定制
- 使用 Ant Design 的主题配置系统
- 在 `vite.config.ts` 中配置主题变量
- 保持与设计规范的一致性

#### 响应式设计
- 使用 Ant Design 的栅格系统 (`Row`, `Col`)
- 利用 CSS Media Queries 适配不同屏幕
- 测试移动端和桌面端的显示效果

#### 组件样式管理
- 使用 CSS Modules 或 styled-components 避免样式冲突
- 保持组件样式的独立性和可重用性
- 遵循 BEM 命名规范或类似的样式约定

## 构建和发布

### 开发构建
```bash
npm run build
npm run tauri:build
```

### 发布准备
- 更新 `package.json` 和 `Cargo.toml` 中的版本号
- 准备应用图标 (存放在 `src-tauri/icons/`)
- 配置代码签名证书 (生产环境)

## 跨平台支持

- **Windows**: .msi 安装包
- **macOS**: .dmg 磁盘映像
- **Linux**: .deb/.rpm 包

每个平台都有对应的原生特性支持：
- Windows: 系统托盘、原生通知
- macOS: 菜单栏集成、通知中心
- Linux: 系统托盘、桌面通知

## 许可证

MIT License