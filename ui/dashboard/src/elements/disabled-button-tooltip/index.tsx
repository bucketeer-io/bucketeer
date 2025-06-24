import { useTranslation } from 'i18n';
import { Tooltip, TooltipProps } from 'components/tooltip';

interface Props extends TooltipProps {
  type?: 'editor' | 'admin';
}

const DisabledButtonTooltip = ({
  align = 'end',
  trigger,
  hidden,
  type = 'editor',
  ...props
}: Props) => {
  const { t } = useTranslation(['common']);
  return (
    <Tooltip
      align={align}
      hidden={hidden}
      content={t(
        type === 'editor'
          ? 'disabled-button-tooltip'
          : 'need-admin-access-tooltip'
      )}
      trigger={trigger}
      className="max-w-[300px]"
      {...props}
    />
  );
};

export default DisabledButtonTooltip;
