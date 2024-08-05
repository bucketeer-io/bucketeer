import { useMemo, useState } from 'react';
import { Popover, Separator } from '@radix-ui/themes';
import { IconCloseFilled } from '@icons';
import List from 'components/list';
import { ListItemProps } from 'components/list/list-item';
import Search from 'components/search';

type ListOption = ListItemProps & {
  id: string;
};

const MyProject = () => {
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

  return (
    <div className="w-[600px]">
      <div className="flex h-16 items-center justify-between p-5">
        <h1 className="typo-head-bold-huge text-gray-900">My Projects</h1>
        <Popover.Close>
          <button className="hover:cursor-pointer">
            <IconCloseFilled />
          </button>
        </Popover.Close>
      </div>
      <Separator size={'4'} color="gray" />
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
  );
};

export default MyProject;
