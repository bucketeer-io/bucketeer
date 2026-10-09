// Includes the 140px chip and its horizontal margins.
const VALUE_WIDTH = 144;
// Leave a small buffer so fractional widths cannot clip the final chip.
const WIDTH_BUFFER = 12;
const MIN_VISIBLE_VALUES = 2;

export const getVisibleValueCount = (width: number, total: number) => {
  const availableWidth = Math.max(0, width - WIDTH_BUFFER);
  if (total * VALUE_WIDTH <= availableWidth) return total;

  const counterWidth = 24 + String(total).length * 10;
  return Math.min(
    total,
    Math.max(
      MIN_VISIBLE_VALUES,
      Math.floor((availableWidth - counterWidth) / VALUE_WIDTH)
    )
  );
};
