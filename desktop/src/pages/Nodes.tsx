import React, { useEffect, useState } from 'react';
import {
  Table,
  Card,
  Space,
  Button,
  Tag,
  Input,
  Select,
  Modal,
  Tooltip,
  Progress,
  message,
  Drawer,
  Descriptions,
  List,
  Typography,
  Alert,
  Popconfirm
} from 'antd';
import {
  SearchOutlined,
  ReloadOutlined,
  ThunderboltOutlined,
  EyeOutlined,
  PlayCircleOutlined,
  HistoryOutlined,
  FilterOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useNodes } from '@/store';
import { nodeAPI } from '@/utils/api';
import { formatTime, formatDelay, formatSpeed, getProtocolName } from '@/utils/helpers';
import StatusBadge from '@/components/StatusBadge';
import CountryFlag from '@/components/CountryFlag';
import { Node } from '@/types';

const { Title, Text } = Typography;
const { Option } = Select;
const { Search } = Input;

const Nodes: React.FC = () => {
  const {
    nodes,
    loading,
    error,
    setNodes,
    updateNode,
    setLoading,
    setError
  } = useNodes();

  const [selectedNodes, setSelectedNodes] = useState<number[]>([]);
  const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
  const [selectedNode, setSelectedNode] = useState<Node | null>(null);
  const [filters, setFilters] = useState({
    status: '',
    country: '',
    protocol: '',
    search: ''
  });
  const [pagination, setPagination] = useState({
    page: 1,
    size: 20,
    total: 0
  });

  useEffect(() => {
    fetchNodes();
  }, [filters, pagination.page, pagination.size]);

  const fetchNodes = async () => {
    setLoading(true);
    try {
      const params = {
        page: pagination.page,
        size: pagination.size,
        ...(filters.status && { status: filters.status }),
        ...(filters.country && { country: filters.country }),
        ...(filters.protocol && { protocol: filters.protocol }),
      };
      
      const response = await nodeAPI.getList(params);
      const data = response.data.data!;
      setNodes(data.nodes);
      setPagination(prev => ({
        ...prev,
        total: data.total
      }));
    } catch (err) {
      setError('获取节点列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleBatchTest = async () => {
    if (selectedNodes.length === 0) {
      message.warning('请先选择要测试的节点');
      return;
    }

    try {
      const response = await nodeAPI.batchTest({
        node_ids: selectedNodes,
        test_types: ['delay', 'speed'],
        test_config: {
          delay_url: 'http://www.gstatic.com/generate_204',
          timeout: 5000,
          concurrent: 10
        }
      });
      
      message.success(`测试任务已创建，任务ID: ${response.data.data!.task_id}`);
      setSelectedNodes([]);
      
      // 定期刷新节点状态
      const timer = setInterval(() => {
        fetchNodes();
      }, 5000);
      
      setTimeout(() => {
        clearInterval(timer);
      }, 60000);
      
    } catch (err) {
      message.error('创建测试任务失败');
    }
  };

  const handleSingleTest = async (nodeId: number) => {
    try {
      message.loading(`正在测试节点 ${nodeId}...`);
      await nodeAPI.test(nodeId, {
        test_types: ['delay', 'speed'],
        test_config: {
          delay_url: 'http://www.gstatic.com/generate_204',
          timeout: 5000
        }
      });
      message.destroy();
      message.success('测试完成');
      fetchNodes();
    } catch (err) {
      message.destroy();
      message.error('测试失败');
    }
  };

  const showNodeDetail = (node: Node) => {
    setSelectedNode(node);
    setDetailDrawerVisible(true);
  };

  const handleFilterChange = (key: string, value: string) => {
    setFilters(prev => ({ ...prev, [key]: value }));
    setPagination(prev => ({ ...prev, page: 1 }));
  };

  const handleTableChange = (paginationConfig: any) => {
    setPagination(prev => ({
      ...prev,
      page: paginationConfig.current,
      size: paginationConfig.pageSize
    }));
  };

  const getDelayColor = (delay?: number) => {
    if (!delay) return '#d9d9d9';
    if (delay < 100) return '#52c41a';
    if (delay < 300) return '#faad14';
    return '#f5222d';
  };

  const columns: ColumnsType<Node> = [
    {
      title: '节点名称',
      dataIndex: 'name',
      key: 'name',
      width: 250,
      ellipsis: true,
      render: (name, record) => (
        <div>
          <div style={{ marginBottom: '4px' }}>
            <Text strong>{name}</Text>
          </div>
          <div>
            <Text type="secondary" style={{ fontSize: '12px' }}>
              {record.server}:{record.port}
            </Text>
          </div>
        </div>
      ),
    },
    {
      title: '地区',
      dataIndex: 'country',
      key: 'country',
      width: 100,
      render: (country) => <CountryFlag country={country} size="small" />,
    },
    {
      title: '协议',
      dataIndex: 'protocol',
      key: 'protocol',
      width: 100,
      render: (protocol) => (
        <Tag color="blue">{getProtocolName(protocol)}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => <StatusBadge status={status} />,
    },
    {
      title: '延迟',
      dataIndex: 'delay',
      key: 'delay',
      width: 80,
      render: (delay) => (
        <Text style={{ color: getDelayColor(delay) }}>
          {formatDelay(delay)}
        </Text>
      ),
      sorter: true,
    },
    {
      title: '速度',
      key: 'speed',
      width: 120,
      render: (_, record) => (
        <div style={{ fontSize: '12px' }}>
          <div>↑ {formatSpeed(record.upload_speed)}</div>
          <div>↓ {formatSpeed(record.download_speed)}</div>
        </div>
      ),
    },
    {
      title: '最后测试',
      dataIndex: 'last_test',
      key: 'last_test',
      width: 120,
      render: (time) => (
        <Text type="secondary" style={{ fontSize: '12px' }}>
          {formatTime(time, 'MM-DD HH:mm')}
        </Text>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="测试">
            <Button
              type="text"
              size="small"
              icon={<PlayCircleOutlined />}
              onClick={() => handleSingleTest(record.id)}
            />
          </Tooltip>
          
          <Tooltip title="详情">
            <Button
              type="text"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => showNodeDetail(record)}
            />
          </Tooltip>
          
          <Tooltip title="测试历史">
            <Button
              type="text"
              size="small"
              icon={<HistoryOutlined />}
              onClick={() => {/* TODO: 显示测试历史 */}}
            />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const rowSelection = {
    selectedRowKeys: selectedNodes,
    onChange: (selectedRowKeys: React.Key[]) => {
      setSelectedNodes(selectedRowKeys as number[]);
    },
    onSelectAll: (selected: boolean, selectedRows: Node[], changeRows: Node[]) => {
      if (selected) {
        setSelectedNodes(prev => [...prev, ...changeRows.map(row => row.id)]);
      } else {
        const removeIds = changeRows.map(row => row.id);
        setSelectedNodes(prev => prev.filter(id => !removeIds.includes(id)));
      }
    },
  };

  return (
    <div style={{ padding: '24px' }}>
      <div className="flex-between" style={{ marginBottom: '24px' }}>
        <div>
          <Title level={2}>节点管理</Title>
          <Text type="secondary">管理和测试代理节点</Text>
        </div>
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
        {/* 过滤器 */}
        <div style={{ marginBottom: '16px' }}>
          <Space wrap>
            <Search
              placeholder="搜索节点名称"
              style={{ width: 200 }}
              onSearch={(value) => handleFilterChange('search', value)}
              allowClear
            />
            
            <Select
              placeholder="状态"
              style={{ width: 120 }}
              allowClear
              onChange={(value) => handleFilterChange('status', value || '')}
            >
              <Option value="online">在线</Option>
              <Option value="offline">离线</Option>
              <Option value="testing">测试中</Option>
              <Option value="unknown">未知</Option>
            </Select>
            
            <Select
              placeholder="地区"
              style={{ width: 120 }}
              allowClear
              onChange={(value) => handleFilterChange('country', value || '')}
            >
              <Option value="HK">香港</Option>
              <Option value="SG">新加坡</Option>
              <Option value="JP">日本</Option>
              <Option value="US">美国</Option>
              <Option value="UK">英国</Option>
            </Select>
            
            <Select
              placeholder="协议"
              style={{ width: 120 }}
              allowClear
              onChange={(value) => handleFilterChange('protocol', value || '')}
            >
              <Option value="vmess">VMess</Option>
              <Option value="ss">Shadowsocks</Option>
              <Option value="trojan">Trojan</Option>
              <Option value="hysteria">Hysteria</Option>
            </Select>
          </Space>
        </div>

        {/* 批量操作 */}
        <div className="flex-between" style={{ marginBottom: '16px' }}>
          <Space>
            {selectedNodes.length > 0 && (
              <>
                <Text>已选择 {selectedNodes.length} 个节点</Text>
                <Button
                  type="primary"
                  icon={<ThunderboltOutlined />}
                  onClick={handleBatchTest}
                >
                  批量测试
                </Button>
                <Button onClick={() => setSelectedNodes([])}>
                  取消选择
                </Button>
              </>
            )}
          </Space>
          
          <Button
            icon={<ReloadOutlined />}
            onClick={fetchNodes}
            loading={loading}
          >
            刷新
          </Button>
        </div>

        <Table
          rowSelection={rowSelection}
          columns={columns}
          dataSource={nodes}
          rowKey="id"
          loading={loading}
          pagination={{
            current: pagination.page,
            pageSize: pagination.size,
            total: pagination.total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
          }}
          onChange={handleTableChange}
          scroll={{ x: 1200 }}
        />
      </Card>

      {/* 节点详情抽屉 */}
      <Drawer
        title="节点详情"
        open={detailDrawerVisible}
        onClose={() => {
          setDetailDrawerVisible(false);
          setSelectedNode(null);
        }}
        width={600}
      >
        {selectedNode && (
          <div>
            <Descriptions column={1} bordered>
              <Descriptions.Item label="节点名称">
                {selectedNode.name}
              </Descriptions.Item>
              <Descriptions.Item label="服务器地址">
                {selectedNode.server}:{selectedNode.port}
              </Descriptions.Item>
              <Descriptions.Item label="协议">
                <Tag color="blue">{getProtocolName(selectedNode.protocol)}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                <StatusBadge status={selectedNode.status} />
              </Descriptions.Item>
              <Descriptions.Item label="地区">
                <CountryFlag country={selectedNode.country} />
              </Descriptions.Item>
              <Descriptions.Item label="延迟">
                <Text style={{ color: getDelayColor(selectedNode.delay) }}>
                  {formatDelay(selectedNode.delay)}
                </Text>
              </Descriptions.Item>
              <Descriptions.Item label="上传速度">
                {formatSpeed(selectedNode.upload_speed)}
              </Descriptions.Item>
              <Descriptions.Item label="下载速度">
                {formatSpeed(selectedNode.download_speed)}
              </Descriptions.Item>
              <Descriptions.Item label="丢包率">
                {selectedNode.loss_rate ? `${selectedNode.loss_rate}%` : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="连续失败">
                {selectedNode.continuous_failures} 次
              </Descriptions.Item>
              <Descriptions.Item label="最后测试">
                {formatTime(selectedNode.last_test)}
              </Descriptions.Item>
              <Descriptions.Item label="最后在线">
                {formatTime(selectedNode.last_online)}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {formatTime(selectedNode.created_at)}
              </Descriptions.Item>
            </Descriptions>

            {selectedNode.streaming_unlock && (
              <div style={{ marginTop: '24px' }}>
                <Title level={4}>流媒体解锁</Title>
                <List
                  size="small"
                  bordered
                  dataSource={Object.entries(selectedNode.streaming_unlock)}
                  renderItem={([service, info]: [string, any]) => (
                    <List.Item>
                      <div className="flex-between" style={{ width: '100%' }}>
                        <span style={{ textTransform: 'capitalize' }}>
                          {service.replace('_', ' ')}
                        </span>
                        <div>
                          <Tag color={info.available ? 'green' : 'red'}>
                            {info.available ? '支持' : '不支持'}
                          </Tag>
                          {info.region && (
                            <Text type="secondary">({info.region})</Text>
                          )}
                        </div>
                      </div>
                    </List.Item>
                  )}
                />
              </div>
            )}
          </div>
        )}
      </Drawer>
    </div>
  );
};

export default Nodes;