import { Evaluator, getFeatureIDsDependsOn } from './evaluation';
import { SegmentUser, SegmentUsers } from './proto/feature/segment_pb';
import { Feature } from './proto/feature/feature_pb';
import { NewUserEvaluations } from './userEvaluation';
import { User } from './proto/user/user_pb';
import { Evaluation, UserEvaluations } from './proto/feature/evaluation_pb';
import { Strategy } from './proto/feature/strategy_pb';
import { Clause } from './proto/feature/clause_pb';
import { Reason } from './proto/feature/reason_pb';
import {
  createClause,
  createFixedStrategy,
  createRolloutStrategy,
  createReason,
  createRule,
  createPrerequisite,
  createFeature,
  createEvaluation,
  createSegmentUser,
  createStrategy,
  createTarget,
  createUser,
  createVariation,
} from './modelFactory';
import {
  GetFeatureFlagsRequest,
  GetFeatureFlagsResponse,
  GetSegmentUsersRequest,
  GetSegmentUsersResponse,
} from './proto/gateway/service_pb';
import { SourceId } from './proto/event/client/event_pb';
import { GatewayClient, ServiceError } from './proto/gateway/service_pb_service';

export { Evaluator, NewUserEvaluations, Evaluation, UserEvaluations, getFeatureIDsDependsOn };
export { User, SegmentUser, SegmentUsers, Feature };
export { Strategy, Clause, Reason };
export {
  createClause,
  createFixedStrategy,
  createRolloutStrategy,
  createReason,
  createRule,
  createPrerequisite,
  createFeature,
  createEvaluation,
  createSegmentUser,
  createStrategy,
  createTarget,
  createUser,
  createVariation,
};
export {
  GetFeatureFlagsRequest,
  GetFeatureFlagsResponse,
  GetSegmentUsersRequest,
  GetSegmentUsersResponse,
};
export { SourceId };
export { GatewayClient, ServiceError };
