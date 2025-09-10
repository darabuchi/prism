import dayjs from 'dayjs';

// 格式化时间
export const formatTime = (time: string | undefined, format = 'YYYY-MM-DD HH:mm:ss') => {
  if (!time) return '-';
  return dayjs(time).format(format);
};

// 格式化相对时间
export const formatRelativeTime = (time: string | undefined) => {
  if (!time) return '-';
  return dayjs(time).fromNow();
};

// 格式化文件大小
export const formatFileSize = (bytes: number | undefined) => {
  if (!bytes || bytes === 0) return '-';
  
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
};

// 格式化网络速度
export const formatSpeed = (bps: number | undefined) => {
  if (!bps || bps === 0) return '-';
  
  const speeds = ['bps', 'Kbps', 'Mbps', 'Gbps'];
  const i = Math.floor(Math.log(bps) / Math.log(1024));
  return `${(bps / Math.pow(1024, i)).toFixed(1)} ${speeds[i]}`;
};

// 格式化延迟
export const formatDelay = (ms: number | undefined) => {
  if (!ms) return '-';
  return `${ms}ms`;
};

// 格式化百分比
export const formatPercentage = (value: number | undefined, decimals = 1) => {
  if (value === undefined || value === null) return '-';
  return `${value.toFixed(decimals)}%`;
};

// 获取状态颜色
export const getStatusColor = (status: string) => {
  switch (status) {
    case 'online':
    case 'active':
      return '#52c41a';
    case 'offline':
    case 'inactive':
    case 'error':
      return '#ff4d4f';
    case 'testing':
      return '#1677ff';
    default:
      return '#d9d9d9';
  }
};

// 获取状态文本
export const getStatusText = (status: string) => {
  switch (status) {
    case 'online':
      return '在线';
    case 'offline':
      return '离线';
    case 'testing':
      return '测试中';
    case 'unknown':
      return '未知';
    case 'active':
      return '启用';
    case 'inactive':
      return '禁用';
    case 'error':
      return '错误';
    default:
      return status;
  }
};

// 获取协议显示名称
export const getProtocolName = (protocol: string) => {
  const protocolMap: Record<string, string> = {
    'ss': 'Shadowsocks',
    'ssr': 'ShadowsocksR',
    'vmess': 'VMess',
    'vless': 'VLESS',
    'trojan': 'Trojan',
    'hysteria': 'Hysteria',
    'hysteria2': 'Hysteria2',
    'tuic': 'TUIC',
  };
  return protocolMap[protocol] || protocol.toUpperCase();
};

// 获取国家显示名称
export const getCountryName = (country: string) => {
  const countryMap: Record<string, string> = {
    'HK': '香港',
    'TW': '台湾',
    'SG': '新加坡',
    'JP': '日本',
    'KR': '韩国',
    'US': '美国',
    'UK': '英国',
    'DE': '德国',
    'FR': '法国',
    'CA': '加拿大',
    'AU': '澳大利亚',
    'RU': '俄罗斯',
    'IN': '印度',
    'BR': '巴西',
    'AR': '阿根廷',
    'MX': '墨西哥',
    'TH': '泰国',
    'MY': '马来西亚',
    'ID': '印度尼西亚',
    'VN': '越南',
    'PH': '菲律宾',
  };
  return countryMap[country] || country;
};

// 获取国家旗帜 emoji
export const getCountryFlag = (country: string) => {
  const flagMap: Record<string, string> = {
    'HK': '🇭🇰',
    'TW': '🇹🇼',
    'SG': '🇸🇬',
    'JP': '🇯🇵',
    'KR': '🇰🇷',
    'US': '🇺🇸',
    'UK': '🇬🇧',
    'DE': '🇩🇪',
    'FR': '🇫🇷',
    'CA': '🇨🇦',
    'AU': '🇦🇺',
    'RU': '🇷🇺',
    'IN': '🇮🇳',
    'BR': '🇧🇷',
    'AR': '🇦🇷',
    'MX': '🇲🇽',
    'TH': '🇹🇭',
    'MY': '🇲🇾',
    'ID': '🇮🇩',
    'VN': '🇻🇳',
    'PH': '🇵🇭',
  };
  return flagMap[country] || '🌐';
};

// 复制到剪贴板
export const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch (error) {
    console.error('复制失败:', error);
    return false;
  }
};

// 下载文件
export const downloadFile = (data: any, filename: string, type = 'application/json') => {
  const blob = new Blob([JSON.stringify(data, null, 2)], { type });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
};

// 验证 URL
export const isValidUrl = (url: string) => {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
};

// 延迟函数
export const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// 防抖函数
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
};

// 节流函数
export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean;
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
};