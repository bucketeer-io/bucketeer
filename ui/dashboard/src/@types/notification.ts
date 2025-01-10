export type RecipientType = 'SlackChannel';
export type NotificationLanguage = 'ENGLISH' | 'JAPANESE';

export type SourceType =
  | 'DOMAIN_EVENT_FEATURE'
  | 'DOMAIN_EVENT_GOAL'
  | 'DOMAIN_EVENT_EXPERIMENT'
  | 'DOMAIN_EVENT_ACCOUNT'
  | 'DOMAIN_EVENT_APIKEY'
  | 'DOMAIN_EVENT_SEGMENT'
  | 'DOMAIN_EVENT_ENVIRONMENT'
  | 'DOMAIN_EVENT_ADMIN_ACCOUNT'
  | 'DOMAIN_EVENT_AUTOOPS_RULE'
  | 'DOMAIN_EVENT_PUSH'
  | 'DOMAIN_EVENT_SUBSCRIPTION'
  | 'DOMAIN_EVENT_ADMIN_SUBSCRIPTION'
  | 'DOMAIN_EVENT_PROJECT'
  | 'DOMAIN_EVENT_WEBHOOK'
  | 'DOMAIN_EVENT_PROGRESSIVE_ROLLOUT'
  | 'DOMAIN_EVENT_ORGANIZATION'
  | 'DOMAIN_EVENT_FLAG_TRIGGER'
  | 'FEATURE_STALE'
  | 'EXPERIMENT_RUNNING'
  | 'MAU_COUNT';

export interface NotificationRecipient {
  type: RecipientType;
  slackChannelRecipient: {
    webhookUrl: string;
  };
  language: NotificationLanguage;
}

export interface Notification {
  id: string;
  createdAt: string;
  updatedAt: string;
  disabled: boolean;
  sourceTypes: SourceType[];
  recipient: NotificationRecipient;
  name: string;
  environmentName: string;
  environmentId: string;
}

export interface NotificationsCollection {
  subscriptions: Array<Notification>;
  cursor: string;
  totalCount: string;
}
