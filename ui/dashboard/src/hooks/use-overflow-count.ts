import { useEffect, useLayoutEffect, useRef, useState } from 'react';

// Space reserved for the "+n" badge rendered after the last visible item
export const OVERFLOW_BADGE_WIDTH = 40;

interface UseOverflowCountOptions {
  labels: string[];
  // Row that holds the items on a single line
  getContainer: () => Element | null | undefined;
  itemSelector: string;
  enabled?: boolean;
  // All items are shown (e.g. while editing), only measure them
  expanded?: boolean;
  // Extra space to keep free in the row (e.g. an inline input)
  reservedWidth?: number;
}

const getOuterWidth = (el: Element) => {
  const style = getComputedStyle(el);
  return (
    (el as HTMLElement).offsetWidth +
    parseFloat(style.marginLeft) +
    parseFloat(style.marginRight)
  );
};

/**
 * Counts how many items fit on one line of the container, leaving room for
 * a "+n" badge. Items must render in `labels` order and stay unshrunk.
 */
export const useOverflowCount = ({
  labels,
  getContainer,
  itemSelector,
  enabled = true,
  expanded = false,
  reservedWidth = 0
}: UseOverflowCountOptions) => {
  const widthCache = useRef(new Map<string, number>());
  const observerRef = useRef<ResizeObserver>(null);
  const observedRef = useRef<Element>(null);
  // Only used to re-render (and re-measure) when the container resizes
  const [, setContainerWidth] = useState(0);
  const [visibleCount, setVisibleCount] = useState(Infinity);

  useEffect(() => () => observerRef.current?.disconnect(), []);

  // Runs after every render: the container may mount later or be replaced
  useLayoutEffect(() => {
    const container = enabled ? getContainer() : null;

    if (container !== observedRef.current) {
      observerRef.current?.disconnect();
      observedRef.current = container ?? null;
      if (container) {
        observerRef.current ??= new ResizeObserver(([entry]) => {
          const width = Math.round(entry.contentRect.width);
          setContainerWidth(prev => {
            // Cached item widths may be clipped by the old width
            if (prev !== width) widthCache.current.clear();
            return width;
          });
        });
        observerRef.current.observe(container);
      }
    }
    if (!container) return;

    const cache = widthCache.current;
    // Items render in label order, so only measure when all of them are shown
    if (expanded || visibleCount >= labels.length) {
      container.querySelectorAll(itemSelector).forEach((item, i) => {
        if (labels[i] !== undefined) cache.set(labels[i], getOuterWidth(item));
      });
    }
    if (expanded) return;

    if (labels.some(label => !cache.has(label))) {
      // Render every item once so the missing widths can be measured
      if (visibleCount !== Infinity) setVisibleCount(Infinity);
      return;
    }

    const style = getComputedStyle(container);
    const gap = parseFloat(style.columnGap) || 0;
    const available =
      container.clientWidth -
      parseFloat(style.paddingLeft) -
      parseFloat(style.paddingRight) -
      reservedWidth +
      gap;
    const widths = labels.map(label => (cache.get(label) as number) + gap);
    const total = widths.reduce((sum, w) => sum + w, 0);

    let next = labels.length;
    if (total > available) {
      let used = 0;
      next = 0;
      while (
        next < widths.length &&
        used + widths[next] + OVERFLOW_BADGE_WIDTH <= available
      ) {
        used += widths[next];
        next++;
      }
      // Always keep one (truncated) item visible
      next = Math.max(1, next);
    }
    if (next !== visibleCount) setVisibleCount(next);
  });

  const effectiveCount =
    !enabled || expanded || visibleCount >= labels.length
      ? Infinity
      : visibleCount;

  return {
    visibleCount: effectiveCount,
    hiddenLabels:
      effectiveCount === Infinity ? [] : labels.slice(effectiveCount)
  };
};
