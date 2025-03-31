import { useMemo } from 'react';
import { AutoOpsRule, Rollout } from '@types';

const ActiveContent = ({
  rollouts,
  operations
}: {
  rollouts: Rollout[];
  operations: AutoOpsRule[];
}) => {
  const activeStatuses = useMemo(() => ['WAITING', 'RUNNING'], []);

  const operationActiveData = useMemo(
    () =>
      operations?.filter(item => activeStatuses.includes(item.autoOpsStatus)),
    [operations, activeStatuses]
  );

  const rolloutActiveData = useMemo(
    () => rollouts?.filter(item => activeStatuses.includes(item.status)),
    [rollouts, activeStatuses]
  );

  console.log({ operationActiveData });
  console.log({ rolloutActiveData });
  return <div>ActiveContent</div>;
};

export default ActiveContent;
