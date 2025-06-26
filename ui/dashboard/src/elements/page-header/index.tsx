import PageLayout from 'elements/page-layout';
import SupportPopover from './support';

interface PageHeaderProps {
  title: string;
  description: string;
}

const PageHeader = ({ title, description }: PageHeaderProps) => {
  return (
    <PageLayout.Header>
      <div className="flex justify-between items-center">
        <h1 className="text-gray-900 typo-head-bold-huge">{title}</h1>
        <div className="flex items-center gap-3 text-gray-500">
          <SupportPopover />
        </div>
      </div>
      <p className="text-gray-600 mt-3 text-sm">{description}</p>
    </PageLayout.Header>
  );
};

export default PageHeader;
