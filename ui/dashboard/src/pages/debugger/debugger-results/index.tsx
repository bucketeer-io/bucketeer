import { useQueryAccounts } from '@queries/accounts';
import { getCurrentEnvironment, useAuth } from 'auth';
import FormLoading from 'elements/form-loading';
import TableListContainer from 'elements/table-list-container';
import { GroupByType } from '../page-content';
import { EvaluationFeature } from '../types';
import ActionBar from './action-bar';
import ResultList from './result-list';

interface Props {
  groupBy: GroupByType;
  isExpandAll: boolean;
  expandedItems: number[];
  groupByEvaluateFeatures: EvaluationFeature[][];
  onResetFields: () => void;
  onEditFields: () => void;
  onChangeGroupBy: (groupBy: GroupByType) => void;
  onToggleExpandItem: (index: number) => void;
  onToggleExpandAll: () => void;
}

const DebuggerResults = ({
  groupBy,
  isExpandAll,
  expandedItems,
  groupByEvaluateFeatures,
  onToggleExpandItem,
  onResetFields,
  onEditFields,
  onChangeGroupBy,
  onToggleExpandAll
}: Props) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: accountCollection, isLoading: accountLoading } =
    useQueryAccounts({
      params: {
        organizationId: currentEnvironment?.organizationId,
        cursor: String(0)
      }
    });

  const accounts = accountCollection?.accounts || [];

  return (
    <TableListContainer className="pt-0 gap-y-6">
      <ActionBar
        isExpandAll={isExpandAll}
        groupBy={groupBy}
        onChangeGroupBy={onChangeGroupBy}
        onResetFields={onResetFields}
        onEditFields={onEditFields}
        onToggleExpandAll={onToggleExpandAll}
      />
      {accountLoading ? (
        <FormLoading />
      ) : (
        <ResultList
          accounts={accounts}
          groupBy={groupBy}
          groupByEvaluateFeatures={groupByEvaluateFeatures}
          expandedItems={expandedItems}
          onToggleExpandItem={onToggleExpandItem}
        />
      )}
    </TableListContainer>
  );
};

export default DebuggerResults;
