import axiosClient from '@api/axios-client';

export interface ProjectUpdaterParams {
  id: string;
  changeDescriptionCommand: {
    description: string;
  };
  renameCommand: {
    name: string;
  };
}

export const projectUpdater = async (params?: ProjectUpdaterParams) => {
  return axiosClient
    .post('/v1/environment/update_project', params)
    .then(response => response.data);
};
