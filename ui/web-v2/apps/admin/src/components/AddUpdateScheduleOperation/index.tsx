import { messages } from '@/lang/messages';
import React, { FC, memo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { DatetimePicker } from '../DatetimePicker';

interface AddUpdateScheduleOperationProps {
  isSeeDetailsSelected: boolean;
}

export const AddUpdateScheduleOperation: FC<AddUpdateScheduleOperationProps> =
  memo(({ isSeeDetailsSelected }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext<any>();

    const {
      formState: { errors },
    } = methods;

    return (
      <div className="mt-1">
        <span className="input-label">{f(messages.autoOps.startDate)}</span>
        <DatetimePicker
          name="datetime.time"
          dateFormat="yyyy/MM/dd HH:mm"
          disabled={isSeeDetailsSelected}
        />
        <p className="input-error">
          {errors.datetime?.time?.message && (
            <span role="alert">{errors.datetime?.time?.message}</span>
          )}
        </p>
      </div>
    );
  });
