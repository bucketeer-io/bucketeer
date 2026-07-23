import axiosClient from '@api/axios-client';
import pickBy from 'lodash/pickBy';
import { CollectionParams, NotificationCenterDraftCollection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams } from 'utils/search-params';

export interface NotificationDraftsFetcherParams extends CollectionParams {
  environmentId?: string;
}

// GET /v1/notifications/drafts — list draft notifications. System admin only.
export const notificationDraftsFetcher = async (
  params?: NotificationDraftsFetcherParams
): Promise<NotificationCenterDraftCollection> => {
  const requestParams = stringifyParams(pickBy(params, v => isNotEmpty(v)));

  return axiosClient
    .get<NotificationCenterDraftCollection>(
      `/v1/notifications/drafts?${requestParams}`
    )
    .then(response => response.data);
};
