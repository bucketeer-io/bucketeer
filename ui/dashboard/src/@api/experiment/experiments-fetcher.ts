import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import {
  CollectionParams,
  ExperimentCollection,
  ExperimentStatus
} from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface ExperimentsFetcherParams extends CollectionParams {
  environmentId: string;
  featureId?: string;
  featureVersion?: number;
  from?: string;
  to?: string;
  statuses?: ExperimentStatus[];
  maintainer?: string;
  archived?: boolean;
}

export const experimentsFetcher = async (
  params?: ExperimentsFetcherParams
): Promise<ExperimentCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<ExperimentCollection>(`/v1/experiments?${requestParams}`)
    .then(response => response.data);
};
