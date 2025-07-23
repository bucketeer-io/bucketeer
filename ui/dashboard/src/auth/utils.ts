import { isNil } from 'lodash';
import { unwrapUndefinable } from 'option-t/undefinable';
import { getCurrentEnvIdStorage } from 'storage/environment';
import { getCurrentProjectEnvironmentStorage } from 'storage/project-environment';
import { ConsoleAccount, Environment, EnvironmentRole, Project } from '@types';
import { checkEnvironmentEmptyId } from 'utils/function';

export const currentEnvironmentRole = (
  account: ConsoleAccount
): EnvironmentRole => {
  const currentEnvId = getCurrentEnvIdStorage();
  const projectEnvironment = getCurrentProjectEnvironmentStorage();
  const curEnvId = !isNil(currentEnvId)
    ? currentEnvId
    : account.environmentRoles[0].environment.id;

  let curEnvRole = undefined;
  const checkEmptyEnvironmentId = checkEnvironmentEmptyId(curEnvId);
  const checkEmptyProjectEnvironmentId =
    projectEnvironment &&
    checkEnvironmentEmptyId(projectEnvironment?.environmentId);
  if (checkEmptyEnvironmentId === checkEmptyProjectEnvironmentId) {
    curEnvRole = account.environmentRoles.find(environmentRole => {
      const { environment, project } = environmentRole || {};
      return (
        project.id === projectEnvironment?.projectId &&
        environment.id === checkEmptyEnvironmentId
      );
    });
  } else {
    curEnvRole = account.environmentRoles.find(environmentRole => {
      const { environment } = environmentRole || {};
      const environmentId = environment?.id;
      return environmentId === checkEmptyEnvironmentId;
    });
  }
  if (!curEnvRole) {
    curEnvRole = account.environmentRoles[0];
  }

  curEnvRole.environment = {
    ...curEnvRole.environment,
    id: checkEnvironmentEmptyId(curEnvRole.environment.id)
  };
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
    const projectEnvironment = getCurrentProjectEnvironmentStorage();
    const checkEmptyEnvironmentId = checkEnvironmentEmptyId(currentEnvId);
    const checkEmptyProjectEnvironmentId =
      projectEnvironment &&
      checkEnvironmentEmptyId(projectEnvironment?.environmentId);
    if (checkEmptyEnvironmentId === checkEmptyProjectEnvironmentId) {
      return unwrapUndefinable(
        roles.find(item => item.project.id === projectEnvironment?.projectId)
      )?.project;
    }
    return unwrapUndefinable(
      roles.find(role => {
        const { environment } = role || {};
        const environmentId = environment?.id;
        return environmentId === checkEmptyEnvironmentId;
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

export const getEditorEnvironments = (account: ConsoleAccount) => {
  const environmentsEditorRole = account.environmentRoles.filter(
    item => item.role === 'Environment_EDITOR'
  );
  const editorEnvironments = environmentsEditorRole.map(
    item => item.environment
  );
  const editorEnvironmentIDs = editorEnvironments.map(item => item.id);
  const projects = environmentsEditorRole.map(item => item.project);

  return {
    editorEnvironments,
    editorEnvironmentIDs,
    projects
  };
};

export const getAccountAccess = (account: ConsoleAccount) => {
  const envEditable = hasEditable(account);
  const isOrganizationAdmin = [
    'Organization_OWNER',
    'Organization_ADMIN'
  ].includes(account?.organizationRole);
  return { envEditable, isOrganizationAdmin };
};
