import { FunctionComponent } from 'react';
import { Link } from 'react-router-dom';
import { cn } from 'utils/style';
import {
  IconArrowDown,
  IconCalendar,
  IconFlagOperation,
  IconInfo,
  IconInfoFilled,
  IconOperationArrow,
  IconUserSettings
} from '@icons';
import { FlagsViewType, FlagStatusType } from 'pages/feature-flags/types';
import Icon, { IconProps } from 'components/icon';

interface FlagNameElementType {
  icon: FunctionComponent;
  name: string;
  link: string;
  status: FlagStatusType;
  id: string;
  viewType: FlagsViewType;
}

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
    <Icon icon={icon} size={'xs'} color={color} />
  </div>
);

const FlagStatus = ({ status }: { status: FlagStatusType }) => (
  <div
    className={cn(
      'flex items-center w-fit min-w-fit gap-x-1 px-2 py-1.5 rounded-[3px] relative',
      {
        'bg-accent-green-50 text-accent-green-500': status === 'active',
        'bg-accent-yellow-50 text-accent-yellow-500': status === 'no_activity',
        'bg-accent-blue-50 text-accent-blue-500': status === 'new'
      }
    )}
  >
    {status === 'no_activity' && (
      <Icon icon={IconInfoFilled} color="accent-yellow-500" size={'xxs'} />
    )}
    <p className="typo-para-small leading-[14px] capitalize whitespace-nowrap">
      {status.replace('_', ' ')}
    </p>
  </div>
);

export const FlagTag = ({ tag }: { tag: string }) => {
  return (
    <div
      className={
        'flex-center w-fit px-2 py-[5px] typo-para-small leading-[14px] text-center rounded capitalize bg-primary-50 text-primary-500'
      }
    >
      {tag}
    </div>
  );
};

export const FlagVariationPolygon = ({
  color = 'blue',
  className
}: {
  color?: 'blue' | 'pink' | 'green';
  className?: string;
}) => (
  <div
    className={cn(
      'flex-center size-[14px] border border-white rounded-sm rotate-45',
      {
        'bg-accent-blue-500': color === 'blue',
        'bg-accent-pink-500': color === 'pink',
        'bg-accent-green-500': color === 'green'
      },
      className
    )}
  />
);

export const FlagNameElement = ({
  icon,
  name,
  link,
  status,
  id,
  viewType
}: FlagNameElementType) => (
  <div className="flex items-center w-full gap-x-4 min-w-[350px] max-w-[440px]">
    {viewType === 'LIST_VIEW' && <FlagDataTypeIcon icon={icon} />}
    <div className="flex flex-col flex-1 w-full gap-y-2">
      <div className="flex items-center w-full gap-x-2">
        {viewType === 'GRID_VIEW' && (
          <FlagDataTypeIcon icon={icon} className="size-[26px]" />
        )}
        <Link
          to={link}
          className="typo-para-medium text-primary-500 line-clamp-1 break-all underline"
        >
          {name}
        </Link>
        <FlagIconWrapper icon={IconUserSettings} />
        <FlagStatus status={status} />
      </div>
      <p className="typo-para-tiny leading-[14px] text-gray-500 truncate">{id}</p>
    </div>
  </div>
);

export const FlagVariationsElement = () => (
  <div className="flex items-center gap-x-2">
    <div className="flex items-center">
      <FlagVariationPolygon />
      <FlagVariationPolygon color="pink" className="z-[1] -ml-0.5" />
      <FlagVariationPolygon color="green" className="z-[2] -ml-0.5" />
    </div>
    <p className="typo-para-small whitespace-nowrap text-gray-700">
      3 Variations
    </p>
    <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
  </div>
);

export const FlagTagsElement = ({ tags }: { tags: string[] }) => (
  <div className="flex items-center gap-x-2">
    {tags.slice(0, 3).map((tag, index) => (
      <FlagTag key={index} tag={tag} />
    ))}

    {tags.length > 3 && <FlagTag tag={`+${tags.length - 3}`} />}
    {tags.length > 3 && (
      <Icon
        icon={IconArrowDown}
        size={'sm'}
        color="gray-500"
        className="cursor-pointer"
      />
    )}
  </div>
);

export const FlagOperationsElement = () => (
  <div className="flex items-center gap-x-2">
    <FlagIconWrapper
      icon={IconFlagOperation}
      color="accent-pink-500"
      className="bg-accent-pink-50"
    />
    <FlagIconWrapper
      icon={IconCalendar}
      color="primary-500"
      className="bg-primary-100"
    />
    <FlagIconWrapper
      icon={IconOperationArrow}
      color="accent-blue-500"
      className="bg-accent-blue-50"
    />
  </div>
);
