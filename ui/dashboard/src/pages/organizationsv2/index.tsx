import { IconAddOutlined } from 'react-icons-material-design';
import Filter from 'containers/filter';
import PageHeader from 'containers/page-header';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import Icon from 'components/icon';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const OrganizationPage = () => {
  const { t } = useTranslation(['common']);

  return (
    <PageLayout.Root title="Organizations">
      <PageHeader
        title={t('organizations')}
        description={t('organization-subtitle')}
      />
      <div className="p-6 flex flex-col flex-1">
        <Filter
          additionalActions={
            <Button className="flex-1 lg:flex-none" onClick={() => {}}>
              <Icon icon={IconAddOutlined} size="sm" />
              {t(`new-org`)}
            </Button>
          }
          searchValue={''}
          onChangeSearchValue={() => {}}
          onKeyDown={() => {}}
        />
        <PageLoader />
      </div>
    </PageLayout.Root>
  );
};

export default OrganizationPage;
