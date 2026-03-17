import { Trans } from 'react-i18next';
import {
  IconCloudDownloadOutlined,
  IconDeleteOutlined,
  IconEditOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { UserSegment } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { IconMember, IconSwitch, IconUser, IconWatch } from '@icons';
import { UserSegmentsActionsType } from 'pages/user-segments/types';
import Divider from 'components/divider';
import Icon from 'components/icon';
import Spinner from 'components/spinner';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

interface UserSegmentCardProps {
  data: UserSegment;
  getUploadingStatus: (segmet: UserSegment) => boolean | undefined;
  onActionHandler: (value: UserSegment, type: UserSegmentsActionsType) => void;
}

export const UserSegmentCard: React.FC<UserSegmentCardProps> = ({
  data,
  getUploadingStatus,
  onActionHandler
}) => {
  const { t } = useTranslation(['common', 'table']);
  const isUploading = getUploadingStatus(data);

  const formatDateTime = useFormatDateTime();
  return (
    <Card>
      <Card.Header
        icon={<Icon icon={IconUser} />}
        triger={
          <div className="flex flex-col gap-y-1">
            <div
              onClick={() =>
                onActionHandler(data, isUploading ? 'UPLOADING' : 'EDIT')
              }
              className="flex items-center gap-x-2 cursor-pointer min-w-[150px]"
            >
              <NameWithTooltip
                id={data.id}
                content={
                  <NameWithTooltip.Content content={data.name} id={data.id} />
                }
                trigger={
                  <NameWithTooltip.Trigger id={data.id} name={data.name} />
                }
              />
              {isUploading && <Spinner />}
            </div>
            <div
              className={cn(
                'typo-para-small text-accent-green-500 bg-accent-green-50 px-2 py-[3px] w-fit text-center whitespace-nowrap rounded',
                {
                  'bg-gray-200 text-gray-600': !data.isInUseStatus,
                  'bg-accent-orange-50 text-accent-orange-500': isUploading
                }
              )}
            >
              {isUploading
                ? t('uploading')
                : data.isInUseStatus
                  ? t('in-use')
                  : t('not-in-use')}
            </div>
          </div>
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            icon={IconMoreVertOutlined}
            options={compact([
              {
                label: `${t('table:popover.download-segment')}`,
                icon: IconCloudDownloadOutlined,
                value: 'DOWNLOAD',
                disabled: !Number(data.includedUserCount)
              },
              !getUploadingStatus(data) && {
                label: `${t('table:popover.edit-segment')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              {
                label: `${t('table:popover.delete-segment')}`,
                icon: IconDeleteOutlined,
                value: 'DELETE'
              }
            ])}
            onClick={value =>
              onActionHandler(data, value as UserSegmentsActionsType)
            }
            align="end"
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex h-full w-full flex-wrap items-stretch justify-between gap-3 pb-3">
          <div className="flex-1 typo-para-tiny font-bold bg-gray-100 p-3 rounded-xl min-w-[130px]">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('user')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1 text-primary-500">
              <Icon icon={IconMember} size="sm" />
              <p className="text-nowrap font-bold typo-para-medium">
                {data.includedUserCount}
              </p>
            </div>
          </div>
          <div className="flex-1 typo-para-tiny font-bold bg-gray-100 p-3 rounded-xl min-w-[130px]">
            <p className="flex items-center gap-1 uppercase text-gray-500">
              <span>{t('table:goals.connections')}</span>
            </p>
            <div className="mt-3 flex items-center gap-1">
              <Icon icon={IconSwitch} size="sm" className="text-primary-500" />
              <p
                className={cn(
                  'flex items-center text-primary-500 typo-para-medium font-bold',
                  {
                    'cursor-pointer': data?.features?.length
                  }
                )}
                onClick={() =>
                  data?.features?.length && onActionHandler(data, 'FLAG')
                }
              >
                {data?.features?.length}
                {` ${data?.features?.length === 1 ? t('flag') : t('table:flags')}`}
              </p>
            </div>
          </div>
        </div>
        <Divider />
      </Card.Meta>
      <Card.Footer
        left={
          <div className="flex-center gap-2">
            <Icon icon={IconWatch} size={'xxs'} />
            <DateTooltip
              trigger={
                <div className="text-gray-500 typo-para-small whitespace-nowrap">
                  {Number(data.updatedAt) === 0 ? (
                    t('never')
                  ) : (
                    <Trans
                      i18nKey={'common:time-updated'}
                      values={{
                        time: formatDateTime(data.updatedAt)
                      }}
                    />
                  )}
                </div>
              }
              date={Number(data.updatedAt) === 0 ? null : data.updatedAt}
            />
          </div>
        }
      />
    </Card>
  );
};
