import { useCallback, useState } from 'react';

export const useToggleOpen = (
  initialState: boolean | (() => boolean)
): [boolean, () => void, () => void, (value: boolean) => void] => {
  const [open, onOpenChange] = useState(initialState);
  const onOpen = useCallback(() => onOpenChange(true), []);
  const onClose = useCallback(() => onOpenChange(false), []);
  return [open, onOpen, onClose, onOpenChange];
};
