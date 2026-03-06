import { memo } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';

interface PresetQuestionsProps {
  onSelect: (question: string) => void;
}

const presetQuestionKeys = [
  'ai-chat:presets.what-are-feature-flags',
  'ai-chat:presets.how-to-create-flag',
  'ai-chat:presets.ab-testing',
  'ai-chat:presets.progressive-rollout'
] as const;

const PresetQuestions = memo(({ onSelect }: PresetQuestionsProps) => {
  const { t } = useTranslation(['ai-chat']);

  return (
    <div className="flex flex-col gap-2">
      {presetQuestionKeys.map(key => {
        const text = t(key);
        return (
          <button
            key={key}
            type="button"
            className={cn(
              'rounded-lg border border-gray-200 px-3 py-2 text-left',
              'typo-para-tiny text-gray-600',
              'transition-colors duration-200 hover:border-primary-300 hover:bg-primary-50',
              'motion-reduce:transition-none',
              'focus:outline-none focus:ring-2 focus:ring-primary-300 focus:ring-offset-2'
            )}
            onClick={() => onSelect(text)}
          >
            {text}
          </button>
        );
      })}
    </div>
  );
});

PresetQuestions.displayName = 'PresetQuestions';

export default PresetQuestions;
