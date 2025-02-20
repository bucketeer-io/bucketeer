import axiosClient from '@api/axios-client';
import { Project } from '@types';

export interface ProjectCreatorPayload {
  name: string;
  urlCode: string;
  organizationId: string;
  description?: string;
}

export interface ProjectResponse {
  project: Project;
}

export const projectCreator = async (
  params?: ProjectCreatorPayload
): Promise<ProjectResponse> => {
  return axiosClient
    .post<ProjectResponse>('/v1/environment/create_project', params)
    .then(response => response.data);
};
