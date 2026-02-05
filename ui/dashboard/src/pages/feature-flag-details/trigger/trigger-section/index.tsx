import { Fragment, useCallback, useMemo, useRef, useState } from 'react';
import { Trans } from 'react-i18next';
import { IconLaunchOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { triggerDelete } from '@api/trigger/triggers-delete';
import {
  triggerUpdate,
  TriggerUpdateParams
} from '@api/trigger/triggers-update';
import { invalidateTriggers, useQueryTriggers } from '@queries/triggers';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { Feature, TriggerItemType } from '@types';
import { IconPlus } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import Card from 'elements/card';
import ConfirmModal from 'elements/confirm-modal';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import FormLoading from 'elements/form-loading';
import CreateTriggerForm from '../create-trigger-form';
import { TriggerAction } from '../types';
import TriggerItem from './trigger-item';

interface ActionState {
  action?: TriggerAction;
  trigger?: TriggerItemType;
}

const TriggerList = ({
  feature,
  editable
}: {
  feature: Feature;
  editable: boolean;
}) => {
  const { t } = useTranslation(['table', 'message', 'common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const formRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();
  const { notify, errorNotify } = useToast();

  const [triggerNewlyCreated, setTriggerNewlyCreated] = useState<
    TriggerItemType | undefined
  >(undefined);
  const [isShowCreateForm, setIsShowCreateForm] = useState(false);
  const [actionState, setActionState] = useState<ActionState>({
    action: undefined,
    trigger: undefined
  });

  const { data: triggerCollection, isLoading } = useQueryTriggers({
    params: {
      environmentId: currentEnvironment.id,
      featureId: feature.id,
      cursor: String(0)
    }
  });

  const isDisabledCreate = useMemo(() => !editable, [editable]);

  const triggers = triggerCollection?.flagTriggers || [];
  const { EDIT, RESET, DISABLE, ENABLE, DELETE } = TriggerAction;
  const isEdit = useMemo(() => actionState?.action === EDIT, [actionState]);

  const isReset = useMemo(() => actionState?.action === RESET, [actionState]);
  const isDisable = useMemo(
    () => actionState?.action === DISABLE,
    [actionState]
  );
  const isEnable = useMemo(() => actionState?.action === ENABLE, [actionState]);
  const isDelete = useMemo(() => actionState?.action === DELETE, [actionState]);

  const confirmModalTitle = useMemo(() => {
    const key = isReset
      ? 'reset-trigger-url'
      : isDisable
        ? 'disable-trigger'
        : isEnable
          ? 'enable-trigger'
          : 'delete-trigger';
    return t(`trigger.${key}`);
  }, [isReset, isDisable, isEnable]);

  const confirmModalDesc = useMemo(() => {
    const key = isReset
      ? 'reset'
      : isDisable
        ? 'disable'
        : isEnable
          ? 'enable'
          : 'delete';
    return t(`trigger.${key}-trigger-desc`);
  }, [isReset, isDisable, isEnable]);

  const mutationState = useMutation({
    mutationFn: async (params: TriggerUpdateParams) => {
      const { id, environmentId, ...rest } = params || {};
      return isDelete
        ? await triggerDelete({
            id,
            environmentId
          })
        : await triggerUpdate({
            id,
            environmentId,
            ...rest
          });
    },
    onSuccess: (data, params) => {
      if (params?.reset) {
        const trigger = triggers.find(
          item => item.flagTrigger.id === params.id
        );
        if (trigger)
          setTriggerNewlyCreated({
            ...trigger,
            url: data.url
          });
      }
      notify({
        message: t('message:collection-action-success', {
          collection: t('feature-flags.trigger'),
          action: t(isDelete ? 'common:delete' : 'common:updated')
        })
      });
      onReset();
      invalidateTriggers(queryClient);
      mutationState.reset();
    },
    onError: error => errorNotify(error)
  });

  const onActions = useCallback(
    (trigger: TriggerItemType, action: TriggerAction) => {
      setActionState({
        action,
        trigger
      });
    },
    [mutationState, currentEnvironment]
  );

  const onReset = useCallback(() => {
    setActionState({
      action: undefined,
      trigger: undefined
    });
    setIsShowCreateForm(false);
  }, []);

  return (
    <Card className="gap-y-6">
      <p className="typo-head-bold-small text-gray-800">{t('trigger.title')}</p>
      <div className="inline typo-para-small text-gray-500 gap-x-1">
        <Trans
          i18nKey={'table:trigger.description'}
          components={{
            comp: (
              <Link
                to={DOCUMENTATION_LINKS.FLAG_TRIGGER}
                target="_blank"
                className="flex items-center gap-x-1 text-primary-500 underline"
              />
            ),
            icon: <Icon icon={IconLaunchOutlined} size="sm" />
          }}
        />
      </div>
      {isLoading ? (
        <FormLoading />
      ) : (
        <>
          {triggers.map((trigger, index) => (
            <Fragment key={index}>
              {isEdit &&
              actionState?.trigger?.flagTrigger?.id ===
                trigger?.flagTrigger?.id ? (
                <CreateTriggerForm
                  disabled={isDisabledCreate}
                  ref={formRef}
                  selectedTrigger={actionState?.trigger}
                  featureId={feature.id}
                  environmentId={currentEnvironment.id}
                  onCancel={onReset}
                  setTriggerNewlyCreated={setTriggerNewlyCreated}
                />
              ) : (
                <TriggerItem
                  disabledAction={isDisabledCreate}
                  trigger={trigger}
                  triggerNewlyCreated={triggerNewlyCreated}
                  onActions={action => onActions(trigger, action)}
                />
              )}
            </Fragment>
          ))}
          {isShowCreateForm && !isEdit ? (
            <CreateTriggerForm
              disabled={isDisabledCreate}
              featureId={feature.id}
              environmentId={currentEnvironment.id}
              onCancel={onReset}
              setTriggerNewlyCreated={setTriggerNewlyCreated}
            />
          ) : (
            <DisabledButtonTooltip
              align="start"
              hidden={!isDisabledCreate}
              trigger={
                <Button
                  variant="text"
                  className="h-8 w-fit p-0"
                  disabled={isDisabledCreate}
                  onClick={() => {
                    setActionState({
                      action: undefined,
                      trigger: undefined
                    });
                    setIsShowCreateForm(true);
                  }}
                >
                  <Icon icon={IconPlus} size="md" />
                  {t('trigger.add-trigger')}
                </Button>
              }
            />
          )}
        </>
      )}

      {!isEdit && !!actionState?.action && !!actionState?.trigger && (
        <ConfirmModal
          isOpen={!isEdit && !!actionState?.action && !!actionState?.trigger}
          title={confirmModalTitle}
          description={confirmModalDesc}
          loading={mutationState.isPending}
          onClose={onReset}
          onSubmit={() =>
            mutationState.mutate({
              id: actionState.trigger!.flagTrigger.id,
              environmentId: currentEnvironment.id,
              reset: isReset,
              disabled: isReset
                ? actionState.trigger!.flagTrigger.disabled
                : isDisable
            })
          }
        />
      )}
    </Card>
  );
};

export default TriggerList;
