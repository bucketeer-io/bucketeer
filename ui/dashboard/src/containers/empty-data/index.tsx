import { ReactNode } from 'react';
import { IconNoData } from '@icons';

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
          <IconNoData />
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
