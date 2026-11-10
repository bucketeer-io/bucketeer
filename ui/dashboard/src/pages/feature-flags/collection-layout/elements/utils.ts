import { FEATURE_RECENT_TRAFFIC_DAYS } from 'constants/feature-flag';
import dayjs from 'dayjs';
import { Feature, FeatureVariationType } from '@types';
import {
  IconFlagJSON,
  IconFlagNumber,
  IconFlagString,
  IconFlagSwitch,
  IconFlagYaml
} from '@icons';
import { FeatureActivityStatus } from 'pages/feature-flags/types';

export const getDataTypeIcon = (type: FeatureVariationType) => {
  if (type === 'BOOLEAN') return IconFlagSwitch;
  if (type === 'STRING') return IconFlagString;
  if (type === 'NUMBER') return IconFlagNumber;
  if (type === 'YAML') return IconFlagYaml;
  return IconFlagJSON;
};

export function getFlagStatus(feature: Feature): FeatureActivityStatus {
  const lastUsedAt = Number(feature.lastUsedAt);
  if (!lastUsedAt) {
    return FeatureActivityStatus.NEVER_USED;
  }

  const daysDifference = dayjs(new Date(lastUsedAt * 1000)).diff(
    dayjs(),
    'day'
  );
  if (daysDifference > -FEATURE_RECENT_TRAFFIC_DAYS) {
    return FeatureActivityStatus.RECEIVING_TRAFFIC;
  }
  return FeatureActivityStatus.NO_RECENT_TRAFFIC;
}
