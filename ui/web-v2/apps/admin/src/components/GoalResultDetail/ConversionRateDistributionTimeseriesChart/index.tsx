import { unwrapUndefinable } from 'option-t/lib/Undefinable/unwrap';
import { FC } from 'react';

import { GoalResult } from '../../../proto/eventcounter/goal_result_pb';
import { Variation } from '../../../proto/feature/variation_pb';
import { TimeseriesAreaLineChart } from '../../TimeseriesAreaLineChart';

interface ConversionRateDistributionTimeseriesChartProps {
  goalResult: GoalResult.AsObject;
  variations: Map<string, Variation.AsObject>;
}

export const ConversionRateDistributionTimeseriesChart: FC<ConversionRateDistributionTimeseriesChartProps> =
  ({ goalResult, variations }) => {
    const variationValues = goalResult.variationResultsList.map((vr) => {
      return unwrapUndefinable(variations.get(vr.variationId)).value;
    });
    const timeseries = unwrapUndefinable(
      goalResult.variationResultsList[0].goalEventCountTimeseries
        ?.timestampsList
    );
    const upperPercentiles = goalResult.variationResultsList.map((vr) => {
      return unwrapUndefinable(vr.cvrPercentile975Timeseries).valuesList;
    });
    const lowerPercentiles = goalResult.variationResultsList.map((vr) => {
      return unwrapUndefinable(vr.cvrPercentile025Timeseries).valuesList;
    });
    const medians = goalResult.variationResultsList.map((vr) => {
      return unwrapUndefinable(vr.cvrMedianTimeseries).valuesList;
    });

    return (
      <TimeseriesAreaLineChart
        label={''}
        dataLabels={variationValues}
        timeseries={timeseries}
        upperBoundaries={upperPercentiles}
        lowerBoundaries={lowerPercentiles}
        representatives={medians}
        height={300}
      />
    );
  };
