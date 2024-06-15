import { Option } from '../components/CheckBoxList';
import { intl } from '../lang';
import { messages } from '../lang/messages';
import { Subscription } from '../proto/notification/subscription_pb';

export const NOTIFICATION_LIST_PAGE_SIZE = 50;
export const NOTIFICATION_NAME_MAX_LENGTH = 100;
export const NOTIFICATION_SOURCE_TYPES_MIN_LENGTH = 1;

export const SOURCE_TYPE_ITEMS: Option[] = [
  {
    label: intl.formatMessage(messages.sourceType.account),
    description: intl.formatMessage(messages.sourceType.accountDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_ACCOUNT.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.apiKey),
    description: intl.formatMessage(messages.sourceType.apiKeyDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_APIKEY.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.autoOps),
    description: intl.formatMessage(messages.sourceType.autoOpsDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_AUTOOPS_RULE.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.experiment),
    description: intl.formatMessage(messages.sourceType.experimentDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_EXPERIMENT.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.featureFlag),
    description: intl.formatMessage(messages.sourceType.featureFlagDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_FEATURE.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.goal),
    description: intl.formatMessage(messages.sourceType.goalDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_GOAL.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.mauCount),
    description: intl.formatMessage(messages.sourceType.mauCountDescription),
    value: Subscription.SourceType.MAU_COUNT.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.notification),
    description: intl.formatMessage(
      messages.sourceType.notificationDescription
    ),
    value: Subscription.SourceType.DOMAIN_EVENT_SUBSCRIPTION.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.push),
    description: intl.formatMessage(messages.sourceType.pushDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_PUSH.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.runningExperiments),
    description: intl.formatMessage(
      messages.sourceType.runningExperimentsDescription
    ),
    value: Subscription.SourceType.EXPERIMENT_RUNNING.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.segment),
    description: intl.formatMessage(messages.sourceType.segmentDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_SEGMENT.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.staleFeatureFlag),
    description: intl.formatMessage(
      messages.sourceType.staleFeatureFlagDescription
    ),
    value: Subscription.SourceType.FEATURE_STALE.toString(),
  },
];
