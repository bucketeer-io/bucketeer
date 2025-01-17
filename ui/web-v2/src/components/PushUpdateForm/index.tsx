import { AppState } from '../../modules';
import { Tag } from '../../proto/tag/tag_pb';
import { Dialog } from '@headlessui/react';
import { FC, memo } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { useIsEditable } from '../../modules/me';
import { selectAll as selectAllTags } from '../../modules/tags';
import { Option } from '../CreatableSelect';
import { Select } from '../Select';

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
      formState: { errors, isValid, isDirty, isSubmitted }
    } = methods;

    const tagsList = useSelector<AppState, Tag.AsObject[]>(
      (state) => selectAllTags(state.tags),
      shallowEqual
    );
    const featureFlagTagsList = tagsList.filter(
      (tag) => tag.entityType === Tag.EntityType.FEATURE_FLAG
    );

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
                        <Select
                          isMulti
                          onChange={(options: Option[]) => {
                            field.onChange(options.map((o) => o.value));
                          }}
                          disabled={!editable || isSubmitted}
                          value={field.value?.map((tag: string) => {
                            return {
                              value: tag,
                              label: tag
                            };
                          })}
                          options={featureFlagTagsList.map((tag) => ({
                            label: tag.name,
                            value: tag.name
                          }))}
                          closeMenuOnSelect={false}
                          placeholder={f(messages.tags.tagsPlaceholder)}
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
