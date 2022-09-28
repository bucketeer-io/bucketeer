import { FC, memo } from 'react';

export interface CheckBoxProps {
  id: string;
  value: string;
  onChange: (value: string, checked: boolean) => void;
  defaultChecked?: boolean;
  disabled?: boolean;
}

export const CheckBox: FC<CheckBoxProps> = memo(
  ({ id, value, defaultChecked, onChange, disabled }) => {
    const handleOnChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      onChange(e.target.value, e.target.checked);
    };

    return (
      <input
        id={id}
        type="checkbox"
        defaultChecked={defaultChecked}
        onChange={handleOnChange}
        value={value}
        disabled={disabled}
        className="input-checkbox"
      />
    );
  }
);
