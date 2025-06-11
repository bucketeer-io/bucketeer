import { useCallback, useEffect, useState } from 'react';
import { IconCloseRound } from 'react-icons-material-design';
import { useNavigate } from 'react-router-dom';
import * as Popover from '@radix-ui/react-popover';
import {
  getCurrentEnvironment,
  getEnvironmentsByProjectId,
  getCurrentProject,
  getUniqueProjects,
  useAuth
} from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import {
  clearCurrentEnvIdStorage,
  getCurrentEnvIdStorage,
  setCurrentEnvIdStorage
} from 'storage/environment';
import { clearOrgIdStorage } from 'storage/organization';
import { Environment, Project } from '@types';
import { cn } from 'utils/style';
import { IconChecked, IconChevronRight, IconFolder, IconNoData } from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import List from 'components/list';
import { ScrollArea } from 'components/scroll-area';
import SearchInput from 'components/search-input';

const MyProjects = () => {
  const { t } = useTranslation(['common']);
  const navigate = useNavigate();
  const { consoleAccount, logout } = useAuth();
  const { errorNotify } = useToast();
  const [isShowProjectsList, setIsShowProjectsList] = useState(false);
  const [searchValue, setSearchValue] = useState('');
  const [projects, setProjects] = useState<Project[]>();
  const [selectedProject, setSelectedProject] = useState<Project>();
  const [selectedEnvironment, setSelectedEnvironment] = useState<Environment>();
  const [environments, setEnvironments] = useState<Environment[]>();

  const handleInitLoadData = useCallback(() => {
    const lastEnvironmentId = getCurrentEnvIdStorage();
    handleChangeData(lastEnvironmentId as string);
  }, [consoleAccount]);

  const handleChangeData = useCallback(
    (lastEnvironmentId?: string) => {
      const { environmentRoles } = consoleAccount!;
      const currentEnvironment = getCurrentEnvironment(consoleAccount!);
      const { id, urlCode } = currentEnvironment || {};

      const environmentId = lastEnvironmentId || id || urlCode;
      const currentProjects = getUniqueProjects(environmentRoles);
      const currentProject = getCurrentProject(environmentRoles, environmentId);

      if (!currentProject) {
        clearCurrentEnvIdStorage();
        clearOrgIdStorage();
        errorNotify(null, 'The environment is not found.');
        return logout();
      }
      const currentEnvironments = getEnvironmentsByProjectId(
        environmentRoles,
        currentProject.id
      );

      setCurrentEnvIdStorage(environmentId);
      setProjects(currentProjects);
      setSelectedProject(currentProject);
      setSelectedEnvironment(currentEnvironment);
      setEnvironments(currentEnvironments);
    },
    [consoleAccount]
  );

  const onOpenChange = useCallback(
    (v: boolean) => {
      if (!v) onClearSearch();
      setIsShowProjectsList(v);
    },
    [consoleAccount]
  );

  const onClearSearch = useCallback(() => {
    setSearchValue('');
    handleChangeData();
  }, [consoleAccount]);

  const onSearchProject = (value: string) => {
    if (!value) {
      onClearSearch();
    } else {
      const regex = new RegExp(value, 'i');
      const projectFiltered =
        projects?.filter(item => regex.test(item.name)) || [];
      if (projectFiltered.length > 0) {
        setSelectedProject(projectFiltered[0]);
      }
      setProjects(projectFiltered);
      setSearchValue(value);
    }
  };

  const onHandleChange = useCallback(
    (value: Environment) => {
      setSelectedEnvironment(value);
      setCurrentEnvIdStorage(value.id || value.urlCode);
      navigate(`/${value.urlCode}${PAGE_PATH_FEATURES}`);
      setIsShowProjectsList(false);
      onClearSearch();
    },
    [setSelectedEnvironment, consoleAccount]
  );

  const onChangeProject = useCallback(
    (project: Project) => {
      const { environmentRoles } = consoleAccount!;

      const currentEnvironments = getEnvironmentsByProjectId(
        environmentRoles,
        project.id
      );

      const environmentFiltered = currentEnvironments?.filter(
        item => item.projectId === project.id
      );
      setEnvironments(environmentFiltered);
      setSelectedProject(project);
    },
    [environments]
  );

  useEffect(() => {
    if (consoleAccount) {
      const timerId = setTimeout(() => handleInitLoadData(), 50);
      return () => clearTimeout(timerId);
    }
  }, [consoleAccount]);

  return (
    <Popover.Root onOpenChange={onOpenChange} open={isShowProjectsList}>
      <Popover.Portal>
        <Popover.Content align="start" className="border-none mt-2 z-20">
          <div className="w-[600px] bg-white rounded-lg shadow-menu">
            <div className="flex items-center justify-between px-5 py-4">
              <h1 className="typo-head-bold-huge text-gray-900 capitalize">
                {t(`navigation.my-projects`)}
              </h1>
              <Popover.Close>
                <Icon icon={IconCloseRound} size="sm" color="gray-500" />
              </Popover.Close>
            </div>
            <Divider />
            <div className="p-5">
              <SearchInput
                placeholder={t(`search`)}
                value={searchValue}
                onChange={onSearchProject}
              />
              {projects && projects?.length > 0 ? (
                <div className="mt-5 grid grid-cols-2 gap-4">
                  <div className="flex flex-col gap-3">
                    <List.Title>{t(`projects`)}</List.Title>
                    <ScrollArea className="h-[120px] pr-2">
                      <List
                        items={
                          projects?.map(item => ({
                            label: item.name,
                            value: item.id,
                            selected: item.id === selectedProject?.id,
                            expanded: item.id === selectedProject?.id,
                            onSelect: () => onChangeProject(item)
                          })) || []
                        }
                      />
                    </ScrollArea>
                  </div>
                  <div className="flex flex-col gap-3">
                    <List.Title>{t(`environments`)}</List.Title>
                    <ScrollArea className="h-[120px] pr-2">
                      <List
                        items={
                          environments?.map(item => ({
                            label: item.name,
                            value: item.id,
                            selected: item.id === selectedEnvironment?.id,
                            icon: IconChecked,
                            onSelect: () => onHandleChange(item)
                          })) || []
                        }
                      />
                    </ScrollArea>
                  </div>
                </div>
              ) : (
                <div className="flex flex-col justify-center items-center gap-3 pt-10 pb-4">
                  <IconNoData />
                  <div className="typo-para-medium text-gray-500">
                    {t(`navigation.no-projects`)}
                  </div>
                </div>
              )}
            </div>
          </div>
        </Popover.Content>
      </Popover.Portal>
      <Popover.Trigger className="w-full">
        <div
          className={cn(
            'flex items-center w-full text-primary-50 hover:bg-primary-400',
            'pl-3 pr-1.5 py-3 rounded-lg typo-para-medium justify-between',
            { 'bg-primary-400': isShowProjectsList }
          )}
        >
          <div className="flex items-center gap-x-2 truncate">
            <Icon color="primary-50" icon={IconFolder} size="sm" />
            <span className="truncate text-ellipsis">
              {selectedEnvironment?.name}
            </span>
          </div>
          <Icon color="primary-50" size="sm" icon={IconChevronRight} />
        </div>
      </Popover.Trigger>
    </Popover.Root>
  );
};

export default MyProjects;
