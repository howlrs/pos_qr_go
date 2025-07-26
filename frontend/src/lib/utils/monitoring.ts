// Application monitoring and error tracking utilities

export interface ErrorInfo {
  message: string;
  stack?: string;
  componentStack?: string;
  errorBoundary?: string;
  userId?: string;
  sessionId?: string;
  url?: string;
  userAgent?: string;
  timestamp: string;
}

export interface PerformanceInfo {
  metric: string;
  value: number;
  url?: string;
  userId?: string;
  sessionId?: string;
  timestamp: string;
}

export class ApplicationMonitor {
  private static instance: ApplicationMonitor;
  private isEnabled: boolean;

  private constructor() {
    this.isEnabled = process.env.NODE_ENV === 'production';
  }

  public static getInstance(): ApplicationMonitor {
    if (!ApplicationMonitor.instance) {
      ApplicationMonitor.instance = new ApplicationMonitor();
    }
    return ApplicationMonitor.instance;
  }

  // Error tracking
  public trackError(error: Error, errorInfo?: any): void {
    if (!this.isEnabled) {
      console.error('Error tracked:', error, errorInfo);
      return;
    }

    const errorData: ErrorInfo = {
      message: error.message,
      stack: error.stack,
      componentStack: errorInfo?.componentStack,
      errorBoundary: errorInfo?.errorBoundary,
      url: typeof window !== 'undefined' ? window.location.href : undefined,
      userAgent: typeof window !== 'undefined' ? window.navigator.userAgent : undefined,
      timestamp: new Date().toISOString(),
    };

    // Send to error tracking service (e.g., Sentry, LogRocket)
    this.sendToErrorService(errorData);
  }

  // Performance tracking
  public trackPerformance(metric: string, value: number): void {
    if (!this.isEnabled) {
      console.log(`Performance metric - ${metric}: ${value}ms`);
      return;
    }

    const performanceData: PerformanceInfo = {
      metric,
      value,
      url: typeof window !== 'undefined' ? window.location.href : undefined,
      timestamp: new Date().toISOString(),
    };

    // Send to analytics service
    this.sendToAnalyticsService(performanceData);
  }

  // User interaction tracking
  public trackUserAction(action: string, properties?: Record<string, any>): void {
    if (!this.isEnabled) {
      console.log(`User action - ${action}:`, properties);
      return;
    }

    const eventData = {
      action,
      properties: {
        ...properties,
        url: typeof window !== 'undefined' ? window.location.href : undefined,
        timestamp: new Date().toISOString(),
      },
    };

    // Send to analytics service
    this.sendToAnalyticsService(eventData);
  }

  // Business metrics tracking
  public trackBusinessMetric(metric: string, value: number, properties?: Record<string, any>): void {
    if (!this.isEnabled) {
      console.log(`Business metric - ${metric}: ${value}`, properties);
      return;
    }

    const metricData = {
      metric,
      value,
      properties: {
        ...properties,
        timestamp: new Date().toISOString(),
      },
    };

    // Send to business intelligence service
    this.sendToAnalyticsService(metricData);
  }

  // Health check
  public async performHealthCheck(): Promise<boolean> {
    try {
      // Check API connectivity
      const apiResponse = await fetch('/api/health', {
        method: 'GET',
        timeout: 5000,
      } as any);

      if (!apiResponse.ok) {
        throw new Error(`API health check failed: ${apiResponse.status}`);
      }

      // Check local storage
      if (typeof window !== 'undefined') {
        try {
          localStorage.setItem('health-check', 'test');
          localStorage.removeItem('health-check');
        } catch (e) {
          throw new Error('Local storage not available');
        }
      }

      return true;
    } catch (error) {
      this.trackError(error as Error, { context: 'health-check' });
      return false;
    }
  }

  private sendToErrorService(errorData: ErrorInfo): void {
    // Example: Send to Sentry
    // Sentry.captureException(new Error(errorData.message), {
    //   extra: errorData,
    // });

    // Example: Send to custom error service
    fetch('/api/errors', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(errorData),
    }).catch(console.error);
  }

  private sendToAnalyticsService(data: any): void {
    // Example: Send to Google Analytics
    // gtag('event', data.action || data.metric, data.properties);

    // Example: Send to custom analytics service
    fetch('/api/analytics', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    }).catch(console.error);
  }
}

// Singleton instance
export const monitor = ApplicationMonitor.getInstance();

// React hook for monitoring
export const useMonitoring = () => {
  return {
    trackError: monitor.trackError.bind(monitor),
    trackPerformance: monitor.trackPerformance.bind(monitor),
    trackUserAction: monitor.trackUserAction.bind(monitor),
    trackBusinessMetric: monitor.trackBusinessMetric.bind(monitor),
    performHealthCheck: monitor.performHealthCheck.bind(monitor),
  };
};

// HOC for error boundary monitoring
export function withErrorMonitoring<P extends object>(
  WrappedComponent: React.ComponentType<P>,
  componentName: string
) {
  return function MonitoredComponent(props: P) {
    React.useEffect(() => {
      // Track component mount
      monitor.trackUserAction('component_mount', { component: componentName });

      return () => {
        // Track component unmount
        monitor.trackUserAction('component_unmount', { component: componentName });
      };
    }, []);

    return React.createElement(WrappedComponent, props);
  };
}

// Utility functions for common tracking scenarios
export const trackOrderPlaced = (orderId: string, amount: number, items: number) => {
  monitor.trackBusinessMetric('order_placed', amount, {
    orderId,
    itemCount: items,
  });
  monitor.trackUserAction('place_order', {
    orderId,
    amount,
    itemCount: items,
  });
};

export const trackCartAction = (action: 'add' | 'remove' | 'update', itemId: string, quantity: number) => {
  monitor.trackUserAction(`cart_${action}`, {
    itemId,
    quantity,
  });
};

export const trackPageView = (page: string, properties?: Record<string, any>) => {
  monitor.trackUserAction('page_view', {
    page,
    ...properties,
  });
};

export const trackSearchQuery = (query: string, resultsCount: number) => {
  monitor.trackUserAction('search', {
    query,
    resultsCount,
  });
};