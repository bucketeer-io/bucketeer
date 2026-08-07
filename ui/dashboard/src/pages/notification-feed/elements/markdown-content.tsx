import MDEditor from '@uiw/react-md-editor';
import { remark } from 'remark';
import strip from 'strip-markdown';
import { visit } from 'unist-util-visit';
import { cn } from 'utils/style';
import './markdown-content.css';

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
      <MDEditor.Markdown source={source} className="!bg-transparent" />
    </div>
  );
};

const stripProcessor = remark().use(strip);

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
