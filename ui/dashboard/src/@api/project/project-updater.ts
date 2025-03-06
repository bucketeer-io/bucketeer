import axiosClient from '@api/axios-client';

export interface ProjectUpdaterPayload {
  id: string;
  description?: string;
  name: string;
}

export const projectUpdater = async (params?: ProjectUpdaterPayload) => {
  return axiosClient
    .post('/v1/environment/update_project', params)
    .then(response => response.data);
};
