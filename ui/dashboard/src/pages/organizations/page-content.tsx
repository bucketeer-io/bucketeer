import { IconAddOutlined } from 'react-icons-material-design';
import { organizationArchive, organizationUnArchive } from '@api/organization';
import { invalidateOrganizations } from '@queries/organizations';
import { useQueryClient } from '@tanstack/react-query';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import {
  CollectionStatusType,
  OrderBy,
  OrderDirection,
  Organization
} from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsList, TabsTrigger, TabsContent } from 'components/tabs';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import { OrganizationFilters } from './types';

const PageContent = ({
  onAdd,
  onEdit
}: {
  onAdd: () => void;
  onEdit: (v: Organization) => void;
}) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common']);

  // NOTE: Need improve search options
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const [filters, setFilters] = usePartialState<OrganizationFilters>({
    page: Number(searchOptions.page) || 1,
    orderBy: (searchOptions.orderBy as OrderBy) || 'DEFAULT',
    orderDirection: (searchOptions.orderDirection as OrderDirection) || 'ASC',
    searchQuery: (searchOptions.searchQuery as string) || '',
    status: (searchOptions.status as CollectionStatusType) || 'ACTIVE'
  });

  const onChangeFilters = (values: Partial<OrganizationFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  const onArchivedOrganization = (organization: Organization) => {
    organizationArchive({
      id: organization.id,
      command: {}
    }).then(() => {
      invalidateOrganizations(queryClient);
    });
  };

  const onUnArchiveOrganization = (organization: Organization) => {
    organizationUnArchive({
      id: organization.id,
      command: {}
    }).then(() => {
      invalidateOrganizations(queryClient);
    });
  };

  const onActionHandler = (type: string, organization: Organization) => {
    if (type === 'ARCHIVED_ORGANIZATION') {
      onArchivedOrganization(organization);
    } else if (type === 'UNARCHIVE_ORGANIZATION') {
      onUnArchiveOrganization(organization);
    } else {
      onEdit(organization);
    }
  };

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
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        defaultValue={filters.status}
        onValueChange={value =>
          onChangeFilters({ status: value as CollectionStatusType })
        }
      >
        <TabsList>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.status}>
          <CollectionLoader
            onAdd={onAdd}
            filters={filters}
            setFilters={onChangeFilters}
            onActionHandler={onActionHandler}
          />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
