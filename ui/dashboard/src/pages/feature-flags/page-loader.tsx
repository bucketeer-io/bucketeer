import { useState } from 'react';
import { Trans } from 'react-i18next';
import { useLocation, useNavigate } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURE_CLONE, PAGE_PATH_FEATURES } from 'constants/routing';
import { useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import AddFlagModal from './flags-modal/add-flag-modal';
import ArchiveModal from './flags-modal/archive-modal';
import ArchiveWarning from './flags-modal/archive-modal/archive-warning';
import CloneFlagModal from './flags-modal/clone-flag-modal';
import ConfirmationRequiredModal from './flags-modal/confirm-required-modal';
import PageContent from './page-content';
import { FlagActionType, FlagsTemp } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { isAdd, isClone, onCloseActionModal, onOpenAddModal } =
    useActionWithURL({
      closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}`
    });
  const isLoading = false,
    isError = false,
    isEmpty = false;

  const navigate = useNavigate();
  const location = useLocation();
  const [selectedFlag, setSelectedFlag] = useState<FlagsTemp>();
  const [isArchiving, setIsArchiving] = useState(false);

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const [
    openConfirmRequiredModal,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    onOpenConfirmRequiredModal,
    onCloseConfirmRequiredModal
  ] = useToggleOpen(false);

  const onHandleActions = (flag: FlagsTemp, type: FlagActionType) => {
    if (type === 'CLONE') {
      return navigate(
        `${location.pathname}${PAGE_PATH_FEATURE_CLONE}/${flag.id}`
      );
    }
    if (['ARCHIVE', 'UNARCHIVE'].includes(type)) {
      setIsArchiving(type === 'ARCHIVE' ? true : false);
      onOpenConfirmModal();
    }
    setSelectedFlag(flag);
  };

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={() => {}} />
      ) : isEmpty ? (
        <PageLayout.EmptyState>
          <EmptyCollection onAdd={onOpenAddModal} />
        </PageLayout.EmptyState>
      ) : (
        <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      )}
      {openConfirmModal && (
        <ArchiveModal
          isOpen={openConfirmModal}
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
          onSubmit={() => {}}
        >
          {<ArchiveWarning days={1} />}
        </ArchiveModal>
      )}
      {isAdd && <AddFlagModal isOpen={isAdd} onClose={onCloseActionModal} />}
      {isClone && (
        <CloneFlagModal isOpen={isClone} onClose={onCloseActionModal} />
      )}
      {openConfirmRequiredModal && (
        <ConfirmationRequiredModal
          isOpen={openConfirmRequiredModal}
          onClose={onCloseConfirmRequiredModal}
          onSubmit={() => {}}
        />
      )}
    </>
  );
};

export default PageLoader;
