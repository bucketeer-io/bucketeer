import { ReactNode } from 'react';
import * as TooltipRadix from '@radix-ui/react-tooltip';

export type TooltipProps = {
  content?: string;
  trigger: ReactNode;
};

const Tooltip = ({ content, trigger }: TooltipProps) => {
  return (
    <TooltipRadix.Provider>
      <TooltipRadix.Root>
        <TooltipRadix.Trigger>{trigger}</TooltipRadix.Trigger>
        {content && (
          <TooltipRadix.Portal>
            <TooltipRadix.Content
              className="data-[state=delayed-open]:data-[side=top]:animate-slideDownAndFade data-[state=delayed-open]:data-[side=right]:animate-slideLeftAndFade data-[state=delayed-open]:data-[side=left]:animate-slideRightAndFade data-[state=delayed-open]:data-[side=bottom]:animate-slideUpAndFade select-none rounded-md px-3 py-1.5 typo-para-medium shadow-tooltip will-change-[transform,opacity] bg-gray-700 text-white"
              sideOffset={5}
            >
              {content}
              <TooltipRadix.Arrow className="fill-gray-700" />
            </TooltipRadix.Content>
          </TooltipRadix.Portal>
        )}
      </TooltipRadix.Root>
    </TooltipRadix.Provider>
  );
};

export default Tooltip;
