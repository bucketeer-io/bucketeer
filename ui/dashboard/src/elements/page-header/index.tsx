import {
  IconHelpOutlineOutlined,
  IconNotificationsNoneOutlined
} from 'react-icons-material-design';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import PageLayout from 'elements/page-layout';
import Notifications from './notifications';

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
          <button className="flex-center size-fit">
            <Icon icon={IconHelpOutlineOutlined} size="sm" />
          </button>
          <Popover
            trigger={
              <Icon
                icon={IconNotificationsNoneOutlined}
                size="sm"
                color="gray-500"
              />
            }
            align="end"
            sideOffset={24}
            closeBtnCls="!flex-center absolute right-4 top-6"
            className="w-[494px] p-0 max-h-[500px]"
          >
            <Notifications />
          </Popover>
        </div>
      </div>
      <p className="text-gray-600 mt-3 text-sm">{description}</p>
    </PageLayout.Header>
  );
};

export default PageHeader;
