import { create } from 'zustand';

// UI state interface
export interface UIState {
  // Sidebar state
  sidebarCollapsed: boolean;
  sidebarVisible: boolean;

  // Modal state
  modals: Record<string, boolean>;

  // Loading state
  globalLoading: boolean;
  loadingStates: Record<string, boolean>;

  // Theme state
  theme: 'light' | 'dark';

  // Mobile state
  isMobile: boolean;

  // Actions
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  setSidebarVisible: (visible: boolean) => void;

  openModal: (modalId: string) => void;
  closeModal: (modalId: string) => void;
  toggleModal: (modalId: string) => void;
  isModalOpen: (modalId: string) => boolean;

  setGlobalLoading: (loading: boolean) => void;
  setLoading: (key: string, loading: boolean) => void;
  isLoading: (key: string) => boolean;

  setTheme: (theme: 'light' | 'dark') => void;
  toggleTheme: () => void;

  setIsMobile: (isMobile: boolean) => void;
}

// Create UI store
export const useUIStore = create<UIState>((set, get) => ({
  // Initial state
  sidebarCollapsed: false,
  sidebarVisible: true,
  modals: {},
  globalLoading: false,
  loadingStates: {},
  theme: 'light',
  isMobile: false,

  // Sidebar actions
  toggleSidebar: () => {
    set((state) => ({
      sidebarCollapsed: !state.sidebarCollapsed,
    }));
  },

  setSidebarCollapsed: (collapsed: boolean) => {
    set({ sidebarCollapsed: collapsed });
  },

  setSidebarVisible: (visible: boolean) => {
    set({ sidebarVisible: visible });
  },

  // Modal actions
  openModal: (modalId: string) => {
    set((state) => ({
      modals: { ...state.modals, [modalId]: true },
    }));
  },

  closeModal: (modalId: string) => {
    set((state) => ({
      modals: { ...state.modals, [modalId]: false },
    }));
  },

  toggleModal: (modalId: string) => {
    const { modals } = get();
    const isOpen = modals[modalId] ?? false;
    set((state) => ({
      modals: { ...state.modals, [modalId]: !isOpen },
    }));
  },

  isModalOpen: (modalId: string): boolean => {
    const { modals } = get();
    return modals[modalId] ?? false;
  },

  // Loading actions
  setGlobalLoading: (globalLoading: boolean) => {
    set({ globalLoading });
  },

  setLoading: (key: string, loading: boolean) => {
    set((state) => ({
      loadingStates: { ...state.loadingStates, [key]: loading },
    }));
  },

  isLoading: (key: string): boolean => {
    const { loadingStates } = get();
    return loadingStates[key] ?? false;
  },

  // Theme actions
  setTheme: (theme: 'light' | 'dark') => {
    set({ theme });
  },

  toggleTheme: () => {
    const { theme } = get();
    set({ theme: theme === 'light' ? 'dark' : 'light' });
  },

  // Mobile actions
  setIsMobile: (isMobile: boolean) => {
    set({ isMobile });
  },
}));

// UI store selectors
export const useSidebar = () => useUIStore((state) => ({
  collapsed: state.sidebarCollapsed,
  visible: state.sidebarVisible,
  toggle: state.toggleSidebar,
  setCollapsed: state.setSidebarCollapsed,
  setVisible: state.setSidebarVisible,
}));

export const useModal = (modalId: string) => useUIStore((state) => ({
  isOpen: state.isModalOpen(modalId),
  open: () => state.openModal(modalId),
  close: () => state.closeModal(modalId),
  toggle: () => state.toggleModal(modalId),
}));

export const useLoading = (key?: string) => {
  if (key) {
    return useUIStore((state) => ({
      isLoading: state.isLoading(key),
      setLoading: (loading: boolean) => state.setLoading(key, loading),
    }));
  }
  
  return useUIStore((state) => ({
    isLoading: state.globalLoading,
    setLoading: state.setGlobalLoading,
  }));
};

export const useTheme = () => useUIStore((state) => ({
  theme: state.theme,
  setTheme: state.setTheme,
  toggleTheme: state.toggleTheme,
}));

export const useIsMobile = () => useUIStore((state) => state.isMobile);