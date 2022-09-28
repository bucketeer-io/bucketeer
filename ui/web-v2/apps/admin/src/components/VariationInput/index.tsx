import { MinusCircleIcon } from '@heroicons/react/solid';
import { FC, memo, useCallback } from 'react';
import { Controller, useFieldArray, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { v4 as uuid } from 'uuid';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { useIsEditable } from '../../modules/me';
import { Feature } from '../../proto/feature/feature_pb';
import { Option, Select } from '../Select';

const variationTypeOptionsBoolean: Option = {
  value: Feature.VariationType.BOOLEAN.toString(),
  label: intl.formatMessage(messages.feature.type.boolean),
};
const variationTypeOptionsString: Option = {
  value: Feature.VariationType.STRING.toString(),
  label: intl.formatMessage(messages.feature.type.string),
};
const variationTypeOptionsNumber: Option = {
  value: Feature.VariationType.NUMBER.toString(),
  label: intl.formatMessage(messages.feature.type.number),
};
const variationTypeOptionsJson: Option = {
  value: Feature.VariationType.JSON.toString(),
  label: intl.formatMessage(messages.feature.type.json),
};

export const variationTypeOptions: Option[] = [
  variationTypeOptionsBoolean,
  variationTypeOptionsString,
  variationTypeOptionsNumber,
  variationTypeOptionsJson,
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

export interface VariationInputProps {
  removeDisabledIndexes: Set<number>;
  typeDisabled: boolean;
  onRemoveVariation: (idx: number) => void;
}

export const VariationInput: FC<VariationInputProps> = memo(
  ({ removeDisabledIndexes, typeDisabled, onRemoveVariation }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const methods = useFormContext();
    const {
      register,
      control,
      getValues,
      watch,
      reset,
      formState: { errors },
    } = methods;
    const {
      fields: variations,
      append,
      remove,
    } = useFieldArray({
      control,
      name: 'variations',
      keyName: 'key', // the default keyName is "id" and it conflicts with the variation id field
    });
    const variationType = watch('variationType');
    const disabledAddBtn =
      variationType == Feature.VariationType.BOOLEAN.toString();

    const handleChange = useCallback((type: string) => {
      reset(
        {
          ...getValues(),
          variationType: type,
          variations: [
            {
              id: uuid(),
              value:
                type === Feature.VariationType.BOOLEAN.toString() ? 'true' : '',
              name: '',
              description: '',
            },
            {
              id: uuid(),
              value:
                type === Feature.VariationType.BOOLEAN.toString()
                  ? 'false'
                  : '',
              name: '',
              description: '',
            },
          ],
        },
        { keepDirty: true }
      );
    }, []);

    const handleAddVariation = useCallback(() => {
      append({
        id: uuid(),
        value: '',
        name: '',
        description: '',
      });
    }, []);

    const handleRemoveVariation = useCallback((idx) => {
      remove(idx);
      onRemoveVariation(idx);
    }, []);

    return (
      <div className="space-y-4 flex flex-col">
        <div className="mb-1">
          <span className="input-label">
            {f(messages.feature.variationType)}
          </span>
          {typeDisabled ? (
            <span className="text-sm text-gray-600">{`: ${
              variationTypeOptions.find((o) => o.value === variationType).label
            }`}</span>
          ) : (
            <div>
              <Controller
                name="variationType"
                control={control}
                render={({ field }) => (
                  <Select
                    options={variationTypeOptions}
                    disabled={!editable}
                    value={variationTypeOptions.find(
                      (o) => o.value === variationType
                    )}
                    onChange={(option: Option) => {
                      handleChange(option.value);
                      field.onChange(option.value);
                    }}
                  />
                )}
              />
            </div>
          )}
        </div>
        <div className="space-y-5 flex flex-col">
          {variations.map((variation: any, idx) => {
            const disableRemoveBtn =
              variationType == Feature.VariationType.BOOLEAN.toString() ||
              variations.length <= 2 ||
              removeDisabledIndexes.has(idx);
            return (
              <div key={idx} className="flex flex-row flex-wrap mb-2">
                {getValues('variationType') ==
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
                {getValues('variationType') ==
                  Feature.VariationType.STRING.toString() ||
                getValues('variationType') ==
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
                        disabled={!editable}
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
                {getValues('variationType') ==
                Feature.VariationType.JSON.toString() ? (
                  <div className="w-full">
                    <label className="input-label" htmlFor="variation">
                      {`${f(messages.feature.variation)} ${idx + 1}`}
                    </label>
                    <div className="space-x-2 flex mt-1">
                      <textarea
                        {...register(`variations.${idx}.value`)}
                        className="input-text w-full"
                        disabled={!editable}
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
                    <span className="input-label-optional">
                      {' '}
                      {f(messages.input.optional)}
                    </span>
                  </label>
                  <div className="w-full mt-1">
                    <input
                      {...register(`variations[${idx}].name`)}
                      type="text"
                      className="w-full input-text"
                      disabled={!editable}
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
                    <span className="input-label">
                      {f(messages.description)}
                    </span>
                    <span className="input-label-optional">
                      {' '}
                      {f(messages.input.optional)}
                    </span>
                  </label>
                  <div className="w-full mt-1">
                    <input
                      {...register(`variations[${idx}].description`)}
                      type="text"
                      className="w-full input-text"
                      disabled={!editable}
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
                {editable && (
                  <div className="flex items-end py-3 ml-3">
                    <button
                      type="button"
                      onClick={() => handleRemoveVariation(idx)}
                      className="minus-circle-icon"
                      disabled={disableRemoveBtn}
                    >
                      <MinusCircleIcon aria-hidden="true" />
                    </button>
                  </div>
                )}
              </div>
            );
          })}
          {editable && (
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
          )}
        </div>
      </div>
    );
  }
);
