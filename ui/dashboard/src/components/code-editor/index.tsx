import { useEffect } from 'react';
import { Editor, EditorProps, useMonaco } from '@monaco-editor/react';
import './style.css';

export default function ReactCodeEditor(props: EditorProps) {
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
          'editorSelection.foreground': '#000000',
          'editor.selectionBackground': '#64748B88',
          'editor.selectionHighlightBackground': '#add6ff26',
          'editor.inactiveSelectionBackground': '#64748B44',
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
        scrollBeyondLastLine: false,
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
        selection: true,
        selectionHighlight: true,
        occurrencesHighlight: true,
        selectionClipboard: true
      }}
      loading={
        <div className="flex-center w-full h-[170px] bg-gray-100 animate-pulse duration-200" />
      }
      {...props}
    />
  );
}
