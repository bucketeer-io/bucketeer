import {
  KeyboardEvent,
  ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState
} from 'react';
import { useTranslation } from 'react-i18next';
import {
  FixedSizeList,
  ListChildComponentProps,
  FixedSizeListProps
} from 'react-window';
import { cn } from 'utils/style';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSearch,
  DropdownMenuTrigger,
  DropdownOption,
  DropdownValue
} from 'components/dropdown';

interface RowWithDataProps {
  items: DropdownOption[];
  selectedFieldValue: string;
  isMultiselect?: boolean;
  selectedOptions?: string[];
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
    items,
    isMultiselect,
    selectedOptions,
    selectedFieldValue,
    additionalElement,
    onSelectOption
  } = data;
  const currentItem = items[index];
  return (
    <DropdownMenuItem
      key={index}
      style={style}
      isSelected={selectedOptions?.includes(
        currentItem[selectedFieldValue] as string
      )}
      isMultiselect={isMultiselect}
      value={currentItem.value}
      label={currentItem.label}
      icon={currentItem?.icon}
      disabled={currentItem?.disabled}
      additionalElement={additionalElement && additionalElement(currentItem)}
      onSelectOption={() =>
        onSelectOption(currentItem[selectedFieldValue] as string)
      }
      className="justify-between gap-x-4 [&>div:last-child]:mb-[2px]"
    />
  );
};

const DropdownMenuWithSearch = ({
  align,
  hidden,
  label,
  placeholder,
  isLoading,
  options,
  selectedOptions,
  isMultiselect,
  createNewOption,
  triggerClassName,
  contentClassName,
  isExpand,
  disabled,
  trigger,
  showArrow,
  showClear,
  ariaLabel,
  inputPlaceholder,
  selectedFieldValue = 'value',
  notFoundOption,
  additionalElement,
  onSelectOption,
  onKeyDown,
  onClear
}: {
  align?: 'start' | 'center' | 'end';
  hidden?: boolean;
  label?: ReactNode;
  placeholder?: string;
  isLoading?: boolean;
  isMultiselect?: boolean;
  options: DropdownOption[];
  selectedOptions?: string[];
  createNewOption?: ReactNode;
  triggerClassName?: string;
  contentClassName?: string;
  isExpand?: boolean;
  disabled?: boolean;
  trigger?: ReactNode;
  showArrow?: boolean;
  showClear?: boolean;
  ariaLabel?: string;
  inputPlaceholder?: string;
  selectedFieldValue?: string;
  notFoundOption?: (
    value: string,
    onChangeValue: (value: string) => void
  ) => ReactNode;
  additionalElement?: (item: DropdownOption) => ReactNode;
  onSelectOption: (value: DropdownValue) => void;
  onKeyDown?: ({
    event,
    searchValue,
    matchOptions,
    onClearSearchValue
  }: {
    event: KeyboardEvent<HTMLInputElement>;
    searchValue: string;
    matchOptions: DropdownOption[];
    onClearSearchValue: () => void;
  }) => void;
  onClear?: () => void;
}) => {
  const { t } = useTranslation(['common']);
  const inputSearchRef = useRef<HTMLInputElement>(null);
  const contentRef = useRef<HTMLInputElement>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);

  const [isOpen, setIsOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');

  const dropdownOptions = useMemo(
    () =>
      options?.filter(item => {
        return !searchValue
          ? item
          : (typeof item?.label === 'object'
              ? item?.labelText
              : (item?.label as string)
            )
              ?.toLowerCase()
              ?.includes(searchValue?.toLowerCase());
      }),
    [options, searchValue]
  );

  let timerId: NodeJS.Timeout | null = null;
  if (timerId) clearTimeout(timerId);
  timerId = setTimeout(() => inputSearchRef?.current?.focus(), 50);
  const handleFocusSearchInput = useCallback(() => {}, []);

  const onClearSearchValue = useCallback(() => {
    setSearchValue('');
  }, []);

  useEffect(() => {
    if (hidden) {
      setIsOpen(false);
      onClearSearchValue();
    }
  }, [hidden]);

  return (
    <DropdownMenu
      open={isOpen}
      onOpenChange={open => {
        setIsOpen(open);
        if (open) return handleFocusSearchInput();
        onClearSearchValue();
      }}
    >
      <DropdownMenuTrigger
        ref={triggerRef}
        showClear={showClear}
        disabled={isLoading || disabled}
        placeholder={placeholder}
        label={label}
        trigger={trigger}
        showArrow={showArrow}
        ariaLabel={ariaLabel}
        variant="secondary"
        className={cn('w-full', triggerClassName)}
        onClear={onClear}
      />
      <DropdownMenuContent
        ref={contentRef}
        align={align}
        className={cn(
          'w-[500px] py-0',
          { 'hidden-scroll': dropdownOptions?.length > 15 },
          contentClassName
        )}
        style={
          isExpand
            ? {
                width: triggerRef.current?.offsetWidth,
                maxWidth: triggerRef.current?.offsetWidth
              }
            : {}
        }
      >
        <DropdownMenuSearch
          ref={inputSearchRef}
          value={searchValue}
          placeholder={inputPlaceholder}
          onChange={value => {
            contentRef.current?.scrollTo({
              top: 0,
              behavior: 'smooth'
            });
            setSearchValue(value);
            handleFocusSearchInput();
          }}
          onKeyDown={event =>
            onKeyDown?.({
              event,
              searchValue,
              matchOptions: dropdownOptions,
              onClearSearchValue
            })
          }
        />
        {dropdownOptions?.length > 0 ? (
          <List
            height={
              dropdownOptions.length > 15 ? 200 : dropdownOptions.length * 44
            }
            width={'100%'}
            itemSize={44}
            itemCount={dropdownOptions.length}
            itemData={{
              items: dropdownOptions,
              className: 'justify-between gap-x-4 [&>div:last-child]:mb-[2px]',
              isMultiselect: isMultiselect,
              selectedOptions,
              selectedFieldValue,
              additionalElement,
              onSelectOption
            }}
            className={
              dropdownOptions?.length < 15 ? 'hidden-scroll' : 'small-scroll'
            }
          >
            {RowWithData}
          </List>
        ) : notFoundOption ? (
          notFoundOption(searchValue, value => {
            setSearchValue(value);
            handleFocusSearchInput();
          })
        ) : (
          <div className="flex-center py-2.5 typo-para-medium text-gray-600">
            {t('no-options-found')}
          </div>
        )}
        {createNewOption}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default DropdownMenuWithSearch;
