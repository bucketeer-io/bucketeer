import { forwardRef, ReactNode, useEffect, useState } from 'react';
import { cn } from 'utils/style';
import { Tooltip } from 'components/tooltip';

interface Props {
  id: string;
  maxLines?: number;
  content: ReactNode;
  trigger: ReactNode;
  align?: 'start' | 'center' | 'end';
  asChild?: boolean;
}

interface DefaultTriggerProps {
  id: string;
  name: ReactNode;
  maxLines?: number;
  className?: string;
  haveAction?: boolean;
  onClick?: () => void;
}

interface DefaultContentProps {
  id: string;
  content: ReactNode;
  className?: string;
}

const DefaultTrigger = forwardRef<HTMLDivElement, DefaultTriggerProps>(
  ({ id, name, maxLines = 2, className, haveAction = true, onClick }, ref) => {
    const parentId = `parent-${id}`;
    const childrenId = `children-${id}`;

    return (
      <div
        ref={ref}
        id={parentId}
        className={cn(
          'typo-para-medium text-gray-700',
          { 'text-primary-500 underline cursor-pointer': haveAction },
          className
        )}
        onClick={onClick}
      >
        <div
          className="text-start break-all"
          style={{
            display: '-webkit-box',
            WebkitLineClamp: maxLines,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden'
          }}
        >
          {name}
        </div>
        <div
          id={childrenId}
          className="w-full break-all -z-10 invisible absolute left-0 right-0 text-transparent"
        >
          {name}
        </div>
      </div>
    );
  }
);

const DefaultContent = forwardRef<HTMLDivElement, DefaultContentProps>(
  ({ id, content, className }, ref) => {
    return (
      <div
        ref={ref}
        style={{
          maxWidth: document?.getElementById(`parent-${id}`)?.offsetWidth,
          wordBreak: 'break-all'
        }}
        className={className}
      >
        {content}
      </div>
    );
  }
);

const NameWithTooltip = ({
  id,
  maxLines = 2,
  content,
  trigger,
  align,
  asChild = true
}: Props) => {
  const [isTruncate, setIsTruncate] = useState(false);
  const childrenId = `children-${id}`;

  useEffect(() => {
    const hasMoreThanMaxLines = () => {
      const pElement = document.getElementById(childrenId);
      if (pElement) {
        pElement.style.display = 'none';
        void pElement.offsetHeight;
        pElement.style.display = '';

        const style = window.getComputedStyle(pElement);
        const lineHeight = parseFloat(style.lineHeight);
        const height = pElement.clientHeight;
        return setIsTruncate(height / lineHeight > maxLines);
      }
      return setIsTruncate(false);
    };
    hasMoreThanMaxLines();

    window.addEventListener('resize', hasMoreThanMaxLines);
    return () => window.removeEventListener('resize', hasMoreThanMaxLines);
  }, [childrenId, maxLines]);

  return (
    <Tooltip
      asChild={asChild}
      align={align}
      content={content}
      hidden={!isTruncate}
      trigger={trigger}
      triggerCls={'relative cursor-default select-text'}
    />
  );
};

NameWithTooltip.Trigger = DefaultTrigger;
NameWithTooltip.Content = DefaultContent;

export default NameWithTooltip;
