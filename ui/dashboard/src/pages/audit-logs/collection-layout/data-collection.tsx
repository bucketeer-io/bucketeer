import { memo, useMemo } from 'react';
import { AuditLog } from '@types';
import TableListContent from 'elements/table-list-content';
import { ExpandOrCollapse } from '../types';
import AuditLogList from './audit-log-list';
import { EmptyCollection } from './empty-collection';

export type FormattedAuditLogs = Map<string, AuditLog[]>;

export const DataCollection = memo(
  ({
    auditLogs,
    expandOrCollapseAllState,
    expandedItems,
    onToggleExpandItem
  }: {
    auditLogs: AuditLog[];
    expandOrCollapseAllState?: ExpandOrCollapse;
    expandedItems: string[];
    onToggleExpandItem: (id: string) => void;
  }) => {
    const isEmpty = useMemo(() => !auditLogs.length, [auditLogs]);
    const formattedAuditLogs: FormattedAuditLogs = useMemo(() => {
      const auditLogMap = new Map();

      if (auditLogs.length) {
        auditLogs.forEach(item => {
          const { timestamp } = item;
          const date = new Date(+timestamp * 1000);
          const key = `${date.getMonth()}-${date.getDate()}-${date.getFullYear()}`;
          const isExistedKey = auditLogMap.has(key);
          auditLogMap.set(key, [
            ...(isExistedKey ? auditLogMap.get(key) : []),
            item
          ]);
        });
        return auditLogMap;
      }
      return auditLogMap;
    }, [auditLogs]);

    if (isEmpty) return <EmptyCollection />;
    return (
      <TableListContent>
        <AuditLogList
          expandOrCollapseAllState={expandOrCollapseAllState}
          expandedItems={expandedItems}
          formattedAuditLogs={formattedAuditLogs}
          onToggleExpandItem={onToggleExpandItem}
        />
      </TableListContent>
    );
  }
);
