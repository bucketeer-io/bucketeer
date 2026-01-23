import {
  ComponentPropsWithoutRef,
  ElementRef,
  forwardRef,
  FunctionComponent,
  ReactNode,
  Ref,
  useCallback,
  useMemo,
  useRef,
  useState
} from 'react';
import { IconExpandMoreRound } from 'react-icons-material-design';
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { cva } from 'class-variance-authority';
import { t } from 'i18next';
import { capitalize, debounce, isEmpty } from 'lodash';
import { cn } from 'utils/style';
import { IconChecked, IconClose, IconSearch } from '@icons';
import Checkbox from 'components/checkbox';
import Icon from 'components/icon';
import Input, { InputProps } from 'components/input';
import Spinner from 'components/spinner';
import { Tooltip } from 'components/tooltip';
import DropdownList from 'elements/dropdown-list';
import NameWithTooltip from 'elements/name-with-tooltip';

export type DropdownValue = number | string;

export type DropdownOption = {
  label: ReactNode;
  value: DropdownValue;
  icon?: FunctionComponent;
  iconElement?: ReactNode;
  additionalElement?: ReactNode;
  description?: string;
  tooltip?: ReactNode;
  haveCheckbox?: boolean;
  disabled?: boolean;
  labelText?: string;
  [key: string]:
    | DropdownValue
    | boolean
    | FunctionComponent
    | undefined
    | ReactNode;
};

const DropdownMenu = DropdownMenuPrimitive.Root;

const DropdownMenuGroup = DropdownMenuPrimitive.Group;

const DropdownMenuPortal = DropdownMenuPrimitive.Portal;

const triggerVariants = cva(
  [
    'flex items-center px-3 py-[11px] gap-x-3 w-fit border rounded-lg bg-white',
    'disabled:cursor-not-allowed disabled:border-gray-400 disabled:bg-gray-100 disabled:!shadow-none'
  ],
  {
    variants: {
      variant: {
        primary:
          'border-primary-500 hover:shadow-border-primary-500 [&>*]:text-primary-500',
        secondary:
          'border-gray-400 hover:shadow-border-gray-400 [&_div]:text-gray-700 [&_span]:text-gray-600 [&>i]:text-gray-500'
      }
    },
    defaultVariants: {
      variant: 'secondary'
    }
  }
);

const DropdownMenuTrigger = forwardRef<
  ElementRef<typeof DropdownMenuPrimitive.Trigger>,
  ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Trigger> & {
    label?: ReactNode;
    description?: string;
    isExpand?: boolean;
    placeholder?: ReactNode;
    variant?: 'primary' | 'secondary';
    showArrow?: boolean;
    showClear?: boolean;
    trigger?: ReactNode;
    ariaLabel?: string;
    loading?: boolean;
    onClear?: () => void;
  }
>(
  (
    {
      className,
      variant,
      label,
      description,
      isExpand,
      placeholder = '',
      showArrow = true,
      showClear = false,
      trigger,
      ariaLabel,
      loading,
      onClear,
      ...props
    },
    ref
  ) => {
    const clearRef = useRef<HTMLDivElement>(null);
    const handleTriggerMouseDown = useCallback(
      (e: React.MouseEvent) => {
        const currentTarget = e.target as HTMLElement;
        if (
          (clearRef.current && clearRef.current.contains(e.target as Node)) ||
          (ariaLabel && currentTarget?.ariaLabel === ariaLabel)
        ) {
          e.preventDefault();
        }
      },
      [ariaLabel]
    );

    return (
      <DropdownMenuPrimitive.Trigger
        type="button"
        ref={ref}
        className={cn(
          triggerVariants({
            variant
          }),
          {
            'justify-between w-full': isExpand
          },
          className,
          'group'
        )}
        onPointerDown={e => handleTriggerMouseDown(e)}
        {...props}
      >
        <>
          <div className="flex items-center w-full justify-between typo-para-medium overflow-hidden">
            {trigger ? (
              trigger
            ) : label ? (
              <div className="max-w-full truncate">
                {label} {description && <span>{description}</span>}
              </div>
            ) : (
              <p className={'!text-gray-500'}>{placeholder}</p>
            )}
          </div>
          {showClear && label && !loading && (
            <div
              ref={clearRef}
              className="size-6 min-w-6 pointer-events-auto"
              onClick={e => {
                e.stopPropagation();
                e.preventDefault();
                if (onClear) onClear();
              }}
            >
              <Icon
                icon={IconClose}
                size={'md'}
                className="text-gray-500 hover:text-gray-900"
              />
            </div>
          )}
          {showArrow && !loading && (
            <div className="size-6 min-w-6 transition-all duration-200 group-data-[state=closed]:rotate-0 group-data-[state=open]:rotate-180">
              <Icon icon={IconExpandMoreRound} size={'md'} color="gray-500" />
            </div>
          )}
          {loading && (
            <div className="flex-center size-fit">
              <Spinner className="size-4 border-2" />
            </div>
          )}
        </>
      </DropdownMenuPrimitive.Trigger>
    );
  }
);

const DropdownMenuContent = forwardRef<
  ElementRef<typeof DropdownMenuPrimitive.Content>,
  ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Content> & {
    isExpand?: boolean;
  }
>(({ className, sideOffset = 4, isExpand, ...props }, ref) => (
  <DropdownMenuPrimitive.Portal>
    <DropdownMenuPrimitive.Content
      ref={ref}
      sideOffset={sideOffset}
      className={cn(
        'z-50 min-w-[196px] max-h-[252px] overflow-x-hidden overflow-y-auto rounded-lg border bg-white p-1 shadow-dropdown small-scroll',
        'data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2',
        { 'dropdown-menu-expand': isExpand },
        className
      )}
      {...props}
    />
  </DropdownMenuPrimitive.Portal>
));

const DropdownMenuItem = forwardRef<
  React.ElementRef<typeof DropdownMenuPrimitive.Item>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Item> & {
    icon?: FunctionComponent;
    iconElement?: ReactNode;
    isMultiselect?: boolean;
    isSelected?: boolean;
    isSelectedItem?: boolean;
    label?: ReactNode;
    value: DropdownValue;
    description?: string;
    closeWhenSelected?: boolean;
    additionalElement?: ReactNode;
    disabled?: boolean;
    isNormalItem?: boolean;
    onSelectOption?: (value: DropdownValue, event: Event) => void;
  }
>(
  (
    {
      className,
      icon,
      iconElement,
      label,
      value,
      description,
      isMultiselect,
      isSelected,
      isSelectedItem,
      closeWhenSelected = true,
      additionalElement,
      disabled,
      isNormalItem = false,
      onSelectOption,
      ...props
    },
    ref
  ) => {
    const dropdownMenuItemId = useMemo(
      () => `dropdown-menu-item-${label}-${value}`,
      [label, value]
    );
    return (
      <DropdownMenuPrimitive.Item
        ref={ref}
        disabled={disabled}
        className={cn(
          'relative flex items-center w-full cursor-pointer select-none rounded-[5px] p-2 gap-x-2 outline-none transition-colors hover:bg-gray-100 data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
          { '!bg-gray-100': isSelectedItem },
          className
        )}
        onSelect={
          onSelectOption
            ? event => {
                if (!closeWhenSelected || isMultiselect) event.preventDefault();
                return onSelectOption(value, event);
              }
            : undefined
        }
        {...props}
      >
        {isMultiselect && <Checkbox checked={isSelected} />}
        {(iconElement || icon) &&
          (iconElement ? (
            iconElement
          ) : (
            <div className="flex-center size-5">
              <Icon
                icon={icon as FunctionComponent}
                size={'xs'}
                color="gray-600"
              />
            </div>
          ))}

        <div className="flex flex-col gap-y-1.5 w-full overflow-hidden">
          {isNormalItem ? (
            <div className="typo-para-medium leading-5 text-gray-700 truncate">
              {label}
            </div>
          ) : (
            <NameWithTooltip
              id={dropdownMenuItemId}
              content={
                <NameWithTooltip.Content
                  content={label}
                  id={dropdownMenuItemId}
                />
              }
              trigger={
                <NameWithTooltip.Trigger
                  id={dropdownMenuItemId}
                  name={label as string}
                  haveAction={false}
                  maxLines={1}
                  className="cursor-pointer"
                />
              }
              maxLines={1}
            />
          )}
          {description && (
            <p className="typo-para-small text-gray-500">{description}</p>
          )}
        </div>
        {additionalElement}
        {isSelectedItem && <IconChecked className="text-primary-500 w-6" />}
      </DropdownMenuPrimitive.Item>
    );
  }
);

type DropdownSearchProps = InputProps;

const DropdownMenuSearch = forwardRef(
  (
    { value, onChange, ...props }: DropdownSearchProps,
    ref: Ref<HTMLInputElement>
  ) => {
    return (
      <div className="sticky top-0 left-0 right-0 flex items-center w-full px-3 py-[11.5px] gap-x-2 border-b border-gray-200 bg-white z-50">
        <div className="flex-center size-5">
          <Icon icon={IconSearch} size={'xs'} color="gray-500" />
        </div>
        <Input
          {...props}
          ref={ref}
          value={value}
          onChange={onChange}
          className="p-0 border-none rounded-none"
        />
      </div>
    );
  }
);

interface DropdownProps {
  searchModel?: 'instant' | 'debounce';
  labelCustom?: string | ReactNode;
  isTruncate?: boolean;
  isExpand?: boolean;
  cleanable?: boolean;
  isSearchable?: boolean;
  disabled?: boolean;
  multiselect?: boolean;
  showArrow?: boolean;
  loading?: boolean;
  isListItem?: boolean;
  isTooltip?: boolean;
  value?: DropdownValue | DropdownValue[];
  additionalValue?: DropdownValue | DropdownValue[];
  options?: DropdownOption[];
  addititonOptions?: DropdownOption[];
  placeholder?: string | ReactNode;
  className?: string;
  itemClassName?: string;
  contentClassName?: string;
  sideOffsetContent?: number;
  wrapTriggerStyle?: string;
  variant?: 'primary' | 'secondary';
  menuContentSide?: 'top' | 'bottom' | 'left' | 'right';
  alignContent?: 'center' | 'end' | 'start';
  trigger?: ReactNode;
  additionalElement?: (item: DropdownOption) => ReactNode;
  onChange?: (value: DropdownValue | DropdownValue[]) => void;
  onChangeAdditional?: (value: DropdownValue | DropdownValue[]) => void; // for additional options
  onClear?: () => void;
}

const Dropdown: React.FC<DropdownProps> = ({
  searchModel = 'debounce',
  labelCustom,
  isTruncate = true,
  isExpand,
  cleanable,
  isSearchable = false,
  disabled,
  multiselect = false,
  isTooltip = false,
  showArrow,
  loading,
  isListItem = false,
  options = [],
  addititonOptions,
  value,
  additionalValue,
  placeholder,
  itemClassName,
  contentClassName,
  menuContentSide,
  wrapTriggerStyle,
  variant = 'secondary',
  className = 'w-full',
  alignContent = 'start',
  trigger,

  additionalElement,
  sideOffsetContent,
  onChange,
  onChangeAdditional,
  onClear
}) => {
  const [searchValue, setSearchValue] = useState('');
  const [searchDebounce, setSearchDebounce] = useState('');
  const inputSearchRef = useRef<HTMLInputElement | null>(null);
  const triggerRef = useRef<HTMLButtonElement | null>(null);
  const selectedValues = useMemo(() => {
    return Array.isArray(value) ? value : value ? [value] : [];
  }, [value]);

  const filterOptions = useMemo(() => {
    if (!searchDebounce) return options;
    return (
      options?.filter(opt =>
        String(opt.label).toLowerCase().includes(searchDebounce.toLowerCase())
      ) || []
    );
  }, [options, searchDebounce]);

  const triggerLabel = useMemo(() => {
    if (!selectedValues.length || !!labelCustom) return '';
    const selectedOption = options?.find(o => o.value === selectedValues[0]);
    return selectedOption?.label ?? '';
  }, [selectedValues, options, multiselect, value]);

  const debounceSearch = useMemo(
    () => debounce((val: string) => setSearchDebounce(val), 500),
    []
  );

  const handleFocusSearchInput = useCallback(() => {
    setTimeout(() => inputSearchRef?.current?.focus(), 50);
  }, []);

  const handleSearchChange = (val: string) => {
    setSearchValue(val);
    if (searchModel === 'debounce') {
      debounceSearch(val);
    } else {
      setSearchDebounce(val);
    }
  };

  const handleClear = () => {
    if (onClear) {
      onClear();
    } else {
      onChange?.(multiselect ? [] : '');
    }
  };

  return (
    <DropdownMenu
      onOpenChange={open => {
        if (open) return handleFocusSearchInput();
        setSearchDebounce('');
        setSearchValue('');
      }}
    >
      <div className={cn('w-full', { truncate: isTruncate }, wrapTriggerStyle)}>
        <DropdownMenuTrigger
          ref={triggerRef}
          placeholder={placeholder ?? `${capitalize(t('common:select-placeholder'))}...`}
          label={labelCustom ? labelCustom : triggerLabel}
          trigger={trigger}
          disabled={disabled}
          loading={loading}
          variant={variant}
          showArrow={showArrow}
          showClear={(!!selectedValues.length && multiselect) || cleanable}
          onClear={handleClear}
          className={className}
        />
      </div>

      <DropdownMenuContent
        style={isExpand ? { width: triggerRef.current?.offsetWidth } : {}}
        align={alignContent}
        className={contentClassName}
        sideOffset={sideOffsetContent}
        side={menuContentSide}
      >
        {isSearchable && (
          <DropdownMenuSearch
            ref={inputSearchRef}
            value={searchValue}
            onChange={handleSearchChange}
            placeholder={`${t('common:search')}...`}
            onKeyDown={e => e.stopPropagation()}
          />
        )}

        <div className={cn('w-full', { 'pb-1': !!addititonOptions })}>
          {filterOptions && !isEmpty(filterOptions) ? (
            isListItem || multiselect ? (
              <DropdownList
                isMultiselect={multiselect}
                itemSelected={value as string}
                selectedOptions={
                  (Array.isArray(value) ? value : []) as string[]
                }
                additionalElement={additionalElement}
                options={filterOptions}
                onSelectOption={val => onChange?.(val)}
              />
            ) : (
              filterOptions?.map((opt, index) =>
                isTooltip && !!opt.tooltip ? (
                  <Tooltip
                    side="right"
                    sideOffset={10}
                    align="start"
                    className="w-[180px] p-3 bg-white typo-para-small text-gray-600 shadow-card"
                    key={index}
                    content={opt.tooltip}
                    showArrow={false}
                    trigger={
                      <DropdownMenuItem
                        icon={opt.icon}
                        label={opt.label}
                        value={opt.value}
                        isNormalItem
                        disabled={opt?.disabled}
                        onSelectOption={
                          onChange ? val => onChange(val) : undefined
                        }
                      />
                    }
                  />
                ) : (
                  <DropdownMenuItem
                    key={opt.value}
                    value={opt.value as DropdownValue}
                    label={opt.label}
                    icon={opt.icon}
                    description={opt.description}
                    iconElement={opt.iconElement}
                    additionalElement={
                      additionalElement ? additionalElement(opt) : undefined
                    }
                    isSelectedItem={value === opt.value}
                    onSelectOption={val => onChange?.(val)}
                    className={itemClassName}
                  />
                )
              )
            )
          ) : (
            <div className="flex-center py-2.5 typo-para-medium text-gray-600">
              {t('no-options-found')}
            </div>
          )}
        </div>

        {addititonOptions && !isEmpty(addititonOptions) && (
          <div className="pt-1">
            {addititonOptions.map(({ label, value }, index) => (
              <DropdownMenuItem
                key={index}
                label={label}
                value={value}
                isSelectedItem={additionalValue === value}
                onSelectOption={val => onChangeAdditional?.(val)}
              />
            ))}
          </div>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default Dropdown;
export {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuSearch,
  DropdownMenuTrigger
};
