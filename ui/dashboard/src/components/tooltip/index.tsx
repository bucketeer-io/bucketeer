import React, { forwardRef, ReactNode, Ref } from 'react';
import * as TooltipPrimitive from '@radix-ui/react-tooltip';
import { cn } from 'utils/style';

export type TooltipProps = {
  align?: 'start' | 'center' | 'end';
  delayDuration?: number;
  hidden?: boolean;
  content?: ReactNode;
  trigger: ReactNode;
  className?: string;
  alignOffset?: number;
  asChild?: boolean;
  sideOffset?: number;
  triggerCls?: string;
};

const TooltipProvider = TooltipPrimitive.Provider;

const TooltipRoot = TooltipPrimitive.Root;

const TooltipTrigger = TooltipPrimitive.Trigger;
const TooltipPortal = TooltipPrimitive.Portal;
const TooltipArrow = TooltipPrimitive.Arrow;

const TooltipContent = React.forwardRef<
  React.ElementRef<typeof TooltipPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof TooltipPrimitive.Content>
>(({ className, sideOffset = 5, ...props }, ref) => (
  <TooltipPortal>
    <TooltipPrimitive.Content
      ref={ref}
      sideOffset={sideOffset}
      className={cn(
        'data-[state=delayed-open]:data-[side=top]:animate-slideDownAndFade data-[state=delayed-open]:data-[side=right]:animate-slideLeftAndFade data-[state=delayed-open]:data-[side=left]:animate-slideRightAndFade data-[state=delayed-open]:data-[side=bottom]:animate-slideUpAndFade select-none rounded px-3 py-1.5 text-para-medium will-change-[transform,opacity] bg-gray-700 text-white z-50',
        className
      )}
      {...props}
    />
  </TooltipPortal>
));
TooltipContent.displayName = TooltipPrimitive.Content.displayName;

const Tooltip = forwardRef(
  (
    {
      delayDuration = 200,
      align = 'center',
      hidden,
      content,
      trigger,
      className,
      alignOffset = 0,
      asChild = true,
      sideOffset = 5,
      triggerCls
    }: TooltipProps,
    ref: Ref<HTMLDivElement>
  ) => {
    return (
      <TooltipProvider delayDuration={delayDuration}>
        <TooltipRoot>
          <TooltipTrigger
            type="button"
            asChild={asChild}
            className={triggerCls}
          >
            {trigger}
          </TooltipTrigger>
          {content && (
            <TooltipContent
              hidden={hidden}
              ref={ref}
              className={className}
              sideOffset={sideOffset}
              alignOffset={alignOffset}
              align={align}
            >
              {content}
              <TooltipArrow className="fill-gray-700" />
            </TooltipContent>
          )}
        </TooltipRoot>
      </TooltipProvider>
    );
  }
);

export {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
  TooltipProvider,
  TooltipRoot,
  TooltipArrow,
  TooltipPortal
};
