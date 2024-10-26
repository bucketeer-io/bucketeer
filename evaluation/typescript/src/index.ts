import { Evaluator } from './evaluation';
import { SegmentUser, SegmentUsers } from './proto/feature/segment_pb';
import { Feature } from './proto/feature/feature_pb';
import { NewUserEvaluations } from './userEvaluation';
import {
  createClause,
  createFixedStrategy,
  createRolloutStrategy,
  createReason,
  createRule,
  createPrerequisite,
  creatFeature,
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

export { Evaluator };
export { NewUserEvaluations };
export { SegmentUser, SegmentUsers, Feature };
export {
  createClause,
  createFixedStrategy,
  createRolloutStrategy,
  createReason,
  createRule,
  createPrerequisite,
  creatFeature,
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
