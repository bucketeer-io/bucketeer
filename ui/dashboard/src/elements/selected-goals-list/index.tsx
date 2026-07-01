import { useTranslation } from 'i18n';
import { IconClose, IconInfo, IconOutperformed } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

// Renders the goals attached to an experiment as a list. The first goal is the
// primary goal — it drives the experiment verdict on the results page — so it
// carries a "Primary" badge with an explanatory tooltip, mirroring the
// LaunchDarkly metric list. Shared by the create/update modal and the
// experiment settings page.
const SelectedGoalsList = ({
  goalIds,
  getGoalName,
  editable = false,
  onRemove
}: {
  goalIds: string[];
  getGoalName: (id: string) => string;
  editable?: boolean;
  onRemove?: (id: string) => void;
}) => {
  const { t } = useTranslation(['form', 'common']);
  if (!goalIds?.length) return null;
  return (
    <div className="flex flex-col w-full mt-2 rounded-lg border border-gray-200 overflow-hidden">
      {goalIds.map((id, index) => (
        <div
          key={id}
          className="flex items-center justify-between gap-x-2 px-4 py-3 border-b border-gray-200 last:border-b-0"
        >
          <p className="typo-para-medium text-gray-700 truncate">
            {getGoalName(id)}
          </p>
          <div className="flex items-center gap-x-3 shrink-0">
            {index === 0 && (
              <Tooltip
                content={t('form:experiments.primary-goal-tooltip')}
                trigger={
                  <span className="flex items-center gap-x-1 typo-para-small font-medium text-primary-500 whitespace-nowrap">
                    <Icon
                      icon={IconOutperformed}
                      size="xxs"
                      color="primary-500"
                    />
                    {t('form:experiments.primary-goal')}
                    <Icon icon={IconInfo} size="xxs" color="gray-500" />
                  </span>
                }
                className="max-w-[300px]"
              />
            )}
            {editable && onRemove && (
              <button
                type="button"
                onClick={() => onRemove(id)}
                aria-label={t('form:experiments.remove-goal')}
                className="flex-center rounded p-1 hover:bg-gray-100"
              >
                <Icon icon={IconClose} size="xxs" color="gray-500" />
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
};

export default SelectedGoalsList;
