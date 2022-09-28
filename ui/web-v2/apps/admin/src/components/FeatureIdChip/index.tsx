import { FC, memo, useCallback, useState } from 'react';
import { useIntl } from 'react-intl';

import { classNames } from '../../utils/css';
import { CopyChip } from '../CopyChip';

export interface FeatureIdChipProps {
  featureId: string;
}

export const FeatureIdChip: FC<FeatureIdChipProps> = memo(({ featureId }) => {
  const { formatMessage: f } = useIntl();
  const [featureIdClicked, setFeatureIdClicked] = useState<boolean>(false);

  const handleFeatureIdClick = useCallback(
    (featureId: string) => {
      navigator.clipboard.writeText(featureId);
      setFeatureIdClicked(true);
    },
    [setFeatureIdClicked]
  );
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
