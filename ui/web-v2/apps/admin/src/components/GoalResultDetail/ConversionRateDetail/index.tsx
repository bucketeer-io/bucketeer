import { FC } from 'react';

import { GoalResult } from '../../../proto/eventcounter/goal_result_pb';
import { Variation } from '../../../proto/feature/variation_pb';
import { ConversionRateDistributionChart } from '../ConversionRateDistributionChart';
import { ConversionRateDistributionTimeseriesChart } from '../ConversionRateDistributionTimeseriesChart';
import { ConversionRateTable } from '../ConversionRateTable';

interface ConversionRateDetailProps {
  goalResult: GoalResult.AsObject;
  baseVariationId: string;
  variations: Map<string, Variation.AsObject>;
}

export const ConversionRateDetail: FC<ConversionRateDetailProps> = ({
  goalResult,
  baseVariationId,
  variations,
}) => {
  return (
    <div>
      <ConversionRateTable
        goalResult={goalResult}
        baseVariationId={baseVariationId}
        variations={variations}
      />
      <div>
        {goalResult.variationResultsList[0].cvrMedianTimeseries ? (
          <div>
            <ConversionRateDistributionTimeseriesChart
              goalResult={goalResult}
              variations={variations}
            />
          </div>
        ) : (
          <div>
            <ConversionRateDistributionChart
              goalResult={goalResult}
              variations={variations}
            />
          </div>
        )}
      </div>
    </div>
  );
};
