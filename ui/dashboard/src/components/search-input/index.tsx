import { useEffect, useRef, useState } from 'react';
import { IconSearch } from '@icons';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';

export interface SearchBarProps {
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
}

const SearchInput = ({
  placeholder,
  value: defaultValue,
  onChange,
  disabled
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
    <form
      className="w-full"
      onSubmit={event => {
        event.preventDefault();
        event.stopPropagation();
        onChange(searchValue);
      }}
    >
      <InputGroup
        className="w-full"
        addon={<Icon icon={IconSearch} size="sm" />}
      >
        <Input
          placeholder={placeholder}
          value={searchValue}
          onChange={setSearchValue}
          disabled={disabled}
        />
      </InputGroup>
    </form>
  );
};

export default SearchInput;
