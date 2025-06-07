import { ReactNode } from 'react';
import Navigation from 'components/navigation';

const DefaultEnvLayout = ({ children }: { children: ReactNode }) => {
  return (
    <div className="flex size-full">
      <Navigation onClickNavLink={() => {}} />
      <div className="w-full ml-[248px] shadow-lg overflow-y-auto">
        {children}
      </div>
    </div>
  );
};

export default DefaultEnvLayout;
