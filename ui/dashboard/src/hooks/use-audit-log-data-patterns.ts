import { useMemo } from 'react';
import { AuditLog } from '@types';
import { areJsonStringsEqual, isJsonString } from 'utils/converts';
import {
  convertJSONToRender,
  formatJSONWithIndent
} from '../pages/audit-logs/utils';

/**
 * Custom hook to handle all audit log data patterns and provide consistent logic
 * for displaying audit logs across different components.
 *
 * AUDIT LOG DATA PATTERNS SUPPORTED:
 *
 * 1. CREATION EVENTS:
 *    - Action type ends with "CREATED"
 *    - Shows: "CREATED" with everything in green (added diff)
 *    - Data manipulation: Shows created entity vs empty previous state
 *
 * 2. DELETION EVENTS:
 *    - Action type ends with "DELETED"
 *    - Shows: "DELETED" with everything in red (deletion diff)
 *    - Data manipulation: Shows empty current state vs deleted entity
 *
 * 3. UPDATE EVENTS:
 *    - Action type ends with "UPDATED" or other patterns
 *    - Shows: "UPDATED" with actual diff between old and new
 *    - No data manipulation needed
 */
export const useAuditLogDataPatterns = (auditLog: AuditLog | undefined) => {
  const entityData = auditLog?.entityData || '';
  const previousEntityData = auditLog?.previousEntityData || '';
  const type = auditLog?.type || '';

  // Helper to detect deletion events
  const isDeletionEvent = useMemo(
    () => type.split('_').at(-1) === 'DELETED',
    [type]
  );

  // Helper to detect creation events
  const isCreationEvent = useMemo(
    () => type.split('_').at(-1) === 'CREATED',
    [type]
  );

  const isSameData = useMemo(
    () => areJsonStringsEqual(entityData, previousEntityData),
    [entityData, previousEntityData]
  );

  // Handle creation events - if the action type is CREATED, it's a creation event
  const isAnyCreationEvent = useMemo(() => isCreationEvent, [isCreationEvent]);

  const isHaveEntityData = useMemo(() => {
    if (!auditLog) return false;

    // Check if current entityData is valid JSON
    const hasCurrentData =
      !!entityData &&
      !!isJsonString(entityData) &&
      !!formatJSONWithIndent(entityData) &&
      !!convertJSONToRender(formatJSONWithIndent(entityData));

    // Check if previousEntityData is valid JSON (for deletion cases)
    const hasPreviousData =
      !!previousEntityData &&
      !!isJsonString(previousEntityData) &&
      !!formatJSONWithIndent(previousEntityData) &&
      !!convertJSONToRender(formatJSONWithIndent(previousEntityData));

    // Return true if we have either current data OR previous data
    return hasCurrentData || hasPreviousData;
  }, [auditLog, entityData, previousEntityData]);

  const parsedEntityData = useMemo(() => {
    try {
      // For creation/update events, use entityData
      if (entityData && isJsonString(entityData)) {
        return JSON.parse(entityData);
      }
      // For deletion events where entityData is empty, use previousEntityData
      if (previousEntityData && isJsonString(previousEntityData)) {
        return JSON.parse(previousEntityData);
      }
      return null;
    } catch {
      return null;
    }
  }, [entityData, previousEntityData]);

  // Determine if this should be treated as "same data" for UI purposes
  const effectiveIsSameData = useMemo(() => {
    // For deletions and creations, never treat as "same data" - always show the diff
    if (isDeletionEvent || isAnyCreationEvent) return false;

    return isSameData;
  }, [isDeletionEvent, isAnyCreationEvent, isSameData]);

  // Determine if tabs/changes should be shown
  const shouldShowChanges = useMemo(() => {
    // For deletions and creations, always show changes/tabs
    if (isDeletionEvent || isAnyCreationEvent) return true;

    return !isSameData;
  }, [isDeletionEvent, isAnyCreationEvent, isSameData]);

  // Determine the appropriate display mode for better UX
  const displayMode = useMemo(() => {
    if (isAnyCreationEvent) return 'created'; // Any creation events
    if (isDeletionEvent) return 'deleted'; // Deletion events
    if (effectiveIsSameData) return 'current-version'; // No changes to show
    return 'updated'; // Regular update events
  }, [isAnyCreationEvent, isDeletionEvent, effectiveIsSameData]);

  // Manipulate data for proper diff display
  const diffEntityData = useMemo(() => {
    if (isDeletionEvent) {
      // For deletions, show empty current state to make everything appear as removed
      return '';
    }
    // For other events, use original entityData
    return entityData;
  }, [isDeletionEvent, entityData]);

  const diffPreviousEntityData = useMemo(() => {
    if (isDeletionEvent) {
      // For deletions, use the entity that was deleted as the "previous" state
      // For old deletion data, both fields have the same content, so use either one
      return entityData || previousEntityData;
    }
    if (isAnyCreationEvent) {
      // For any creation events (old and new patterns),
      // show empty previous state to make everything appear as added (green)
      return '';
    }
    return previousEntityData;
  }, [isDeletionEvent, isAnyCreationEvent, entityData, previousEntityData]);

  return {
    isDeletionEvent,
    isCreationEvent,
    isAnyCreationEvent,
    isSameData,
    isHaveEntityData,
    parsedEntityData,
    effectiveIsSameData,
    shouldShowChanges,
    displayMode,
    entityData: diffEntityData,
    previousEntityData: diffPreviousEntityData
  };
};
