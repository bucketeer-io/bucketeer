import axiosClient from '@api/axios-client';
import {
  AuthResponse,
  ForgotPasswordForm,
  ResetPassword,
  UpdatePasswordForm
} from '@types';

export interface ValidateResetPassword {
  isValid: boolean;
  email: string;
}

export const setupPassword = (formValues: ResetPassword): Promise<void> => {
  return axiosClient
    .post('/v1/auth/setup-password', formValues)
    .then(response => response.data);
};

export const forgotPassword = (
  formValue: ForgotPasswordForm
): Promise<{ message: string }> => {
  return axiosClient
    .post('/v1/auth/password/reset/inintate', formValue)
    .then(response => response.data);
};

export const resetPassword = (formValue: ResetPassword): Promise<void> => {
  return axiosClient
    .post('/v1/auth/password/reset', formValue)
    .then(response => response.data);
};

export const validateResetPassword = (
  resetToken: string
): Promise<ValidateResetPassword> => {
  return axiosClient
    .post('/v1/auth/password/reset/validate', { resetToken })
    .then(response => response.data);
};

export const validateSetUpPassword = (
  resetToken: string
): Promise<ValidateResetPassword> => {
  return axiosClient
    .post('/v1/auth/password/setup/validate', { resetToken })
    .then(response => response.data);
};

export const updatePassword = (
  formValue: UpdatePasswordForm
): Promise<AuthResponse> => {
  return axiosClient
    .post('/v1/auth/password', formValue)
    .then(response => response.data);
};
