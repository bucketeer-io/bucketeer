import { FC } from 'react';

import { GoalResult } from '../../../proto/eventcounter/goal_result_pb';
import { Variation } from '../../../proto/feature/variation_pb';
import { ValuePerUserDistributionTimeseriesChart } from '../ValuePerUserDistributionTimeseriesChart';
import { ValuePerUserTable } from '../ValuePerUserTable';

interface ValuePerUserDetailProps {
  goalResult: GoalResult.AsObject;
  baseVariationId: string;
  variations: Map<string, Variation.AsObject>;
}

export const ValuePerUserDetail: FC<ValuePerUserDetailProps> = ({
  goalResult,
  baseVariationId,
  variations,
}) => {
  return (
    <div>
      <ValuePerUserTable
        goalResult={goalResult}
        baseVariationId={baseVariationId}
        variations={variations}
      />
      <div>
        {goalResult.variationResultsList[0]
          .goalValueSumPerUserMedianTimeseries && (
          <div>
            <ValuePerUserDistributionTimeseriesChart
              goalResult={goalResult}
              variations={variations}
            />
          </div>
        )}
      </div>
    </div>
  );
};
