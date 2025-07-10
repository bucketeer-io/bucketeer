import { ReactNode } from 'react';
import { IconArrowBackFilled } from 'react-icons-material-design';
import { cn } from 'utils/style';
import Icon from 'components/icon';
import SupportPopover from 'elements/page-header/support';
import CreatedAtTime from './created-at-time';

export type PageDetailsHeaderProps = {
  title?: string;
  createdAt?: string;
  children?: ReactNode;
  additionElement?: ReactNode;
  onBack: () => void;
};

const PageDetailsHeader = ({
  title,
  createdAt,
  children,
  additionElement,
  onBack
}: PageDetailsHeaderProps) => {
  return (
    <header className="grid pt-7 px-6">
      <div className="flex items-start justify-between gap-x-2">
        <div className="flex items-start gap-4">
          <button
            className={cn(
              'size-6 min-w-6 flex-center rounded hover:shadow-border-gray-500',
              'shadow-border-gray-400 text-gray-600'
            )}
            onClick={onBack}
          >
            <Icon icon={IconArrowBackFilled} size="xxs" />
          </button>
          <div className="flex items-start gap-x-2">
            {title && (
              <h1 className="text-gray-900 flex-1 typo-head-bold-huge -mt-1.5">
                {title}
              </h1>
            )}
            {additionElement}
            {createdAt && <CreatedAtTime createdAt={createdAt} />}
          </div>
        </div>
        <SupportPopover />
      </div>
      {children}
    </header>
  );
};

export default PageDetailsHeader;
