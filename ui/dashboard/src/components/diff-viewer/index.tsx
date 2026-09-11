import { DependencyList, useEffect, useMemo } from 'react';
import DiffViewer, {
  DiffMethod,
  ReactDiffViewerProps
} from 'react-diff-viewer-continued';
import { useTheme } from 'hooks/use-theme';

interface Props extends ReactDiffViewerProps {
  lineNumber: number;
  prefix: string;
  condition?: boolean;
  deps?: DependencyList;
}

const ReactDiffViewer = ({
  oldValue = '',
  newValue = '',
  showDiffOnly,
  lineNumber,
  prefix,
  condition,
  deps = [],
  ...props
}: Props) => {
  const { theme } = useTheme();
  const isDark = theme === 'dark';

  const diffViewerStyles = useMemo(
    () => ({
      variables: {
        light: {
          codeFoldGutterBackground: 'transparent',
          codeFoldBackground: 'transparent',
          diffViewerBackground: isDark ? '#1B1725' : '#F8FAFC',
          addedBackground: isDark ? '#14532d4d' : '#DCF4DE',
          removedBackground: isDark ? '#4a044e4d' : '#FCE3F3',
          wordAddedBackground: 'transparent',
          wordRemovedBackground: 'transparent'
        },
        dark: {
          codeFoldGutterBackground: 'transparent',
          codeFoldBackground: 'transparent',
          diffViewerBackground: isDark ? '#1B1725' : '#F8FAFC',
          addedBackground: isDark ? '#14532d4d' : '#DCF4DE',
          removedBackground: isDark ? '#4a044e4d' : '#FCE3F3',
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
        color: isDark ? '#B5B0C2' : '#64748B',
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
              color: isDark ? '#B5B0C2' : '#64748B',
              minWidth: 42,
              textAlign: 'right',
              background: isDark ? '#2B1F45' : '#64748B1F'
            }
          },
          'tr.first-line-item td:first-child': {
            paddingTop: 12,
            paddingBottom: 2,
            borderTopLeftRadius: 8
          },
          'tr.first-line-item td:last-child': {
            marginTop: '12px !important'
          },
          'tr.last-line-item td:first-child': {
            paddingBottom: 12,
            paddingTop: 2,
            borderBottomLeftRadius: 8
          },
          'tr.last-line-item td:last-child': {
            marginBottom: 12
          },

          'tr.first-line-item.last-line-item td:first-child': {
            display: 'flex',
            alignItems: 'end',
            justifyContent: 'right'
          }
        }
      },
      diffAdded: {
        marginTop: '1px !important'
      }
    }),
    [isDark]
  );
  const handleAddClassForLines = () => {
    try {
      const rows = document.querySelectorAll(`tr:has(td[class^=${prefix}-])`);
      const firstRow = rows[0];
      const lastRow = rows[rows.length - 1];
      firstRow?.classList?.add('first-line-item');
      lastRow?.classList?.add('last-line-item');
    } catch (error) {
      console.log(error);
    }
  };

  const handleRemoveClassForLines = () => {
    try {
      document
        .querySelectorAll(`tr.first-line-item:has(td[class^=${prefix}-])`)
        ?.forEach(element => element.classList.remove('first-line-item'));
      document
        .querySelectorAll(`tr.last-line-item:has(td[class^=${prefix}-])`)
        ?.forEach(element => element.classList.remove('last-line-item'));
    } catch (error) {
      console.log(error);
    }
  };

  useEffect(() => {
    if (condition) {
      handleAddClassForLines();

      return () => {
        handleRemoveClassForLines();
      };
    }
  }, [...deps]);

  return (
    <DiffViewer
      oldValue={oldValue}
      newValue={newValue}
      splitView={false}
      showDiffOnly={showDiffOnly}
      compareMethod={DiffMethod.LINES}
      hideMarkers
      hideLineNumbers
      useDarkTheme={isDark}
      extraLinesSurroundingDiff={1}
      codeFoldMessageRenderer={() => <></>}
      renderGutter={() => {
        ++lineNumber;
        return <td className={`${prefix}-${lineNumber}`}>{lineNumber}</td>;
      }}
      renderContent={str => (
        <span
          className={`typo-para-small font-sofia-pro ${isDark ? 'text-dark-gray-200' : 'text-gray-600'}`}
          dangerouslySetInnerHTML={{ __html: str }}
        />
      )}
      styles={diffViewerStyles}
      {...props}
    />
  );
};

export default ReactDiffViewer;
