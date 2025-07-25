import { useCallback, useEffect, useMemo } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { useAuthAccess } from 'auth';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import isNil from 'lodash/isNil';
import pickBy from 'lodash/pickBy';
import { APIKey } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import FilterAPIKeyModal from './api-key-modal/filter-api-key-modal';
import CollectionLoader from './collection-loader';
import { APIKeyActionsType, APIKeysFilters } from './types';

const PageContent = ({
  onAdd,
  onHandleActions
}: {
  onAdd: () => void;
  onHandleActions: (item: APIKey, type: APIKeyActionsType) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);
  const { envEditable, isOrganizationAdmin } = useAuthAccess();
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<APIKeysFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as APIKeysFilters;

  const [filters, setFilters] = usePartialState<APIKeysFilters>(defaultFilters);

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const filterCount = useMemo(() => {
    const filterKeys = ['disabled', 'environmentIds'];
    const count = filterKeys.reduce((acc, curr) => {
      if (!isNil(filters[curr as keyof APIKeysFilters])) ++acc;
      return acc;
    }, 0);
    return count || undefined;
  }, [filters]);

  const onChangeFilters = useCallback(
    (values: Partial<APIKeysFilters>) => {
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  const onClearFilters = useCallback(
    () =>
      onChangeFilters({
        searchQuery: '',
        disabled: undefined,
        environmentIds: undefined
      }),
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
        link={DOCUMENTATION_LINKS.API_KEYS}
        placeholder={t('form:name-desc-search-placeholder')}
        name="api-keys-list-search"
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
                {t(`new-api-key`)}
              </Button>
            }
          />
        }
        searchValue={filters.searchQuery}
        filterCount={filterCount}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      {openFilterModal && (
        <FilterAPIKeyModal
          isOpen={openFilterModal}
          filters={filters}
          onClose={onCloseFilterModal}
          onSubmit={value => {
            onChangeFilters(value);
            onCloseFilterModal();
          }}
          onClearFilters={() => {
            onChangeFilters({ disabled: undefined, environmentIds: undefined });
            onCloseFilterModal();
          }}
        />
      )}
      <TableListContainer>
        <CollectionLoader
          filters={filters}
          onAdd={onAdd}
          setFilters={onChangeFilters}
          onActions={onHandleActions}
          onClearFilters={onClearFilters}
        />
      </TableListContainer>
    </PageLayout.Content>
  );
};

export default PageContent;
