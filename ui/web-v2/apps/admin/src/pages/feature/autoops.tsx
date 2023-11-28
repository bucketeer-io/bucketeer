import { listProgressiveRollout } from '@/modules/porgressiveRollout';
import { Feature } from '@/proto/feature/feature_pb';
import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useDispatch, useSelector, shallowEqual } from 'react-redux';

import { DetailSkeleton } from '../../components/DetailSkeleton';
import {
  ClauseType,
  createInitialDatetimeClause,
  createInitialOpsEventRateClause,
  FeatureAutoOpsRulesForm,
} from '../../components/FeatureAutoOpsRulesForm';
import { AppState } from '../../modules';
import { listAutoOpsRules } from '../../modules/autoOpsRules';
import { selectById as selectFeatureById } from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import { OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import { AppDispatch } from '../../store';

import { operationFormSchema } from './formSchema';

interface FeatureAutoOpsPageProps {
  featureId: string;
}

export const FeatureAutoOpsPage: FC<FeatureAutoOpsPageProps> = memo(
  ({ featureId }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();
    const isFeatureLoading = useSelector<AppState, boolean>(
      (state) => state.features.loading,
      shallowEqual
    );
    const [feature] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ]);
    const isAutoOpsRuleLoading = useSelector<AppState, boolean>(
      (state) => state.autoOpsRules.loading,
      shallowEqual
    );
    const isProgressiveRolloutsLoading = useSelector<AppState, boolean>(
      (state) => state.progressiveRollout.loading,
      shallowEqual
    );

    const isLoading =
      isFeatureLoading || isAutoOpsRuleLoading || isProgressiveRolloutsLoading;

    const defaultValues = {
      opsType: OpsType.ENABLE_FEATURE,
      clauseType: ClauseType.DATETIME,
      datetime: createInitialDatetimeClause(),
      eventRate: createInitialOpsEventRateClause(feature),
      progressiveRollout: {
        template: {
          datetime: createInitialDatetimeClause(),
          interval: '1',
          increments: 10,
          variationId: feature.variationsList[0].id,
          schedulesList: [],
        },
        manual: {
          variationId: feature.variationsList[0].id,
          schedulesList: [
            {
              executeAt: createInitialDatetimeClause(),
              weight: 20,
            },
          ],
        },
      },
    };

    const methods = useForm({
      resolver: yupResolver(operationFormSchema),
      defaultValues,
      mode: 'onChange',
    });

    const { reset, setValue } = methods;

    const handleRefetchAutoOpsRules = useCallback(() => {
      dispatch(
        listAutoOpsRules({
          featureId: featureId,
          environmentNamespace: currentEnvironment.id,
        })
      );
    }, [dispatch]);

    const handleRefetchProgressiveRollouts = useCallback(() => {
      dispatch(
        listProgressiveRollout({
          featureId: featureId,
          environmentNamespace: currentEnvironment.id,
        })
      );
    }, [dispatch]);

    const handleReset = () => {
      reset(defaultValues);
      setValue(
        'progressiveRollout.manual.schedulesList',
        defaultValues.progressiveRollout.manual.schedulesList
      );
    };

    if (isLoading) {
      return (
        <div className="p-9 bg-gray-100">
          <DetailSkeleton />
        </div>
      );
    }

    return (
      <FormProvider {...methods}>
        <FeatureAutoOpsRulesForm
          featureId={featureId}
          refetchAutoOpsRules={handleRefetchAutoOpsRules}
          refetchProgressiveRollouts={handleRefetchProgressiveRollouts}
          reset={handleReset}
        />
      </FormProvider>
    );
  }
);
