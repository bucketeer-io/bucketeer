import { AppState } from '../../modules';
import { Feature } from '../../proto/feature/feature_pb';
import { Dialog } from '@headlessui/react';
import { FC, memo, useCallback } from 'react';
import { Controller, useFieldArray, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';
import { v4 as uuid } from 'uuid';

import { messages } from '../../lang/messages';
import { selectAll as selectAllTags } from '../../modules/tags';
import { CreatableSelect, Option } from '../CreatableSelect';
import { Select } from '../Select';
import { intl } from '../../lang';
import { AddForm } from '../../pages/feature/formSchema';
import { HoverPopover } from '../HoverPopover';
import { classNames } from '../../utils/css';
import { MinusCircleIcon } from '@heroicons/react/solid';
import { Tag } from '../../proto/tag/tag_pb';

export interface FeatureAddFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const FeatureAddForm: FC<FeatureAddFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext<AddForm>();
    const {
      register,
      control,
      formState: { errors, isSubmitting, isDirty, isValid },
      getValues,
      setValue
    } = methods;

    const tagsList = useSelector<AppState, Tag.AsObject[]>(
      (state) => selectAllTags(state.tags),
      shallowEqual
    );

    const featureFlagTagsList = tagsList.filter(
      (tag) => tag.entityType === Tag.EntityType.FEATURE_FLAG
    );

    const variationsOptions = getValues('variations').map((variation, idx) => {
      return {
        id: variation.id,
        value: idx.toString(),
        label: `${f(messages.feature.variation)} ${idx + 1}`
      };
    });

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
                          options={featureFlagTagsList.map((tag) => ({
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
                <div className="">
                  <VariationAddInput />
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
                            menuPlacement="top"
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
                          menuPlacement="top"
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

const variationTypeOptionsBoolean: Option = {
  value: Feature.VariationType.BOOLEAN.toString(),
  label: intl.formatMessage(messages.feature.type.boolean)
};
const variationTypeOptionsString: Option = {
  value: Feature.VariationType.STRING.toString(),
  label: intl.formatMessage(messages.feature.type.string)
};
const variationTypeOptionsNumber: Option = {
  value: Feature.VariationType.NUMBER.toString(),
  label: intl.formatMessage(messages.feature.type.number)
};
const variationTypeOptionsJson: Option = {
  value: Feature.VariationType.JSON.toString(),
  label: intl.formatMessage(messages.feature.type.json)
};

export const variationTypeOptions: Option[] = [
  variationTypeOptionsBoolean,
  variationTypeOptionsString,
  variationTypeOptionsNumber,
  variationTypeOptionsJson
];

export const getVariationTypeOption = (
  type: Feature.VariationTypeMap[keyof Feature.VariationTypeMap]
): Option => {
  switch (type) {
    case Feature.VariationType.BOOLEAN:
      return variationTypeOptionsBoolean;
    case Feature.VariationType.STRING:
      return variationTypeOptionsString;
    case Feature.VariationType.NUMBER:
      return variationTypeOptionsNumber;
    default:
      return variationTypeOptionsJson;
  }
};

export interface VariationAddInputProps {}

export const VariationAddInput: FC<VariationAddInputProps> = memo(({}) => {
  const { formatMessage: f } = useIntl();
  const methods = useFormContext<AddForm>();
  const {
    register,
    control,
    getValues,
    watch,
    reset,
    formState: { errors }
  } = methods;
  const {
    fields: variations,
    append,
    remove
  } = useFieldArray({
    control,
    name: 'variations',
    keyName: 'key'
  });
  const variationType = watch('variationType');
  const disabledAddBtn =
    variationType === Feature.VariationType.BOOLEAN.toString();

  const { id: onVariationId } = getValues('onVariation');
  const { id: offVariationId } = getValues('offVariation');

  const handleChange = useCallback((type) => {
    const defaultVariationId1 = uuid();
    const defaultVariationId2 = uuid();

    reset(
      {
        ...getValues(),
        variationType: type,
        variations: [
          {
            id: defaultVariationId1,
            value:
              type === Feature.VariationType.BOOLEAN.toString() ? 'true' : '',
            name: '',
            description: ''
          },
          {
            id: defaultVariationId2,
            value:
              type === Feature.VariationType.BOOLEAN.toString() ? 'false' : '',
            name: '',
            description: ''
          }
        ],
        onVariation: {
          id: defaultVariationId1,
          value: '0',
          label: `${f(messages.feature.variation)} 1`
        },
        offVariation: {
          id: defaultVariationId2,
          value: '1',
          label: `${f(messages.feature.variation)} 2`
        }
      },
      { keepDirty: true }
    );
  }, []);

  const handleAddVariation = useCallback(() => {
    append({
      id: uuid(),
      value: '',
      name: '',
      description: ''
    });
  }, []);

  const handleRemoveVariation = useCallback((idx) => {
    remove(idx);
  }, []);

  const getVariationMessage = useCallback(
    (variationId) => {
      // Use a switch statement to determine the message
      switch (true) {
        // Check if the variation is both on and off variations
        case onVariationId === variationId && offVariationId === variationId:
          return f(messages.feature.variationSettings.bothVariations);

        // Check if the variation is the on variation
        case onVariationId === variationId:
          return f(messages.feature.variationSettings.defaultStrategy);

        // Check if the variation is the off variation
        case offVariationId === variationId:
          return f(messages.feature.variationSettings.offVariation);

        // Return null if none of the conditions are met
        default:
          return null;
      }
    },
    [onVariationId, offVariationId]
  );

  return (
    <div className="space-y-4 flex flex-col">
      <div className="mb-1">
        <span className="input-label">{f(messages.feature.variationType)}</span>
        <div>
          <Controller
            name="variationType"
            control={control}
            render={({ field }) => (
              <Select
                options={variationTypeOptions}
                disabled={false}
                value={variationTypeOptions.find(
                  (o) => o.value === variationType.toString()
                )}
                onChange={(option: Option) => {
                  handleChange(option.value);
                  field.onChange(option.value);
                }}
              />
            )}
          />
        </div>
      </div>
      <div className="space-y-5 flex flex-col">
        {variations.map((variation, idx) => {
          const disableRemoveBtn =
            variationType.toString() ==
              Feature.VariationType.BOOLEAN.toString() ||
            variation.id === onVariationId ||
            variation.id === offVariationId;
          return (
            <div key={idx} className="flex flex-row flex-wrap mb-2">
              {getValues('variationType') ===
              Feature.VariationType.BOOLEAN.toString() ? (
                <div>
                  <label className="input-label" htmlFor="variation">
                    {`${f(messages.feature.variation)} ${idx + 1}`}
                  </label>
                  <div className="mr-2 mt-1">
                    <input
                      {...register(`variations.${idx}.value`)}
                      type="text"
                      className="input-text"
                      disabled={true}
                    />
                  </div>
                  <p className="input-error">
                    {errors.variations?.[idx]?.value?.message && (
                      <span role="alert">
                        {errors.variations[idx].value.message}
                      </span>
                    )}
                  </p>
                </div>
              ) : null}
              {getValues('variationType') ===
                Feature.VariationType.STRING.toString() ||
              getValues('variationType') ===
                Feature.VariationType.NUMBER.toString() ? (
                <div>
                  <label className="input-label" htmlFor="variation">
                    {`${f(messages.feature.variation)} ${idx + 1}`}
                  </label>
                  <div className="mr-2 mt-1">
                    <input
                      {...register(`variations.${idx}.value`)}
                      type="text"
                      className="input-text"
                      disabled={false}
                    />
                  </div>
                  <p className="input-error">
                    {errors.variations?.[idx]?.value?.message && (
                      <span role="alert">
                        {errors.variations[idx].value.message}
                      </span>
                    )}
                  </p>
                </div>
              ) : null}
              {getValues('variationType') ===
              Feature.VariationType.JSON.toString() ? (
                <div className="w-full">
                  <label className="input-label" htmlFor="variation">
                    {`${f(messages.feature.variation)} ${idx + 1}`}
                  </label>
                  <div className="space-x-2 flex mt-1">
                    <textarea
                      {...register(`variations.${idx}.value`)}
                      className="input-text w-full"
                      disabled={false}
                      rows={10}
                    />
                  </div>
                  <p className="input-error">
                    {errors.variations?.[idx]?.value?.message && (
                      <span role="alert">
                        {errors.variations[idx].value.message}
                      </span>
                    )}
                  </p>
                </div>
              ) : null}
              <div className="mr-2 flex-grow">
                <label>
                  <span className="input-label">{f(messages.name)}</span>
                </label>
                <div className="w-full mt-1">
                  <input
                    {...register(`variations.${idx}.name`)}
                    type="text"
                    className="w-full input-text"
                    disabled={false}
                  />
                </div>
                <p className="input-error">
                  {errors.variations?.[idx]?.name?.message && (
                    <span role="alert">
                      {errors.variations[idx].name.message}
                    </span>
                  )}
                </p>
              </div>
              <div className="flex-grow">
                <label htmlFor="description">
                  <span className="input-label">{f(messages.description)}</span>
                  <span className="input-label-optional">
                    {' '}
                    {f(messages.input.optional)}
                  </span>
                </label>
                <div className="w-full mt-1">
                  <input
                    {...register(`variations.${idx}.description`)}
                    type="text"
                    className="w-full input-text"
                    disabled={false}
                  />
                </div>
                <p className="input-error">
                  {errors.variations?.[idx]?.description?.message && (
                    <span role="alert">
                      {errors.variations[idx].description.message}
                    </span>
                  )}
                </p>
              </div>
              {variations.length > 2 && (
                <div className="flex items-end py-3 ml-3">
                  {disableRemoveBtn ? (
                    <HoverPopover
                      render={() => {
                        const variationMessage = getVariationMessage(
                          variation.id
                        );
                        return variationMessage ? (
                          <div
                            className={classNames(
                              'bg-gray-900 text-white p-2 text-xs w-[350px]',
                              'rounded cursor-pointer'
                            )}
                          >
                            {variationMessage}
                          </div>
                        ) : null;
                      }}
                    >
                      <button
                        type="button"
                        className="minus-circle-icon"
                        disabled={disableRemoveBtn}
                      >
                        <MinusCircleIcon aria-hidden="true" />
                      </button>
                    </HoverPopover>
                  ) : (
                    <button
                      type="button"
                      onClick={() => handleRemoveVariation(idx)}
                      className="minus-circle-icon"
                    >
                      <MinusCircleIcon aria-hidden="true" />
                    </button>
                  )}
                </div>
              )}
            </div>
          );
        })}
        <div className="py-4 flex">
          <button
            type="button"
            className="btn-submit"
            onClick={handleAddVariation}
            disabled={disabledAddBtn}
          >
            {f(messages.button.addVariation)}
          </button>
        </div>
      </div>
    </div>
  );
});
