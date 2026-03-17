import { useTranslation } from 'react-i18next';
import { IconLaunchOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { useScreen } from 'hooks';
import { cn } from 'utils/style';
import Icon from 'components/icon';
import CreatedAtTime from 'elements/page-details-header/created-at-time';
import PageLayout from 'elements/page-layout';
import SDKApiEndpoint from './sdk-api-endpoint';
import SupportPopover from './support';

interface PageHeaderProps {
  title: string;
  description: string;
  createdAt?: string;
  isShowApiEndpoint?: boolean;
  link?: string;
}

const PageHeader = ({
  title,
  description,
  createdAt,
  isShowApiEndpoint,
  link
}: PageHeaderProps) => {
  const { t } = useTranslation(['common']);
  const { fromTabletScreen } = useScreen();
  return (
    <PageLayout.Header>
      <div
        className={cn('flex justify-between gap-2', {
          'flex-col items-start': isShowApiEndpoint,
          'flex-row items-center': fromTabletScreen
        })}
      >
        <div className="flex items-center gap-2">
          <h1 className="text-gray-900 typo-head-bold-huge">{title}</h1>
          {createdAt && <CreatedAtTime createdAt={createdAt} />}
        </div>
        <div className="flex items-center gap-4 text-gray-500">
          {isShowApiEndpoint && <SDKApiEndpoint />}
          <SupportPopover />
        </div>
      </div>
      <div className="w-full inline">
        <p className="inline text-gray-600 mt-3 typo-para-small">
          {description}
        </p>
        {link && (
          <Link
            className="inline-flex items-center ml-1 typo-para-small text-primary-500 underline sm:hidden"
            target="_blank"
            to={link}
            onClick={e => {
              if (!link) return e.preventDefault();
            }}
          >
            {t('documentation')}
            <Icon icon={IconLaunchOutlined} size="xxs" />
          </Link>
        )}
      </div>
    </PageLayout.Header>
  );
};

export default PageHeader;
