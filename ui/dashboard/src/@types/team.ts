import { AutoOpsChangeType } from './auto-ops';

export interface Team {
  id: string;
  name: string;
  description: string;
  organizationId: string;
  organizationName: string;
  createdAt: string;
  updatedAt: string;
}

export interface TeamCollection {
  teams: Team[];
  nextCursor: string;
  totalCount: string;
}

export interface TeamChange {
  changeType: AutoOpsChangeType;
  team: string;
}
