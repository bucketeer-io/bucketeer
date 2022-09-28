import { FC, memo } from 'react';

import { classNames } from '../../utils/css';

export const DetailSkeleton: FC = memo(() => {
  return (
    <div
      className={classNames(
        'h-full',
        'flex flex-col',
        'animate-pulse justify-center space-y-5'
      )}
    >
      <div className="w-32 h-12 rounded-md bg-gray-300"></div>
      <div className="w-full h-16 rounded-md bg-gray-300"></div>
      <div className="w-full h-48 rounded-md bg-gray-300"></div>
    </div>
  );
});
