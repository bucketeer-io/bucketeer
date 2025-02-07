import { PropsWithChildren } from 'react';
import {
  IconArchiveOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { useFormatDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import { IconTrash, IconWatch } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import { FlagsTemp } from '../types';
import { getDataTypeIcon } from './data-collection';
import {
  FlagNameElement,
  FlagOperationsElement,
  FlagTagsElement,
  FlagVariationsElement
} from './elements';

const GridViewRoot = ({ children }: PropsWithChildren) => (
  <div className="flex flex-col w-full gap-y-4">{children}</div>
);

const GridViewRow = ({ children }: PropsWithChildren) => (
  <div className="flex items-center w-full p-5 gap-x-2 rounded shadow-card bg-white">
    {children}
  </div>
);

const GridViewCollection = ({
  data,
  onActions
}: {
  data: FlagsTemp[];
  onActions: (item: FlagsTemp, type: unknown) => void;
}) => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { searchOptions } = useSearchParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <GridViewRoot>
      {data.map((item, index) => {
        const { id, name, type, status, tags, updatedAt, disabled } = item;
        return (
          <GridViewRow key={index}>
            <FlagNameElement
              link={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${id}`}
              name={name}
              id={id}
              icon={getDataTypeIcon(type)}
              status={status}
              viewType="GRID_VIEW"
            />
            <div className="flex flex-col w-full gap-y-3 max-w-[442px] min-w-[300px]">
              <FlagVariationsElement />
              <div className="flex items-center w-full gap-x-2">
                <FlagTagsElement tags={tags} />
                <FlagOperationsElement />
              </div>
            </div>
            <div className="flex flex-1 justify-end self-start h-full gap-x-2">
              <div className="flex-center">
                <Icon icon={IconWatch} size={'xxs'} />
              </div>
              <div className="text-gray-700 typo-para-medium">
                {Number(updatedAt) === 0
                  ? t('never')
                  : `Updated ${formatDateTime(updatedAt)}`}
              </div>
              <div className="flex-center">
                <Switch checked={!disabled} />
              </div>
              <Popover
                options={compact([
                  searchOptions.status === 'ARCHIVED'
                    ? {
                        label: `${t('table:popover.unarchive-goal')}`,
                        icon: IconArchiveOutlined,
                        value: 'UNARCHIVE'
                      }
                    : {
                        label: `${t('table:popover.archive-goal')}`,
                        icon: IconArchiveOutlined,
                        value: 'ARCHIVE'
                      },
                  {
                    label: `${t('table:popover.delete-goal')}`,
                    icon: IconTrash,
                    value: 'DELETE'
                  }
                ])}
                icon={IconMoreVertOutlined}
                onClick={value => onActions(item, value as unknown)}
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
