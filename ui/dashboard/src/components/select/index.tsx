export type SelectProps = React.DetailedHTMLProps<
  React.SelectHTMLAttributes<HTMLSelectElement>,
  HTMLSelectElement
> & {
  label: string;
  options: { value: string; label: string }[];
  required?: boolean;
  placeholder?: string;
};

const Select = ({
  id,
  label,
  options = [],
  required,
  placeholder,
  ...props
}: SelectProps) => {
  return (
    <div className="grid gap-[6px] mt-2">
      <label htmlFor={id} className="text-gray-600 typo-para-small">
        {label}{' '}
        {required && (
          <span className="text-accent-red-500 typo-para-small">*</span>
        )}
      </label>
      <select
        {...props}
        id={id}
        className="p-3 rounded-lg text-gray-500 typo-para-medium border-gray-400"
      >
        <option value="" disabled selected>
          {placeholder}
        </option>
        {options.map(i => (
          <option key={i.value} value={i.value}>
            {i.label}
          </option>
        ))}
      </select>
    </div>
  );
};

export default Select;