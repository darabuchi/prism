import React from 'react';
import { Badge, BadgeProps } from 'antd';
import { getStatusColor, getStatusText } from '@/utils/helpers';

interface StatusBadgeProps extends Omit<BadgeProps, 'status'> {
  status: string;
  showText?: boolean;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ 
  status, 
  showText = true, 
  ...props 
}) => {
  const color = getStatusColor(status);
  const text = getStatusText(status);

  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
      <Badge 
        color={color} 
        {...props}
      />
      {showText && (
        <span style={{ color, fontSize: '12px' }}>
          {text}
        </span>
      )}
    </div>
  );
};

export default StatusBadge;