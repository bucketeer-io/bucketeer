import { useMemo } from 'react';
import { hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import {
  IconPlus,
  IconPrerequisite,
  IconSetting,
  IconUserOutlined
} from '@icons';
import Dropdown, { DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import { ADD_NODE_LEFT_OFFSET_PX } from '../evaluation-flow';
import { RuleCategory } from '../types';

const AddRule = ({
  isDisableAddPrerequisite,
  isDisableAddIndividualRules,
  isInsertSegmentRule,
  indexInsertSegmentRule,
  onAddRule
}: {
  isDisableAddPrerequisite: boolean;
  isDisableAddIndividualRules: boolean;
  isInsertSegmentRule?: boolean;
  indexInsertSegmentRule?: number;
  onAddRule: (rule: RuleCategory, index?: number) => void;
}) => {
  const { t } = useTranslation(['form', 'table']);
  const { consoleAccount } = useAuth();
  const editable = hasEditable(consoleAccount!);

  const options = useMemo(
    () => [
      {
        label: t('feature-flags.prerequisites'),
        value: RuleCategory.PREREQUISITE,
        tooltip: t('targeting.prerequisite-tooltip'),
        icon: IconPrerequisite,
        disabled: isDisableAddPrerequisite
      },
      {
        label: t('targeting.individual-targeting'),
        value: RuleCategory.INDIVIDUAL,
        tooltip: t('targeting.individual-tooltip'),
        icon: IconUserOutlined,
        disabled: isDisableAddIndividualRules
      },
      {
        label: t('targeting.custom-rule'),
        value: RuleCategory.CUSTOM,
        tooltip: t('targeting.custom-rule-tooltip'),
        icon: IconSetting
      }
    ],
    [isDisableAddIndividualRules, isDisableAddPrerequisite]
  );

  const getRuleCategoryCall = (value: RuleCategory) => {
    if (value === RuleCategory.CUSTOM) {
      if (isInsertSegmentRule) {
        return onAddRule(RuleCategory.CUSTOM, indexInsertSegmentRule);
      }
    }
    return onAddRule(value);
  };

  const sharedDropdownProps = {
    isTooltip: true,
    options,
    disabled: !editable,
    showArrow: false,
    onChange: (value: DropdownValue | DropdownValue[]) =>
      getRuleCategoryCall(value as RuleCategory)
  };

  return (
    <>
      {/* Spine plus — its own dropdown so the menu opens *below the plus
          circle*, not below the centred "+ Add Rule" button. Absolutely
          positioned to align with the EvaluationFlow spine. */}
      <div
        className="absolute top-1/2 -translate-y-1/2 z-10"
        style={{ left: `${ADD_NODE_LEFT_OFFSET_PX}px` }}
      >
        <Dropdown
          {...sharedDropdownProps}
          isTruncate={false}
          alignContent="start"
          wrapTriggerStyle="!w-fit"
          // `[&>div]:overflow-visible` lets the spine plus's `ring-4 ring-white`
          // halo (which masks the spine line behind it) extend beyond the
          // dropdown trigger's default overflow-hidden.
          className="!w-fit !p-0 !border-0 !shadow-none !bg-transparent [&>div]:overflow-visible"
          trigger={
            <DisabledButtonTooltip
              align="center"
              hidden={editable}
              trigger={
                <span
                  className={cn(
                    'flex size-5 items-center justify-center rounded-full bg-white dark:bg-dark-black-800 border border-dashed ring-4 ring-white dark:ring-dark-black-800 transition-colors',
                    editable
                      ? 'border-gray-400 dark:border-dark-black-700 text-gray-500 dark:text-dark-gray-200 cursor-pointer hover:border-primary-500 dark:hover:border-dark-purple-300 hover:text-primary-500 dark:hover:text-dark-purple-700'
                      : 'border-gray-300 dark:border-dark-black-700 text-gray-400 dark:text-dark-gray-200 cursor-not-allowed opacity-60'
                  )}
                >
                  <span className="sr-only">
                    {t('table:feature-flags.add-rule')}
                  </span>
                  <Icon icon={IconPlus} size="xxs" />
                </span>
              }
            />
          }
        />
      </div>

      {/* Visible "+ Add Rule" dashed button. */}
      <Dropdown
        {...sharedDropdownProps}
        alignContent="center"
        className="w-full [&>div]:flex-center border-dashed !shadow-none dark:border-dark-black-700"
        trigger={
          <DisabledButtonTooltip
            align="center"
            hidden={editable}
            trigger={
              <div className="flex items-center gap-x-2 h-6 p-0 typo-para-medium !text-primary-500  dark:!text-dark-purple-400">
                <Icon icon={IconPlus} size={'md'} />
                {t('table:feature-flags.add-rule')}
              </div>
            }
          />
        }
      />
    </>
  );
};

export default AddRule;
