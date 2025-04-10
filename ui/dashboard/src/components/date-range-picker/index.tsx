import { memo, useCallback, useEffect, useMemo, useState } from 'react';
import { DateRangePicker, DateRangePickerProps, Range } from 'react-date-range';
import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css';
import en from 'date-fns/locale/en-US';
import ja from 'date-fns/locale/ja';
import dayjs from 'dayjs';
import { useToggleOpen } from 'hooks';
import { getLanguage, useTranslation } from 'i18n';
import { formatLongDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { IconCalendar } from '@icons';
import { truncNumber } from 'pages/audit-logs/utils';
import Button from 'components/button';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import ActionBar from './action-bar';
import './customize-date-range-picker.css';
import CustomizeNavigator from './customize-navigator';

const ReactDateRangePickerComp =
  DateRangePicker as unknown as React.FC<DateRangePickerProps>;

interface DefaultRangeLabel {
  label: string;
  value: number;
  type: 'today' | 'day' | 'month' | 'all-time';
}

const getDefaultRangesLabels = (
  isLanguageJapanese: boolean
): DefaultRangeLabel[] => [
  { label: isLanguageJapanese ? '今日' : 'Today', value: 0, type: 'today' },
  { label: isLanguageJapanese ? '昨日' : 'Yesterday', value: 1, type: 'day' },
  {
    label: isLanguageJapanese ? '直近7日' : 'Last 7 days',
    value: 7,
    type: 'day'
  },
  {
    label: isLanguageJapanese ? '直近14日' : 'Last 14 days',
    value: 14,
    type: 'day'
  },
  {
    label: isLanguageJapanese ? '直近30日' : 'Last 30 days',
    value: 1,
    type: 'month'
  },
  {
    label: isLanguageJapanese ? '過去3ヶ月' : 'Last 90 days',
    value: 3,
    type: 'month'
  },
  {
    label: isLanguageJapanese ? '過去12ヶ月' : 'Last 12 months',
    value: 12,
    type: 'month'
  },
  {
    label: isLanguageJapanese ? '全期間' : 'All Time',
    value: 0,
    type: 'all-time'
  }
];

export interface StaticRangeOption {
  type: DefaultRangeLabel['type'];
  label: string;
  isSelected: (range: Range) => boolean;
  range: () => Range;
}

interface ReactDateRangePickerProps
  extends Omit<DateRangePickerProps, 'onChange'> {
  from?: string | number;
  to?: string | number;
  isAllTime?: boolean;
  onChange: (startDate?: number, endDate?: number) => void;
}

const defaultClassNames = {
  calendarWrapper: 'range__wrapper',
  dateDisplayWrapper: 'range__date-display-wrapper',
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
    isAllTime,
    months = 2,
    showPreview = true,
    moveRangeOnFirstSelection = false,
    direction = 'horizontal',
    rangeColors = ['#5d3597'],
    showDateDisplay = true,
    onChange,
    ...props
  }) => {
    const { t } = useTranslation(['common']);
    const isLanguageJapanese = getLanguage() === 'ja';

    const [isOpenRangePicker, onOpenRangePicker, onCloseRangePicker] =
      useToggleOpen(false);

    const hasValue = useMemo(() => !!from && !!to, [from, to]);

    const staticRanges = useMemo(() => {
      const defaultRangesLabels = getDefaultRangesLabels(isLanguageJapanese);

      return defaultRangesLabels.map(({ label, type, value }) => {
        const rangeFn = () => ({
          startDate: ['all-time', 'today'].includes(type)
            ? new Date(new Date().setHours(0, 0, 0, 0))
            : new Date(
                dayjs()
                  .subtract(value, type === 'month' ? 'month' : 'day')
                  .toDate()
                  .setHours(0, 0, 0, 0)
              ),
          endDate: new Date(
            new Date().setHours(23, type === 'all-time' ? 58 : 59, 59, 999)
          )
        });
        return {
          type,
          label,
          hasCustomRendering: true,
          range: rangeFn,
          isSelected: (range: Range) => {
            if (range?.startDate && range?.endDate) {
              const currentRange = rangeFn();
              return (
                truncNumber(currentRange.startDate.getTime() / 1000) ===
                  truncNumber(range.startDate.getTime() / 1000) &&
                truncNumber(currentRange.endDate.getTime() / 1000) ===
                  truncNumber(range.endDate.getTime() / 1000)
              );
            }
            return false;
          }
        };
      });
    }, [isLanguageJapanese]);

    const [range, setRange] = useState<Range>({
      ...staticRanges[0].range(),
      key: 'selection'
    });
    const [staticRangeSelected, setStaticRangeSelected] = useState<
      StaticRangeOption | undefined
    >(undefined);

    const triggerLabel = useMemo(() => {
      if (staticRangeSelected?.type === 'all-time')
        return staticRanges.at(-1)?.label;
      if (!hasValue) return t('date-range');
      const selectedRange = staticRanges.find(item => item.isSelected(range));
      if (selectedRange) return selectedRange.label;
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
    }, [staticRanges, hasValue, from, to, range, staticRangeSelected]);

    const handleApply = useCallback(() => {
      if (staticRangeSelected?.type === 'all-time') {
        onChange(undefined, undefined);
      } else {
        const { startDate, endDate } = range;
        if (startDate && endDate)
          onChange(
            truncNumber(startDate.getTime() / 1000),
            truncNumber(endDate.getTime() / 1000)
          );
      }
      onCloseRangePicker();
    }, [range, staticRangeSelected, onChange]);

    const handleSetRange = useCallback(() => {
      if (isAllTime) {
        setStaticRangeSelected(staticRanges.at(-1));
        return setRange({
          ...range,
          startDate: staticRanges.at(-1)!.range().startDate,
          endDate: staticRanges.at(-1)!.range().endDate
        });
      }
      if (staticRangeSelected?.type === 'all-time')
        setStaticRangeSelected(undefined);
      if (from && to) {
        setRange({
          ...range,
          startDate: new Date(truncNumber(+from * 1000)),
          endDate: new Date(truncNumber(+to * 1000))
        });
      }
    }, [from, to, range, staticRanges, staticRangeSelected]);

    useEffect(() => {
      handleSetRange();
    }, [from, to, isAllTime]);

    return (
      <>
        <Button
          variant="secondary-2"
          className="border border-gray-400 shadow-none hover:shadow-border-gray-400 rounded-lg !max-w-[200px] xxl:!max-w-fit"
          onClick={onOpenRangePicker}
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
        </Button>
        <DialogModal
          title="Date Range Picker"
          isShowHeader={false}
          isOpen={isOpenRangePicker}
          onClose={() => {
            handleSetRange();
            onCloseRangePicker();
          }}
          className="w-[820px] p-0 rounded-lg overflow-y-auto"
        >
          <ReactDateRangePickerComp
            {...props}
            onChange={item => setRange(item.selection)}
            showPreview={showPreview}
            moveRangeOnFirstSelection={moveRangeOnFirstSelection}
            months={months}
            ranges={[range]}
            direction={direction}
            rangeColors={rangeColors}
            staticRanges={staticRanges}
            showDateDisplay={showDateDisplay}
            dateDisplayFormat="yyyy/MM/dd"
            inputRanges={[]}
            renderStaticRangeLabel={staticRange => (
              <div
                onClick={() =>
                  setStaticRangeSelected(staticRange as StaticRangeOption)
                }
                className={cn('flex items-center size-full pl-3', {
                  'range-selected': staticRange.isSelected(range)
                })}
              >
                {staticRange.label}
              </div>
            )}
            navigatorRenderer={(currFocusedDate, changeShownDate) => (
              <CustomizeNavigator
                currFocusedDate={currFocusedDate}
                changeShownDate={changeShownDate}
              />
            )}
            locale={isLanguageJapanese ? ja : en}
            classNames={{
              ...defaultClassNames,
              ...props?.classNames
            }}
          />

          <ActionBar
            onCancel={() => {
              handleSetRange();
              onCloseRangePicker();
            }}
            onApply={handleApply}
          />
        </DialogModal>
      </>
    );
  }
);
