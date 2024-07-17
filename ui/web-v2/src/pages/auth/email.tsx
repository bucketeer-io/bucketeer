import { PAGE_PATH_ROOT } from '../../constants/routing';
import { classNames } from '../../utils/css';
import {
  ArrowNarrowLeftIcon,
  EyeIcon,
  EyeOffIcon
} from '@heroicons/react/outline';
import { yupResolver } from '@hookform/resolvers/yup';
import React, { memo, FC, useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import { Link } from 'react-router-dom';

import AuthWrapper from './authWrapper';
import { loginSchema } from './formSchema';
import { signIn } from '../../modules/auth';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../../store';
import { getToken } from '../../storage/token';
import { useHistory } from 'react-router-dom';
import { AppState } from '../../modules';
import {
  DEMO_SIGN_IN_ENABLED,
  DEMO_SIGN_IN_EMAIL,
  DEMO_SIGN_IN_PASSWORD
} from '../../config';
import { IToast, removeToast, selectAll } from '../../modules/toasts';

type Inputs = {
  email: string;
  password: string;
};

const Email: FC = memo(() => {
  const methods = useForm({
    resolver: yupResolver(loginSchema),
    defaultValues: {
      email: '',
      password: ''
    },
    mode: 'onSubmit'
  });
  const {
    register,
    handleSubmit,
    formState: { errors, isDirty }
  } = methods;

  const [isPasswordVisible, setIsPasswordVisible] = useState(false);

  const history = useHistory();
  const dispatch = useDispatch<AppDispatch>();
  const isLoading = useSelector<AppState, boolean>(
    (state) => state.auth.loading,
    shallowEqual
  );
  const toasts = useSelector<AppState, IToast[]>((state) =>
    selectAll(state.toasts)
  );

  const onSubmit: SubmitHandler<Inputs> = (data) => {
    if (!DEMO_SIGN_IN_ENABLED) {
      return;
    }
    dispatch(
      signIn({
        email: data.email,
        password: data.password
      })
    ).then(() => {
      const token = getToken();

      if (token?.accessToken) {
        // Remove all toasts
        toasts.forEach((toast) => {
          dispatch(
            removeToast({
              id: toast.id
            })
          );
        });
        history.push('/');
      }
    });
  };

  return (
    <AuthWrapper>
      <Link to={PAGE_PATH_ROOT}>
        <button className="rounded border w-6 h-6 flex justify-center items-center mt-12">
          <ArrowNarrowLeftIcon width={16} />
        </button>
      </Link>
      <h2 className="font-bold text-xl mt-8">Sign In</h2>
      <p className="mt-3 text-[#64738B] w-[90%]">
        To access our Demo site, please sign in using the following information.
      </p>
      <p className="mt-6 text-[#64738B]">
        Email: {DEMO_SIGN_IN_EMAIL}
        <br />
        Password: {DEMO_SIGN_IN_PASSWORD}
      </p>
      {/* <div className="rounded-xl bg-red-50 p-4 mt-8">
        <div className="flex items-center">
          <div className="flex-shrink-0">
            <ExclamationCircleIcon
              className="h-5 w-5 text-red-400"
              aria-hidden="true"
            />
          </div>
          <div className="ml-3">
            <div className="text-sm text-red-600 font-medium">
              <p>Wrong email or password. Please double-check and try again</p>
            </div>
          </div>
        </div>
      </div> */}
      <form className="mt-6" onSubmit={handleSubmit(onSubmit)}>
        <div className="space-y-1 flex flex-col">
          <label htmlFor="email" className="text-sm text-gray-500">
            Email
          </label>
          <input
            {...register('email')}
            type="email"
            placeholder="Email"
            className={classNames(
              'border border-gray-300',
              errors.email ? 'input-text-error' : 'input-text'
            )}
          />
          <p className="input-error">
            {errors.email && <span role="alert">{errors.email.message}</span>}
          </p>
        </div>
        <div className="space-y-1 flex flex-col mt-3">
          <label htmlFor="email" className="text-sm text-gray-500">
            Password
          </label>
          <div className="w-full">
            <div className="relative">
              <input
                {...register('password')}
                placeholder="Password"
                className={classNames(
                  'border border-gray-300 w-full',
                  errors.password ? 'input-text-error' : 'input-text'
                )}
                type={isPasswordVisible ? 'text' : 'password'}
              />
              <button
                type="button"
                className="absolute right-2 inset-y-0 p-[2px]"
                onClick={() => setIsPasswordVisible(!isPasswordVisible)}
              >
                {isPasswordVisible ? (
                  <EyeOffIcon width={16} />
                ) : (
                  <EyeIcon width={16} />
                )}
              </button>
            </div>
            <p className="input-error">
              {errors.password && (
                <span role="alert">{errors.password.message}</span>
              )}
            </p>
            {/* <p className="text-red-600 text-sm mt-1">
              Wrong email or password. Try again or{' '}
              <strong className="underline cursor-pointer">
                create an account.
              </strong>
            </p> */}
          </div>
        </div>
        <button
          type="submit"
          className="btn btn-submit mt-8 w-full"
          disabled={!isDirty || Object.keys(errors).length > 0 || isLoading}
        >
          Log In
        </button>
      </form>
    </AuthWrapper>
  );
});

export default Email;
