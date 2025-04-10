// theme css file
import { memo, useEffect, useState } from 'react';
import {
  defaultStaticRanges,
  createStaticRanges,
  DateRange,
  DateRangeProps
} from 'react-date-range';
import 'react-date-range/dist/styles.css';
// main css file
import 'react-date-range/dist/theme/default.css';
import { useTranslation } from 'i18n';
import { IconCalendar } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import CustomizeNavigator from './customize-navigator';

const ReactDateRangePickerComp =
  DateRange as unknown as React.FC<DateRangeProps>;

const defaultEnRangesLabels = [
  'Today',
  'Yesterday',
  'This week',
  'Last week',
  'This month',
  'Last month'
];

interface ReactDateRangePickerProps extends DateRangeProps {
  from?: string | number;
  to?: string | number;
}

export const ReactDateRangePicker: React.FC<ReactDateRangePickerProps> = memo(
  ({
    from,
    to,
    months = 2,
    showPreview = true,
    moveRangeOnFirstSelection = false,
    direction = 'horizontal',
    rangeColors = ['#5d3597'],
    showDateDisplay = false,
    onChange,
    ...props
  }) => {
    const { t } = useTranslation(['common']);

    const [ranges, setRanges] = useState([
      {
        startDate: new Date(),
        endDate: new Date(),
        key: 'selection'
      }
    ]);
    const [isDateSelected, setIsDateSelected] = useState(false);

    const staticRanges = defaultStaticRanges.map((range, rangeIdx) => ({
      ...range,
      label: defaultEnRangesLabels[rangeIdx]
    }));
    console.log(staticRanges[0].range());
    useEffect(() => {
      if (from && to) {
        setIsDateSelected(true);
        return setRanges([
          {
            ...ranges[0],
            startDate: new Date(+from * 1000),
            endDate: new Date(+to * 1000)
          }
        ]);
      }

      if (isDateSelected) {
        setIsDateSelected(false);
      }
    }, [from, to]);

    return (
      <Popover
        trigger={
          <div className="flex items-center gap-x-2 px-4 h-12 border border-gray-200 rounded-lg">
            <Icon icon={IconCalendar} color="gray-500" size="sm" />
            <p className="typo-para-medium text-gray-600">{t('date-range')}</p>
          </div>
        }
        className="max-h-[501px] w-[777px] max-w-[777px] p-0"
        align="end"
      >
        <ReactDateRangePickerComp
          {...props}
          onChange={onChange}
          showPreview={showPreview}
          moveRangeOnFirstSelection={moveRangeOnFirstSelection}
          months={months}
          ranges={ranges}
          direction={direction}
          rangeColors={rangeColors}
          showDateDisplay={showDateDisplay}
          navigatorRenderer={(currFocusedDate, changeShownDate) => (
            <CustomizeNavigator
              currFocusedDate={currFocusedDate}
              changeShownDate={changeShownDate}
            />
          )}
          classNames={{
            calendarWrapper: '!font-sofia-pro w-full',
            months:
              'grid grid-cols-2 justify-between w-full divide-x divide-gray-200 overflow-visible',
            month: 'flex flex-col gap-y-6 col-span-1 w-full',
            monthName: 'typo-para-medium text-gray-600 p-0',
            weekDays: 'w-full',
            weekDay:
              'flex-center !min-w-[49px] h-[21px] uppercase typo-para-small text-gray-500 px-1.5 first:pl-0 last:pr-0',
            days: 'w-full gap-y-4',
            day: 'w-[49px] h-8 typo-para-small text-gray-600 px-1.5 first:pl-0 last:pr-0',
            dayNumber: 'size-8 absolute top-0 right-2 bottom-0 left-2',
            dayPassive: 'text-gray-500',
            startEdge:
              'size-8 rounded-full absolute top-0 right-1.5 bottom-0 left-1.5',
            endEdge:
              'size-8 rounded-full absolute top-0 right-1.5 bottom-0 left-1.5'
            // dayInPreview: 'h-8 absolute top-0 right-0 bottom-0 left-0 -translate-x-1/2 -translate-y-1/2'
          }}
        />
      </Popover>
    );
  }
);
