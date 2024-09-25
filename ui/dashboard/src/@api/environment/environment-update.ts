import axiosClient from '@api/axios-client';

export interface EnvironmentUpdateParams {
  id: string;
  renameCommand: {
    name: string;
  };
  changeDescriptionCommand: {
    description: string;
  };
  changeRequireCommentCommand: {
    requireComment: boolean;
  };
}

export const environmentUpdate = async (params?: EnvironmentUpdateParams) => {
  return axiosClient
    .post('/v1/environment/update_environment', params)
    .then(response => response.data);
};
