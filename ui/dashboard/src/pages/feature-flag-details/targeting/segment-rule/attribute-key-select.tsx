import React from 'react';
import { components, MenuListProps } from 'react-select';
import ReactCreatableSelect from 'react-select/creatable';
import { useTranslation } from 'i18n';
import { t } from 'i18next';
import { cn } from 'utils/style';
import { IconChecked } from '@icons';
import { colorStyles, Option, optionStyle } from 'components/creatable-select';
import Divider from 'components/divider';
import { UserMessage } from '../individual-rule';

const CustomMenuList = (props: MenuListProps<Option, false>) => {
  const { getValue, children } = props;
  const selected = getValue()[0] as
    | (Option & { __isNew__?: boolean })
    | undefined;

  let selectedChild: React.ReactNode = null;
  const otherChildren: React.ReactNode[] = [];
  let hasCreateOption = false;

  React.Children.forEach(children, child => {
    if (!React.isValidElement(child)) return;

    const data = (child.props as { data?: Option & { __isNew__?: boolean } })
      .data;

    if (data?.__isNew__) {
      hasCreateOption = true;
    }

    if (selected && data?.value === selected.value) {
      selectedChild = child;
    } else {
      otherChildren.push(child);
    }
  });

  return (
    <components.MenuList {...props}>
      {selectedChild && !selected?.__isNew__ && (
        <div className="flex items-center justify-between w-full px-3 py-1.5">
          <span>{selected?.label}</span>
          <IconChecked className="text-primary-500 w-6" />
        </div>
      )}
      {selectedChild && !selected?.__isNew__ && otherChildren.length > 0 && (
        <Divider className="mt-0.5" />
      )}

      {!hasCreateOption && otherChildren.length > 0 && (
        <>
          <div
            className={cn(
              'typo-para-tiny text-gray-600 bg-gray-100 relative w-full px-3 py-2.5',
              { '-mt-1 rounded-t-md': !selectedChild }
            )}
          >
            {t('form:feature-flags.attribute-key-select-title')}
          </div>
          <Divider className="mb-0.5" />
        </>
      )}

      {otherChildren}
    </components.MenuList>
  );
};

const AttributeKeySelect = ({
  options,
  value,
  onChange
}: {
  options: Option[];
  value: Option;
  onChange: (v: string) => void;
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);

  return (
    <ReactCreatableSelect<Option, false>
      options={options}
      classNamePrefix="react-select"
      components={{ MenuList: CustomMenuList }}
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
        <p>
          {`${t('create-option', {
            option: value
          })}`}
        </p>
      )}
      noOptionsMessage={() => (
        <UserMessage message={t('no-opts-type-to-create')} />
      )}
    />
  );
};

export default AttributeKeySelect;
