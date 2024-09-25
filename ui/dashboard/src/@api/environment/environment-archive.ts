import axiosClient from '@api/axios-client';
import { ArchiveParams } from '@api/organization';

export const environmentArchive = async (params?: ArchiveParams) => {
  return axiosClient
    .post('/v1/environment/archive_environment', params)
    .then(response => response.data);
};
