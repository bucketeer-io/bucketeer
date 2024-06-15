import { unwrapUndefinable } from 'option-t/lib/Undefinable/unwrap';
import { FC } from 'react';

import { GoalResult } from '../../../proto/eventcounter/goal_result_pb';
import { Variation } from '../../../proto/feature/variation_pb';
import { TimeseriesLineChart } from '../../TimeseriesLineChart';

interface GoalUserTimeseriesChartProps {
  goalResult: GoalResult.AsObject;
  variations: Map<string, Variation.AsObject>;
}

export const GoalUserTimeseriesChart: FC<GoalUserTimeseriesChartProps> = ({
  goalResult,
  variations,
}) => {
  const variationValues = goalResult.variationResultsList.map((vr) => {
    return unwrapUndefinable(variations.get(vr.variationId)).value;
  });
  const timeseries = unwrapUndefinable(
    goalResult.variationResultsList[0].goalUserCountTimeseries?.timestampsList
  );
  const data = goalResult.variationResultsList.map((vr) => {
    return unwrapUndefinable(vr.goalUserCountTimeseries).valuesList;
  });

  return (
    <TimeseriesLineChart
      label={''}
      dataLabels={variationValues}
      timeseries={timeseries}
      data={data}
      height={300}
    />
  );
};
