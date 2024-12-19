import type { FunctionComponent, ReactNode } from 'react';
import clsx from 'clsx';
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
      className={clsx(
        'flex items-center justify-center rounded-3xl border text-center',
        'mx-auto h-[220px] max-w-[700px] bg-gray-100',
        color === 'positive' ? 'border-gray-400' : 'border-accent-red-200'
      )}
    >
      <div>
        {icon ? (
          <div
            className={clsx(
              'mb-4 flex items-center justify-center',
              color === 'positive' ? 'text-primary-500' : 'text-accent-red-500'
            )}
          >
            <Icon size="2xl" icon={icon} />
            <img
              className="mx-auto"
              alt="Photo Select Icon"
              src="/assets/photo-upload-icon.svg"
            />
          </div>
        ) : (
          <div className="mb-4 flex items-center justify-center">
            <img
              className="mx-auto"
              alt="Photo Select Icon"
              src="/assets/photo-upload-icon.svg"
            />
          </div>
        )}
        {children}
      </div>
    </div>
  );
};

export default UploadZone;
