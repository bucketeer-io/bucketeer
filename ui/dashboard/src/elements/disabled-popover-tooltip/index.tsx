import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { useAuthAccess } from 'auth';
import { Popover, PopoverProps, PopoverValue } from 'components/popover';
import { formatPopoverOptions } from './utils';

interface Props extends PopoverProps<PopoverValue> {
  isNeedAdminAccess?: boolean;
  content?: string;
}

const DisabledPopoverTooltip = ({
  isNeedAdminAccess = false,
  options,
  icon,
  align = 'end',
  content,
  ...props
}: Props) => {
  const { t } = useTranslation(['common']);
  const { envEditable, isOrganizationAdmin } = useAuthAccess();

  const tooltipContent = useMemo(() => {
    if (content) return content;
    if (envEditable && isOrganizationAdmin) return '';
    if (isNeedAdminAccess && !isOrganizationAdmin)
      return t('need-admin-access-tooltip');
    if (!envEditable) return t('disabled-button-tooltip');
    return '';
  }, [envEditable, isOrganizationAdmin, isNeedAdminAccess, content]);

  const formattedOptions = useMemo(
    () => formatPopoverOptions(options, tooltipContent, content),
    [tooltipContent, options, content]
  );

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
