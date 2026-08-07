import axiosClient from '@api/axios-client';
import { SegmentRulePayload, UserSegment } from '@types';

export interface UserSegmentCreatorPayload {
  environmentId: string;
  name: string;
  description?: string;
  rules?: SegmentRulePayload[];
}

export interface UserSegmentCreatorResponse {
  segment: UserSegment;
}

export const userSegmentCreator = async (
  params?: UserSegmentCreatorPayload
): Promise<UserSegmentCreatorResponse> => {
  return axiosClient
    .post<UserSegmentCreatorResponse>('/v1/segment', params)
    .then(response => response.data);
};
