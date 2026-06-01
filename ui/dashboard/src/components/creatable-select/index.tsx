import { memo, ReactNode } from 'react';
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
import { useTheme } from 'hooks/use-theme';
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

export interface CreatableSelectProps extends CreatableProps<
  Option,
  true,
  GroupBase<Option>
> {
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
  formatCreateLabel?: (v: string) => ReactNode;
  noOptionsMessage?: (props: NoOptionsMessageProps) => ReactNode;
}

const fontSize = '1rem';
const lineHeight = '1.25rem';
const minHeight = '3rem';

export const optionStyle = (
  styles: CSSObjectWithLabel,
  props: OptionProps<Option, boolean, GroupBase<Option>>,
  isHiddenCreateNewOption: boolean,
  isDark = false
) => {
  const { isFocused, data, isSelected } = props;
  const isNewOption = data?.__isNew__;

  return {
    ...styles,
    backgroundColor: isSelected
      ? isDark
        ? '#2B1F45'
        : '#E8E4F1'
      : isFocused
        ? isDark
          ? '#2B1F45'
          : '#FAFAFC'
        : isDark
          ? '#110D1C'
          : 'white',
    color: isDark ? '#F2EDF7' : '#475569',
    ':hover': {
      backgroundColor: `${isDark ? '#2B1F45' : '#FAFAFC'} !important`,
      cursor: 'pointer'
    },
    display: isNewOption && isHiddenCreateNewOption ? 'none' : 'flex'
  };
};

export const buildColorStyles = (
  isDark: boolean
): StylesConfig<Option, boolean> => ({
  control: (styles, { isDisabled }) => ({
    ...styles,
    backgroundColor: isDisabled
      ? isDark
        ? '#1B1725'
        : '#F3F4F6'
      : isDark
        ? '#110D1C'
        : 'white',
    borderColor: isDark ? '#2B1F45' : '#CBD5E1',
    '&:hover': {
      borderColor: isDark ? '#7B4FF5' : '#CBD5E1'
    },
    fontSize,
    lineHeight,
    minHeight,
    boxShadow: 'none !important',
    borderRadius: '8px',
    color: isDark ? '#F2EDF7' : '#475569'
  }),
  menu: base => ({
    ...base,
    fontSize,
    lineHeight,
    backgroundColor: isDark ? '#110D1C' : 'white',
    border: `1px solid ${isDark ? '#1B1725' : '#E2E8F0'}`,
    boxShadow: isDark
      ? '0px 4px 8px 1px rgba(0, 0, 0, 0.4)'
      : '0px 4px 8px rgba(35, 35, 35, 0.1)',
    color: isDark ? '#F2EDF7' : '#475569'
  }),
  menuList: base => ({
    ...base,
    backgroundColor: isDark ? '#110D1C' : 'white',
    padding: '4px'
  }),
  placeholder: styles => ({
    ...styles,
    color: `${isDark ? '#B5B0C2' : '#94A3B8'} !important`
  }),
  multiValue: (base, { isDisabled }) => ({
    ...base,
    color: isDisabled
      ? isDark
        ? '#7D768E'
        : '#6B7280'
      : isDark
        ? '#F2EDF7'
        : '#475569',
    backgroundColor: `${isDark ? '#2B1F45' : '#E8E4F1'} !important`,
    borderRadius: '4px'
  }),
  multiValueLabel: base => ({
    ...base,
    color: `${isDark ? '#B58CFF' : '#573792'} !important`
  }),
  multiValueRemove: base => ({
    ...base,
    color: `${isDark ? '#9A6FFF' : '#9A87BE'} !important`,
    ':hover': {
      color: `${isDark ? '#F2EDF7' : '#292C4C'} !important`,
      backgroundColor: `${isDark ? '#3D2A6B' : '#BCAFD3'} !important`
    }
  }),
  indicatorSeparator: base => ({
    ...base,
    display: 'none'
  }),
  input: base => ({
    ...base,
    color: isDark ? '#F2EDF7' : '#475569'
  }),
  singleValue: (styles, { isDisabled }) => ({
    ...styles,
    color: isDisabled
      ? isDark
        ? '#7D768E'
        : '#6B7280'
      : isDark
        ? '#F2EDF7'
        : '#475569'
  })
});

export const colorStyles = buildColorStyles(false);

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

export const CreatableSelect = memo<CreatableSelectProps>(
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
    const { theme } = useTheme();
    const isDark = theme === 'dark';
    const themedColorStyles = buildColorStyles(isDark);

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
            optionStyle(styles, props, isHiddenCreateNewOption, isDark),
          ...themedColorStyles,
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
