import { useScreen } from 'hooks';
import { cn } from 'utils/style';
import PageLayout from 'elements/page-layout';
import SDKApiEndpoint from './sdk-api-endpoint';
import SupportPopover from './support';

interface PageHeaderProps {
  title: string;
  description: string;
  isShowApiEndpoint?: boolean;
}

const PageHeader = ({
  title,
  description,
  isShowApiEndpoint
}: PageHeaderProps) => {
  const { fromTabletScreen } = useScreen();
  return (
    <PageLayout.Header>
      <div
        className={cn('flex justify-between gap-2', {
          'flex-col items-start': isShowApiEndpoint,
          'flex-row items-center': fromTabletScreen
        })}
      >
        <h1 className="text-gray-900 typo-head-bold-huge">{title}</h1>
        <div className="flex items-center gap-4 text-gray-500">
          {isShowApiEndpoint && <SDKApiEndpoint />}
          <SupportPopover />
        </div>
      </div>
      <p className="text-gray-600 mt-3 typo-para-small">{description}</p>
    </PageLayout.Header>
  );
};

export default PageHeader;
