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
    <header className="p-6 border-b border-gray-200">
      <div className="flex justify-between items-center">
        <h1 className="text-gray-900 typo-head-bold-huge">{title}</h1>
        <div className="flex items-center gap-3 text-gray-500">
          <button>
            <Icon icon={IconHelpOutlineOutlined} size="sm" />
          </button>
          <button>
            <Icon icon={IconNotificationsNoneOutlined} size="sm" />
          </button>
        </div>
      </div>
      <p className="text-gray-600 mt-3 text-sm">{description}</p>
    </header>
  );
};

export default PageHeader;
