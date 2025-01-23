import { yupLocale } from '../../lang/yup';
import * as yup from 'yup';

import {
  NOTIFICATION_NAME_MAX_LENGTH,
  NOTIFICATION_SOURCE_TYPES_MIN_LENGTH
} from '../../constants/notification';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { Subscription } from '../../proto/notification/subscription_pb';

yup.setLocale(yupLocale);

const nameSchema = yup.string().required().max(NOTIFICATION_NAME_MAX_LENGTH);
const sourceTypesSchema = yup
  .array()
  .required()
  .min(
    NOTIFICATION_SOURCE_TYPES_MIN_LENGTH,
    intl.formatMessage(messages.input.error.minSelectOptionLength)
  );

const webhookUrlSchema = yup.string().required().url();
const featureFlagTagsListSchema = yup
  .array()
  .of(yup.string())
  .test(
    'tags-required-for-flag',
    intl.formatMessage(messages.input.error.required),
    function (featureFlagTagsList) {
      const { sourceTypes } = this.parent;
      const hasFlag = sourceTypes?.includes(
        Subscription.SourceType.DOMAIN_EVENT_FEATURE.toString()
      );

      // If 'Flag' is selected, ensure tags are present and not empty.
      if (
        hasFlag &&
        (!featureFlagTagsList || featureFlagTagsList.length === 0)
      ) {
        return false;
      }
      return true;
    }
  );

export const addFormSchema = yup.object().shape({
  name: nameSchema,
  sourceTypes: sourceTypesSchema,
  webhookUrl: webhookUrlSchema,
  featureFlagTagsList: featureFlagTagsListSchema
});

export const updateFormSchema = yup.object().shape({
  name: nameSchema,
  sourceTypes: sourceTypesSchema,
  featureFlagTagsList: featureFlagTagsListSchema
});
