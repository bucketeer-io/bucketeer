import { FC, memo, useLayoutEffect, useRef, useState } from 'react';

import { classNames } from '../../utils/css';

export interface HoverPopoverProps {
  onClick?: () => void;
  onMouseLeave?: () => void;
  render: () => JSX.Element;
  disabled?: boolean;
  alignRight?: boolean;
}

export const HoverPopover: FC<HoverPopoverProps> = memo(
  ({ disabled, onClick, onMouseLeave, render, children, alignRight }) => {
    let timeout;
    const timeoutDuration = 400;
    const [open, setOpen] = useState<boolean>(false);
    const popoverRef = useRef(null);

    const handleMounseEnter = () => {
      clearTimeout(timeout);
      timeout = setTimeout(() => setOpen(true), timeoutDuration);
    };

    const handleMounseLeave = () => {
      clearTimeout(timeout);
      timeout = setTimeout(() => {
        setOpen(false);
        if (onMouseLeave) onMouseLeave();
      }, timeoutDuration);
    };

    useLayoutEffect(() => {
      if (!open) {
        return;
      }
      // Keep the popover inside viewport.
      const bounding = popoverRef.current.getBoundingClientRect();
      const viewRight =
        window.innerWidth || document.documentElement.clientWidth;
      const viewBottom =
        window.innerHeight || document.documentElement.clientHeight;
      if (bounding.right > viewRight) {
        popoverRef.current.style.left = `${viewRight - bounding.right - 10}px`;
      }
      if (bounding.bottom > viewBottom) {
        popoverRef.current.style.top = `${viewBottom - bounding.bottom - 10}px`;
      }
    }, [open, popoverRef]);

    return (
      <div
        onMouseEnter={disabled ? null : handleMounseEnter}
        onMouseLeave={disabled ? null : handleMounseLeave}
        onClick={disabled ? null : onClick}
        className="relative w-fit h-fit"
      >
        {children}
        {open ? (
          <div
            ref={popoverRef}
            className={classNames('absolute z-10', alignRight && 'right-0')}
          >
            {render()}
          </div>
        ) : null}
      </div>
    );
  }
);
