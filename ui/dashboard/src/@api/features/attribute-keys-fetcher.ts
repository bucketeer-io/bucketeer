import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { AttributeKeysResponse } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface AttributeKeysParams {
  environmentId: string;
}

export const attributeKeysFetcher = async (
  params?: AttributeKeysParams
): Promise<AttributeKeysResponse> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<AttributeKeysResponse>(
      `/v1/feature/user-attribute-keys?${requestParams}`
    )
    .then(response => response.data);
};
