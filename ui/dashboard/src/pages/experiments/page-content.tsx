import { useCallback, useEffect, useMemo } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import pickBy from 'lodash/pickBy';
import { Experiment, ExperimentCollection, ExperimentStatus } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
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
  disabled,
  summary,
  onAdd,
  onHandleActions
}: {
  disabled: boolean;
  summary?: ExperimentCollection['summary'];
  onAdd: () => void;
  onHandleActions: (item: Experiment, type: ExperimentActionsType) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);

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

  const filterCount = useMemo(() => {
    if (isNotEmpty(filters?.isFilter || filters?.filterBySummary)) {
      const filterKeys = ['statuses', 'maintainer'];
      const count = filterKeys.reduce((acc, curr) => {
        if (isNotEmpty(filters[curr as keyof ExperimentFilters])) ++acc;
        return acc;
      }, 0);
      return count || undefined;
    }
    return undefined;
  }, [filters]);

  const isHiddenTab = useMemo(() => {
    return (
      (!!filters.searchQuery || filters.isFilter || filters.filterBySummary) &&
      !filters.filterByTab
    );
  }, [filters]);

  const onChangeFilters = useCallback(
    (values: Partial<ExperimentFilters>, isChangeParams = true) => {
      values.page = values?.page || 1;
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));

      if (isChangeParams) {
        onChangSearchParams(options);
      }
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
      filterBySummary: undefined,
      maintainer: undefined
    });
    onCloseFilterModal();
  }, []);

  const onChangeTab = useCallback(
    (status: ExperimentTab) => {
      onChangeFilters({
        status,
        searchQuery: filters.searchQuery ?? '',
        isFilter: undefined,
        maintainer: undefined,
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
      const isSameSummaryValue = filters.filterBySummary === summaryFilterValue;
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
        link={DOCUMENTATION_LINKS.EXPERIMENTS}
        placeholder={t('form:name-desc-search-placeholder')}
        name="experiments-list-search"
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
                {t(`new-experiment`)}
              </Button>
            }
          />
        }
        searchValue={filters.searchQuery}
        filterCount={filterCount}
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
        <TabsList className={isHiddenTab ? 'hidden' : 'px-3 sm:px-6'}>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="FINISHED">{t(`finished`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={(filters.status as string) || ''} className="mt-0">
          <TableListContainer>
            <CollectionLoader
              onAdd={onAdd}
              filters={filters}
              setFilters={onChangeFilters}
              onActions={onHandleActions}
            />
          </TableListContainer>
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
