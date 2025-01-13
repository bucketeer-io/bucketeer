export interface Goal {
  id: string;
  name: string;
  description: string;
  deleted: boolean;
  createdAt: string;
  updatedAt: string;
  isInUseStatus: boolean;
  archived: boolean;
  connections: {
    type: string;
    data: {
      id: string;
      name: string;
    }[];
  } | null;
}

export interface GoalCollection {
  goals: Array<Goal>;
  cursor: string;
  totalCount: string;
}
