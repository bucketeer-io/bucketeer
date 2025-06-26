import axiosClient from '@api/axios-client';
import { Team } from '@types';

export interface TeamCreatorPayload {
  name: string;
  description?: string;
  organizationId: string;
}

export interface TeamCreatorResponse {
  team: Team;
}

export const teamCreator = async (
  params?: TeamCreatorPayload
): Promise<TeamCreatorResponse> => {
  return axiosClient
    .post<TeamCreatorResponse>('/v1/team', params)
    .then(response => response.data);
};
