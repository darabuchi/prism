import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Form, 
  Input, 
  Switch, 
  Button, 
  Select, 
  InputNumber,
  message,
  Divider,
  Space,
  Typography,
  Alert,
  Tabs,
  Table
} from 'antd';
import { SaveOutlined, ReloadOutlined, SettingOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;
const { Option } = Select;
const { TabPane } = Tabs;

interface SystemConfig {
  server: {
    host: string;
    port: number;
  };
  database: {
    driver: string;
    dsn: string;
    debug: boolean;
    max_open_conns: number;
    max_idle_conns: number;
  };
  log: {
    level: string;
    format: string;
  };
  proxy: {
    enabled: boolean;
    port: number;
    mode: string;
  };
  scheduler: {
    auto_update_interval: number;
    test_interval: number;
    cleanup_interval: number;
  };
}

const Settings: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [config, setConfig] = useState<SystemConfig | null>(null);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    // 模拟加载配置
    const mockConfig: SystemConfig = {
      server: {
        host: '0.0.0.0',
        port: 9090
      },
      database: {
        driver: 'sqlite',
        dsn: 'data/prism.db',
        debug: false,
        max_open_conns: 25,
        max_idle_conns: 5
      },
      log: {
        level: 'info',
        format: 'text'
      },
      proxy: {
        enabled: true,
        port: 7890,
        mode: 'rule'
      },
      scheduler: {
        auto_update_interval: 3600,
        test_interval: 1800,
        cleanup_interval: 21600
      }
    };
    
    setConfig(mockConfig);
    form.setFieldsValue(mockConfig);
  };

  const handleSave = async (values: any) => {
    setLoading(true);
    try {
      // 模拟保存配置
      await new Promise(resolve => setTimeout(resolve, 1000));
      setConfig(values);
      message.success('配置保存成功');
    } catch (error) {
      message.error('配置保存失败');
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    if (config) {
      form.setFieldsValue(config);
      message.info('已重置为当前配置');
    }
  };

  const systemInfoData = [
    { key: '版本', value: '1.0.0' },
    { key: '运行时间', value: '2天 3小时 45分钟' },
    { key: '内存使用', value: '128.5 MB' },
    { key: 'CPU 使用率', value: '15.2%' },
    { key: '数据库大小', value: '45.2 MB' },
    { key: '节点数量', value: '156 个' },
    { key: '订阅数量', value: '8 个' },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <div style={{ marginBottom: '24px' }}>
        <Title level={2}>设置</Title>
        <Text type="secondary">系统配置和偏好设置</Text>
      </div>

      <Tabs defaultActiveKey="general">
        <TabPane tab="常规设置" key="general">
          <Card>
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSave}
            >
              <Title level={4}>服务器配置</Title>
              <Form.Item
                label="监听地址"
                name={['server', 'host']}
                rules={[{ required: true, message: '请输入监听地址' }]}
              >
                <Input placeholder="0.0.0.0" />
              </Form.Item>
              
              <Form.Item
                label="端口"
                name={['server', 'port']}
                rules={[{ required: true, message: '请输入端口' }]}
              >
                <InputNumber min={1} max={65535} style={{ width: '100%' }} />
              </Form.Item>

              <Divider />

              <Title level={4}>代理配置</Title>
              <Form.Item
                label="启用代理"
                name={['proxy', 'enabled']}
                valuePropName="checked"
              >
                <Switch />
              </Form.Item>
              
              <Form.Item
                label="代理端口"
                name={['proxy', 'port']}
              >
                <InputNumber min={1} max={65535} style={{ width: '100%' }} />
              </Form.Item>
              
              <Form.Item
                label="代理模式"
                name={['proxy', 'mode']}
              >
                <Select>
                  <Option value="direct">直连</Option>
                  <Option value="global">全局</Option>
                  <Option value="rule">规则</Option>
                </Select>
              </Form.Item>

              <Divider />

              <Title level={4}>日志配置</Title>
              <Form.Item
                label="日志级别"
                name={['log', 'level']}
              >
                <Select>
                  <Option value="debug">调试</Option>
                  <Option value="info">信息</Option>
                  <Option value="warn">警告</Option>
                  <Option value="error">错误</Option>
                </Select>
              </Form.Item>
              
              <Form.Item
                label="日志格式"
                name={['log', 'format']}
              >
                <Select>
                  <Option value="text">文本</Option>
                  <Option value="json">JSON</Option>
                </Select>
              </Form.Item>

              <Form.Item style={{ marginTop: '24px' }}>
                <Space>
                  <Button 
                    type="primary" 
                    icon={<SaveOutlined />}
                    htmlType="submit"
                    loading={loading}
                  >
                    保存配置
                  </Button>
                  <Button 
                    icon={<ReloadOutlined />}
                    onClick={handleReset}
                  >
                    重置
                  </Button>
                </Space>
              </Form.Item>
            </Form>
          </Card>
        </TabPane>

        <TabPane tab="调度配置" key="scheduler">
          <Card>
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSave}
            >
              <Title level={4}>自动任务配置</Title>
              
              <Form.Item
                label="自动更新间隔(秒)"
                name={['scheduler', 'auto_update_interval']}
                extra="订阅自动更新的时间间隔"
              >
                <InputNumber 
                  min={300} 
                  max={86400} 
                  style={{ width: '100%' }}
                  formatter={value => `${value}s`}
                  parser={value => value!.replace('s', '')}
                />
              </Form.Item>
              
              <Form.Item
                label="节点测试间隔(秒)"
                name={['scheduler', 'test_interval']}
                extra="节点连通性测试的时间间隔"
              >
                <InputNumber 
                  min={600} 
                  max={86400} 
                  style={{ width: '100%' }}
                  formatter={value => `${value}s`}
                  parser={value => value!.replace('s', '')}
                />
              </Form.Item>
              
              <Form.Item
                label="数据清理间隔(秒)"
                name={['scheduler', 'cleanup_interval']}
                extra="过期数据清理的时间间隔"
              >
                <InputNumber 
                  min={3600} 
                  max={604800} 
                  style={{ width: '100%' }}
                  formatter={value => `${value}s`}
                  parser={value => value!.replace('s', '')}
                />
              </Form.Item>

              <Form.Item style={{ marginTop: '24px' }}>
                <Space>
                  <Button 
                    type="primary" 
                    icon={<SaveOutlined />}
                    htmlType="submit"
                    loading={loading}
                  >
                    保存配置
                  </Button>
                  <Button 
                    icon={<ReloadOutlined />}
                    onClick={handleReset}
                  >
                    重置
                  </Button>
                </Space>
              </Form.Item>
            </Form>
          </Card>
        </TabPane>

        <TabPane tab="数据库配置" key="database">
          <Card>
            <Alert
              message="数据库配置修改后需要重启应用才能生效"
              type="warning"
              showIcon
              style={{ marginBottom: '16px' }}
            />
            
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSave}
            >
              <Form.Item
                label="数据库驱动"
                name={['database', 'driver']}
              >
                <Select>
                  <Option value="sqlite">SQLite</Option>
                  <Option value="mysql">MySQL</Option>
                  <Option value="postgres">PostgreSQL</Option>
                </Select>
              </Form.Item>
              
              <Form.Item
                label="连接字符串"
                name={['database', 'dsn']}
                extra="数据库连接字符串，SQLite 为文件路径"
              >
                <Input placeholder="data/prism.db" />
              </Form.Item>
              
              <Form.Item
                label="调试模式"
                name={['database', 'debug']}
                valuePropName="checked"
                extra="启用后会打印 SQL 查询语句"
              >
                <Switch />
              </Form.Item>
              
              <Form.Item
                label="最大连接数"
                name={['database', 'max_open_conns']}
              >
                <InputNumber min={1} max={100} style={{ width: '100%' }} />
              </Form.Item>
              
              <Form.Item
                label="最大空闲连接"
                name={['database', 'max_idle_conns']}
              >
                <InputNumber min={1} max={50} style={{ width: '100%' }} />
              </Form.Item>

              <Form.Item style={{ marginTop: '24px' }}>
                <Space>
                  <Button 
                    type="primary" 
                    icon={<SaveOutlined />}
                    htmlType="submit"
                    loading={loading}
                  >
                    保存配置
                  </Button>
                  <Button 
                    icon={<ReloadOutlined />}
                    onClick={handleReset}
                  >
                    重置
                  </Button>
                </Space>
              </Form.Item>
            </Form>
          </Card>
        </TabPane>

        <TabPane tab="系统信息" key="system">
          <Card title="系统状态">
            <Table
              dataSource={systemInfoData.map((item, index) => ({ ...item, index }))}
              columns={[
                {
                  title: '项目',
                  dataIndex: 'key',
                  key: 'key',
                },
                {
                  title: '值',
                  dataIndex: 'value',
                  key: 'value',
                },
              ]}
              rowKey="index"
              pagination={false}
              size="middle"
            />
            
            <div style={{ marginTop: '24px', textAlign: 'center' }}>
              <Space>
                <Button icon={<ReloadOutlined />}>
                  刷新状态
                </Button>
                <Button type="primary" danger>
                  重启服务
                </Button>
              </Space>
            </div>
          </Card>
        </TabPane>
      </Tabs>
    </div>
  );
};

export default Settings;