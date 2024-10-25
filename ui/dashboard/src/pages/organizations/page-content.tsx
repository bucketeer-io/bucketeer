import { useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { organizationArchive, organizationUnArchive } from '@api/organization';
import { invalidateOrganizations } from '@queries/organizations';
import { useQueryClient } from '@tanstack/react-query';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { CollectionStatusType, Organization } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsList, TabsTrigger, TabsContent } from 'components/tabs';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import FilterOrganizationModal from './organization-modal/filter-organization-modal';
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
  const searchFilters: Partial<OrganizationFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    status: 'ACTIVE',
    ...searchFilters
  } as OrganizationFilters;

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const [filters, setFilters] =
    usePartialState<OrganizationFilters>(defaultFilters);

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

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content>
      <Filter
        onOpenFilter={onOpenFilterModal}
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-org`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <FilterOrganizationModal
        isOpen={openFilterModal}
        onClose={onCloseFilterModal}
        onSubmit={value => {
          onChangeFilters(value);
          onCloseFilterModal();
        }}
        onClearFilters={() => {
          onChangeFilters({ disabled: undefined });
          onCloseFilterModal();
        }}
      />
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        value={filters.status}
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
