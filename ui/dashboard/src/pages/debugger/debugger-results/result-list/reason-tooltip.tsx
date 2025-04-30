import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import { EvaluationReason } from '@types';
import { IconInfoFilled } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

const ReasonTooltip = ({ reason }: { reason: EvaluationReason }) => {
  const { t } = useTranslation(['table']);
  const reasonContentKey = useMemo(() => {
    switch (reason.type) {
      case 'TARGET':
        return 'reason-target-tooltip';
      case 'RULE':
        return 'reason-rule-tooltip';
      case 'OFF_VARIATION':
        return 'reason-offVariation-tooltip';
      case 'PREREQUISITE':
        return 'reason-prerequisite-tooltip';
      case 'DEFAULT':
      default:
        return 'reason-default-tooltip';
    }
  }, [reason]);

  return (
    <Tooltip
      align="end"
      content={t(reasonContentKey)}
      trigger={
        <div className="flex items-center w-fit gap-x-2">
          <p className="typo-para-medium text-gray-700 capitalize">
            {(reason.type || 'DEFAULT')?.replace('_', ' ')?.toLowerCase()}
          </p>
          <Icon icon={IconInfoFilled} color="gray-500" />
        </div>
      }
    />
  );
};

export default ReasonTooltip;
