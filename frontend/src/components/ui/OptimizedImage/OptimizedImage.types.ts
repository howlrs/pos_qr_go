import { ImageProps } from 'next/image';

export interface OptimizedImageProps extends Omit<ImageProps, 'src'> {
  src: string;
  alt: string;
  lazy?: boolean;
  fallbackSrc?: string;
}