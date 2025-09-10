import React, { useEffect } from 'react';
import { Row, Col, Card, Statistic, Progress, List, Typography, Spin, Alert } from 'antd';
import { 
  LinkOutlined, 
  NodeExpandOutlined, 
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  ClockCircleOutlined,
  GlobalOutlined
} from '@ant-design/icons';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { useStats } from '@/store';
import { statsAPI } from '@/utils/api';
import { formatPercentage } from '@/utils/helpers';

const { Title, Text } = Typography;

// 模拟数据
const performanceData = [
  { time: '00:00', delay: 120, speed: 45 },
  { time: '04:00', delay: 110, speed: 52 },
  { time: '08:00', delay: 95, speed: 48 },
  { time: '12:00', delay: 88, speed: 55 },
  { time: '16:00', delay: 102, speed: 51 },
  { time: '20:00', delay: 115, speed: 47 },
];

const protocolData = [
  { name: 'VMess', value: 45, color: '#1677ff' },
  { name: 'Shadowsocks', value: 30, color: '#52c41a' },
  { name: 'Trojan', value: 15, color: '#faad14' },
  { name: 'Hysteria', value: 10, color: '#f5222d' },
];

const recentActivities = [
  { title: '订阅 "高速线路" 更新完成', time: '2 分钟前', type: 'success' },
  { title: '节点测试任务已完成', time: '5 分钟前', type: 'info' },
  { title: '发现 3 个新节点', time: '10 分钟前', type: 'success' },
  { title: '订阅 "备用线路" 更新失败', time: '15 分钟前', type: 'error' },
];

const Dashboard: React.FC = () => {
  const { stats, loading, error, setStats, setLoading, setError } = useStats();

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    setLoading(true);
    try {
      const response = await statsAPI.getOverview();
      setStats(response.data.data!);
    } catch (err) {
      setError('获取统计数据失败');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex-center" style={{ height: '100vh' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: '24px' }}>
        <Alert 
          message="数据加载失败" 
          description={error} 
          type="error" 
          showIcon 
          action={
            <button onClick={fetchStats}>重试</button>
          }
        />
      </div>
    );
  }

  return (
    <div style={{ padding: '24px' }}>
      <div style={{ marginBottom: '24px' }}>
        <Title level={2}>仪表盘</Title>
        <Text type="secondary">代理服务总览</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="总订阅数"
              value={stats?.total_subscriptions || 0}
              prefix={<LinkOutlined />}
              valueStyle={{ color: '#1677ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="总节点数"
              value={stats?.total_nodes || 0}
              prefix={<NodeExpandOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="在线节点"
              value={stats?.active_nodes || 0}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="存活率"
              value={stats?.overall_survival_rate || 0}
              precision={1}
              suffix="%"
              prefix={<GlobalOutlined />}
              valueStyle={{ color: '#f5222d' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 详细统计 */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={24} md={12}>
          <Card title="订阅状态">
            <Row gutter={[16, 16]}>
              <Col span={12}>
                <Statistic
                  title="活跃订阅"
                  value={stats?.active_subscriptions || 0}
                  valueStyle={{ color: '#52c41a' }}
                />
              </Col>
              <Col span={12}>
                <Progress 
                  type="circle" 
                  size={80}
                  percent={stats?.total_subscriptions ? 
                    Math.round((stats.active_subscriptions / stats.total_subscriptions) * 100) : 0
                  }
                  format={percent => `${percent}%`}
                />
              </Col>
            </Row>
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="今日测试">
            <Row gutter={[16, 16]}>
              <Col span={12}>
                <Statistic
                  title="总测试次数"
                  value={stats?.total_tests_today || 0}
                  valueStyle={{ color: '#1677ff' }}
                />
              </Col>
              <Col span={12}>
                <Statistic
                  title="成功率"
                  value={stats?.total_tests_today ? 
                    formatPercentage((stats.successful_tests_today / stats.total_tests_today) * 100) : '0%'
                  }
                  valueStyle={{ color: '#52c41a' }}
                />
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>

      {/* 图表 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card title="性能趋势" style={{ height: '400px' }}>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={performanceData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="time" />
                <YAxis yAxisId="delay" orientation="left" />
                <YAxis yAxisId="speed" orientation="right" />
                <Tooltip />
                <Area 
                  yAxisId="delay"
                  type="monotone" 
                  dataKey="delay" 
                  stackId="1"
                  stroke="#1677ff" 
                  fill="#1677ff" 
                  fillOpacity={0.3}
                  name="平均延迟(ms)"
                />
                <Area 
                  yAxisId="speed"
                  type="monotone" 
                  dataKey="speed" 
                  stackId="2"
                  stroke="#52c41a" 
                  fill="#52c41a" 
                  fillOpacity={0.3}
                  name="平均速度(Mbps)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        
        <Col xs={24} lg={8}>
          <Card title="协议分布" style={{ height: '400px' }}>
            <ResponsiveContainer width="100%" height={200}>
              <PieChart>
                <Pie
                  data={protocolData}
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                >
                  {protocolData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
            <div style={{ marginTop: '16px' }}>
              {protocolData.map(item => (
                <div key={item.name} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                  <span>
                    <span style={{ 
                      display: 'inline-block', 
                      width: '8px', 
                      height: '8px', 
                      backgroundColor: item.color,
                      marginRight: '8px',
                      borderRadius: '50%'
                    }} />
                    {item.name}
                  </span>
                  <span>{item.value}个</span>
                </div>
              ))}
            </div>
          </Card>
        </Col>
      </Row>

      {/* 最近活动 */}
      <Row style={{ marginTop: '24px' }}>
        <Col span={24}>
          <Card title="最近活动">
            <List
              dataSource={recentActivities}
              renderItem={item => (
                <List.Item>
                  <List.Item.Meta
                    avatar={
                      item.type === 'success' ? (
                        <CheckCircleOutlined style={{ color: '#52c41a' }} />
                      ) : item.type === 'error' ? (
                        <ExclamationCircleOutlined style={{ color: '#f5222d' }} />
                      ) : (
                        <ClockCircleOutlined style={{ color: '#1677ff' }} />
                      )
                    }
                    title={item.title}
                    description={item.time}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;