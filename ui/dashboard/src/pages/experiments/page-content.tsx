import { useCallback, useEffect, useMemo } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import { Experiment, ExperimentCollection, ExperimentStatus } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import FilterExperimentModal from './experiments-modal/filter-experiment-modal';
import Overview from './overview';
import {
  ExperimentActionsType,
  ExperimentFilters,
  ExperimentTab,
  SummaryType
} from './types';

const PageContent = ({
  summary,
  onAdd,
  onHandleActions
}: {
  summary?: ExperimentCollection['summary'];
  onAdd: () => void;
  onHandleActions: (item: Experiment, type: ExperimentActionsType) => void;
}) => {
  const { t } = useTranslation(['common']);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<ExperimentFilters> = searchOptions;

  const defaultFilters = {
    filterByTab: true,
    filterBySummary: undefined,
    page: 1,
    orderBy: 'NAME',
    orderDirection: 'ASC',
    status: 'ACTIVE',
    statuses: ['WAITING', 'RUNNING'],
    ...searchFilters
  } as ExperimentFilters;

  const [filters, setFilters] =
    usePartialState<ExperimentFilters>(defaultFilters);

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const isHiddenTab = useMemo(
    () =>
      (!!filters.searchQuery ||
        filters?.isFilter ||
        filters?.filterBySummary) &&
      !filters.filterByTab,
    [filters]
  );

  const onChangeFilters = useCallback(
    (values: Partial<ExperimentFilters>, isChangeParams = true) => {
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      if (isChangeParams) onChangSearchParams(options);
      setFilters({ ...values });
    },
    [filters]
  );

  const onClearFilters = useCallback(() => {
    onChangeFilters({
      archived: undefined,
      statuses: ['WAITING', 'RUNNING'],
      isFilter: undefined,
      status: 'ACTIVE',
      filterByTab: true,
      filterBySummary: undefined
    });
    onCloseFilterModal();
  }, []);

  const onChangeTab = useCallback(
    (status: ExperimentTab) => {
      onChangeFilters({
        status,
        searchQuery: filters?.searchQuery ?? '',
        isFilter: undefined,
        filterByTab: true,
        statuses:
          status === 'FINISHED'
            ? ['STOPPED', 'FORCE_STOPPED']
            : status === 'ACTIVE'
              ? ['WAITING', 'RUNNING']
              : []
      });
    },
    [filters]
  );

  const onFilterBySummary = useCallback(
    (statuses: ExperimentStatus[], summaryFilterValue: SummaryType) => {
      const isSameSummaryValue =
        filters?.filterBySummary === summaryFilterValue;
      onChangeFilters({
        statuses: isSameSummaryValue ? ['WAITING', 'RUNNING'] : statuses,
        filterBySummary: isSameSummaryValue ? undefined : summaryFilterValue,
        filterByTab: isSameSummaryValue,
        archived: undefined,
        status: isSameSummaryValue ? 'ACTIVE' : undefined
      });
    },
    [filters]
  );

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content>
      <Overview
        summary={summary}
        filterBySummary={filters?.filterBySummary}
        onChangeFilters={onFilterBySummary}
      />
      <Filter
        onOpenFilter={onOpenFilterModal}
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-experiment`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        filterCount={
          isNotEmpty(filters?.isFilter || filters?.filterBySummary)
            ? 1
            : undefined
        }
        onSearchChange={searchQuery => {
          onChangeFilters(
            {
              searchQuery,
              filterByTab:
                filters?.searchQuery === searchQuery || filters?.filterByTab
            },
            filters?.searchQuery !== searchQuery
          );
        }}
      />
      {openFilterModal && (
        <FilterExperimentModal
          isOpen={openFilterModal}
          filters={filters}
          onClose={onCloseFilterModal}
          onSubmit={value => {
            onChangeFilters({
              ...value,
              archived: undefined,
              status: undefined,
              filterByTab: false,
              filterBySummary: undefined
            });
            onCloseFilterModal();
          }}
          onClearFilters={onClearFilters}
        />
      )}
      <Tabs
        className={cn('flex-1 flex h-full flex-col', {
          'mt-6': !isHiddenTab
        })}
        value={filters.status}
        onValueChange={status => onChangeTab(status as ExperimentTab)}
      >
        <TabsList className={isHiddenTab ? 'hidden' : 'px-6'}>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="FINISHED">{t(`finished`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.status as string} className="px-6">
          <CollectionLoader
            onAdd={onAdd}
            filters={filters}
            setFilters={onChangeFilters}
            onActions={onHandleActions}
          />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
