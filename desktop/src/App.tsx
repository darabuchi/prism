import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { Layout } from 'antd';

import Sidebar from './components/Sidebar';
import Dashboard from './pages/Dashboard';
import Subscriptions from './pages/Subscriptions';
import Nodes from './pages/Nodes';
import Settings from './pages/Settings';

const { Header, Sider, Content } = Layout;

const App: React.FC = () => {
  return (
    <Layout className="full-height">
      <Sider width={200} theme="light">
        <Sidebar />
      </Sider>
      <Layout>
        <Content style={{ margin: 0, minHeight: '100vh' }}>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/subscriptions" element={<Subscriptions />} />
            <Route path="/nodes" element={<Nodes />} />
            <Route path="/settings" element={<Settings />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  );
};

export default App;