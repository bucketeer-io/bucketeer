import { memo } from 'react';
import { Suggestion } from '@types';
import { cn } from 'utils/style';

interface SuggestionCardProps {
  suggestion: Suggestion;
  onClick: (text: string) => void;
}

const SuggestionCard = memo(
  ({ suggestion, onClick }: SuggestionCardProps) => {
    return (
      <button
        type="button"
        className={cn(
          'w-full rounded-lg border border-gray-200 p-3 text-left',
          'transition-colors duration-200 hover:border-primary-300 hover:bg-primary-50',
          'motion-reduce:transition-none',
          'focus:outline-none focus:ring-2 focus:ring-primary-300 focus:ring-offset-2'
        )}
        onClick={() => onClick(suggestion.title)}
      >
        <p className="typo-para-small font-medium text-gray-700">
          {suggestion.title}
        </p>
        <p className="typo-para-tiny mt-1 text-gray-500">
          {suggestion.description}
        </p>
      </button>
    );
  }
);

SuggestionCard.displayName = 'SuggestionCard';

export default SuggestionCard;
