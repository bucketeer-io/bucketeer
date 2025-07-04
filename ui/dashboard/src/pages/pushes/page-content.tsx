import { useCallback, useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { Push } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import CollectionLoader from './collection-loader';
import FilterPushKeyModal from './push-modal/filter-push-modal';
import { PushActionsType, PushFilters } from './types';

const PageContent = ({
  disabled,
  onAdd,
  onHandleActions
}: {
  disabled?: boolean;
  onAdd: () => void;
  onHandleActions: (item: Push, type: PushActionsType) => void;
}) => {
  const { t } = useTranslation(['common']);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<PushFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as PushFilters;

  const [filters, setFilters] = usePartialState<PushFilters>(defaultFilters);

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const onChangeFilters = useCallback(
    (values: Partial<PushFilters>) => {
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
        link={DOCUMENTATION_LINKS.PUSHES}
        onOpenFilter={onOpenFilterModal}
        action={
          <DisabledButtonTooltip
            hidden={!disabled}
            trigger={
              <Button
                className="flex-1 lg:flex-none"
                onClick={onAdd}
                disabled={disabled}
              >
                <Icon icon={IconAddOutlined} size="sm" />
                {t(`new-push`)}
              </Button>
            }
          />
        }
        searchValue={filters.searchQuery}
        filterCount={isNotEmpty(filters.disabled) ? 1 : undefined}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      {openFilterModal && (
        <FilterPushKeyModal
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
      <TableListContainer>
        <CollectionLoader
          onAdd={onAdd}
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
