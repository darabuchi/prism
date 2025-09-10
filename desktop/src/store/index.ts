import { create } from 'zustand';
import { subscript } from 'crypto';
import { 
  Subscription, 
  Node, 
  NodePool, 
  OverviewStats,
  LoadingState,
  ErrorState
} from '@/types';

// 全局状态接口
interface AppState {
  // 订阅状态
  subscriptions: Subscription[];
  subscriptionLoading: LoadingState;
  subscriptionError: ErrorState;
  
  // 节点状态
  nodes: Node[];
  nodeLoading: LoadingState;
  nodeError: ErrorState;
  
  // 节点池状态
  nodePools: NodePool[];
  nodePoolLoading: LoadingState;
  nodePoolError: ErrorState;
  
  // 统计状态
  stats: OverviewStats | null;
  statsLoading: LoadingState;
  statsError: ErrorState;
  
  // 通用操作
  setLoading: (key: string, loading: boolean) => void;
  setError: (key: string, error: string | null) => void;
  
  // 订阅操作
  setSubscriptions: (subscriptions: Subscription[]) => void;
  addSubscription: (subscription: Subscription) => void;
  updateSubscription: (id: number, updates: Partial<Subscription>) => void;
  removeSubscription: (id: number) => void;
  
  // 节点操作
  setNodes: (nodes: Node[]) => void;
  updateNode: (id: number, updates: Partial<Node>) => void;
  
  // 节点池操作
  setNodePools: (nodePools: NodePool[]) => void;
  addNodePool: (nodePool: NodePool) => void;
  updateNodePool: (id: number, updates: Partial<NodePool>) => void;
  removeNodePool: (id: number) => void;
  
  // 统计操作
  setStats: (stats: OverviewStats) => void;
}

// 创建状态存储
export const useAppStore = create<AppState>((set, get) => ({
  // 初始状态
  subscriptions: [],
  subscriptionLoading: false,
  subscriptionError: null,
  
  nodes: [],
  nodeLoading: false,
  nodeError: null,
  
  nodePools: [],
  nodePoolLoading: false,
  nodePoolError: null,
  
  stats: null,
  statsLoading: false,
  statsError: null,
  
  // 通用操作
  setLoading: (key: string, loading: boolean) => {
    set({ [`${key}Loading`]: loading } as any);
  },
  
  setError: (key: string, error: string | null) => {
    set({ [`${key}Error`]: error } as any);
  },
  
  // 订阅操作
  setSubscriptions: (subscriptions: Subscription[]) => {
    set({ subscriptions });
  },
  
  addSubscription: (subscription: Subscription) => {
    const { subscriptions } = get();
    set({ subscriptions: [subscription, ...subscriptions] });
  },
  
  updateSubscription: (id: number, updates: Partial<Subscription>) => {
    const { subscriptions } = get();
    set({
      subscriptions: subscriptions.map(sub => 
        sub.id === id ? { ...sub, ...updates } : sub
      )
    });
  },
  
  removeSubscription: (id: number) => {
    const { subscriptions } = get();
    set({ subscriptions: subscriptions.filter(sub => sub.id !== id) });
  },
  
  // 节点操作
  setNodes: (nodes: Node[]) => {
    set({ nodes });
  },
  
  updateNode: (id: number, updates: Partial<Node>) => {
    const { nodes } = get();
    set({
      nodes: nodes.map(node => 
        node.id === id ? { ...node, ...updates } : node
      )
    });
  },
  
  // 节点池操作
  setNodePools: (nodePools: NodePool[]) => {
    set({ nodePools });
  },
  
  addNodePool: (nodePool: NodePool) => {
    const { nodePools } = get();
    set({ nodePools: [nodePool, ...nodePools] });
  },
  
  updateNodePool: (id: number, updates: Partial<NodePool>) => {
    const { nodePools } = get();
    set({
      nodePools: nodePools.map(pool => 
        pool.id === id ? { ...pool, ...updates } : pool
      )
    });
  },
  
  removeNodePool: (id: number) => {
    const { nodePools } = get();
    set({ nodePools: nodePools.filter(pool => pool.id !== id) });
  },
  
  // 统计操作
  setStats: (stats: OverviewStats) => {
    set({ stats });
  },
}));

// 订阅相关的 hooks
export const useSubscriptions = () => {
  const subscriptions = useAppStore(state => state.subscriptions);
  const loading = useAppStore(state => state.subscriptionLoading);
  const error = useAppStore(state => state.subscriptionError);
  
  const setSubscriptions = useAppStore(state => state.setSubscriptions);
  const addSubscription = useAppStore(state => state.addSubscription);
  const updateSubscription = useAppStore(state => state.updateSubscription);
  const removeSubscription = useAppStore(state => state.removeSubscription);
  const setLoading = useAppStore(state => state.setLoading);
  const setError = useAppStore(state => state.setError);
  
  return {
    subscriptions,
    loading,
    error,
    setSubscriptions,
    addSubscription,
    updateSubscription,
    removeSubscription,
    setLoading: (loading: boolean) => setLoading('subscription', loading),
    setError: (error: string | null) => setError('subscription', error),
  };
};

// 节点相关的 hooks
export const useNodes = () => {
  const nodes = useAppStore(state => state.nodes);
  const loading = useAppStore(state => state.nodeLoading);
  const error = useAppStore(state => state.nodeError);
  
  const setNodes = useAppStore(state => state.setNodes);
  const updateNode = useAppStore(state => state.updateNode);
  const setLoading = useAppStore(state => state.setLoading);
  const setError = useAppStore(state => state.setError);
  
  return {
    nodes,
    loading,
    error,
    setNodes,
    updateNode,
    setLoading: (loading: boolean) => setLoading('node', loading),
    setError: (error: string | null) => setError('node', error),
  };
};

// 节点池相关的 hooks
export const useNodePools = () => {
  const nodePools = useAppStore(state => state.nodePools);
  const loading = useAppStore(state => state.nodePoolLoading);
  const error = useAppStore(state => state.nodePoolError);
  
  const setNodePools = useAppStore(state => state.setNodePools);
  const addNodePool = useAppStore(state => state.addNodePool);
  const updateNodePool = useAppStore(state => state.updateNodePool);
  const removeNodePool = useAppStore(state => state.removeNodePool);
  const setLoading = useAppStore(state => state.setLoading);
  const setError = useAppStore(state => state.setError);
  
  return {
    nodePools,
    loading,
    error,
    setNodePools,
    addNodePool,
    updateNodePool,
    removeNodePool,
    setLoading: (loading: boolean) => setLoading('nodePool', loading),
    setError: (error: string | null) => setError('nodePool', error),
  };
};

// 统计相关的 hooks
export const useStats = () => {
  const stats = useAppStore(state => state.stats);
  const loading = useAppStore(state => state.statsLoading);
  const error = useAppStore(state => state.statsError);
  
  const setStats = useAppStore(state => state.setStats);
  const setLoading = useAppStore(state => state.setLoading);
  const setError = useAppStore(state => state.setError);
  
  return {
    stats,
    loading,
    error,
    setStats,
    setLoading: (loading: boolean) => setLoading('stats', loading),
    setError: (error: string | null) => setError('stats', error),
  };
};