import { useCallback } from 'react';
import {
  IconCloudDownloadOutlined,
  IconDeleteOutlined,
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { Popover } from 'components/popover';
import Spinner from 'components/spinner';
import DateTooltip from 'elements/date-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { UserSegmentsActionsType } from '../types';

export const useColumns = ({
  segmentUploading,
  onActionHandler
}: {
  segmentUploading: UserSegment | null;
  onActionHandler: (value: UserSegment, type: UserSegmentsActionsType) => void;
}): ColumnDef<UserSegment>[] => {
  const { consoleAccount } = useAuth();
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const editable = hasEditable(consoleAccount!);

  const getUploadingStatus = useCallback(
    (segment: UserSegment) => {
      if (segment.status === 'UPLOADING') return true;
      if (segmentUploading?.id === segment.id) return true;
    },
    [segmentUploading]
  );

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 350,
      cell: ({ row }) => {
        const segment = row.original;
        const { id, name } = segment;

        return (
          <div
            onClick={() => onActionHandler(segment, 'EDIT')}
            className="flex items-center gap-x-2 cursor-pointer min-w-[300px]"
          >
            <NameWithTooltip
              id={id}
              content={<NameWithTooltip.Content content={name} id={id} />}
              trigger={<NameWithTooltip.Trigger id={id} name={name} />}
            />
            {getUploadingStatus(segment) && <Spinner />}
          </div>
        );
      }
    },
    {
      accessorKey: 'users',
      header: `${t('users')}`,
      size: 200,
      cell: ({ row }) => {
        const segment = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {segment.includedUserCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'connections',
      header: `${t('connections')}`,
      size: 200,
      cell: ({ row }) => {
        const segment = row.original;
        return (
          <div
            className={cn(
              'flex-center w-fit px-2 py-1.5 rounded bg-primary-50 text-primary-500 typo-para-medium',
              {
                'cursor-pointer': segment?.features?.length
              }
            )}
            onClick={() =>
              segment?.features?.length && onActionHandler(segment, 'FLAG')
            }
          >
            {segment?.features?.length}
            {` ${segment?.features?.length === 1 ? 'Flag' : 'Flags'}`}
          </div>
        );
      }
    },
    {
      accessorKey: 'status',
      header: `${t('status')}`,
      size: 150,
      minSize: 150,
      maxSize: 150,
      cell: ({ row }) => {
        const segment = row.original;
        const isUploading = getUploadingStatus(segment);
        return (
          <div
            className={cn(
              'typo-para-small text-accent-green-500 bg-accent-green-50 px-2 py-[3px] w-fit text-center whitespace-nowrap rounded',
              {
                'bg-gray-200 text-gray-600': !segment.isInUseStatus,
                'bg-accent-orange-50 text-accent-orange-500': isUploading
              }
            )}
          >
            {isUploading
              ? 'Uploading'
              : segment.isInUseStatus
                ? 'In Use'
                : 'Not In Use'}
          </div>
        );
      }
    },
    {
      accessorKey: 'updatedAt',
      header: t('table:updated-at'),
      size: 200,
      cell: ({ row }) => {
        const segment = row.original;
        const isNever = Number(segment.updatedAt) === 0;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {Number(segment.updatedAt) === 0
                  ? t('never')
                  : formatDateTime(segment.updatedAt)}
              </div>
            }
            date={isNever ? null : segment.updatedAt}
          />
        );
      }
    },
    {
      accessorKey: 'action',
      size: 60,
      header: '',
      meta: {
        align: 'center',
        style: { textAlign: 'center', fitContent: true }
      },
      enableSorting: false,
      cell: ({ row }) => {
        const segment = row.original;

        return (
          editable && (
            <Popover
              options={[
                {
                  label: `${t('table:popover.download-segment')}`,
                  icon: IconCloudDownloadOutlined,
                  value: 'DOWNLOAD',
                  disabled: !Number(segment.includedUserCount)
                },
                {
                  label: `${t('table:popover.edit-segment')}`,
                  icon: IconEditOutlined,
                  value: 'EDIT'
                },
                {
                  label: `${t('table:popover.delete-segment')}`,
                  icon: IconDeleteOutlined,
                  value: 'DELETE'
                }
              ]}
              icon={IconMoreHorizOutlined}
              onClick={value =>
                onActionHandler(segment, value as UserSegmentsActionsType)
              }
              align="end"
            />
          )
        );
      }
    }
  ];
};
