import { useCallback, useMemo, useRef, useState } from 'react';
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { t } from 'i18next';
import { capitalize, isEmpty } from 'lodash';
import { cn } from 'utils/style';
import { DropdownMenuContent } from './content';
import { DropdownMenuItem } from './item';
import { OptionsList } from './options-list';
import { DropdownMenuSearch } from './search';
import { DropdownMenuTrigger } from './trigger';
import { DropdownOption, DropdownOptionGroup, DropdownValue } from './types';

export type { DropdownOption, DropdownOptionGroup, DropdownValue };

const DropdownMenu = DropdownMenuPrimitive.Root;
const DropdownMenuGroup = DropdownMenuPrimitive.Group;
const DropdownMenuPortal = DropdownMenuPrimitive.Portal;

const DEBOUNCE_DELAY_MS = 500;

interface DropdownProps {
  searchModel?: 'instant' | 'debounce';
  labelCustom?: string | React.ReactNode;
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
  options?: DropdownOption[];
  groups?: DropdownOptionGroup[];
  placeholder?: string | React.ReactNode;
  className?: string;
  itemClassName?: string;
  contentClassName?: string;
  sideOffsetContent?: number;
  wrapTriggerStyle?: string;
  variant?: 'primary' | 'secondary';
  menuContentSide?: 'top' | 'bottom' | 'left' | 'right';
  alignContent?: 'center' | 'end' | 'start';
  trigger?: React.ReactNode;
  additionalElement?: (item: DropdownOption) => React.ReactNode;
  onChange?: (value: DropdownValue | DropdownValue[]) => void;
  onSearch?: (query: string) => void;
  onClear?: () => void;
  additionalOptions?: DropdownOption[];
  additionalValue?: DropdownValue | DropdownValue[];
  onChangeAdditional?: (value: DropdownValue | DropdownValue[]) => void;
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
  groups,
  value,
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
  onSearch,
  onClear,
  additionalOptions,
  additionalValue,
  onChangeAdditional
}) => {
  const [searchValue, setSearchValue] = useState('');
  const [debouncedQuery, setDebouncedQuery] = useState('');
  const [contentWidth, setContentWidth] = useState<number | undefined>();

  const inputSearchRef = useRef<HTMLInputElement | null>(null);
  const triggerRef = useRef<HTMLButtonElement | null>(null);

  const debounceSetQuery = useRef(
    (() => {
      let timer: ReturnType<typeof setTimeout>;
      return (val: string) => {
        clearTimeout(timer);
        timer = setTimeout(() => setDebouncedQuery(val), DEBOUNCE_DELAY_MS);
      };
    })()
  ).current;

  const selectedValues = useMemo(
    () => (Array.isArray(value) ? value : value !== undefined ? [value] : []),
    [value]
  );

  const filteredOptions = useMemo(() => {
    if (onSearch || !debouncedQuery) return options;
    const query = debouncedQuery.toLowerCase();
    return options.filter(opt => {
      const text = opt.labelText ?? String(opt.label);
      return text.toLowerCase().includes(query);
    });
  }, [options, debouncedQuery, onSearch]);

  const triggerLabel = useMemo(() => {
    if (!selectedValues.length || labelCustom) return '';
    return options.find(o => o.value === selectedValues[0])?.label ?? '';
  }, [selectedValues, options, labelCustom]);

  const resolvedGroups = useMemo<DropdownOptionGroup[]>(() => {
    const base = groups ?? [];
    if (!isEmpty(additionalOptions)) {
      return [
        ...base,
        {
          options: additionalOptions!,
          value: additionalValue,
          onChange: onChangeAdditional
        }
      ];
    }
    return base;
  }, [groups, additionalOptions, additionalValue, onChangeAdditional]);

  const useVirtualList = isListItem || multiselect;

  const handleSearchChange = useCallback(
    (val: string) => {
      setSearchValue(val);
      if (onSearch) {
        onSearch(val);
        return;
      }
      if (searchModel === 'instant') {
        setDebouncedQuery(val);
      } else {
        debounceSetQuery(val);
      }
    },
    [onSearch, searchModel, debounceSetQuery]
  );

  const handleClear = useCallback(() => {
    if (onClear) {
      onClear();
    } else {
      onChange?.(multiselect ? [] : '');
    }
  }, [onClear, onChange, multiselect]);

  const handleOpenChange = useCallback(
    (open: boolean) => {
      if (open) {
        if (isExpand) {
          setContentWidth(triggerRef.current?.offsetWidth);
        }
        requestAnimationFrame(() => inputSearchRef.current?.focus());
        return;
      }
      setDebouncedQuery('');
      setSearchValue('');
    },
    [isExpand]
  );

  const contentStyle = isExpand && contentWidth ? { width: contentWidth } : {};

  return (
    <DropdownMenu onOpenChange={handleOpenChange}>
      <div className={cn('w-full', { truncate: isTruncate }, wrapTriggerStyle)}>
        <DropdownMenuTrigger
          ref={triggerRef}
          placeholder={
            placeholder ?? `${capitalize(t('common:select-placeholder'))}...`
          }
          label={labelCustom ?? triggerLabel}
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
        style={contentStyle}
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

        <div className={cn('w-full', { 'pb-1': resolvedGroups.length > 0 })}>
          <OptionsList
            options={filteredOptions}
            value={value}
            multiselect={multiselect}
            useVirtualList={useVirtualList}
            isTooltip={isTooltip}
            itemClassName={itemClassName}
            additionalElement={additionalElement}
            onChange={onChange}
          />
        </div>

        {resolvedGroups.map((group, index) => (
          <div key={index} className="pt-1 border-t border-gray-100">
            {group.options.map(opt => (
              <DropdownMenuItem
                key={opt.value}
                label={opt.label}
                value={opt.value}
                isSelectedItem={
                  Array.isArray(group.value)
                    ? group.value.includes(opt.value)
                    : group.value === opt.value
                }
                onSelectOption={val => group.onChange?.(val)}
              />
            ))}
          </div>
        ))}
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
