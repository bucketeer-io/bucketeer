import { useMemo } from 'react';
import { AuditLog } from '@types';
import { areJsonStringsEqual, isJsonString } from 'utils/converts';
import { convertJSONToRender, formatJSONWithIndent } from '../utils';

/**
 * Custom hook to handle all audit log data patterns and provide consistent logic
 * for displaying audit logs across different components.
 *
 * AUDIT LOG DATA PATTERNS SUPPORTED:
 *
 * 1. NEW CREATION EVENTS (after backend fix):
 *    - entityData: actual entity JSON
 *    - previousEntityData: "" (empty, from nil)
 *    - Shows: "CREATED" with everything in green (added diff)
 *
 * 2. OLD CREATION EVENTS (before backend fix - backward compatibility):
 *    - entityData: actual entity JSON (same as previous)
 *    - previousEntityData: actual entity JSON (copy of same)
 *    - Shows: "CREATED" with everything in green (added diff)
 *    - Note: Data is manipulated for diff viewer (created entity vs empty)
 *
 * 3. NEW DELETION EVENTS (after backend fix):
 *    - entityData: "" (empty, from nil)
 *    - previousEntityData: actual deleted entity JSON
 *    - Shows: "DELETED" with everything in red (deletion diff)
 *
 * 4. OLD DELETION EVENTS (before backend fix - backward compatibility):
 *    - entityData: actual entity JSON (same as previous)
 *    - previousEntityData: actual entity JSON (copy of same)
 *    - Shows: "DELETED" with everything in red (deletion diff)
 *    - Note: Data is manipulated for diff viewer (empty vs deleted entity)
 *
 * 5. REGULAR UPDATE EVENTS:
 *    - entityData: updated entity JSON
 *    - previousEntityData: old entity JSON
 *    - Shows: "UPDATED" with actual diff between old and new
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

  // Handle creation events (both new and old patterns)
  // New pattern: previousEntityData is empty (nil from backend)
  // Old pattern: previousEntityData equals entityData but event type is CREATED
  const isOldDataIssue = useMemo(
    () => !previousEntityData && !!entityData,
    [entityData, previousEntityData]
  );

  const isOldCreationPattern = useMemo(
    () => isCreationEvent && isSameData && !!entityData && !!previousEntityData,
    [isCreationEvent, isSameData, entityData, previousEntityData]
  );

  const isAnyCreationEvent = useMemo(
    () => isOldDataIssue || isOldCreationPattern,
    [isOldDataIssue, isOldCreationPattern]
  );

  // BACKWARD COMPATIBILITY: Handle old deletion data that incorrectly had the same entity
  // in both entityData and previousEntityData fields. For these cases, we should still
  // show it as a deletion (everything deleted) rather than "CURRENT VERSION".
  const shouldShowAsDeletion = useMemo(
    () => isDeletionEvent && isSameData && !!entityData && !!previousEntityData,
    [isDeletionEvent, isSameData, entityData, previousEntityData]
  );

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

    return isSameData && !shouldShowAsDeletion;
  }, [isDeletionEvent, isAnyCreationEvent, isSameData, shouldShowAsDeletion]);

  // Determine if tabs/changes should be shown
  const shouldShowChanges = useMemo(() => {
    // For deletions and creations, always show changes/tabs
    if (isDeletionEvent || isAnyCreationEvent) return true;

    return !isSameData || shouldShowAsDeletion;
  }, [isDeletionEvent, isAnyCreationEvent, isSameData, shouldShowAsDeletion]);

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
    isOldDataIssue,
    isOldCreationPattern,
    isAnyCreationEvent,
    isSameData,
    shouldShowAsDeletion,
    isHaveEntityData,
    parsedEntityData,
    effectiveIsSameData,
    shouldShowChanges,
    displayMode,
    entityData: diffEntityData,
    previousEntityData: diffPreviousEntityData
  };
};
