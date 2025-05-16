import { Feature } from '@types';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import { FlagStatus } from 'pages/feature-flags/collection-layout/elements';
import { getFlagStatus } from 'pages/feature-flags/collection-layout/elements/utils';
import Icon from 'components/icon';

export const Polygon = ({ className }: { className: string }) => (
  <div
    className={cn(
      'flex-center size-[14px] border rounded-sm rotate-45',
      className
    )}
  />
);

const HeaderDetails = ({ feature }: { feature: Feature }) => {
  return (
    <div className="flex flex-col w-full gap-y-4 mt-4">
      <div className="flex items-center gap-x-2">
        <h1 className="text-gray-900 typo-head-bold-huge">{feature.name}</h1>
        <FlagStatus status={getFlagStatus(feature)} />
        <Icon
          icon={IconChevronRight}
          className="rotate-90"
          color="gray-500"
          size={'sm'}
        />
      </div>
    </div>
  );
};

export default HeaderDetails;
