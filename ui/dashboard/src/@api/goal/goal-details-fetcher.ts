import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { Goal } from '@types';
import { isNotEmpty } from 'utils/data-type';

export interface GoalDetailsFetcherParams {
  id: string;
  environmentId: string;
}

export interface GoalDetailsResponse {
  goal: Goal;
}

export const goalDetailsFetcher = async (
  params?: GoalDetailsFetcherParams
): Promise<GoalDetailsResponse> => {
  return axiosClient
    .get<GoalDetailsResponse>('/v1/goal', {
      params: pickBy(params, v => isNotEmpty(v))
    })
    .then(response => response.data);
};
