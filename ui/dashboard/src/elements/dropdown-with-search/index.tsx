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
import { cn } from 'utils/style';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuSearch,
  DropdownMenuTrigger,
  DropdownOption,
  DropdownValue
} from 'components/dropdown';
import DropdownList from 'elements/dropdown-list';

export interface DropdownMenuWithSearchProps {
  align?: 'start' | 'center' | 'end';
  hidden?: boolean;
  label?: ReactNode;
  placeholder?: string;
  isLoading?: boolean;
  isMultiselect?: boolean;
  options: DropdownOption[];
  selectedOptions?: string[];
  itemSelected?: string;
  createNewOption?: ReactNode;
  triggerClassName?: string;
  contentClassName?: string;
  itemClassName?: string;
  isExpand?: boolean;
  disabled?: boolean;
  trigger?: ReactNode;
  showArrow?: boolean;
  showClear?: boolean;
  ariaLabel?: string;
  inputPlaceholder?: string;
  selectedFieldValue?: string;
  itemSize?: number;
  maxOptions?: number;
  isHasMore?: boolean;
  isLoadingMore?: boolean;
  onHasMoreOptions?: () => void;
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
  onSearchChange?: (value: string) => void;
}

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
  itemClassName,
  isExpand,
  disabled,
  trigger,
  showArrow,
  showClear,
  ariaLabel,
  inputPlaceholder,
  selectedFieldValue = 'value',
  itemSelected,
  itemSize = 44,
  maxOptions = 15,
  isHasMore,
  isLoadingMore,
  onHasMoreOptions,
  notFoundOption,
  additionalElement,
  onSelectOption,
  onKeyDown,
  onClear,
  onSearchChange
}: DropdownMenuWithSearchProps) => {
  const { t } = useTranslation(['common']);
  const inputSearchRef = useRef<HTMLInputElement>(null);
  const contentRef = useRef<HTMLInputElement>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);

  const [isOpen, setIsOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');

  const dropdownOptions = useMemo(
    () =>
      onSearchChange
        ? options
        : options?.filter(item => {
            return !searchValue
              ? item
              : (typeof item?.label === 'object'
                  ? item?.labelText
                  : (item?.label as string)
                )
                  ?.toLowerCase()
                  ?.includes(searchValue?.toLowerCase());
          }),
    [options, searchValue, onSearchChange]
  );

  let timerId: NodeJS.Timeout | null = null;
  if (timerId) clearTimeout(timerId);
  timerId = setTimeout(() => inputSearchRef?.current?.focus(), 50);
  const handleFocusSearchInput = useCallback(() => {}, []);

  const onClearSearchValue = useCallback(() => {
    setSearchValue('');
    onSearchChange?.('');
  }, [onSearchChange]);

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
          { 'hidden-scroll': dropdownOptions?.length > maxOptions },
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
            onSearchChange?.(value);
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
          <DropdownList
            options={dropdownOptions}
            itemSize={itemSize}
            itemSelected={itemSelected}
            maxOptions={maxOptions}
            isMultiselect={isMultiselect}
            selectedOptions={selectedOptions}
            selectedFieldValue={selectedFieldValue}
            additionalElement={additionalElement}
            onSelectOption={onSelectOption}
            className={itemClassName}
            isHasMore={isHasMore}
            isLoadingMore={isLoadingMore}
            onHasMoreOptions={onHasMoreOptions}
          />
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
