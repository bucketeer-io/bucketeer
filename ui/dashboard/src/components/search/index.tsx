import { TextField } from '@radix-ui/themes';
import { Responsive } from '@radix-ui/themes/props';
import { IconSearch } from '@icons';

export type SearchProps = TextField.RootProps & {
  size?: Responsive<'1' | '2' | '3'> | undefined;
};

const Search = ({ size = '2', className, ...props }: SearchProps) => {
  return (
    <TextField.Root
      placeholder="Search"
      size={size}
      className={className}
      {...props}
    >
      <TextField.Slot>
        <IconSearch />
      </TextField.Slot>
    </TextField.Root>
  );
};

export default Search;
