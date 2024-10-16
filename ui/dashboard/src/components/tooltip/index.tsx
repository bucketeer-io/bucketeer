import React, { forwardRef, ReactNode, Ref } from 'react';
import * as TooltipPrimitive from '@radix-ui/react-tooltip';
import { cn } from 'utils/style';

export type TooltipProps = {
  content?: string;
  trigger: ReactNode;
  className?: string;
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
        'data-[state=delayed-open]:data-[side=top]:animate-slideDownAndFade data-[state=delayed-open]:data-[side=right]:animate-slideLeftAndFade data-[state=delayed-open]:data-[side=left]:animate-slideRightAndFade data-[state=delayed-open]:data-[side=bottom]:animate-slideUpAndFade select-none rounded-md px-3 py-1.5 text-para-medium will-change-[transform,opacity] bg-additional-gray-150 text-white',
        className
      )}
      {...props}
    />
  </TooltipPortal>
));
TooltipContent.displayName = TooltipPrimitive.Content.displayName;

const Tooltip = forwardRef(
  ({ content, trigger, className }: TooltipProps, ref: Ref<HTMLDivElement>) => {
    return (
      <TooltipProvider>
        <TooltipRoot>
          <TooltipTrigger type="button" asChild>
            {trigger}
          </TooltipTrigger>
          {content && (
            <TooltipContent ref={ref} className={className} sideOffset={5}>
              {content}
              <TooltipArrow className="fill-additional-gray-150" />
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
