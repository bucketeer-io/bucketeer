import axiosClient from '@api/axios-client';
import { UserSegment } from '@types';

export interface UserSegmentUpdaterParams {
  id: string;
  environmentId: string;
  name?: string;
  description?: string;
}

export interface UserSegmentUpdaterResponse {
  segment: UserSegment;
}

export const userSegmentUpdater = async (
  params?: UserSegmentUpdaterParams
): Promise<UserSegmentUpdaterResponse> => {
  return axiosClient
    .patch<UserSegmentUpdaterResponse>('/v1/segment', params)
    .then(response => response.data);
};
