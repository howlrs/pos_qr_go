'use client';

import React from 'react';
import { Button, Result } from 'antd';
import { ReloadOutlined, HomeOutlined } from '@ant-design/icons';

import { ErrorHandler } from '@/lib/utils/errorHandler';
import { isDevelopment } from '@/lib/config/env';

interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
  errorInfo?: React.ErrorInfo;
}

interface ErrorBoundaryProps {
  children: React.ReactNode;
  fallback?: React.ComponentType<{ error?: Error; retry: () => void }>;
  onError?: (error: Error, errorInfo: React.ErrorInfo) => void;
}

export class ErrorBoundary extends React.Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    this.setState({
      error,
      errorInfo,
    });

    // Handle error
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    } else {
      ErrorHandler.createErrorBoundaryHandler('ErrorBoundary')(error, {
        componentStack: errorInfo.componentStack || '',
      });
    }
  }

  handleRetry = () => {
    this.setState({ hasError: false, error: undefined, errorInfo: undefined });
  };

  handleGoHome = () => {
    window.location.href = '/';
  };

  render() {
    if (this.state.hasError) {
      // Custom fallback component
      if (this.props.fallback) {
        const FallbackComponent = this.props.fallback;
        return <FallbackComponent error={this.state.error} retry={this.handleRetry} />;
      }

      // Default error UI
      return (
        <div className="min-h-screen flex items-center justify-center p-4">
          <Result
            status="error"
            title="予期しないエラーが発生しました"
            subTitle={
              isDevelopment && this.state.error
                ? this.state.error.message
                : 'アプリケーションでエラーが発生しました。ページを再読み込みするか、ホームに戻ってください。'
            }
            extra={[
              <Button
                key="retry"
                type="primary"
                icon={<ReloadOutlined />}
                onClick={this.handleRetry}
              >
                再試行
              </Button>,
              <Button
                key="home"
                icon={<HomeOutlined />}
                onClick={this.handleGoHome}
              >
                ホームに戻る
              </Button>,
            ]}
          />
          {isDevelopment && this.state.errorInfo && (
            <details className="mt-4 p-4 bg-gray-100 rounded">
              <summary className="cursor-pointer font-semibold">
                開発者情報 (本番環境では表示されません)
              </summary>
              <pre className="mt-2 text-sm overflow-auto">
                {this.state.error?.stack}
                {this.state.errorInfo.componentStack}
              </pre>
            </details>
          )}
        </div>
      );
    }

    return this.props.children;
  }
}

// Hook for error boundary
export const useErrorHandler = () => {
  const [error, setError] = React.useState<Error | null>(null);

  const handleError = React.useCallback((error: Error) => {
    setError(error);
    ErrorHandler.handleAppError(error);
  }, []);

  const clearError = React.useCallback(() => {
    setError(null);
  }, []);

  // Throw error to be caught by error boundary
  if (error) {
    throw error;
  }

  return { handleError, clearError };
};

export default ErrorBoundary;