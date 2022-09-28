import { FC, memo } from 'react';

import { classNames } from '../../utils/css';

export const ListSkeleton: FC = memo(() => {
  return (
    <fieldset>
      <div className="divide-y divide-gray-300">
        {Array.from(Array(5).keys()).map((idx) => (
          <div key={idx} className="items-start px-5 py-2">
            <div className={classNames('h-full', 'animate-pulse')}>
              <div className="w-[500px] h-3 my-2 bg-gray-300 rounded-md"></div>
              <div className="w-[300px] h-3 my-2 bg-gray-300 rounded-md"></div>
              <div className="w-[500px] h-3 my-2 bg-gray-300 rounded-md"></div>
            </div>
          </div>
        ))}
      </div>
    </fieldset>
  );
});
