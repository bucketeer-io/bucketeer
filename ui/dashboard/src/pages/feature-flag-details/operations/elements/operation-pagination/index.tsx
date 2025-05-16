import { memo, FC } from 'react';
import { cn } from 'utils/style';
import { IconArrowLeft, IconArrowRight } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';

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
        <Button
          variant="secondary-2"
          className="!size-6 !p-0"
          disabled={page === 0}
          onClick={() => onPageChange(page - 1)}
        >
          <Icon icon={IconArrowLeft} size="xxs" />
        </Button>
        <div className="flex space-x-2">
          {Array(count)
            .fill('')
            .map((_, i) => (
              <div
                key={i}
                className={cn('size-2 rounded-full bg-gray-300', {
                  'w-6 bg-gray-600': page === i
                })}
              />
            ))}
        </div>
        <Button
          variant="secondary-2"
          className="!size-6 !p-0"
          disabled={page === count - 1}
          onClick={() => onPageChange(page + 1)}
        >
          <Icon icon={IconArrowRight} size="xxs" />
        </Button>
      </div>
    );
  }
);

export default OperationPagination;
