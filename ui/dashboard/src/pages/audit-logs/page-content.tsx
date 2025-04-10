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

  const [isExpandOrCollapseAll, setIsExpandOrCollapseAll] =
    useState<ExpandOrCollapse>(ExpandOrCollapse.COLLAPSE);

  const isExpandAll = useMemo(
    () => isExpandOrCollapseAll === ExpandOrCollapse.EXPAND,
    [isExpandOrCollapseAll]
  );

  const onChangeFilters = (values: Partial<AuditLogsFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  const handleExpandOrCollapseAll = useCallback(() => {
    setIsExpandOrCollapseAll(
      isExpandAll ? ExpandOrCollapse.COLLAPSE : ExpandOrCollapse.EXPAND
    );
  }, [isExpandAll]);

  return (
    <PageLayout.Content className='gap-y-6'>
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
            <Button variant={'secondary'} onClick={handleExpandOrCollapseAll}>
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
      <CollectionLoader filters={filters} onChangeFilters={onChangeFilters} />
    </PageLayout.Content>
  );
};

export default PageContent;
