import { IconAddOutlined } from 'react-icons-material-design';
import Filter from 'containers/filter';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import { CollectionStatusType } from '@types';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsList, TabsTrigger, TabsContent } from 'components/tabs/tabs';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import { OrganizationFilters } from './types';

const PageContent = ({ onAdd }: { onAdd: () => void }) => {
  const { t } = useTranslation(['common']);

  const [filters, setFilters] = usePartialState<OrganizationFilters>({
    searchQuery: '',
    status: 'ACTIVE'
  });

  return (
    <PageLayout.Content>
      <Filter
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-org`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        onSearchChange={v => setFilters({ ...filters, searchQuery: v })}
      />
      <Tabs
        defaultValue={filters.status}
        onValueChange={v =>
          setFilters({ ...filters, status: v as CollectionStatusType })
        }
        className="mt-6"
      >
        <TabsList>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>
        <TabsContent value={filters.status}>
          <CollectionLoader />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
