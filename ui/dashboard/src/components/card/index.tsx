import { PropsWithChildren } from 'react';
import { cn } from 'utils/style';

export type CardProps = PropsWithChildren & {
  tag?: keyof JSX.IntrinsicElements;
  className?: string;
};

const Card = ({ children, tag: Tag = 'div', className }: CardProps) => {
  return (
    <Tag className={cn('shadow-card rounded-lg bg-white', className)}>
      {children}
    </Tag>
  );
};

export default Card;
