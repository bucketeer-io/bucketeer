import { useCallback, useMemo, useState } from 'react';
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

type ListOption = ListItemProps & {
  id: string;
};

const ProjectList = () => {
  const [isShowProjectsList, setIsShowProjectsList] = useState(false);
  const [selectedProject, setSelectedProject] = useState('');
  const [projects] = useState<ListOption[]>([
    {
      id: '1',
      text: 'Default Project',
      type: 'icon'
    },
    {
      id: '2',
      text: 'Bucketeer demo application'
    },
    {
      id: '3',
      text: 'Yuichi'
    }
  ]);

  const handleSelectedProject = (value: string) => {
    setSelectedProject(value);
  };

  const menuProjects: ListItemProps[] = useMemo(() => {
    return projects.map(i => {
      i.selected = selectedProject === i.id;
      i.onClick = () => handleSelectedProject(i.id);
      return i;
    });
  }, [projects, selectedProject]);

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
            <InputGroup
              className="w-full"
              addon={<Icon icon={IconSearch} size="sm" />}
            >
              <Input placeholder="Search" className="w-full" />
            </InputGroup>
            <div className="mt-5 grid grid-cols-2 gap-4">
              <List title="Projects" options={menuProjects} />
              {selectedProject && (
                <List
                  title="Environment"
                  options={[
                    {
                      text: 'Test',
                      type: 'icon'
                    },
                    {
                      text: 'Production'
                    }
                  ]}
                />
              )}
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
