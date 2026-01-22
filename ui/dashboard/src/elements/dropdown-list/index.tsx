import {
  DropdownMenuItem,
  DropdownOption,
  DropdownValue
} from 'components/dropdown';
import Spinner from 'components/spinner';
import { ReactNode, useCallback, useEffect, useMemo, useRef } from 'react';
import {
  FixedSizeList,
  FixedSizeListProps,
  ListChildComponentProps
} from 'react-window';
import { cn } from 'utils/style';

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

const LOADING_ROW_HEIGHT = 50;

const List = FixedSizeList as unknown as React.FC<FixedSizeListProps>;

const RowWithData = ({
  index,
  style,
  data
}: ListChildComponentProps<
  RowWithDataProps & {
    isLoadingMore?: boolean;
  }
>) => {
  const {
    options,
    isMultiselect,
    selectedOptions,
    selectedFieldValue = 'value',
    className,
    itemSelected,
    additionalElement,
    onSelectOption,
    isLoadingMore
  } = data;

  // Show loading indicator as the last row when loading more
  if (isLoadingMore && index === options.length) {
    return (
      <div
        style={style}
        className="flex items-center justify-center p-2 border-t"
      >
        <Spinner />
      </div>
    );
  }

  const currentItem = options[index];

  return (
    <DropdownMenuItem
      key={index}
      style={style}
      isSelectedItem={itemSelected === currentItem?.value}
      isSelected={selectedOptions?.includes(
        currentItem[selectedFieldValue] as string
      )}
      isMultiselect={isMultiselect}
      value={currentItem?.value}
      label={currentItem?.label}
      description={currentItem?.description}
      icon={currentItem?.icon}
      disabled={currentItem?.disabled}
      additionalElement={additionalElement && additionalElement(currentItem)}
      onSelectOption={() =>
        onSelectOption(currentItem[selectedFieldValue] as string)
      }
      className={className}
    />
  );
};

interface DropdownListProps extends RowWithDataProps {
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
  const itemCount = useMemo(
    () => options.length + (isHasMore ? 1 : 0),
    [options.length, isHasMore]
  );

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
      height={maxHeightList + (isLoadingMore ? LOADING_ROW_HEIGHT : 0)}
      width={width}
      itemSize={itemSize}
      itemCount={itemCount}
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
        onSelectOption,
        isLoadingMore
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
