import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { DetailSkeleton } from '../../components/DetailSkeleton';
import { FeatureEvaluation } from '../../components/FeatureEvaluation';
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
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.evaluationTimeseriesCount.loading,
      shallowEqual
    );

    useEffect(() => {
      dispatch(
        getEvaluationTimeseriesCount({
          featureId: featureId,
          environmentNamespace: currentEnvironment.namespace,
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
    return <FeatureEvaluation featureId={featureId} />;
  }
);
