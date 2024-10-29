import axiosClient from '@api/axios-client';
import { Project } from '@types';

export interface ProjectCreatorCommand {
  id: string;
  name: string;
  urlCode: string;
  description?: string;
}

export interface ProjectCreatorParams {
  command: ProjectCreatorCommand;
}

export interface ProjectResponse {
  project: Array<Project>;
}

export const projectCreator = async (
  params?: ProjectCreatorParams
): Promise<ProjectResponse> => {
  return axiosClient
    .post<ProjectResponse>('/v1/environment/create_project', params)
    .then(response => response.data);
};
