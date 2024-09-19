export interface Project {
  createdAt: string;
  creatorEmail: string;
  description: string;
  disabled: boolean;
  id: string;
  name: string;
  organizationId: string;
  environmentCount: number;
  featureFlagCount: number;
  trial: boolean;
  updatedAt: string;
  urlCode: string;
}

export interface ProjectCollection {
  projects: Array<Project>;
  cursor: string;
  totalCount: string;
}
