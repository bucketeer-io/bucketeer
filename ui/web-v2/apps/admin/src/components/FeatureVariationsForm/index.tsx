import 'react-datepicker/dist/react-datepicker.css';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { useIsEditable } from '../../modules/me';
import { Feature } from '../../proto/feature/feature_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
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
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      control,
      formState: { errors, isSubmitting, isDirty },
    } = methods;
    const [feature, _] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);
    const isValid = Object.keys(errors).length == 0;

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

    return (
      <div className="p-10 bg-gray-100">
        <form className="">
          <div className="bg-white border border-gray-300 rounded p-5 ">
            <VariationInput
              typeDisabled={true}
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
                disabled={!isDirty || !isValid}
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
