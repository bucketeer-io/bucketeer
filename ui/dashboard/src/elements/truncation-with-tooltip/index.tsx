import { ReactNode, useEffect, useState } from 'react';
import { Tooltip } from 'components/tooltip';

type Props = {
  trigger: ReactNode;
  elementId: string;
  maxSize: number;
  content: string;
};

const TruncationWithTooltip = ({
  elementId,
  trigger,
  maxSize,
  content
}: Props) => {
  const [isTruncate, setIsTruncate] = useState(false);

  useEffect(() => {
    const checkTruncation = () => {
      const element = document.getElementById(elementId);
      if (element) {
        const { offsetWidth, parentElement } = element;
        const size = maxSize - 32;
        if (
          offsetWidth > size ||
          (parentElement && parentElement?.offsetWidth < size)
        ) {
          element.classList.add(...['w-full', 'max-w-full', 'truncate']);
          return setIsTruncate(true);
        }
        return setIsTruncate(false);
      }
    };

    checkTruncation();

    // Optional: Recheck on window resize
    window.addEventListener('resize', checkTruncation);
    return () => window.removeEventListener('resize', checkTruncation);
  }, [elementId, maxSize]);

  return <Tooltip hidden={!isTruncate} trigger={trigger} content={content} />;
};

export default TruncationWithTooltip;
