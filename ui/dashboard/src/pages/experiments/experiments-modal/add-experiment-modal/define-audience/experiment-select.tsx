import { cn } from 'utils/style';

const ExperimentSelect = ({
  label,
  value,
  isActive,
  onSelect
}: {
  label: string;
  value: string | number;
  isActive: boolean;
  onSelect: (value: string | number) => void;
}) => {
  return (
    <div
      className={cn(
        'flex-center size-fit min-w-[53px] py-[14px] px-3 border border-gray-400 rounded-lg typo-para-medium leading-5 text-gray-700 capitalize cursor-pointer',
        {
          'text-primary-500 border-primary-500': isActive
        }
      )}
      onClick={() => onSelect(value)}
    >
      {label}
    </div>
  );
};

export default ExperimentSelect;
