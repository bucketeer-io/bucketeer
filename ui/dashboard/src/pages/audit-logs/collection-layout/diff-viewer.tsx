import { memo, useEffect, useMemo } from 'react';
import DiffViewer, { DiffMethod } from 'react-diff-viewer-continued';
import { DomainEventType } from '@types';
import { AuditLogTab } from '../types';
import { convertJSONToRender, formatJSONWithIndent } from '../utils';

const ReactDiffViewer = memo(
  ({
    prefix,
    type,
    lineNumber,
    currentTab,
    previousEntityData,
    entityData
  }: {
    prefix: string;
    type: DomainEventType;
    lineNumber: number;
    currentTab: AuditLogTab;
    previousEntityData: string;
    entityData: string;
  }) => {
    const isChangesTab = useMemo(
      () => currentTab === AuditLogTab.CHANGES,
      [currentTab]
    );

    const entityDataFormatted = formatJSONWithIndent(entityData);
    const prevEntityDataFormatted = formatJSONWithIndent(previousEntityData);

    const isDeleted = useMemo(() => type?.includes('DELETED'), [type]);
    const isCreated = useMemo(() => type?.includes('CREATED'), [type]);
    const isSameData = useMemo(
      () => entityData === previousEntityData,
      [entityData, previousEntityData]
    );
    const oldValue = useMemo(() => {
      if (!isChangesTab) return entityDataFormatted;
      if (isCreated) return '';
      return prevEntityDataFormatted;
    }, [
      isChangesTab,
      entityDataFormatted,
      prevEntityDataFormatted,
      isSameData,
      isCreated
    ]);
    const newValue = useMemo(() => {
      if (!isChangesTab || !isSameData || !isDeleted || isCreated)
        return entityDataFormatted;
      return '';
    }, [
      isChangesTab,
      entityDataFormatted,
      prevEntityDataFormatted,
      isSameData,
      isDeleted
    ]);

    useEffect(() => {
      const handleAddStyleForLine = () => {
        const rows = document.querySelectorAll(`tr:has(td[class^=${prefix}-])`);
        const firstRow = rows[0];
        const lastRow = rows[rows.length - 1];
        firstRow?.classList?.add('first-line-item');
        lastRow?.classList?.add('last-line-item');
      };

      handleAddStyleForLine();

      return () => {
        document
          .querySelectorAll(`tr.first-line-item:has(td[class^=${prefix}-])`)
          ?.forEach(element => element.classList.remove('first-line-item'));
        document
          .querySelectorAll(`tr.last-line-item:has(td[class^=${prefix}-])`)
          ?.forEach(element => element.classList.remove('last-line-item'));
      };
    }, [prefix, currentTab]);

    return (
      <DiffViewer
        oldValue={convertJSONToRender(oldValue)}
        newValue={convertJSONToRender(newValue)}
        splitView={false}
        showDiffOnly={isChangesTab ? true : false}
        compareMethod={DiffMethod.LINES}
        hideMarkers
        hideLineNumbers
        useDarkTheme={false}
        extraLinesSurroundingDiff={1}
        codeFoldMessageRenderer={() => <></>}
        renderGutter={() => {
          ++lineNumber;
          return <td className={`${prefix}-${lineNumber}`}>{lineNumber}</td>;
        }}
        renderContent={str => (
          <span
            className="typo-para-small font-sofia-pro text-gray-600"
            dangerouslySetInnerHTML={{ __html: str }}
          />
        )}
        styles={{
          variables: {
            light: {
              codeFoldGutterBackground: 'transparent',
              codeFoldBackground: 'transparent',
              diffViewerBackground: '#F8FAFC',
              addedBackground: '#DCF4DE',
              removedBackground: '#FCE3F3',
              wordAddedBackground: 'transparent',
              wordRemovedBackground: 'transparent'
            },
            dark: {
              codeFoldGutterBackground: 'transparent',
              codeFoldBackground: 'transparent',
              diffViewerBackground: '#F8FAFC',
              addedBackground: '#DCF4DE',
              removedBackground: '#FCE3F3',
              wordAddedBackground: 'transparent',
              wordRemovedBackground: 'transparent'
            }
          },
          codeFold: {
            display: 'none'
          },
          line: {
            display: 'flex',
            width: 'fit-content',
            className: 'line'
          },
          content: {
            margin: '0 12px'
          },
          contentText: {
            display: 'flex',
            color: '#64748B',
            fontSize: 14,
            fontFamily: 'Sofia Pro',
            width: 'fit-content',
            wordBreak: 'break-all'
          },
          codeFoldGutter: {
            display: 'none'
          },
          diffContainer: {
            display: 'flex',
            paddingLeft: 0,
            borderRadius: 8,
            tbody: {
              tr: {
                minHeight: 24,
                'td:first-child': {
                  padding: '0 12px',
                  fontSize: 14,
                  color: '#64748B',
                  minWidth: 38,
                  textAlign: 'right',
                  background: '#64748B1F'
                }
              },
              'tr.first-line-item td:first-child': {
                paddingTop: 12,
                paddingBottom: 2,
                borderTopLeftRadius: 8
              },
              'tr.first-line-item td:last-child': {
                marginTop: 12
              },
              'tr.last-line-item td:first-child': {
                paddingBottom: 12,
                paddingTop: 2,
                borderBottomLeftRadius: 8
              },
              'tr.last-line-item td:last-child': {
                marginBottom: 12
              }
            }
          },
          diffAdded: {
            marginTop: '1px !important'
          }
        }}
      />
    );
  }
);

export default ReactDiffViewer;
