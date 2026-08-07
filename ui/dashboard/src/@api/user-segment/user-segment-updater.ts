import axiosClient from '@api/axios-client';
import { SegmentRuleListValue, UserSegment } from '@types';

export interface UserSegmentUpdaterPayload {
  id: string;
  environmentId: string;
  name?: string;
  description?: string;
  // Absent = rules unchanged; present = full replacement of all rules.
  rules?: SegmentRuleListValue;
}

export interface UserSegmentUpdaterResponse {
  segment: UserSegment;
}

export const userSegmentUpdater = async (
  params?: UserSegmentUpdaterPayload
): Promise<UserSegmentUpdaterResponse> => {
  return axiosClient
    .patch<UserSegmentUpdaterResponse>('/v1/segment', params)
    .then(response => response.data);
};
