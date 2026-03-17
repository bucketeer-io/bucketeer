import {
  ReactNode,
  Ref,
  useCallback,
  useEffect,
  useImperativeHandle,
  useMemo,
  useRef
} from 'react';
import {
  FixedSizeList,
  ListChildComponentProps,
  FixedSizeListProps
} from 'react-window';
import { cn } from 'utils/style';
import {
  DropdownMenuItem,
  DropdownOption,
  DropdownValue
} from 'components/dropdown';

interface RowWithDataProps {
  options: DropdownOption[];
  itemSelected?: string;
  selectedFieldValue?: string;
  isMultiselect?: boolean;
  selectedOptions?: string[];
  className?: string;
  additionalElement?: (item: DropdownOption) => ReactNode;
  onSelectOption: (value: DropdownValue) => void;
}

const List = FixedSizeList as unknown as React.ComponentType<
  FixedSizeListProps & { ref?: Ref<FixedSizeList> }
>;

export interface DropdownListHandle {
  focusFirstEnabled: () => void;
}

const RowWithData = ({
  index,
  style,
  data
}: ListChildComponentProps<RowWithDataProps>) => {
  const {
    options,
    isMultiselect,
    selectedOptions,
    selectedFieldValue = 'value',
    className,
    itemSelected,
    additionalElement,
    onSelectOption
  } = data;
  const currentItem = options[index];

  return (
    <DropdownMenuItem
      key={index}
      data-option-index={index}
      style={style}
      isSelectedItem={itemSelected === currentItem?.value}
      isSelected={selectedOptions?.includes(
        currentItem[selectedFieldValue as keyof DropdownOption] as string
      )}
      isMultiselect={isMultiselect}
      value={currentItem?.value}
      label={currentItem?.label}
      description={currentItem?.description}
      icon={currentItem?.icon}
      disabled={currentItem?.disabled}
      additionalElement={additionalElement && additionalElement(currentItem)}
      onSelectOption={() =>
        onSelectOption(
          currentItem[selectedFieldValue as keyof DropdownOption] as string
        )
      }
      onPointerMove={event => {
        if (
          event.currentTarget
            .closest('[role="menu"]')
            ?.querySelector('input:focus')
        )
          event.preventDefault();
      }}
      onPointerLeave={event => {
        if (
          event.currentTarget
            .closest('[role="menu"]')
            ?.querySelector('input:focus')
        )
          event.preventDefault();
      }}
      className={className}
    />
  );
};

interface DropdownListProps extends RowWithDataProps {
  ref?: Ref<DropdownListHandle>;
  itemSelected?: string;
  height?: number;
  maxHeight?: number;
  width?: string | number;
  itemSize?: number;
  maxOptions?: number;
  isHasMore?: boolean;
  isLoadingMore?: boolean;
  onHasMoreOptions?: () => void;
}

const DropdownList = ({
  ref,
  height,
  maxHeight = 200,
  width = '100%',
  itemSize = 40,
  options,
  maxOptions = 15,
  isMultiselect = false,
  itemSelected,
  selectedOptions,
  selectedFieldValue = 'value',
  className,
  additionalElement,
  isHasMore,
  isLoadingMore,
  onHasMoreOptions,
  onSelectOption
}: DropdownListProps) => {
  const listRef = useRef<FixedSizeList>(null);
  const outerRef = useRef<HTMLDivElement>(null);
  const focusFrameRef = useRef<number | null>(null);

  useEffect(
    () => () => {
      if (focusFrameRef.current !== null)
        cancelAnimationFrame(focusFrameRef.current);
    },
    []
  );

  useImperativeHandle(ref, () => ({
    focusFirstEnabled: () => {
      const index = options.findIndex(option => !option.disabled);
      if (index === -1) return;
      if (focusFrameRef.current !== null)
        cancelAnimationFrame(focusFrameRef.current);
      listRef.current?.scrollToItem(index, 'start');
      // Wait for react-window to mount the row after scrolling.
      focusFrameRef.current = requestAnimationFrame(() => {
        focusFrameRef.current = null;
        outerRef.current
          ?.querySelector<HTMLElement>(
            `[data-option-index="${index}"]:not([data-disabled])`
          )
          ?.focus();
      });
    }
  }));

  const isLoadingMoreRef = useRef(isLoadingMore);
  isLoadingMoreRef.current = isLoadingMore;
  const isFetchingRef = useRef(false);
  const prevOptionsLengthRef = useRef(options.length);

  useEffect(() => {
    if (options.length !== prevOptionsLengthRef.current) {
      prevOptionsLengthRef.current = options.length;
      isFetchingRef.current = false;
    }
  }, [options.length]);

  const maxHeightList = useMemo(
    () =>
      height ??
      (isHasMore || options.length >= maxOptions
        ? maxHeight
        : options.length * itemSize),
    [options, maxOptions, height, itemSize, isHasMore, maxHeight]
  );

  const handleScroll = useCallback(
    ({ scrollOffset }: { scrollOffset: number }) => {
      if (!isHasMore || isLoadingMoreRef.current || isFetchingRef.current)
        return;
      const totalHeight = options.length * itemSize;
      const isNearBottom =
        scrollOffset + maxHeightList >= totalHeight - itemSize;
      if (isNearBottom) {
        isFetchingRef.current = true;
        onHasMoreOptions?.();
      }
    },
    [isHasMore, options.length, itemSize, maxHeightList, onHasMoreOptions]
  );

  return (
    <List
      ref={listRef}
      outerRef={outerRef}
      height={maxHeightList}
      width={width}
      itemSize={itemSize}
      itemCount={options.length}
      itemData={{
        options,
        className: cn(
          'justify-between gap-x-4 [&>div:last-child]:mb-[2px]',
          className
        ),
        itemSelected,
        isMultiselect,
        selectedOptions,
        selectedFieldValue,
        additionalElement,
        onSelectOption
      }}
      className={
        options?.length <= maxOptions ? 'hidden-scroll' : 'small-scroll'
      }
      onScroll={handleScroll}
    >
      {RowWithData}
    </List>
  );
};

export default DropdownList;
