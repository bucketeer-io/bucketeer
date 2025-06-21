import { unwrapUndefinable } from 'option-t/undefinable';
import { getCurrentEnvIdStorage } from 'storage/environment';
import { ConsoleAccount, Environment, EnvironmentRole, Project } from '@types';
import { isNotEmpty } from 'utils/data-type';

export const currentEnvironmentRole = (
  account: ConsoleAccount
): EnvironmentRole => {
  const currentEnvId = getCurrentEnvIdStorage();
  const curEnvId = isNotEmpty(currentEnvId)
    ? currentEnvId
    : account.environmentRoles[0].environment.id;

  let curEnvRole = account.environmentRoles.find(environmentRole => {
    const { environment } = environmentRole || {};
    const environmentId = environment?.id || environment?.urlCode;
    return environmentId === curEnvId;
  });
  if (!curEnvRole) {
    curEnvRole = account.environmentRoles[0];
  }
  return curEnvRole;
};

export const getCurrentEnvironment = (account: ConsoleAccount): Environment => {
  const envRole = currentEnvironmentRole(account);

  return envRole.environment;
};

export const getCurrentProject = (
  roles: EnvironmentRole[],
  currentEnvId: string
) => {
  try {
    return unwrapUndefinable(
      roles.find(role => {
        const { environment } = role || {};
        const environmentId = environment?.id || environment?.urlCode;
        return environmentId === currentEnvId;
      })
    )?.project;
  } catch {
    return null;
  }
};

export const hasEditable = (account: ConsoleAccount): boolean => {
  if (account.isSystemAdmin) return true;

  const envRole = currentEnvironmentRole(account);
  return envRole.role === 'Environment_EDITOR';
};

export const getUniqueProjects = (roles: EnvironmentRole[]): Project[] => {
  const projectMap = new Map<string, Project>();

  roles.forEach(role => {
    if (!projectMap.has(role.project.id)) {
      projectMap.set(role.project.id, role.project);
    }
  });

  return Array.from(projectMap.values());
};

export const getEnvironmentsByProjectId = (
  roles: EnvironmentRole[],
  projectId: string
): Environment[] => {
  return roles
    .filter(role => role.environment.projectId === projectId)
    .map(role => role.environment);
};

export const getEditorEnvironmentIds = (account: ConsoleAccount): string[] => {
  const ids = account.environmentRoles
    .filter(item => item.role === 'Environment_EDITOR')
    .map(item => item.environment.id);
  return ids;
};
