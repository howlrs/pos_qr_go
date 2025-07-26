import React from 'react';
import { Spin, SpinProps } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export interface LoadingProps extends SpinProps {
  fullScreen?: boolean;
  message?: string;
  overlay?: boolean;
}

export const Loading: React.FC<LoadingProps> = ({
  fullScreen = false,
  message,
  overlay = false,
  className = '',
  ...props
}) => {
  const loadingIcon = <LoadingOutlined style={{ fontSize: 24 }} spin />;

  const spinElement = (
    <Spin
      indicator={loadingIcon}
      tip={message}
      className={className}
      {...props}
    />
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-white bg-opacity-75 z-50">
        {spinElement}
      </div>
    );
  }

  if (overlay) {
    return (
      <div className="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10">
        {spinElement}
      </div>
    );
  }

  return spinElement;
};

// Page loading component
export const PageLoading: React.FC<{ message?: string }> = ({ 
  message = '読み込み中...' 
}) => (
  <Loading fullScreen message={message} size="large" />
);

// Inline loading component
export const InlineLoading: React.FC<{ message?: string }> = ({ 
  message = '読み込み中...' 
}) => (
  <div className="flex items-center justify-center p-4">
    <Loading message={message} />
  </div>
);

// Button loading component
export const ButtonLoading: React.FC = () => (
  <LoadingOutlined style={{ fontSize: 16 }} spin />
);

export default Loading;