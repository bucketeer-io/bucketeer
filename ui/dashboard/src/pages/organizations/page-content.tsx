import { useCallback, useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
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
import TableListContainer from 'elements/table-list-container';
import CollectionLoader from './collection-loader';
import FilterOrganizationModal from './organization-modal/filter-organization-modal';
import { OrganizationActionsType, OrganizationFilters } from './types';

const PageContent = ({
  onAdd,
  onHandleActions
}: {
  onAdd: () => void;
  onHandleActions: (item: Organization, type: OrganizationActionsType) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);

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

  const onChangeFilters = useCallback(
    (values: Partial<OrganizationFilters>) => {
      values.page = values?.page || 1;
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  const onClearFilters = useCallback(
    () => setFilters({ searchQuery: '', disabled: undefined }),
    []
  );

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content>
      <Filter
        isShowDocumentation={false}
        placeholder={t('form:name-search-placeholder')}
        name="organization-list-search"
        onOpenFilter={onOpenFilterModal}
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-org`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        filterCount={isNotEmpty(filters.disabled) ? 1 : undefined}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      {openFilterModal && (
        <FilterOrganizationModal
          isOpen={openFilterModal}
          filters={filters}
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
      )}
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        value={filters.status}
        onValueChange={value => {
          const status = value as CollectionStatusType;
          onChangeFilters({ status, searchQuery: '', disabled: undefined });
        }}
      >
        <TabsList className="px-3 sm:px-6">
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.status} className="w-full mt-0">
          <TableListContainer>
            <CollectionLoader
              filters={filters}
              onAdd={onAdd}
              setFilters={onChangeFilters}
              onActions={onHandleActions}
              onClearFilters={onClearFilters}
            />
          </TableListContainer>
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
