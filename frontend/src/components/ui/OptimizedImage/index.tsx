'use client';

import React, { useState, useRef, useEffect } from 'react';
import Image from 'next/image';
import { Skeleton } from 'antd';
import { createIntersectionObserver } from '@/lib/utils/performance';
import { OptimizedImageProps } from './OptimizedImage.types';

export const OptimizedImage: React.FC<OptimizedImageProps> = ({
  src,
  alt,
  width,
  height,
  className,
  priority = false,
  quality = 75,
  placeholder = 'blur',
  blurDataURL,
  sizes,
  fill = false,
  lazy = true,
  onLoad,
  onError,
  fallbackSrc = '/images/placeholder.jpg',
  ...props
}) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const [hasError, setHasError] = useState(false);
  const [isInView, setIsInView] = useState(!lazy || priority);
  const imgRef = useRef<HTMLDivElement>(null);

  // Intersection Observer for lazy loading
  useEffect(() => {
    if (!lazy || priority || isInView) return;

    const observer = createIntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsInView(true);
            observer?.disconnect();
          }
        });
      },
      { rootMargin: '100px' }
    );

    if (observer && imgRef.current) {
      observer.observe(imgRef.current);
    }

    return () => observer?.disconnect();
  }, [lazy, priority, isInView]);

  const handleLoad = (event: React.SyntheticEvent<HTMLImageElement>) => {
    setIsLoaded(true);
    onLoad?.(event);
  };

  const handleError = (event: React.SyntheticEvent<HTMLImageElement>) => {
    setHasError(true);
    onError?.(event);
  };

  // Generate blur data URL for placeholder
  const generateBlurDataURL = (w: number, h: number) => {
    if (blurDataURL) return blurDataURL;
    
    const canvas = document.createElement('canvas');
    canvas.width = w;
    canvas.height = h;
    const ctx = canvas.getContext('2d');
    
    if (ctx) {
      ctx.fillStyle = '#f0f0f0';
      ctx.fillRect(0, 0, w, h);
    }
    
    return canvas.toDataURL();
  };

  const imageProps = {
    src: hasError ? fallbackSrc : src,
    alt,
    quality,
    priority,
    onLoad: handleLoad,
    onError: handleError,
    className: `transition-opacity duration-300 ${
      isLoaded ? 'opacity-100' : 'opacity-0'
    } ${className || ''}`,
    ...props,
  };

  if (fill) {
    return (
      <div ref={imgRef} className="relative overflow-hidden">
        {!isLoaded && (
          <Skeleton.Image
            active
            style={{
              width: '100%',
              height: '100%',
              position: 'absolute',
              top: 0,
              left: 0,
            }}
          />
        )}
        {isInView && (
          <Image
            {...imageProps}
            fill
            sizes={sizes || '(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw'}
            placeholder={placeholder}
            blurDataURL={
              placeholder === 'blur' 
                ? generateBlurDataURL(400, 300)
                : undefined
            }
          />
        )}
      </div>
    );
  }

  return (
    <div ref={imgRef} className="relative" style={{ width, height }}>
      {!isLoaded && (
        <Skeleton.Image
          active
          style={{ width, height }}
          className="absolute top-0 left-0"
        />
      )}
      {isInView && (
        <Image
          {...imageProps}
          width={width}
          height={height}
          sizes={sizes || '(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw'}
          placeholder={placeholder}
          blurDataURL={
            placeholder === 'blur' 
              ? generateBlurDataURL(Number(width) || 400, Number(height) || 300)
              : undefined
          }
        />
      )}
    </div>
  );
};