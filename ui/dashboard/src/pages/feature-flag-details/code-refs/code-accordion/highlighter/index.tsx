import { useMemo } from 'react';
import { Highlight, themes } from 'prism-react-renderer';
import { CodeReference } from '@types';
import { cn } from 'utils/style';
import './style.css';

const supportedExtensions = ['kt', 'swift', 'go', 'dart', 'js', 'ts'];

const CodeHighlighter = ({
  featureId,
  codeRef
}: {
  featureId: string;
  codeRef: CodeReference;
}) => {
  const language = useMemo(() => {
    let lang = codeRef.fileExtension.replace('.', '');

    if (!supportedExtensions.includes(lang)) {
      lang = 'js';
    }
    return lang;
  }, [codeRef]);

  return (
    <Highlight
      theme={themes.github}
      code={codeRef.codeSnippet}
      language={language}
    >
      {({ style, tokens, getLineProps, getTokenProps }) => (
        <pre
          className="w-max min-w-full"
          style={{
            ...style,
            backgroundColor: '#F8FAFC'
          }}
        >
          {tokens.map((line, i) => {
            const lineProps = getLineProps({ line });

            return (
              <div
                {...lineProps}
                key={i}
                style={{
                  backgroundColor: line.some(token =>
                    token.content.includes(featureId)
                  )
                    ? '#e8e4f1'
                    : 'transparent',
                  fontSize: 14
                }}
              >
                <span
                  className={cn(
                    'inline-block w-16 text-right pr-4 select-none bg-primary-400 text-white text-opacity-90',
                    i === 0 && 'pt-3',
                    i === tokens.length - 1 && 'pb-3'
                  )}
                >
                  {codeRef.lineNumber + i}
                </span>
                {line.map((token, key) => {
                  const tokenProps = getTokenProps({ token });

                  const tokenContent = token.content;
                  const parts = tokenContent.split(
                    new RegExp(`(${featureId})`, 'gi')
                  );
                  return (
                    <span
                      {...tokenProps}
                      key={key}
                      style={{ ...tokenProps.style }}
                      className={tokenProps?.className}
                    >
                      {parts.map((part, index) =>
                        part.toLowerCase() === featureId.toLowerCase() ? (
                          <span
                            key={index}
                            className="text-primary-500"
                            style={{
                              backgroundColor: '#d6cee5'
                            }}
                          >
                            {part}
                          </span>
                        ) : (
                          <span key={index}>{part}</span>
                        )
                      )}
                    </span>
                  );
                })}
              </div>
            );
          })}
        </pre>
      )}
    </Highlight>
  );
};

export default CodeHighlighter;
