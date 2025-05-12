import { useEffect } from 'react';
import { Editor, EditorProps, useMonaco } from '@monaco-editor/react';
import Spinner from 'components/spinner';
import './style.css';

interface ReactCodeEditorProps extends EditorProps {
  readOnly?: boolean;
}

export default function ReactCodeEditor(props: ReactCodeEditorProps) {
  const monaco = useMonaco();

  useEffect(() => {
    if (monaco) {
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
          'editorBracketHighlight.foreground1': '#64748B',
          'editorBracketHighlight.foreground2': '#64748B',
          'editorBracketHighlight.foreground3': '#64748B',
          'editorBracketHighlight.foreground4': '#64748B',
          'editorBracketHighlight.foreground5': '#64748B',
          'editorBracketHighlight.foreground6': '#64748B',
          'editorBracketMatch.border': '#64748B1F'
        }
      });
      monaco.editor.setTheme('customLight');
    }
  }, [monaco]);

  return (
    <Editor
      height={170}
      width={'100%'}
      defaultLanguage="json"
      theme="customLight"
      wrapperProps={{
        className:
          'flex border-none outline-none rounded-lg overflow-hidden editor-wrapper'
      }}
      options={{
        minimap: { enabled: false },
        fontSize: 14,
        fontFamily: 'Sofia Pro, san-serif',
        padding: {
          top: 12,
          bottom: 12
        },
        lineNumbersMinChars: 3,
        bracketPairColorization: { enabled: true },
        scrollBeyondLastLine: true,
        smoothScrolling: true,
        wordWrap: 'on',
        automaticLayout: true,
        renderLineHighlight: 'all',
        cursorBlinking: 'smooth',
        cursorSmoothCaretAnimation: true,
        tabSize: 4,
        renderWhitespace: 'boundary',
        folding: true,
        foldingHighlight: true,
        showFoldingControls: 'always',
        highlightActiveIndentGuide: true,
        scrollbar: {
          vertical: 'visible',
          horizontal: 'visible',
          verticalScrollbarSize: 10,
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
      }}
      onMount={editor => {
        editor.onDidChangeCursorSelection(() => {
          const selection = editor.getSelection();
          const editorDomNode = editor.getDomNode();
          const highlightedLines =
            editorDomNode.querySelectorAll('.highlighted-line');
          highlightedLines?.forEach((line: HTMLElement) =>
            line.classList.remove('highlighted-line')
          );
          if (selection && !selection.isEmpty()) {
            const selectionLineNumbers = [];
            for (
              let i = selection.startLineNumber;
              i <= selection.endLineNumber;
              i++
            ) {
              selectionLineNumbers.push(i);
            }
            if (editorDomNode) {
              selectionLineNumbers.forEach(lineNumber => {
                const lineElement = editorDomNode.querySelector(
                  `.view-line:nth-child(${lineNumber})`
                );
                if (lineElement) {
                  lineElement.classList.add('highlighted-line');
                }
              });
            }
          }
        });
      }}
      loading={
        <div className="flex-center w-full gap-x-2 h-[170px] bg-gray-100 animate-pulse duration-200">
          <p className="typo-para-medium text-gray-600 animate-pulse duration-500">
            Loading...
          </p>
          <Spinner />
        </div>
      }
      {...props}
    />
  );
}
