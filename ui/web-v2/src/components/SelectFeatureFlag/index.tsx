import React, { FC, memo } from 'react';
import ReactSelect from 'react-select';

export interface OptionFeatureFlag {
  value: string;
  label: string;
  enabled: boolean;
}

export interface SelectFeatureFlagProps {
  options: OptionFeatureFlag[];
  value?: OptionFeatureFlag;
  className?: string;
  onChange:
    | ((option: OptionFeatureFlag) => void)
    | ((option: OptionFeatureFlag[]) => void);
  placeholder?: string;
}

export const SelectFeatureFlag: FC<SelectFeatureFlagProps> = memo(
  ({ className, onChange, options, value, placeholder }) => {
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
        width: '100%',
      }),
    };

    const formatOptionLabel = ({ label, enabled, ...rest }) => {
      return (
        <div className="flex justify-between space-x-4 pr-2">
          <span className="flex-1 truncate">{label}</span>
          <span
            className={`border rounded-lg text-sm w-11 flex justify-center ${
              enabled
                ? 'bg-primary border-primary text-white'
                : 'bg-gray-100 border-gray-300'
            }`}
          >
            {enabled ? 'On' : 'Off'}
          </span>
        </div>
      );
    };

    return (
      <ReactSelect
        options={options}
        className={className}
        classNamePrefix="react-select"
        styles={colourStyles}
        formatOptionLabel={formatOptionLabel}
        placeholder={placeholder ? placeholder : ''}
        value={value}
        onChange={onChange}
      />
    );
  }
);
