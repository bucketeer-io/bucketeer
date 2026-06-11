/**
 * Single source of truth for cross-entity cache invalidation.
 *
 * Every successful non-GET response that goes through the shared `axiosClient`
 * is matched against this table by the response interceptor in
 * `axios-client.ts`, and every matching query key is invalidated on the shared
 * `queryClient`.
 *
 * Adding a new mutating endpoint or a new cross-entity dependency only
 * requires editing this file.
 *
 * Notes:
 * - Keys are the literal strings used as the first element of `queryKey`
 *   tuples in `src/@queries/*.ts`. Adding a new query module means the new key
 *   should be referenced here for any endpoint that mutates data it depends
 *   on.
 * - `audit-logs` (and `admin-audit-logs` for org-level admin actions) is
 *   listed for almost every mutation because the audit log is a global feed
 *   of every change; under-invalidating leads to stale audit log views.
 */
export type CacheInvalidationRule = {
  match: RegExp;
  keys: readonly string[];
};

export const URL_TO_KEYS: readonly CacheInvalidationRule[] = [
  // Feature flags (create / update / clone / bulk_clone)
  {
    match: /\/v1\/feature(\/clone|\/bulk_clone)?(\?|$)/,
    keys: [
      'features',
      'feature-details',
      'experiments',
      'experiment-details',
      'evaluation-timeseries',
      'auto-ops-rules',
      'auto-ops-count',
      'rollouts',
      'triggers',
      'scheduled-flag-changes-query-key',
      'scheduled-flag-change-query-key',
      'scheduled-flag-change-summary-query-key',
      'schedule-flags-query-key',
      'tags',
      'user-segments',
      'segment-details',
      'code-refs',
      'histories',
      'audit-logs'
    ]
  },

  // Scheduled flag change (create / update / delete / execute)
  {
    match: /\/v1\/scheduled_flag_change(\/execute)?(\?|$)/,
    keys: [
      'scheduled-flag-changes-query-key',
      'scheduled-flag-change-query-key',
      'scheduled-flag-change-summary-query-key',
      'features',
      'feature-details',
      'evaluation-timeseries',
      'histories',
      'audit-logs'
    ]
  },

  // Legacy schedule flag (create / update / delete)
  {
    match: /\/v1\/schedule_flag(\?|$)/,
    keys: [
      'schedule-flags-query-key',
      'features',
      'feature-details',
      'evaluation-timeseries',
      'histories',
      'audit-logs'
    ]
  },

  // Goals
  {
    match: /\/v1\/goal(\?|$)/,
    keys: [
      'goals',
      'goal-details',
      'experiments',
      'experiment-details',
      'experiment-result-details',
      'histories',
      'audit-logs'
    ]
  },

  // Experiments
  {
    match: /\/v1\/experiment(\?|$)/,
    keys: [
      'experiments',
      'experiment-details',
      'experiment-result-details',
      'evaluation-timeseries',
      'features',
      'feature-details',
      'histories',
      'audit-logs'
    ]
  },

  // User segments
  {
    match: /\/v1\/segment(\?|$)/,
    keys: [
      'user-segments',
      'segment-details',
      'features',
      'feature-details',
      'histories',
      'audit-logs'
    ]
  },
  {
    match: /\/v1\/segment_users\//,
    keys: [
      'user-segments',
      'segment-details',
      'user-attribute-keys',
      'histories',
      'audit-logs'
    ]
  },

  // Auto-ops rules (create / update / delete / stop)
  {
    match: /\/v1\/auto_ops_rule(\/stop)?(\?|$)/,
    keys: [
      'auto-ops-rules',
      'auto-ops-count',
      'features',
      'feature-details',
      'evaluation-timeseries',
      'histories',
      'audit-logs'
    ]
  },

  // Progressive rollouts (create / stop / execute / delete)
  {
    match: /\/v1\/progressive_rollout(\/stop|\/execute)?(\?|$)/,
    keys: [
      'rollouts',
      'features',
      'feature-details',
      'evaluation-timeseries',
      'auto-ops-count',
      'histories',
      'audit-logs'
    ]
  },

  // Triggers
  {
    match: /\/v1\/flag_trigger(\?|$)/,
    keys: ['triggers', 'features', 'feature-details', 'histories', 'audit-logs']
  },

  // Tags
  {
    match: /\/v1\/tag(\?|$)/,
    keys: ['tags', 'features', 'feature-details', 'audit-logs']
  },

  // Notification subscriptions
  {
    match: /\/v1\/subscription(\?|$)/,
    keys: ['notifications', 'notification-details', 'audit-logs']
  },

  // Pushes
  {
    match: /\/v1\/push(\?|$)/,
    keys: ['pushes', 'push-details', 'audit-logs']
  },

  // Environments (create / update / archive / unarchive all hit these endpoints)
  {
    match: /\/v1\/environment\/(create|update)_environment/,
    keys: [
      'environments',
      'environment-details',
      'environments-multiple-ids',
      'projects',
      'project-details',
      'organizations',
      'organization-details',
      'accounts',
      'api-keys',
      'histories',
      'audit-logs'
    ]
  },

  // Projects
  {
    match: /\/v1\/environment\/(create|update)_project/,
    keys: [
      'projects',
      'project-details',
      'environments',
      'environment-details',
      'organizations',
      'organization-details',
      'accounts',
      'audit-logs'
    ]
  },

  // Organizations (incl. archive / unarchive / demo)
  {
    match: /\/v1\/environment\/(create|update|archive|unarchive)_organization/,
    keys: [
      'organizations',
      'organization-details',
      'projects',
      'project-details',
      'accounts',
      'audit-logs',
      'admin-audit-logs'
    ]
  },
  {
    match: /\/v1\/environment\/create_demo_organization/,
    keys: [
      'organizations',
      'organization-details',
      'projects',
      'project-details',
      'accounts',
      'audit-logs',
      'admin-audit-logs'
    ]
  },

  // Accounts (create / update / enable / disable / delete)
  {
    match: /\/v1\/account\/(create|update|enable|disable|delete)_account/,
    keys: ['accounts', 'teams', 'audit-logs']
  },

  // API keys
  {
    match: /\/v1\/account\/(create|update)_api_key/,
    keys: ['api-keys', 'api-keys-details', 'audit-logs']
  },

  // Teams
  {
    match: /\/v1\/team(\?|$)/,
    keys: ['teams', 'accounts', 'audit-logs']
  }
];

/**
 * Resolve the union of query keys to invalidate for a given request URL.
 * Returns an empty array when no rule matches (e.g. read-only endpoints,
 * auth/exchange flows, AI chat suggestions, etc.).
 */
export const resolveInvalidationKeys = (url: string): string[] => {
  const keys = new Set<string>();
  for (const { match, keys: ruleKeys } of URL_TO_KEYS) {
    if (match.test(url)) {
      ruleKeys.forEach(key => keys.add(key));
    }
  }
  return Array.from(keys);
};
