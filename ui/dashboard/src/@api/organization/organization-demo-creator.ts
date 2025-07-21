import axios, { AxiosInstance } from 'axios';
import { urls } from 'configs';
import { getDemoTokenStorage } from 'storage/demo-token';
import { Organization } from '@types';

export interface OrganizationDemoCreatorPayload {
  name: string;
  urlCode: string;
  description?: string;
  ownerEmail?: string;
}

export interface OrganizationDemoResponse {
  organization: Organization;
}

const axiosClient: AxiosInstance = axios.create({
  baseURL: urls.WEB_API_ENDPOINT
});

axiosClient.interceptors.request.use(
  config => {
    const authToken = getDemoTokenStorage();
    if (authToken) {
      config.headers['Authorization'] = `Bearer ${authToken.accessToken}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

export const organizationDemoCreator = async (
  payload?: OrganizationDemoCreatorPayload
): Promise<OrganizationDemoResponse> => {
  return axiosClient
    .post<OrganizationDemoResponse>(
      '/v1/environment/create_demo_organization',
      payload
    )
    .then(response => response.data);
};
