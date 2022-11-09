import { Variation } from '@/proto/feature/variation_pb';

export const createVariationLabel = (variation: Variation.AsObject) => {
  if (variation == null) {
    return 'None';
  }
  const maxLength = 150;
  const ellipsis = '...';
  const label = variation.name
    ? variation.name + ' - ' + variation.value
    : variation.value;
  if (label.length > maxLength) {
    return label.slice(0, maxLength - ellipsis.length) + ellipsis;
  }
  return label;
};
