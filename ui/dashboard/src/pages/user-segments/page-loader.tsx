import { useState } from 'react';
// import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
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

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);
  const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
    useToggleOpen(false);
  const [isOpenFlagModal, onOpenFlagModal, onCloseFlagModal] =
    useToggleOpen(false);
  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);

  // const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

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
      {isOpenEditModal && (
        <EditUserSegmentModal
          isOpen={isOpenEditModal}
          onClose={onCloseEditModal}
          userSegment={selectedSegment!}
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
          onSubmit={() => {}}
        />
      )}
    </>
  );
};

export default PageLoader;
