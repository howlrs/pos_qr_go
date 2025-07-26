import { theme } from 'antd';

// Theme configuration
export const lightTheme = {
  algorithm: theme.defaultAlgorithm,
  token: {
    colorPrimary: '#1890ff',
    colorSuccess: '#52c41a',
    colorWarning: '#faad14',
    colorError: '#ff4d4f',
    colorInfo: '#1890ff',
    colorBgBase: '#ffffff',
    colorTextBase: '#000000',
    borderRadius: 6,
    wireframe: false,
  },
  components: {
    Layout: {
      headerBg: '#ffffff',
      siderBg: '#ffffff',
      bodyBg: '#f5f5f5',
    },
    Menu: {
      itemBg: 'transparent',
      itemSelectedBg: '#e6f7ff',
      itemHoverBg: '#f5f5f5',
    },
    Button: {
      borderRadius: 6,
    },
    Card: {
      borderRadius: 8,
    },
  },
};

export const darkTheme = {
  algorithm: theme.darkAlgorithm,
  token: {
    colorPrimary: '#1890ff',
    colorSuccess: '#52c41a',
    colorWarning: '#faad14',
    colorError: '#ff4d4f',
    colorInfo: '#1890ff',
    colorBgBase: '#141414',
    colorTextBase: '#ffffff',
    borderRadius: 6,
    wireframe: false,
  },
  components: {
    Layout: {
      headerBg: '#001529',
      siderBg: '#001529',
      bodyBg: '#000000',
    },
    Menu: {
      itemBg: 'transparent',
      itemSelectedBg: '#1890ff',
      itemHoverBg: '#262626',
    },
    Button: {
      borderRadius: 6,
    },
    Card: {
      borderRadius: 8,
    },
  },
};

// Responsive breakpoints
export const breakpoints = {
  xs: 480,
  sm: 576,
  md: 768,
  lg: 992,
  xl: 1200,
  xxl: 1600,
} as const;

// Media query helpers
export const mediaQueries = {
  xs: `@media (max-width: ${breakpoints.xs}px)`,
  sm: `@media (max-width: ${breakpoints.sm}px)`,
  md: `@media (max-width: ${breakpoints.md}px)`,
  lg: `@media (max-width: ${breakpoints.lg}px)`,
  xl: `@media (max-width: ${breakpoints.xl}px)`,
  xxl: `@media (max-width: ${breakpoints.xxl}px)`,
  
  minXs: `@media (min-width: ${breakpoints.xs + 1}px)`,
  minSm: `@media (min-width: ${breakpoints.sm + 1}px)`,
  minMd: `@media (min-width: ${breakpoints.md + 1}px)`,
  minLg: `@media (min-width: ${breakpoints.lg + 1}px)`,
  minXl: `@media (min-width: ${breakpoints.xl + 1}px)`,
  minXxl: `@media (min-width: ${breakpoints.xxl + 1}px)`,
} as const;

// Theme utilities
export const getThemeConfig = (isDark: boolean) => {
  return isDark ? darkTheme : lightTheme;
};

// Responsive utilities
export const isMobile = (width: number): boolean => width <= breakpoints.sm;
export const isTablet = (width: number): boolean => 
  width > breakpoints.sm && width <= breakpoints.lg;
export const isDesktop = (width: number): boolean => width > breakpoints.lg;

export default {
  lightTheme,
  darkTheme,
  breakpoints,
  mediaQueries,
  getThemeConfig,
  isMobile,
  isTablet,
  isDesktop,
};