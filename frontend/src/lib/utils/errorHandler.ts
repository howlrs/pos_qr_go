import { message } from 'antd';

import { ApiError } from '@/lib/api';
import { isDevelopment } from '@/lib/config/env';

// Error types
export interface AppError extends Error {
  code?: string;
  statusCode?: number;
  details?: Record<string, unknown>;
}

// Error handler class
export class ErrorHandler {
  /**
   * Handle API errors
   */
  static handleApiError(error: ApiError | Error): void {
    if (isDevelopment) {
      // eslint-disable-next-line no-console
      console.error('API Error:', error);
    }

    // Check if it's an ApiError
    if ('code' in error && 'message' in error) {
      const apiError = error as ApiError;
      
      switch (apiError.code) {
        case 'NETWORK_ERROR':
          message.error('ネットワークエラーが発生しました。接続を確認してください。');
          break;
        case 'HTTP_401':
          message.error('認証が必要です。再度ログインしてください。');
          break;
        case 'HTTP_403':
          message.error('この操作を実行する権限がありません。');
          break;
        case 'HTTP_404':
          message.error('要求されたリソースが見つかりません。');
          break;
        case 'HTTP_500':
          message.error('サーバーエラーが発生しました。しばらく後に再試行してください。');
          break;
        default:
          message.error(apiError.message || '予期しないエラーが発生しました。');
      }
    } else {
      // Generic error
      message.error(error.message || '予期しないエラーが発生しました。');
    }
  }

  /**
   * Handle form validation errors
   */
  static handleValidationError(errors: Record<string, string[]>): void {
    const errorMessages = Object.entries(errors)
      .map(([field, messages]) => `${field}: ${messages.join(', ')}`)
      .join('\n');
    
    message.error(`入力エラー:\n${errorMessages}`);
  }

  /**
   * Handle generic application errors
   */
  static handleAppError(error: AppError): void {
    if (isDevelopment) {
      // eslint-disable-next-line no-console
      console.error('App Error:', error);
    }

    const errorMessage = error.message || '予期しないエラーが発生しました。';
    message.error(errorMessage);
  }

  /**
   * Handle async operation errors
   */
  static async handleAsyncError<T>(
    operation: () => Promise<T>,
    errorMessage?: string
  ): Promise<T | null> {
    try {
      return await operation();
    } catch (error) {
      if (errorMessage) {
        message.error(errorMessage);
      } else {
        this.handleApiError(error as Error);
      }
      return null;
    }
  }

  /**
   * Create error boundary error handler
   */
  static createErrorBoundaryHandler(componentName: string) {
    return (error: Error, errorInfo: { componentStack: string }) => {
      if (isDevelopment) {
        // eslint-disable-next-line no-console
        console.error(`Error in ${componentName}:`, error, errorInfo);
      }

      // Log error to monitoring service in production
      if (!isDevelopment) {
        // TODO: Integrate with error monitoring service (e.g., Sentry)
        // logErrorToService(error, errorInfo, componentName);
      }

      message.error(`${componentName}でエラーが発生しました。ページを再読み込みしてください。`);
    };
  }
}

// Utility functions
export const handleApiError = ErrorHandler.handleApiError;
export const handleValidationError = ErrorHandler.handleValidationError;
export const handleAppError = ErrorHandler.handleAppError;
export const handleAsyncError = ErrorHandler.handleAsyncError;

export default ErrorHandler;