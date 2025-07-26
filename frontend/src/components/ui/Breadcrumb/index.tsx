import React from 'react';
import { Breadcrumb as AntBreadcrumb } from 'antd';
import { HomeOutlined } from '@ant-design/icons';
import Link from 'next/link';

export interface BreadcrumbItem {
  title: string;
  href?: string;
  icon?: React.ReactNode;
}

export interface BreadcrumbProps {
  items: BreadcrumbItem[];
  showHome?: boolean;
  homeHref?: string;
  className?: string;
}

export const Breadcrumb: React.FC<BreadcrumbProps> = ({
  items,
  showHome = true,
  homeHref = '/',
  className = '',
}) => {
  const breadcrumbItems = [
    ...(showHome
      ? [
          {
            title: (
              <Link href={homeHref} className="flex items-center space-x-1">
                <HomeOutlined />
                <span>ホーム</span>
              </Link>
            ),
          },
        ]
      : []),
    ...items.map((item) => ({
      title: item.href ? (
        <Link href={item.href} className="flex items-center space-x-1">
          {item.icon}
          <span>{item.title}</span>
        </Link>
      ) : (
        <span className="flex items-center space-x-1">
          {item.icon}
          <span>{item.title}</span>
        </span>
      ),
    })),
  ];

  return (
    <AntBreadcrumb
      items={breadcrumbItems}
      className={`mb-4 ${className}`}
    />
  );
};

export default Breadcrumb;