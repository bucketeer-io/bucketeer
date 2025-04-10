import { memo, useCallback, useMemo } from 'react';
import DiffViewer, { DiffMethod } from 'react-diff-viewer-continued';
import { DomainEventType } from '@types';
import { AuditLogTab } from '../types';

const ReactDiffViewer = memo(
  ({
    type,
    lineNumber,
    currentTab,
    previousEntityData,
    entityData
  }: {
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

    const formatJSONWithIndent = useCallback((json: string) => {
      const parsedJSON = JSON.parse(json) || {};
      return JSON.stringify(parsedJSON, null, 4);
    }, []);

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

    const convertJSONToRender = (json: string) => {
      if (typeof json != 'string') {
        json = JSON.stringify(json, null, 4);
      }
      json = json
        .replace(
          /("(\\u[\da-fA-F]{4}|\\[^u]|[^\\"])*"(?=\s*:))/g,
          '<span style="color: #e439ac">$1</span>'
        ) // keys
        .replace(
          /(:\s*)("(\\u[\da-fA-F]{4}|\\[^u]|[^\\"])*")/g,
          '$1<span style="color: #40BF42">$2</span>'
        ) // string values
        .replace(/(:\s*)(\d+)/g, '$1<span style="color: #64748B">$2</span>'); // numeric values
      return json;
    };

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
          return <td>{lineNumber}</td>;
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
            }
          },
          codeFold: {
            display: 'none'
          },
          line: {
            display: 'flex',
            width: 'fit-content'
          },
          content: {
            marginLeft: 12
          },
          contentText: {
            display: 'flex',
            color: '#64748B',
            fontSize: 14,
            fontFamily: 'Sofia Pro',
            width: 'fit-content'
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
              ...(isChangesTab
                ? {
                    'tr:nth-child(2) td': {
                      paddingTop: 12,
                      paddingBottom: 2,
                      borderTopLeftRadius: 8
                    },
                    'tr:nth-last-child(2) td': {
                      paddingBottom: 12,
                      borderBottomLeftRadius: 8
                    }
                  }
                : {
                    'tr:first-child td': {
                      paddingTop: 12,
                      paddingBottom: 2,
                      borderTopLeftRadius: 8
                    },
                    'tr:last-child td': {
                      paddingBottom: 12,
                      borderBottomLeftRadius: 8
                    }
                  })
            }
          },
          diffAdded: {
            marginTop: 1
          }
        }}
      />
    );
  }
);

export default ReactDiffViewer;
