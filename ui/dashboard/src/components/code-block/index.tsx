import { Prism, SyntaxHighlighterProps } from 'react-syntax-highlighter';

const SyntaxHighlighter = Prism as unknown as React.FC<SyntaxHighlighterProps>;

const CodeBlock = ({ ...props }: SyntaxHighlighterProps) => {
  return <SyntaxHighlighter {...props}>{props.children}</SyntaxHighlighter>;
};

export default CodeBlock;
