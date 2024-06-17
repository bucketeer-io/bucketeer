import { Dialog } from '@headlessui/react';
import React, { FC, memo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { useIsOwner } from '../../modules/me';

export interface APIKeyUpdateFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const APIKeyUpdateForm: FC<APIKeyUpdateFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsOwner();
    const methods = useFormContext();
    const {
      register,
      formState: { errors, isSubmitting, isDirty, isValid },
      getValues
    } = methods;

    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.apiKey.update.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.apiKey.update.header.description)}
                </p>
              </div>
            </div>
            <div
              className="
                flex-1
                flex flex-col
                justify-between
              "
            >
              <div
                className="
                  space-y-6 px-5 pt-6 pb-5
                  flex flex-col
                "
              >
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
                      disabled={!editable}
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="pt-6">
                  <p className="font-bold text-md text-gray-600">
                    {f(messages.apiKey.section.title)}
                  </p>
                  <div className="divide-y opacity-50">
                    {[
                      {
                        id: 'client-sdk',
                        label: f(messages.apiKey.section.clientSdk),
                        description: f(
                          messages.apiKey.section.clientSdkDescription
                        ),
                        value: 1
                      },
                      {
                        id: 'server-sdk',
                        label: f(messages.apiKey.section.serverSdk),
                        description: f(
                          messages.apiKey.section.serverSdkDescription
                        ),
                        value: 2
                      }
                    ].map(({ id, label, description, value }) => (
                      <div
                        key={id}
                        className="flex items-center py-4 space-x-5"
                      >
                        <label htmlFor={id} className="flex-1">
                          <p className="font-bold text-sm text-gray-600 text-md text-gray-500">
                            {label}
                          </p>
                          <p className="text-sm text-gray-500">{description}</p>
                        </label>
                        <input
                          id={id}
                          type="radio"
                          className="h-4 w-4 text-primary focus:ring-primary border-gray-300 mt-1"
                          checked={getValues('role') === value}
                          disabled={true}
                        />
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="flex-shrink-0 px-4 py-4 flex justify-end">
            <div className="mr-3">
              <button
                type="button"
                className="btn-cancel"
                disabled={false}
                onClick={onCancel}
              >
                {f(messages.button.cancel)}
              </button>
            </div>
            {editable && (
              <button
                type="button"
                className="btn-submit"
                disabled={!isDirty || !isValid || isSubmitting}
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
