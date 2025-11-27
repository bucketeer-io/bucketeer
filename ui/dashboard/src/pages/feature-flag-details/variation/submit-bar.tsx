import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { IconLaunchOutlined } from 'react-icons-material-design';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import useOptions from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
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
  const { flagTypeOptions } = useOptions();
  const {
    formState: { isDirty, isValid }
  } = useFormContext();

  const currentOption = useMemo(
    () => flagTypeOptions.find(item => item.value === feature.variationType),
    [feature, flagTypeOptions]
  );

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
        {currentOption && <Icon icon={currentOption?.icon} size="sm" />}
        <p
          className={cn('typo-para-small text-gray-600', {
            uppercase: feature.variationType === 'JSON'
          })}
        >
          {currentOption?.label}
        </p>
      </div>
      <div className="flex items-center gap-x-3">
        <Button
          type="button"
          variant="text"
          onClick={() =>
            window.open(DOCUMENTATION_LINKS.FLAG_VARIATIONS, '_blank')
          }
        >
          <Icon icon={IconLaunchOutlined} size="sm" />
          {t('documentation')}
        </Button>
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
    </div>
  );
};

export default SubmitBar;
