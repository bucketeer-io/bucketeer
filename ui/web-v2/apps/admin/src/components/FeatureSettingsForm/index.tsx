import { FC, memo } from 'react';
import { useFormContext, Controller } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { useIsEditable } from '../../modules/me';
import { classNames } from '../../utils/css';
import { CreatableSelect, Option } from '../CreatableSelect';

interface FeatureSettingsFormProps {
  onOpenConfirmDialog: () => void;
}

export const FeatureSettingsForm: FC<FeatureSettingsFormProps> = memo(
  ({ onOpenConfirmDialog }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors, isDirty },
    } = methods;
    const isValid = Object.keys(errors).length == 0;

    return (
      <div className="p-10 bg-gray-100">
        <form className="w-[600px]">
          <div
            className={classNames(
              'grid grid-cols-1 gap-y-6 gap-x-4 w-[600px]',
              'bg-white p-5 border border-gray-300 rounded'
            )}
          >
            <div className="col-span-1">
              <label htmlFor="name" className="input-label">
                {f(messages.name)}
              </label>
              <div className="">
                <input
                  {...register('name')}
                  type="text"
                  id="name"
                  autoComplete="name"
                  className="input-text w-full"
                  disabled={!editable}
                />
              </div>
              <p className="input-error">
                {errors.name && <span role="alert">{errors.name.message}</span>}
              </p>
            </div>
            <div className="col-span-1">
              <label
                htmlFor="about"
                className="block text-sm font-medium text-gray-700"
              >
                {f(messages.description)}
              </label>
              <div className="mt-1">
                <textarea
                  {...register('description')}
                  id="description"
                  rows={5}
                  className="input-textarea w-full"
                  disabled={!editable}
                />
                <p className="input-error">
                  {errors.description && (
                    <span role="alert">{errors.description.message}</span>
                  )}
                </p>
              </div>
            </div>
            <div className="col-span-1">
              <label htmlFor="tags" className="input-label">
                {f(messages.tags)}
              </label>
              <Controller
                name="tags"
                control={control}
                render={({ field }) => {
                  return (
                    <CreatableSelect
                      disabled={!editable}
                      defaultValues={field.value.map((tag) => {
                        return {
                          value: tag,
                          label: tag,
                        };
                      })}
                      onChange={(options: Option[]) => {
                        field.onChange(options.map((o) => o.value));
                      }}
                    />
                  );
                }}
              />
              <p className="input-error">
                {errors.tags && <span role="alert">{errors.tags.message}</span>}
              </p>
            </div>
          </div>
          {editable && (
            <div className="pt-10">
              <div className="flex justify-end">
                <button
                  type="button"
                  className="btn-submit"
                  disabled={!isDirty || !isValid}
                  onClick={onOpenConfirmDialog}
                >
                  {f(messages.button.saveWithComment)}
                </button>
              </div>
            </div>
          )}
        </form>
      </div>
    );
  }
);
