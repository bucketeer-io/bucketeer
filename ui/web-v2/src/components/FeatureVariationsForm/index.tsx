import 'react-datepicker/dist/react-datepicker.css';
import {
  listProgressiveRollout,
  selectAll as selectAllProgressiveRollouts,
} from '../../modules/porgressiveRollout';
import { ProgressiveRollout } from '../../proto/autoops/progressive_rollout_pb';
import { AppDispatch } from '../../store';
import { InformationCircleIcon } from '@heroicons/react/solid';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useEffect } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { Feature } from '../../proto/feature/feature_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
import { isProgressiveRolloutsRunningWaiting } from '../AddProgressiveRolloutOperation';
import { VariationInput } from '../VariationInput';

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
    const methods = useFormContext();
    const {
      control,
      formState: { errors, isDirty },
      watch,
    } = methods;
    const [feature, _] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
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
          environmentNamespace: currentEnvironment.id,
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
              }}
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
