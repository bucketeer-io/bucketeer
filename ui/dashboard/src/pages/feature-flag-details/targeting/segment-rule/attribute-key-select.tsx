import { components, GroupBase } from 'react-select';
import ReactCreatableSelect from 'react-select/creatable';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconChecked } from '@icons';
import { colorStyles, Option, optionStyle } from 'components/creatable-select';
import { UserMessage } from '../individual-rule';

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
            className="!max-h-[250px] !-mt-1 overflow-x-hidden overflow-y-auto small-scroll"
          />
        ),
        Option: props => (
          <components.Option
            {...props}
            className={cn(
              'flex items-center justify-between w-full px-3 py-1.5 mb-0.5',
              props.isSelected && 'bg-gray-100'
            )}
          >
            <span>{props.data.label}</span>
            {props.isSelected && (
              <IconChecked className="text-primary-500 w-6" />
            )}
          </components.Option>
        ),
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
        ...colorStyles
      }}
      value={value}
      onChange={option => {
        const newValue = option as Option;
        onChange(newValue.value);
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
