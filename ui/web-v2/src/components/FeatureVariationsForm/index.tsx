import 'react-datepicker/dist/react-datepicker.css';
import {
  listProgressiveRollout,
  selectAll as selectAllProgressiveRollouts
} from '../../modules/porgressiveRollout';
import { ProgressiveRollout } from '../../proto/autoops/progressive_rollout_pb';
import { AppDispatch } from '../../store';
import { InformationCircleIcon } from '@heroicons/react/solid';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useEffect } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { Feature } from '../../proto/feature/feature_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
import { isProgressiveRolloutsRunningWaiting } from '../ProgressiveRolloutAddForm';
import { classNames } from '../../utils/css';
import { MinusCircleIcon } from '@heroicons/react/solid';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { v4 as uuid } from 'uuid';

import { intl } from '../../lang';
import { HoverPopover } from '../HoverPopover';
import { Option, Select } from '../Select';
import { VariationForm } from '../../pages/feature/formSchema';

export type ClauseType = 'compare' | 'segment' | 'date';
export interface ClauseTypeOption {
  value: ClauseType;
  label: string;
}

interface FeatureVariationsFormProps {
  featureId: string;
  onOpenConfirmDialog: () => void;
}

export const FeatureVariationsForm: FC<FeatureVariationsFormProps> = memo(
  ({ featureId, onOpenConfirmDialog }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext<VariationForm>();
    const {
      formState: { errors, isDirty },
      watch
    } = methods;
    const [feature] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError
    ]);

    const progressiveRollouts = useSelector<
      AppState,
      ProgressiveRollout.AsObject[]
    >(
      (state) =>
        selectAllProgressiveRollouts(state.progressiveRollout).filter(
          (rule) => rule.featureId === featureId
        ),
      shallowEqual
    );

    useEffect(() => {
      dispatch(
        listProgressiveRollout({
          featureId: featureId,
          environmentId: currentEnvironment.id
        })
      );
    }, []);

    const variations = watch('variations');

    const isValid = Object.keys(errors).length == 0;
    const isEmpty = variations.some((v) => !v.value || !v.name);

    const onVariationIds = [];
    if (feature.defaultStrategy.type === Strategy.Type.FIXED) {
      onVariationIds.push(feature.defaultStrategy.fixedStrategy.variation);
    } else if (feature.defaultStrategy.type === Strategy.Type.ROLLOUT) {
      feature.defaultStrategy.rolloutStrategy.variationsList.forEach((v) => {
        if (v.weight > 0) {
          onVariationIds.push(v.variation);
        }
      });
    }


    const ruleVariationIds = [];
    feature.rulesList.forEach((rule) => {
      if (rule.strategy.type === Strategy.Type.FIXED) {
        ruleVariationIds.push(rule.strategy.fixedStrategy.variation);
      } else if (rule.strategy.type === Strategy.Type.ROLLOUT) {
        rule.strategy.rolloutStrategy.variationsList.forEach((v) => {
          if (v.weight > 0) {
            ruleVariationIds.push(v.variation);
          }
        });
      }
    });

    const isProgressiveRolloutsRunning =
      progressiveRollouts.filter((p) =>
        isProgressiveRolloutsRunningWaiting(p.status)
      ).length > 0;

    return (
      <div className="p-10 bg-gray-100">
        {isProgressiveRolloutsRunning && (
          <div className="bg-blue-50 p-4 border-l-4 border-blue-400 mb-7 inline-block">
            <div className="flex">
              <div className="flex-shrink-0">
                <InformationCircleIcon
                  className="h-5 w-5 text-blue-400"
                  aria-hidden="true"
                />
              </div>
              <div className="ml-3 flex-1">
                <p className="text-sm text-blue-700">
                  {f(messages.feature.isProgressiveRolloutsRunning)}
                </p>
              </div>
            </div>
          </div>
        )}
        <form className="">
          <div className="bg-white border border-gray-300 rounded p-5 ">
            <VariationInput
              typeDisabled={true}
              isProgressiveRolloutsRunning={isProgressiveRolloutsRunning}
              rulesAppliedVariationList={{
                onVariationIds,
                offVariationId: feature.offVariation,
                ruleVariationIds
              }}
              featureId={featureId}
            />
          </div>
          <div className="flex justify-end mt-5">
            {editable && (
              <button
                type="button"
                className="btn-submit"
                disabled={!isDirty || !isValid || isEmpty}
                onClick={onOpenConfirmDialog}
              >
                {f(messages.button.saveWithComment)}
              </button>
            )}
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

type RulesAppliedVariationList = {
  onVariationId?: string;
  onVariationIds?: string[];
  offVariationId: string;
  ruleVariationIds?: string[];
};
export interface VariationInputProps {
  typeDisabled: boolean;
  isProgressiveRolloutsRunning?: boolean;
  rulesAppliedVariationList?: RulesAppliedVariationList;
  featureId: string;
}

export const VariationInput: FC<VariationInputProps> = memo(
  ({
    typeDisabled,
    rulesAppliedVariationList,
    isProgressiveRolloutsRunning,
    featureId
  }) => {
    const { formatMessage: f } = useIntl();
    const [feature] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError
    ]);
    const editable = useIsEditable();
    const methods = useFormContext<VariationForm>();
    const {
      register,
      control,
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
    const disabledAddBtn =
      feature.variationType == Feature.VariationType.BOOLEAN ||
      isProgressiveRolloutsRunning;

    const { onVariationId, onVariationIds, offVariationId, ruleVariationIds } =
      rulesAppliedVariationList;

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

          // Check if the variation is used in targeting rules
          case ruleVariationIds.includes(variationId):
            return f(messages.feature.variationSettings.targetingRule);

          // Return null if none of the conditions are met
          default:
            return null;
        }
      },
      [onVariationId, onVariationIds, offVariationId, ruleVariationIds]
    );

    return (
      <div className="space-y-4 flex flex-col">
        <div className="mb-1">
          <span className="input-label">
            {f(messages.feature.variationType)}
          </span>
          {typeDisabled ? (
            <span className="text-sm text-gray-600">{`: ${
              variationTypeOptions.find(
                (o) => o.value === feature.variationType.toString()
              ).label
            }`}</span>
          ) : (
            <div>
              <Select
                options={variationTypeOptions}
                disabled={!editable}
                value={variationTypeOptions.find(
                  (o) => o.value === feature.variationType.toString()
                )}
                onChange={() => {}}
              />
            </div>
          )}
        </div>
        <div className="space-y-5 flex flex-col">
          {variations.map((variation, idx) => {
            const disableRemoveBtn =
              feature.variationType.toString() ==
                Feature.VariationType.BOOLEAN.toString() ||
              variation.id === onVariationId ||
              variation.id === offVariationId ||
              (typeDisabled && onVariationIds.includes(variation.id)) ||
              ruleVariationIds.includes(variation.id);
            return (
              <div key={idx} className="flex flex-row flex-wrap mb-2">
                {feature.variationType === Feature.VariationType.BOOLEAN ? (
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
                {feature.variationType === Feature.VariationType.STRING ||
                feature.variationType === Feature.VariationType.NUMBER ? (
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
                {feature.variationType === Feature.VariationType.JSON ? (
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
                      {...register(`variations.${idx}.name`)}
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
                      {...register(`variations.${idx}.description`)}
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
