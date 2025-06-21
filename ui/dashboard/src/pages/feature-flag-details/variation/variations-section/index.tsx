import { VariationProps } from '..';
import Card from '../../elements/card';
import VariationList from './variation-list';

const VariationsSection = ({
  editable,
  feature,
  isRunningExperiment
}: VariationProps) => {
  return (
    <Card className="divide-y divide-gray-900/10">
      <VariationList
        editable={editable}
        feature={feature}
        isRunningExperiment={isRunningExperiment}
      />
    </Card>
  );
};

export default VariationsSection;
