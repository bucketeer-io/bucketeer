import React from 'react';
import { useTranslation } from 'react-i18next';
import { getCurrentEnvironment, useAuth } from 'auth';
import { usePartialState, useToggleOpen } from 'hooks';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { ReactDateRangePicker } from 'components/date-range-picker';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import EntityTypeDropdown from './elements/entity-type-dropdown';
import { AuditLogsFilters } from './types';

const PageContent = () => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<AuditLogsFilters> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'TIMESTAMP',
    orderDirection: 'DESC',
    ...searchFilters
  } as AuditLogsFilters;

  const [filters, setFilters] =
    usePartialState<AuditLogsFilters>(defaultFilters);

  const onChangeFilters = (values: Partial<AuditLogsFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  return (
    <PageLayout.Content>
      <Filter
        action={
          <>
            <EntityTypeDropdown
              isSystemAdmin={!!consoleAccount?.isSystemAdmin}
              entityType={filters?.entityType}
              onChangeFilters={onChangeFilters}
            />
            <ReactDateRangePicker />
          </>
        }
        searchValue={filters.searchQuery as string}
        filterCount={isNotEmpty(filters.entityType) ? 1 : undefined}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
