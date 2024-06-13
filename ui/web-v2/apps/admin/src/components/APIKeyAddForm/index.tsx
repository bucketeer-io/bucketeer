import { Dialog } from '@headlessui/react';
import React, { FC, memo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';

export interface APIKeyAddFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const APIKeyAddForm: FC<APIKeyAddFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      register,
      formState: { errors, isSubmitting, isValid },
      getValues,
    } = methods;

    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.apiKey.add.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.apiKey.add.header.description)}
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
                "
              >
                <div className="">
                  <label htmlFor="name">
                    <span className="input-label">{f(messages.name)}</span>
                    <span className="text-red-500">*</span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('name')}
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="pt-6">
                  <p className="font-bold text-lg text-gray-600">
                    {f(messages.apiKey.section.title)}
                  </p>
                  <div className="divide-y">
                    {[
                      {
                        id: 'client-sdk',
                        label: f(messages.apiKey.section.clientSdk),
                        description: f(
                          messages.apiKey.section.clientSdkDescription
                        ),
                        value: 1,
                      },
                      {
                        id: 'server-sdk',
                        label: f(messages.apiKey.section.serverSdk),
                        description: f(
                          messages.apiKey.section.serverSdkDescription
                        ),
                        value: 2,
                      },
                    ].map(({ id, label, description, value }) => (
                      <div
                        key={id}
                        className="flex items-center py-4 space-x-5"
                      >
                        <label htmlFor={id} className="flex-1 cursor-pointer">
                          <p className="text-base">{label}</p>
                          <p className="text-sm">{description}</p>
                        </label>
                        <input
                          {...register('role')}
                          id={id}
                          type="radio"
                          value={value}
                          className="h-4 w-4 text-primary focus:ring-primary border-gray-300 mt-1"
                          defaultChecked={getValues('role') === value}
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
            <button
              type="button"
              className="btn-submit"
              disabled={!isValid || isSubmitting}
              onClick={onSubmit}
            >
              {f(messages.button.submit)}
            </button>
          </div>
        </form>
      </div>
    );
  }
);
