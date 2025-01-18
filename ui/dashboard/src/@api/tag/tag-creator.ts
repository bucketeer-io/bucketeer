import axiosClient from '@api/axios-client';
import { EntityType, Tag } from '@types';

export type TagCreatorPayload = {
  name: string;
  environmentId: string;
  entityType: EntityType;
};

export interface TagCreatorResponse {
  tag: Tag;
}

export const tagCreator = async (
  payload?: TagCreatorPayload
): Promise<TagCreatorResponse> => {
  return axiosClient
    .post<TagCreatorResponse>('/v1/tag', payload)
    .then(response => response.data);
};
