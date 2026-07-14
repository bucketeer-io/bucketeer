import { useState } from 'react';
import MDEditor, { commands, ICommand } from '@uiw/react-md-editor';
import '@uiw/react-md-editor/markdown-editor.css';
import { useTranslation } from 'i18n';
import { TFunction } from 'i18next';
import { AtSign } from 'lucide-react';
import { cn } from 'utils/style';
import Button from 'components/button';
import '../markdown-content.css';
import './markdown-editor.css';

interface MarkdownEditorProps {
  // stored as Markdown text
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

type Mode = 'edit' | 'preview';

// Inserts an "@" at the cursor to start a mention, mirroring the previous
// editor's behavior. There is no built-in mention command in the library.
const createMention = (t: TFunction): ICommand => ({
  name: 'mention',
  keyCommand: 'mention',
  buttonProps: {
    'aria-label': t('form:insert-mention'),
    title: t('form:mention')
  },
  icon: <AtSign size={12} />,
  execute: (_state, api) => api.replaceSelection('@')
});

// The toolbar buttons, grouped and divider-separated to match the design:
//   Heading | Bold | Italic  ·  Quote | Code | Link  ·  Lists  ·  Attach | @
const createToolbar = (t: TFunction): ICommand[] => [
  // Text styles
  commands.heading,
  commands.bold,
  commands.italic,
  commands.divider,
  // Blocks
  commands.quote,
  commands.code,
  commands.link,
  commands.divider,
  // Lists
  commands.orderedListCommand,
  commands.unorderedListCommand,
  commands.checkedListCommand,
  commands.divider,
  // Attachments & mentions
  commands.image,
  createMention(t)
];

// GitHub-style Markdown editor: a raw-Markdown textarea with a formatting
// toolbar and a custom Write/Preview tab bar, backed by @uiw/react-md-editor.
const MarkdownEditor = ({
  value,
  onChange,
  placeholder
}: MarkdownEditorProps) => {
  const { t } = useTranslation(['common', 'form']);
  const [mode, setMode] = useState<Mode>('edit');

  const tabs: { key: Mode; label: string }[] = [
    { key: 'edit', label: t('common:write') },
    { key: 'preview', label: t('common:preview') }
  ];

  return (
    <div
      data-color-mode="light"
      className="overflow-hidden rounded-lg border border-gray-300"
    >
      {/* Custom Write/Preview tab bar. Drives the editor's `preview` prop. */}
      <div className="flex items-center gap-1 border-b border-gray-200 px-2 pt-2">
        {tabs.map(tab => (
          <Button
            key={tab.key}
            type="button"
            variant="text"
            onClick={() => setMode(tab.key)}
            className={cn(
              'h-auto rounded-t-md border-b-2 px-3 py-1.5 typo-para-medium transition-colors',
              mode === tab.key
                ? 'border-primary-500 text-primary-500'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            )}
          >
            {tab.label}
          </Button>
        ))}
      </div>

      <MDEditor
        value={value}
        onChange={next => onChange(next ?? '')}
        height={320}
        preview={mode}
        visibleDragbar={false}
        // The custom tab bar above replaces the built-in Edit/Preview toggle,
        // and removes the trailing help/mode badge on the right.
        extraCommands={[]}
        // Curated toolbar. Kept visible in preview mode too — the library
        // automatically renders the buttons in a disabled state there.
        commands={createToolbar(t)}
        textareaProps={{
          placeholder: placeholder ?? t('form:description-placeholder')
        }}
      />
    </div>
  );
};

export default MarkdownEditor;
