import { useId } from 'react';
import { cn } from 'utils/style';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import NameWithTooltip from 'elements/name-with-tooltip';

const VariationLabel = ({
  label,
  index,
  className,
  asChild = true
}: {
  label: string;
  index: number;
  className?: string;
  asChild?: boolean;
}) => {
  // Unique per instance: the same index is rendered in several places on one page
  const id = `variation-label-${useId()}`;

  return (
    <div className={cn('flex items-center gap-x-2 pl-0.5', className)}>
      <FlagVariationPolygon index={index} />

      <NameWithTooltip
        id={id}
        maxLines={1}
        content={
          <div className="flex items-center gap-x-2">
            <FlagVariationPolygon index={index} />
            <p>{label}</p>
          </div>
        }
        asChild={asChild}
        align="start"
        trigger={
          <NameWithTooltip.Trigger
            id={id}
            name={label}
            maxLines={1}
            className="-mt-0.5"
            haveAction={false}
          />
        }
      />
    </div>
  );
};

export default VariationLabel;
