import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { UserSegment } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface UserSegmentFetcherParams {
  id: string;
  environmentId: string;
}

export interface UserSegmentResponse {
  segment: Array<UserSegment>;
}

export const userSegmentFetcher = async (
  params?: UserSegmentFetcherParams
): Promise<UserSegmentResponse> => {
  return axiosClient
    .get<UserSegmentResponse>('/v1/segment', {
      params: pickBy(params, v => isNotEmpty(v))
    })
    .then(response => response.data);
};
