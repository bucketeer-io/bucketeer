import * as React from 'react';
import { Link, LinkComponentProps } from '@tanstack/react-router';
import { cn } from 'utils/style';

const Tabs = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('flex-1 flex h-full flex-col', className)}
    {...props}
  />
));

const TabsList = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      'inline-flex w-full border-b border-gray-300 items-center',
      className
    )}
    {...props}
  />
));

const TabsLink = React.forwardRef<HTMLAnchorElement, LinkComponentProps>(
  ({ className, to, ...props }, ref) => (
    <Link
      to={to}
      className={cn(
        'flex-center whitespace-nowrap typo-para-medium text-gray-500',
        'py-1 px-4 border-b-2 border-transparent transition-all tab-item',
        className
      )}
      {...props}
      ref={ref}
    />
  )
);

const TabsContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('mt-4 flex flex-col flex-1', className)}
    {...props}
  />
));

export { Tabs, TabsList, TabsLink, TabsContent };
