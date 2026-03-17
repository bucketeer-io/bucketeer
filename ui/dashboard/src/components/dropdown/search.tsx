import { forwardRef, Ref } from 'react';
import { cn } from 'utils/style';
import { IconSearch } from '@icons';
import Icon from 'components/icon';
import Input, { InputProps } from 'components/input';

export type DropdownMenuSearchProps = InputProps;

export const DropdownMenuSearch = forwardRef(
  (
    { value, onChange, className, ...props }: DropdownMenuSearchProps,
    ref: Ref<HTMLInputElement>
  ) => (
    <div className="sticky top-0 flex items-center w-full px-3 py-[11.5px] gap-x-2 border-b border-gray-200 bg-white z-50">
      <div className="flex-center size-5">
        <Icon icon={IconSearch} size="xs" color="gray-500" />
      </div>
      <Input
        {...props}
        ref={ref}
        value={value}
        onChange={onChange}
        className={cn('p-0 border-none rounded-none', className)}
      />
    </div>
  )
);
