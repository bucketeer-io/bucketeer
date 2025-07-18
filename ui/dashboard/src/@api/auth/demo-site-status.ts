import axiosClient from '@api/axios-client';
import { DemoSiteStatusType } from '@types';

export const demoSiteStatusFetcher = async (): Promise<DemoSiteStatusType> => {
  return axiosClient
    .get(`/v1/demo_site_status`)
    .then(response => response.data);
};
