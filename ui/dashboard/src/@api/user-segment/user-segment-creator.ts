import axiosClient from '@api/axios-client';
import { UserSegment } from '@types';

export interface UserSegmentCreatorParams {
  environmentId: string;
  name: string;
  description?: string;
}

export interface UserSegmentCreatorResponse {
  segment: UserSegment;
}

export const userSegmentCreator = async (
  params?: UserSegmentCreatorParams
): Promise<UserSegmentCreatorResponse> => {
  return axiosClient
    .post<UserSegmentCreatorResponse>('/v1/segment', params)
    .then(response => response.data);
};
