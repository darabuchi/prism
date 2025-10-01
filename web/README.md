# Web 前端

## 技术栈

- **框架**: Next.js 14+ (App Router)
- **语言**: TypeScript
- **包管理器**: pnpm
- **UI 组件库**:
  - Ant Design 5.x（基础组件）
  - ProComponents（优先使用，高级业务组件）
- **状态管理**: Zustand / Redux Toolkit
- **数据请求**: SWR / TanStack Query

## 目录说明

### frontend/
用户前端界面，提供：
- 订阅管理
- 节点查看与测试
- 实时监控
- 使用统计

### admin/
管理后台界面，提供：
- 系统配置
- 用户管理
- 日志查看
- 性能监控
- 数据统计

### shared/
前端共享代码：
- TypeScript 类型定义
- 工具函数
- 常量定义
- 通用组件

## 开发

```bash
# 安装依赖
cd frontend && pnpm install
cd admin && pnpm install

# 开发模式
cd frontend && pnpm dev      # 用户前端
cd admin && pnpm dev          # 管理后台

# 构建生产版本
pnpm build

# 启动生产服务器
pnpm start
```

## ProComponents 使用

优先使用 ProComponents 而非基础的 Ant Design 组件：

- 使用 `ProTable` 代替 `Table`
- 使用 `ProForm` 代替 `Form`
- 使用 `ProLayout` 代替自定义布局
- 使用 `ProCard` 代替 `Card`

参考文档：https://procomponents.ant.design/
