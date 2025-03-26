import { Feature } from '@types';
import Card from '../../elements/card';
import VariationList from './variation-list';
import VariationType from './variation-type';

const VariationsSection = ({ feature }: { feature: Feature }) => {
  return (
    <Card className="divide-y divide-gray-900/10">
      <VariationType feature={feature} />
      <VariationList feature={feature} />
    </Card>
  );
};

export default VariationsSection;
