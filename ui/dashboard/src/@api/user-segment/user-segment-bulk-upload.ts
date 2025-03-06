import axiosClient from '@api/axios-client';

export interface UserSegmentBulkUploadPayload {
  segmentId: string;
  environmentId?: string;
  data?: string;
  state?: 'INCLUDED' | 'EXCLUDED';
}

export const userSegmentBulkUpload = async (
  params?: UserSegmentBulkUploadPayload
) => {
  return axiosClient
    .post('/v1/segment_users/bulk_upload', params)
    .then(response => response.data);
};
