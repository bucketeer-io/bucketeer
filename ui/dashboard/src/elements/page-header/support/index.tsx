import { useCallback } from 'react';
import {
  IconHelpOutlineOutlined,
  IconLaunchOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { urls } from 'configs';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { copyToClipBoard } from 'utils/function';
import { IconCopy } from '@icons';
import Button from 'components/button';
import Divider from 'components/divider';
import Icon from 'components/icon';
import { Popover } from 'components/popover';

const SupportPopoverItem = ({ to, title }: { to: string; title: string }) => (
  <Link
    to={to}
    target="_blank"
    className="flex items-center w-full gap-x-2 px-6 py-1 typo-para-small text-gray-700 hover:bg-gray-200 hover:text-primary-500 hover:underline transition-colors rounded"
  >
    {title}
    <Icon icon={IconLaunchOutlined} size="xxs" />
  </Link>
);

const SupportPopover = () => {
  const { t } = useTranslation(['form', 'message', 'common']);
  const { notify } = useToast();

  const handleCopy = useCallback(() => {
    copyToClipBoard(urls.API_ENDPOINT || '');
    notify({
      message: t('message:copied')
    });
  }, [urls]);

  return (
    <Popover
      align="end"
      sideOffset={6}
      className="py-1.5 shadow-card-secondary"
      trigger={
        <div className="flex-center size-fit">
          <Icon icon={IconHelpOutlineOutlined} size="sm" color="gray-500" />
        </div>
      }
    >
      <div className="flex flex-col w-full gap-y-0.5">
        {urls.API_ENDPOINT && (
          <>
            <div className="flex flex-col w-full gap-y-0.5 px-3 py-1 typo-para-small">
              <p className="text-gray-700">{t('sdk-api-endpoint')}</p>
              <div className="flex items-center gap-x-1.5">
                <p className="text-gray-600">{urls.API_ENDPOINT || ''}</p>
                <Button
                  variant="grey"
                  className="size-fit flex-center"
                  onClick={handleCopy}
                >
                  <Icon icon={IconCopy} size="sm" />
                </Button>
              </div>
            </div>
            <Divider />
          </>
        )}
        <div className="px-3">
          <p className="typo-para-small font-semibold text-gray-700 pt-1">
            {t('common:docs')}
          </p>
          <SupportPopoverItem
            to={DOCUMENTATION_LINKS.GETTING_STARTED}
            title={t('getting-started')}
          />
          <SupportPopoverItem
            to={DOCUMENTATION_LINKS.SDK_CONFIGURATION}
            title={t('sdk-configuration')}
          />
        </div>
      </div>
    </Popover>
  );
};

export default SupportPopover;
