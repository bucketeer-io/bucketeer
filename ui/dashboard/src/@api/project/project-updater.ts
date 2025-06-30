import axiosClient from '@api/axios-client';
import { ProjectResponse } from './project-creator';

export interface ProjectUpdaterPayload {
  id: string;
  description?: string;
  name: string;
}

export const projectUpdater = async (
  params?: ProjectUpdaterPayload
): Promise<ProjectResponse> => {
  return axiosClient
    .post<ProjectResponse>('/v1/environment/update_project', params)
    .then(response => response.data);
};
