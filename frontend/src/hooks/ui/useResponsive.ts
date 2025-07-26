import { useState, useEffect } from 'react';

import { breakpoints, isMobile, isTablet, isDesktop } from '@/lib/utils/theme';
import { useUIStore } from '@/store';

export interface ResponsiveInfo {
  width: number;
  height: number;
  isMobile: boolean;
  isTablet: boolean;
  isDesktop: boolean;
  breakpoint: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | 'xxl';
}

export const useResponsive = (): ResponsiveInfo => {
  const setIsMobile = useUIStore((state) => state.setIsMobile);
  
  const [responsive, setResponsive] = useState<ResponsiveInfo>({
    width: 0,
    height: 0,
    isMobile: false,
    isTablet: false,
    isDesktop: true,
    breakpoint: 'lg',
  });

  useEffect(() => {
    const updateResponsive = () => {
      const width = window.innerWidth;
      const height = window.innerHeight;
      
      const mobile = isMobile(width);
      const tablet = isTablet(width);
      const desktop = isDesktop(width);

      // Determine breakpoint
      let breakpoint: ResponsiveInfo['breakpoint'] = 'lg';
      if (width <= breakpoints.xs) breakpoint = 'xs';
      else if (width <= breakpoints.sm) breakpoint = 'sm';
      else if (width <= breakpoints.md) breakpoint = 'md';
      else if (width <= breakpoints.lg) breakpoint = 'lg';
      else if (width <= breakpoints.xl) breakpoint = 'xl';
      else breakpoint = 'xxl';

      const newResponsive: ResponsiveInfo = {
        width,
        height,
        isMobile: mobile,
        isTablet: tablet,
        isDesktop: desktop,
        breakpoint,
      };

      setResponsive(newResponsive);
      
      // Update global mobile state
      setIsMobile(mobile);
    };

    // Initial update
    updateResponsive();

    // Add event listener
    window.addEventListener('resize', updateResponsive);

    // Cleanup
    return () => {
      window.removeEventListener('resize', updateResponsive);
    };
  }, [setIsMobile]);

  return responsive;
};

// Hook for specific breakpoint checks
export const useBreakpoint = () => {
  const responsive = useResponsive();
  
  return {
    ...responsive,
    isXs: responsive.breakpoint === 'xs',
    isSm: responsive.breakpoint === 'sm',
    isMd: responsive.breakpoint === 'md',
    isLg: responsive.breakpoint === 'lg',
    isXl: responsive.breakpoint === 'xl',
    isXxl: responsive.breakpoint === 'xxl',
    
    // Utility functions
    isAtLeast: (breakpoint: ResponsiveInfo['breakpoint']) => {
      const order = ['xs', 'sm', 'md', 'lg', 'xl', 'xxl'];
      const currentIndex = order.indexOf(responsive.breakpoint);
      const targetIndex = order.indexOf(breakpoint);
      return currentIndex >= targetIndex;
    },
    
    isAtMost: (breakpoint: ResponsiveInfo['breakpoint']) => {
      const order = ['xs', 'sm', 'md', 'lg', 'xl', 'xxl'];
      const currentIndex = order.indexOf(responsive.breakpoint);
      const targetIndex = order.indexOf(breakpoint);
      return currentIndex <= targetIndex;
    },
  };
};

export default useResponsive;