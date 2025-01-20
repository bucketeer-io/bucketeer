import axiosClient from '@api/axios-client';

export interface ProjectUpdaterParams {
  id: string;
  description?: string;
  name: string;
}

export const projectUpdater = async (params?: ProjectUpdaterParams) => {
  return axiosClient
    .post('/v1/environment/update_project', params)
    .then(response => response.data);
};
