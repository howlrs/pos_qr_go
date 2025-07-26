import { CardProps as AntCardProps } from 'antd';

export interface CardProps extends AntCardProps {
  shadow?: 'none' | 'small' | 'medium' | 'large';
  padding?: 'none' | 'small' | 'medium' | 'large';
}

export type CardShadow = CardProps['shadow'];
export type CardPadding = CardProps['padding'];