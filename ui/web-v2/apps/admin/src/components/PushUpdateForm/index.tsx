import { Dialog } from '@headlessui/react';
import { FC, memo } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { useIsEditable } from '../../modules/me';
import { CreatableSelect, Option } from '../CreatableSelect';

export interface PushUpdateFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const PushUpdateForm: FC<PushUpdateFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const editable = useIsEditable();
    const {
      register,
      control,
      formState: { errors, isValid, isDirty, isSubmitted },
    } = methods;

    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.push.update.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.push.update.header.description)}
                </p>
              </div>
            </div>
            <div className="flex-1 flex flex-col justify-between">
              <div className="space-y-6 px-5 pt-6 pb-5 flex flex-col">
                <div className="">
                  <label htmlFor="name">
                    <span className="input-label">{f(messages.name)}</span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('name')}
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                      disabled={!editable || isSubmitted}
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="fcmApiKey" className="block">
                    <span className="input-label">
                      {f(messages.push.input.fcmApiKey)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <textarea
                      {...register('fcmApiKey')}
                      name="fcmApiKey"
                      id="fcmApiKey"
                      rows={1}
                      className="input-text w-full h-48 break-all"
                      disabled={true}
                    />
                    <p className="input-error">
                      {errors.fcmApiKey && (
                        <span role="alert">{errors.fcmApiKey.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="tags">
                    <span className="input-label">{f(messages.tags)}</span>
                  </label>
                  <Controller
                    name="tags"
                    control={control}
                    render={({ field }) => {
                      return (
                        <CreatableSelect
                          onChange={(options: Option[]) => {
                            field.onChange(options.map((o) => o.value));
                          }}
                          disabled={!editable || isSubmitted}
                          defaultValues={field.value.map((tag: string) => {
                            return {
                              value: tag,
                              label: tag,
                            };
                          })}
                        />
                      );
                    }}
                  />
                  <p className="input-error">
                    {errors.tags && (
                      <span role="alert">{errors.tags.message}</span>
                    )}
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div className="flex-shrink-0 px-4 py-4 flex justify-end">
            <div className="mr-3">
              <button
                type="button"
                className="btn-cancel"
                disabled={isSubmitted}
                onClick={onCancel}
              >
                {f(messages.button.cancel)}
              </button>
            </div>
            {editable && (
              <button
                type="button"
                className="btn-submit"
                disabled={!isValid || !isDirty || isSubmitted}
                onClick={onSubmit}
              >
                {f(messages.button.submit)}
              </button>
            )}
          </div>
        </form>
      </div>
    );
  }
);
