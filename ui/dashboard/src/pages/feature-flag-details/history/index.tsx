import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import {
  PAGE_PATH_FEATURE_HISTORY,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import dayjs from 'dayjs';
import { usePartialState } from 'hooks';
import pickBy from 'lodash/pickBy';
import { AuditLog, Feature } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { IconCollapse, IconExpand } from '@icons';
import AuditLogDetailsModal from 'pages/audit-logs/elements/audit-logs-modal/audit-log-details';
import { truncNumber } from 'pages/audit-logs/utils';
import Button from 'components/button';
import { ReactDateRangePicker } from 'components/date-range-picker';
import Icon from 'components/icon';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import { HistoriesFilters, ExpandOrCollapse } from './types';

export type ExpandOrCollapseRef = {
  toggle: () => void;
};

const HistoryPage = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['common', 'form']);
  const params = useParams();
  const navigate = useNavigate();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<HistoriesFilters> = searchOptions;

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
    featureId: feature.id,
    ...searchFilters
  } as HistoriesFilters;

  const [filters, setFilters] =
    usePartialState<HistoriesFilters>(defaultFilters);

  const [expandOrCollapseAllState, setExpandOrCollapseAllState] = useState<
    ExpandOrCollapse | undefined
  >(undefined);

  const [expandedItems, setExpandedItems] = useState<string[]>([]);

  const expandOfCollapseRef = useRef<ExpandOrCollapseRef>(null);
  const isExpandAll = useMemo(
    () => expandOrCollapseAllState === ExpandOrCollapse.EXPAND,
    [expandOrCollapseAllState]
  );

  const pathName = useMemo(
    () =>
      `/${params?.envUrlCode || currentEnvironment?.urlCode}${PAGE_PATH_FEATURES}/${feature.id}${PAGE_PATH_FEATURE_HISTORY}`,
    [params, currentEnvironment, feature]
  );

  const auditLogId = useMemo(() => {
    const id = params['*'];
    return id && id !== pathName ? id : '';
  }, [pathName, params]);

  const onChangeFilters = useCallback(
    (values: Partial<HistoriesFilters>) => {
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
    <PageLayout.Content className="pt-0 gap-y-6">
      <Filter
        link={DOCUMENTATION_LINKS.FLAG_HISTORY}
        placeholder={t('form:name-email-search-placeholder')}
        actionClassName="flex-nowrap"
        action={
          <>
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
              className="w-full"
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
        feature={feature}
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
          title={t('form:history-details')}
          auditLogId={auditLogId}
          isOpen={!!auditLogId}
          onClose={() => {
            onChangeFilters({});
            navigate(pathName);
          }}
        />
      )}
    </PageLayout.Content>
  );
};

export default HistoryPage;
