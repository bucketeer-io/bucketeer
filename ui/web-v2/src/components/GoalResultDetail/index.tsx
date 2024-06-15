import { SerializedError } from '@reduxjs/toolkit';
import { FC, memo, useState } from 'react';
import { useSelector, shallowEqual } from 'react-redux';

import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { selectById as selectExperimentResultById } from '../../modules/experimentResult';
import { selectById as selectExperimentById } from '../../modules/experiments';
import { ExperimentResult } from '../../proto/eventcounter/experiment_result_pb';
import { GoalResult } from '../../proto/eventcounter/goal_result_pb';
import { VariationResult } from '../../proto/eventcounter/variation_result_pb';
import { Experiment } from '../../proto/experiment/experiment_pb';
import { Variation } from '../../proto/feature/variation_pb';
import { classNames } from '../../utils/css';
import { Option, Select } from '../Select';

import { ConversionRateDetail } from './ConversionRateDetail';
import { ConversionRateTimeseriesChart } from './ConversionRateTimeseriesChart';
import { EvaluationUserTimeseriesChart } from './EvaluationUserTimeseriesChart';
import { GoalResultTable } from './GoalResultTable';
import { GoalTotalTimeseriesChart } from './GoalTotalTimeseriesChart';
import { GoalUserTimeseriesChart } from './GoalUserTimeseriesChart';
import { ValuePerUserDetail } from './ValuePerUserDetail';
import { ValuePerUserTimeseriesChart } from './ValuePerUserTimeseriesChart';
import { ValueTotalTimeseriesChart } from './ValueTotalTimeseriesChart';

const CHART_EVALUATION_USER = 'Evaluation user';
const CHART_GOAL_TOTAL = 'Goal total';
const CHART_GOAL_USER = 'Goal user';
const CHART_CONVERSION_RATE = 'Conversion rate';
const CHART_VALUE_TOTAL = 'Value total';
const CHART_VALUE_PER_USER = 'Value/User';

const chartOptions = [
  {
    value: CHART_EVALUATION_USER,
    label: intl.formatMessage(messages.experiment.result.evaluationUser.label),
  },
  {
    value: CHART_GOAL_TOTAL,
    label: intl.formatMessage(messages.experiment.result.goals.label),
  },
  {
    value: CHART_GOAL_USER,
    label: intl.formatMessage(messages.experiment.result.goalUser.label),
  },
  {
    value: CHART_CONVERSION_RATE,
    label: intl.formatMessage(messages.experiment.result.conversionRate.label),
  },
  {
    value: CHART_VALUE_TOTAL,
    label: intl.formatMessage(messages.experiment.result.valueSum.label),
  },
  {
    value: CHART_VALUE_PER_USER,
    label: intl.formatMessage(messages.experiment.result.valuePerUser.label),
  },
];

const ANALYSIS_CONVERSION_RATE = 'Conversion Rate Analysis';
const ANALYSIS_VALUE_PER_USER = 'Value Per User Analysis';
const analysisOptions = [
  {
    value: ANALYSIS_CONVERSION_RATE,
    label: intl.formatMessage(messages.experiment.result.conversionRate.label),
  },
  {
    value: ANALYSIS_VALUE_PER_USER,
    label: intl.formatMessage(messages.experiment.result.valuePerUser.label),
  },
];

interface GoalResultDetailProps {
  experimentId: string;
  goalId: string;
}

export const GoalResultDetail: FC<GoalResultDetailProps> = memo(
  ({ experimentId, goalId }) => {
    const [experiment, getExperimentError] = useSelector<
      AppState,
      [Experiment.AsObject | undefined, SerializedError | null]
    >(
      (state) => [
        selectExperimentById(state.experiments, experimentId),
        state.experiments.getExperimentError,
      ],
      shallowEqual
    );
    const [experimentResult, getExperimentResultError] = useSelector<
      AppState,
      [ExperimentResult.AsObject | undefined, SerializedError | null]
    >(
      (state) => [
        selectExperimentResultById(state.experimentResults, experimentId),
        state.experimentResults.getExperimentResultError,
      ],
      shallowEqual
    );
    const isExperimentLoading = useSelector<AppState, boolean>(
      (state) => state.experiments.loading,
      shallowEqual
    );
    const isExperimentResultLoading = useSelector<AppState, boolean>(
      (state) => state.experimentResults.loading,
      shallowEqual
    );
    const isLoading = isExperimentLoading || isExperimentResultLoading;
    const goalResult = goalId
      ? experimentResult.goalResultsList.find((gr) => gr.goalId === goalId)
      : experimentResult.goalResultsList[0];
    const variationMap = new Map<string, Variation.AsObject>();
    experiment.variationsList.forEach((v) => {
      variationMap.set(v.id, v);
    });

    // For old data before the Timeseries implementation,
    // we need to check the variation result contains the Timeseries data or not.
    // For this case we just need to check one.
    const containsTimeseriesData = (
      variationResult: VariationResult.AsObject[]
    ): boolean => {
      let containsTimeseries = true;
      variationResult.forEach((element) => {
        if (!element.evaluationUserCountTimeseries) {
          containsTimeseries = false;
          return;
        }
      });
      return containsTimeseries;
    };

    return isLoading || !goalResult ? null : (
      <div>
        <GoalResultTable goalResult={goalResult} variations={variationMap} />
        {containsTimeseriesData(goalResult.variationResultsList) && (
          <div>
            <div className="my-6 border-b border-gray-300" />
            <GoalResultDetailChart
              goalResult={goalResult}
              variationMap={variationMap}
            />
            <div className="my-6 border-b border-gray-300" />
            <GoalResultDetailAnalysis
              experiment={experiment}
              goalResult={goalResult}
              variationMap={variationMap}
            />
          </div>
        )}
      </div>
    );
  }
);

interface GoalResultDetailChartProps {
  goalResult: GoalResult.AsObject;
  variationMap: Map<string, Variation.AsObject>;
}

export const GoalResultDetailChart: FC<GoalResultDetailChartProps> = ({
  goalResult,
  variationMap,
}) => {
  const [chart, setChart] = useState<string>(CHART_EVALUATION_USER);

  return (
    <div>
      <Select
        options={chartOptions}
        className={classNames('text-sm w-[300px]')}
        value={chartOptions.find((c) => c.value === chart)}
        isSearchable={false}
        onChange={(e) => {
          setChart(e.value);
        }}
      />
      <div>
        {chart == CHART_EVALUATION_USER && (
          <EvaluationUserTimeseriesChart
            goalResult={goalResult}
            variations={variationMap}
          />
        )}
        {chart == CHART_GOAL_TOTAL && (
          <GoalTotalTimeseriesChart
            goalResult={goalResult}
            variations={variationMap}
          />
        )}
        {chart == CHART_GOAL_USER && (
          <GoalUserTimeseriesChart
            goalResult={goalResult}
            variations={variationMap}
          />
        )}
        {chart == CHART_CONVERSION_RATE && (
          <ConversionRateTimeseriesChart
            goalResult={goalResult}
            variations={variationMap}
          />
        )}
        {chart == CHART_VALUE_TOTAL && (
          <ValueTotalTimeseriesChart
            goalResult={goalResult}
            variations={variationMap}
          />
        )}
        {chart == CHART_VALUE_PER_USER && (
          <ValuePerUserTimeseriesChart
            goalResult={goalResult}
            variations={variationMap}
          />
        )}
      </div>
    </div>
  );
};

interface GoalResultDetailAnalysisProps {
  experiment: Experiment.AsObject;
  goalResult: GoalResult.AsObject;
  variationMap: Map<string, Variation.AsObject>;
}

export const GoalResultDetailAnalysis: FC<GoalResultDetailAnalysisProps> = ({
  experiment,
  goalResult,
  variationMap,
}) => {
  const [analysis, setAnalysis] = useState<string>(ANALYSIS_CONVERSION_RATE);

  return (
    <div>
      <Select
        options={analysisOptions}
        className={classNames('text-sm w-[300px] mb-6')}
        value={analysisOptions.find((a) => a.value === analysis)}
        isSearchable={false}
        onChange={(e) => {
          setAnalysis(e.value);
        }}
      />
      <div className="my-3">
        {analysis == ANALYSIS_CONVERSION_RATE && (
          <ConversionRateDetail
            goalResult={goalResult}
            baseVariationId={experiment.baseVariationId}
            variations={variationMap}
          />
        )}
        {analysis == ANALYSIS_VALUE_PER_USER && (
          <ValuePerUserDetail
            goalResult={goalResult}
            baseVariationId={experiment.baseVariationId}
            variations={variationMap}
          />
        )}
      </div>
    </div>
  );
};
