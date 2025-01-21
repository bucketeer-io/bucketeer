import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import Icon from 'components/icon';

const Polygon = ({ className }: { className: string }) => (
  <div
    className={cn(
      'flex-center size-[14px] border rounded-sm rotate-45',
      className
    )}
  />
);

const HeaderDetails = () => {
  return (
    <div className="flex flex-col w-full gap-y-4 mt-4">
      <div className="flex items-center gap-x-1.5">
        <div className="flex items-center">
          <Polygon className="bg-accent-blue-500 border-transparent size-3" />
          <Polygon className="bg-accent-pink-500 border-white -ml-0.5 relative z-10" />
        </div>
        <p className="typo-para-small text-gray-700">2 Variations</p>
        <Icon
          icon={IconInfo}
          color="gray-600"
          size={'xxs'}
          className="flex-center"
        />
      </div>
      <div className="flex items-center gap-x-2">
        <h1 className="text-gray-900 typo-head-bold-huge">
          This is a big experiment name
        </h1>
        <div className="flex-center w-fit px-2 py-1.5 bg-accent-green-50 text-accent-green-500 typo-para-small leading-[14px] rounded whitespace-nowrap">
          Running
        </div>
      </div>
    </div>
  );
};

export default HeaderDetails;
