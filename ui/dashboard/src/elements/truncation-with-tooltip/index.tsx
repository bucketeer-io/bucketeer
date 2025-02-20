import { ReactNode, useEffect, useState } from 'react';
import { cn } from 'utils/style';
import { Tooltip } from './tooltip';

type Props = {
  children: ReactNode;
  elementId: string;
  maxSize: number;
  content: string;
  className?: string;
  additionalClassName?: string[];
  tooltipWrapperCls?: string;
  tooltipContentCls?: string;
};

const TruncationWithTooltip = ({
  elementId,
  children,
  maxSize,
  content,
  className = '',
  additionalClassName = ['w-full', 'max-w-full', 'truncate'],
  tooltipWrapperCls = '',
  tooltipContentCls = ''
}: Props) => {
  const [isTruncate, setIsTruncate] = useState(false);

  useEffect(() => {
    const checkTruncation = () => {
      const element = document.getElementById(elementId);
      if (element) {
        const { offsetWidth } = element;
        const size = maxSize - 32;

        if (offsetWidth >= size) {
          element.classList.add(...additionalClassName);
          return setIsTruncate(true);
        }
        return setIsTruncate(false);
      }
    };

    checkTruncation();

    window.addEventListener('resize', checkTruncation);
    return () => window.removeEventListener('resize', checkTruncation);
  }, [elementId, maxSize, additionalClassName]);

  return (
    <div className={cn('w-full relative group', className)}>
      {isTruncate && (
        <Tooltip
          tooltipWrapperCls={tooltipWrapperCls}
          tooltipContentCls={tooltipContentCls}
          content={content}
        />
      )}
      {children}
    </div>
  );
};

export default TruncationWithTooltip;
