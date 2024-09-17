import { ReactNode } from 'react';
import { IconFolder2Outlined } from '@icons';
import Icon from 'components/icon';

export type EmptyProps = {
  title: string;
  description: string;
  emptyActions?: ReactNode;
};

const EmptyData = ({ title, description, emptyActions }: EmptyProps) => {
  return (
    <div className="grid gap-6 w-[352px]">
      <div className="grid gap-4">
        <div className="flex justify-center">
          <div className="flex-center size-16 rounded-full bg-additional-gray-50">
            <div className="flex-center size-12 rounded-full bg-additional-gray-100 ">
              <Icon icon={IconFolder2Outlined} />
            </div>
          </div>
        </div>
        <div className="text-center">
          <h3 className="text-additional-gray-300 typo-head-bold-medium mb-2">
            {title}
          </h3>
          <p className="text-additional-gray-200 typo-para-small">
            {description}
          </p>
        </div>
      </div>
      {emptyActions}
    </div>
  );
};

export default EmptyData;
