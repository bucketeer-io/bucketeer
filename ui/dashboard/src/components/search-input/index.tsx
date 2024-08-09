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
  const searchAtionRef = useRef<HTMLDivElement>(null);
  const [searchValue, setSearchValue] = useState(defaultValue);
  const searchValueRef = useRef(searchValue);

  useEffect(() => {
    searchValueRef.current = searchValue;
  }, [searchValue]);

  useEffect(() => {
    setSearchValue(defaultValue);
  }, [defaultValue]);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (searchValueRef.current === '') {
        setSearchValue('');
        onChange('');
      } else if (
        searchAtionRef.current &&
        !searchAtionRef.current.contains(event.target as Node)
      ) {
        if (defaultValue !== searchValueRef.current) {
          setSearchValue(defaultValue);
        }
      }
    }

    document.addEventListener('mousedown', handleClickOutside);

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [defaultValue]);

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
