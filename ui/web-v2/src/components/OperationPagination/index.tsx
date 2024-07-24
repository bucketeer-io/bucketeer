import React, { memo, FC } from 'react';
import { classNames } from '../../utils/css';
import {
  ArrowNarrowLeftIcon,
  ArrowNarrowRightIcon
} from '@heroicons/react/solid';

interface OperationPaginationProps {
  page: number;
  count: number;
  onPageChange: (page: number) => void;
}

const OperationPagination: FC<OperationPaginationProps> = memo(
  ({ page, count, onPageChange }) => {
    if (count <= 1) {
      return null;
    }
    return (
      <div className="mt-4 flex justify-between items-center">
        <button
          className={classNames(
            'p-1.5 rounded border',
            page === 0 && 'opacity-50 cursor-not-allowed'
          )}
          disabled={page === 0}
          onClick={() => onPageChange(page - 1)}
        >
          <ArrowNarrowLeftIcon width={16} className="text-gray-400" />
        </button>
        <div className="flex space-x-2">
          {Array(count)
            .fill('')
            .map((_, i) =>
              page === i ? (
                <div
                  key={i}
                  className="w-[24px] h-[8px] rounded-full bg-gray-400"
                />
              ) : (
                <div
                  key={i}
                  className="w-[8px] h-[8px] rounded-full bg-gray-200"
                />
              )
            )}
        </div>
        <button
          className={classNames(
            'p-1.5 rounded border',
            page === count - 1 && 'opacity-50 cursor-not-allowed'
          )}
          disabled={page === count - 1}
          onClick={() => onPageChange(page + 1)}
        >
          <ArrowNarrowRightIcon width={16} className="text-gray-400" />
        </button>
      </div>
    );
  }
);

export default OperationPagination;
