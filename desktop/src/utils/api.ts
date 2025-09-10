import axios, { AxiosResponse, AxiosError } from 'axios';
import { message } from 'antd';
import { 
  APIResponse, 
  Subscription, 
  CreateSubscriptionRequest,
  UpdateSubscriptionRequest,
  ListSubscriptionsResponse,
  UpdateResult,
  SubscriptionStats,
  SubscriptionLog,
  Node,
  ListNodesResponse,
  NodePool,
  OverviewStats,
  GeoDistribution,
  ProtocolDistribution
} from '@/types';

// 创建 axios 实例
const api = axios.create({
  baseURL: 'http://localhost:9090/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 可以在这里添加 token 等认证信息
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse<APIResponse>) => {
    const { data } = response;
    if (data.code !== 0) {
      message.error(data.message || '请求失败');
      return Promise.reject(new Error(data.message || '请求失败'));
    }
    return response;
  },
  (error: AxiosError<APIResponse>) => {
    if (error.response?.data?.message) {
      message.error(error.response.data.message);
    } else if (error.message) {
      message.error(`网络错误: ${error.message}`);
    } else {
      message.error('未知错误');
    }
    return Promise.reject(error);
  }
);

// 订阅管理 API
export const subscriptionAPI = {
  // 获取订阅列表
  getList: (params?: {
    page?: number;
    size?: number;
    status?: string;
    auto_update?: boolean;
  }) => {
    return api.get<APIResponse<ListSubscriptionsResponse>>('/subscriptions', { params });
  },

  // 获取单个订阅
  getById: (id: number) => {
    return api.get<APIResponse<Subscription>>(`/subscriptions/${id}`);
  },

  // 创建订阅
  create: (data: CreateSubscriptionRequest) => {
    return api.post<APIResponse<Subscription>>('/subscriptions', data);
  },

  // 更新订阅
  update: (id: number, data: UpdateSubscriptionRequest) => {
    return api.put<APIResponse<Subscription>>(`/subscriptions/${id}`, data);
  },

  // 删除订阅
  delete: (id: number) => {
    return api.delete<APIResponse>(`/subscriptions/${id}`);
  },

  // 启用订阅
  enable: (id: number) => {
    return api.post<APIResponse<Subscription>>(`/subscriptions/${id}/enable`);
  },

  // 禁用订阅
  disable: (id: number) => {
    return api.post<APIResponse<Subscription>>(`/subscriptions/${id}/disable`);
  },

  // 手动更新订阅
  update_content: (id: number) => {
    return api.post<APIResponse<UpdateResult>>(`/subscriptions/${id}/update`);
  },

  // 获取订阅统计
  getStats: (id: number) => {
    return api.get<APIResponse<SubscriptionStats>>(`/subscriptions/${id}/stats`);
  },

  // 获取订阅日志
  getLogs: (id: number, params?: {
    page?: number;
    size?: number;
    success?: boolean;
    update_type?: string;
  }) => {
    return api.get<APIResponse<{ logs: SubscriptionLog[] }>>(`/subscriptions/${id}/logs`, { params });
  },

  // 导入订阅
  import: (data: { subscriptions: CreateSubscriptionRequest[] }) => {
    return api.post<APIResponse>('/subscriptions/import', data);
  },

  // 导出订阅
  export: (format?: 'json' | 'yaml') => {
    return api.get<APIResponse>('/subscriptions/export', { params: { format } });
  },
};

// 节点管理 API
export const nodeAPI = {
  // 获取节点列表
  getList: (params?: {
    page?: number;
    size?: number;
    subscription_id?: number;
    node_pool_id?: number;
    country?: string;
    protocol?: string;
    status?: string;
    sort?: string;
    order?: string;
  }) => {
    return api.get<APIResponse<ListNodesResponse>>('/nodes', { params });
  },

  // 获取单个节点
  getById: (id: number) => {
    return api.get<APIResponse<Node>>(`/nodes/${id}`);
  },

  // 测试单个节点
  test: (id: number, data: {
    test_types: string[];
    test_config?: Record<string, any>;
  }) => {
    return api.post<APIResponse>(`/nodes/${id}/test`, data);
  },

  // 批量测试节点
  batchTest: (data: {
    node_ids: number[];
    test_types: string[];
    test_config?: Record<string, any>;
  }) => {
    return api.post<APIResponse<{ task_id: string; total_nodes: number; status: string }>>('/nodes/batch-test', data);
  },

  // 获取测试任务状态
  getTestTaskStatus: (taskId: string) => {
    return api.get<APIResponse>(`/nodes/test-tasks/${taskId}`);
  },

  // 获取节点测试历史
  getTestHistory: (id: number, params?: {
    test_type?: string;
    start_time?: string;
    end_time?: string;
    limit?: number;
  }) => {
    return api.get<APIResponse>(`/nodes/${id}/test-history`, { params });
  },

  // 智能节点选择
  getBestSelection: (params?: {
    node_pool_id?: number;
    country?: string;
    protocol?: string;
    streaming_unlock?: string[];
    count?: number;
  }) => {
    return api.get<APIResponse>('/nodes/best-selection', { params });
  },
};

// 节点池管理 API
export const nodePoolAPI = {
  // 获取节点池列表
  getList: () => {
    return api.get<APIResponse<{ node_pools: NodePool[] }>>('/nodepools');
  },

  // 获取单个节点池
  getById: (id: number) => {
    return api.get<APIResponse<NodePool>>(`/nodepools/${id}`);
  },

  // 创建节点池
  create: (data: {
    name: string;
    description?: string;
    enabled?: boolean;
    priority?: number;
  }) => {
    return api.post<APIResponse<NodePool>>('/nodepools', data);
  },

  // 更新节点池
  update: (id: number, data: {
    name?: string;
    description?: string;
    enabled?: boolean;
    priority?: number;
  }) => {
    return api.put<APIResponse<NodePool>>(`/nodepools/${id}`, data);
  },

  // 删除节点池
  delete: (id: number) => {
    return api.delete<APIResponse>(`/nodepools/${id}`);
  },

  // 关联订阅到节点池
  associateSubscriptions: (id: number, data: {
    subscription_ids: number[];
    enabled?: boolean;
    priority?: number;
  }) => {
    return api.post<APIResponse>(`/nodepools/${id}/subscriptions`, data);
  },
};

// 统计分析 API
export const statsAPI = {
  // 获取总体统计
  getOverview: () => {
    return api.get<APIResponse<OverviewStats>>('/stats/overview');
  },

  // 获取地区分布
  getGeoDistribution: () => {
    return api.get<APIResponse<GeoDistribution[]>>('/stats/geo-distribution');
  },

  // 获取协议分布
  getProtocolDistribution: () => {
    return api.get<APIResponse<ProtocolDistribution[]>>('/stats/protocol-distribution');
  },

  // 获取性能趋势
  getPerformanceTrend: (params: {
    period: 'hour' | 'day' | 'week' | 'month';
    node_pool_id?: number;
    country?: string;
  }) => {
    return api.get<APIResponse>('/stats/performance-trend', { params });
  },
};

// 系统管理 API
export const systemAPI = {
  // 获取系统状态
  getStatus: () => {
    return api.get<APIResponse>('/system/status');
  },

  // 获取系统配置
  getConfig: () => {
    return api.get<APIResponse>('/system/config');
  },

  // 更新系统配置
  updateConfig: (data: Record<string, any>) => {
    return api.put<APIResponse>('/system/config', data);
  },
};

// 任务管理 API
export const taskAPI = {
  // 获取自动更新任务状态
  getAutoUpdateStatus: () => {
    return api.get<APIResponse>('/tasks/auto-update');
  },

  // 触发自动更新
  triggerAutoUpdate: () => {
    return api.post<APIResponse>('/tasks/auto-update/trigger');
  },

  // 获取定时测试任务状态
  getScheduledTestStatus: () => {
    return api.get<APIResponse>('/tasks/scheduled-test');
  },
};

export default api;