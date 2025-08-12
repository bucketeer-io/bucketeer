import { FC, memo, ReactNode } from 'react';
import {
  StylesConfig,
  ActionMeta,
  MultiValue,
  GroupBase,
  components,
  MenuListProps,
  CSSObjectWithLabel,
  OptionProps
} from 'react-select';
import ReactCreatableSelect, { CreatableProps } from 'react-select/creatable';
import Spinner from 'components/spinner';

export interface Option {
  value: string;
  label: string;
  [key: string]: string | number | boolean;
}

export interface NoOptionsMessageProps {
  inputValue: string;
  [key: string]: string | number | boolean;
}

export interface CreatableSelectProps
  extends CreatableProps<Option, true, GroupBase<Option>> {
  isMulti?: true;
  loading?: boolean;
  options?: Option[];
  disabled?: boolean;
  isSearchable?: boolean;
  defaultValues?: MultiValue<Option>;
  closeMenuOnSelect?: boolean;
  value?: MultiValue<Option>;
  placeholder?: string;
  className?: string;
  isHiddenCreateNewOption?: boolean;
  onChange: (
    options: MultiValue<Option>,
    actionMeta: ActionMeta<Option>
  ) => void;
  onCreateOption?: (v: string) => void;
  formatCreateLabel?: (v: string) => JSX.Element;
  noOptionsMessage?: (props: NoOptionsMessageProps) => ReactNode;
}

const textColor = '#475569';
const textColorDisabled = '#6B7280';
const backgroundColor = 'white';
const backgroundColorDisabled = '#F3F4F6';
const borderColor = '#CBD5E1';
const fontSize = '1rem';
const lineHeight = '1.25rem';
const minHeight = '3rem';

export const optionStyle = (
  styles: CSSObjectWithLabel,
  props: OptionProps<Option, boolean, GroupBase<Option>>,
  isHiddenCreateNewOption: boolean
) => {
  const { isFocused, data } = props;

  const isNewOption = data?.__isNew__;

  return {
    ...styles,
    backgroundColor: isFocused ? backgroundColor : undefined,
    color: textColor,
    ':hover': {
      backgroundColor: '#FAFAFC !important',
      cursor: 'pointer'
    },
    display: isNewOption && isHiddenCreateNewOption ? 'none' : 'flex'
  };
};

export const colorStyles: StylesConfig<Option, boolean> = {
  control: (styles, { isDisabled }) => ({
    ...styles,
    backgroundColor: isDisabled ? backgroundColorDisabled : backgroundColor,
    borderColor: borderColor,
    '&:hover': {
      borderColor: borderColor
    },
    fontSize: fontSize,
    lineHeight: lineHeight,
    minHeight: minHeight,
    boxShadow: 'none !important',
    borderRadius: '8px'
  }),
  menu: base => ({
    ...base,
    fontSize: fontSize,
    lineHeight: lineHeight,
    color: textColor
  }),
  placeholder: styles => ({
    ...styles,
    color: '#94A3B8 !important'
  }),
  multiValue: (base, { isDisabled }) => ({
    ...base,
    color: isDisabled ? textColorDisabled : textColor,
    backgroundColor: '#E8E4F1 !important',
    borderRadius: '4px'
  }),
  multiValueLabel: base => ({
    ...base,
    color: '#573792 !important'
  }),
  multiValueRemove: base => ({
    ...base,
    color: '#9A87BE !important',
    ':hover': {
      color: '#292C4C !important'
    }
  }),
  indicatorSeparator: base => ({
    ...base,
    display: 'none'
  }),
  singleValue: (styles, { isDisabled }) => ({
    ...styles,
    color: isDisabled ? textColorDisabled : textColor
  })
};

type CustomMenuListProps = MenuListProps<Option, true, GroupBase<Option>> & {
  createNewOption?: ReactNode;
};

export const CustomMenuList = ({
  children,
  createNewOption,
  ...props
}: CustomMenuListProps) => {
  return (
    <components.MenuList {...props} className="!pb-0">
      {children}
      {createNewOption}
    </components.MenuList>
  );
};

export const CreatableSelect: FC<CreatableSelectProps> = memo(
  ({
    isMulti = true,
    loading = false,
    disabled,
    isSearchable,
    className,
    onChange,
    options,
    defaultValues,
    closeMenuOnSelect,
    value,
    placeholder = '',
    onCreateOption,
    formatCreateLabel,
    noOptionsMessage,
    components,
    styles,
    isHiddenCreateNewOption = false,
    ...props
  }) => {
    return (
      <ReactCreatableSelect
        {...props}
        isMulti={isMulti}
        options={options}
        placeholder={placeholder}
        className={className}
        classNamePrefix="react-select"
        styles={{
          option: (styles, props) =>
            optionStyle(styles, props, isHiddenCreateNewOption),
          ...colorStyles,
          ...styles
        }}
        components={{
          ...components,
          DropdownIndicator: null,
          LoadingIndicator: () => <Spinner className="size-5 mr-4" />
        }}
        isDisabled={disabled}
        isLoading={loading}
        isSearchable={isSearchable}
        value={value}
        defaultValue={defaultValues}
        onCreateOption={onCreateOption}
        onChange={(newValue, actionMeta) =>
          onChange(newValue as MultiValue<Option>, actionMeta)
        }
        closeMenuOnSelect={closeMenuOnSelect}
        formatCreateLabel={formatCreateLabel}
        noOptionsMessage={noOptionsMessage}
      />
    );
  }
);
