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
  DropdownMenuItem,
  DropdownMenuSearch,
  DropdownMenuTrigger,
  DropdownOption,
  DropdownValue
} from 'components/dropdown';

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
  notFoundOption,
  additionalElement,
  onSelectOption,
  onKeyDown
}: {
  align?: 'start' | 'center' | 'end';
  hidden?: boolean;
  label: ReactNode;
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
}) => {
  const { t } = useTranslation(['common']);

  const inputSearchRef = useRef<HTMLInputElement>(null);
  const contentRef = useRef<HTMLInputElement>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);

  const [isOpen, setIsOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');

  const dropdownOptions = useMemo(
    () =>
      options?.filter(item =>
        !searchValue
          ? item
          : (item.label as string)
              .toLowerCase()
              .includes(searchValue.toLowerCase())
      ),
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
        disabled={isLoading || disabled}
        placeholder={placeholder}
        label={label}
        variant="secondary"
        className={cn('w-full', triggerClassName)}
      />
      <DropdownMenuContent
        ref={contentRef}
        align={align}
        className={cn('w-[500px] py-0', contentClassName)}
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
          dropdownOptions.map((item, index) => (
            <DropdownMenuItem
              key={index}
              isSelected={selectedOptions?.includes(item.value as string)}
              isMultiselect={isMultiselect}
              value={item.value}
              label={item.label}
              icon={item?.icon}
              disabled={item?.disabled}
              additionalElement={additionalElement && additionalElement(item)}
              onSelectOption={onSelectOption}
              className="justify-between gap-x-4 [&>div:last-child]:mb-[2px]"
            />
          ))
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
