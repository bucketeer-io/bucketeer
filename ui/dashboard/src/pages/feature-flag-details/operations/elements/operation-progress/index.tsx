import { useMemo } from 'react';
import { AutoOpsCount, RuleStrategyVariation } from '@types';
import { OperationCombinedType, OpsTypeMap } from '../../types';
import EventRateProgress from './event-rate-progress';
import RolloutProgress from './rollout-progress';
import ScheduleProgress from './schedule-progress';

const OperationProgress = ({
  operation,
  opsCounts,
  currentAllocationPercentage
}: {
  currentAllocationPercentage: RuleStrategyVariation[];
  operation: OperationCombinedType;
  opsCounts: AutoOpsCount[];
}) => {
  const isSchedule = useMemo(
    () => operation.opsType === OpsTypeMap.SCHEDULE,
    [operation]
  );
  const isEventRate = useMemo(
    () => operation.opsType === OpsTypeMap.EVENT_RATE,
    [operation]
  );
  if (isSchedule) return <ScheduleProgress operation={operation} />;
  if (isEventRate)
    return <EventRateProgress operation={operation} opsCounts={opsCounts} />;
  return (
    <RolloutProgress
      currentAllocationPercentage={currentAllocationPercentage}
      operation={operation}
    />
  );
};

export default OperationProgress;
