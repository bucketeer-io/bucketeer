import { forwardRef, Ref, useImperativeHandle } from 'react';
import { useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { AuditLog } from '@types';
import Pagination from 'components/pagination';
import FormLoading from 'elements/form-loading';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import { DataCollection } from '../collection-layout/data-collection';
import { ExpandOrCollapseRef } from '../page-content';
import { AuditLogsFilters, ExpandOrCollapse } from '../types';
import { useFetchAuditLogs } from './use-fetch-audit-logs';

const CollectionLoader = forwardRef(
  (
    {
      expandOrCollapseAllState,
      expandedItems,
      filters,
      onChangeFilters,
      onToggleExpandItem,
      handleExpandOrCollapseAll
    }: {
      expandOrCollapseAllState?: ExpandOrCollapse;
      expandedItems: string[];
      filters: AuditLogsFilters;
      onChangeFilters: (filters: Partial<AuditLogsFilters>) => void;
      onToggleExpandItem: (id: string, auditLogs: AuditLog[]) => void;
      handleExpandOrCollapseAll: (auditLogs: AuditLog[]) => void;
    },
    ref: Ref<ExpandOrCollapseRef>
  ) => {
    const { consoleAccount } = useAuth();
    const params = useParams();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);

    const {
      data: auditLogCollection,
      isLoading,
      refetch,
      isError
    } = useFetchAuditLogs({
      ...filters,
      environmentId: currentEnvironment?.id,
      enabledFetching: params?.envUrlCode === currentEnvironment?.urlCode
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

    return isError ? (
      <PageLayout.ErrorState onRetry={refetch} />
    ) : (
      <TableListContainer className="px-6 gap-y-6">
        {isLoading ? (
          <FormLoading />
        ) : (
          <DataCollection
            auditLogs={auditLogs}
            expandOrCollapseAllState={expandOrCollapseAllState}
            expandedItems={expandedItems}
            onToggleExpandItem={id => onToggleExpandItem(id, auditLogs)}
          />
        )}
        {!isLoading && (
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
