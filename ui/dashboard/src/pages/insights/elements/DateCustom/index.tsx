import { useMemo, useState } from 'react';
import { Range } from 'react-date-range';
import dayjs from 'dayjs';
import { useTranslation } from 'i18n';
import { truncNumber } from 'pages/audit-logs/utils';
import Button from 'components/button';
import { ReactDateRangePicker } from 'components/date-range-picker';

const MAX_RANGE_DAYS = 31;

interface InsightsDateRangePickerProps {
  onApply: (startAt: string, endAt: string) => void;
  onLabelChange: (label: string) => void;
  isOpen: boolean;
  onClose?: () => void;
}

const InsightsDateRangePicker = ({
  onApply,
  onLabelChange,
  isOpen,
  onClose
}: InsightsDateRangePickerProps) => {
  const { t } = useTranslation(['common']);

  const initRange = useMemo(
    () => ({
      from: truncNumber(
        new Date(
          dayjs()
            .subtract(MAX_RANGE_DAYS - 1, 'day')
            .toDate()
            .setHours(0, 0, 0, 0)
        ).getTime() / 1000
      ),
      to: truncNumber(
        new Date(new Date().setHours(23, 59, 59, 999)).getTime() / 1000
      )
    }),
    []
  );

  // today at end-of-day so today itself is fully selectable
  const maxDate = useMemo(
    () => new Date(new Date().setHours(23, 59, 59, 999)),
    []
  );

  const [selectedRange, setSelectedRange] = useState(initRange);
  const [pendingRange, setPendingRange] = useState<Range | null>(null);

  const formatRangeLabel = (start: number, end: number) =>
    `${dayjs(start * 1000).format('MMM D, YYYY')} - ${dayjs(end * 1000).format('MMM D, YYYY')}`;

  const rangeExceedsLimit = useMemo(() => {
    const start =
      pendingRange?.startDate ?? new Date(selectedRange.from * 1000);
    const end = pendingRange?.endDate ?? new Date(selectedRange.to * 1000);
    const diffDays = dayjs(end).diff(dayjs(start), 'day');
    return diffDays > MAX_RANGE_DAYS - 1;
  }, [pendingRange, selectedRange]);

  return (
    <ReactDateRangePicker
      className="insights-date-picker"
      from={selectedRange.from}
      to={selectedRange.to}
      maxDate={maxDate}
      getTriggerLabel={() => {}}
      isShowRange={isOpen}
      onClose={onClose}
      onChange={(startDate, endDate) => {
        if (startDate && endDate) {
          // If same day selected, extend end to 23:59:59 of that day
          const startDay = dayjs(startDate * 1000).startOf('day');
          const endDay = dayjs(endDate * 1000).startOf('day');
          const isSameDay = startDay.isSame(endDay, 'day');
          const normalizedEnd = isSameDay
            ? truncNumber(
                dayjs(startDate * 1000)
                  .endOf('day')
                  .valueOf() / 1000
              )
            : endDate;
          setSelectedRange({ from: startDate, to: normalizedEnd });
          onLabelChange(formatRangeLabel(startDate, normalizedEnd));
          onApply(startDate.toString(), normalizedEnd.toString());
        } else {
          onApply(
            startDate ? startDate.toString() : '',
            endDate ? endDate.toString() : ''
          );
        }
      }}
      onRangeChange={setPendingRange}
      renderActionBar={({ onApply: handleApply, onCancel }) => (
        <div className="sticky bottom-0 left-0 right-0 flex flex-col w-full border-t border-gray-200 bg-white">
          {rangeExceedsLimit && (
            <p className="px-5 pt-3 typo-para-small text-accent-red-500">
              {t('insights.date-range-limit', { days: MAX_RANGE_DAYS })}
            </p>
          )}
          <div className="flex items-center justify-end gap-x-4 p-5">
            <Button variant="secondary" onClick={onCancel}>
              {t('cancel')}
            </Button>
            <Button disabled={rangeExceedsLimit} onClick={handleApply}>
              {t('apply')}
            </Button>
          </div>
        </div>
      )}
    />
  );
};

export default InsightsDateRangePicker;
