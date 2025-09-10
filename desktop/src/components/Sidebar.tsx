import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Menu, Typography } from 'antd';
import {
  DashboardOutlined,
  LinkOutlined,
  NodeExpandOutlined,
  SettingOutlined,
  AppstoreOutlined,
} from '@ant-design/icons';

const { Title } = Typography;

const Sidebar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: '仪表盘',
    },
    {
      key: '/subscriptions',
      icon: <LinkOutlined />,
      label: '订阅管理',
    },
    {
      key: '/nodes',
      icon: <NodeExpandOutlined />,
      label: '节点管理',
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: '设置',
    },
  ];

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
  };

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <div style={{ padding: '16px', borderBottom: '1px solid #f0f0f0' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <AppstoreOutlined style={{ fontSize: '24px', color: '#1677ff' }} />
          <Title level={4} style={{ margin: 0, color: '#1677ff' }}>
            Prism
          </Title>
        </div>
      </div>
      <div style={{ flex: 1, paddingTop: '8px' }}>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{ border: 'none' }}
        />
      </div>
    </div>
  );
};

export default Sidebar;