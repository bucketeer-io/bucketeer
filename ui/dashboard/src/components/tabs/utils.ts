import type { TabContainedValue } from './types';

export const getTabListContainerCls = (
  contained?: TabContainedValue
): string => {
  if (!contained) return '';
  if (contained === true) return 'container';
  switch (contained) {
    case 'sm':
      return 'px-4 sm:px-0 sm:container';
    case 'md':
      return 'px-4 md:px-0 md:container';
    default:
      return '';
  }
};

export const getTabPanesContainerCls = (
  contained?: TabContainedValue
): string => {
  if (!contained) return '';
  if (contained === true) return 'container';
  switch (contained) {
    case 'sm':
      return 'sm:container';
    case 'md':
      return 'md:container';
    default:
      return '';
  }
};
