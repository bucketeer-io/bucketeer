import { useEffect, useState, DependencyList, RefObject } from 'react';

export const useIsTruncated = (
  ref: RefObject<HTMLElement>,
  deps: DependencyList = []
) => {
  const [truncated, setTruncated] = useState(false);

  useEffect(() => {
    const el = ref.current;
    if (el) {
      setTruncated(el.scrollWidth > el.clientWidth);
    }
  }, deps);

  return truncated;
};
