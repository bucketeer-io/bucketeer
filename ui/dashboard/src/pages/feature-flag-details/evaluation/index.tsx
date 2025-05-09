import React, { useState } from 'react';
import { useQueryEvaluation } from '@queries/evaluation';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { EvaluationTimeRange, Feature } from '@types';
import PageLayout from 'elements/page-layout';
import FilterBar from './filter-bar';

const EvaluationPage = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [timeRange, setTimeRange] = useState<EvaluationTimeRange>(
    EvaluationTimeRange.THIRTY_DAYS
  );

  const { data: evaluationCollection, isLoading } = useQueryEvaluation({
    params: {
      environmentId: currentEnvironment.id,
      featureId: feature.id,
      timeRange: EvaluationTimeRange.THIRTY_DAYS
    }
  });

  console.log({ evaluationCollection });

  return (
    <PageLayout.Content className="p-6 pt-0 gap-y-6">
      <FilterBar timeRange={timeRange} />
    </PageLayout.Content>
  );
};

export default EvaluationPage;
