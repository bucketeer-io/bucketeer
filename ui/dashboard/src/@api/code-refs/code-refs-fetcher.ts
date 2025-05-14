import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CodeRefCollection, CollectionParams, RepositoryType } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface CodeRefsFetcherParams extends CollectionParams {
  environmentId: string;
  featureId: string;
  repositoryName: string;
  repositoryOwner?: string;
  repositoryType?: RepositoryType;
  repositoryBranch?: string;
  fileExtension?: string;
}

export const codeRefsFetcher = async (
  _params?: CodeRefsFetcherParams
): Promise<CodeRefCollection> => {
  const params = pickBy(_params, v => isNotEmpty(v));

  return axiosClient
    .get<CodeRefCollection>(`/v1/code_references?${stringifyParams(params)}`)
    .then(response => response.data);
};
