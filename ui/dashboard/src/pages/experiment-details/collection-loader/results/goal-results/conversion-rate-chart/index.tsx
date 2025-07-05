import { forwardRef, Ref, useMemo } from 'react';
import { GoalResult } from '@types';
import { getTimeSeries } from 'utils/chart';
import { GoalResultState } from '../..';
import { HistogramChart } from '../histogram-chart';
import {
  ChartToggleLegendRef,
  DatasetReduceType,
  TimeseriesAreaLineChart
} from '../timeseries-area-line-chart';
import { DataLabel } from '../timeseries-line-chart';

const ConversionRateChart = forwardRef(
  (
    {
      variationValues,
      goalResult,
      goalResultState,
      setConversionRateDataSets
    }: {
      variationValues: DataLabel[];
      goalResult: GoalResult;
      goalResultState: GoalResultState;
      setConversionRateDataSets: (datasets: DatasetReduceType[]) => void;
    },
    ref: Ref<ChartToggleLegendRef>
  ) => {
    const chartType = useMemo(
      () => goalResultState?.chartType,
      [goalResultState]
    );

    const isConversionRateChart = useMemo(
      () => chartType === 'conversion-rate',
      [chartType]
    );

    const timeseries = getTimeSeries(
      goalResult?.variationResults,
      goalResultState?.chartType,
      goalResultState?.tab
    );
    const upperBoundaries = useMemo(
      () =>
        goalResult?.variationResults?.map(item =>
          isConversionRateChart
            ? item?.cvrPercentile975Timeseries?.values.map(item => item * 100)
            : item?.goalValueSumPerUserPercentile025Timeseries?.values
        ) || [],
      [goalResult, isConversionRateChart]
    );
    const lowerBoundaries = useMemo(
      () =>
        goalResult?.variationResults?.map(item =>
          isConversionRateChart
            ? item?.cvrPercentile025Timeseries?.values.map(item => item * 100)
            : item?.goalValueSumPerUserPercentile025Timeseries?.values
        ) || [],
      [goalResult, isConversionRateChart]
    );

    const representatives = useMemo(
      () =>
        goalResult?.variationResults?.map(item =>
          isConversionRateChart
            ? item?.cvrMedianTimeseries?.values.map(item => item * 100)
            : item?.goalValueSumPerUserMedianTimeseries?.values
        ) || [],
      [goalResult, isConversionRateChart]
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

    bins = useMemo(() => bins?.map(b => Math.round(b * 10000) / 100), [bins]);

    return isConversionRateChart &&
      !goalResult.variationResults[0]?.cvrMedianTimeseries ? (
      <HistogramChart
        dataLabels={variationValues}
        bins={bins}
        hist={hist}
        label="Posterior Distribution"
      />
    ) : (
      <TimeseriesAreaLineChart
        ref={ref}
        chartType={chartType}
        dataLabels={variationValues}
        timeseries={timeseries}
        upperBoundaries={upperBoundaries}
        lowerBoundaries={lowerBoundaries}
        representatives={representatives}
        setDataSets={setConversionRateDataSets}
      />
    );
  }
);

export default ConversionRateChart;
