import { Feature } from '@/proto/feature/feature_pb';
import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { FC, memo, useCallback, useEffect } from 'react';
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
import { listGoals } from '../../modules/goals';
import { useCurrentEnvironment } from '../../modules/me';
import { OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import { OpsEventRateClause } from '../../proto/autoops/clause_pb';
import { ListGoalsRequest } from '../../proto/experiment/service_pb';
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
    const [feature, getFeatureError] = useSelector<
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
    const isLoading = isFeatureLoading || isAutoOpsRuleLoading;

    const defaultValues = {
      opsType: OpsType.ENABLE_FEATURE,
      clauseType: ClauseType.DATETIME,
      datetime: createInitialDatetimeClause(),
      eventRate: createInitialOpsEventRateClause(feature),
    };

    const methods = useForm({
      resolver: yupResolver(operationFormSchema),
      defaultValues,
      mode: 'onChange',
    });

    const handleRefetchAutoOpsRules = useCallback(() => {
      dispatch(
        listAutoOpsRules({
          featureId: featureId,
          environmentNamespace: currentEnvironment.id,
        })
      );
    }, [dispatch]);

    useEffect(() => {
      dispatch(
        listGoals({
          environmentNamespace: currentEnvironment.id,
          pageSize: 99999,
          cursor: '',
          searchKeyword: '',
          status: null,
          orderBy: ListGoalsRequest.OrderBy.NAME,
          orderDirection: ListGoalsRequest.OrderDirection.ASC,
        })
      );
    }, [dispatch, featureId, currentEnvironment]);

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
        />
      </FormProvider>
    );
  }
);

interface OpsEventRateClauseSchema {
  variation: string;
  goal: string;
  minCount: number;
  threadsholdRate: number;
  operator: string;
}

export function createOpsEventRateClause(
  oerc: OpsEventRateClauseSchema
): OpsEventRateClause {
  const clause = new OpsEventRateClause();
  clause.setVariationId(oerc.variation);
  clause.setGoalId(oerc.goal);
  clause.setMinCount(oerc.minCount);
  clause.setThreadsholdRate(oerc.threadsholdRate / 100);
  clause.setOperator(createOpsEventRateOperator(oerc.operator));
  return clause;
}

export function createOpsEventRateOperator(
  value: string
): OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap] {
  if (value === OpsEventRateClause.Operator.GREATER_OR_EQUAL.toString()) {
    return OpsEventRateClause.Operator.GREATER_OR_EQUAL;
  }
  return OpsEventRateClause.Operator.LESS_OR_EQUAL;
}
