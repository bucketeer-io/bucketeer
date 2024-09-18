import { listProgressiveRollout } from '../../modules/porgressiveRollout';
import { Feature } from '../../proto/feature/feature_pb';
import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useEffect } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useDispatch, useSelector, shallowEqual } from 'react-redux';
import { DetailSkeleton } from '../../components/DetailSkeleton';
import { FeatureAutoOpsRulesForm } from '../../components/FeatureAutoOpsRulesForm';
import { AppState } from '../../modules';
import { listAutoOpsRules } from '../../modules/autoOpsRules';
import { selectById as selectFeatureById } from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import { OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import { ProgressiveRollout } from '../../proto/autoops/progressive_rollout_pb';
import { AppDispatch } from '../../store';
import dayjs from 'dayjs';
import { operationFormSchema } from './formSchema';
import { v4 as uuid } from 'uuid';
import { ActionType, OpsEventRateClause } from '../../proto/autoops/clause_pb';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';

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
      state.features.getFeatureError
    ]);

    const defaultValues = {
      opsType: OpsType.SCHEDULE,
      datetimeClausesList: createInitialDatetimeClausesList(),
      eventRate: createInitialOpsEventRateClause(feature),
      progressiveRolloutType: ProgressiveRollout.Type.TEMPLATE_SCHEDULE,
      progressiveRollout: {
        template: {
          datetime: createInitialDatetimeClause(),
          interval: '1',
          increments: 10,
          variationId: feature.variationsList[0].id,
          schedulesList: []
        },
        manual: {
          variationId: feature.variationsList[0].id,
          schedulesList: [
            {
              executeAt: createInitialDatetimeClause(),
              weight: 10
            }
          ]
        }
      }
    };

    const methods = useForm({
      resolver: yupResolver(operationFormSchema),
      defaultValues,
      mode: 'onChange'
    });

    const { reset, setValue } = methods;

    useEffect(() => {
      fetchProgressiveRollouts();
      fetchAutoOpsRules();
    }, []);

    const fetchAutoOpsRules = useCallback(() => {
      dispatch(
        listAutoOpsRules({
          featureId: featureId,
          environmentId: currentEnvironment.id
        })
      );
    }, [dispatch]);

    const fetchProgressiveRollouts = useCallback(() => {
      dispatch(
        listProgressiveRollout({
          featureId: featureId,
          environmentId: currentEnvironment.id
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

    if (isFeatureLoading) {
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
          refetchAutoOpsRules={fetchAutoOpsRules}
          refetchProgressiveRollouts={fetchProgressiveRollouts}
          reset={handleReset}
        />
      </FormProvider>
    );
  }
);

const createInitialOpsEventRateClause = (feature: Feature.AsObject) => {
  return {
    variation: feature.variationsList[0].id,
    goal: null,
    minCount: 50,
    threadsholdRate: 50,
    operator: operatorOptions[0].value
  };
};

const createInitialDatetimeClause = () => {
  return {
    time: dayjs().add(1, 'hour').toDate()
  };
};

const createInitialDatetimeClausesList = () => {
  return [
    {
      id: uuid(),
      actionType: ActionType.ENABLE,
      time: dayjs().add(1, 'hour').toDate()
    }
  ];
};

export const operatorOptions = [
  {
    value: OpsEventRateClause.Operator.GREATER_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.greaterOrEqual)
  },
  {
    value: OpsEventRateClause.Operator.LESS_OR_EQUAL.toString(),
    label: intl.formatMessage(messages.feature.clause.operator.lessOrEqual)
  }
];
