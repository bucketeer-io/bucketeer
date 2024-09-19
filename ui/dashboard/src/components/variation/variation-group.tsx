import { IconInfoOutlined } from 'react-icons-material-design';
import Icon from 'components/icon';
import Variation, { VariationProps } from './variation';

export type VariationGroupProps = {
  variations?: VariationProps[];
};

const VariationGroup = ({ variations = [] }: VariationGroupProps) => {
  return (
    <div className="flex items-center relative">
      {variations.map((i, index) => (
        <Variation
          key={index}
          text={variations.length === 1 ? i.text : ''}
          className={`relative`}
          style={{
            left: `-${index * 4}px`
          }}
          variant={i.variant}
        />
      ))}
      {variations.length > 1 && (
        <div className="flex gap-1">
          <p className="text-gray-700 typo-para-small">
            {variations.length} Variations
          </p>
          <span className="text-gray-500 grid place-items-center">
            <Icon icon={IconInfoOutlined} size="xs" />
          </span>
        </div>
      )}
    </div>
  );
};

export default VariationGroup;
