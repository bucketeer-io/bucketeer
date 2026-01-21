import { ReactNode, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  FixedSizeList,
  ListChildComponentProps,
  FixedSizeListProps
} from 'react-window';
import { cn } from 'utils/style';
import Button from 'components/button';
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

const LOAD_MORE_BUTTON_HEIGHT = 50;

const List = FixedSizeList as unknown as React.FC<FixedSizeListProps>;

const RowWithData = ({
  index,
  style,
  data
}: ListChildComponentProps<
  RowWithDataProps & {
    isHasMore?: boolean;
    isLoadingMore?: boolean;
    onHasMoreOptions?: () => void;
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
    isHasMore,
    isLoadingMore,
    onHasMoreOptions
  } = data;
  const { t } = useTranslation(['common']);
  if (isHasMore && index === options.length) {
    return (
      <div
        style={style}
        className="flex items-center justify-center p-2 border-t"
      >
        <Button
          loading={isLoadingMore}
          variant="text"
          onClick={onHasMoreOptions}
          className="w-full"
        >
          {!isLoadingMore && t('load-more')}
        </Button>
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
  onHasMoreOptions,
  isHasMore,
  isLoadingMore,
  additionalElement,
  onSelectOption
}: DropdownListProps) => {
  const itemCount = useMemo(
    () => options.length + (isHasMore ? 1 : 0),
    [options.length, isHasMore]
  );

  const maxHeightList = useMemo(
    () =>
      height || options.length > maxOptions
        ? maxHeight
        : options.length * itemSize,
    [options, maxOptions, height, itemSize]
  );

  return (
    <List
      height={maxHeightList + (isHasMore ? LOAD_MORE_BUTTON_HEIGHT : 0)}
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
        isLoadingMore,
        additionalElement,
        onSelectOption,
        isHasMore,
        onHasMoreOptions
      }}
      className={cn(
        options?.length < maxOptions ? 'hidden-scroll' : 'small-scroll'
      )}
    >
      {RowWithData}
    </List>
  );
};

export default DropdownList;
