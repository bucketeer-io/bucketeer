import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import {
  CollectionParams,
  RolloutCollection,
  OperationStatus,
  RolloutType
} from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface RolloutsFetcherParams extends CollectionParams {
  environmentId?: string;
  featureIds?: string[];
  status?: OperationStatus;
  type?: RolloutType;
}

export const progressiveRolloutsFetcher = async (
  _params?: RolloutsFetcherParams
): Promise<RolloutCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<RolloutCollection>(
      `/v1/progressive_rollouts?${stringifyParams(params)}`
    )
    .then(response => response.data);
};
