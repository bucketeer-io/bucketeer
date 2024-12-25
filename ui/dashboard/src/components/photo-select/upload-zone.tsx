import type { FunctionComponent, ReactNode } from 'react';
import { cn } from 'utils/style';
import { IconUpload } from '@icons';
import Icon from 'components/icon';

const UploadZone = ({
  color,
  icon,
  children
}: {
  color: 'positive' | 'negative';
  icon?: FunctionComponent;
  children?: ReactNode;
}) => {
  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center rounded-3xl border text-center',
        'mx-auto h-[220px] max-w-[700px] bg-gray-100',
        color === 'positive' ? 'border-gray-400' : 'border-accent-red-200'
      )}
    >
      <Icon
        size="2xl"
        icon={icon || IconUpload}
        color={color === 'positive' ? 'primary-600' : 'accent-red-500'}
      />

      {children}
    </div>
  );
};

export default UploadZone;
