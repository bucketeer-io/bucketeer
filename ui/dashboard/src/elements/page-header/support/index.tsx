import {
  IconHelpOutlineOutlined,
  IconLaunchOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { urls } from 'configs';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useTranslation } from 'i18n';
import Divider from 'components/divider';
import Icon from 'components/icon';
import { Popover } from 'components/popover';

const SupportPopoverItem = ({ to, title }: { to: string; title: string }) => (
  <Link
    to={to}
    target="_blank"
    className="flex items-center w-full gap-x-2 px-3 py-1 typo-para-small text-gray-700 hover:bg-gray-200 hover:text-primary-500 hover:underline transition-colors rounded"
  >
    {title}
    <Icon icon={IconLaunchOutlined} size="xxs" />
  </Link>
);

const SupportPopover = () => {
  const { t } = useTranslation(['form']);
  return (
    <Popover
      align="end"
      className="py-1.5 shadow-card-secondary"
      trigger={
        <div className="flex-center size-fit">
          <Icon icon={IconHelpOutlineOutlined} size="sm" color="gray-500" />
        </div>
      }
    >
      <div className="flex flex-col w-full gap-y-0.5">
        <div className="flex flex-col w-full gap-y-0.5 px-3 py-1 typo-para-small">
          <p className="text-gray-700">{t('sdk-api-endpoint')}</p>
          <p className="text-gray-600">{urls.SDK_API_ENDPOINT || ''}</p>
        </div>
        <Divider />
        <SupportPopoverItem
          to={DOCUMENTATION_LINKS.GETTING_STARTED}
          title={t('getting-started')}
        />
        <SupportPopoverItem
          to={DOCUMENTATION_LINKS.SDK_CONFIGURATION}
          title={t('sdk-configuration')}
        />
      </div>
    </Popover>
  );
};

export default SupportPopover;
