import { useCallback, useEffect, useState } from 'react';
import { userSegmentBulkDownload } from '@api/user-segment';
import { userSegmentDelete } from '@api/user-segment/user-segment-delete';
import { useQueryUserSegment } from '@queries/user-segment-details';
import { invalidateUserSegments } from '@queries/user-segments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_USER_SEGMENTS } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import PageContent from './page-content';
import { UserSegmentsActionsType } from './types';
import DeleteUserSegmentModal from './user-segment-modal/delete-segment-modal';
import FlagsConnectedModal from './user-segment-modal/flags-connected-modal';
import SegmentCreateUpdateModal from './user-segment-modal/segment-create-update-form';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'message']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  const {
    id,
    isAdd,
    isEdit,
    onOpenAddModal,
    onOpenEditModal,
    onCloseActionModal
  } = useActionWithURL({
    closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}`
  });

  const [isOpenFlagModal, onOpenFlagModal, onCloseFlagModal] =
    useToggleOpen(false);
  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);
  const [selectedSegment, setSelectedSegment] = useState<UserSegment>();
  const [segmentUploading, setSegmentUploading] = useState<UserSegment | null>(
    null
  );

  const {
    data: segmentCollection,
    isLoading: isLoadingSegment,
    isError,
    error
  } = useQueryUserSegment({
    params: {
      environmentId: currentEnvironment.id,
      id: id as string
    },
    enabled: !!isEdit && !!id && !selectedSegment
  });

  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (selectedSegment: UserSegment) => {
      if (editable)
        return userSegmentDelete({
          id: selectedSegment.id,
          environmentId: currentEnvironment.id
        });
    },
    onSuccess: () => {
      onCloseDeleteModal();
      invalidateUserSegments(queryClient);
      notify({
        message: t('message:collection-action-success', {
          collection: t('source-type.segment'),
          action: t('deleted')
        })
      });
      mutation.reset();
    },
    onError: error => errorNotify(error)
  });

  const onDeleteSegment = () => {
    if (selectedSegment && editable) {
      mutation.mutate(selectedSegment);
    }
  };

  const onActionHandler = useCallback(
    (segment: UserSegment, type: UserSegmentsActionsType) => {
      if (type !== 'DOWNLOAD') setSelectedSegment(segment);
      switch (type) {
        case 'EDIT':
          return onOpenEditModal(
            `/${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}/${segment.id}`
          );
        case 'FLAG':
          return onOpenFlagModal();
        case 'DELETE':
          return onOpenDeleteModal();
        default:
          return onBulkDownloadSegment(segment);
      }
    },
    [currentEnvironment]
  );

  const onBulkDownloadSegment = useCallback(
    async (segment: UserSegment) => {
      const resp = await userSegmentBulkDownload({
        segmentId: segment.id,
        environmentId: currentEnvironment.id
      });
      if (resp.data) {
        const url = window.URL.createObjectURL(
          new Blob([atob(String(resp.data))])
        );
        const link = window.document.createElement('a');
        link.href = url;
        link.setAttribute(
          'download',
          `${currentEnvironment.name}-${segment.name}.csv`
        );
        window.document.body.appendChild(link);
        link.click();
        if (link.parentNode) {
          link.parentNode.removeChild(link);
        }
      }
    },
    [currentEnvironment]
  );

  useEffect(() => {
    if (segmentCollection) {
      setSelectedSegment(segmentCollection.segment);
    }
  }, [segmentCollection]);

  useEffect(() => {
    if (isError && error) {
      errorNotify(error);
      onCloseActionModal();
    }
  }, [isError, error]);

  return (
    <>
      <PageContent
        editable={editable}
        segmentUploading={segmentUploading}
        onAdd={onOpenAddModal}
        onActionHandler={onActionHandler}
      />

      {(!!isAdd || !!isEdit) && (
        <SegmentCreateUpdateModal
          isUpdate={!!isEdit || !!selectedSegment || isLoadingSegment}
          isDisabled={!editable}
          isOpen={!!isAdd || !!isEdit}
          isLoadingSegment={isLoadingSegment}
          userSegment={selectedSegment!}
          onClose={() => {
            setSelectedSegment(undefined);
            onCloseActionModal();
          }}
          setSegmentUploading={setSegmentUploading}
        />
      )}

      {isOpenFlagModal && selectedSegment && (
        <FlagsConnectedModal
          segment={selectedSegment}
          isOpen={isOpenFlagModal}
          onClose={onCloseFlagModal}
        />
      )}
      {isOpenDeleteModal && selectedSegment && (
        <DeleteUserSegmentModal
          isDisabled={!editable}
          isOpen={isOpenDeleteModal}
          loading={false}
          userSegment={selectedSegment}
          onClose={onCloseDeleteModal}
          onSubmit={onDeleteSegment}
        />
      )}
    </>
  );
};

export default PageLoader;
