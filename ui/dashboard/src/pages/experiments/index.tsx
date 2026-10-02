import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const ExperimentsPage = () => {
  const { t } = useTranslation(['common']);
  return (
    <PageLayout.Root title={t('navigation.experiments')}>
      <PageHeader
        title={`${t('navigation.experiments')}`}
        description={t('experiments-subtitle')}
        link={DOCUMENTATION_LINKS.EXPERIMENTS}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default ExperimentsPage;
