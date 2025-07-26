import React from 'react';
import { Card as AntCard, CardProps as AntCardProps } from 'antd';

export interface CardProps extends AntCardProps {
  shadow?: 'none' | 'small' | 'medium' | 'large';
  padding?: 'none' | 'small' | 'medium' | 'large';
}

export const Card: React.FC<CardProps> = ({
  shadow = 'small',
  padding = 'medium',
  className = '',
  bodyStyle,
  ...props
}) => {
  const getShadowClass = (): string => {
    switch (shadow) {
      case 'none':
        return 'shadow-none';
      case 'small':
        return 'shadow-sm';
      case 'medium':
        return 'shadow-md';
      case 'large':
        return 'shadow-lg';
      default:
        return 'shadow-sm';
    }
  };

  const getPaddingStyle = (): React.CSSProperties => {
    const paddingMap = {
      none: 0,
      small: 12,
      medium: 24,
      large: 32,
    };

    return {
      padding: paddingMap[padding],
      ...bodyStyle,
    };
  };

  const combinedClassName = `${getShadowClass()} ${className}`.trim();

  return (
    <AntCard
      className={combinedClassName}
      bodyStyle={getPaddingStyle()}
      {...props}
    />
  );
};

export default Card;