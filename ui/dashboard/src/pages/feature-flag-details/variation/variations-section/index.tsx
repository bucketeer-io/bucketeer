import { VariationProps } from '..';
import Card from '../../elements/card';
import VariationList from './variation-list';

const VariationsSection = ({
  editable,
  feature,
  isRunningExperiment
}: VariationProps) => {
  return (
    <Card className="divide-y divide-gray-900/10 dark:divide-dark-black-700">
      <VariationList
        editable={editable}
        feature={feature}
        isRunningExperiment={isRunningExperiment}
      />
    </Card>
  );
};

export default VariationsSection;
