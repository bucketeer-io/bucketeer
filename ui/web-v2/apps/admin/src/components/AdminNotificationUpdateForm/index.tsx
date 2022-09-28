import { Dialog } from '@headlessui/react';
import { FC, memo, useState } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { SOURCE_TYPE_ITEMS } from '../../constants/adminNotification';
import { messages } from '../../lang/messages';
import { useIsEditable } from '../../modules/me';
import { CheckBoxList } from '../CheckBoxList';

export interface AdminNotificationUpdateFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const AdminNotificationUpdateForm: FC<AdminNotificationUpdateFormProps> =
  memo(({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const editable = useIsEditable();
    const {
      register,
      control,
      getValues,
      formState: { errors, isValid, isDirty, isSubmitted },
    } = methods;

    const [defaultValues] = useState(() =>
      SOURCE_TYPE_ITEMS.filter((item) =>
        getValues().sourceTypes.includes(Number(item.value))
      )
    );

    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.notification.update.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.notification.update.header.description)}
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
              </div>
            </div>
            <div className="flex-1 flex flex-col justify-between">
              <div className="space-y-6 px-5 pt-6 pb-5 flex flex-col">
                <div className="">
                  <label htmlFor="name">
                    <span className="input-label">{f(messages.type)}</span>
                  </label>
                  <div className="mt-1">
                    <Controller
                      name="sourceTypes"
                      control={control}
                      render={({ field }) => {
                        return (
                          <CheckBoxList
                            options={SOURCE_TYPE_ITEMS}
                            onChange={(values) => {
                              const convList = values.map((value) =>
                                Number(value)
                              );
                              field.onChange(convList.sort());
                            }}
                            disabled={!editable || isSubmitted}
                            defaultValues={defaultValues}
                          />
                        );
                      }}
                    />
                    <p className="input-error mt-2">
                      {errors.sourceTypes && (
                        <span role="alert">{errors.sourceTypes.message}</span>
                      )}
                    </p>
                  </div>
                </div>
              </div>
              <div className="flex-1 flex flex-col justify-between">
                <div className="space-y-6 px-5 pt-6 pb-5 flex flex-col">
                  <div className="">
                    <label htmlFor="webhookUrl">
                      <span className="input-label">
                        {f(messages.notification.slackIncomingWebhookUrl)}
                      </span>
                    </label>
                    <div className="mt-1">
                      <input
                        {...register('webhookUrl')}
                        type="text"
                        name="webhookUrl"
                        id="webhookUrl"
                        className="input-text w-full"
                        disabled={true}
                      />
                      <p className="input-error">
                        {errors.webhookUrl && (
                          <span role="alert">{errors.webhookUrl.message}</span>
                        )}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="flex-shrink-0 px-4 pt-4 pb-10 flex justify-end">
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
  });
