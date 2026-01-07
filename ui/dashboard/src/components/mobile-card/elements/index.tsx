import React, { ReactNode, HTMLAttributes } from 'react';
import clsx from 'clsx';
import { cn } from 'utils/style';

export interface CardProps extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
}

const CardRoot: React.FC<CardProps> = ({ children, className, ...props }) => {
  return (
    <div
      className={cn(
        'bg-white rounded-2xl p-5 shadow-card relative overflow-hidden transition-all active:scale-[0.99]',
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
};

export interface CardStatusBarProps {
  colorClass: string;
}

const StatusBar: React.FC<CardStatusBarProps> = ({ colorClass }) => (
  <div className={clsx('absolute left-0 top-0 bottom-0 w-1', colorClass)} />
);

export interface CardHeaderProps {
  icon?: ReactNode;
  title?: string;
  subtitle?: string;
  children?: ReactNode;
  triger?: ReactNode;
}

const Header: React.FC<CardHeaderProps> = ({
  icon,
  title,
  subtitle,
  triger,
  children
}) => (
  <div className="flex items-center justify-between gap-4 mb-4 typo-head-semi-medium">
    <div className="flex items-center gap-x-4">
      {icon && (
        <div className="min-w-12 h-12 bg-primary-50 text-primary-500 rounded-xl flex items-center justify-center shadow-sm">
          {icon}
        </div>
      )}
      {triger ? (
        triger
      ) : (
        <div className="flex gap-x-4">
          <div className="flex-1 min-w-0">
            <h3 className="truncate">{title}</h3>
            {subtitle && (
              <div className="typo-para-small text-gray-500 line-clamp-2 mt-1">
                {subtitle}
              </div>
            )}
          </div>
        </div>
      )}
    </div>

    {children}
  </div>
);

export interface CardIconProps {
  icon: string;
}

export interface CardTitleProps {
  title: string;
  subtitle?: string;
}

export interface CardActionProps {
  children: ReactNode;
}

const Action: React.FC<CardActionProps> = ({ children }) => (
  <div className="-mr-2 -mt-2">{children}</div>
);

export interface CardMetaProps {
  children: ReactNode;
}

const Meta: React.FC<CardMetaProps> = ({ children }) => (
  <div className="rounded-xl py-3 mb-4">{children}</div>
);

export interface CardFooterProps {
  left?: ReactNode;
  right?: ReactNode;
}

const Footer: React.FC<CardFooterProps> = ({ left, right }) => (
  <div
    className={cn(
      'flex items-center pl-1',
      !!right && !left ? 'justify-end' : 'justify-between'
    )}
  >
    {left}
    {right}
  </div>
);

export const Card = Object.assign(CardRoot, {
  StatusBar,
  Header,
  Action,
  Meta,
  Footer
});
