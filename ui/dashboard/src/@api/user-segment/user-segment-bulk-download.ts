import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { isNotEmpty } from 'utils/data-type';

export interface UserSegmentBulkDownloadParams {
  segmentId: string;
  environmentId: string;
  state?: 'INCLUDED' | 'EXCLUDED';
}

export interface UserSegmentBulkDownloadResponse {
  data: string;
}

export const userSegmentBulkDownload = async (
  params?: UserSegmentBulkDownloadParams
): Promise<UserSegmentBulkDownloadResponse> => {
  return axiosClient
    .get<UserSegmentBulkDownloadResponse>('/v1/segment_users/bulk_download', {
      params: pickBy(params, v => isNotEmpty(v))
    })
    .then(response => response.data);
};
