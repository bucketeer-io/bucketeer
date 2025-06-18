import { useEffect, useState } from 'react';
import { userSegmentBulkDownload } from '@api/user-segment';
import { userSegmentDelete } from '@api/user-segment/user-segment-delete';
import { useQueryUserSegment } from '@queries/user-segment-details';
import { invalidateUserSegments } from '@queries/user-segments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_USER_SEGMENTS } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { UserSegment } from '@types';
import PageContent from './page-content';
import { UserSegmentsActionsType } from './types';
import AddUserSegmentModal from './user-segment-modal/add-segment-modal';
import DeleteUserSegmentModal from './user-segment-modal/delete-segment-modal';
import EditUserSegmentModal from './user-segment-modal/edit-segment-modal';
import FlagsConnectedModal from './user-segment-modal/flags-connected-modal';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

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
      return userSegmentDelete({
        id: selectedSegment.id,
        environmentId: currentEnvironment.id
      });
    },
    onSuccess: () => {
      onCloseDeleteModal();
      invalidateUserSegments(queryClient);
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{selectedSegment?.name}</b>
            {` has been deleted successfully!`}
          </span>
        )
      });
      mutation.reset();
    }
  });

  const onDeleteSegment = () => {
    if (selectedSegment) {
      mutation.mutate(selectedSegment);
    }
  };

  const onActionHandler = (
    segment: UserSegment,
    type: UserSegmentsActionsType
  ) => {
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
  };

  const onBulkDownloadSegment = async (segment: UserSegment) => {
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
  };

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
        segmentUploading={segmentUploading}
        onAdd={onOpenAddModal}
        onActionHandler={onActionHandler}
      />
      {isAdd && (
        <AddUserSegmentModal isOpen={isAdd} onClose={onCloseActionModal} />
      )}
      {isEdit && (
        <EditUserSegmentModal
          isOpen={isEdit}
          isLoadingSegment={isLoadingSegment}
          userSegment={selectedSegment!}
          onClose={onCloseActionModal}
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
