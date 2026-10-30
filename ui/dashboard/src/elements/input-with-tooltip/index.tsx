import { forwardRef, useCallback, useEffect, useRef, useState } from 'react';
import type { FocusEvent, Ref } from 'react';
import { cn } from 'utils/style';
import Input, { InputProps } from 'components/input';
import { Tooltip } from 'components/tooltip';

/**
 * Input that shows an ellipsis when its value overflows and,
 * while not focused, a tooltip with the full value on hover.
 */
const InputWithTooltip = forwardRef(
  (
    { className, value, onFocus, onBlur, ...props }: InputProps,
    ref: Ref<HTMLInputElement>
  ) => {
    const inputRef = useRef<HTMLInputElement | null>(null);
    const [isOverflow, setIsOverflow] = useState(false);
    const [isFocused, setIsFocused] = useState(false);

    const setRefs = useCallback(
      (node: HTMLInputElement | null) => {
        inputRef.current = node;
        if (typeof ref === 'function') ref(node);
        else if (ref) ref.current = node;
      },
      [ref]
    );

    useEffect(() => {
      const element = inputRef.current;
      if (!element) return;

      const checkOverflow = () =>
        setIsOverflow(element.scrollWidth > element.clientWidth);
      checkOverflow();

      const resizeObserver = new ResizeObserver(checkOverflow);
      resizeObserver.observe(element);
      return () => resizeObserver.disconnect();
    }, [value]);

    const showTooltip = isOverflow && !isFocused && !!value;

    return (
      <Tooltip
        align="start"
        className="max-w-[500px] break-all pointer-events-none"
        disableHoverableContent
        content={showTooltip ? String(value) : undefined}
        trigger={
          <div className="w-full">
            <Input
              ref={setRefs}
              value={value}
              className={cn('text-ellipsis', className)}
              onFocus={(event: FocusEvent<HTMLInputElement>) => {
                setIsFocused(true);
                onFocus?.(event);
              }}
              onBlur={(inputValue, event) => {
                setIsFocused(false);
                onBlur?.(inputValue, event);
              }}
              {...props}
            />
          </div>
        }
      />
    );
  }
);

export default InputWithTooltip;
