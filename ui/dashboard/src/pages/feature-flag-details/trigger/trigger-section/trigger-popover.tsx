import React from 'react';
import {
  IconEditOutlined,
  IconMoreHorizOutlined,
  IconRefreshOutlined
} from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { IconDisable, IconTrash } from '@icons';
import { Popover } from 'components/popover';
import { TriggerAction } from '../types';

const TriggerPopover = ({
  onActions
}: {
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
          label: `${t('trigger.disable-trigger')}`,
          icon: IconDisable,
          value: TriggerAction.DISABLE
        },
        {
          label: `${t('trigger.reset-url')}`,
          icon: IconRefreshOutlined,
          value: TriggerAction.RESET
        },
        {
          label: `${t('trigger.delete-trigger')}`,
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
