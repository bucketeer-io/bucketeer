import { useCallback, useMemo, useState } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryFeatures } from '@queries/features';
import { getCurrentEnvironment, useAuth } from 'auth';
import { Evaluation } from '@types';
import FormLoading from 'elements/form-loading';
import TableListContainer from 'elements/table-list-container';
import { EvaluationFeature } from '../types';
import ActionBar from './action-bar';
import ResultList from './result-list';

export type GroupByType = 'FLAG' | 'USER';

interface Props {
  evaluations: Evaluation[];
  onResetFields: () => void;
  onEditFields: () => void;
}

const DebuggerResults = ({
  evaluations,
  onResetFields,
  onEditFields
}: Props) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [groupBy, setGroupBy] = useState<GroupByType>('FLAG');
  const [expandedItems, setExpandedItems] = useState<number[]>([]);

  const { data: accountCollection, isLoading: accountLoading } =
    useQueryAccounts({
      params: {
        organizationId: currentEnvironment?.organizationId,
        cursor: String(0)
      }
    });

  const accounts = accountCollection?.accounts || [];

  const { data: featureCollection, isLoading: featureLoading } =
    useQueryFeatures({
      params: {
        cursor: String(0),
        environmentId: currentEnvironment.id
      }
    });

  const features = featureCollection?.features || [];

  const groupByEvaluateFeatures: EvaluationFeature[][] = useMemo(() => {
    const isFlag = groupBy === 'FLAG';
    const data = new Map();

    evaluations.forEach(item => {
      const groupByField = isFlag ? item.featureId : item.userId;
      data.set(groupByField, [
        ...(data.get(groupByField) || []),
        {
          ...item,
          feature: features.find(feature => feature.id === item.featureId)
        }
      ]);
    });
    const results: EvaluationFeature[][] = [];
    data.forEach(evaluations => results.push(evaluations));
    return results;
  }, [groupBy, features, evaluations]);

  const onToggleCollapseItem = useCallback(
    (index: number) => {
      const isExistedItem = expandedItems.includes(index);
      setExpandedItems(
        isExistedItem
          ? expandedItems.filter(item => item !== index)
          : [...expandedItems, index]
      );
    },
    [expandedItems]
  );

  return (
    <TableListContainer className="mt-0 gap-y-6">
      <ActionBar
        groupBy={groupBy}
        setGroupBy={setGroupBy}
        onResetFields={onResetFields}
        onEditFields={onEditFields}
        onResetExpandItems={() => setExpandedItems([])}
      />
      {accountLoading || featureLoading ? (
        <FormLoading />
      ) : (
        <ResultList
          accounts={accounts}
          groupBy={groupBy}
          groupByEvaluateFeatures={groupByEvaluateFeatures}
          expandedItems={expandedItems}
          onToggleCollapseItem={onToggleCollapseItem}
        />
      )}
    </TableListContainer>
  );
};

export default DebuggerResults;
