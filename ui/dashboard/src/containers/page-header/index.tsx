import {
  IconHelpOutlineOutlined,
  IconNotificationsNoneOutlined
} from 'react-icons-material-design';
import Icon from 'components/icon';

type PageHeaderProps = {
  title: string;
  description: string;
};

const PageHeader = ({ title, description }: PageHeaderProps) => {
  return (
    <header className="py-8 px-6 border-b border-gray-200">
      <div className="flex justify-between mb-3">
        <div className="flex items-center">
          <h1 className="text-gray-900 typo-head-light-huge">{title}</h1>
        </div>
        <div className="flex items-center gap-4 text-gray-500">
          <Icon icon={IconHelpOutlineOutlined} size="xs" />
          <Icon icon={IconNotificationsNoneOutlined} size="xs" />
        </div>
      </div>
      <p className="text-gray-600">{description}</p>
    </header>
  );
};

export default PageHeader;
