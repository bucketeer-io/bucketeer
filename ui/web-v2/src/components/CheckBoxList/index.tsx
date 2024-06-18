import { FC, memo, useState } from 'react';

import { CheckBox } from '../CheckBox';

export interface Option {
  value: string;
  label: string;
  description?: string;
}

export interface CheckBoxListProps {
  onChange: (values: string[]) => void;
  options: Option[];
  defaultValues?: Option[];
  disabled?: boolean;
}

export const CheckBoxList: FC<CheckBoxListProps> = memo(
  ({ onChange, options, defaultValues, disabled }) => {
    const [checkedItems] = useState(() => {
      const items = new Map();
      defaultValues &&
        defaultValues.forEach((item, _) => {
          items.set(item.value, item.value);
        });
      return items;
    });

    const handleOnChange = (value: string, checked: boolean) => {
      if (checked) {
        checkedItems.set(value, value);
      } else {
        checkedItems.delete(value);
      }
      const valueList = [];
      checkedItems.forEach((v, _) => {
        valueList.push(v);
      });
      onChange(valueList);
    };

    return (
      <div>
        <fieldset className="border-t border-b border-gray-300">
          <div className="divide-y divide-gray-300">
            {options.map((item, index) => {
              return (
                <div
                  key={item.label}
                  className="relative flex items-start py-4"
                >
                  <div className="min-w-0 flex-1 text-sm">
                    <label htmlFor={`id_${index}`} key={`key_${index}`}>
                      <p className="text-sm font-medium text-gray-700">
                        {item.label}
                      </p>
                      {item.description && (
                        <p className="text-sm text-gray-500">
                          {item.description}
                        </p>
                      )}
                    </label>
                  </div>
                  <div className="ml-3 flex items-center h-5">
                    <CheckBox
                      id={`id_${index}`}
                      value={item.value}
                      onChange={handleOnChange}
                      defaultChecked={checkedItems.has(item.value)}
                      disabled={disabled}
                    />
                  </div>
                </div>
              );
            })}
          </div>
        </fieldset>
      </div>
    );
  }
);
