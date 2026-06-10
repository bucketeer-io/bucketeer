import { memo, useRef, useState, useCallback } from 'react';
import DatePicker, { DatePickerProps } from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { useTranslation } from 'react-i18next';
import { useScreen } from 'hooks/use-screen';
import { cn } from 'utils/style';
import { IconCalendar } from '@icons';
import Drawer from 'components/drawer';
import Icon from 'components/icon';
import './custom-datepicker.css';
import CustomizeHeader from './customize-header';

type SingleDatePickerProps = Extract<
  DatePickerProps,
  { selectsRange?: never; selectsMultiple?: never }
>;
type SingleDateOnChange = NonNullable<SingleDatePickerProps['onChange']>;

type ReactDatetimePickerProps = DatePickerProps;

export const ReactDatePicker = memo<ReactDatetimePickerProps>(
  ({
    selected,
    disabled,
    dateFormat = 'yyyy/MM/dd HH:mm',
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
    onChange,
    ...props
  }) => {
    const { t } = useTranslation(['form']);
    const ref = useRef<DatePicker>(null);
    const { isMobile } = useScreen();
    const [drawerOpen, setDrawerOpen] = useState(false);
    const handleOpenDrawer = useCallback(() => {
      if (!disabled) setDrawerOpen(true);
    }, [disabled]);

    if (isMobile) {
      return (
        <>
          <DatePicker
            {...(props as SingleDatePickerProps)}
            ref={ref}
            selected={selected}
            dateFormat={dateFormat}
            showTimeSelect={showTimeSelect}
            showIcon={showIcon}
            icon={
              <div className="flex-center" onClick={handleOpenDrawer}>
                {icon}
              </div>
            }
            calendarIconClassName={cn(
              'flex-center top-1/2 -translate-y-1/2 right-1',
              calendarIconClassName
            )}
            className={cn(
              'typo-para-medium !py-[11px] !pl-4 !pr-10 w-full disabled:border-gray-400 disabled:bg-gray-100',
              className
            )}
            placeholderText={placeholderText || t('select-date')}
            wrapperClassName={cn('flex items-center w-full', wrapperClassName)}
            disabled={disabled}
            open={false}
            onChange={onChange as SingleDateOnChange}
          />

          <Drawer
            open={drawerOpen}
            onClose={() => setDrawerOpen(false)}
            side="bottom"
          >
            <DatePicker
              {...(props as SingleDatePickerProps)}
              ref={ref}
              selected={selected}
              onChange={(date, event) => {
                (onChange as SingleDateOnChange | undefined)?.(date, event);
                setDrawerOpen(false);
              }}
              inline
              showTimeSelect={showTimeSelect}
              timeIntervals={timeIntervals}
              dateFormat={dateFormat}
              renderCustomHeader={headerProps => (
                <CustomizeHeader {...headerProps} />
              )}
            />
          </Drawer>
        </>
      );
    }

    return (
      <DatePicker
        {...(props as SingleDatePickerProps)}
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
              'cursor-pointer': toggleCalendarOnIconClick && !disabled
            })}
          >
            {icon}
          </div>
        }
        calendarIconClassName={cn(
          'flex-center top-1/2 -translate-y-1/2 right-1',
          calendarIconClassName
        )}
        className={cn(
          'typo-para-medium !py-[11px] !pl-4 !pr-10 w-full disabled:border-gray-400 disabled:bg-gray-100',
          className
        )}
        timeIntervals={timeIntervals}
        placeholderText={placeholderText || t('select-date')}
        wrapperClassName={cn('flex items-center w-full', wrapperClassName)}
        disabled={disabled}
        popperPlacement={popperPlacement}
        toggleCalendarOnIconClick={toggleCalendarOnIconClick}
        renderCustomHeader={headerProps => <CustomizeHeader {...headerProps} />}
        onChange={onChange as SingleDateOnChange}
      />
    );
  }
);
