import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useTranslation } from 'i18n';
import PageHeader from 'elements/page-header';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const MembersPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title={t('members')}>
      <PageHeader
        title={t('members')}
        description={t('member-subtitle')}
        link={DOCUMENTATION_LINKS.MEMBERS}
      />
      <PageLoader />
    </PageLayout.Root>
  );
};

export default MembersPage;
