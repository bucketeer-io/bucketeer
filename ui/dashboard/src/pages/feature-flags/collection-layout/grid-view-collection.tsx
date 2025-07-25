import { ReactNode, useCallback, useMemo } from 'react';
import { Trans } from 'react-i18next';
import {
  IconArchiveOutlined,
  IconSaveAsFilled
} from 'react-icons-material-design';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useScreen } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Account, AutoOpsRule, Feature, Rollout } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import { IconWatch } from '@icons';
import Icon from 'components/icon';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import ExpandableTag from 'elements/expandable-tag';
import TableListContent from 'elements/table-list-content';
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
  filterTags,
  autoOpsRules,
  rollouts,
  accounts,
  data,
  emptyState,
  onActions,
  handleTagFilters
}: {
  filterTags?: string[];
  autoOpsRules: AutoOpsRule[];
  rollouts: Rollout[];
  accounts: Account[];
  data: Feature[];
  emptyState: ReactNode;
  onActions: (item: Feature, type: FlagActionType) => void;
  handleTagFilters: (tag: string) => void;
}) => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { searchOptions } = useSearchParams();
  const { fromXLScreen } = useScreen();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  const popoverOptions = useMemo(
    () =>
      compact([
        searchOptions.tab === 'ARCHIVED'
          ? {
              label: `${t('unarchive-flag')}`,
              icon: IconArchiveOutlined,
              value: 'UNARCHIVE'
            }
          : {
              label: `${t('archive-flag')}`,
              icon: IconArchiveOutlined,
              value: 'ARCHIVE'
            },
        {
          label: `${t('clone-flag')}`,
          icon: IconSaveAsFilled,
          value: 'CLONE'
        }
      ]),
    [searchOptions]
  );

  const handleGetMaintainerInfo = useCallback(
    (email: string) => {
      const existedAccount = accounts?.find(account => account.email === email);
      if (
        !existedAccount ||
        !existedAccount?.firstName ||
        !existedAccount?.lastName
      )
        return email;
      return `${existedAccount.firstName} ${existedAccount.lastName}`;
    },
    [accounts]
  );

  if (!data?.length) return <div className="pt-32">{emptyState}</div>;

  return (
    <TableListContent>
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
            variations
          } = item;
          return (
            <GridViewRow key={index}>
              <FlagNameElement
                id={id}
                link={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${id}/targeting`}
                name={name}
                maintainer={handleGetMaintainerInfo(maintainer)}
                variationType={variationType}
                icon={getDataTypeIcon(variationType)}
                status={getFlagStatus(item)}
              />
              <div
                id="variations-wrapper"
                className="flex flex-col gap-y-3 col-span-4 flex-1"
              >
                <FlagVariationsElement variations={variations} />
                <div className="flex items-center flex-wrap w-full gap-2">
                  <ExpandableTag
                    tags={tags}
                    filterTags={filterTags}
                    rowId={item.id}
                    className={cn('!max-w-[200px] truncate cursor-pointer', {
                      '!max-w-[350px]': fromXLScreen
                    })}
                    wrapperClassName="w-fit"
                    maxSize={fromXLScreen ? 382 : 232}
                    onTagClick={tag => handleTagFilters(tag)}
                  />
                  <FlagOperationsElement
                    autoOpsRules={autoOpsRules}
                    rollouts={rollouts}
                    featureId={id}
                  />
                </div>
              </div>
              <div className="flex col-span-3 justify-end self-start gap-x-2">
                <div className="flex-center">
                  <Icon icon={IconWatch} size={'xxs'} />
                </div>
                <DateTooltip
                  trigger={
                    <div className="text-gray-700 typo-para-small whitespace-nowrap">
                      {Number(updatedAt) === 0 ? (
                        t('never')
                      ) : (
                        <Trans
                          i18nKey={'common:time-updated'}
                          values={{
                            time: formatDateTime(updatedAt)
                          }}
                        />
                      )}
                    </div>
                  }
                  date={Number(updatedAt) === 0 ? null : updatedAt}
                />
                <DisabledButtonTooltip
                  hidden={editable}
                  trigger={
                    <div className="flex-center">
                      <Switch
                        disabled={!editable}
                        checked={enabled}
                        onCheckedChange={() =>
                          onActions(item, enabled ? 'INACTIVE' : 'ACTIVE')
                        }
                      />
                    </div>
                  }
                />
                <DisabledPopoverTooltip
                  onClick={value => onActions(item, value as FlagActionType)}
                  options={popoverOptions}
                />
              </div>
            </GridViewRow>
          );
        })}
      </GridViewRoot>
    </TableListContent>
  );
};

export default GridViewCollection;
