import { useCallback, useMemo, useRef, useState } from 'react';
import { defaultStaticRanges } from 'react-date-range';
import { useTranslation } from 'react-i18next';
import { useAuth } from 'auth';
import { usePartialState } from 'hooks';
import { pickBy } from 'lodash';
import { AuditLog } from '@types';
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
import { truncNumber } from './utils';

export type ExpandOrCollapseRef = {
  toggle: () => void;
};

const PageContent = () => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();

  const expandOfCollapseRef = useRef<ExpandOrCollapseRef>(null);
  const initRange = useMemo(() => {
    const range = defaultStaticRanges[defaultStaticRanges.length - 1].range();
    return {
      from: range.startDate
        ? truncNumber(range.startDate?.getTime() / 1000)
        : undefined,
      to: range.endDate
        ? truncNumber(range.endDate?.getTime() / 1000)
        : undefined
    };
  }, [defaultStaticRanges]);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<AuditLogsFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'TIMESTAMP',
    orderDirection: 'DESC',
    from: initRange.from,
    to: initRange.to,
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

  const onChangeFilters = useCallback(
    (values: Partial<AuditLogsFilters>) => {
      const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
      onChangSearchParams(options);
      setFilters({ ...values });
      handleResetExpandAuditLog();
    },
    [filters]
  );

  const handleExpandOrCollapseAll = useCallback(
    (auditLogs: AuditLog[]) => {
      setExpandOrCollapseAllState(
        isExpandAll ? ExpandOrCollapse.COLLAPSE : ExpandOrCollapse.EXPAND
      );
      setExpandedItems(isExpandAll ? [] : auditLogs.map(item => item.id));
    },
    [isExpandAll]
  );

  const onToggleExpandItem = useCallback(
    (id: string, auditLogs: AuditLog[]) => {
      const isExistedItem = expandedItems.find(item => item === id);
      const newExpandItems = isExistedItem
        ? expandedItems.filter(item => item !== id)
        : [...expandedItems, id];
      setExpandedItems(newExpandItems);

      if (newExpandItems.length === auditLogs.length)
        return setExpandOrCollapseAllState(ExpandOrCollapse.EXPAND);
      if (expandOrCollapseAllState) setExpandOrCollapseAllState(undefined);
    },
    [expandedItems, expandOrCollapseAllState]
  );

  const handleResetExpandAuditLog = useCallback(() => {
    setExpandOrCollapseAllState(undefined);
    setExpandedItems([]);
  }, []);

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
              onClick={() => expandOfCollapseRef.current?.toggle()}
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
        ref={expandOfCollapseRef}
        expandOrCollapseAllState={expandOrCollapseAllState}
        expandedItems={expandedItems}
        filters={filters}
        onChangeFilters={onChangeFilters}
        onToggleExpandItem={onToggleExpandItem}
        handleExpandOrCollapseAll={handleExpandOrCollapseAll}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
