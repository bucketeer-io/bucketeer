import React from 'react';
import clsx from 'clsx';

export interface StatusBadgeProps {
  label: string;
  className: string;
  ping?: boolean;
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({
  label,
  ping = false
}) => {
  return (
    <span
      className={clsx(
        'inline-flex items-center gap-1.5 rounded-md typo-para-small font-bold'
      )}
    >
      <span className="relative flex h-2 w-2">
        {ping && (
          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-accent-red-500 opacity-75" />
        )}
        <span className="relative inline-flex rounded-full h-2 w-2 bg-accent-red-500" />
      </span>
      {label}
    </span>
  );
};
