import { ReactNode, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { IconCloseFilled } from 'react-icons-material-design';
import {
  OVERFLOW_BADGE_WIDTH,
  useOverflowCount
} from 'hooks/use-overflow-count';
import { cn } from 'utils/style';
import { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';
import OverflowValuesPopover from 'components/overflow-values-popover';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import NameWithTooltip from 'elements/name-with-tooltip';

const SelectMenu = ({
  fieldValues,
  options,
  disabled,
  selectedFieldValue = 'value',
  inputPlaceholderKey = 'search-or-create-tags',
  dropdownPlaceholderKey = 'placeholder-tags',
  collapseOverflow = false,
  overflowSummary,
  onChange,
  onChangeOptions
}: {
  fieldValues: string[];
  options: DropdownOption[];
  disabled?: boolean;
  selectedFieldValue?: string;
  inputPlaceholderKey?: string;
  dropdownPlaceholderKey?: string;
  // Keep values on one line and show "+n" for the ones that don't fit
  collapseOverflow?: boolean;
  // Text below the input, shown when values are hidden, that opens the full list
  overflowSummary?: (count: number) => ReactNode;
  onChange: (values: string[]) => void;
  onChangeOptions?: (options: DropdownOption[]) => void;
}) => {
  const { t } = useTranslation(['form']);
  const rowRef = useRef<HTMLDivElement>(null);
  const getLabel = (item: string) =>
    options.find(({ label, value }) => item === value || item === label)
      ?.label || item;
  const labels = (fieldValues || []).map(item => String(getLabel(item)));
  const { visibleCount, hiddenLabels } = useOverflowCount({
    labels,
    getContainer: () => rowRef.current,
    itemSelector: '[data-select-menu-item]',
    enabled: collapseOverflow
  });
  const isCollapsed = visibleCount !== Infinity;

  const dropdown = (
    <DropdownMenuWithSearch
      label={fieldValues?.length ? 'clear' : ''}
      showClear={!!fieldValues?.length}
      showArrow={!fieldValues?.length}
      ariaLabel={'tag-delete-btn'}
      disabled={disabled}
      inputPlaceholder={t(inputPlaceholderKey)}
      selectedFieldValue={selectedFieldValue}
      trigger={
        fieldValues?.length ? (
          <div className="flex items-center justify-between w-full min-w-0">
            <div
              ref={rowRef}
              className={cn('flex items-center w-full gap-1', {
                'flex-wrap': !collapseOverflow,
                'min-w-0 overflow-hidden': collapseOverflow
              })}
            >
              {fieldValues?.map((item: string, index: number) => {
                if (index >= visibleCount) return null;
                const content = getLabel(item);
                const id = `menu-item-${index}`;
                return (
                  <div
                    key={index}
                    data-select-menu-item
                    className={cn(
                      'flex items-center max-w-full gap-x-2 px-1.5 rounded bg-primary-100',
                      { 'shrink-0': collapseOverflow }
                    )}
                    style={
                      isCollapsed
                        ? { maxWidth: `calc(100% - ${OVERFLOW_BADGE_WIDTH}px)` }
                        : undefined
                    }
                  >
                    <NameWithTooltip
                      asChild
                      id={id}
                      maxLines={1}
                      content={
                        <NameWithTooltip.Content content={content} id={id} />
                      }
                      trigger={
                        <NameWithTooltip.Trigger
                          name={content}
                          id={id}
                          className="typo-para-small py-1 [&>div]:!text-primary-500"
                          maxLines={1}
                          haveAction={false}
                        />
                      }
                    />
                    {!(overflowSummary && isCollapsed) && (
                      <div
                        aria-label="tag-delete-btn"
                        className="flex-center w-3 min-w-3 min-h-full self-stretch cursor-pointer hover:text-gray-900"
                        onClick={() =>
                          onChange(
                            fieldValues.filter((tag: string) => tag !== item)
                          )
                        }
                      >
                        <Icon
                          icon={IconCloseFilled}
                          className="h-full w-3 pointer-events-none"
                        />
                      </div>
                    )}
                  </div>
                );
              })}
              {hiddenLabels.length > 0 && (
                <span
                  title={hiddenLabels.join(', ')}
                  className="flex-center shrink-0 px-1.5 py-1 rounded bg-gray-200 typo-para-small text-gray-600"
                >
                  +{hiddenLabels.length}
                </span>
              )}
            </div>
          </div>
        ) : undefined
      }
      isExpand
      isMultiselect
      placeholder={t(dropdownPlaceholderKey)}
      options={options}
      selectedOptions={fieldValues}
      onClear={() => onChange([])}
      onKeyDown={({ event, searchValue, matchOptions, onClearSearchValue }) => {
        const value: string = matchOptions?.length
          ? (matchOptions[0].value as string)
          : searchValue;
        if (event.key === 'Enter' && !fieldValues?.includes(value)) {
          if (!matchOptions?.length) {
            if (onChangeOptions) {
              onChangeOptions([
                ...options,
                {
                  label: value,
                  value
                }
              ]);
            } else {
              options.push({
                label: value,
                value
              });
            }
          }
          onChange([...fieldValues, value]);
          onClearSearchValue();
        }
      }}
      onSelectOption={value => {
        const isExisted = fieldValues?.find((item: string) => item === value);
        onChange(
          isExisted
            ? fieldValues?.filter((item: string) => item !== value)
            : [...fieldValues, value as string]
        );
      }}
      notFoundOption={(searchValue, onChangeValue) => {
        const isExisted = fieldValues?.find(
          (item: string) => item === searchValue
        );
        return (
          searchValue && (
            <div
              className={cn(
                'flex items-center py-2 px-4 my-1 rounded pointer-events-none',
                {
                  'hover:bg-gray-100 cursor-pointer pointer-events-auto':
                    !isExisted
                }
              )}
              onClick={() => {
                onChange([...fieldValues, searchValue]);
                onChangeValue('');
                if (onChangeOptions) {
                  onChangeOptions([
                    ...options,
                    {
                      label: searchValue,
                      value: searchValue
                    }
                  ]);
                } else {
                  options.push({
                    label: searchValue,
                    value: searchValue
                  });
                }
              }}
            >
              <p className="text-gray-700">
                {t('create-tag-name', {
                  name: searchValue
                })}
              </p>
            </div>
          )
        );
      }}
    />
  );

  if (!overflowSummary) return dropdown;

  return (
    <div className="flex flex-col w-full min-w-0">
      {dropdown}
      {hiddenLabels.length > 0 && (
        <OverflowValuesPopover
          labels={labels}
          text={overflowSummary(labels.length)}
          onRemove={
            disabled
              ? undefined
              : index => onChange(fieldValues.filter((_, i) => i !== index))
          }
        />
      )}
    </div>
  );
};

export default SelectMenu;
