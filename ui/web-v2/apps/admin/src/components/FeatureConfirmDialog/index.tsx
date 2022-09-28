import { Dialog } from '@headlessui/react';
import { FC } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { FEATURE_UPDATE_COMMENT_MAX_LENGTH } from '../../constants/feature';
import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';
import { CheckBox } from '../CheckBox';
import { Modal } from '../Modal';

interface FeatureConfirmDialogProps {
  open: boolean;
  handleSubmit: () => void;
  onClose: () => void;
  title: string;
  description: string;
  displayResetSampling?: boolean;
}

export const FeatureConfirmDialog: FC<FeatureConfirmDialogProps> = ({
  open,
  handleSubmit,
  onClose,
  title,
  description,
  displayResetSampling,
}) => {
  const { formatMessage: f } = useIntl();
  const methods = useFormContext();
  const {
    register,
    control,
    formState: { errors, isSubmitting, isDirty, isValid },
  } = methods;
  return (
    <Modal open={open} onClose={onClose}>
      <Dialog.Title
        as="h3"
        className="text-lg font-medium leading-6 text-gray-900"
      >
        {title}
      </Dialog.Title>
      <div className="mt-2">
        <p className="text-sm text-gray-500">{description}</p>
      </div>

      <div className="mt-5">
        <label
          htmlFor="about"
          className="block text-sm font-medium text-gray-700"
        >
          {f(messages.feature.updateComment)}
        </label>
        <div className="mt-1">
          <textarea
            {...register('comment', {
              required: true,
              maxLength: FEATURE_UPDATE_COMMENT_MAX_LENGTH,
            })}
            id="description"
            rows={3}
            className="input-text w-full"
          />
          <p className="input-error">
            {errors.comment && (
              <span role="alert">{errors.comment.message}</span>
            )}
          </p>
        </div>

        {displayResetSampling && (
          <div className="mt-3 flex space-x-2 items-center">
            <Controller
              name="resetSampling"
              control={control}
              render={({ field }) => {
                return (
                  <CheckBox
                    id="resample"
                    value={'resample'}
                    onChange={(_: string, checked: boolean): void =>
                      field.onChange(checked)
                    }
                    defaultChecked={false}
                  />
                );
              }}
            />
            <label htmlFor="resample" className={classNames('input-label')}>
              {f(messages.feature.resetRandomSampling)}
            </label>
          </div>
        )}
      </div>

      <div className="pt-5">
        <div className="flex justify-end">
          <button
            type="button"
            className="btn-cancel mr-3"
            disabled={false}
            onClick={onClose}
          >
            {f(messages.button.cancel)}
          </button>
          <button
            type="button"
            className="btn-submit"
            disabled={!isDirty || !isValid || isSubmitting}
            onClick={handleSubmit}
          >
            {f(messages.button.submit)}
          </button>
        </div>
      </div>
    </Modal>
  );
};
