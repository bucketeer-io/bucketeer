import { getCurrentEnvIdStorage } from 'storage/environment';
import { ConsoleAccount, Environment, EnvironmentRole } from '@types';

export const currentEnvironmentRole = (
  consoleAccount: ConsoleAccount
): EnvironmentRole => {
  const currentEnvId = getCurrentEnvIdStorage();
  const curEnvId =
    currentEnvId != null && currentEnvId != undefined
      ? currentEnvId
      : consoleAccount.environmentRoles[0].environment.id;
  let curEnvRole = consoleAccount.environmentRoles.find(
    environmentRole => environmentRole.environment.id === curEnvId
  );
  if (!curEnvRole) {
    curEnvRole = consoleAccount.environmentRoles[0];
  }
  return curEnvRole;
};

export const useCurrentEnvironment = (
  consoleAccount: ConsoleAccount
): Environment => {
  return currentEnvironmentRole(consoleAccount).environment;
};

export const useIsEditable = (consoleAccount: ConsoleAccount): boolean => {
  if (consoleAccount.isSystemAdmin) return true;

  const envRole = currentEnvironmentRole(consoleAccount);
  return envRole.role === 'Environment_EDITOR';
};
