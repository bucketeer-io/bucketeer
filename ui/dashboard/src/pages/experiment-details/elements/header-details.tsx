import { CSSProperties } from 'react';
import { Experiment } from '@types';
import { cn, getVariationColor } from 'utils/style';
import { IconInfo } from '@icons';
import Icon from 'components/icon';

export const Polygon = ({
  className,
  style
}: {
  className?: string;
  style?: CSSProperties;
}) => (
  <div
    className={cn(
      'flex-center size-[14px] border border-white rounded-sm rotate-45',
      className
    )}
    style={{
      ...style
    }}
  />
);

const HeaderDetails = ({ experiment }: { experiment: Experiment }) => {
  return (
    <div className="flex flex-col w-full gap-y-4 mt-4">
      <div className="flex items-center gap-x-1.5">
        <div className="flex items-center">
          {experiment.variations?.map((_, index) => (
            <Polygon
              key={index}
              className="size-3"
              style={{
                background: getVariationColor(index),
                zIndex: index
              }}
            />
          ))}
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
