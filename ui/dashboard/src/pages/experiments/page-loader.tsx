import { useCallback, useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { experimentUpdater, ExperimentUpdaterParams } from '@api/experiment';
import { invalidateExperiments } from '@queries/experiments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { isNotEmptyObject } from 'utils/data-type';
import { stringifyParams, useSearchParams } from 'utils/search-params';
import ConfirmModal from 'elements/confirm-modal';
import PageLayout from 'elements/page-layout';
import { useFetchExperiments } from './collection-loader/use-fetch-experiment';
import ExperimentCreateUpdateModal from './experiments-modal/experiment-create-update';
import GoalsConnectionModal from './experiments-modal/goals-connection-modal';
import PageContent from './page-content';
import { ExperimentActionsType } from './types';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'table', 'message']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);
  const { notify, errorNotify } = useToast();
  const { searchOptions } = useSearchParams();

  const queryString = useMemo(
    () =>
      isNotEmptyObject(searchOptions)
        ? `?${decodeURIComponent(stringifyParams(searchOptions))}`
        : '',
    [searchOptions]
  );

  const [selectedExperiment, setSelectedExperiment] = useState<Experiment>();
  const [isArchiving, setIsArchiving] = useState<boolean>();
  const [isStop, setIsStop] = useState<boolean>();

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);
  const [openGoalsModal, onOpenGoalsModal, onCloseGoalsModal] =
    useToggleOpen(false);
  const [
    openToggleExperimentModal,
    onOpenToggleExperimentModal,
    onCloseToggleExperimentModal
  ] = useToggleOpen(false);

  const { isAdd, isEdit, onOpenAddModal, onOpenEditModal, onCloseActionModal } =
    useActionWithURL({
      closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}${queryString}`
    });

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchExperiments({
    environmentId: currentEnvironment.id
  });

  const summary = useMemo(() => collection?.summary, [collection]);

  const mutation = useMutation({
    mutationFn: async (params: ExperimentUpdaterParams) => {
      return experimentUpdater(params);
    },
    onSuccess: () => {
      onCloseConfirmModal();
      onCloseToggleExperimentModal();
      invalidateExperiments(queryClient);
      mutation.reset();
      notify({
        message: t('message:collection-action-success', {
          collection: t('source-type.experiment'),
          action: t('updated')
        })
      });
    },
    onError: error => errorNotify(error)
  });

  const onHandleArchive = () => {
    if (selectedExperiment?.id) {
      mutation.mutate({
        id: selectedExperiment?.id,
        archived: isArchiving,
        environmentId: currentEnvironment.id
      });
    }
  };

  const onToggleExperiment = () => {
    if (selectedExperiment?.id) {
      mutation.mutate({
        id: selectedExperiment?.id,
        environmentId: currentEnvironment.id,
        startAt: selectedExperiment.startAt,
        stopAt: selectedExperiment.stopAt,
        status: {
          status: isStop ? 'FORCE_STOPPED' : 'RUNNING'
        }
      });
    }
  };

  const onHandleActions = useCallback(
    (item: Experiment, type: ExperimentActionsType) => {
      if (type === 'EDIT') {
        return onOpenEditModal(
          `/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${item.id}${queryString}`
        );
      }
      setSelectedExperiment(item);
      if (type === 'GOALS-CONNECTION') {
        return onOpenGoalsModal();
      }
      if (['START', 'STOP'].includes(type)) {
        setIsStop(type === 'STOP');
        return onOpenToggleExperimentModal();
      }
      setIsArchiving(type === 'ARCHIVE');
      return onOpenConfirmModal();
    },
    [queryString]
  );

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        <PageContent
          disabled={!editable}
          summary={summary}
          onAdd={onOpenAddModal}
          onHandleActions={onHandleActions}
        />
      )}
      {(!!isAdd || !!isEdit) && (
        <ExperimentCreateUpdateModal
          disabled={!editable}
          isOpen={!!isAdd || !!isEdit}
          onClose={onCloseActionModal}
        />
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
      {openToggleExperimentModal && (
        <ConfirmModal
          isOpen={openToggleExperimentModal}
          onClose={onCloseToggleExperimentModal}
          onSubmit={onToggleExperiment}
          title={
            isStop
              ? t(`table:popover.stop-experiment`)
              : t(`table:popover.start-experiment`)
          }
          description={
            <Trans
              i18nKey={
                isStop
                  ? 'table:experiment.confirm-stop-desc'
                  : 'table:experiment.confirm-start-desc'
              }
              values={{ name: selectedExperiment?.name }}
              components={{ bold: <strong /> }}
            />
          }
          loading={mutation.isPending}
          disabled={!editable}
        />
      )}
      {openGoalsModal && selectedExperiment && (
        <GoalsConnectionModal
          isOpen={openGoalsModal}
          experiment={selectedExperiment}
          onClose={onCloseGoalsModal}
        />
      )}
    </>
  );
};

export default PageLoader;
