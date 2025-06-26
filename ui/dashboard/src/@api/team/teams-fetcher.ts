import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, TeamCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface TeamsFetcherParams extends CollectionParams {
  organizationId: string;
}

export const teamsFetcher = async (
  params?: TeamsFetcherParams
): Promise<TeamCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<TeamCollection>(`/v1/teams?${requestParams}`)
    .then(response => response.data);
};
