import { ReactNode } from 'react';
import { cn } from 'utils/style';

export type TextProps = {
  text?: string;
  description?: string;
  sub?: ReactNode;
  isLink?: boolean;
};

const Text = ({ text, description, sub, isLink }: TextProps) => {
  return (
    <div>
      {text && (
        <div className="flex items-center">
          <p
            className={cn(
              'text-gray-700 typo-para-medium mr-2',
              isLink && 'underline text-primary-500'
            )}
          >
            {text}
          </p>
          {sub}
        </div>
      )}
      {description && (
        <p className="typo-para-tiny text-gray-700 mt-1">{description}</p>
      )}
    </div>
  );
};

export default Text;
