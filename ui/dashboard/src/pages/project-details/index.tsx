import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const ProjectDetailsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('projects')}>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default ProjectDetailsPage;
