import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';

const VariationLabel = ({ label, index }: { label: string; index: number }) => {
  return (
    <div className="flex items-center gap-x-2 pl-0.5">
      <FlagVariationPolygon index={index} />
      <p className="-mt-0.5">{label}</p>
    </div>
  );
};

export default VariationLabel;
