import { FC, memo } from 'react';

import { classNames } from '../../utils/css';
import { CopyChip } from '../CopyChip';

export interface FeatureIdChipProps {
  featureId: string;
}

export const FeatureIdChip: FC<FeatureIdChipProps> = memo(({ featureId }) => {
  return (
    <CopyChip text={featureId}>
      <span
        className={classNames(
          'p-1.5 rounded-lg text-xs text-gray-800',
          'bg-gray-200 cursor-pointer'
        )}
      >
        {featureId}
      </span>
    </CopyChip>
  );
});
