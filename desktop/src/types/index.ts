// API 响应类型
export interface APIResponse<T = any> {
  code: number;
  message: string;
  data?: T;
  timestamp: number;
}

// 订阅相关类型
export interface Subscription {
  id: number;
  name: string;
  url: string;
  user_agent: string;
  auto_update: boolean;
  update_interval: number;
  total_nodes: number;
  active_nodes: number;
  unique_new_nodes: number;
  status: 'active' | 'inactive' | 'error';
  last_update?: string;
  last_success?: string;
  error_message: string;
  error_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateSubscriptionRequest {
  name: string;
  url: string;
  user_agent?: string;
  auto_update?: boolean;
  update_interval?: number;
  node_pool_ids?: number[];
}

export interface UpdateSubscriptionRequest {
  name?: string;
  user_agent?: string;
  auto_update?: boolean;
  update_interval?: number;
  status?: string;
}

export interface ListSubscriptionsResponse {
  total: number;
  page: number;
  size: number;
  subscriptions: Subscription[];
}

export interface UpdateResult {
  subscription_id: number;
  start_time: string;
  end_time: string;
  duration: number;
  total_fetched: number;
  valid_nodes: number;
  new_nodes: number;
  global_new_nodes: number;
  updated_nodes: number;
  removed_nodes: number;
}

export interface SubscriptionStats {
  subscription_id: number;
  total_nodes: number;
  active_nodes: number;
  survival_rate: number;
  protocol_distribution: Array<{
    protocol: string;
    count: number;
  }>;
  country_distribution: Array<{
    country: string;
    count: number;
  }>;
  recent_logs: SubscriptionLog[];
}

export interface SubscriptionLog {
  id: number;
  subscription_id: number;
  update_type: string;
  success: boolean;
  total_fetched: number;
  valid_nodes: number;
  new_nodes: number;
  global_new_nodes: number;
  updated_nodes: number;
  removed_nodes: number;
  error_message: string;
  http_status?: number;
  response_time?: number;
  created_at: string;
}

// 节点相关类型
export interface Node {
  id: number;
  subscription_id: number;
  node_pool_id?: number;
  name: string;
  hash: string;
  server: string;
  port: number;
  protocol: string;
  clash_config: Record<string, any>;
  country?: string;
  country_name?: string;
  city?: string;
  isp?: string;
  delay?: number;
  upload_speed?: number;
  download_speed?: number;
  loss_rate?: number;
  status: 'online' | 'offline' | 'testing' | 'unknown';
  last_test?: string;
  last_online?: string;
  continuous_failures: number;
  streaming_unlock?: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface ListNodesResponse {
  total: number;
  page: number;
  size: number;
  nodes: Node[];
}

// 节点池相关类型
export interface NodePool {
  id: number;
  name: string;
  description: string;
  total_subscriptions: number;
  total_nodes: number;
  active_nodes: number;
  survival_rate: number;
  enabled: boolean;
  priority: number;
  created_at: string;
  updated_at: string;
}

// 统计相关类型
export interface OverviewStats {
  total_subscriptions: number;
  active_subscriptions: number;
  total_node_pools: number;
  total_nodes: number;
  active_nodes: number;
  overall_survival_rate: number;
  total_tests_today: number;
  successful_tests_today: number;
}

export interface GeoDistribution {
  country: string;
  country_name: string;
  node_count: number;
  active_count: number;
}

export interface ProtocolDistribution {
  protocol: string;
  node_count: number;
  active_count: number;
}

// 通用类型
export interface PaginationParams {
  page?: number;
  size?: number;
}

export interface FilterParams {
  status?: string;
  country?: string;
  protocol?: string;
  auto_update?: boolean;
}

// 状态类型
export type LoadingState = boolean;
export type ErrorState = string | null;