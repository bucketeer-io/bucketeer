import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const DebuggerPage = () => {
  const { t } = useTranslation(['common', 'form', 'table']);

  return (
    <PageLayout.Root title={t('navigation.debugger')}>
      <PageHeader
        title={t('navigation.debugger')}
        description={t('debugger-subtitle')}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default DebuggerPage;
