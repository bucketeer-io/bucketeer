import { useMemo } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { pickBy } from 'lodash';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams, useSearchParams } from 'utils/search-params';
import AuditLogDetailsModal from './elements/audit-logs-modal/audit-log-details';
import PageContent from './page-content';

const PageLoader = () => {
  const params = useParams();
  const { searchOptions } = useSearchParams();
  const requestParams = stringifyParams(
    pickBy(searchOptions, v => isNotEmpty(v as string))
  );
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const navigate = useNavigate();

  const auditLogId = useMemo(
    () =>
      params['*'] && params['*'] !== `${currentEnvironment?.urlCode}/audit-logs`
        ? params['*']
        : '',
    [params, currentEnvironment]
  );

  return (
    <>
      <PageContent />
      {!!auditLogId && (
        <AuditLogDetailsModal
          auditLogId={auditLogId}
          isOpen={!!auditLogId}
          onClose={() =>
            navigate(
              `${currentEnvironment.urlCode}/audit-logs${requestParams ? `?${requestParams}` : ''}`
            )
          }
        />
      )}
    </>
  );
};

export default PageLoader;
