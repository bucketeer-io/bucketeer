import { useCallback, useEffect, useMemo } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { useAuth } from 'auth';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { Account } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import CollectionLoader from './collection-loader';
import FilterMemberModal from './member-modal/filter-member-modal';
import { MemberActionsType, MembersFilters } from './types';

const PageContent = ({
  onAdd,
  onHandleActions
}: {
  onAdd: () => void;
  onHandleActions: (item: Account, type: MemberActionsType) => void;
}) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const isOrganizationAdmin =
    consoleAccount?.organizationRole === 'Organization_ADMIN';

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<MembersFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as MembersFilters;

  const [filters, setFilters] = usePartialState<MembersFilters>(defaultFilters);

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const filterCount = useMemo(() => {
    const { disabled, organizationRole, tags } = filters;
    return isNotEmpty(disabled || organizationRole || tags) ? 1 : undefined;
  }, [filters]);

  const onChangeFilters = (values: Partial<MembersFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  const onClearFilters = useCallback(() => {
    onChangeFilters({
      ...filters,
      searchQuery: '',
      disabled: undefined,
      organizationRole: undefined,
      tags: undefined
    });
    onCloseFilterModal();
  }, [filters]);

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
          isOrganizationAdmin && (
            <Button className="flex-1 lg:flex-none" onClick={onAdd}>
              <Icon icon={IconAddOutlined} size="sm" />
              {t(`invite-member`)}
            </Button>
          )
        }
        searchValue={filters.searchQuery}
        filterCount={filterCount}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      {openFilterModal && (
        <FilterMemberModal
          isOpen={openFilterModal}
          filters={filters}
          onClose={onCloseFilterModal}
          onSubmit={value => {
            onChangeFilters(value);
            onCloseFilterModal();
          }}
          onClearFilters={onClearFilters}
        />
      )}
      <TableListContainer>
        <CollectionLoader
          onAdd={() => {
            if (isOrganizationAdmin) onAdd();
          }}
          filters={filters}
          setFilters={onChangeFilters}
          onActions={onHandleActions}
          onClearFilters={onClearFilters}
        />
      </TableListContainer>
    </PageLayout.Content>
  );
};

export default PageContent;
