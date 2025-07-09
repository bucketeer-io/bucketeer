import { type Nullable } from 'option-t/nullable';

const KEY = 'project-environment';

interface ProjectEnvironment {
  environmentId: string;
  projectId: string;
}

export const getCurrentProjectEnvironmentStorage =
  (): Nullable<ProjectEnvironment> => {
    try {
      const data = window.localStorage.getItem(KEY);
      if (data) {
        return JSON.parse(data) as ProjectEnvironment;
      }
      return null;
    } catch (error) {
      console.error(error);
    }
    return null;
  };

export const setCurrentProjectEnvironmentStorage = (
  data: ProjectEnvironment
): void => {
  try {
    window.localStorage.setItem(KEY, JSON.stringify(data));
  } catch (error) {
    console.error(error);
  }
};

export const clearCurrentProjectEnvironmentStorage = (): void => {
  try {
    window.localStorage.removeItem(KEY);
  } catch (error) {
    console.error(error);
  }
};
