import { useState } from 'react';
import { userSegmentBulkDownload } from '@api/user-segment';
import { userSegmentDelete } from '@api/user-segment/user-segment-delete';
import { invalidateUserSegments } from '@queries/user-segments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import PageContent from './page-content';
import AddUserSegmentModal from './user-segment-modal/add-segment-modal';
import DeleteUserSegmentModal from './user-segment-modal/delete-segment-modal';
import EditUserSegmentModal from './user-segment-modal/edit-segment-modal';
import FlagsConnectedModal from './user-segment-modal/flags-connected-modal';

const PageLoader = () => {
  const { t } = useTranslation(['common', 'message']);
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
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
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
        message: t('message:collection-action-success', {
          collection: t('source-type.segment'),
          action: t('deleted').toLowerCase()
        })
      });
      mutation.reset();
    },
    onError: error => errorNotify(error)
  });

  const onDeleteSegment = () => {
    if (selectedSegment) {
      mutation.mutate(selectedSegment);
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

  return (
    <>
      <PageContent
        segmentUploading={segmentUploading}
        onAdd={onOpenAddModal}
        onEdit={segment => {
          setSelectedSegment(segment);
          onOpenEditModal();
        }}
        onOpenFlagModal={segment => {
          setSelectedSegment(segment);
          onOpenFlagModal();
        }}
        onDelete={segment => {
          setSelectedSegment(segment);
          onOpenDeleteModal();
        }}
        onDownload={onBulkDownloadSegment}
      />
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
