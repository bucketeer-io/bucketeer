import { FC, memo } from 'react';
import { MultiValue } from 'react-select';
import ReactCreatableSelect from 'react-select/creatable';
import { cn } from 'utils/style';

export interface Option {
  value: string;
  label: string;
}

export interface CreatableSelectProps {
  options?: Option[];
  disabled?: boolean;
  isSearchable?: boolean;
  defaultValues?: Option[];
  closeMenuOnSelect?: boolean;
  className?: string;
  onChange: (options: MultiValue<Option>) => void;
  value?: Option;
  placeholder?: string;
}

export const CreatableSelect: FC<CreatableSelectProps> = memo(
  ({
    disabled,
    isSearchable,
    className,
    onChange,
    options,
    defaultValues,
    closeMenuOnSelect,
    value,
    placeholder = ''
  }) => {
    return (
      <ReactCreatableSelect
        isMulti
        options={options}
        placeholder={placeholder}
        className={className}
        classNamePrefix="react-select"
        classNames={{
          control: ({ menuIsOpen, isFocused, isDisabled }) =>
            cn(
              'flex items-center px-3 py-1.5 gap-x-3 w-full border !outline-none !rounded-lg bg-white !border-gray-400 hover:!shadow-border-gray-400 focus:!shadow-border-gray-400',
              {
                '!shadow-border-gray-400 !outline-none':
                  menuIsOpen || isFocused,
                'disabled:!cursor-not-allowed disabled:!border-gray-400 disabled:!bg-gray-100 disabled:!shadow-none':
                  isDisabled
              }
            ),
          option: ({ isSelected }) =>
            cn(
              '!typo-para-medium !leading-5 !text-gray-700 hover:!bg-gray-100',
              {
                '!bg-gray-100': isSelected
              }
            ),
          input: () => 'm-0 p-0',
          placeholder: () => '!text-gray-500',
          valueContainer: () => '!p-0 !m-0',
          multiValueLabel: () =>
            '!typo-para-medium !text-gray-700 !bg-primary-100/70 !rounded truncate',
          singleValue: () => '!typo-para-medium !text-gray-700 truncate'
        }}
        components={{
          DropdownIndicator: null
        }}
        isDisabled={disabled}
        isSearchable={isSearchable}
        value={value}
        defaultValue={defaultValues}
        onChange={onChange}
        closeMenuOnSelect={closeMenuOnSelect}
      />
    );
  }
);
