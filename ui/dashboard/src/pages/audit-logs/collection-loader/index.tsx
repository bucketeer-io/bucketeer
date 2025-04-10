import { getCurrentEnvironment, useAuth } from 'auth';
import Pagination from 'components/pagination';
import FormLoading from 'elements/form-loading';
import PageLayout from 'elements/page-layout';
import TableListContainer from 'elements/table-list-container';
import { DataCollection } from '../collection-layout/data-collection';
import { AuditLogsFilters } from '../types';
import { useFetchAuditLogs } from './use-fetch-audit-logs';

const CollectionLoader = ({
  filters,
  onChangeFilters
}: {
  filters: AuditLogsFilters;
  onChangeFilters: (filters: Partial<AuditLogsFilters>) => void;
}) => {
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: auditLogCollection,
    isLoading,
    refetch,
    isError
  } = useFetchAuditLogs({
    ...filters,
    isSystemAdmin: consoleAccount?.isSystemAdmin,
    environmentId: currenEnvironment?.id
  });

  const auditLogs = auditLogCollection?.auditLogs || [];
  const totalCount = Number(auditLogCollection?.totalCount) || 0;

  return isError ? (
    <PageLayout.ErrorState onRetry={refetch} />
  ) : (
    <TableListContainer className="px-6 gap-y-6">
      {isLoading ? <FormLoading /> : <DataCollection auditLogs={auditLogs} />}
      {!isLoading && (
        <Pagination
          page={filters.page as number}
          totalCount={totalCount}
          onChange={page => onChangeFilters({ page })}
        />
      )}
    </TableListContainer>
  );
};

export default CollectionLoader;
