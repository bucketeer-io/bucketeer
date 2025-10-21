import { ReactNode, useMemo } from 'react';
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

const List = FixedSizeList as unknown as React.FC<FixedSizeListProps>;

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
  onSelectOption
}: DropdownListProps) => {
  const maxHeightList = useMemo(
    () =>
      height || options.length > maxOptions
        ? maxHeight
        : options.length * itemSize,
    [options, maxOptions, height, itemSize]
  );

  return (
    <List
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
        options?.length < maxOptions ? 'hidden-scroll' : 'small-scroll'
      }
    >
      {RowWithData}
    </List>
  );
};

export default DropdownList;
