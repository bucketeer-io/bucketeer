import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { getCurrentEnvironment, useAuth } from 'auth';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import { CollectionStatusType, Feature, FeatureCountByStatus } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import FilterFlagModal from './flags-modal/filter-flag-modal';
import Overview from './overview';
import SortBy from './sort-by';
import { FlagActionType, FlagFilters } from './types';

const PageContent = ({
  onAdd,
  onHandleActions
}: {
  onAdd: () => void;
  onHandleActions: (item: Feature, type: FlagActionType) => void;
}) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const isFiltered = useRef(false);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const [summary, setSummary] = useState<FeatureCountByStatus>();

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const searchFilters: Partial<FlagFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    status: 'ACTIVE',
    ...searchFilters
  } as FlagFilters;

  const [filters, setFilters] = usePartialState<FlagFilters>(defaultFilters);

  const filterCount = useMemo(() => {
    const { hasExperiment, hasPrerequisites, maintainer, enabled, tags } =
      filters || {};
    return isNotEmpty(
      enabled ?? hasExperiment ?? hasPrerequisites ?? maintainer ?? tags
    )
      ? 1
      : undefined;
  }, [filters]);

  const isHiddenTab = useMemo(
    () => !!filterCount || !!filters?.searchQuery,
    [filterCount, filters]
  );

  const onChangeFilters = (values: Partial<FlagFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
    isFiltered.current = true;
  };

  const onClearFilters = useCallback(() => {
    onChangeFilters({
      ...filters,
      searchQuery: '',
      hasExperiment: undefined,
      hasPrerequisites: undefined,
      maintainer: undefined,
      enabled: undefined,
      archived: undefined,
      tags: undefined,
      status: 'ACTIVE'
    });
    onCloseFilterModal();
  }, [filters]);

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  useEffect(() => {
    if (isFiltered.current) {
      onClearFilters();
      isFiltered.current = false;
    }
  }, [currentEnvironment]);

  return (
    <PageLayout.Content>
      <Overview summary={summary} onChangeFilters={() => {}} />
      <Filter
        action={
          <>
            <SortBy filters={filters} setFilters={setFilters} />
            <Button className="flex-1 lg:flex-none" onClick={onAdd}>
              <Icon icon={IconAddOutlined} size="sm" />
              {t(`create-flag`)}
            </Button>
          </>
        }
        filterCount={filterCount}
        searchValue={filters.searchQuery}
        onOpenFilter={onOpenFilterModal}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      {openFilterModal && (
        <FilterFlagModal
          filters={filters}
          isOpen={openFilterModal}
          onClearFilters={onClearFilters}
          onClose={onCloseFilterModal}
          onSubmit={v => {
            onChangeFilters(v);
            onCloseFilterModal();
          }}
        />
      )}
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        value={filters.status}
        onValueChange={value => {
          const status = value as CollectionStatusType;
          onChangeFilters({
            searchQuery: '',
            status,
            archived: status === 'ARCHIVED' || undefined
          });
        }}
      >
        {!isHiddenTab && (
          <TabsList>
            <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
            <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
          </TabsList>
        )}

        <TabsContent value={filters.status}>
          <CollectionLoader
            filters={filters}
            onAdd={onAdd}
            setFilters={setFilters}
            setSummary={setSummary}
            onHandleActions={onHandleActions}
            onClearFilters={onClearFilters}
            onChangeFilters={onChangeFilters}
          />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
