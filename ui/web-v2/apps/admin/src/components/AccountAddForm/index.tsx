import { Dialog } from '@headlessui/react';
import React, { FC, memo } from 'react';
import { useFormContext, Controller } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { Account } from '../../proto/account/account_pb';
import { Select } from '../Select';

export interface AccountAddFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export interface Option {
  value: string;
  label: string;
}

export const roleOptions: Option[] = [
  {
    value: Account.Role.VIEWER.toString(),
    label: intl.formatMessage(messages.account.role.viewer),
  },
  {
    value: Account.Role.EDITOR.toString(),
    label: intl.formatMessage(messages.account.role.editor),
  },
  {
    value: Account.Role.OWNER.toString(),
    label: intl.formatMessage(messages.account.role.owner),
  },
];

export const AccountAddForm: FC<AccountAddFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors, isSubmitting, isValid },
    } = methods;

    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.account.add.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.account.add.header.description)}
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
                  <label htmlFor="email">
                    <span className="input-label">
                      {f(messages.input.email)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('email')}
                      type="text"
                      name="email"
                      id="email"
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.email && (
                        <span role="alert">{errors.email.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="role" className="block">
                    <span className="input-label">
                      {f(messages.input.role)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <Controller
                      name="role"
                      control={control}
                      render={({ field }) => {
                        return (
                          <Select
                            onChange={(o: Option) => field.onChange(o.value)}
                            options={roleOptions}
                            isSearchable={false}
                          />
                        );
                      }}
                    />
                    <p className="input-error">
                      {errors.role && (
                        <span role="alert">{errors.role.message}</span>
                      )}
                    </p>
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
