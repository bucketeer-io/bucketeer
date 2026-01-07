import { useCallback, useMemo } from 'react';
import { Trans } from 'react-i18next';
import {
  IconArchiveOutlined,
  IconSaveAsFilled
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useScreen, useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { ArrowRight } from 'lucide-react';
import { Account, AutoOpsRule, Feature, Rollout } from '@types';
import { truncateBySide } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import { IconCopy, IconUserSettings, IconWatch } from '@icons';
import {
  FlagIconWrapper,
  FlagOperationsElement,
  FlagStatus,
  FlagVariationsElement
} from 'pages/feature-flags/collection-layout/elements';
import {
  getDataTypeIcon,
  getFlagStatus
} from 'pages/feature-flags/collection-layout/elements/utils';
import { FlagActionType } from 'pages/feature-flags/types';
import Icon from 'components/icon';
import Switch from 'components/switch';
import { Tooltip } from 'components/tooltip';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import ExpandableTag from 'elements/expandable-tag';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

interface FeatureCardProps {
  data: Feature;
  accounts: Account[];
  filterTags?: string[];
  rollouts: Rollout[];
  autoOpsRules: AutoOpsRule[];
  handleTagFilters: (tag: string) => void;
  onActions: (item: Feature, type: FlagActionType) => void;
}

const FlagIdElement = ({ id }: { id: string }) => {
  const { t } = useTranslation(['message']);
  const { notify } = useToast();

  const handleCopyId = (id: string) => {
    copyToClipBoard(id);
    notify({
      message: t('copied')
    });
  };
  return (
    <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
      <p className="truncate">{truncateBySide(id, 55)}</p>
      <div onClick={() => handleCopyId(id)}>
        <Icon
          icon={IconCopy}
          size={'sm'}
          className="opacity-0 group-hover:opacity-100 cursor-pointer"
        />
      </div>
    </div>
  );
};

export const FeatureCard: React.FC<FeatureCardProps> = ({
  data,
  rollouts,
  autoOpsRules,
  accounts,
  filterTags,
  handleTagFilters,
  onActions
}) => {
  const { t } = useTranslation(['common', 'table']);
  const { fromXLScreen } = useScreen();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { searchOptions } = useSearchParams();
  const formatDateTime = useFormatDateTime();
  const editable = hasEditable(consoleAccount!);
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
  const maintainer = useMemo(
    () => handleGetMaintainerInfo(data.maintainer),
    [data]
  );
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
  return (
    <Card>
      <Card.Header
        icon={<Icon icon={getDataTypeIcon(data.variationType)} />}
        title={data.name}
        subtitle={data.id}
        triger={
          <div className="flex flex-col gap-y-1">
            <NameWithTooltip
              id={data.id}
              content={
                <NameWithTooltip.Content content={data.name} id={data.id} />
              }
              trigger={
                <Link
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${data.id}/targeting`}
                >
                  <NameWithTooltip.Trigger
                    id={data.id}
                    name={data.name}
                    maxLines={1}
                  />
                </Link>
              }
              maxLines={1}
            />
            {<FlagIdElement id={data.id} />}
          </div>
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            onClick={value => onActions(data, value as FlagActionType)}
            options={popoverOptions}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {maintainer && (
              <Tooltip
                asChild={false}
                align="start"
                trigger={<FlagIconWrapper icon={IconUserSettings} />}
                content={maintainer}
              />
            )}
            {getFlagStatus(data) && (
              <Tooltip
                asChild={false}
                align="start"
                trigger={<FlagStatus status={getFlagStatus(data)} />}
                content={t(`feature-flags.${getFlagStatus(data)}-description`)}
                className="max-w-[300px]"
              />
            )}
          </div>
          <DisabledButtonTooltip
            hidden={editable}
            trigger={
              <div className="flex-center">
                <Switch
                  disabled={!editable}
                  checked={data.enabled}
                  onCheckedChange={() =>
                    onActions(data, data.enabled ? 'INACTIVE' : 'ACTIVE')
                  }
                />
              </div>
            }
          />
        </div>
        <div className="flex flex-col gap-2 rounded-lg bg-gray-100 p-3 mt-6">
          <div className="flex items-start justify-between bg-gray-100 rounded-lg">
            <div
              id="variations-wrapper"
              className="flex flex-col gap-y-3 col-span-4 flex-1"
            >
              <FlagVariationsElement variations={data.variations} />
            </div>
          </div>
          <div className="flex items-center flex-wrap w-full gap-2 rounded-lg p-2">
            <ExpandableTag
              tags={data.tags}
              filterTags={filterTags}
              rowId={data.id}
              className={cn('!max-w-[220px] truncate cursor-pointer', {
                '!max-w-[350px]': fromXLScreen
              })}
              wrapperClassName="w-fit"
              maxSize={fromXLScreen ? 382 : 232}
              onTagClick={tag => handleTagFilters(tag)}
            />
            <FlagOperationsElement
              autoOpsRules={autoOpsRules}
              rollouts={rollouts}
              featureId={data.id}
            />
          </div>
        </div>
      </Card.Meta>
      <Card.Footer
        left={
          <div className="flex-center gap-2">
            <Icon icon={IconWatch} size={'xxs'} />
            <DateTooltip
              trigger={
                <div className="text-gray-700 typo-para-small whitespace-nowrap">
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
        right={
          <div className="flex items-center typo-para-tiny font-bold text-primary-500 gap-1">
            <Link
              to={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${data.id}/targeting`}
            >
              {t('common:detail')}
            </Link>
            <Icon icon={ArrowRight} size="xxs" />
          </div>
        }
      />
    </Card>
  );
};
