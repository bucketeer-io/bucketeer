import { forwardRef, ReactNode, Ref } from 'react';
import { cn } from 'utils/style';

export type TextProps = {
  text?: string;
  description?: string;
  sub?: ReactNode;
  isLink?: boolean;
  className?: string;
  descClassName?: string;
  isReverse?: boolean;
};

const Text = forwardRef(
  (
    {
      text,
      description,
      sub,
      isLink,
      className,
      descClassName,
      isReverse
    }: TextProps,
    ref: Ref<HTMLDivElement>
  ) => {
    return (
      <div
        ref={ref}
        className={cn('flex flex-col', {
          'flex-col-reverse': isReverse
        })}
      >
        {text && (
          <div className="flex items-center">
            <p
              className={cn(
                'text-gray-700 typo-para-medium mr-2',
                isLink && 'underline text-primary-500',
                className
              )}
            >
              {text}
            </p>
            {sub}
          </div>
        )}
        {description && (
          <p
            className={cn(
              'typo-para-tiny text-gray-500 mt-1 text-left',
              descClassName
            )}
          >
            {description}
          </p>
        )}
      </div>
    );
  }
);

export default Text;
