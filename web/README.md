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

## 功能模块

### 订阅管理
- 添加订阅（支持多种格式：Clash、V2Ray、Surge）
- 编辑订阅信息
- 删除订阅
- 手动更新订阅
- 查看订阅详情

### 节点管理
- 查看所有节点列表
- 测试节点延迟
- 检测节点解锁能力
- 节点筛选与排序
- 速度图表展示

### 系统设置
- 代理配置
- 路由规则管理
- 主题设置（亮色/暗色）
- 语言切换

## 开发

```bash
# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 构建生产版本
pnpm build

# 启动生产服务器
pnpm start

# 类型检查
pnpm type-check

# 代码检查
pnpm lint
```

## 项目结构

```
web/
├── app/              # Next.js App Router
│   ├── layout.tsx    # 根布局
│   ├── page.tsx      # 首页
│   ├── subscription/ # 订阅管理页面
│   ├── node/         # 节点管理页面
│   ├── settings/     # 设置页面
│   └── api/          # API Routes
├── components/       # React 组件
│   ├── common/       # 通用组件
│   └── business/     # 业务组件
├── lib/              # 工具库
│   ├── api/          # API 客户端
│   ├── hooks/        # 自定义 Hooks
│   ├── store/        # 状态管理
│   ├── types/        # 类型定义
│   ├── utils/        # 工具函数
│   └── constants/    # 常量定义
└── public/           # 静态资源
```

## ProComponents 使用

优先使用 ProComponents 而非基础的 Ant Design 组件：

- ✅ 使用 `ProTable` 代替 `Table`
- ✅ 使用 `ProForm` 代替 `Form`
- ✅ 使用 `ProLayout` 代替自定义布局
- ✅ 使用 `ProCard` 代替 `Card`

参考文档：https://procomponents.ant.design/

## 共享代码

Web 应用使用 `@prism/shared` 包中的共享代码：

```typescript
import type { Subscription, Node } from '@prism/shared/types';
import { formatBytes, formatDate } from '@prism/shared/utils';
import { API_BASE_URL } from '@prism/shared/constants';
```

## 环境变量

创建 `.env.local` 文件：

```bash
# API 地址
NEXT_PUBLIC_API_URL=http://localhost:8080

# 其他配置...
```
