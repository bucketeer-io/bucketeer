import { useCallback, useEffect, useMemo } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { useAuthAccess } from 'auth';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { Account } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
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
  const { envEditable, isOrganizationAdmin } = useAuthAccess();
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
    const filterKeys = ['disabled', 'organizationRole', 'teams'];
    const count = filterKeys.reduce((acc, curr) => {
      if (isNotEmpty(filters[curr as keyof MembersFilters])) ++acc;
      return acc;
    }, 0);
    return count || undefined;
  }, [filters]);

  const onChangeFilters = useCallback(
    (values: Partial<MembersFilters>) => {
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  const onAddMember = useCallback(() => {
    if (!envEditable || !isOrganizationAdmin) return undefined;
    return onAdd;
  }, [isOrganizationAdmin, envEditable]);

  const onClearFilters = useCallback(() => {
    onChangeFilters({
      ...filters,
      searchQuery: '',
      disabled: undefined,
      organizationRole: undefined,
      teams: undefined
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
        link={DOCUMENTATION_LINKS.MEMBERS}
        onOpenFilter={onOpenFilterModal}
        action={
          <DisabledButtonTooltip
            type={!isOrganizationAdmin ? 'admin' : 'editor'}
            hidden={envEditable && isOrganizationAdmin}
            trigger={
              <Button
                className="flex-1 lg:flex-none"
                onClick={onAdd}
                disabled={!envEditable || !isOrganizationAdmin}
              >
                <Icon icon={IconAddOutlined} size="sm" />
                {t(`invite-member`)}
              </Button>
            }
          />
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
          onAdd={onAddMember}
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
