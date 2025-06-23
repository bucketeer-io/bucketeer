import { useCallback, useState } from 'react';
import { Trans } from 'react-i18next';
import { useLocation, useNavigate } from 'react-router-dom';
import { autoOpsCreator } from '@api/auto-ops';
import { featureUpdater } from '@api/features';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures } from '@queries/features';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURE_CLONE, PAGE_PATH_FEATURES } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import { Feature, FeatureUpdaterParams } from '@types';
import ConfirmationRequiredModal, {
  ConfirmRequiredValues
} from 'pages/feature-flag-details/elements/confirm-required-modal';
import { getFlagStatus } from './collection-layout/elements/utils';
import AddFlagModal from './flags-modal/add-flag-modal';
import ArchiveModal from './flags-modal/archive-modal';
import CloneFlagModal from './flags-modal/clone-flag-modal';
import PageContent from './page-content';
import { FeatureActivityStatus, FlagActionType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'table', 'message']);
  const queryClient = useQueryClient();
  const { notify, errorNotify } = useToast();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    id: flagId,
    isAdd,
    isClone,
    onCloseActionModal,
    onOpenAddModal
  } = useActionWithURL({
    idKey: 'flagId',
    closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}`
  });

  const navigate = useNavigate();
  const location = useLocation();
  const [selectedFlag, setSelectedFlag] = useState<Feature>();
  const [isArchiving, setIsArchiving] = useState(false);
  const [isEnabling, setIsEnabling] = useState(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const [
    openConfirmRequiredModal,
    onOpenConfirmRequiredModal,
    onCloseConfirmRequiredModal
  ] = useToggleOpen(false);

  const onHandleActions = (flag: Feature, type: FlagActionType) => {
    if (type === 'CLONE') {
      return navigate(
        `${location.pathname}${PAGE_PATH_FEATURE_CLONE}/${flag.id}`
      );
    }
    setSelectedFlag(flag);
    if (['ARCHIVE', 'UNARCHIVE'].includes(type)) {
      setIsArchiving(type === 'ARCHIVE');
      return onOpenConfirmModal();
    }
    if (['ACTIVE', 'INACTIVE'].includes(type)) {
      setIsEnabling(type === 'ACTIVE');
      onOpenConfirmRequiredModal();
    }
  };

  const mutation = useMutation({
    mutationFn: async (params: Partial<FeatureUpdaterParams>) => {
      return featureUpdater(params);
    },
    onSuccess: () => {
      onCloseConfirmModal();
      notify({
        message: t('message:collection-action-success', {
          collection: t('source-type.feature-flag'),
          action: t('updated')
        })
      });
      invalidateFeatures(queryClient);
      invalidateFeature(queryClient);
      mutation.reset();
    },
    onError: error => errorNotify(error)
  });

  const handleUpdateFeature = useCallback(
    async (params: Partial<FeatureUpdaterParams>) => {
      if (selectedFlag?.id) {
        mutation.mutate({
          id: selectedFlag.id,
          environmentId: currentEnvironment.id,
          ...params
        });
      }
    },
    [selectedFlag]
  );

  const handleToggleFeatureState = useCallback(
    async (additionalValues?: ConfirmRequiredValues) => {
      try {
        if (selectedFlag) {
          const { scheduleType, comment, scheduleAt } = additionalValues || {};
          let resp;
          if (['ENABLE', 'DISABLE'].includes(scheduleType as string)) {
            resp = await featureUpdater({
              id: selectedFlag.id,
              environmentId: currentEnvironment.id,
              enabled: !selectedFlag.enabled,
              comment
            });
          } else {
            resp = await autoOpsCreator({
              environmentId: currentEnvironment.id,
              featureId: selectedFlag.id,
              opsType: 'SCHEDULE',
              datetimeClauses: [
                {
                  actionType: selectedFlag.enabled ? 'DISABLE' : 'ENABLE',
                  time: scheduleAt as string
                }
              ]
            });
          }
          if (resp) {
            notify({
              message: t('message:collection-action-success', {
                collection: t('source-type.feature-flag'),
                action: t('updated')
              })
            });
            invalidateFeatures(queryClient);
            onCloseConfirmRequiredModal();
          }
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [isEnabling, selectedFlag, currentEnvironment]
  );

  return (
    <>
      <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      {openConfirmModal && selectedFlag && (
        <ArchiveModal
          isArchiving={isArchiving}
          isOpen={openConfirmModal}
          isLoading={mutation.isPending}
          isShowWarning={
            isArchiving &&
            getFlagStatus(selectedFlag) === FeatureActivityStatus.ACTIVE
          }
          onClose={onCloseConfirmModal}
          className="py-5"
          title={
            isArchiving
              ? t(`table:popover.archive-flag`)
              : t(`table:popover.unarchive-flag`)
          }
          description={
            <Trans
              i18nKey={
                isArchiving
                  ? 'table:feature-flags.confirm-archive-desc'
                  : 'table:feature-flags.confirm-unarchive-desc'
              }
              values={{ name: selectedFlag?.name }}
              components={{ text: <span /> }}
            />
          }
          onSubmit={({ comment }) =>
            handleUpdateFeature({
              archived: isArchiving,
              comment
            })
          }
        />
      )}
      {isAdd && <AddFlagModal isOpen={isAdd} onClose={onCloseActionModal} />}
      {isClone && flagId && (
        <CloneFlagModal
          flagId={flagId}
          isOpen={isClone}
          onClose={onCloseActionModal}
        />
      )}
      {openConfirmRequiredModal && selectedFlag && (
        <ConfirmationRequiredModal
          isOpen={openConfirmRequiredModal}
          feature={selectedFlag}
          isShowScheduleSelect={true}
          isShowRolloutWarning={selectedFlag.enabled}
          onClose={onCloseConfirmRequiredModal}
          onSubmit={handleToggleFeatureState}
        />
      )}
    </>
  );
};

export default PageLoader;
