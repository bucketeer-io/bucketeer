import { ReactNode, useEffect, useRef, useState } from 'react';
import { cn } from 'utils/style';
import { Tooltip } from 'components/tooltip';

interface Props {
  text: ReactNode;
  maxLines?: number;
  className?: string;
  align?: 'start' | 'center' | 'end';
}

const TruncateWithTooltip = ({
  text,
  maxLines = 2,
  className,
  align
}: Props) => {
  const textRef = useRef<HTMLDivElement>(null);
  const [isTruncated, setIsTruncated] = useState(false);

  useEffect(() => {
    const check = () => {
      const el = textRef.current;
      if (!el) return;
      setIsTruncated(el.scrollHeight > el.clientHeight);
    };
    check();
    window.addEventListener('resize', check);
    return () => window.removeEventListener('resize', check);
  }, [text, maxLines]);

  const trigger = (
    <div
      ref={textRef}
      className={cn('break-all text-start', className)}
      style={{
        display: '-webkit-box',
        WebkitLineClamp: maxLines,
        WebkitBoxOrient: 'vertical',
        overflow: 'hidden'
      }}
    >
      {text}
    </div>
  );

  return (
    <Tooltip
      align={align}
      content={<div style={{ wordBreak: 'break-all' }}>{text}</div>}
      hidden={!isTruncated}
      trigger={trigger}
      triggerCls="relative cursor-default select-text"
    />
  );
};

export default TruncateWithTooltip;
