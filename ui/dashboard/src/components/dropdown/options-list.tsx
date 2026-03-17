import { memo, ReactNode } from 'react';
import { t } from 'i18next';
import { Tooltip } from 'components/tooltip';
import DropdownList from 'elements/dropdown-list';
import { DropdownMenuItem } from './item';
import { DropdownOption, DropdownValue } from './types';

interface OptionsListProps {
  options: DropdownOption[];
  value?: DropdownValue | DropdownValue[];
  multiselect: boolean;
  useVirtualList: boolean;
  isTooltip: boolean;
  itemClassName?: string;
  additionalElement?: (item: DropdownOption) => ReactNode;
  onChange?: (value: DropdownValue | DropdownValue[]) => void;
}

export const OptionsList = memo(
  ({
    options,
    value,
    multiselect,
    useVirtualList,
    isTooltip,
    itemClassName,
    additionalElement,
    onChange
  }: OptionsListProps) => {
    if (options.length === 0) {
      return (
        <div className="flex-center py-2.5 typo-para-medium text-gray-600">
          {t('no-options-found')}
        </div>
      );
    }

    if (useVirtualList) {
      return (
        <DropdownList
          isMultiselect={multiselect}
          itemSelected={value as string}
          selectedOptions={(Array.isArray(value) ? value : []) as string[]}
          additionalElement={additionalElement}
          options={options}
          onSelectOption={val => onChange?.(val)}
        />
      );
    }

    return (
      <>
        {options.map(opt =>
          isTooltip && opt.tooltip ? (
            <Tooltip
              key={opt.value}
              side="right"
              sideOffset={10}
              align="start"
              className="w-[180px] p-3 bg-white typo-para-small text-gray-600 shadow-card"
              content={opt.tooltip}
              showArrow={false}
              trigger={
                <DropdownMenuItem
                  icon={opt.icon}
                  label={opt.label}
                  value={opt.value}
                  disabled={opt.disabled}
                  onSelectOption={onChange ? val => onChange(val) : undefined}
                />
              }
            />
          ) : (
            <DropdownMenuItem
              key={opt.value}
              value={opt.value}
              label={opt.label}
              icon={opt.icon}
              description={opt.description}
              iconElement={opt.iconElement}
              additionalElement={
                additionalElement ? additionalElement(opt) : undefined
              }
              isSelectedItem={value === opt.value}
              onSelectOption={val => onChange?.(val)}
              className={itemClassName}
            />
          )
        )}
      </>
    );
  }
);
