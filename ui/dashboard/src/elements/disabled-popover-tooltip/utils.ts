import type { PopoverOption, PopoverValue } from 'components/popover';

export const formatPopoverOptions = (
  options: PopoverOption<PopoverValue>[] | undefined,
  tooltipContent: string,
  content?: string
): PopoverOption<PopoverValue>[] => {
  if (!tooltipContent) return options || [];
  return (options || []).map(item =>
    // Read-only options stay available for accounts without edit access
    item.allowViewer
      ? item
      : {
          ...item,
          tooltip: tooltipContent || item?.tooltip,
          disabled:
            item.disabled || (!!tooltipContent && !content) || !!item?.tooltip
        }
  );
};
