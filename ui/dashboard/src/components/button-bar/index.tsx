import { ReactNode } from 'react';

export const ButtonBar = ({
  primaryButton,
  secondaryButton
}: {
  primaryButton: ReactNode;
  secondaryButton?: ReactNode;
}) => {
  return (
    <div className="p-5 border-t border-gray-200 w-full">
      <div className="flex items-center gap-4 justify-end">
        {primaryButton}
        {secondaryButton}
      </div>
    </div>
  );
};
