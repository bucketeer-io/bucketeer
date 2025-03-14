import { ReactNode } from 'react';
import {
  IconArchiveOutlined,
  IconMoreVertOutlined,
  IconSaveAsFilled
} from 'react-icons-material-design';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { Feature } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import { IconWatch } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import ExpandableTag from 'elements/expandable-tag';
import { FlagActionType } from '../types';
import {
  FlagNameElement,
  FlagOperationsElement,
  FlagVariationsElement,
  GridViewRoot,
  GridViewRow
} from './elements';
import { getDataTypeIcon, getFlagStatus } from './elements/utils';

const GridViewCollection = ({
  data,
  emptyState,
  onActions
}: {
  data: Feature[];
  emptyState: ReactNode;
  onActions: (item: Feature, type: FlagActionType) => void;
}) => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { searchOptions } = useSearchParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  if (!data?.length) return <div className="pt-32">{emptyState}</div>;

  return (
    <GridViewRoot>
      {data.map((item, index) => {
        const {
          id,
          name,
          maintainer,
          tags,
          updatedAt,
          enabled,
          variationType,
          variations,
          autoOpsSummary
        } = item;
        return (
          <GridViewRow key={index}>
            <FlagNameElement
              id={id}
              link={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${id}/targeting`}
              name={name}
              maintainer={maintainer}
              icon={getDataTypeIcon(variationType)}
              status={getFlagStatus(item)}
            />
            <div className="flex flex-col w-fit gap-y-3 min-w-[300px]">
              <FlagVariationsElement variations={variations} />
              <div className="flex items-center flex-wrap w-[410px] max-w-[410px] gap-2">
                <ExpandableTag
                  tags={tags}
                  rowId={item.id}
                  className={cn('!max-w-[180px] truncate', {
                    '!max-w-[250px]': tags.length <= 1
                  })}
                  wrapperClassName="w-fit"
                  maxSize={tags.length <= 1 ? 250 : 180}
                />
                <FlagOperationsElement autoOpsSummary={autoOpsSummary} />
              </div>
            </div>
            <div className="flex flex-1 justify-end self-start h-full gap-x-2">
              <div className="flex-center">
                <Icon icon={IconWatch} size={'xxs'} />
              </div>
              <div className="text-gray-700 typo-para-small whitespace-nowrap">
                {Number(updatedAt) === 0
                  ? t('never')
                  : `Updated ${formatDateTime(updatedAt)}`}
              </div>
              <div className="flex-center">
                <Switch
                  checked={enabled}
                  onCheckedChange={() =>
                    onActions(item, enabled ? 'INACTIVE' : 'ACTIVE')
                  }
                />
              </div>
              <Popover
                options={compact([
                  searchOptions.status === 'ARCHIVED'
                    ? {
                        label: `${t('table:popover.unarchive-flag')}`,
                        icon: IconArchiveOutlined,
                        value: 'UNARCHIVE'
                      }
                    : {
                        label: `${t('table:popover.archive-flag')}`,
                        icon: IconArchiveOutlined,
                        value: 'ARCHIVE'
                      },
                  {
                    label: `${t('table:popover.clone-flag')}`,
                    icon: IconSaveAsFilled,
                    value: 'CLONE'
                  }
                ])}
                icon={IconMoreVertOutlined}
                onClick={value => onActions(item, value as FlagActionType)}
                align="end"
              />
            </div>
          </GridViewRow>
        );
      })}
    </GridViewRoot>
  );
};

export default GridViewCollection;
