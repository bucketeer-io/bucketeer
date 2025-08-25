export interface Environment {
  archived: boolean;
  createdAt: string;
  description: string;
  id: string;
  name: string;
  organizationId: string;
  projectId: string;
  requireComment: boolean;
  featureFlagCount: number;
  updatedAt: string;
  urlCode: string;
}

export interface EnvironmentCollection {
  environments: Array<Environment>;
  cursor: string;
  totalCount: string;
}

export interface EnvironmentResponse {
  environment: Environment;
}
