import { useCallback } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { EvaluationFeatureAccount } from 'pages/debugger/types';
import { GroupByType } from '..';
import ResultItem from './result-item';

const ResultList = ({
  groupBy,
  evaluationFeatures
}: {
  groupBy: GroupByType;
  evaluationFeatures: EvaluationFeatureAccount[];
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: accountCollection } = useQueryAccounts({
    params: {
      organizationId: currentEnvironment?.organizationId,
      cursor: String(0)
    }
  });

  const accounts = accountCollection?.accounts || [];

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
    <div className="flex flex-col w-full gap-y-2">
      {evaluationFeatures.map((item, index) => (
        <ResultItem
          groupBy={groupBy}
          key={index}
          featureId={item.id}
          feature={item.feature}
          maintainer={item.feature.maintainer}
          userId={item.user_id}
          isExpand={false}
          handleGetMaintainerInfo={handleGetMaintainerInfo}
        />
      ))}
    </div>
  );
};

export default ResultList;
