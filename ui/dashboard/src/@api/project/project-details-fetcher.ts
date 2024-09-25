import axiosClient from '@api/axios-client';
import { Project } from '@types';

export type DetailsFetcherParams = {
  id: string;
};

export type ProjectDetailsCollection = {
  project: Project;
};

export const projectDetailsFetcher = async (
  params?: DetailsFetcherParams
): Promise<ProjectDetailsCollection> => {
  return axiosClient
    .post<ProjectDetailsCollection>('/v1/environment/get_project', params)
    .then(response => response.data);
};
