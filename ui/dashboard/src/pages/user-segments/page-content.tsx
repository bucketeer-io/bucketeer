import { useCallback, useEffect, useMemo } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { getCurrentEnvironment, useAuth } from 'auth';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { UserSegment } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import CollectionLoader from './collection-loader';
import { UserSegmentsActionsType, UserSegmentsFilters } from './types';
import FilterUserSegmentModal from './user-segment-modal/filter-segment-modal';

const PageContent = ({
  editable,
  segmentUploading,
  onAdd,
  onActionHandler
}: {
  editable: boolean;
  segmentUploading: UserSegment | null;
  onAdd: () => void;
  onActionHandler: (
    segment: UserSegment,
    type: UserSegmentsActionsType
  ) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const organizationIds = useMemo(
    () => [currentEnvironment.organizationId],
    [currentEnvironment]
  );

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<UserSegmentsFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as UserSegmentsFilters;

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const [filters, setFilters] =
    usePartialState<UserSegmentsFilters>(defaultFilters);

  const onChangeFilters = useCallback(
    (values: Partial<UserSegmentsFilters>) => {
      values.page = values?.page || 1;
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  const onClearFilters = useCallback(() => {
    setFilters({ searchQuery: '', isInUseStatus: undefined });
  }, []);

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content>
      <Filter
        link={DOCUMENTATION_LINKS.SEGMENT}
        placeholder={t('form:name-desc-search-placeholder')}
        name="segment-list-search"
        onOpenFilter={onOpenFilterModal}
        action={
          <DisabledButtonTooltip
            align="end"
            hidden={editable}
            trigger={
              <Button
                className="flex-1 lg:flex-none"
                onClick={onAdd}
                disabled={!editable}
              >
                <Icon icon={IconAddOutlined} size="sm" />
                {t(`new-user-segment`)}
              </Button>
            }
          />
        }
        searchValue={filters.searchQuery as string}
        filterCount={
          isNotEmpty(filters.isInUseStatus as boolean) ? 1 : undefined
        }
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      {openFilterModal && (
        <FilterUserSegmentModal
          isOpen={openFilterModal}
          filters={filters}
          onClose={onCloseFilterModal}
          onSubmit={value => {
            onChangeFilters(value);
            onCloseFilterModal();
          }}
          onClearFilters={() => {
            onChangeFilters({ isInUseStatus: undefined });
            onCloseFilterModal();
          }}
        />
      )}
      <TableListContainer>
        <CollectionLoader
          segmentUploading={segmentUploading}
          organizationIds={organizationIds}
          filters={filters}
          onAdd={onAdd}
          setFilters={onChangeFilters}
          onActionHandler={onActionHandler}
          onClearFilters={onClearFilters}
        />
      </TableListContainer>
    </PageLayout.Content>
  );
};

export default PageContent;
