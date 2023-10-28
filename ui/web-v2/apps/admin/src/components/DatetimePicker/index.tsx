import React, { FC, memo } from 'react';
import RDatePicker from 'react-datepicker';
import { Controller, useFormContext } from 'react-hook-form';

import { classNames } from '../../utils/css';
import 'react-datepicker/dist/react-datepicker.css';
import './custom-datepicker.css';

export interface DatetimePickerProps {
  name: string;
  disabled?: boolean;
}

export const DatetimePicker: FC<DatetimePickerProps> = memo(
  ({ name, disabled }) => {
    const methods = useFormContext();
    const { control } = methods;

    return (
      <Controller
        control={control}
        name={name}
        render={({ field: { onChange, value } }) => (
          <RDatePicker
            dateFormat="yyyy-MM-dd HH:mm"
            showTimeSelect
            timeIntervals={60}
            placeholderText=""
            className={classNames('input-text w-full')}
            wrapperClassName="w-full"
            onChange={onChange}
            selected={value as Date}
            disabled={disabled}
          />
        )}
      />
    );
  }
);

export interface ReactDatetimePickerProps {
  value: Date;
  disabled?: boolean;
}

export const ReactDatePicker: FC<ReactDatetimePickerProps> = memo(
  ({ value, disabled }) => {
    return (
      <RDatePicker
        selected={value as Date}
        dateFormat="yyyy-MM-dd HH:mm"
        showTimeSelect
        timeIntervals={60}
        placeholderText=""
        className={classNames('input-text w-full')}
        wrapperClassName="w-full"
        disabled={disabled}
      />
    );
  }
);
