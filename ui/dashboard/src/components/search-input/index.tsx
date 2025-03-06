import { KeyboardEvent, useEffect, useRef, useState } from 'react';
import { IconSearch } from '@icons';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';

export interface SearchBarProps {
  placeholder: string;
  value: string;
  disabled?: boolean;
  variant?: 'primary' | 'secondary';
  onChange: (value: string) => void;
  onKeyDown?: (e: KeyboardEvent<HTMLInputElement>) => void;
}

const SearchInput = ({
  placeholder,
  value: defaultValue,
  disabled,
  variant = 'primary',
  onChange,
  onKeyDown
}: SearchBarProps) => {
  const [searchValue, setSearchValue] = useState(defaultValue);
  const searchValueRef = useRef(false);

  useEffect(() => {
    setSearchValue(defaultValue);
  }, [defaultValue]);

  useEffect(() => {
    if (searchValueRef.current) {
      const timeout = setTimeout(() => onChange(searchValue), 500);
      return () => {
        clearTimeout(timeout);
      };
    }
  }, [searchValue]);

  useEffect(() => {
    searchValueRef.current = true;
  }, []);

  return (
    <fieldset
      className="w-full"
      onSubmit={event => {
        event.preventDefault();
        event.stopPropagation();
        onChange(searchValue);
      }}
    >
      <InputGroup
        className="w-full"
        addon={
          <Icon
            icon={IconSearch}
            size="sm"
            color={variant === 'primary' ? 'gray-500' : 'primary-500'}
          />
        }
      >
        <Input
          variant={variant}
          placeholder={placeholder}
          value={searchValue}
          disabled={disabled}
          onChange={setSearchValue}
          onKeyDown={e => onKeyDown && onKeyDown(e)}
        />
      </InputGroup>
    </fieldset>
  );
};

export default SearchInput;
