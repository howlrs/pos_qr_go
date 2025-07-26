import React from 'react';
import { Button as AntButton, ButtonProps as AntButtonProps } from 'antd';

export interface ButtonProps extends Omit<AntButtonProps, 'variant'> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  fullWidth?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  fullWidth = false,
  className = '',
  ...props
}) => {
  const getButtonType = (): AntButtonProps['type'] => {
    switch (variant) {
      case 'primary':
        return 'primary';
      case 'secondary':
        return 'default';
      case 'danger':
        return 'primary';
      case 'ghost':
        return 'text';
      default:
        return 'primary';
    }
  };

  const getDanger = (): boolean => {
    return variant === 'danger';
  };

  const combinedClassName = `${fullWidth ? 'w-full' : ''} ${className}`.trim();

  return (
    <AntButton
      type={getButtonType()}
      danger={getDanger()}
      className={combinedClassName}
      block={fullWidth}
      {...props}
    />
  );
};

export default Button;