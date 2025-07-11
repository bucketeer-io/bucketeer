import React, { memo } from 'react';
import authBackground from 'assets/logos/auth-background.svg';
import logo from 'assets/logos/logo-primary.svg';
import FormLoading from 'elements/form-loading';

interface AuthWrapperProps {
  children: React.ReactNode;
  isLoading?: boolean;
}

export const AuthWrapper = memo(({ children, isLoading }: AuthWrapperProps) => {
  return (
    <div className="flex min-h-screen overflow-auto">
      <div className="flex-1 grid place-items-center p-5 pr-0">
        <div className="w-full max-w-[442px]">
          {isLoading ? (
            <FormLoading />
          ) : (
            <div className="mb-12">
              <img src={logo} alt="bucketer-logo" />
            </div>
          )}
          {children}
        </div>
      </div>
      <div className="flex-1 p-4">
        <div className="bg-primary-additional flex items-center h-full rounded-3xl justify-end">
          <img
            src={authBackground}
            alt="feature flags dashboard"
            className="w-[92%]"
          />
        </div>
      </div>
    </div>
  );
});

export default AuthWrapper;
