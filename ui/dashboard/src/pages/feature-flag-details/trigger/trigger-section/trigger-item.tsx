import { useCallback, useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { TriggerActionType, TriggerItemType } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { IconCopy, IconToastWarning, IconWatch, IconWebhook } from '@icons';
import Button from 'components/button';
import Divider from 'components/divider';
import Icon from 'components/icon';
import Input from 'components/input';
import DateTooltip from 'elements/date-tooltip';
import FeatureFlagStatus from 'elements/feature-flag-status';
import { TriggerAction } from '../types';
import TriggerPopover from './trigger-popover';

const TriggerItem = ({
  trigger,
  triggerNewlyCreated,
  onActions
}: {
  trigger: TriggerItemType;
  triggerNewlyCreated?: TriggerItemType;
  onActions: (action: TriggerAction) => void;
}) => {
  const { t } = useTranslation(['table', 'common', 'message', 'form']);
  const formatDateTime = useFormatDateTime();
  const { notify } = useToast();
  const { flagTrigger, url } = trigger;
  const {
    updatedAt,
    description,
    action,
    triggerCount,
    lastTriggeredAt,
    disabled,
    id
  } = flagTrigger;
  const isOnFlag = useMemo(() => action === TriggerActionType.ON, [action]);

  const isNewlyCreated = useMemo(
    () => triggerNewlyCreated?.flagTrigger?.id === id,
    [id, triggerNewlyCreated]
  );

  const handleCopy = useCallback((url: string) => {
    copyToClipBoard(url);
    notify({
      message: t('message:copied')
    });
  }, []);

  return (
    <div className="flex w-full min-w-fit gap-x-3 p-6 bg-white border border-gray-400 rounded-lg">
      <Icon icon={IconWebhook} />
      <div className="flex flex-col flex-1 gap-y-4">
        <div className="flex items-center w-full justify-between">
          <div className="flex items-center gap-x-2">
            <p className="typo-para-medium text-gray-700">
              {t('trigger.generic-trigger')}
            </p>
            <FeatureFlagStatus
              status={t(
                !disabled ? 'form:experiments.on' : 'form:experiments.off'
              )}
              enabled={!disabled}
            />
          </div>
          <div className="flex items-center gap-x-4">
            <div className="flex items-center gap-x-1.5">
              <Icon icon={IconWatch} size={'xxs'} />
              <DateTooltip
                trigger={
                  <div className="text-gray-500 typo-para-small whitespace-nowrap">
                    {Number(updatedAt) === 0 ? (
                      t('never')
                    ) : (
                      <Trans
                        i18nKey={'common:time-updated'}
                        values={{
                          time: formatDateTime(updatedAt)
                        }}
                      />
                    )}
                  </div>
                }
                date={Number(updatedAt) === 0 ? null : updatedAt}
              />
            </div>
            <TriggerPopover
              trigger={trigger.flagTrigger}
              onActions={onActions}
            />
          </div>
        </div>
        {description && (
          <p className="typo-para-medium text-gray-500">{description}</p>
        )}

        <Divider className="border-gray-300" />
        {isNewlyCreated ? (
          <div className="flex flex-col gap-y-3">
            <p className="typo-para-medium text-gray-500 uppercase">
              {t('trigger.trigger-url')}
            </p>
            <div className="flex items-center gap-x-2">
              <Input value={triggerNewlyCreated?.url} readOnly />
              <Button
                variant={'secondary-2'}
                className="px-3"
                onClick={() => handleCopy(triggerNewlyCreated?.url || url)}
              >
                <Icon icon={IconCopy} size="sm" />
              </Button>
            </div>
            <div className="flex items-center gap-x-2">
              <Icon icon={IconToastWarning} size="xxs" />
              <p className="typo-para-small text-accent-orange-600">
                {t('message:copy-trigger-url')}
              </p>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-12 w-full gap-x-4">
            <div className="flex flex-col gap-y-3 col-span-2">
              <p className="typo-para-medium text-gray-500 uppercase">
                {t('trigger.flag-target')}
              </p>
              <p className="typo-para-medium text-gray-700">
                {t(`trigger.${isOnFlag ? 'turn-on-flag' : 'turn-off-flag'}`)}
              </p>
            </div>
            <div className="flex flex-col gap-y-3 col-span-4">
              <p className="typo-para-medium text-gray-500 uppercase">
                {t('trigger.trigger-url')}
              </p>
              <p className="typo-para-medium text-primary-500 truncate">
                {url}
              </p>
            </div>
            <div className="flex flex-col gap-y-3 col-span-3">
              <p className="typo-para-medium text-gray-500 uppercase">
                {t('trigger.triggered-times')}
              </p>
              <p className="typo-para-medium text-gray-700">{triggerCount}</p>
            </div>
            <div className="flex flex-col gap-y-3 col-span-3">
              <p className="typo-para-medium text-gray-500 uppercase">
                {t('trigger.last-triggered')}
              </p>
              <p className="typo-para-medium text-gray-700">
                {Number(lastTriggeredAt) === 0
                  ? '-'
                  : formatDateTime(lastTriggeredAt)}
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default TriggerItem;
