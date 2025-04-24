import { isSameOrBeforeDate } from 'utils/function';
import { cn } from 'utils/style';
import { getDateTimeDisplay } from '../../utils';

export const ProgressDateTimePoint = ({
  displayTime,
  displayLabel,
  className,
  conditionDate
}: {
  displayTime: string;
  displayLabel: string;
  className?: string;
  conditionDate?: Date;
}) => {
  const isSameOrBefore = isSameOrBeforeDate(
    new Date(+displayTime * 1000),
    conditionDate
  );
  return (
    <div
      className={cn(
        'flex items-center bg-gray-200',
        { 'bg-accent-pink-500 border-accent-pink-500': isSameOrBefore },
        className
      )}
    >
      <div
        className={cn(
          'size-2 relative rounded-full border bg-white border-gray-400',
          { 'bg-accent-pink-500 border-accent-pink-500': isSameOrBefore }
        )}
      >
        <span className="typo-para-medium text-gray-700 absolute -top-8 left-1/2 -translate-x-1/2 whitespace-nowrap">
          {displayLabel}
        </span>
        <div className="typo-para-tiny text-gray-500 absolute space-y-[2px] left-1/2 -translate-x-1/2 whitespace-nowrap text-center top-[18px]">
          <p>{getDateTimeDisplay(displayTime).time}</p>
          <p>{getDateTimeDisplay(displayTime).date}</p>
        </div>
      </div>
    </div>
  );
};
