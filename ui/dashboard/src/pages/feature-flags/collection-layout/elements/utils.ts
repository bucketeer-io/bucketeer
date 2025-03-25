import { isNumber } from 'lodash';
import { Feature, FeatureVariationType } from '@types';
import {
  IconFlagJSON,
  IconFlagNumber,
  IconFlagString,
  IconFlagSwitch
} from '@icons';
import { FeatureActivityStatus } from 'pages/feature-flags/types';

export const getDataTypeIcon = (type: FeatureVariationType) => {
  if (type === 'BOOLEAN') return IconFlagSwitch;
  if (type === 'STRING') return IconFlagString;
  if (type === 'NUMBER') return IconFlagNumber;
  return IconFlagJSON;
};

export function getFlagStatus(feature: Feature): FeatureActivityStatus {
  if (!feature.lastUsedInfo) {
    return FeatureActivityStatus.NEW;
  }
  const { lastUsedAt } = feature?.lastUsedInfo || {};

  if (lastUsedAt && isNumber(+lastUsedAt)) {
    const _lastUsedAt = new Date(+lastUsedAt * 1000);
    if (_lastUsedAt.getDate() - new Date().getDate() > -7)
      return FeatureActivityStatus.ACTIVE;
  }
  return FeatureActivityStatus.INACTIVE;
}
