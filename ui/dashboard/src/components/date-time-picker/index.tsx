import { FC, memo, useRef } from 'react';
import DatePicker, { DatePickerProps } from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { useTranslation } from 'react-i18next';
import { cn } from 'utils/style';
import { IconCalendar } from '@icons';
import Icon from 'components/icon';
import './custom-datepicker.css';

type ReactDatetimePickerProps = DatePickerProps;

export const ReactDatePicker: FC<ReactDatetimePickerProps> = memo(
  ({
    selected,
    disabled,
    dateFormat = 'yyyy-MM-dd HH:mm',
    timeIntervals = 60,
    showIcon = true,
    showTimeSelect = true,
    icon = (
      <Icon icon={IconCalendar} color="gray-600" className="flex-center" />
    ),
    wrapperClassName,
    placeholderText,
    calendarIconClassName,
    popperPlacement = 'bottom-start',
    toggleCalendarOnIconClick = true,
    className,
    ...props
  }) => {
    const { t } = useTranslation(['form']);
    const ref = useRef<DatePicker>(null);

    return (
      <DatePicker
        ref={ref}
        selected={selected}
        dateFormat={dateFormat}
        showTimeSelect={showTimeSelect}
        showIcon={showIcon}
        icon={
          <div
            onClick={() =>
              toggleCalendarOnIconClick && ref.current?.setOpen(true)
            }
            className={cn('flex-center', {
              'cursor-pointer': toggleCalendarOnIconClick
            })}
          >
            {icon}
          </div>
        }
        calendarIconClassName={cn(
          'flex-center top-1/2 -translate-y-1/2 right-1',
          calendarIconClassName
        )}
        className={cn('!py-[11px] !pl-4 !pr-10 w-full', className)}
        timeIntervals={timeIntervals}
        placeholderText={placeholderText || t('select-date')}
        wrapperClassName={cn('flex items-center w-full', wrapperClassName)}
        disabled={disabled}
        popperPlacement={popperPlacement}
        toggleCalendarOnIconClick={toggleCalendarOnIconClick}
        {...props}
      />
    );
  }
);
