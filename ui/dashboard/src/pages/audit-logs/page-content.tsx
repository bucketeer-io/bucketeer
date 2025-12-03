import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import dayjs from 'dayjs';
import { usePartialState } from 'hooks';
import pickBy from 'lodash/pickBy';
import { AuditLog } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { IconCollapse, IconExpand } from '@icons';
import Button from 'components/button';
import { ReactDateRangePicker } from 'components/date-range-picker';
import Icon from 'components/icon';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import AuditLogDetailsModal from './elements/audit-logs-modal/audit-log-details';
import EntityTypeDropdown from './elements/entity-type-dropdown';
import { AuditLogsFilters, ExpandOrCollapse } from './types';
import { truncNumber } from './utils';

export type ExpandOrCollapseRef = {
  toggle: () => void;
};

const PageContent = () => {
  const { t } = useTranslation(['common', 'form']);
  const params = useParams();
  const navigate = useNavigate();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<AuditLogsFilters> = searchOptions;

  const initRange = useMemo(() => {
    return {
      from:
        searchFilters?.range === 'all-time'
          ? undefined
          : truncNumber(
              new Date(
                dayjs().subtract(1, 'month').toDate().setHours(0, 0, 0, 0)
              ).getTime() / 1000
            ),
      to:
        searchFilters?.range === 'all-time'
          ? undefined
          : truncNumber(
              new Date(new Date().setHours(23, 59, 59, 999)).getTime() / 1000
            )
    };
  }, [searchFilters]);

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

  const expandOfCollapseRef = useRef<ExpandOrCollapseRef>(null);
  const isExpandAll = useMemo(
    () => expandOrCollapseAllState === ExpandOrCollapse.EXPAND,
    [expandOrCollapseAllState]
  );
  const auditLogId = useMemo(() => {
    const id = params['*'];
    return id && id !== `${currentEnvironment?.urlCode}/audit-logs` ? id : '';
  }, [params, currentEnvironment]);

  const onChangeFilters = useCallback(
    (values: Partial<AuditLogsFilters>) => {
      const options = pickBy(
        { ...filters, ...values, page: values?.page || 1 },
        v => isNotEmpty(v)
      );
      onChangSearchParams(options);
      setFilters({ ...values, page: values?.page || 1 });
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

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content className="gap-y-6">
      <Filter
        link={DOCUMENTATION_LINKS.AUDIT_LOGS}
        placeholder={t('form:name-email-search-placeholder')}
        name="audit-logs-search"
        action={
          <>
            <EntityTypeDropdown
              className="w-fit"
              isSystemAdmin={!!consoleAccount?.isSystemAdmin}
              entityType={filters?.entityType}
              onChangeFilters={onChangeFilters}
            />
            <ReactDateRangePicker
              from={filters?.from}
              to={filters?.to}
              isAllTime={[filters?.range, searchFilters?.range].includes(
                'all-time'
              )}
              onChange={(startDate, endDate) => {
                onChangeFilters({
                  from: startDate ? startDate?.toString() : undefined,
                  to: endDate ? endDate?.toString() : undefined,
                  range: !startDate && !endDate ? 'all-time' : undefined
                });
              }}
            />
            <Button
              variant={'secondary'}
              onClick={() => expandOfCollapseRef.current?.toggle()}
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
      {!!auditLogId && (
        <AuditLogDetailsModal
          auditLogId={auditLogId}
          isOpen={!!auditLogId}
          onClose={() => {
            onChangeFilters({});
            navigate(
              `/${params?.envUrlCode || currentEnvironment.urlCode}/audit-logs`
            );
          }}
        />
      )}
    </PageLayout.Content>
  );
};

export default PageContent;
