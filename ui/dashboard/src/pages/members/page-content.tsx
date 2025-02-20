import { useEffect } from 'react';
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

  const onChangeFilters = (values: Partial<MembersFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
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
          isOrganizationAdmin && (
            <Button className="flex-1 lg:flex-none" onClick={onAdd}>
              <Icon icon={IconAddOutlined} size="sm" />
              {t(`invite-member`)}
            </Button>
          )
        }
        searchValue={filters.searchQuery}
        filterCount={
          isNotEmpty(filters.disabled || filters.organizationRole)
            ? 1
            : undefined
        }
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
          onClearFilters={() => {
            onChangeFilters({
              disabled: undefined,
              organizationRole: undefined
            });
            onCloseFilterModal();
          }}
        />
      )}
      <div className="mt-5 flex flex-col flex-1">
        <CollectionLoader
          onAdd={() => {
            if (isOrganizationAdmin) onAdd();
          }}
          filters={filters}
          setFilters={onChangeFilters}
          onActions={onHandleActions}
        />
      </div>
    </PageLayout.Content>
  );
};

export default PageContent;
