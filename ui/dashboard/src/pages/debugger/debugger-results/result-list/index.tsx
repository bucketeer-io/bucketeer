import { useCallback } from 'react';
import { Account } from '@types';
import { cn } from 'utils/style';
import { GroupByType } from 'pages/debugger/page-content';
import { EvaluationFeature } from 'pages/debugger/types';
import TableListContent from 'elements/table-list-content';
import ResultItem from './result-item';

const ResultList = ({
  groupBy,
  accounts,
  groupByEvaluateFeatures,
  expandedItems,
  className,
  onToggleExpandItem
}: {
  groupBy: GroupByType;
  accounts: Account[];
  groupByEvaluateFeatures: EvaluationFeature[][];
  expandedItems: number[];
  className?: string;
  onToggleExpandItem?: (index: number) => void;
}) => {
  const handleGetMaintainerInfo = useCallback(
    (email: string) => {
      const existedAccount = accounts?.find(account => account.email === email);
      if (
        !existedAccount ||
        !existedAccount?.firstName ||
        !existedAccount?.lastName
      )
        return email;
      return `${existedAccount.firstName} ${existedAccount.lastName}`;
    },
    [accounts]
  );

  return (
    <TableListContent className={cn('gap-y-2', className)}>
      {groupByEvaluateFeatures.map((group, index) => (
        <ResultItem
          groupBy={groupBy}
          key={index}
          group={group}
          isExpand={expandedItems.includes(index)}
          handleGetMaintainerInfo={handleGetMaintainerInfo}
          onToggleExpandItem={
            onToggleExpandItem ? () => onToggleExpandItem(index) : undefined
          }
        />
      ))}
    </TableListContent>
  );
};

export default ResultList;
