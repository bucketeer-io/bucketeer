import { useCallback, useEffect, useState } from 'react';
import { Trans } from 'react-i18next';
import { experimentUpdater } from '@api/experiment';
import { invalidateExperiments } from '@queries/experiments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchExperiments } from './collection-loader/use-fetch-experiment';
import AddExperimentModal from './experiments-modal/add-experiment-modal';
import EditExperimentModal from './experiments-modal/edit-experiment-modal';
import PageContent from './page-content';
import { ExperimentActionsType } from './types';

const PageLoader = ({
  setTotalCount
}: {
  setTotalCount: (value: string | number) => void;
}) => {
  const { t } = useTranslation(['common', 'table']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { notify } = useToast();

  const [selectedExperiment, setSelectedExperiment] = useState<Experiment>();
  const [isArchiving, setIsArchiving] = useState<boolean>();

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const { isAdd, isEdit, onOpenAddModal, onOpenEditModal, onCloseActionModal } =
    useActionWithURL({
      closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}`
    });

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchExperiments({ environmentId: currentEnvironment.id });

  const isEmpty = collection?.experiments?.length === 0;

  const mutation = useMutation({
    mutationFn: async (id: string) => {
      return experimentUpdater({
        id,
        archived: isArchiving,
        environmentId: currentEnvironment.id
      });
    },
    onSuccess: () => {
      onCloseConfirmModal();
      invalidateExperiments(queryClient);
      mutation.reset();
    },
    onError: error => {
      notify({
        toastType: 'toast',
        messageType: 'error',
        message: error?.message || 'Something went wrong.'
      });
    }
  });

  const onHandleArchive = () => {
    if (selectedExperiment?.id) {
      mutation.mutate(selectedExperiment?.id);
    }
  };

  const onHandleActions = useCallback(
    (item: Experiment, type: ExperimentActionsType) => {
      if (type === 'EDIT') {
        return onOpenEditModal(
          `/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${item.id}`
        );
      }
      setSelectedExperiment(item);
      if (type === 'ARCHIVE') {
        setIsArchiving(true);
        return onOpenConfirmModal();
      }
      setIsArchiving(false);
      return onOpenConfirmModal();
    },
    []
  );

  useEffect(() => {
    if (collection?.experiments) setTotalCount(collection?.experiments?.length);
  }, [collection]);

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : isEmpty ? (
        <PageLayout.EmptyState>
          <EmptyCollection onAdd={onOpenAddModal} />
        </PageLayout.EmptyState>
      ) : (
        <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      )}
      {isAdd && (
        <AddExperimentModal isOpen={isAdd} onClose={onCloseActionModal} />
      )}
      {isEdit && (
        <EditExperimentModal isOpen={isEdit} onClose={onCloseActionModal} />
      )}

      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          onClose={onCloseConfirmModal}
          onSubmit={onHandleArchive}
          title={
            isArchiving
              ? t(`table:popover.archive-experiment`)
              : t(`table:popover.unarchive-experiment`)
          }
          description={
            <Trans
              i18nKey={
                isArchiving
                  ? 'table:experiment.confirm-archive-desc'
                  : 'table:experiment.confirm-unarchive-desc'
              }
              values={{ name: selectedExperiment?.name }}
              components={{ bold: <strong /> }}
            />
          }
          loading={mutation.isPending}
        />
      )}
    </>
  );
};

export default PageLoader;
