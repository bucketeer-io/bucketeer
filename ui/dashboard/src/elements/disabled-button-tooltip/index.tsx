import { useTranslation } from 'i18n';
import { Tooltip, TooltipProps } from 'components/tooltip';

const DisabledButtonTooltip = ({
  align = 'end',
  trigger,
  hidden,
  ...props
}: TooltipProps) => {
  const { t } = useTranslation(['common']);
  return (
    <Tooltip
      align={align}
      hidden={hidden}
      content={t('disabled-button-tooltip')}
      trigger={trigger}
      className="max-w-[300px]"
      {...props}
    />
  );
};

export default DisabledButtonTooltip;
