# 桌面应用

## 技术栈

- **后端**: Tauri 1.x (Rust)
- **前端**: Next.js + TypeScript
- **UI 库**:
  - shadcn/ui（基于 Radix UI 和 Tailwind CSS）
  - Ant Design + ProComponents
- **样式方案**: Tailwind CSS
- **构建工具**: Next.js
- **包管理**: pnpm (前端) + Cargo (Rust)

## 特性

- ✅ 跨平台支持（Windows、macOS、Linux）
- ✅ 系统托盘集成
- ✅ 开机自启动
- ✅ 本地配置管理
- ✅ 与后端服务通信（通过 Tauri Commands）
- ✅ 轻量级（相比 Electron 更小的安装包）

## 目录结构

```
desktop/
├── src-tauri/          # Tauri Rust 后端
│   ├── src/           # Rust 源码
│   ├── icons/         # 应用图标
│   ├── Cargo.toml     # Rust 依赖
│   └── tauri.conf.json # Tauri 配置
├── src/               # React 前端
├── public/            # 静态资源
└── package.json       # 前端依赖
```

## 开发

```bash
# 安装依赖
pnpm install

# 开发模式
pnpm tauri dev

# 构建应用
pnpm tauri build
```

## 构建特定平台

```bash
# Windows
pnpm tauri build --target x86_64-pc-windows-msvc

# macOS Intel
pnpm tauri build --target x86_64-apple-darwin

# macOS Apple Silicon
pnpm tauri build --target aarch64-apple-darwin

# Linux
pnpm tauri build --target x86_64-unknown-linux-gnu
```

## Tauri Commands

Rust 后端通过 Tauri Commands 与前端通信：

```rust
// src-tauri/src/commands.rs
#[tauri::command]
fn get_subscriptions() -> Result<Vec<Subscription>, String> {
    // 实现
}
```

```typescript
// 前端调用
import { invoke } from '@tauri-apps/api/tauri';

const subscriptions = await invoke('get_subscriptions');
```

## 参考文档

- Tauri 官方文档: https://tauri.app/
- Tauri API: https://tauri.app/v1/api/js/
