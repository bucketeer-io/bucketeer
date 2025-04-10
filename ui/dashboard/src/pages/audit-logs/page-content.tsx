import { useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuth } from 'auth';
import { usePartialState } from 'hooks';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { IconCollapse, IconExpand } from '@icons';
import Button from 'components/button';
import { ReactDateRangePicker } from 'components/date-range-picker';
import Icon from 'components/icon';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import EntityTypeDropdown from './elements/entity-type-dropdown';
import { AuditLogsFilters, ExpandOrCollapse } from './types';

const PageContent = () => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<AuditLogsFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'TIMESTAMP',
    orderDirection: 'DESC',
    from: Math.trunc(new Date().getTime() / 1000),
    to: Math.trunc(new Date().getTime() / 1000),
    ...searchFilters
  } as AuditLogsFilters;

  const [filters, setFilters] =
    usePartialState<AuditLogsFilters>(defaultFilters);

  const [expandOrCollapseAllState, setExpandOrCollapseAllState] = useState<
    ExpandOrCollapse | undefined
  >(undefined);

  const [expandedItems, setExpandedItems] = useState<string[]>([]);

  const isExpandAll = useMemo(
    () => expandOrCollapseAllState === ExpandOrCollapse.EXPAND,
    [expandOrCollapseAllState]
  );

  const onChangeFilters = (values: Partial<AuditLogsFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  const handleExpandOrCollapseAll = useCallback(() => {
    setExpandOrCollapseAllState(
      isExpandAll ? ExpandOrCollapse.COLLAPSE : ExpandOrCollapse.EXPAND
    );
    setExpandedItems([]);
  }, [isExpandAll]);

  const onToggleExpandItem = useCallback(
    (id: string) => {
      const isExistedItem = expandedItems.find(item => item === id);
      setExpandedItems(
        isExistedItem
          ? expandedItems.filter(item => item !== id)
          : [...expandedItems, id]
      );
    },
    [expandedItems]
  );

  return (
    <PageLayout.Content className="gap-y-6">
      <Filter
        action={
          <>
            <EntityTypeDropdown
              isSystemAdmin={!!consoleAccount?.isSystemAdmin}
              entityType={filters?.entityType}
              onChangeFilters={onChangeFilters}
            />
            <ReactDateRangePicker
              from={filters?.from}
              to={filters?.to}
              onChange={(startDate, endDate) =>
                onChangeFilters({
                  from: startDate ? startDate?.toString() : undefined,
                  to: endDate ? endDate?.toString() : undefined
                })
              }
            />
            <Button
              variant={'secondary'}
              onClick={handleExpandOrCollapseAll}
              className="max-w-[154px]"
            >
              <Icon
                icon={isExpandAll ? IconCollapse : IconExpand}
                size="sm"
                color="primary-500"
              />
              {t(isExpandAll ? 'collapse-all' : 'expand-all')}
            </Button>
          </>
        }
        searchValue={filters.searchQuery as string}
        filterCount={isNotEmpty(filters.entityType) ? 1 : undefined}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <CollectionLoader
        expandOrCollapseAllState={expandOrCollapseAllState}
        expandedItems={expandedItems}
        filters={filters}
        onChangeFilters={onChangeFilters}
        onToggleExpandItem={onToggleExpandItem}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
