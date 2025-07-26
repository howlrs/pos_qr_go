export interface StoreInfoProps {
  visible: boolean;
  store?: {
    id: string;
    name: string;
    address: string;
    phone?: string;
  };
  onClose: () => void;
}