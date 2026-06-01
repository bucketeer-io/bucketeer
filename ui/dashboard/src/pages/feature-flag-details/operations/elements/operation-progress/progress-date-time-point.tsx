import { isSameOrBeforeDate } from 'utils/function';
import { cn } from 'utils/style';
import { getDateTimeDisplay } from '../../utils';

export const ProgressDateTimePoint = ({
  displayTime,
  displayLabel,
  className,
  conditionDate,
  isCurrentActive
}: {
  displayTime: string;
  displayLabel: string;
  className?: string;
  conditionDate?: Date;
  isCurrentActive?: boolean;
}) => {
  const isSameOrBefore = isSameOrBeforeDate(
    new Date(+displayTime * 1000),
    conditionDate
  );
  return (
    <div
      className={cn(
        'flex items-center bg-gray-200 dark:bg-gray-600',
        {
          'bg-accent-pink-500 dark:bg-accent-pink-500 border-accent-pink-500':
            isSameOrBefore
        },
        className
      )}
    >
      <div
        className={cn(
          'size-2 relative rounded-full border bg-white dark:bg-dark-black-800 border-gray-400 dark:border-dark-gray-200',
          { 'bg-accent-pink-500 border-accent-pink-500': isSameOrBefore }
        )}
      >
        <span
          className={cn(
            'typo-para-medium text-gray-700 dark:text-dark-gray-400 absolute -top-8 left-1/2 -translate-x-1/2 whitespace-nowrap',
            {
              'text-accent-pink-500': isCurrentActive
            }
          )}
        >
          {displayLabel}
        </span>
        <div className="typo-para-tiny text-gray-500 dark:text-dark-gray-200 absolute space-y-[2px] left-1/2 -translate-x-1/2 whitespace-nowrap text-center top-[18px]">
          <p>{getDateTimeDisplay(displayTime).time}</p>
          <p>{getDateTimeDisplay(displayTime).date}</p>
        </div>
      </div>
    </div>
  );
};
