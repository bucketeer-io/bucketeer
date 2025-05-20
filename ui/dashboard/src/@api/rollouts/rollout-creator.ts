import axiosClient from '@api/axios-client';
import {
  Rollout,
  RolloutManualScheduleClause,
  RolloutTemplateScheduleClause
} from '@types';

export interface RolloutCreatorParams {
  environmentId: string;
  featureId: string;
  progressiveRolloutManualScheduleClause?: RolloutManualScheduleClause;
  progressiveRolloutTemplateScheduleClause?: RolloutTemplateScheduleClause;
}

export interface RolloutCreatorResponse {
  progressiveRollout: Rollout;
}

export const rolloutCreator = async (
  params?: RolloutCreatorParams
): Promise<RolloutCreatorResponse> => {
  return axiosClient
    .post<RolloutCreatorResponse>('/v1/progressive_rollout', params)
    .then(response => response.data);
};
