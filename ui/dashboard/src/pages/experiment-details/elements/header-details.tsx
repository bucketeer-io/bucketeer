import { Experiment } from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import Icon from 'components/icon';

export const Polygon = ({ className }: { className: string }) => (
  <div
    className={cn(
      'flex-center size-[14px] border rounded-sm rotate-45',
      className
    )}
  />
);

const HeaderDetails = ({ experiment }: { experiment: Experiment }) => {
  return (
    <div className="flex flex-col w-full gap-y-4 mt-4">
      <div className="flex items-center gap-x-1.5">
        <div className="flex items-center">
          <Polygon className="bg-accent-blue-500 border-transparent size-3" />
          <Polygon className="bg-accent-pink-500 border-white -ml-0.5 relative z-10" />
        </div>
        <p className="typo-para-small text-gray-700">
          {experiment.variations?.length} Variations
        </p>
        <Icon
          icon={IconInfo}
          color="gray-600"
          size={'xxs'}
          className="flex-center"
        />
      </div>
      <div className="flex items-center">
        <h1 className="text-gray-900 typo-head-bold-huge">{experiment.name}</h1>
      </div>
    </div>
  );
};

export default HeaderDetails;
