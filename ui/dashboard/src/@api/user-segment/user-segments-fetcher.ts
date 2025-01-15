import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import {
  UserSegmentCollection,
  CollectionParams,
  FeatureSegmentStatus
} from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface UserSegmentsFetcherParams extends CollectionParams {
  environmentId?: string;
  isInUseStatus?: boolean;
  status?: FeatureSegmentStatus;
}

export const userSegmentsFetcher = async (
  params?: UserSegmentsFetcherParams
): Promise<UserSegmentCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<UserSegmentCollection>(`/v1/segments?${requestParams}`)
    .then(response => response.data);
};
