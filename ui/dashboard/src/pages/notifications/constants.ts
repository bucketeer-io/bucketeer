import { i18n } from 'i18n';
import { NotificationOption } from './types';

const translate = i18n.t;

export const SOURCE_TYPE_ITEMS: NotificationOption[] = [
  {
    label: translate(`source-type.account`),
    description: translate(`source-type.account-description`),
    value: 'DOMAIN_EVENT_ACCOUNT'
  },
  {
    label: translate(`source-type.api-key`),
    description: translate(`source-type.api-key-description`),
    value: 'DOMAIN_EVENT_APIKEY'
  },
  {
    label: translate(`source-type.auto-ops`),
    description: translate(`source-type.auto-ops-description`),
    value: 'DOMAIN_EVENT_AUTOOPS_RULE'
  },
  {
    label: translate(`source-type.experiment`),
    description: translate(`source-type.experiment-description`),
    value: 'DOMAIN_EVENT_EXPERIMENT'
  },
  {
    label: translate(`source-type.feature-flag`),
    description: translate(`source-type.feature-flag-description`),
    value: 'DOMAIN_EVENT_FEATURE'
  },
  {
    label: translate(`source-type.goal`),
    description: translate(`source-type.goal-description`),
    value: 'DOMAIN_EVENT_GOAL'
  },
  {
    label: translate(`source-type.mau-count`),
    description: translate(`source-type.mau-count-description`),
    value: 'MAU_COUNT'
  },
  {
    label: translate(`source-type.notification`),
    description: translate(`source-type.notification-description`),
    value: 'DOMAIN_EVENT_SUBSCRIPTION'
  },
  {
    label: translate(`source-type.push`),
    description: translate(`source-type.push-description`),
    value: 'DOMAIN_EVENT_PUSH'
  },
  {
    label: translate(`source-type.running-experiments`),
    description: translate(`source-type.running-experiments-description`),
    value: 'EXPERIMENT_RUNNING'
  },
  {
    label: translate(`source-type.segment`),
    description: translate(`source-type.segment-description`),
    value: 'DOMAIN_EVENT_SEGMENT'
  },
  {
    label: translate(`source-type.stale-feature-flag`),
    description: translate(`source-type.stale-feature-flag-description`),
    value: 'FEATURE_STALE'
  }
];
