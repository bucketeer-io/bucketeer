import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { cn } from 'utils/style';
import {
  IconFlagJSON,
  IconFlagNumber,
  IconFlagString,
  IconFlagSwitch,
  IconInfo
} from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';

const SubmitBar = ({
  editable,
  feature,
  onShowConfirmDialog
}: {
  editable: boolean;
  feature: Feature;
  onShowConfirmDialog: () => void;
}) => {
  const { t } = useTranslation(['common', 'table']);

  const {
    formState: { isDirty, isValid }
  } = useFormContext();

  const flagTypeIcon = useMemo(() => {
    const { variationType } = feature;
    if (variationType === 'JSON') return IconFlagJSON;
    if (variationType === 'BOOLEAN') return IconFlagSwitch;
    if (variationType === 'NUMBER') return IconFlagNumber;
    return IconFlagString;
  }, [feature]);

  return (
    <div className="flex items-center justify-between w-full gap-x-6">
      <div className="flex items-center gap-x-2">
        <h3 className="typo-head-bold-small text-gray-800">
          {t('table:feature-flags.variation')}
        </h3>
        <Tooltip
          align="start"
          alignOffset={-50}
          content={t('table:feature-flags.variation-type-tooltip')}
          trigger={
            <div className="flex-center h-full">
              <Icon
                icon={IconInfo}
                size={'xxs'}
                color="gray-500"
                className="flex-center -mb-0.5"
              />
            </div>
          }
          className="max-w-[300px]"
        />
        <Icon icon={flagTypeIcon} size="sm" />
        <p
          className={cn('typo-para-small text-gray-600 capitalize', {
            uppercase: feature.variationType === 'JSON'
          })}
        >
          {feature?.variationType?.toLowerCase()}
        </p>
      </div>
      <DisabledButtonTooltip
        hidden={editable}
        trigger={
          <Button
            type="button"
            disabled={!isDirty || !isValid || !editable}
            onClick={onShowConfirmDialog}
          >
            {t('save-with-comment')}
          </Button>
        }
      />
    </div>
  );
};

export default SubmitBar;
