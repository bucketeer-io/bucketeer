import { memo, useRef } from 'react';
import { components, GroupBase, OptionProps } from 'react-select';
import ReactCreatableSelect from 'react-select/creatable';
import { useIsTruncated } from 'hooks/use-is-truncated';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconChecked } from '@icons';
import { colorStyles, Option, optionStyle } from 'components/creatable-select';
import { Tooltip } from 'components/tooltip';
import { UserMessage } from '../individual-rule';

const CustomOption = memo((props: OptionProps<Option>) => {
  const label = props.data.label;
  const spanRef = useRef<HTMLSpanElement>(null);
  const isTruncated = useIsTruncated(spanRef, [props.data.label]);

  const labelNode = (
    <span ref={spanRef} className="truncate block max-w-[350px]">
      {label}
    </span>
  );

  return (
    <components.Option
      {...props}
      className={cn(
        'flex items-center justify-between w-full gap-2 px-3 py-1.5 mb-0.5',
        props.isSelected && 'bg-gray-100'
      )}
    >
      {isTruncated ? (
        <Tooltip align="start" content={label} trigger={labelNode} />
      ) : (
        labelNode
      )}

      {props.isSelected && <IconChecked className="text-primary-500 w-6" />}
    </components.Option>
  );
});

const AttributeKeySelect = ({
  createdOptions,
  sdkOptions,
  value,
  onChange,
  onCreateOption
}: {
  createdOptions: Option[];
  sdkOptions: Option[];
  value: Option;
  onChange: (v: string) => void;
  onCreateOption: (v: string) => void;
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);

  return (
    <ReactCreatableSelect<Option, false, GroupBase<Option>>
      options={[
        { label: '', options: createdOptions },
        {
          label: t('form:feature-flags.attribute-key-select-title'),
          options: sdkOptions
        }
      ]}
      classNamePrefix="react-select"
      components={{
        MenuList: props => (
          <components.MenuList
            {...props}
            className="!max-h-[250px] overflow-x-hidden overflow-y-auto small-scroll"
          />
        ),
        Option: CustomOption,
        GroupHeading: props =>
          props.children && (
            <div
              className={cn(
                'typo-para-tiny text-gray-600 bg-gray-100 relative w-full px-2 py-2.5 mb-2'
              )}
            >
              {props.children}
            </div>
          )
      }}
      styles={{
        option: (styles, props) => optionStyle(styles, props, false),
        ...colorStyles,
        menu: base => ({
          ...base,
          width: 'auto',
          minWidth: base.width,
          maxWidth: '400px'
        }),
        menuPortal: base => ({
          ...base,
          zIndex: 9999
        })
      }}
      value={value}
      onChange={(option: Option | null) => {
        if (option) onChange(option.value);
      }}
      formatCreateLabel={value => (
        <p>{`${t('create-option', { option: value })}`}</p>
      )}
      onCreateOption={onCreateOption}
      noOptionsMessage={() => (
        <UserMessage message={t('no-opts-type-to-create')} />
      )}
      isValidNewOption={(inputValue, _value, options) => {
        if (!inputValue) return false;
        const existingOptions = options.flatMap(
          group => group.options as Option[]
        );

        const exists = existingOptions.some(
          opt => opt.label.toLowerCase() === inputValue.toLowerCase()
        );

        return !exists;
      }}
    />
  );
};

export default AttributeKeySelect;
