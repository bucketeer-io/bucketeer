import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { PageContext, PageType, Suggestion } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface SuggestionsFetcherParams {
  environmentId: string;
  pageContext?: PageContext;
}

export interface SuggestionsFetcherResponse {
  suggestions: Suggestion[];
}

const protoPageTypeMap: Record<PageType, string> = {
  feature_flags: 'PAGE_TYPE_FEATURE_FLAGS',
  targeting: 'PAGE_TYPE_TARGETING',
  experiments: 'PAGE_TYPE_EXPERIMENTS',
  segments: 'PAGE_TYPE_SEGMENTS',
  autoops: 'PAGE_TYPE_AUTOOPS'
};

export const suggestionsFetcher = async (
  params: SuggestionsFetcherParams
): Promise<SuggestionsFetcherResponse> => {
  const queryParams = pickBy(
    {
      environmentId: params.environmentId,
      'pageContext.pageType': params.pageContext?.pageType
        ? protoPageTypeMap[params.pageContext.pageType]
        : undefined,
      'pageContext.featureId': params.pageContext?.featureId
    },
    v => isNotEmpty(v)
  );

  return axiosClient
    .get<SuggestionsFetcherResponse>(
      `/v1/aichat/suggestions?${stringifyParams(queryParams)}`
    )
    .then(response => response.data);
};
