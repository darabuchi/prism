# Web 前端

## 技术栈

- **框架**: Next.js 14+ (App Router)
- **语言**: TypeScript
- **包管理器**: pnpm
- **UI 组件库**:
  - shadcn/ui（基于 Radix UI 和 Tailwind CSS 的组件系统）
  - Ant Design 5.x（基础组件）
  - ProComponents（优先使用，高级业务组件）
- **样式方案**: Tailwind CSS
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

## UI 组件使用优先级

### 1. shadcn/ui（最优先）

用于通用 UI 组件，提供完全可定制的组件：

```tsx
// 示例：使用 shadcn/ui 组件
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
```

**特点**:
- ✅ 基于 Radix UI（无障碍性强）
- ✅ 使用 Tailwind CSS（高度可定制）
- ✅ 组件代码直接复制到项目中（完全控制）
- ✅ 支持暗色模式

**安装组件**:
```bash
pnpm dlx shadcn-ui@latest add button card dialog tabs
```

### 2. ProComponents（次优先）

用于复杂的业务场景和数据展示：

```tsx
// 示例：使用 ProComponents
import { ProTable } from '@ant-design/pro-table';
import { ProForm, ProFormText } from '@ant-design/pro-form';
import { ProLayout } from '@ant-design/pro-layout';
```

**适用场景**:
- ✅ 数据表格：使用 `ProTable` 代替 `Table`
- ✅ 表单：使用 `ProForm` 代替 `Form`
- ✅ 页面布局：使用 `ProLayout`
- ✅ 卡片：使用 `ProCard`

参考文档：https://procomponents.ant.design/

### 3. Ant Design（基础组件）

仅在 shadcn/ui 和 ProComponents 都不适用时使用：

```tsx
import { message, notification } from 'antd';
```

## shadcn/ui 配置

### 初始化

```bash
cd web
pnpm dlx shadcn-ui@latest init
```

### components.json 配置

```json
{
  "style": "default",
  "rsc": true,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.ts",
    "css": "app/globals.css",
    "baseColor": "slate",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils"
  }
}
```

### 常用组件列表

```bash
# 布局组件
pnpm dlx shadcn-ui@latest add card separator

# 表单组件
pnpm dlx shadcn-ui@latest add button input label form select checkbox

# 反馈组件
pnpm dlx shadcn-ui@latest add alert dialog toast

# 导航组件
pnpm dlx shadcn-ui@latest add tabs navigation-menu

# 数据展示
pnpm dlx shadcn-ui@latest add table badge avatar
```

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
