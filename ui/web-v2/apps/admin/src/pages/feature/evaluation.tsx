import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { DetailSkeleton } from '../../components/DetailSkeleton';
import {
  FeatureEvaluation,
  TimeRange,
} from '../../components/FeatureEvaluation';
import { AppState } from '../../modules';
import { getEvaluationTimeseriesCount } from '../../modules/evaluationTimeseriesCount';
import { useCurrentEnvironment } from '../../modules/me';
import { AppDispatch } from '../../store';

interface FeatureEvaluationPageProps {
  featureId: string;
}

export const FeatureEvaluationPage: FC<FeatureEvaluationPageProps> = memo(
  ({ featureId }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    useEffect(() => {
      dispatch(
        getEvaluationTimeseriesCount({
          featureId: featureId,
          environmentNamespace: currentEnvironment.namespace,
          timeRange: TimeRange.LAST_THIRTY_DAYS,
        })
      );
    }, [dispatch, featureId, currentEnvironment]);

    return <FeatureEvaluation featureId={featureId} />;
  }
);
