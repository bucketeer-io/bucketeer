import MDEditor from '@uiw/react-md-editor';
import { remark } from 'remark';
import strip from 'strip-markdown';
import { cn } from 'utils/style';
import './markdown-content.css';

// Renders author-supplied Markdown from the publish editor. `MDEditor.Markdown`
// sanitizes the output before it reaches the DOM to prevent XSS, and keeps the
// rendered result visually consistent with the editor's preview tab.
export const MarkdownContent = ({
  source,
  className
}: {
  source: string;
  className?: string;
}) => {
  return (
    // `data-color-mode="light"` keeps MDEditor.Markdown from applying its
    // default dark theme wherever this is used (e.g. inside the draft card).
    // `!bg-transparent` lets it inherit the surrounding surface color.
    <div
      data-color-mode="light"
      className={cn('markdown-content max-w-none text-sm', className)}
    >
      <MDEditor.Markdown source={source} className="!bg-transparent" />
    </div>
  );
};

// Strips Markdown syntax to a plain-text snippet for compact list previews and
// search. Uses remark to parse the Markdown and `strip-markdown` to remove the
// formatting, then collapses whitespace into a single line.
const stripProcessor = remark().use(strip);

export const markdownToText = (markdown: string): string =>
  String(stripProcessor.processSync(markdown)).replace(/\s+/g, ' ').trim();
