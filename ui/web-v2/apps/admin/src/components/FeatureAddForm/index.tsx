import { AppState } from '@/modules';
import { Tag } from '@/proto/feature/feature_pb';
import { Dialog } from '@headlessui/react';
import { FC, memo } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { selectAll as selectAllTags } from '../../modules/tags';
import { CreatableSelect, Option } from '../CreatableSelect';
import { Select } from '../Select';
import { VariationInput } from '../VariationInput';

export interface FeatureAddFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const FeatureAddForm: FC<FeatureAddFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors, isSubmitting, isDirty, isValid },
      getValues,
      setValue,
      watch,
    } = methods;
    const tagsList = useSelector<AppState, Tag.AsObject[]>(
      (state) => selectAllTags(state.tags),
      shallowEqual
    );

    const variationsOptions = getValues('variations').map((variation, idx) => {
      return {
        id: variation.id,
        value: idx.toString(),
        label: `${f(messages.feature.variation)} ${idx + 1}`,
      };
    });

    const onVariation = watch('onVariation');
    const offVariation = watch('offVariation');

    return (
      <div className="w-[600px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.feature.add.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.feature.add.header.description)}
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
                  <label htmlFor="id">
                    <span className="input-label">
                      {f(messages.feature.id)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('id')}
                      type="text"
                      className="input-text w-full"
                      placeholder={''}
                    />
                    <p className="input-error">
                      {errors.id && (
                        <span role="alert">{errors.id.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="name">
                    <span className="input-label">{f({ id: 'name' })}</span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('name')}
                      type="text"
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
                  <label htmlFor="description" className="block">
                    <span className="input-label">
                      {f(messages.description)}
                    </span>
                    <span className="input-label-optional">
                      {' '}
                      {f(messages.input.optional)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <textarea
                      {...register('description')}
                      rows={4}
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.description && (
                        <span role="alert">{errors.description.message}</span>
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
                          options={tagsList.map((tag) => ({
                            label: tag.id,
                            value: tag.id,
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
                <div className="">
                  <VariationInput
                    typeDisabled={false}
                    rulesAppliedVariationList={{
                      onVariationId: onVariation.id,
                      offVariationId: offVariation.id,
                    }}
                  />
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label htmlFor="onVariation" className="input-label">
                      {f(messages.feature.onVariation)}
                    </label>
                    <Controller
                      name="onVariation"
                      control={control}
                      render={({ field }) => {
                        return (
                          <Select
                            onChange={(variation) => {
                              field.onChange(variation);
                              setValue('onVariation', variation);
                            }}
                            options={variationsOptions}
                            className="w-full"
                            value={variationsOptions.find(
                              (v) => v.id === field.value.id
                            )}
                            isSearchable={false}
                          />
                        );
                      }}
                    />
                  </div>
                  <div>
                    <label htmlFor="offVariation" className="input-label">
                      {f(messages.feature.offVariation)}
                    </label>
                    <Controller
                      name="offVariation"
                      control={control}
                      render={({ field }) => (
                        <Select
                          onChange={(variation) => {
                            field.onChange(variation);
                            setValue('offVariation', variation);
                          }}
                          options={variationsOptions}
                          className="w-full"
                          value={variationsOptions.find(
                            (v) => v.id === field.value.id
                          )}
                          isSearchable={false}
                        />
                      )}
                    />
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
              disabled={!isDirty || !isValid || isSubmitting}
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
