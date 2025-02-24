import { useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import { Experiment, ExperimentCollection } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
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
  ExperimentTab
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
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    status: 'ACTIVE',
    statuses: ['WAITING', 'RUNNING'],
    ...searchFilters
  } as ExperimentFilters;

  const [filters, setFilters] =
    usePartialState<ExperimentFilters>(defaultFilters);

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const onChangeFilters = (values: Partial<ExperimentFilters>) => {
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
      <Overview
        summary={summary}
        onChangeFilters={statuses =>
          onChangeFilters({
            statuses
          })
        }
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
        filterCount={isNotEmpty(filters.isFilter) ? 1 : undefined}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      {openFilterModal && (
        <FilterExperimentModal
          isOpen={openFilterModal}
          filters={filters}
          onClose={onCloseFilterModal}
          onSubmit={value => {
            onChangeFilters(value);
            onCloseFilterModal();
          }}
          onClearFilters={() => {
            onChangeFilters({
              archived: undefined,
              statuses: [],
              isFilter: undefined,
              status: 'ACTIVE'
            });
            onCloseFilterModal();
          }}
        />
      )}
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        value={filters.status}
        onValueChange={value => {
          const status = value as ExperimentTab;
          onChangeFilters({
            status,
            searchQuery: '',
            isFilter: undefined,
            statuses:
              status === 'FINISHED'
                ? ['STOPPED', 'FORCE_STOPPED']
                : status === 'ACTIVE'
                  ? ['WAITING', 'RUNNING']
                  : []
          });
        }}
      >
        <TabsList>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="FINISHED">{t(`finished`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.status as string}>
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
