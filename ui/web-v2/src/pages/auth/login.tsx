import React, { memo, FC } from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
// import { Link } from 'react-router-dom';
import GoogleIconSvg from '../../assets/svg/google-icon.svg';

import AuthWrapper from './authWrapper';

// import { PAGE_PATH_AUTH_LOGIN } from '@/constants/routing';
import { AppState } from '../..//modules';
import { setupAuthToken } from '../..//modules/auth';
import { AppDispatch } from '../..//store';

const Login: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const isLoading = useSelector<AppState, boolean>(
    (state) => state.auth.loading,
    shallowEqual
  );

  const handleGoogleLogin = () => {
    dispatch(setupAuthToken());
  };

  return (
    <AuthWrapper>
      {/* <h2 className="font-bold text-xl">Log In</h2>
      <p className="mt-4 text-[#64738B] w-[90%]">
        To access our Demo site, please log in using the following information.
      </p>
      <p className="mt-6 text-[#64738B]">
        Email: demo@bucketeer.io
        <br />
        Password: demo
      </p> */}
      <div className="mt-8 space-y-4">
        {/* <Link to={PAGE_PATH_AUTH_LOGIN}>
          <button className="flex h-10 justify-center border items-center rounded w-full space-x-2 hover:border-gray-500 hover:bg-gray-50 transition-all duration-300">
            <img
              src="/assets/svg/email-icon.svg"
              alt="email logo"
              className="w-[18px]"
            />
            <span className="text-sm text-gray-600">Log in With Email</span>
          </button>
        </Link> */}
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
              <span className="text-sm text-gray-600">Log in With Google</span>
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

export default Login;
