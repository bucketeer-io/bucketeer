import { Fragment, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

const AutoWrapText = ({
  text,
  width,
  lineHeight = 20,
  sparate = false
}: {
  text: string;
  width: number;
  lineHeight?: number;
  sparate?: boolean;
}) => {
  const containerRef = useRef(null);
  const [lines, setLines] = useState<string[]>([]);
  const { t } = useTranslation('common');
  useEffect(() => {
    if (!containerRef.current) return;
    const tempLines = [];
    let currentLine = '';

    const testSpan = document.createElement('span');
    testSpan.style.visibility = 'hidden';
    testSpan.style.position = 'absolute';
    testSpan.style.whiteSpace = 'nowrap';
    document.body.appendChild(testSpan);

    [text].forEach(word => {
      const testLine = currentLine ? currentLine + ' ' + word : word;
      testSpan.innerText = testLine;
      if (testSpan.offsetWidth > width) {
        if (currentLine) tempLines.push(currentLine);
        currentLine = word;
      } else {
        currentLine = testLine;
      }
    });

    if (currentLine) tempLines.push(currentLine);
    setLines(tempLines);

    document.body.removeChild(testSpan);
  }, [text, width]);

  return (
    <div ref={containerRef} className="w-full text-xs break-all ">
      {lines.map((line, index) => (
        <Fragment key={index}>
          <div
            key={index}
            style={{ lineHeight: `${lineHeight}px` }}
            className="px-1 pb-[1px] border bg-gray-100 border-gray-200 rounded-sm inline"
          >
            {line}
          </div>
          {sparate && index === lines.length - 1 && t('comma')}
        </Fragment>
      ))}
    </div>
  );
};

export default AutoWrapText;
