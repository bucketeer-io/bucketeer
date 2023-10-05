import { classNames } from '@/utils/css';
import { MinusCircleIcon } from '@heroicons/react/solid';
import { FC, memo, useCallback } from 'react';
import { Controller, useFieldArray, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { v4 as uuid } from 'uuid';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { useIsEditable } from '../../modules/me';
import { Feature } from '../../proto/feature/feature_pb';
import { HoverPopover } from '../HoverPopover';
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

type RulesAppliedVariationList = {
  onVariationId?: string;
  onVariationIds?: string[];
  offVariationId: string;
};
export interface VariationInputProps {
  typeDisabled: boolean;
  rulesAppliedVariationList?: RulesAppliedVariationList;
}

export const VariationInput: FC<VariationInputProps> = memo(
  ({ typeDisabled, rulesAppliedVariationList }) => {
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
      keyName: 'key',
      // keyName: 'key', // the default keyName is "id" and it conflicts with the variation id field
    });
    const variationType = watch('variationType');
    const disabledAddBtn =
      variationType == Feature.VariationType.BOOLEAN.toString();
    const { onVariationId, onVariationIds, offVariationId } =
      rulesAppliedVariationList;

    const handleChange = useCallback((type: string) => {
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
              description: '',
            },
            {
              id: defaultVariationId2,
              value:
                type === Feature.VariationType.BOOLEAN.toString()
                  ? 'false'
                  : '',
              name: '',
              description: '',
            },
          ],
          onVariation: {
            id: defaultVariationId1,
            value: 0,
            label: `${f(messages.feature.variation)} 1`,
          },
          offVariation: {
            id: defaultVariationId2,
            value: 1,
            label: `${f(messages.feature.variation)} 2`,
          },
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
    }, []);

    const getVariationMessage = useCallback(
      (variationId) => {
        // Use a switch statement to determine the message
        switch (true) {
          // Check if the variation is both on and off variations
          case typeDisabled &&
            onVariationIds.includes(variationId) &&
            offVariationId === variationId:
            return f(messages.feature.variationSettings.bothVariations);

          // Check if the variation is the on variation
          case typeDisabled && onVariationIds.includes(variationId):
            return f(messages.feature.variationSettings.defaultStrategy);

          // Check if the variation is the off variation
          case typeDisabled && variationId === offVariationId:
            return f(messages.feature.variationSettings.offVariation);

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
      [onVariationId, onVariationIds, offVariationId]
    );

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
              variation.id === onVariationId ||
              variation.id === offVariationId ||
              (typeDisabled && onVariationIds.includes(variation.id));
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
                {editable && variations.length > 2 && (
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
