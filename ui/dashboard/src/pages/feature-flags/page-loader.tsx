import { useCallback, useState } from 'react';
import { Trans } from 'react-i18next';
import { useLocation, useNavigate } from 'react-router-dom';
import { featureUpdater } from '@api/features';
import { invalidateFeatures } from '@queries/features';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURE_CLONE, PAGE_PATH_FEATURES } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import { Feature, FeatureUpdaterParams } from '@types';
import { getFlagStatus } from './collection-layout/elements/utils';
import AddFlagModal from './flags-modal/add-flag-modal';
import ArchiveModal from './flags-modal/archive-modal';
import CloneFlagModal from './flags-modal/clone-flag-modal';
import ConfirmationRequiredModal from './flags-modal/confirm-required-modal';
import PageContent from './page-content';
import { FeatureActivityStatus, FlagActionType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'table']);
  const queryClient = useQueryClient();
  const { notify } = useToast();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    id: flagId,
    isAdd,
    isClone,
    onCloseActionModal,
    onOpenAddModal,
    errorToast
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
        message: 'Updated feature flag successfully.'
      });
      invalidateFeatures(queryClient);
      mutation.reset();
    },
    onError: error => errorToast(error)
  });

  const handleUpdateFeature = useCallback(
    (params: Partial<FeatureUpdaterParams>) => {
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

  return (
    <>
      <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      {openConfirmModal && selectedFlag && (
        <ArchiveModal
          isArchiving={isArchiving}
          isOpen={openConfirmModal}
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
          onSubmit={() =>
            handleUpdateFeature({
              archived: isArchiving
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
          errorToast={errorToast}
        />
      )}
      {openConfirmRequiredModal && selectedFlag && (
        <ConfirmationRequiredModal
          isEnabling={isEnabling}
          selectedFlag={selectedFlag}
          isOpen={openConfirmRequiredModal}
          onClose={onCloseConfirmRequiredModal}
        />
      )}
    </>
  );
};

export default PageLoader;
