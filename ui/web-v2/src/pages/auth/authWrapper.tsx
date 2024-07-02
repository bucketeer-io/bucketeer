import React, { memo, FC } from 'react';
import logo2 from '../../assets/logo-2.png';
import featureFlagsDashboard from '../../assets/feature-flags-dashboard-list.png';

interface AuthWrapperProps {
  children: React.ReactNode;
}

const AuthWrapper: FC = memo(({ children }: AuthWrapperProps) => {
  return (
    <div className="flex h-full">
      <div className="flex-1 flex justify-center items-center">
        <div className="w-[400px]">
          <img src={logo2} alt="logo" className="w-44 mb-8" />
          {children}
        </div>
      </div>
      <div className="flex-1 p-4">
        <div className="bg-[#EFEFFE] flex items-center h-full rounded-3xl justify-end">
          <img
            src={featureFlagsDashboard}
            alt="feature flags dashboard"
            className="w-[92%]"
          />
        </div>
      </div>
    </div>
  );
});

export default AuthWrapper;
