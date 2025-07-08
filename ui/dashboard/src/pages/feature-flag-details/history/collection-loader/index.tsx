import {
  forwardRef,
  Ref,
  useImperativeHandle,
  useLayoutEffect,
  useState
} from 'react';
import { useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { AuditLog, Feature } from '@types';
import { DataCollection } from 'pages/audit-logs/collection-layout/data-collection';
import Pagination from 'components/pagination';
import FormLoading from 'elements/form-loading';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import { ExpandOrCollapseRef } from '..';
import { HistoriesFilters, ExpandOrCollapse } from '../types';
import { useFetchHistories } from './use-fetch-histories';

const CollectionLoader = forwardRef(
  (
    {
      feature,
      expandOrCollapseAllState,
      expandedItems,
      filters,
      onChangeFilters,
      onToggleExpandItem,
      handleExpandOrCollapseAll
    }: {
      feature: Feature;
      expandOrCollapseAllState?: ExpandOrCollapse;
      expandedItems: string[];
      filters: HistoriesFilters;
      onChangeFilters: (filters: Partial<HistoriesFilters>) => void;
      onToggleExpandItem: (id: string, auditLogs: AuditLog[]) => void;
      handleExpandOrCollapseAll: (auditLogs: AuditLog[]) => void;
    },
    ref: Ref<ExpandOrCollapseRef>
  ) => {
    const { consoleAccount } = useAuth();
    const params = useParams();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    const [isPending, setIsPending] = useState(true);

    const {
      data: auditLogCollection,
      isLoading,
      refetch,
      isError
    } = useFetchHistories({
      ...filters,
      featureId: feature.id,
      environmentId: currentEnvironment?.id,
      enabledFetching:
        params?.envUrlCode === currentEnvironment?.urlCode && !isPending
    });

    const auditLogs = auditLogCollection?.auditLogs || [];
    const totalCount = auditLogs.length
      ? Number(auditLogCollection?.totalCount) || 0
      : 0;

    useImperativeHandle(ref, () => {
      return {
        toggle() {
          handleExpandOrCollapseAll(auditLogs);
        }
      };
    }, [auditLogs, handleExpandOrCollapseAll]);

    useLayoutEffect(() => {
      const timerId = setTimeout(() => {
        setIsPending(false);
      }, 1000);
      return () => clearTimeout(timerId);
    }, []);

    return isError ? (
      <PageLayout.ErrorState onRetry={refetch} />
    ) : (
      <TableListContainer className="px-6 gap-y-6">
        {isLoading || isPending ? (
          <FormLoading />
        ) : (
          <DataCollection
            auditLogs={auditLogs}
            expandOrCollapseAllState={expandOrCollapseAllState}
            expandedItems={expandedItems}
            onToggleExpandItem={id => onToggleExpandItem(id, auditLogs)}
          />
        )}
        {!isLoading && !isPending && (
          <Pagination
            page={filters.page as number}
            totalCount={totalCount}
            onChange={page => onChangeFilters({ page })}
          />
        )}
      </TableListContainer>
    );
  }
);

export default CollectionLoader;
