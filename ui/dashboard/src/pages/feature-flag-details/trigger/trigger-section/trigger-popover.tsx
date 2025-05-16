import {
  IconCheckCircleOutlineOutlined,
  IconEditOutlined,
  IconMoreHorizOutlined,
  IconRefreshOutlined
} from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { Trigger } from '@types';
import { IconDisable, IconTrash } from '@icons';
import { Popover } from 'components/popover';
import { TriggerAction } from '../types';

const TriggerPopover = ({
  trigger,
  onActions
}: {
  trigger: Trigger;
  onActions: (action: TriggerAction) => void;
}) => {
  const { t } = useTranslation(['table']);

  return (
    <Popover
      options={compact([
        {
          label: `${t('trigger.edit-desc')}`,
          icon: IconEditOutlined,
          value: TriggerAction.EDIT,
          color: 'gray-600'
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
            <span className="text-accent-red-500">
              {t('trigger.delete-trigger')}
            </span>
          ),
          icon: IconTrash,
          value: TriggerAction.DELETE,
          color: 'accent-red-500'
        }
      ])}
      icon={IconMoreHorizOutlined}
      onClick={value => onActions(value as TriggerAction)}
      align="end"
    />
  );
};

export default TriggerPopover;
