import { unwrapUndefinable } from 'option-t/lib/Undefinable/unwrap';
import { FC } from 'react';

import { GoalResult } from '../../../proto/eventcounter/goal_result_pb';
import { Variation } from '../../../proto/feature/variation_pb';
import { HistogramChart } from '../../HistogramChart';

interface ConversionRateDistributionChartProps {
  goalResult: GoalResult.AsObject;
  variations: Map<string, Variation.AsObject>;
}

export const ConversionRateDistributionChart: FC<ConversionRateDistributionChartProps> =
  ({ goalResult, variations }) => {
    const variationValues = goalResult.variationResultsList.map((vr) => {
      return unwrapUndefinable(variations.get(vr.variationId)).value;
    });
    let bins = Array<number>();
    const hist = goalResult.variationResultsList.map((vr) => {
      const cvrProb = vr.cvrProb;
      if (!cvrProb) {
        return [];
      }
      const histogram = unwrapUndefinable(cvrProb.histogram);
      if (bins.length == 0) {
        bins = histogram.binsList;
      }
      return histogram.histList;
    });
    bins = bins.map((b) => {
      return Math.round(b * 10000) / 100;
    });

    return (
      <HistogramChart
        label={'Posterior Distribution'}
        dataLabels={variationValues}
        hist={hist}
        bins={bins}
      />
    );
  };
