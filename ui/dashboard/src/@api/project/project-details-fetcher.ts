import axiosClient from '@api/axios-client';
import { Project } from '@types';

export interface ProjectDetailsFetcherParams {
  id: string;
}
export interface ProjectDetailsResponse {
  project: Project;
}

export const projectDetailsFetcher = async (
  params?: ProjectDetailsFetcherParams
): Promise<ProjectDetailsResponse> => {
  return axiosClient
    .post<ProjectDetailsResponse>('/v1/environment/get_project', params)
    .then(response => response.data);
};
