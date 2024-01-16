import React, { FC, memo } from 'react';
import ReactSelect, { components } from 'react-select';

export interface Option {
  value: string;
  label: string;
}

export interface SelectProps {
  options: Option[];
  disabled?: boolean;
  clearable?: boolean;
  isLoading?: boolean;
  isMulti?: boolean;
  isSearchable?: boolean;
  value?: Option;
  className?: string;
  onChange: ((option: Option) => void) | ((option: Option[]) => void);
  placeholder?: string;
  customControl?: React.ReactNode;
}

export const Select: FC<SelectProps> = memo(
  ({
    disabled,
    className,
    clearable,
    isLoading,
    isMulti,
    isSearchable,
    onChange,
    options,
    value,
    placeholder,
    customControl,
  }) => {
    const textColor = '#3F3F46';
    const textColorDisabled = '#6B7280';
    const backgroundColor = 'white';
    const backgroundColorDisabled = '#F3F4F6';
    const borderColor = '#D1D5DB';
    const fontSize = '0.875rem';
    const lineHeight = '1.25rem';
    const minHeight = '2.5rem';
    const colourStyles = {
      control: (styles, { isDisabled }) => ({
        ...styles,
        backgroundColor: isDisabled ? backgroundColorDisabled : backgroundColor,
        borderColor: borderColor,
        '&:hover': {
          borderColor: borderColor,
        },
        fontSize: fontSize,
        lineHeight: lineHeight,
        minHeight: minHeight,
        '*': {
          boxShadow: 'none !important',
        },
      }),
      option: (styles, { isFocused, isSelected }) => {
        return {
          ...styles,
          backgroundColor: isFocused
            ? backgroundColorDisabled
            : isSelected
            ? backgroundColor
            : null,
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
          color: textColor,
          ':active': {
            backgroundColor: backgroundColor,
          },
        };
      },
      menu: (base) => ({
        ...base,
        fontSize: fontSize,
        lineHeight: lineHeight,
        color: textColor,
      }),
      multiValueLabel: (base, { isDisabled }) => ({
        ...base,
        color: isDisabled ? textColorDisabled : textColor,
      }),
      singleValue: (base, { isDisabled }) => ({
        ...base,
        color: isDisabled ? textColorDisabled : textColor,
      }),
    };

    if (customControl) {
      return (
        <ReactSelect
          options={options}
          className={className}
          classNamePrefix="react-select"
          styles={colourStyles}
          components={{
            Control: customControl,
          }}
          isDisabled={isLoading || disabled}
          isClearable={clearable}
          isMulti={isMulti}
          isSearchable={isSearchable}
          isLoading={isLoading}
          placeholder={placeholder ? placeholder : ''}
          value={value}
          onChange={onChange}
        />
      );
    }

    return (
      <ReactSelect
        options={options}
        className={className}
        classNamePrefix="react-select"
        styles={colourStyles}
        components={
          disabled && {
            DropdownIndicator: () => null,
            IndicatorSeparator: () => null,
          }
        }
        isDisabled={isLoading || disabled}
        isClearable={clearable}
        isMulti={isMulti}
        isSearchable={isSearchable}
        isLoading={isLoading}
        placeholder={placeholder ? placeholder : ''}
        value={value}
        onChange={onChange}
      />
    );
  }
);
