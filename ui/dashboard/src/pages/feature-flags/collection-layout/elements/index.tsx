import {
  FunctionComponent,
  PropsWithChildren,
  ReactNode,
  useMemo
} from 'react';
import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import {
  AutoOpsRule,
  FeatureVariation,
  FeatureVariationType,
  Rollout
} from '@types';
import { truncateBySide } from 'utils/converts';
import { copyToClipBoard } from 'utils/function';
import { cn, getVariationColor } from 'utils/style';
import {
  IconCalendar,
  IconCopy,
  IconFlagOperation,
  IconInfo,
  IconInfoFilled,
  IconOperationArrow,
  IconUserSettings
} from '@icons';
import {
  FeatureActivityStatus,
  FlagOperationType
} from 'pages/feature-flags/types';
import Icon, { IconProps } from 'components/icon';
import { Tooltip, TooltipProps } from 'components/tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';

interface FlagNameElementType {
  id: string;
  icon: FunctionComponent;
  name: string;
  link: string;
  status?: FeatureActivityStatus;
  variationType: FeatureVariationType;
  maintainer?: string;
  className?: string;
  iconElement?: ReactNode;
  maxLines?: number;
  variationCls?: string;
  variant?: 'primary' | 'secondary';
}

export const GridViewRoot = ({ children }: PropsWithChildren) => (
  <div className="flex flex-col w-full gap-y-4">{children}</div>
);

export const GridViewRow = ({ children }: PropsWithChildren) => (
  <div className="grid grid-cols-12 items-center w-full max-w-full p-5 gap-x-10 xxl:gap-x-10 rounded shadow-card bg-white self-stretch">
    {children}
  </div>
);

export const FlagDataTypeIcon = ({
  icon,
  className
}: {
  icon: FunctionComponent;
  className?: string;
}) => (
  <div className={cn('flex-center size-8 bg-primary-50 rounded-md', className)}>
    <Icon icon={icon} size={'xxs'} color="primary-500" />
  </div>
);

export const FlagIconWrapper = ({
  icon,
  className,
  color = 'primary-500'
}: IconProps) => (
  <div
    className={cn(
      'flex-center size-[26px] min-w-[26px] bg-primary-50 rounded-md',
      className
    )}
  >
    <Icon icon={icon} size={'xs'} color={color} className="flex-center" />
  </div>
);

export const FlagStatus = ({ status }: { status: FeatureActivityStatus }) => {
  const { t } = useTranslation(['common']);

  const isActive = status === FeatureActivityStatus.ACTIVE;
  const isNew = status === FeatureActivityStatus.NEW;
  const isInActive = !isActive && !isNew;
  const statusKey = isActive ? 'active' : isNew ? 'new' : 'no-activity';

  return (
    <div
      className={cn(
        'flex items-center w-fit min-w-fit gap-x-1 px-2 py-1.5 rounded-[3px] relative',
        {
          'bg-accent-green-50 text-accent-green-500': isActive,
          'bg-accent-yellow-50 text-accent-yellow-500': isInActive,
          'bg-accent-blue-50 text-accent-blue-500': isNew
        }
      )}
    >
      {isInActive && (
        <Icon icon={IconInfoFilled} color="accent-yellow-500" size={'xxs'} />
      )}
      <p className="typo-para-small leading-[14px] capitalize whitespace-nowrap">
        {t(statusKey)}
      </p>
    </div>
  );
};

export const FlagVariationPolygon = ({
  index,
  className
}: {
  index: number;
  className?: string;
}) => {
  const color = getVariationColor(index);
  return (
    <div
      style={{
        background: color,
        zIndex: index
      }}
      className={cn(
        'flex-center size-[14px] min-w-[14px] border border-gray-200 rounded rotate-45',
        className
      )}
    />
  );
};

export const VariationTypeTooltip = ({
  trigger,
  variationType,
  asChild = false,
  className,
  align = 'start'
}: {
  trigger: ReactNode;
  variationType: FeatureVariationType;
  asChild?: boolean;
  className?: string;
  align?: TooltipProps['align'];
}) => (
  <Tooltip
    asChild={asChild}
    align={align}
    trigger={trigger}
    content={
      <Trans
        i18nKey={'table:feature-flags.variation-type'}
        values={{
          type:
            variationType === 'JSON'
              ? variationType
              : variationType?.toLowerCase()
        }}
        components={{
          text: <span className="capitalize" />
        }}
        className={className}
      />
    }
  />
);

const FlagIdElement = ({ id }: { id: string }) => {
  const { notify } = useToast();

  const handleCopyId = (id: string) => {
    copyToClipBoard(id);
    notify({
      toastType: 'toast',
      message: (
        <span>
          <b>ID</b> {` has been successfully copied!`}
        </span>
      )
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

export const FlagNameElement = ({
  id,
  icon,
  name,
  maintainer,
  link,
  status,
  variationType,
  className,
  iconElement,
  maxLines = 2,
  variationCls,
  variant = 'primary'
}: FlagNameElementType) => {
  const { t } = useTranslation(['table']);

  return (
    <div
      className={cn(
        'flex items-center col-span-5 w-full max-w-full gap-x-4 overflow-hidden',
        className
      )}
    >
      <div className="flex flex-col w-full max-w-full gap-y-2">
        <div className="flex items-start w-full gap-x-2">
          <div className="flex-center self-stretch">
            {iconElement || (
              <VariationTypeTooltip
                trigger={
                  <FlagDataTypeIcon
                    icon={icon}
                    className={cn('size-[26px]', variationCls)}
                  />
                }
                variationType={variationType}
              />
            )}
          </div>
          <div className="flex flex-col gap-y-1">
            <NameWithTooltip
              id={id}
              content={<NameWithTooltip.Content content={name} id={id} />}
              trigger={
                <Link to={link}>
                  <NameWithTooltip.Trigger
                    id={id}
                    name={name}
                    maxLines={maxLines}
                  />
                </Link>
              }
              maxLines={maxLines}
            />
            {variant === 'secondary' && <FlagIdElement id={id} />}
          </div>

          {maintainer && (
            <Tooltip
              asChild={false}
              align="start"
              trigger={<FlagIconWrapper icon={IconUserSettings} />}
              content={maintainer}
            />
          )}
          {status && (
            <Tooltip
              asChild={false}
              align="start"
              trigger={<FlagStatus status={status} />}
              content={t(
                `feature-flags.${status === 'active' ? 'active-description' : status === 'in-active' ? 'inactive-description' : 'new-description'}`
              )}
            />
          )}
        </div>
        {variant === 'primary' && <FlagIdElement id={id} />}
      </div>
    </div>
  );
};

export const FlagVariationsElement = ({
  variations
}: {
  variations: FeatureVariation[];
}) => {
  const { t } = useTranslation(['common', 'table']);

  const variationsWrapperWidth =
    document.getElementById('variations-wrapper')?.offsetWidth;

  const variationCount = variations?.length;

  if (!variationCount)
    return (
      <p className="typo-para-small text-gray-700">{t('no-variations')}</p>
    );
  if (variationCount === 1)
    return (
      <div className="flex items-center gap-x-2 w-full overflow-hidden">
        <div className="flex-center size-4">
          <FlagVariationPolygon index={0} />
        </div>
        <p className="typo-para-small text-gray-700 truncate flex-1">
          {variations[variationCount].name}
        </p>
      </div>
    );
  return (
    <div className="flex w-fit max-w-full">
      <Tooltip
        asChild={false}
        align="start"
        trigger={
          <div className="flex items-center w-full gap-2">
            <div className="flex items-center w-full flex-wrap gap-y-1">
              {variations.map((_, index) => (
                <FlagVariationPolygon key={index} index={index} />
              ))}
            </div>
            <p className="typo-para-small whitespace-nowrap text-gray-700">
              {`${variationCount} ${variationCount > 1 ? t('variations') : t('table:results.variation')}`}
            </p>
            <div className="flex-center size-fit">
              <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
            </div>
          </div>
        }
        content={
          <div
            style={{
              maxWidth: variationsWrapperWidth
            }}
            className="flex flex-col gap-y-2 w-full"
          >
            {variations.map((item, index) => (
              <div
                className={'flex items-center gap-x-1 max-w-full'}
                key={index}
              >
                <div className="flex-center size-4">
                  <FlagVariationPolygon
                    index={index === 0 ? index : index + index + 1}
                    className="border-white/10"
                  />
                </div>
                <p className="typo-para-small text-white break-all truncate">
                  {item.name || item.value}
                </p>
              </div>
            ))}
          </div>
        }
      />
    </div>
  );
};

export const FlagOperationsElement = ({
  autoOpsRules,
  rollouts,
  featureId
}: {
  autoOpsRules: AutoOpsRule[];
  rollouts: Rollout[];
  featureId: string;
}) => {
  const { t } = useTranslation(['table']);

  const operationType: FlagOperationType | null = useMemo(() => {
    if (rollouts?.find(item => item.featureId === featureId))
      return FlagOperationType.ROLLOUT;
    const operation = autoOpsRules?.find(item => item.featureId === featureId);
    if (operation?.opsType === 'SCHEDULE') return FlagOperationType.SCHEDULED;
    if (operation?.opsType === 'EVENT_RATE')
      return FlagOperationType.KILL_SWITCH;
    return null;
  }, [autoOpsRules, rollouts, featureId]);

  if (!operationType) return <></>;

  return (
    <div className="flex items-center gap-x-2">
      {operationType === FlagOperationType.ROLLOUT && (
        <Tooltip
          asChild={false}
          trigger={
            <FlagIconWrapper
              icon={IconFlagOperation}
              color="accent-pink-500"
              className="bg-accent-pink-50"
            />
          }
          content={t('feature-flags.progressive-description')}
        />
      )}
      {operationType === FlagOperationType.SCHEDULED && (
        <Tooltip
          asChild={false}
          trigger={
            <FlagIconWrapper
              icon={IconCalendar}
              color="primary-500"
              className="bg-primary-50"
            />
          }
          content={t('feature-flags.scheduled-description')}
        />
      )}
      {operationType === FlagOperationType.KILL_SWITCH && (
        <Tooltip
          asChild={false}
          trigger={
            <FlagIconWrapper
              icon={IconOperationArrow}
              color="accent-blue-500"
              className="bg-accent-blue-50"
            />
          }
          content={t('feature-flags.kill-description')}
        />
      )}
    </div>
  );
};
