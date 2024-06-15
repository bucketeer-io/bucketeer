import React, { FC, memo, useEffect, useState } from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import {
  FeatureEvaluation,
  TimeRange,
  timeRangeOptions,
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
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.evaluationTimeseriesCount.loading,
      shallowEqual
    );
    const [selectedTimeRange, setSelectedTimeRange] = useState(
      timeRangeOptions[0]
    );

    useEffect(() => {
      dispatch(
        getEvaluationTimeseriesCount({
          featureId: featureId,
          environmentNamespace: currentEnvironment.id,
          timeRange: TimeRange.LAST_THIRTY_DAYS,
        })
      );
    }, [dispatch, featureId, currentEnvironment]);

    if (isLoading) {
      return (
        <div className="flex pt-60 justify-center bg-gray-100">
          <div className="w-6 h-6 border-4 border-t-primary rounded-full animate-spin"></div>
        </div>
      );
    }

    return (
      <FeatureEvaluation
        featureId={featureId}
        selectedTimeRange={selectedTimeRange}
        setSelectedTimeRange={setSelectedTimeRange}
      />
    );
  }
);
