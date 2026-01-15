import { useCallback } from 'react';
import { cn } from 'utils/style';

const AudienceSelect = ({
  label,
  value,
  isActive,
  disabled,
  onSelect
}: {
  label: string;
  value: string | number;
  isActive: boolean;
  disabled?: boolean;
  onSelect: (value: string | number) => void;
}) => {
  const handleClick = useCallback(() => {
    if (!disabled) {
      onSelect(value);
    }
  }, [disabled, onSelect, value]);
  return (
    <div
      className={cn(
        'flex-center size-fit min-w-20 py-[14px] px-3 border border-gray-400 rounded-lg typo-para-medium leading-5 text-gray-700 capitalize',
        disabled ? 'cursor-not-allowed bg-gray-100' : 'cursor-pointer',
        {
          'text-primary-500 border-primary-500': isActive && !disabled
        }
      )}
      onClick={handleClick}
    >
      {label}
    </div>
  );
};

export default AudienceSelect;
