import { useCallback, useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { userSegmentBulkDownload } from '@api/user-segment';
import { userSegmentDelete } from '@api/user-segment/user-segment-delete';
import { useMutation } from '@tanstack/react-query';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_NEW, PAGE_PATH_USER_SEGMENTS } from 'constants/routing';
import { useToast } from 'hooks';
import { useToggleOpen } from 'hooks/use-toggle-open';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import PageContent from './page-content';
import { UserSegmentsActionsType } from './types';
import DeleteUserSegmentModal from './user-segment-modal/delete-segment-modal';
import FlagsConnectedModal from './user-segment-modal/flags-connected-modal';
import SegmentUploadingModal from './user-segment-modal/segment-uploading-modal';

// How long the freshly saved segment keeps its local "uploading" indicator
// before we rely solely on the status reported by the server.
const UPLOADING_INDICATOR_TIMEOUT = 10000;

const PageLoader = () => {
  const { t } = useTranslation(['common', 'message']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);
  const navigate = useNavigate();
  const location = useLocation();

  const [isOpenFlagModal, onOpenFlagModal, onCloseFlagModal] =
    useToggleOpen(false);
  const [isOpenDeleteModal, onOpenDeleteModal, onCloseDeleteModal] =
    useToggleOpen(false);
  const [isOpenUploadingModal, onOpenUploadingdModal, onCloseUploadingdModal] =
    useToggleOpen(false);
  const [selectedSegment, setSelectedSegment] = useState<UserSegment>();
  // The segment form page passes the segment being uploaded via navigation
  // state, so the list can show the uploading indicator right away.
  const [segmentUploading, setSegmentUploading] = useState<UserSegment | null>(
    (location.state as { segmentUploading?: UserSegment } | null)
      ?.segmentUploading || null
  );

  const { notify, errorNotify } = useToast();

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

  const onHandleAddSegment = () => {
    navigate(
      `/${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}${PAGE_PATH_NEW}`
    );
  };

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
          return navigate(
            `/${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}/${segment.id}`
          );
        case 'FLAG':
          return onOpenFlagModal();
        case 'DELETE':
          return onOpenDeleteModal();
        case 'UPLOADING':
          return onOpenUploadingdModal();
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
    if (segmentUploading) {
      const timerId = setTimeout(
        () => setSegmentUploading(null),
        UPLOADING_INDICATOR_TIMEOUT
      );
      return () => clearTimeout(timerId);
    }
  }, [segmentUploading]);

  return (
    <>
      <PageContent
        editable={editable}
        segmentUploading={segmentUploading}
        onAdd={onHandleAddSegment}
        onActionHandler={onActionHandler}
      />

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
      {isOpenUploadingModal && (
        <SegmentUploadingModal
          isOpen={isOpenUploadingModal}
          onClose={onCloseUploadingdModal}
        />
      )}
    </>
  );
};

export default PageLoader;
