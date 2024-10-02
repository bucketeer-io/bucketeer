import { ReactNode } from 'react';
import {
  IconAccessTimeOutlined,
  IconArrowBackFilled
} from 'react-icons-material-design';
import { useNavigate } from 'react-router-dom';
import Button from 'components/button';
import Icon from 'components/icon';
import Tab from 'components/tab';
import { TabItemProps } from 'components/tab/tab-item';
import { StatusTag } from 'components/tag';
import { StatusTagType } from 'components/tag/status-tag';

export type PageDetailHeaderProps = {
  title: string;
  description: string;
  tabs: TabItemProps[];
  targetTab: string;
  navigateRoute: string;
  titleActions?: ReactNode;
  status?: StatusTagType;
  onSelectTab: (value: string) => void;
};

const PageDetailHeader = ({
  title,
  description,
  tabs,
  targetTab,
  navigateRoute,
  titleActions,
  status,
  onSelectTab
}: PageDetailHeaderProps) => {
  const navigate = useNavigate();

  const handleBack = () => navigate(navigateRoute);

  return (
    <header className="grid gap-6 pt-8 px-6">
      <div className="grid gap-4">
        <Button
          variant="secondary"
          className="size-6 p-0 shadow-border-gray-400 text-gray-600"
          onClick={handleBack}
        >
          <Icon icon={IconArrowBackFilled} size="xxs" />
        </Button>
        <div className="text-gray-500 flex items-center gap-2">
          <Icon icon={IconAccessTimeOutlined} size="xxs" />
          <p className="typo-para-small">{description}</p>
        </div>
        <div className="flex items-center justify-between">
          <div className="flex items-end gap-2">
            <h1 className="text-gray-900 typo-head-light-huge">{title}</h1>
            {status && <StatusTag variant={status} />}
          </div>
          {titleActions}
        </div>
      </div>
      <div>
        <Tab options={tabs} value={targetTab} onSelect={onSelectTab} />
      </div>
    </header>
  );
};

export default PageDetailHeader;
