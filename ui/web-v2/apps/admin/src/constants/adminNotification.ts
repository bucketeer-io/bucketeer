import { Option } from '../components/CheckBoxList';
import { intl } from '../lang';
import { messages } from '../lang/messages';
import { Subscription } from '../proto/notification/subscription_pb';

export const NOTIFICATION_LIST_PAGE_SIZE = 50;
export const NOTIFICATION_NAME_MAX_LENGTH = 100;
export const NOTIFICATION_SOURCE_TYPES_MIN_LENGTH = 1;

export const SOURCE_TYPE_ITEMS: Option[] = [
  {
    label: intl.formatMessage(messages.sourceType.project),
    description: intl.formatMessage(messages.sourceType.projectDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_PROJECT.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.environment),
    description: intl.formatMessage(messages.sourceType.environmentDescription),
    value: Subscription.SourceType.DOMAIN_EVENT_ENVIRONMENT.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.adminAccount),
    description: intl.formatMessage(
      messages.sourceType.adminAccountDescription
    ),
    value: Subscription.SourceType.DOMAIN_EVENT_ADMIN_ACCOUNT.toString(),
  },
  {
    label: intl.formatMessage(messages.sourceType.adminNotification),
    description: intl.formatMessage(
      messages.sourceType.adminNotificationDescription
    ),
    value: Subscription.SourceType.DOMAIN_EVENT_ADMIN_SUBSCRIPTION.toString(),
  },
];
