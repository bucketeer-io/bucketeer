import { useMemo } from 'react';
import { GoalResult } from '@types';
import { getTimeSeries } from 'utils/chart';
import { GoalResultState } from '../..';
import { HistogramChart } from '../histogram-chart';
import { TimeseriesAreaLineChart } from '../timeseries-area-line-chart';

const ConversionRateChart = ({
  variationValues,
  goalResult,
  goalResultState
}: {
  variationValues: string[];
  goalResult: GoalResult;
  goalResultState: GoalResultState;
}) => {
  const chartType = useMemo(
    () => goalResultState?.chartType,
    [goalResultState]
  );

  const timeseries = getTimeSeries(
    goalResult?.variationResults,
    goalResultState?.chartType,
    goalResultState?.tab
  );

  const upperBoundaries = useMemo(
    () =>
      goalResult?.variationResults?.map(item =>
        chartType === 'conversion-rate'
          ? item.cvrPercentile975Timeseries.values.map(item => item * 100)
          : item.goalValueSumPerUserPercentile025Timeseries?.values
      ) || [],
    [goalResult, chartType]
  );
  const lowerBoundaries = useMemo(
    () =>
      goalResult?.variationResults?.map(item =>
        chartType === 'conversion-rate'
          ? item.cvrPercentile025Timeseries.values.map(item => item * 100)
          : item.goalValueSumPerUserPercentile025Timeseries?.values
      ) || [],
    [goalResult, chartType]
  );

  const representatives = useMemo(
    () =>
      goalResult?.variationResults?.map(item =>
        chartType === 'conversion-rate'
          ? item.cvrMedianTimeseries.values.map(item => item * 100)
          : item.goalValueSumPerUserMedianTimeseries?.values
      ) || [],
    [goalResult, chartType]
  );

  let bins: number[] = [];

  const hist = useMemo(
    () =>
      goalResult.variationResults.map(vr => {
        const cvrProb = vr?.cvrProb;
        if (!cvrProb) {
          return [];
        }
        const histogram = cvrProb?.histogram;
        bins = !bins.length ? histogram?.bins || [] : bins;

        return histogram?.hist || [];
      }),
    [goalResult, bins]
  );

  bins = bins?.map(b => Math.round(b * 10000) / 100);

  return chartType === 'conversion-rate' &&
    !goalResult.variationResults[0]?.cvrMedianTimeseries ? (
    <HistogramChart
      dataLabels={variationValues}
      bins={bins}
      hist={hist}
      label="Posterior Distribution"
    />
  ) : (
    <TimeseriesAreaLineChart
      dataLabels={variationValues}
      timeseries={timeseries}
      upperBoundaries={upperBoundaries}
      lowerBoundaries={lowerBoundaries}
      representatives={representatives}
    />
  );
};

export default ConversionRateChart;
