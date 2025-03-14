import { ReactNode, useCallback, useState } from 'react';
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
import { IconWatch } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import { FlagActionType } from '../types';
import {
  FlagNameElement,
  FlagOperationsElement,
  FlagTagsElement,
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

  const [tagsExpanded, setTagsExpanded] = useState<string[]>([]);

  const onToggleExpandTag = useCallback(
    (id: string) => {
      const isExpanded = tagsExpanded.includes(id);
      setTagsExpanded(
        isExpanded
          ? tagsExpanded.filter(item => item !== id)
          : [...tagsExpanded, id]
      );
    },
    [tagsExpanded]
  );

  if (!data?.length) return emptyState;

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
              <div className="flex items-center w-full gap-x-2">
                <FlagTagsElement
                  tags={tags}
                  isExpanded={tagsExpanded.includes(id)}
                  onToggleExpandTag={() => onToggleExpandTag(id)}
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
