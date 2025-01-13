import { useState } from 'react';
import { userSegmentDelete } from '@api/user-segment/user-segment-delete';
import { invalidateUserSegments } from '@queries/user-segments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { UserSegment } from '@types';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import { useFetchSegments } from './collection-loader/use-fetch-segment';
import PageContent from './page-content';
import AddUserSegmentModal from './user-segment-modal/add-segment-modal';
import DeleteUserSegmentModal from './user-segment-modal/delete-segment-modal';
import EditUserSegmentModal from './user-segment-modal/edit-segment-modal';
import FlagsConnectedModal from './user-segment-modal/flags-connected-modal';

const PageLoader = () => {
  const [selectedSegment, setSelectedSegment] = useState<UserSegment>();
  const [segmentUploading, setSegmentUploading] = useState<UserSegment | null>(
    null
  );

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);
  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);
  const [isOpenFlagModal, onOpenFlagModal, onCloseFlagModal] =
    useToggleOpen(false);
  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);

  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);
  const { notify } = useToast();
  const queryClient = useQueryClient();

  const {
    data: collection,
    isLoading,
    refetch,
    isError
  } = useFetchSegments({
    pageSize: 1,
    environmentId: currenEnvironment.id
  });

  const isEmpty = collection?.segments.length === 0;

  const mutation = useMutation({
    mutationFn: async (selectedSegment: UserSegment) => {
      return userSegmentDelete({
        id: selectedSegment.id,
        environmentId: currenEnvironment.id
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
        <PageContent
          segmentUploading={segmentUploading}
          onAdd={onOpenAddModal}
          onEdit={value => {
            setSelectedSegment(value);
            onOpenEditModal();
          }}
          onOpenFlagModal={value => {
            setSelectedSegment(value);
            onOpenFlagModal();
          }}
          onDelete={value => {
            setSelectedSegment(value);
            onOpenDeleteModal();
          }}
        />
      )}
      {isOpenAddModal && (
        <AddUserSegmentModal
          isOpen={isOpenAddModal}
          onClose={onCloseAddModal}
        />
      )}
      {isOpenEditModal && selectedSegment && (
        <EditUserSegmentModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          userSegment={selectedSegment!}
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
