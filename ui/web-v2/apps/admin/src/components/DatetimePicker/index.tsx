import React, { FC, memo } from 'react';
import ReactDatePicker from 'react-datepicker';
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
          <ReactDatePicker
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
