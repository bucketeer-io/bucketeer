import React, { memo, FC } from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';
import EmailIconSvg from '../../assets/svg/email-icon.svg';
import GoogleIconSvg from '../../assets/svg/google-icon.svg';

import AuthWrapper from './authWrapper';

import { PAGE_PATH_AUTH_SIGNIN } from '../../constants/routing';
import { AppState } from '../../modules';
import { redirectToAuthUrl } from '../../modules/auth';
import { AppDispatch } from '../../store';
import { DEMO_SIGN_IN_ENABLED } from '../../config';

const SignIn: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const isLoading = useSelector<AppState, boolean>(
    (state) => state.auth.loading,
    shallowEqual
  );

  const handleGoogleLogin = () => {
    dispatch(redirectToAuthUrl());
  };

  return (
    <AuthWrapper>
      <h3 className="font-semibold text-xl">Sign in to Bucketeer</h3>
      <div className="mt-6 space-y-4">
        {DEMO_SIGN_IN_ENABLED && (
          <Link to={PAGE_PATH_AUTH_SIGNIN}>
            <button className="flex h-10 justify-center border items-center rounded w-full space-x-2 hover:border-gray-500 hover:bg-gray-50 transition-all duration-300">
              <EmailIconSvg />
              <span className="text-sm text-gray-600">Sign in with email</span>
            </button>
          </Link>
        )}
        <button
          className="flex h-10 justify-center border items-center rounded w-full space-x-2 hover:border-gray-500 hover:bg-gray-50 transition-all duration-300"
          onClick={handleGoogleLogin}
          disabled={isLoading}
        >
          {isLoading ? (
            <div className="spinner" />
          ) : (
            <>
              <GoogleIconSvg />
              <span className="text-sm text-gray-600">Sign in with Google</span>
            </>
          )}
        </button>
        {/* <button className="flex h-10 justify-center border items-center rounded w-full space-x-2">
          <img
            src="/assets/svg/github-icon.svg"
            alt="email logo"
            className="w-[18px]"
          />
          <span className="text-sm text-gray-600">Log in With Github</span>
        </button> */}
      </div>
    </AuthWrapper>
  );
});

export default SignIn;
