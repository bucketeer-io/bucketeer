import axiosClient from '@api/axios-client';

export interface UserSegmentBulkUploadParams {
  segmentId: string;
  environmentId?: string;
  data?: string;
  state?: 'INCLUDED' | 'EXCLUDED';
}

export const userSegmentBulkUpload = async (
  params?: UserSegmentBulkUploadParams
) => {
  return axiosClient
    .post('/v1/segment_users/bulk_upload', params)
    .then(response => response.data);
};
