import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const CreateUpdateSegmentPage = () => {
  const { t } = useTranslation(['common']);
  return (
    <PageLayout.Root title={t('user-segments')}>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default CreateUpdateSegmentPage;
