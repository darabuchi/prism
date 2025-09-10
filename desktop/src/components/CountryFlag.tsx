import React from 'react';
import { Tooltip } from 'antd';
import { getCountryFlag, getCountryName } from '@/utils/helpers';

interface CountryFlagProps {
  country?: string;
  showName?: boolean;
  size?: 'small' | 'medium' | 'large';
}

const CountryFlag: React.FC<CountryFlagProps> = ({ 
  country, 
  showName = true, 
  size = 'medium' 
}) => {
  if (!country) return null;

  const flag = getCountryFlag(country);
  const name = getCountryName(country);

  const sizeMap = {
    small: '14px',
    medium: '16px',
    large: '20px',
  };

  const flagElement = (
    <span 
      style={{ 
        fontSize: sizeMap[size], 
        lineHeight: 1,
        marginRight: showName ? '4px' : 0
      }}
    >
      {flag}
    </span>
  );

  return (
    <div style={{ display: 'flex', alignItems: 'center' }}>
      {showName ? (
        <>
          {flagElement}
          <Tooltip title={`${name} (${country})`}>
            <span style={{ fontSize: '12px' }}>
              {name}
            </span>
          </Tooltip>
        </>
      ) : (
        <Tooltip title={`${name} (${country})`}>
          {flagElement}
        </Tooltip>
      )}
    </div>
  );
};

export default CountryFlag;