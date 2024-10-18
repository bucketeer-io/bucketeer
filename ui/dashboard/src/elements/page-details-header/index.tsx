import {
  IconAccessTimeOutlined,
  IconArrowBackFilled
} from 'react-icons-material-design';
import { useNavigate } from 'react-router-dom';
import { cn } from 'utils/style';
import Icon from 'components/icon';

export type PageDetailHeaderProps = {
  title: string;
  description: string;
};

const PageDetailHeader = ({ title, description }: PageDetailHeaderProps) => {
  const navigate = useNavigate();

  return (
    <header className="grid pt-7 px-6">
      <button
        className={cn(
          'size-6 flex-center rounded hover:shadow-border-gray-500',
          'shadow-border-gray-400 text-gray-600'
        )}
        onClick={() => navigate(-1)}
      >
        <Icon icon={IconArrowBackFilled} size="xxs" />
      </button>
      <div className="text-gray-500 flex items-center gap-1.5 mt-4">
        <Icon icon={IconAccessTimeOutlined} size="xxs" />
        <p className="typo-para-small">{description}</p>
      </div>
      <h1 className="text-gray-900 typo-head-bold-huge mt-2">{title}</h1>
    </header>
  );
};

export default PageDetailHeader;
