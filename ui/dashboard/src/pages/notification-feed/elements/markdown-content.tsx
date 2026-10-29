import MDEditor, { MDEditorProps } from '@uiw/react-md-editor';
import rehypeSanitize from 'rehype-sanitize';
import { remark } from 'remark';
import remarkGfm from 'remark-gfm';
import strip from 'strip-markdown';
import { visit } from 'unist-util-visit';
import { cn } from 'utils/style';
import './markdown-content.css';

// Shared by the feed viewer and the publisher's live preview so both render the
// same output. The preview component parses inline HTML found in the markdown,
// so the resulting tree is passed through rehype-sanitize (GitHub-style
// allowlist): only plain formatting elements and attributes are kept, and
// link/image URLs are restricted to http(s)/mailto-like protocols.
export const markdownPreviewOptions: NonNullable<
  MDEditorProps['previewOptions']
> = {
  rehypePlugins: [rehypeSanitize]
};

export const MarkdownContent = ({
  source,
  className
}: {
  source: string;
  className?: string;
}) => {
  return (
    <div
      data-color-mode="light"
      className={cn('markdown-content max-w-none text-sm', className)}
    >
      <MDEditor.Markdown
        source={source}
        className="!bg-transparent"
        {...markdownPreviewOptions}
      />
    </div>
  );
};

const stripProcessor = remark().use(remarkGfm).use(strip);

export const markdownToText = (markdown: string): string =>
  String(stripProcessor.processSync(markdown)).replace(/\s+/g, ' ').trim();

const linkProcessor = remark();

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
