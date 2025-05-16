import { useMemo } from 'react';
import { OperationCombinedType, OpsTypeMap } from '../../types';
import EventRateProgress from './event-rate-progress';
import RolloutProgress from './rollout-progress';
import ScheduleProgress from './schedule-progress';

const OperationProgress = ({
  operation
}: {
  operation: OperationCombinedType;
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
  if (isEventRate) return <EventRateProgress operation={operation} />;
  return <RolloutProgress operation={operation} />;
};

export default OperationProgress;
