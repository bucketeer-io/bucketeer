import MDEditor from '@uiw/react-md-editor';
import { remark } from 'remark';
import strip from 'strip-markdown';
import { visit } from 'unist-util-visit';
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

const linkProcessor = remark();

// Finds the first real Markdown link (not an image) in `markdown`, via the
// same remark AST used by `markdownToText`, so link extraction can't
// misfire on image syntax, links inside code spans, or reference-style
// links the way a hand-rolled regex would.
export const firstMarkdownLink = (
  markdown: string
): { label: string; url: string } | null => {
  const tree = linkProcessor.parse(markdown);
  let found: { label: string; url: string } | null = null;
  visit(tree, 'link', node => {
    if (found) return;
    const label = node.children
      .map(child => ('value' in child ? child.value : ''))
      .join('');
    found = { label: label || node.url, url: node.url };
  });
  return found;
};
