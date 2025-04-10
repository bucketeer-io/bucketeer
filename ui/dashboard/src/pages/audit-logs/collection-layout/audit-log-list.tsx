import { useCallback, useMemo } from 'react';
import { formatLongDateTime } from 'utils/date-time';
import { ExpandOrCollapse } from '../types';
import AuditLogItem from './audit-log-item';
import { FormattedAuditLogs } from './data-collection';

const AuditLogList = ({
  formattedAuditLogs,
  expandOrCollapseAllState,
  expandedItems,
  onToggleExpandItem
}: {
  formattedAuditLogs: FormattedAuditLogs;
  expandOrCollapseAllState?: ExpandOrCollapse;
  expandedItems: string[];
  onToggleExpandItem: (id: string) => void;
}) => {
  const auditLogDates = useMemo(
    () => [...formattedAuditLogs.keys()],
    [formattedAuditLogs]
  );

  const getDateLabel = useCallback((auditLogKey: string) => {
    const date = new Date(auditLogKey);
    const currentDate = new Date();

    if (date.getDate() === currentDate.getDate()) return 'Today';
    if (date.getDate() === currentDate.getDate() - 1) return 'Yesterday';

    return formatLongDateTime({
      value: Math.trunc(date.getTime() / 1000).toString(),
      overrideOptions: { month: 'long', day: 'numeric' }
    });
  }, []);

  return (
    <div className="flex flex-col w-full gap-y-6">
      {auditLogDates?.map(item => {
        return (
          <div key={item} className="flex flex-col items-center w-full gap-y-6">
            <p className="typo-para-medium text-gray-600">
              {getDateLabel(item)}
            </p>
            <div className="flex flex-col w-full gap-y-2">
              {formattedAuditLogs
                .get(item)
                ?.map(item => (
                  <AuditLogItem
                    isExpanded={
                      expandOrCollapseAllState === ExpandOrCollapse.EXPAND ||
                      expandedItems.includes(item.id)
                    }
                    key={item.id}
                    auditLog={item}
                    type={item.type}
                    onClick={() => onToggleExpandItem(item.id)}
                  />
                ))}
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default AuditLogList;
