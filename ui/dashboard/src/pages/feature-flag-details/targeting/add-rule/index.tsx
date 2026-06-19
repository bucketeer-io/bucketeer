import { useMemo } from 'react';
import { hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import {
  IconPlus,
  IconPrerequisite,
  IconSetting,
  IconUserOutlined
} from '@icons';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
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
    onChange: (value: string | number) =>
      getRuleCategoryCall(value as RuleCategory)
  };

  return (
    <>
      {/* Spine plus — its own dropdown so the menu opens *below the plus
          circle*, not below the centred "+ Add Rule" button. Absolutely
          positioned to align with the EvaluationFlow spine. */}
      <div
        className="absolute top-1/2 -translate-y-1/2 z-10"
        style={{ left: '-52px' }}
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
                  aria-label={t('table:feature-flags.add-rule')}
                  className="flex size-5 items-center justify-center rounded-full bg-white border border-dashed border-gray-400 ring-4 ring-white hover:border-primary-500 hover:text-primary-500 transition-colors cursor-pointer"
                >
                  <Icon icon={IconPlus} size="xxs" color="gray-500" />
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
        className="w-full [&>div]:flex-center border-dashed !shadow-none"
        trigger={
          <DisabledButtonTooltip
            align="center"
            hidden={editable}
            trigger={
              <div className="flex items-center gap-x-2 h-6 p-0 typo-para-medium !text-primary-500">
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
