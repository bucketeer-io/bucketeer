import { memo, useMemo } from 'react';
import { format, getYear, setMonth, setYear } from 'date-fns';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';

interface Props {
  currFocusedDate: Date;
  changeShownDate: (
    value: Date | number | string,
    mode?: 'set' | 'setYear' | 'setMonth' | 'monthOffset'
  ) => void;
}

const months = Array.from({ length: 12 }, (_, i) =>
  format(new Date(2000, i), 'MMMM')
);

const years = Array.from(
  { length: 30 },
  (_, i) => new Date().getFullYear() - 15 + i
);

const CustomizeNavigator = memo(
  ({ currFocusedDate, changeShownDate }: Props) => {
    const prevMonthButtonDisabled = useMemo(
      () =>
        currFocusedDate.getMonth() === 0 &&
        currFocusedDate.getFullYear() === years[0],
      [currFocusedDate]
    );
    const nextMonthButtonDisabled = useMemo(
      () =>
        currFocusedDate.getMonth() === 11 &&
        currFocusedDate.getFullYear() === years.at(-1),
      [currFocusedDate]
    );

    return (
      <div className={cn('w-full relative z-[1000] p-5 pb-0')}>
        <div className="flex items-center justify-between w-full border-b border-gray-200 pb-5">
          <Button
            type="button"
            size={'icon-sm'}
            variant={'secondary-2'}
            onClick={() => {
              changeShownDate(
                currFocusedDate.setMonth(currFocusedDate.getMonth() - 1)
              );
            }}
            disabled={prevMonthButtonDisabled}
            className="size-8"
          >
            <Icon
              icon={IconChevronRight}
              className="rotate-180"
              size={'xs'}
              color="gray-500"
            />
          </Button>
          <div className="flex items-center gap-x-2">
            <Dropdown
              options={months.map((item, index) => ({
                label: item,
                value: index
              }))}
              labelCustom={months[currFocusedDate.getMonth()]}
              value={months[currFocusedDate.getMonth()]}
              onChange={value => {
                const newDate = setMonth(currFocusedDate, +value);
                changeShownDate(newDate);
              }}
              isTruncate={false}
              className="!shadow-none !border-none p-0"
              contentClassName="min-w-[120px]"
              sideOffsetContent={-8}
            />
            <Dropdown
              options={years.map(item => ({
                label: item.toString(),
                value: item
              }))}
              value={getYear(currFocusedDate)}
              onChange={value => {
                const newDate = setYear(currFocusedDate, +value);
                changeShownDate(newDate);
              }}
              className="!shadow-none !border-none p-0"
              contentClassName="min-w-[120px]"
              sideOffsetContent={-8}
            />
          </div>
          <Button
            type="button"
            size={'icon-sm'}
            variant={'secondary-2'}
            onClick={() =>
              changeShownDate(
                currFocusedDate.setMonth(currFocusedDate.getMonth() + 1)
              )
            }
            disabled={nextMonthButtonDisabled}
            className="size-8"
          >
            <Icon icon={IconChevronRight} size={'xs'} color="gray-500" />
          </Button>
        </div>
      </div>
    );
  }
);

export default CustomizeNavigator;
