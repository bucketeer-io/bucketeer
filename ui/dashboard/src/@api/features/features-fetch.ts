import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, FeatureCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';
import { StatusFilterType } from 'pages/feature-flags/types';

export interface FeaturesFetcherParams extends CollectionParams {
  environmentId: string;
  maintainer?: string;
  archived?: boolean;
  hasExperiment?: boolean;
  enabled?: boolean;
  hasPrerequisites?: boolean;
  tags?: string[];
  status?: StatusFilterType;
  hasFeatureFlagAsRule?: boolean;
}

export const featuresFetcher = async (
  params?: FeaturesFetcherParams
): Promise<FeatureCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<FeatureCollection>(`/v1/features?${requestParams}`)
    .then(response => response.data);
};
