# 前端共享代码

## 说明

此目录包含 Web 和 Desktop 应用共享的前端代码，使用 TypeScript 编写。

## 目录结构

```
shared/
├── types/        # TypeScript 类型定义
├── utils/        # 工具函数
├── constants/    # 常量定义
├── hooks/        # React Hooks（仅 React 项目使用）
├── package.json  # 依赖配置
└── tsconfig.json # TypeScript 配置
```

## 设计原则

1. **平台无关**: 不依赖特定框架或平台特性
2. **类型安全**: 使用 TypeScript，提供完整的类型定义
3. **可复用**: 代码可被 Web 和 Desktop 项目直接引用
4. **独立性**: 不依赖业务逻辑，保持纯粹的工具性质

## 使用方式

### 作为 pnpm workspace 包

在项目根目录的 `pnpm-workspace.yaml` 中配置：

```yaml
packages:
  - 'web'
  - 'desktop'
  - 'shared'
```

### 在其他项目中引用

```json
// web/package.json 或 desktop/package.json
{
  "dependencies": {
    "@prism/shared": "workspace:*"
  }
}
```

### 导入使用

```typescript
// 导入类型
import type { Subscription, Node } from '@prism/shared/types';

// 导入工具函数
import { formatBytes, formatDate } from '@prism/shared/utils';

// 导入常量
import { API_BASE_URL } from '@prism/shared/constants';

// 导入 Hooks
import { useLocalStorage } from '@prism/shared/hooks';
```

## 包含内容

### types/
TypeScript 类型定义：
- `subscription.ts`: 订阅相关类型
- `node.ts`: 节点相关类型
- `config.ts`: 配置相关类型
- `index.ts`: 统一导出

### utils/
工具函数：
- `format.ts`: 格式化工具（日期、大小、速度等）
- `validate.ts`: 验证工具（URL、配置等）
- `storage.ts`: 存储工具（localStorage、sessionStorage）

### constants/
常量定义：
- `api.ts`: API 端点常量
- `config.ts`: 配置常量

### hooks/
React Hooks（仅在 React 项目中使用）：
- `useLocalStorage.ts`: localStorage Hook
- `useDebounce.ts`: 防抖 Hook
- `useThrottle.ts`: 节流 Hook

## 示例

### 类型定义示例

```typescript
// types/subscription.ts
export interface Subscription {
  id: number;
  url: string;
  title: string;
  created_at: number;
  updated_at: number;
}

export interface Node {
  id: number;
  name: string;
  type: string;
  latency: number;
}
```

### 工具函数示例

```typescript
// utils/format.ts
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}
```

## 开发规范

1. 所有代码必须使用 TypeScript
2. 提供完整的类型定义
3. 添加 JSDoc 注释
4. 编写单元测试
5. 遵循项目代码规范
