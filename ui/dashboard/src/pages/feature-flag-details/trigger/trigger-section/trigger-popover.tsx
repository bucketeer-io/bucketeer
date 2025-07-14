import {
  IconCheckCircleOutlineOutlined,
  IconEditOutlined,
  IconRefreshOutlined
} from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Trigger } from '@types';
import { IconDisable, IconTrash } from '@icons';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import { TriggerAction } from '../types';

const TriggerPopover = ({
  disabled,
  trigger,
  onActions
}: {
  disabled: boolean;
  trigger: Trigger;
  onActions: (action: TriggerAction) => void;
}) => {
  const { t } = useTranslation(['table']);

  return (
    <DisabledPopoverTooltip
      options={compact([
        {
          label: `${t('trigger.edit-desc')}`,
          icon: IconEditOutlined,
          value: TriggerAction.EDIT,
          color: disabled ? undefined : 'gray-600'
        },
        {
          label: `${t(`trigger.${trigger.disabled ? 'enable-trigger' : 'disable-trigger'}`)}`,
          icon: trigger.disabled ? IconCheckCircleOutlineOutlined : IconDisable,
          value: trigger.disabled ? TriggerAction.ENABLE : TriggerAction.DISABLE
        },
        {
          label: `${t('trigger.reset-url')}`,
          icon: IconRefreshOutlined,
          value: TriggerAction.RESET
        },
        {
          label: (
            <span className={disabled ? '' : 'text-accent-red-500'}>
              {t('trigger.delete-trigger')}
            </span>
          ),
          icon: IconTrash,
          value: TriggerAction.DELETE,
          color: disabled ? undefined : 'accent-red-500'
        }
      ])}
      onClick={value => onActions(value as TriggerAction)}
    />
  );
};

export default TriggerPopover;
