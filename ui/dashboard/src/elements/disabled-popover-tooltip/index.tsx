import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { useAuthAccess } from 'auth';
import { Popover, PopoverProps, PopoverValue } from 'components/popover';

interface Props extends PopoverProps<PopoverValue> {
  isNeedAdminAccess?: boolean;
}

const DisabledPopoverTooltip = ({
  isNeedAdminAccess = false,
  options,
  icon,
  align = 'end',
  ...props
}: Props) => {
  const { t } = useTranslation(['common']);
  const { envEditable, isOrganizationAdmin } = useAuthAccess();

  const tooltipContent = useMemo(() => {
    if (envEditable && isOrganizationAdmin) return '';
    if (!envEditable) return t('disabled-button-tooltip');
    return isNeedAdminAccess && !isOrganizationAdmin
      ? t('need-admin-access-tooltip')
      : '';
  }, [envEditable, isOrganizationAdmin, isNeedAdminAccess]);

  const formattedOptions = useMemo(() => {
    if (!tooltipContent) return options;
    return (options || []).map(item => ({
      ...item,
      tooltip: tooltipContent || item?.tooltip,
      disabled: item.disabled || !!tooltipContent || !!item?.tooltip
    }));
  }, [tooltipContent, options]);

  return (
    <Popover
      {...props}
      options={formattedOptions}
      icon={icon || IconMoreHorizOutlined}
      align={align}
    />
  );
};

export default DisabledPopoverTooltip;
