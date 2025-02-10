import { ReactNode } from 'react';
import Button from 'components/button';

interface Props {
  title: string;
  description: string;
  btnText: string;
  disabled?: boolean;
  children?: ReactNode;
  onClick: () => void;
}

const GoalActions = ({
  title,
  description,
  btnText,
  disabled = false,
  children,
  onClick
}: Props) => {
  return (
    <div className="flex flex-col w-full p-5 gap-y-5 shadow-card rounded-lg bg-white">
      <p className="text-gray-800 typo-head-bold-small">{title}</p>
      {description && (
        <p className="typo-para-small leading-[14px] text-gray-600">
          {description}
        </p>
      )}
      {children}
      <Button
        className="w-fit"
        type="button"
        variant={'secondary'}
        disabled={disabled}
        onClick={onClick}
      >
        {btnText}
      </Button>
    </div>
  );
};

export default GoalActions;
