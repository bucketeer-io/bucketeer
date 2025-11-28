import { useTranslation } from 'react-i18next';
import { IconCloseFilled } from 'react-icons-material-design';
import { cn } from 'utils/style';
import { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import NameWithTooltip from 'elements/name-with-tooltip';

const SelectMenu = ({
  fieldValues,
  options,
  disabled,
  selectedFieldValue = 'value',
  inputPlaceholderKey = 'search-or-create-tags',
  dropdownPlaceholderKey = 'placeholder-tags',
  onChange,
  onChangeOptions
}: {
  fieldValues: string[];
  options: DropdownOption[];
  disabled?: boolean;
  selectedFieldValue?: string;
  inputPlaceholderKey?: string;
  dropdownPlaceholderKey?: string;
  onChange: (values: string[]) => void;
  onChangeOptions?: (options: DropdownOption[]) => void;
}) => {
  const { t } = useTranslation(['form']);
  return (
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
          <div className="flex items-center justify-between w-full ">
            <div className="flex items-center w-full gap-1 flex-wrap">
              {fieldValues?.map((item: string, index: number) => {
                const content =
                  options.find(
                    ({ label, value }) => item === value || item === label
                  )?.label || item;
                const id = `menu-item-${index}`;
                return (
                  <div
                    key={index}
                    className="flex items-center max-w-full gap-x-2 px-1.5 rounded bg-primary-100"
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
                  </div>
                );
              })}
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
};

export default SelectMenu;
