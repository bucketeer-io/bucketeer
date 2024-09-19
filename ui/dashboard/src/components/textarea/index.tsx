import { forwardRef, FunctionComponent, MouseEvent, Ref } from 'react';
import Icon from 'components/icon';

export type TextAreaProps = React.DetailedHTMLProps<
  React.TextareaHTMLAttributes<HTMLTextAreaElement>,
  HTMLTextAreaElement
> & {
  iconLeft?: FunctionComponent;
  onClickIcon?: (e: MouseEvent) => void;
};

const TextArea = forwardRef(
  (
    { id, iconLeft: IconLeft, onClickIcon, ...props }: TextAreaProps,
    ref: Ref<HTMLTextAreaElement>
  ) => {
    return (
      <div className="relative w-full">
        <textarea
          ref={ref}
          {...props}
          id={id}
          className="p-3 border border-gray-400 rounded-lg w-full text-gray-700 typo-para-medium resize-none"
        />
        {IconLeft && (
          <button onClick={onClickIcon}>
            <div className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
              <Icon icon={IconLeft} size="sm" />
            </div>
          </button>
        )}
      </div>
    );
  }
);

export default TextArea;
