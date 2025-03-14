import { AnyObject } from 'yup';
import { OperationStatus } from './goal';

export type AutoOpsType = 'TYPE_UNKNOWN' | 'SCHEDULE' | 'EVENT_RATE';
export type ClauseActionType = 'UNKNOWN' | 'ENABLE' | 'DISABLE';

export interface AutoOpsRule {
  id: string;
  featureId: string;
  opsType: AutoOpsType;
  clauses: [
    {
      id: string;
      clause: AnyObject;
      actionType: ClauseActionType;
    }
  ];
  createdAt: string;
  updatedAt: string;
  deleted: boolean;
  autoOpsStatus: OperationStatus;
  featureName: string;
}

export interface AutoOpsRuleCollection {
  autoOpsRules: AutoOpsRule[];
  cursor: string;
}
