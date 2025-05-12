import { useMemo } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { Evaluation, Feature } from '@types';
import ResultList from 'pages/debugger/debugger-results/result-list';
import { EvaluationFeature } from 'pages/debugger/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import FormLoading from 'elements/form-loading';

const TargetingDebuggerResults = ({
  isOpen,
  features,
  evaluations,
  onClose,
  onClearFields
}: {
  isOpen: boolean;
  features: Feature[];
  evaluations: Evaluation[];
  onClose: () => void;
  onClearFields: () => void;
}) => {
  const { t } = useTranslation(['common']);
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
  const evaluationData = useMemo(() => {
    const data = new Map();

    evaluations.forEach(item => {
      data.set(item.featureId, [
        ...(data.get(item.featureId) || []),
        {
          ...item,
          feature: features.find(feature => feature.id === item.featureId)
        }
      ]);
    });
    const results: EvaluationFeature[][] = [];
    data.forEach(evaluations => results.push(evaluations));
    return results;
  }, [features, evaluations]);

  return (
    <DialogModal
      isOpen={isOpen}
      title={t('navigation.debugger')}
      className="w-[940px] max-w-[940px]"
      onClose={onClose}
    >
      <div className="w-full max-h-[600px] overflow-auto py-3 [&>div>div]:!shadow-none">
        {accountLoading ? (
          <FormLoading />
        ) : (
          <ResultList
            className="min-w-full"
            accounts={accounts}
            groupBy={'FLAG'}
            groupByEvaluateFeatures={evaluationData}
            expandedItems={evaluationData.map((_, index) => index)}
          />
        )}
      </div>
      <ButtonBar
        primaryButton={
          <Button variant={'secondary'} className="w-fit" onClick={onClose}>
            {t('close')}
          </Button>
        }
        secondaryButton={
          <Button className="w-fit" onClick={onClearFields}>
            {t('clear-all-fields')}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default TargetingDebuggerResults;
