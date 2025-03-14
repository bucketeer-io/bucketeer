import { FunctionComponent, PropsWithChildren } from 'react';
import { Link } from 'react-router-dom';
import { COLORS } from 'constants/styles';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { AutoOpsSummary, FeatureVariation } from '@types';
import { truncateTextCenter } from 'utils/converts';
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
import Icon, { IconProps } from 'components/icon';

interface FlagNameElementType {
  id: string;
  icon: FunctionComponent;
  name: string;
  link: string;
  status: FeatureActivityStatus;
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
        'flex-center size-[14px] border border-white rounded rotate-45',
        className
      )}
    />
  );
};

export const FlagNameElement = ({
  id,
  icon,
  name,
  maintainer,
  link,
  status
}: FlagNameElementType) => {
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
    <div className="flex items-center w-full min-w-[400px] max-w-[400px] xxl:min-w-[500px] gap-x-4">
      <div className="flex flex-col flex-1 w-full gap-y-2">
        <div className="flex items-center w-full gap-x-2">
          <FlagDataTypeIcon icon={icon} className="size-[26px]" />
          <Link
            to={link}
            className="typo-para-medium text-primary-500 line-clamp-1 break-all underline"
          >
            {name}
          </Link>
          {maintainer && <FlagIconWrapper icon={IconUserSettings} />}
          <FlagStatus status={status} />
        </div>
        <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
          {truncateTextCenter(id)}
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
  if (!variationCount)
    return (
      <p className="typo-para-small text-gray-700">{t('no-variations')}</p>
    );
  if (variationCount === 1)
    return (
      <div className="flex items-center gap-x-2">
        <FlagVariationPolygon index={0} />
        <p className="typo-para-small text-gray-700">
          {variations[variationCount]?.name}
        </p>
      </div>
    );
  return (
    <div className="flex items-center gap-x-2">
      <div className="flex items-center">
        {variations.map((_, index) => (
          <FlagVariationPolygon key={index} index={index} />
        ))}
      </div>
      <p className="typo-para-small whitespace-nowrap text-gray-700">
        {`${variationCount} ${variationCount > 1 ? t('variations') : t('table:results.variation')}`}
      </p>
      <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
    </div>
  );
};

export const FlagOperationsElement = ({
  autoOpsSummary
}: {
  autoOpsSummary: AutoOpsSummary;
}) => (
  <div className="flex items-center gap-x-2">
    {!!autoOpsSummary?.progressiveRolloutCount && (
      <FlagIconWrapper
        icon={IconFlagOperation}
        color="accent-pink-500"
        className="bg-accent-pink-50"
      />
    )}
    {!!autoOpsSummary?.scheduleCount && (
      <FlagIconWrapper
        icon={IconCalendar}
        color="primary-500"
        className="bg-primary-50"
      />
    )}
    {!!autoOpsSummary?.killSwitchCount && (
      <FlagIconWrapper
        icon={IconOperationArrow}
        color="accent-blue-500"
        className="bg-accent-blue-50"
      />
    )}
  </div>
);
