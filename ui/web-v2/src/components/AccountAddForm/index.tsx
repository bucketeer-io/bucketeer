import { Select, Option } from '../Select';
import { messages } from '../../lang/messages';
import { getRoleListV1 } from '../../pages/account';
import { Dialog } from '@headlessui/react';
import React, { FC, memo } from 'react';
import { useFormContext, Controller } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { CreatableSelect } from '../CreatableSelect';
import { shallowEqual, useSelector } from 'react-redux';
import { AppState } from '../../modules';
import { Tag } from '../../proto/tag/tag_pb';
import { selectAll as selectAllTags } from '../../modules/tags';

export interface AccountAddFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const AccountAddForm: FC<AccountAddFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors, isSubmitting, isValid }
    } = methods;

    const tagsList = useSelector<AppState, Tag.AsObject[]>(
      (state) => selectAllTags(state.tags),
      shallowEqual
    );

    const accountTagsList = tagsList.filter(
      (tag) => tag.entityType === Tag.EntityType.ACCOUNT
    );

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
                  <label htmlFor="name">
                    <span className="input-label">
                      {f(messages.input.name)}
                    </span>
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
                            options={getRoleListV1()}
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
                <div className="">
                  <label htmlFor="tags">
                    <span className="input-label">
                      {f(messages.tags.title)}
                    </span>
                  </label>
                  <Controller
                    name="tags"
                    control={control}
                    render={({ field }) => {
                      return (
                        <CreatableSelect
                          options={accountTagsList.map((tag) => ({
                            label: tag.name,
                            value: tag.name
                          }))}
                          onChange={(options: Option[]) => {
                            field.onChange(options.map((o) => o.value));
                          }}
                          closeMenuOnSelect={false}
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
