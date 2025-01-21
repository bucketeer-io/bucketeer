import { useState } from 'react';
import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const ExperimentsPage = () => {
  const { t } = useTranslation(['common']);
  const [totalCount, setTotalCount] = useState<string | number>(0);
  return (
    <PageLayout.Root title={t('navigation.experiments')}>
      <PageHeader
        title={`${t('navigation.experiments')} (${totalCount})`}
        description={t('experiments-subtitle')}
      />
      <PageLoader setTotalCount={setTotalCount} />
    </PageLayout.Root>
  );
};

export default ExperimentsPage;
