import { Dialog } from '@headlessui/react';
import { FC, memo } from 'react';
import { useFormContext, Controller } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { useIsEditable, useIsOwner } from '../../modules/me';
import { Account } from '../../proto/account/account_pb';
import { Option, Select } from '../Select';

export interface AccountUpdateFormProps {
  onSubmit: () => void;
  onCancel: () => void;
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

export const AccountUpdateForm: FC<AccountUpdateFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsOwner();
    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors, isSubmitting, isDirty, isValid },
    } = methods;

    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.account.update.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.account.update.header.description)}
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
                      disabled={!editable}
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
                            disabled={!editable}
                            value={roleOptions.find(
                              (o) => field.value == o.value
                            )}
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
