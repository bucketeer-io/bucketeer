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
        title="Organizations"
        description="You can see all your clients data"
      />
      <div className="py-8 px-6">
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
      </div>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default OrganizationPage;
