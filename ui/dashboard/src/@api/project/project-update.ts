import axiosClient from '@api/axios-client';

export interface ProjectUpdateParams {
  id: string;
  changeDescriptionCommand: {
    description: string;
  };
  renameCommand: {
    name: string;
  };
}

export const projectUpdate = async (params?: ProjectUpdateParams) => {
  return axiosClient
    .post('/v1/environment/update_project', params)
    .then(response => response.data);
};
