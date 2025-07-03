import dayjs from 'dayjs';
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
    return FeatureActivityStatus.NEVER_USED;
  }

  const { lastUsedAt } = feature?.lastUsedInfo || {};

  if (lastUsedAt && isNumber(+lastUsedAt)) {
    const _lastUsedAt = new Date(+lastUsedAt * 1000);
    const daysDifference = dayjs(_lastUsedAt).diff(dayjs(), 'day');

    if (daysDifference > -7) return FeatureActivityStatus.RECEIVING_TRAFFIC;
  }
  return FeatureActivityStatus.NO_RECENT_TRAFFIC;
}
