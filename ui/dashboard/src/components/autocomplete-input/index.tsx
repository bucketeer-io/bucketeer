import React, {
  useState,
  ChangeEvent,
  InputHTMLAttributes,
  forwardRef,
  Ref,
  useRef,
  useImperativeHandle
} from 'react';
import { cn } from 'utils/style';

type AutocompleteInputProps = {
  options: string[];
} & InputHTMLAttributes<HTMLInputElement>;

const AutocompleteInput: React.FC<AutocompleteInputProps> = forwardRef(
  ({ options, className, value, onChange }, ref: Ref<HTMLInputElement>) => {
    const isControlled = value !== undefined;
    const [uncontrolledValue, setUncontrolledValue] = useState('');
    const inputValue = isControlled ? (value as string) : uncontrolledValue;

    const [suggestions, setSuggestions] = useState<string[]>([]);
    const [showSuggestions, setShowSuggestions] = useState(false);

    const inputRef = useRef<HTMLInputElement>(null);
    useImperativeHandle(ref, () => inputRef.current as HTMLInputElement);

    const filterSuggestions = (input: string) => {
      return input.trim() === ''
        ? options
        : options.filter(option =>
            option.toLowerCase().startsWith(input.toLowerCase())
          );
    };

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value;

      if (!isControlled) {
        setUncontrolledValue(newValue);
      }

      if (onChange) {
        onChange(e);
      }

      setSuggestions(filterSuggestions(newValue));
      setShowSuggestions(true);
    };

    const handleFocus = () => {
      setSuggestions(filterSuggestions(inputValue));
      setShowSuggestions(true);
    };

    const handleSelect = (selected: string) => {
      if (!isControlled) {
        setUncontrolledValue(selected);
      }

      if (inputRef.current) {
        const nativeInput = inputRef.current;
        nativeInput.value = selected;

        if (onChange) {
          const event = {
            target: nativeInput
          } as ChangeEvent<HTMLInputElement>;

          onChange(event);
        }
      }

      setSuggestions([]);
      setShowSuggestions(false);
    };

    const handleBlur = () => setShowSuggestions(false);

    return (
      <div className={cn('relative w-full', className)}>
        <input
          ref={inputRef}
          type="text"
          value={inputValue}
          onChange={handleChange}
          onFocus={handleFocus}
          onBlur={handleBlur}
          className="border border-gray-400 p-2 typo-para-medium w-full px-4 py-[11px] text-gray-700"
        />
        {showSuggestions && suggestions.length > 0 && (
          <ul className="absolute p-1 max-h-[252px] bg-white border rounded-lg w-full shadow-dropdown overflow-x-hidden overflow-y-auto small-scroll mt-1 z-10">
            {suggestions.map((suggestion, index) => (
              <li
                key={index}
                onMouseDown={() => handleSelect(suggestion)}
                className="p-2 hover:bg-gray-100 cursor-pointer rounded-md typo-para-medium text-gray-700"
              >
                {suggestion}
              </li>
            ))}
          </ul>
        )}
      </div>
    );
  }
);

export default AutocompleteInput;
