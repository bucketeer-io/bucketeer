import { ReactNode } from 'react';
import {
  IconLaunchOutlined,
  IconFilterListOutlined
} from 'react-icons-material-design';
import { Button } from 'components/button';
import Icon from 'components/icon';
import SearchInput from 'components/search-input';

type FilterProps = { additionalActions?: ReactNode };

const Filter = ({ additionalActions }: FilterProps) => {
  return (
    <div className="flex lg:items-center justify-between flex-col lg:flex-row">
      <div className="w-full lg:w-[365px]">
        <SearchInput placeholder="Search Input" value="" onChange={() => {}} />
      </div>
      <div className="flex items-center gap-4 mt-3 lg:mt-0">
        <Button variant="text" className="flex-1 lg:flex-none">
          <Icon icon={IconLaunchOutlined} size="sm" />
          Documentation
        </Button>
        <Button
          variant="secondary"
          className="text-gray-600 shadow-border-gray-400 flex-1 lg:flex-none"
        >
          <Icon icon={IconFilterListOutlined} size="sm" />
          Filter
        </Button>
        {additionalActions}
      </div>
    </div>
  );
};

export default Filter;
