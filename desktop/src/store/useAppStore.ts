import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { invoke } from '@tauri-apps/api/tauri'

export interface SystemStatus {
  version: string
  uptime: number
  proxy_status: 'running' | 'stopped' | 'error'
  memory_usage: number
  cpu_usage: number
  total_upload: number
  total_download: number
}

export interface NodePool {
  id: string
  name: string
  description: string
  node_count: number
  active_nodes: number
  status: 'active' | 'inactive' | 'error'
  created_at: string
  updated_at: string
}

export interface ProxyNode {
  id: string
  name: string
  server: string
  port: number
  protocol: string
  country: string
  city: string
  delay: number
  status: 'online' | 'offline' | 'testing'
  last_test: string
}

export interface AppState {
  // 系统状态
  systemStatus: SystemStatus | null
  isConnected: boolean
  isDarkMode: boolean
  
  // 代理相关
  currentNode: ProxyNode | null
  proxyMode: 'rule' | 'global' | 'direct'
  
  // 数据
  nodePools: NodePool[]
  selectedNodePool: string | null
  
  // UI 状态
  sidebarCollapsed: boolean
  loading: boolean
  
  // Actions
  setSystemStatus: (status: SystemStatus) => void
  setConnected: (connected: boolean) => void
  toggleDarkMode: () => void
  setCurrentNode: (node: ProxyNode | null) => void
  setProxyMode: (mode: 'rule' | 'global' | 'direct') => void
  setNodePools: (pools: NodePool[]) => void
  setSelectedNodePool: (poolId: string | null) => void
  setSidebarCollapsed: (collapsed: boolean) => void
  setLoading: (loading: boolean) => void
  
  // 异步 Actions
  initializeApp: () => Promise<void>
  checkCoreConnection: () => Promise<boolean>
  fetchSystemStatus: () => Promise<void>
  fetchNodePools: () => Promise<void>
}

export const useAppStore = create<AppState>()(
  persist(
    (set, get) => ({
      // 初始状态
      systemStatus: null,
      isConnected: false,
      isDarkMode: false,
      currentNode: null,
      proxyMode: 'rule',
      nodePools: [],
      selectedNodePool: null,
      sidebarCollapsed: false,
      loading: false,
      
      // 同步 Actions
      setSystemStatus: (status) => set({ systemStatus: status }),
      
      setConnected: (connected) => set({ isConnected: connected }),
      
      toggleDarkMode: () => {
        const { isDarkMode } = get()
        const newMode = !isDarkMode
        set({ isDarkMode: newMode })
        document.body.setAttribute('data-theme', newMode ? 'dark' : 'light')
      },
      
      setCurrentNode: (node) => set({ currentNode: node }),
      
      setProxyMode: (mode) => set({ proxyMode: mode }),
      
      setNodePools: (pools) => set({ nodePools: pools }),
      
      setSelectedNodePool: (poolId) => set({ selectedNodePool: poolId }),
      
      setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
      
      setLoading: (loading) => set({ loading: loading }),
      
      // 异步 Actions
      initializeApp: async () => {
        const { isDarkMode, checkCoreConnection, fetchSystemStatus } = get()
        
        // 设置主题
        document.body.setAttribute('data-theme', isDarkMode ? 'dark' : 'light')
        
        // 检查核心连接
        const connected = await checkCoreConnection()
        set({ isConnected: connected })
        
        // 如果连接成功，获取系统状态
        if (connected) {
          await fetchSystemStatus()
        }
      },
      
      checkCoreConnection: async () => {
        try {
          const result = await invoke<boolean>('check_core_connection', {
            url: 'http://localhost:9090'
          })
          set({ isConnected: result })
          return result
        } catch (error) {
          console.error('Failed to check core connection:', error)
          set({ isConnected: false })
          return false
        }
      },
      
      fetchSystemStatus: async () => {
        try {
          set({ loading: true })
          // 这里应该调用实际的 API
          // const response = await apiClient.get('/system/status')
          // set({ systemStatus: response.data.data })
          
          // 模拟数据
          const mockStatus: SystemStatus = {
            version: '1.0.0',
            uptime: Date.now() - 1000000,
            proxy_status: 'running',
            memory_usage: 128.5,
            cpu_usage: 15.2,
            total_upload: 1048576,
            total_download: 2097152,
          }
          set({ systemStatus: mockStatus })
        } catch (error) {
          console.error('Failed to fetch system status:', error)
        } finally {
          set({ loading: false })
        }
      },
      
      fetchNodePools: async () => {
        try {
          set({ loading: true })
          // 这里应该调用实际的 API
          // const response = await apiClient.get('/nodepools')
          // set({ nodePools: response.data.data })
          
          // 模拟数据
          const mockPools: NodePool[] = [
            {
              id: 'pool-1',
              name: '高速节点池',
              description: '优质高速节点',
              node_count: 50,
              active_nodes: 45,
              status: 'active',
              created_at: '2024-01-01T00:00:00Z',
              updated_at: '2024-01-10T00:00:00Z',
            },
            {
              id: 'pool-2',
              name: '备用节点池',
              description: '备用节点池',
              node_count: 30,
              active_nodes: 28,
              status: 'active',
              created_at: '2024-01-01T00:00:00Z',
              updated_at: '2024-01-10T00:00:00Z',
            },
          ]
          set({ nodePools: mockPools })
        } catch (error) {
          console.error('Failed to fetch node pools:', error)
        } finally {
          set({ loading: false })
        }
      },
    }),
    {
      name: 'prism-app-store',
      partialize: (state) => ({
        isDarkMode: state.isDarkMode,
        proxyMode: state.proxyMode,
        selectedNodePool: state.selectedNodePool,
        sidebarCollapsed: state.sidebarCollapsed,
      }),
    }
  )
)