import { memo, useCallback, useEffect, useMemo, useState } from 'react';
import {
  defaultStaticRanges,
  createStaticRanges,
  DateRangePicker,
  DateRangePickerProps,
  Range
} from 'react-date-range';
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

const defaultEnRangesLabels = ['Today', 'Yesterday', 'Last week', 'Last month'];
const defaultJaRangesLabels = ['今日', '昨日', '先週', '先月'];

const extraEnRangesLabels = [
  { label: 'Last 3 months', month: 3 },
  { label: 'Last 6 months', month: 6 },
  { label: 'Last 12 months', month: 12 }
];

const extraJaRangesLabels = [
  { label: '過去3ヶ月', month: 3 },
  { label: '過去6ヶ月', month: 6 },
  { label: '過去12ヶ月', month: 12 }
];

export interface StaticRangeOption {
  label: string;
  isSelected: (range: Range) => boolean;
  range: () => Range;
}

interface ReactDateRangePickerProps
  extends Omit<DateRangePickerProps, 'onChange'> {
  from?: string | number;
  to?: string | number;
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

    const defaultLabels = useMemo(
      () =>
        isLanguageJapanese ? defaultJaRangesLabels : defaultEnRangesLabels,
      [isLanguageJapanese]
    );
    const extraLabels = useMemo(
      () => (isLanguageJapanese ? extraJaRangesLabels : extraEnRangesLabels),
      [isLanguageJapanese]
    );

    const staticRanges = useMemo(
      () =>
        [
          ...defaultStaticRanges
            .filter(
              range =>
                !['This Week', 'This Month'].includes(range?.label as string)
            )
            .map((range, rangeIdx) => ({
              ...range,
              label: defaultLabels[rangeIdx]
            })),
          ...createStaticRanges(
            extraLabels.map(({ label, month }) => ({
              label,
              range: () => ({
                startDate: dayjs().subtract(month, 'month').toDate(),
                endDate: new Date()
              })
            }))
          )
        ] as StaticRangeOption[],
      [defaultLabels, extraLabels]
    );

    const triggerLabel = useMemo(() => {
      if (!hasValue) return t('date-range');
      const selectedRange = staticRanges.find(item => {
        const { startDate, endDate } = item.range();

        if (startDate && endDate) {
          return (
            +from! === truncNumber(startDate.getTime() / 1000) &&
            +to! === truncNumber(endDate.getTime() / 1000)
          );
        }
      });

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
    }, [staticRanges, hasValue, from, to]);

    const [range, setRange] = useState<Range>({
      ...staticRanges[0].range(),
      key: 'selection'
    });

    const handleApply = useCallback(() => {
      const { startDate, endDate } = range;

      if (startDate && endDate)
        onChange(
          truncNumber(startDate.getTime() / 1000),
          truncNumber(endDate.getTime() / 1000)
        );
      onCloseRangePicker();
    }, [range, onChange]);

    useEffect(() => {
      if (from && to) {
        return setRange({
          ...range,
          startDate: new Date(truncNumber(+from * 1000)),
          endDate: new Date(truncNumber(+to * 1000))
        });
      }
    }, [from, to]);

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
          onClose={onCloseRangePicker}
          className="w-[820px] p-0 rounded-lg overflow-y-auto"
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
            staticRanges={staticRanges}
            showDateDisplay={showDateDisplay}
            dateDisplayFormat="yyyy/MM/dd"
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
          <ActionBar onCancel={onCloseRangePicker} onApply={handleApply} />
        </DialogModal>
      </>
    );
  }
);
