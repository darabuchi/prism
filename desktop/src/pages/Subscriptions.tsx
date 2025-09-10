import React, { useEffect, useState } from 'react';
import { 
  Table, 
  Button, 
  Card, 
  Space, 
  Tag, 
  Modal, 
  Form, 
  Input, 
  Switch, 
  InputNumber,
  message,
  Popconfirm,
  Tooltip,
  Progress,
  Drawer,
  Descriptions,
  List,
  Typography,
  Alert
} from 'antd';
import { 
  PlusOutlined,
  ReloadOutlined,
  EditOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  BarChartOutlined,
  HistoryOutlined,
  ImportOutlined,
  ExportOutlined,
  SyncOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useSubscriptions } from '@/store';
import { subscriptionAPI } from '@/utils/api';
import { formatTime, formatRelativeTime, getStatusColor } from '@/utils/helpers';
import StatusBadge from '@/components/StatusBadge';
import { Subscription, CreateSubscriptionRequest, SubscriptionStats, SubscriptionLog } from '@/types';

const { Title, Text } = Typography;
const { TextArea } = Input;

const Subscriptions: React.FC = () => {
  const {
    subscriptions,
    loading,
    error,
    setSubscriptions,
    addSubscription,
    updateSubscription,
    removeSubscription,
    setLoading,
    setError
  } = useSubscriptions();

  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [statsDrawerVisible, setStatsDrawerVisible] = useState(false);
  const [logsDrawerVisible, setLogsDrawerVisible] = useState(false);
  const [selectedSubscription, setSelectedSubscription] = useState<Subscription | null>(null);
  const [subscriptionStats, setSubscriptionStats] = useState<SubscriptionStats | null>(null);
  const [subscriptionLogs, setSubscriptionLogs] = useState<SubscriptionLog[]>([]);
  const [form] = Form.useForm();
  const [editForm] = Form.useForm();

  useEffect(() => {
    fetchSubscriptions();
  }, []);

  const fetchSubscriptions = async () => {
    setLoading(true);
    try {
      const response = await subscriptionAPI.getList();
      setSubscriptions(response.data.data!.subscriptions);
    } catch (err) {
      setError('获取订阅列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (values: CreateSubscriptionRequest) => {
    try {
      const response = await subscriptionAPI.create(values);
      addSubscription(response.data.data!);
      setCreateModalVisible(false);
      form.resetFields();
      message.success('创建订阅成功');
    } catch (err) {
      message.error('创建订阅失败');
    }
  };

  const handleEdit = async (values: any) => {
    if (!selectedSubscription) return;
    
    try {
      const response = await subscriptionAPI.update(selectedSubscription.id, values);
      updateSubscription(selectedSubscription.id, response.data.data!);
      setEditModalVisible(false);
      editForm.resetFields();
      setSelectedSubscription(null);
      message.success('更新订阅成功');
    } catch (err) {
      message.error('更新订阅失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await subscriptionAPI.delete(id);
      removeSubscription(id);
      message.success('删除订阅成功');
    } catch (err) {
      message.error('删除订阅失败');
    }
  };

  const handleToggleStatus = async (subscription: Subscription) => {
    try {
      if (subscription.status === 'active') {
        await subscriptionAPI.disable(subscription.id);
        updateSubscription(subscription.id, { status: 'inactive' });
        message.success('订阅已禁用');
      } else {
        await subscriptionAPI.enable(subscription.id);
        updateSubscription(subscription.id, { status: 'active' });
        message.success('订阅已启用');
      }
    } catch (err) {
      message.error('操作失败');
    }
  };

  const handleUpdateContent = async (id: number) => {
    try {
      message.loading('正在更新订阅...', 0);
      const response = await subscriptionAPI.update_content(id);
      const result = response.data.data!;
      
      message.destroy();
      message.success(`更新完成！新增 ${result.new_nodes} 个节点`);
      
      // 刷新订阅列表
      fetchSubscriptions();
    } catch (err) {
      message.destroy();
      message.error('更新订阅失败');
    }
  };

  const showStats = async (subscription: Subscription) => {
    setSelectedSubscription(subscription);
    try {
      const response = await subscriptionAPI.getStats(subscription.id);
      setSubscriptionStats(response.data.data!);
      setStatsDrawerVisible(true);
    } catch (err) {
      message.error('获取统计信息失败');
    }
  };

  const showLogs = async (subscription: Subscription) => {
    setSelectedSubscription(subscription);
    try {
      const response = await subscriptionAPI.getLogs(subscription.id);
      setSubscriptionLogs(response.data.data!.logs);
      setLogsDrawerVisible(true);
    } catch (err) {
      message.error('获取日志失败');
    }
  };

  const openEditModal = (subscription: Subscription) => {
    setSelectedSubscription(subscription);
    editForm.setFieldsValue({
      name: subscription.name,
      auto_update: subscription.auto_update,
      update_interval: subscription.update_interval,
    });
    setEditModalVisible(true);
  };

  const columns: ColumnsType<Subscription> = [
    {
      title: '订阅名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => <StatusBadge status={status} />,
    },
    {
      title: '节点统计',
      key: 'nodes',
      width: 150,
      render: (_, record) => (
        <div>
          <div>
            总数: <Text strong>{record.total_nodes}</Text>
          </div>
          <div>
            在线: <Text type="success">{record.active_nodes}</Text>
          </div>
          <Progress
            percent={record.total_nodes > 0 ? Math.round((record.active_nodes / record.total_nodes) * 100) : 0}
            size="small"
            showInfo={false}
            strokeColor="#52c41a"
          />
        </div>
      ),
    },
    {
      title: '自动更新',
      dataIndex: 'auto_update',
      key: 'auto_update',
      width: 100,
      render: (auto_update) => (
        <Tag color={auto_update ? 'green' : 'default'}>
          {auto_update ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '更新间隔',
      dataIndex: 'update_interval',
      key: 'update_interval',
      width: 100,
      render: (interval) => `${Math.floor(interval / 3600)}小时`,
    },
    {
      title: '最后更新',
      dataIndex: 'last_update',
      key: 'last_update',
      width: 120,
      render: (time) => (
        <Tooltip title={formatTime(time)}>
          {formatRelativeTime(time)}
        </Tooltip>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 300,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="手动更新">
            <Button
              type="text"
              size="small"
              icon={<SyncOutlined />}
              onClick={() => handleUpdateContent(record.id)}
            />
          </Tooltip>
          
          <Tooltip title={record.status === 'active' ? '禁用' : '启用'}>
            <Button
              type="text"
              size="small"
              icon={record.status === 'active' ? <PauseCircleOutlined /> : <PlayCircleOutlined />}
              onClick={() => handleToggleStatus(record)}
            />
          </Tooltip>
          
          <Tooltip title="编辑">
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() => openEditModal(record)}
            />
          </Tooltip>
          
          <Tooltip title="统计">
            <Button
              type="text"
              size="small"
              icon={<BarChartOutlined />}
              onClick={() => showStats(record)}
            />
          </Tooltip>
          
          <Tooltip title="日志">
            <Button
              type="text"
              size="small"
              icon={<HistoryOutlined />}
              onClick={() => showLogs(record)}
            />
          </Tooltip>
          
          <Popconfirm
            title="确定要删除这个订阅吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button
                type="text"
                size="small"
                icon={<DeleteOutlined />}
                danger
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <div className="flex-between" style={{ marginBottom: '24px' }}>
        <div>
          <Title level={2}>订阅管理</Title>
          <Text type="secondary">管理您的代理订阅</Text>
        </div>
        
        <Space>
          <Button icon={<ImportOutlined />}>导入</Button>
          <Button icon={<ExportOutlined />}>导出</Button>
          <Button 
            type="primary" 
            icon={<PlusOutlined />}
            onClick={() => setCreateModalVisible(true)}
          >
            添加订阅
          </Button>
        </Space>
      </div>

      {error && (
        <Alert 
          message="错误" 
          description={error} 
          type="error" 
          showIcon 
          closable 
          style={{ marginBottom: '16px' }}
          onClose={() => setError(null)}
        />
      )}

      <Card>
        <div className="flex-between" style={{ marginBottom: '16px' }}>
          <Text strong>订阅列表</Text>
          <Button 
            icon={<ReloadOutlined />} 
            onClick={fetchSubscriptions}
            loading={loading}
          >
            刷新
          </Button>
        </div>
        
        <Table
          columns={columns}
          dataSource={subscriptions}
          rowKey="id"
          loading={loading}
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
          }}
        />
      </Card>

      {/* 创建订阅弹窗 */}
      <Modal
        title="添加订阅"
        open={createModalVisible}
        onCancel={() => {
          setCreateModalVisible(false);
          form.resetFields();
        }}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreate}
        >
          <Form.Item
            label="订阅名称"
            name="name"
            rules={[{ required: true, message: '请输入订阅名称' }]}
          >
            <Input placeholder="请输入订阅名称" />
          </Form.Item>
          
          <Form.Item
            label="订阅链接"
            name="url"
            rules={[
              { required: true, message: '请输入订阅链接' },
              { type: 'url', message: '请输入有效的URL' }
            ]}
          >
            <Input placeholder="https://example.com/subscribe" />
          </Form.Item>
          
          <Form.Item
            label="User-Agent"
            name="user_agent"
            initialValue="clash"
          >
            <Input placeholder="clash" />
          </Form.Item>
          
          <Form.Item
            label="自动更新"
            name="auto_update"
            valuePropName="checked"
            initialValue={true}
          >
            <Switch />
          </Form.Item>
          
          <Form.Item
            label="更新间隔(秒)"
            name="update_interval"
            initialValue={3600}
          >
            <InputNumber min={300} max={86400} style={{ width: '100%' }} />
          </Form.Item>
          
          <Form.Item style={{ marginTop: '24px', textAlign: 'right' }}>
            <Space>
              <Button onClick={() => {
                setCreateModalVisible(false);
                form.resetFields();
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                创建
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* 编辑订阅弹窗 */}
      <Modal
        title="编辑订阅"
        open={editModalVisible}
        onCancel={() => {
          setEditModalVisible(false);
          editForm.resetFields();
          setSelectedSubscription(null);
        }}
        footer={null}
        width={600}
      >
        <Form
          form={editForm}
          layout="vertical"
          onFinish={handleEdit}
        >
          <Form.Item
            label="订阅名称"
            name="name"
            rules={[{ required: true, message: '请输入订阅名称' }]}
          >
            <Input placeholder="请输入订阅名称" />
          </Form.Item>
          
          <Form.Item
            label="自动更新"
            name="auto_update"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
          
          <Form.Item
            label="更新间隔(秒)"
            name="update_interval"
          >
            <InputNumber min={300} max={86400} style={{ width: '100%' }} />
          </Form.Item>
          
          <Form.Item style={{ marginTop: '24px', textAlign: 'right' }}>
            <Space>
              <Button onClick={() => {
                setEditModalVisible(false);
                editForm.resetFields();
                setSelectedSubscription(null);
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                保存
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* 统计信息抽屉 */}
      <Drawer
        title={`订阅统计 - ${selectedSubscription?.name}`}
        open={statsDrawerVisible}
        onClose={() => {
          setStatsDrawerVisible(false);
          setSubscriptionStats(null);
          setSelectedSubscription(null);
        }}
        width={600}
      >
        {subscriptionStats && (
          <div>
            <Descriptions column={2} style={{ marginBottom: '24px' }}>
              <Descriptions.Item label="总节点数">
                {subscriptionStats.total_nodes}
              </Descriptions.Item>
              <Descriptions.Item label="在线节点">
                {subscriptionStats.active_nodes}
              </Descriptions.Item>
              <Descriptions.Item label="存活率">
                {subscriptionStats.survival_rate.toFixed(1)}%
              </Descriptions.Item>
            </Descriptions>
            
            <Title level={4}>协议分布</Title>
            <List
              size="small"
              dataSource={subscriptionStats.protocol_distribution}
              renderItem={(item: any) => (
                <List.Item>
                  <div className="flex-between" style={{ width: '100%' }}>
                    <span>{item.protocol}</span>
                    <span>{item.count}个</span>
                  </div>
                </List.Item>
              )}
              style={{ marginBottom: '24px' }}
            />
            
            <Title level={4}>地区分布</Title>
            <List
              size="small"
              dataSource={subscriptionStats.country_distribution}
              renderItem={(item: any) => (
                <List.Item>
                  <div className="flex-between" style={{ width: '100%' }}>
                    <span>{item.country}</span>
                    <span>{item.count}个</span>
                  </div>
                </List.Item>
              )}
            />
          </div>
        )}
      </Drawer>

      {/* 日志抽屉 */}
      <Drawer
        title={`更新日志 - ${selectedSubscription?.name}`}
        open={logsDrawerVisible}
        onClose={() => {
          setLogsDrawerVisible(false);
          setSubscriptionLogs([]);
          setSelectedSubscription(null);
        }}
        width={700}
      >
        <List
          dataSource={subscriptionLogs}
          renderItem={(log) => (
            <List.Item>
              <div style={{ width: '100%' }}>
                <div className="flex-between" style={{ marginBottom: '8px' }}>
                  <div>
                    <Tag color={log.success ? 'green' : 'red'}>
                      {log.success ? '成功' : '失败'}
                    </Tag>
                    <Text type="secondary">{formatTime(log.created_at)}</Text>
                  </div>
                  <Text type="secondary">{log.update_type}</Text>
                </div>
                
                {log.success ? (
                  <div style={{ fontSize: '12px' }}>
                    <Text type="secondary">
                      获取: {log.total_fetched} | 
                      有效: {log.valid_nodes} | 
                      新增: {log.new_nodes} | 
                      全局新增: {log.global_new_nodes}
                    </Text>
                  </div>
                ) : (
                  <Text type="danger" style={{ fontSize: '12px' }}>
                    {log.error_message}
                  </Text>
                )}
              </div>
            </List.Item>
          )}
        />
      </Drawer>
    </div>
  );
};

export default Subscriptions;