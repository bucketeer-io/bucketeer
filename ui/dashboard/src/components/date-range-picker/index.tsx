import { memo, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import {
  defaultStaticRanges,
  createStaticRanges,
  DateRange,
  DateRangeProps,
  Range
} from 'react-date-range';
import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css';
import dayjs from 'dayjs';
import { useTranslation } from 'i18n';
import { formatLongDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { IconCalendar, IconClose } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import ActionBar from './action-bar';
import './customize-date-range-picker.css';
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

const extraEnRangesLabels = [
  { label: 'Last 3 months', month: 3 },
  { label: 'Last 6 months', month: 6 },
  { label: 'Last 12 months', month: 12 }
];

export interface StaticRangeOption {
  label: string;
  isSelected: (range: Range) => boolean;
  range: () => Range;
}

interface ReactDateRangePickerProps extends Omit<DateRangeProps, 'onChange'> {
  from?: string | number;
  to?: string | number;
  onChange: (startDate?: number, endDate?: number) => void;
}

const defaultClassNames = {
  calendarWrapper: 'range__wrapper',
  months: 'range__months',
  month: 'range__months--month',
  monthName: 'range__months--monthName',
  weekDays: 'w-full',
  weekDay: 'range__weekDay',
  days: 'range__days',
  day: 'range__days--day',
  dayNumber: cn('absolute__center', 'range__days--number'),
  dayPassive: 'text-gray-500',
  dayEndOfWeek: 'range__days--dayEndOfWeek',
  startEdge: cn('absolute__center', 'range__days--edge'),
  endEdge: cn('absolute__center', 'range__days--edge'),
  dayInPreview: 'range__days--dayPreview',
  dayStartPreview: 'range__days--dayPreview',
  dayEndPreview: 'range__days--dayPreview',
  inRange: 'range__days--dayPreview'
};

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
    const popoverCloseRef = useRef<HTMLButtonElement>(null);
    const hasValue = useMemo(() => !!from && !!to, [from, to]);

    const triggerLabel = useMemo(() => {
      if (!hasValue) return t('date-range');
      const fromFormatted = formatLongDateTime({
        value: from as string,
        overrideOptions: {
          month: 'long',
          day: 'numeric',
          year: 'numeric'
        }
      });

      const toFormatted = formatLongDateTime({
        value: to as string,
        overrideOptions: {
          month: 'long',
          day: 'numeric',
          year: 'numeric'
        }
      });
      return `${fromFormatted} - ${toFormatted}`;
    }, [hasValue, from, to]);

    const staticRanges = useMemo(
      () =>
        [
          ...defaultStaticRanges.map((range, rangeIdx) => ({
            ...range,
            label: defaultEnRangesLabels[rangeIdx]
          })),
          ...createStaticRanges(
            extraEnRangesLabels.map(({ label, month }) => ({
              label,
              range: () => ({
                startDate: dayjs().subtract(month, 'month').toDate(),
                endDate: new Date()
              })
            }))
          )
        ] as StaticRangeOption[],
      []
    );
    const [range, setRange] = useState<Range>({
      ...staticRanges[0].range(),
      key: 'selection'
    });

    const staticRangeSelected = useMemo(
      () => staticRanges.find(item => item.isSelected(range)),
      [staticRanges, range]
    );

    const [isDateSelected, setIsDateSelected] = useState(false);

    const handleClear = useCallback(() => {
      setRange({
        startDate: new Date(),
        endDate: new Date(),
        key: 'selection'
      });
      onChange();
    }, [onChange]);

    const handleApply = useCallback(() => {
      const { startDate, endDate } = range;
      popoverCloseRef?.current?.click();

      if (startDate && endDate)
        onChange(
          Math.trunc(startDate.getTime() / 1000),
          Math.trunc(endDate.getTime() / 1000)
        );
      popoverCloseRef?.current?.click();
    }, [range, onChange]);

    useEffect(() => {
      if (from && to) {
        setIsDateSelected(true);
        return setRange({
          ...range,
          startDate: new Date(+from * 1000),
          endDate: new Date(+to * 1000)
        });
      }

      if (isDateSelected) {
        setIsDateSelected(false);
      }
    }, [from, to]);

    return (
      <Popover
        closeRef={popoverCloseRef}
        trigger={
          <div
            className={cn(
              'flex items-center gap-x-2 px-4 h-12 border border-gray-400 hover:shadow-border-gray-400 rounded-lg max-w-[200px] xxl:max-w-fit'
            )}
          >
            {!hasValue && (
              <Icon
                icon={IconCalendar}
                color="gray-500"
                size="sm"
                className="flex-center"
              />
            )}
            <p className="typo-para-medium text-gray-600 truncate">
              {triggerLabel}
            </p>
            {hasValue && (
              <div
                className="flex-center cursor-pointer"
                onClick={e => {
                  e.stopPropagation();
                  handleClear();
                }}
              >
                <Icon icon={IconClose} size="sm" className="flex-center" />
              </div>
            )}
          </div>
        }
        className="max-h-[501px] w-[793px] p-0"
        triggerCls="hover:!drop-shadow-none"
        align="end"
      >
        <ReactDateRangePickerComp
          {...props}
          onChange={item => {
            setRange(item.selection);
          }}
          showPreview={showPreview}
          moveRangeOnFirstSelection={moveRangeOnFirstSelection}
          months={months}
          ranges={[range]}
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
            ...defaultClassNames,
            ...props?.classNames
          }}
        />
        <ActionBar
          staticRanges={staticRanges}
          staticRangeSelected={staticRangeSelected}
          setRange={setRange}
          onCancel={() => popoverCloseRef?.current?.click()}
          onApply={handleApply}
        />
      </Popover>
    );
  }
);
