import { cn } from 'utils/style';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import TruncateWithTooltip from 'elements/truncate-with-tooltip';

const VariationLabel = ({
  label,
  index,
  className
}: {
  label: string;
  index: number;
  className?: string;
}) => {
  return (
    <div className={cn('flex items-center gap-x-2 pl-0.5', className)}>
      <FlagVariationPolygon index={index} />
      <TruncateWithTooltip
        text={label}
        maxLines={1}
        align="start"
        className="-mt-0.5"
      />
    </div>
  );
};

export default VariationLabel;
