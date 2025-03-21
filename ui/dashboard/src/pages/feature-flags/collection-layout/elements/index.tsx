import {
  FunctionComponent,
  PropsWithChildren,
  ReactNode,
  useMemo
} from 'react';
import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { COLORS } from 'constants/styles';
import { useScreen, useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { AutoOpsSummary, FeatureVariation, FeatureVariationType } from '@types';
import { truncateBySide } from 'utils/converts';
import { copyToClipBoard } from 'utils/function';
import { cn } from 'utils/style';
import {
  IconCalendar,
  IconCopy,
  IconFlagOperation,
  IconInfo,
  IconInfoFilled,
  IconOperationArrow,
  IconUserSettings
} from '@icons';
import { FeatureActivityStatus } from 'pages/feature-flags/types';
import Divider from 'components/divider';
import Icon, { IconProps } from 'components/icon';
import { Tooltip } from 'components/tooltip';
import TruncationWithTooltip from 'elements/truncation-with-tooltip';

interface FlagNameElementType {
  id: string;
  icon: FunctionComponent;
  name: string;
  link: string;
  status: FeatureActivityStatus;
  variationType: FeatureVariationType;
  maintainer?: string;
}

export const GridViewRoot = ({ children }: PropsWithChildren) => (
  <div className="flex flex-col w-ful min-w-max overflow-visible gap-y-4">
    {children}
  </div>
);

export const GridViewRow = ({ children }: PropsWithChildren) => (
  <div className="flex items-center w-full min-w-fit p-5 gap-x-4 xxl:gap-x-10 rounded shadow-card bg-white self-stretch">
    {children}
  </div>
);

const FlagDataTypeIcon = ({
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
  const colorIndex = index > 20 ? index % 20 : index;
  const color = COLORS[colorIndex];
  return (
    <div
      style={{
        background: color,
        zIndex: index
      }}
      className={cn(
        'flex-center size-[14px] min-w-[14px] border border-white rounded rotate-45',
        className
      )}
    />
  );
};

export const VariationTypeTooltip = ({
  trigger,
  variationType,
  asChild = false,
  className
}: {
  trigger: ReactNode;
  variationType: FeatureVariationType;
  asChild?: boolean;
  className?: string;
}) => (
  <Tooltip
    asChild={asChild}
    align="start"
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

export const FlagNameElement = ({
  id,
  icon,
  name,
  maintainer,
  link,
  status,
  variationType
}: FlagNameElementType) => {
  const { notify } = useToast();
  const { t } = useTranslation(['table']);
  const { isXXLScreen } = useScreen();

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
    <div className="flex items-center w-full min-w-[400px] max-w-[400px] xxl:min-w-[500px] gap-x-4">
      <div className="flex flex-col flex-1 w-full gap-y-2">
        <div className="flex items-center w-full gap-x-2">
          <div className="flex-center size-fit">
            <VariationTypeTooltip
              trigger={<FlagDataTypeIcon icon={icon} className="size-[26px]" />}
              variationType={variationType}
            />
          </div>
          <div>
            <TruncationWithTooltip
              content={name}
              elementId={id}
              maxSize={isXXLScreen ? 370 : 270}
              className="w-fit max-w-[270px] xxl:max-w-[370px]"
              tooltipWrapperCls="left-0 translate-x-0"
            >
              <Link
                id={id}
                to={link}
                className="typo-para-medium text-primary-500 underline w-full"
              >
                <p className="truncate">{name}</p>
              </Link>
            </TruncationWithTooltip>
          </div>
          {maintainer && (
            <Tooltip
              asChild={false}
              align="start"
              trigger={<FlagIconWrapper icon={IconUserSettings} />}
              content={maintainer}
            />
          )}
          <Tooltip
            asChild={false}
            align="start"
            trigger={<FlagStatus status={status} />}
            content={t(
              `feature-flags.${status === 'active' ? 'active-description' : status === 'in-active' ? 'inactive-description' : 'new-description'}`
            )}
          />
        </div>
        <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
          {truncateBySide(id, 20)}
          <div onClick={() => handleCopyId(id)}>
            <Icon
              icon={IconCopy}
              size={'sm'}
              className="opacity-0 group-hover:opacity-100 cursor-pointer"
            />
          </div>
        </div>
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

  const variationCount = variations?.length;

  const _variations = useMemo(() => {
    if (variationCount <= 2) return [variations?.map(item => item.name)];
    const vars: string[][] = [];

    variations.forEach(item => {
      if (!vars.length || vars[vars.length - 1]?.length > 1) {
        vars.push([item.name]);
      } else vars[vars.length - 1] = [...vars[vars.length - 1], item.name];
    });
    return vars;
  }, [variations, variationCount]);

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
          {_variations[variationCount]}
        </p>
      </div>
    );
  return (
    <div className="flex items-center gap-x-2 w-full">
      <div className="flex items-center">
        {variations.map((_, index) => (
          <FlagVariationPolygon key={index} index={index} />
        ))}
      </div>
      <p className="typo-para-small whitespace-nowrap text-gray-700">
        {`${variationCount} ${variationCount > 1 ? t('variations') : t('table:results.variation')}`}
      </p>
      <Tooltip
        asChild={false}
        trigger={
          <div className="flex-center size-fit">
            <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
          </div>
        }
        content={
          <div className="flex flex-col gap-y-3 w-full max-w-[420px]">
            {_variations.map((item, index) => (
              <div className="flex items-center w-full gap-3" key={index}>
                {item.map((variation, variationIndex) => (
                  <div
                    className={cn('flex items-center gap-x-1 max-w-[140px]', {
                      'w-[140px]': _variations.length > 1
                    })}
                    key={variationIndex}
                  >
                    {variationIndex !== 0 && (
                      <Divider
                        className={cn('h-2 min-w-px bg-white/15 border-none', {
                          'mr-2.5': variationIndex !== 0
                        })}
                      />
                    )}
                    <div className="flex-center size-4">
                      <FlagVariationPolygon
                        index={
                          index === 0
                            ? variationIndex
                            : variationIndex + index + 1
                        }
                        className="border-none"
                      />
                    </div>
                    <p className="typo-para-small text-white break-all truncate">
                      {variation}
                    </p>
                  </div>
                ))}
              </div>
            ))}
          </div>
        }
      />
    </div>
  );
};

export const FlagOperationsElement = ({
  autoOpsSummary
}: {
  autoOpsSummary: AutoOpsSummary;
}) => {
  const { t } = useTranslation(['table']);

  return (
    <div className="flex items-center gap-x-2">
      {!!autoOpsSummary?.progressiveRolloutCount && (
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
      {!!autoOpsSummary?.scheduleCount && (
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
      {!!autoOpsSummary?.killSwitchCount && (
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
