import { useCallback, useState } from 'react';
import { IconCloseRound } from 'react-icons-material-design';
import * as Popover from '@radix-ui/react-popover';
import {
  getCurrentEnvironment,
  getEnvironmentsByProjectId,
  getUniqueProjects,
  useAuth
} from 'auth';
import { cn } from 'utils/style';
import { IconChevronRight, IconFolder } from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import List from 'components/list';
import { ScrollArea } from 'components/scroll-area';
import SearchInput from 'components/search-input';

const MyProjects = () => {
  const { consoleAccount } = useAuth();
  const [isShowProjectsList, setIsShowProjectsList] = useState(false);
  const [searchValue, setSearchValue] = useState('');

  const envRoles = consoleAccount?.environmentRoles || [];
  const projects = getUniqueProjects(envRoles);
  const currentEnv = getCurrentEnvironment(consoleAccount!);

  const [selectedProjectId, setSelectedProjectId] = useState<string>(
    projects[0].id
  );
  const [selectedEnvId, setSelectedEnvId] = useState<string>(currentEnv.id);
  const environments = getEnvironmentsByProjectId(envRoles, selectedProjectId);

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
                <ScrollArea className="h-[120px] pr-2">
                  <List
                    items={projects.map(item => ({
                      label: item.name,
                      value: item.id,
                      selected: item.id === selectedProjectId,
                      expanded: item.id === selectedProjectId,
                      onSelect: setSelectedProjectId
                    }))}
                  />
                </ScrollArea>
              </div>

              <div className="flex flex-col gap-3">
                <List.Title>{`Environment`}</List.Title>
                <ScrollArea className="h-[120px] pr-2">
                  <List
                    items={environments.map(item => ({
                      label: item.name,
                      value: item.id,
                      selected: item.id === selectedEnvId,
                      onSelect: setSelectedEnvId
                    }))}
                  />
                </ScrollArea>
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

export default MyProjects;
