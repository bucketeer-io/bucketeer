import axiosClient from '@api/axios-client';
import { DebuggerCollection } from '@types';

interface UserPayload {
  id: string;
  data?: {
    [key: string]: string | string[];
  };
}

export interface EvaluationCreatorPayload {
  users: UserPayload[];
  featureIds: string[];
  environmentId: string;
}

export const debuggerEvaluate = async (
  params?: EvaluationCreatorPayload
): Promise<DebuggerCollection> => {
  return axiosClient
    .post<DebuggerCollection>('/v1/debug_evaluate_features', params)
    .then(response => response.data);
};
