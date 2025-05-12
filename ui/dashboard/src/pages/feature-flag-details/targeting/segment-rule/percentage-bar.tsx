import { COLORS } from 'constants/styles';
import { cn } from 'utils/style';

interface Props {
  weight?: number;
  currentIndex: number;
  isRoundedFull: boolean;
}

const PercentageBar = ({ weight, currentIndex, isRoundedFull }: Props) => {
  return (
    <div
      className={cn('first:rounded-l-full last:rounded-r-full h-2', {
        'rounded-full': isRoundedFull
      })}
      style={{
        width: `${weight || 0}%`,
        backgroundColor: COLORS[currentIndex % COLORS.length]
      }}
    />
  );
};

export default PercentageBar;
