import { PropsWithChildren } from 'react';
import { IconKeyboardArrowDownFilled } from 'react-icons-material-design';
import Icon from 'components/icon';

export type TagGroupProps = PropsWithChildren & {
  expandable?: boolean;
};

const TagGroup = ({ children, expandable }: TagGroupProps) => {
  return (
    <div className="flex gap-2 items-center">
      <div className="flex items-center gap-2">{children}</div>
      {expandable && (
        <div className="flex-center size-fit text-gray-500">
          <Icon icon={IconKeyboardArrowDownFilled} size="xs" />
        </div>
      )}
    </div>
  );
};

export default TagGroup;
