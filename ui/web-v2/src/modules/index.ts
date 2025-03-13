import { AnyAction, combineReducers } from 'redux';
import { ThunkAction } from 'redux-thunk';

import { accountsSlice } from './accounts';
import { adminNotificationSlice } from './adminNotifications';
import { apiKeySlice } from './apiKeys';
import { codeRefsSlice } from './codeRefs';
import { auditLogSlice } from './auditLogs';
import { authSlice } from './auth';
import { autoOpsRulesSlice } from './autoOpsRules';
import { environmentsSlice } from './environments';
import { evaluationTimeseriesCountSlice } from './evaluationTimeseriesCount';
import { experimentResultSlice } from './experimentResult';
import { experimentsSlice } from './experiments';
import { featuresSlice } from './features';
import { flagTriggersSlice } from './flagTriggers';
import { goalsSlice } from './goals';
import { meSlice } from './me';
import { myOrganizationSlice } from './myOrganization';
import { notificationSlice } from './notifications';
import { opsCountsSlice } from './opsCounts';
import { progressiveRolloutSlice } from './porgressiveRollout';
import { projectsSlice } from './projects';
import { pushSlice } from './pushes';
import { segmentsSlice } from './segments';
import { tagsSlice } from './tags';
import { toastsSlice } from './toasts';

export const reducers = combineReducers({
  accounts: accountsSlice.reducer,
  adminNotification: adminNotificationSlice.reducer,
  auditLog: auditLogSlice.reducer,
  apiKeys: apiKeySlice.reducer,
  codeRefs: codeRefsSlice.reducer,
  auth: authSlice.reducer,
  autoOpsRules: autoOpsRulesSlice.reducer,
  progressiveRollout: progressiveRolloutSlice.reducer,
  opsCounts: opsCountsSlice.reducer,
  environments: environmentsSlice.reducer,
  evaluationTimeseriesCount: evaluationTimeseriesCountSlice.reducer,
  experiments: experimentsSlice.reducer,
  experimentResults: experimentResultSlice.reducer,
  features: featuresSlice.reducer,
  goals: goalsSlice.reducer,
  notification: notificationSlice.reducer,
  projects: projectsSlice.reducer,
  push: pushSlice.reducer,
  segments: segmentsSlice.reducer,
  me: meSlice.reducer,
  toasts: toastsSlice.reducer,
  tags: tagsSlice.reducer,
  myOrganization: myOrganizationSlice.reducer,
  flagTriggers: flagTriggersSlice.reducer
});

export type AppState = ReturnType<typeof reducers>;
export type AppThunk = ThunkAction<
  Promise<unknown>,
  AppState,
  unknown,
  AnyAction
>;
