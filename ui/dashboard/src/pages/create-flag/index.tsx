import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const CreateFlagPage = () => {
  const { t } = useTranslation(['common', 'table', 'form']);
  return (
    <PageLayout.Root title={t('navigation.feature-flags')}>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default CreateFlagPage;
