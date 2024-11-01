import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const ProjectsPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('projects')}>
      <PageHeader title={t('projects')} description={t('project-subtitle')} />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default ProjectsPage;
