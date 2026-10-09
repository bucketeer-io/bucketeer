import { RefObject, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Editor, EditorProps, Monaco } from '@monaco-editor/react';
import { useTheme } from 'hooks/use-theme';
import { cn } from 'utils/style';
import { IconExpandSquar } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import Spinner from 'components/spinner';
import { Tooltip } from 'components/tooltip';
import './style.css';

interface ReactCodeEditorProps extends EditorProps {
  readOnly?: boolean;
  isDefaulScroll?: boolean;
  defaultLanguage?: string;
  wrapCls?: string;
  className?: string;
  isResize?: boolean;
  isExpand?: boolean;
  lastLine?: number;
  scrollParent?: RefObject<HTMLElement | null>;
  onExpand?: () => void;
}

export default function ReactCodeEditor({
  onExpand,
  wrapCls,
  isExpand = true,
  isResize = true,
  lastLine = 5,
  scrollParent,
  ...props
}: ReactCodeEditorProps) {
  const { t } = useTranslation('common');
  const { theme } = useTheme();
  const isDark = theme === 'dark';

  const editorOptions = useMemo(
    () => ({
      minimap: { enabled: false },
      fontSize: 14,
      fontFamily: 'Sofia Pro, sans-serif',
      padding: {
        top: 12,
        bottom: lastLine * 12
      },
      lineNumbersMinChars: 3,
      bracketPairColorization: { enabled: true },
      scrollBeyondLastLine: false,
      renderIndentGuides: false,
      renderWhitespace: 'none',
      smoothScrolling: true,
      wordWrap: 'on',
      automaticLayout: true,
      renderLineHighlight: 'all',
      cursorBlinking: 'smooth',
      cursorSmoothCaretAnimation: 'on',
      tabSize: 4,
      folding: true,
      foldingHighlight: true,
      showFoldingControls: 'always',
      scrollbar: {
        vertical: 'visible',
        horizontal: 'visible',
        verticalScrollbarSize: 3,
        horizontalScrollbarSize: 10
      },
      quickSuggestions: true,
      suggestOnTriggerCharacters: true,
      acceptSuggestionOnEnter: 'on',
      tabCompletion: 'on',
      formatOnType: true,
      formatOnPaste: true,
      hideCursorInOverviewRuler: true,
      overviewRulerLanes: 0,
      overviewRulerBorder: false,
      columnSelection: true,
      readOnly: props?.readOnly
    }),
    [props?.readOnly, lastLine]
  );

  const handleBeforeMount = useCallback((monaco: Monaco) => {
    monaco.editor.defineTheme('customLight', {
      base: 'vs',
      inherit: true,
      rules: [
        { token: 'comment', foreground: '008000' },
        { token: 'keyword', foreground: '0000FF' },
        { token: 'string.key.json', foreground: 'E439AC' },
        { token: 'string.value.json', foreground: '40BF42' }
      ],
      colors: {
        'editor.background': '#FAFAFA',
        'editor.foreground': '#475569',
        'editorLineNumber.foreground': '#64748B',
        'editorCursor.foreground': '#000000',
        'editor.selectionBackground': '#64748B1F',
        'editorBracketHighlight.foreground1': '#64748B',
        'editorBracketHighlight.foreground2': '#64748B',
        'editorBracketHighlight.foreground3': '#64748B',
        'editorBracketHighlight.foreground4': '#64748B',
        'editorBracketHighlight.foreground5': '#64748B',
        'editorBracketHighlight.foreground6': '#64748B',
        'editorBracketMatch.border': '#64748B1F'
      }
    });

    monaco.editor.defineTheme('customDark', {
      base: 'vs-dark',
      inherit: true,
      rules: [
        { token: 'comment', foreground: '7D768E' },
        { token: 'keyword', foreground: 'B58CFF' },
        { token: 'string.key.json', foreground: 'E961BD' },
        { token: 'string.value.json', foreground: '66CC68' }
      ],
      colors: {
        'editor.background': '#110D1C',
        'editor.foreground': '#B5B0C2',
        'editorLineNumber.foreground': '#B5B0C2',
        'editorCursor.foreground': '#F2EDF7',
        'editor.selectionBackground': '#2B1F45',
        'editorBracketHighlight.foreground1': '#B5B0C2',
        'editorBracketHighlight.foreground2': '#B5B0C2',
        'editorBracketHighlight.foreground3': '#B5B0C2',
        'editorBracketHighlight.foreground4': '#B5B0C2',
        'editorBracketHighlight.foreground5': '#B5B0C2',
        'editorBracketHighlight.foreground6': '#B5B0C2',
        'editorBracketMatch.border': '#2B1F45'
      }
    });
  }, []);

  return (
    <div
      className={cn(
        'relative w-full min-h-[170px] h-[170px] max-h-[600px] overflow-hidden rounded-lg dark:ring-1 dark:ring-dark-purple-100',
        { 'resize-y': isResize },
        props.className
      )}
    >
      {isExpand && !!onExpand && (
        <Tooltip
          content={t('common:full-screen-mode')}
          trigger={
            <Button
              onClick={onExpand}
              variant="secondary"
              className="absolute w-5 h-5 top-3 right-3 z-10 !p-1 !rounded-md"
            >
              <Icon size="sm" className="w-3 h-3" icon={IconExpandSquar} />
            </Button>
          }
        />
      )}
      <Editor
        {...props}
        height={'100%'}
        width={'100%'}
        defaultLanguage={props.defaultLanguage || 'json'}
        theme={isDark ? 'customDark' : 'customLight'}
        wrapperProps={{
          className: cn(
            'flex border-none outline-none rounded-lg overflow-hidden editor-wrapper',
            wrapCls
          )
        }}
        options={editorOptions}
        beforeMount={handleBeforeMount}
        onMount={editor => {
          const domNode = editor.getDomNode();
          if (!domNode || !!props.isDefaulScroll) return;

          const wheelHandler = (event: WheelEvent) => {
            const scrollTop = editor.getScrollTop();
            const scrollHeight = editor.getScrollHeight();
            const height = editor.getLayoutInfo().height;

            const isAtTop = scrollTop <= 0.5;
            const isAtBottom = scrollTop + height >= scrollHeight - 0.5;
            const scrollingUp = event.deltaY < 0;
            const scrollingDown = event.deltaY > 0;
            const editorHasNoScroll = scrollHeight <= height + 0.5;
            if (
              editorHasNoScroll ||
              (isAtTop && scrollingUp) ||
              (isAtBottom && scrollingDown)
            ) {
              event.preventDefault();
              event.stopImmediatePropagation();
              if (scrollParent) {
                scrollParent?.current?.scrollBy({ top: event.deltaY });
              } else {
                window.scrollBy({ top: event.deltaY });
              }
            }
          };

          domNode.addEventListener('wheel', wheelHandler, {
            passive: false,
            capture: true
          });

          editor.onDidDispose(() => {
            domNode.removeEventListener(
              'wheel',
              wheelHandler as EventListener,
              { capture: true }
            );
          });
        }}
        loading={
          <div className="flex-center w-full gap-x-2 h-[170px] bg-gray-100 dark:bg-dark-black-900 animate-pulse duration-200">
            <p className="typo-para-medium text-gray-600 dark:text-dark-gray-200 animate-pulse duration-500">
              {t('common:loading')}
            </p>
            <Spinner />
          </div>
        }
      />
    </div>
  );
}
