import React, { useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { ConfigProvider, theme } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { useAppStore } from './store/useAppStore'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import NodePools from './pages/NodePools'
import Rules from './pages/Rules'
import Settings from './pages/Settings'
import './App.css'

const App: React.FC = () => {
  const { isDarkMode, initializeApp } = useAppStore()

  useEffect(() => {
    initializeApp()
  }, [initializeApp])

  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        algorithm: isDarkMode ? theme.darkAlgorithm : theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1890ff',
        },
      }}
    >
      <Router>
        <Layout>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/nodepools" element={<NodePools />} />
            <Route path="/rules" element={<Rules />} />
            <Route path="/settings" element={<Settings />} />
          </Routes>
        </Layout>
      </Router>
    </ConfigProvider>
  )
}

export default App