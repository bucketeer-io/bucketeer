import { memo, Suspense, useMemo } from 'react';
import { useInView } from 'react-intersection-observer';
import ReactDiffViewer from 'components/diff-viewer';
import { AuditLogTab } from '../types';
import { convertJSONToRender, formatJSONWithIndent } from '../utils';

const AuditLogJSONCompare = memo(
  ({
    isSameData,
    prefix,
    lineNumber,
    currentTab,
    previousEntityData,
    entityData
  }: {
    isSameData: boolean;
    prefix: string;
    lineNumber: number;
    currentTab: AuditLogTab;
    previousEntityData: string;
    entityData: string;
  }) => {
    const { ref, inView } = useInView({
      triggerOnce: true,
      threshold: 0.1
    });

    const isChangesTab = useMemo(
      () => currentTab === AuditLogTab.CHANGES,
      [currentTab]
    );

    const entityDataFormatted = formatJSONWithIndent(entityData);
    const prevEntityDataFormatted = formatJSONWithIndent(previousEntityData);

    const oldValue = useMemo(() => {
      if (
        !isChangesTab ||
        isSameData ||
        entityDataFormatted === prevEntityDataFormatted
      )
        return entityDataFormatted;
      return prevEntityDataFormatted;
    }, [
      isChangesTab,
      entityDataFormatted,
      prevEntityDataFormatted,
      isSameData
    ]);

    const newValue = useMemo(() => entityDataFormatted, [entityDataFormatted]);

    const memoizedOldValue = useMemo(
      () => convertJSONToRender(oldValue),
      [oldValue]
    );
    const memoizedNewValue = useMemo(
      () => convertJSONToRender(newValue),
      [newValue]
    );

    return (
      <div ref={ref}>
        {inView && (
          <Suspense>
            <ReactDiffViewer
              oldValue={memoizedOldValue || ''}
              newValue={memoizedNewValue || ''}
              splitView={false}
              condition={inView}
              deps={[prefix, currentTab, inView]}
              showDiffOnly={
                isChangesTab && !isSameData && oldValue !== newValue
              }
              lineNumber={lineNumber}
              prefix={prefix}
            />
          </Suspense>
        )}
      </div>
    );
  }
);

export default AuditLogJSONCompare;
