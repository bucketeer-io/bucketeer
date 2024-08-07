import { useCallback, useMemo, useState } from 'react';
import { IconCloseRound } from 'react-icons-material-design';
import * as Popover from '@radix-ui/react-popover';
import { Button } from 'components/button';
import Divider from 'components/divider';
import Icon from 'components/icon';
import List from 'components/list';
import { ListItemProps } from 'components/list/list-item';
import Search from 'components/search';

type ListOption = ListItemProps & {
  id: string;
};

const ProjectList = ({
  isOpen,
  onClose
}: {
  isOpen: boolean;
  onClose: () => void;
}) => {
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
    if (v === false) onClose();
  }, []);

  return (
    <Popover.Root open={isOpen} onOpenChange={onOpenChange}>
      <Popover.Content align="start" className="border-none mt-2">
        <div className="w-[600px] bg-white rounded-lg">
          <div className="flex items-center justify-between px-5 py-3.5">
            <h1 className="typo-head-bold-huge text-gray-900">{`My Projects`}</h1>
            <Button size="icon-sm" variant="grey" onClick={onClose}>
              <Icon icon={IconCloseRound} size="sm" />
            </Button>
          </div>
          <Divider />
          <div className="p-5">
            <Search size={'3'} className="h-12 rounded-lg" />
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
      <Popover.Trigger />
    </Popover.Root>
  );
};

export default ProjectList;
