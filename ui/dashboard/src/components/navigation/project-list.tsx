import { useCallback, useState } from 'react';
import { IconCloseRound } from 'react-icons-material-design';
import * as Popover from '@radix-ui/react-popover';
import { cn } from 'utils/style';
import { IconChevronRight, IconFolder, IconSearch } from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import List from 'components/list';
import { ListItemProps } from 'components/list/list-item';
import SearchInput from 'components/search-input';

const ProjectList = () => {
  const [isShowProjectsList, setIsShowProjectsList] = useState(false);
  const [searchValue, setSearchValue] = useState('');
  const [projects] = useState<ListItemProps[]>([
    {
      label: 'Default Project',
      icon: IconChevronRight,
      selected: true
    },
    {
      label: 'Bucketeer demo application'
    },
    {
      label: 'Yuichi'
    }
  ]);
  const [environments] = useState<ListItemProps[]>([
    {
      label: 'Test',
      selected: true
    },
    {
      label: 'Production'
    }
  ]);

  const onOpenChange = useCallback((v: boolean) => {
    setIsShowProjectsList(v);
  }, []);

  return (
    <Popover.Root onOpenChange={onOpenChange}>
      <Popover.Content align="start" className="border-none mt-2">
        <div className="w-[600px] bg-white rounded-lg shadow">
          <div className="flex items-center justify-between px-5 py-4">
            <h1 className="typo-head-bold-huge text-gray-900">
              {`My Projects`}
            </h1>
            <Popover.Close>
              <Icon icon={IconCloseRound} size="sm" color="gray-500" />
            </Popover.Close>
          </div>
          <Divider />
          <div className="p-5">
            <SearchInput
              placeholder="Search"
              value={searchValue}
              onChange={setSearchValue}
            />
            <div className="mt-5 grid grid-cols-2 gap-4">
              <div className="flex flex-col gap-3">
                <List.Title>{`Projects`}</List.Title>
                <List items={projects} />
              </div>
              <div className="flex flex-col gap-3">
                <List.Title>{`Environment`}</List.Title>
                <List items={environments} />
              </div>
            </div>
          </div>
        </div>
      </Popover.Content>
      <Popover.Trigger>
        <div
          className={cn(
            'flex items-center gap-x-2 w-full text-primary-50',
            'px-3 py-3 rounded-lg typo-para-medium justify-between',
            'hover:bg-primary-400 hover:opacity-80',
            { 'bg-primary-400': isShowProjectsList }
          )}
        >
          <div className="flex items-center gap-x-2">
            <Icon color="primary-50" icon={IconFolder} size="sm" />
            {`Abematv`}
          </div>
          <Icon color="primary-50" icon={IconChevronRight} />
        </div>
      </Popover.Trigger>
    </Popover.Root>
  );
};

export default ProjectList;
