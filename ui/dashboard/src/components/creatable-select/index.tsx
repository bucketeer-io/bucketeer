import {
  createContext,
  memo,
  ReactNode,
  useContext,
  useLayoutEffect,
  useRef,
  useState
} from 'react';
import {
  StylesConfig,
  ActionMeta,
  MultiValue,
  GroupBase,
  components,
  MenuListProps,
  MultiValueProps,
  CSSObjectWithLabel,
  OptionProps
} from 'react-select';
import ReactCreatableSelect, { CreatableProps } from 'react-select/creatable';
import {
  OVERFLOW_BADGE_WIDTH,
  useOverflowCount
} from 'hooks/use-overflow-count';
import OverflowValuesPopover from 'components/overflow-values-popover';
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
  // Keep values on one line when unfocused and show "+n" for the ones that don't fit
  collapseOverflow?: boolean;
  // Text below the input, shown when values are hidden, that opens the full list
  overflowSummary?: (count: number) => ReactNode;
  onChange: (
    options: MultiValue<Option>,
    actionMeta: ActionMeta<Option>
  ) => void;
  onCreateOption?: (v: string) => void;
  formatCreateLabel?: (v: string) => ReactNode;
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
  const { isFocused, data, isSelected } = props;

  const isNewOption = data?.__isNew__;

  return {
    ...styles,
    backgroundColor: isFocused && !isSelected ? backgroundColor : undefined,
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

// Space kept free for the inline search input
const INPUT_RESERVE_WIDTH = 8;
const VALUE_CONTAINER_SELECTOR = '.react-select__value-container';

const OverflowContext = createContext<{
  visibleCount: number;
  hiddenLabels: string[];
  // Removing is done from the overflow popup instead
  hideRemove: boolean;
}>({ visibleCount: Infinity, hiddenLabels: [], hideRemove: false });

const NoRemove = () => null;

const CollapsibleMultiValue = (
  props: MultiValueProps<Option, true, GroupBase<Option>>
) => {
  const { visibleCount, hiddenLabels, hideRemove } =
    useContext(OverflowContext);
  if (props.index >= visibleCount) return null;

  return (
    <>
      <components.MultiValue
        {...props}
        components={
          hideRemove
            ? { ...props.components, Remove: NoRemove }
            : props.components
        }
      />
      {props.index === visibleCount - 1 && hiddenLabels.length > 0 && (
        <span
          title={hiddenLabels.join(', ')}
          className="flex-center shrink-0 h-[22px] px-1.5 mx-0.5 rounded bg-gray-200 typo-para-small text-gray-600"
        >
          +{hiddenLabels.length}
        </span>
      )}
    </>
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
    collapseOverflow = false,
    overflowSummary,
    ...props
  }) => {
    const [isFocused, setIsFocused] = useState(false);
    const wrapperRef = useRef<HTMLDivElement>(null);
    // Track the selection when uncontrolled so overflow measurement and
    // popover removals always see what is actually rendered
    const [uncontrolledValue, setUncontrolledValue] = useState<
      MultiValue<Option>
    >(defaultValues ?? []);
    const selectedValues = value ?? uncontrolledValue;
    const handleChange = (
      newValue: MultiValue<Option>,
      actionMeta: ActionMeta<Option>
    ) => {
      if (!value) setUncontrolledValue(newValue);
      onChange(newValue, actionMeta);
    };
    const labels = selectedValues.map(item => item.label);
    const getValueContainer = () =>
      wrapperRef.current?.querySelector(VALUE_CONTAINER_SELECTOR);
    const { visibleCount, hiddenLabels } = useOverflowCount({
      labels,
      getContainer: getValueContainer,
      itemSelector: '.react-select__multi-value',
      enabled: collapseOverflow,
      expanded: isFocused,
      reservedWidth: INPUT_RESERVE_WIDTH
    });
    const labelsKey = labels.join('\u0000');

    // Keep the input in view while editing on a single scrollable line
    const scrollToEnd = () => {
      const valueContainer = getValueContainer();
      if (valueContainer)
        valueContainer.scrollLeft = valueContainer.scrollWidth;
    };
    useLayoutEffect(() => {
      if (collapseOverflow && isFocused) scrollToEnd();
    }, [collapseOverflow, isFocused, labelsKey]);

    const isCollapsed = visibleCount !== Infinity;

    const select = (
      <ReactCreatableSelect<Option, true>
        {...props}
        onFocus={e => {
          setIsFocused(true);
          props.onFocus?.(e);
        }}
        onBlur={e => {
          setIsFocused(false);
          props.onBlur?.(e);
        }}
        onInputChange={(newValue, actionMeta) => {
          if (collapseOverflow) requestAnimationFrame(scrollToEnd);
          return props.onInputChange?.(newValue, actionMeta);
        }}
        isMulti={isMulti}
        options={options}
        placeholder={placeholder}
        className={className}
        classNamePrefix="react-select"
        styles={{
          option: (styles, props) =>
            optionStyle(styles, props, isHiddenCreateNewOption),
          ...colorStyles,
          ...(collapseOverflow && {
            // Always a single line: collapsed to "+n" when blurred,
            // horizontally scrollable while focused
            valueContainer: base => ({
              ...base,
              flexWrap: 'nowrap',
              overflowX: isFocused ? 'auto' : 'hidden',
              scrollbarWidth: 'none',
              '&::-webkit-scrollbar': { display: 'none' },
              '& .react-select__multi-value': {
                flexShrink: 0,
                maxWidth: isCollapsed
                  ? `calc(100% - ${OVERFLOW_BADGE_WIDTH}px)`
                  : '100%'
              }
            })
          }),
          ...styles
        }}
        components={{
          ...components,
          ...(collapseOverflow && { MultiValue: CollapsibleMultiValue }),
          DropdownIndicator: null,
          LoadingIndicator: () => <Spinner className="size-5 mr-4" />
        }}
        isDisabled={disabled}
        isLoading={loading}
        isSearchable={isSearchable}
        value={selectedValues}
        onCreateOption={onCreateOption}
        onChange={(newValue, actionMeta) =>
          handleChange(newValue as MultiValue<Option>, actionMeta)
        }
        closeMenuOnSelect={closeMenuOnSelect}
        formatCreateLabel={formatCreateLabel}
        noOptionsMessage={noOptionsMessage}
      />
    );

    if (!collapseOverflow) return select;

    return (
      <OverflowContext.Provider
        value={{
          visibleCount,
          hiddenLabels,
          hideRemove: !!overflowSummary && isCollapsed
        }}
      >
        <div ref={wrapperRef} className="w-full min-w-0">
          {select}
          {overflowSummary && hiddenLabels.length > 0 && (
            <OverflowValuesPopover
              labels={labels}
              text={overflowSummary(labels.length)}
              onRemove={
                disabled
                  ? undefined
                  : index => {
                      const removedValue = selectedValues[index];
                      handleChange(
                        selectedValues.filter((_, i) => i !== index),
                        { action: 'remove-value', removedValue }
                      );
                    }
              }
            />
          )}
        </div>
      </OverflowContext.Provider>
    );
  }
);
