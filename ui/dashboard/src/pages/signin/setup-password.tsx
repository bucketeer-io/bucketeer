import { useCallback, useEffect, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import {
  IconRemoveRedEyeOutlined,
  IconVisibilityOffOutlined
} from 'react-icons-material-design';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { setupPassword, validateSetUpPassword } from '@api/auth/setup-password';
import { yupResolver } from '@hookform/resolvers/yup';
import { PAGE_PATH_AUTH_SIGNIN } from 'constants/routing';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { SetupPasswordForm } from '@types';
import AuthWrapper from 'pages/signin/elements/auth-wrapper';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    password: yup
      .string()
      .required(requiredMessage)
      .min(
        4,
        translation('message:validation.name-at-least-characters', {
          name: translation('auth:password').toLowerCase(),
          count: 4
        })
      ),
    confirmPassword: yup
      .string()
      .oneOf(
        [yup.ref('password'), ''],
        translation('auth:password-confirm-not-match')
      )
      .required(requiredMessage)
  });

const SetupPassword = () => {
  const { t } = useTranslation(['auth', 'form', 'common', 'message']);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const { notify, errorNotify } = useToast();
  const [searchParams] = useSearchParams();
  const resetToken = searchParams.get('resetToken');
  const navigate = useNavigate();

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      password: '',
      confirmPassword: ''
    },
    mode: 'onSubmit'
  });

  const {
    formState: { isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<SetupPasswordForm> = useCallback(
    async values => {
      try {
        await setupPassword({
          resetToken: resetToken || '',
          newPassword: values.password
        });
        notify({
          toastType: 'info-message',
          message: t('message:setup-password-success')
        });
        navigate(PAGE_PATH_AUTH_SIGNIN);
      } catch (error) {
        errorNotify(error, t('message:reset-password-faild'));
      }
    },
    [resetToken]
  );

  // const PasswordAddonAction = ({
  //   type
  // }: {
  //   type: 'password' | 'confirm-password';
  // }) => (
  //   <Button
  //     type="button"
  //     variant="grey"
  //     className="text-gray-500 size-6"
  //     onClick={() =>
  //       type === 'password'
  //         ? setShowPassword(!showPassword)
  //         : setShowConfirmPassword(!showConfirmPassword)
  //     }
  //   >
  //     <Icon
  //       icon={
  //         (showPassword && type === 'password') ||
  //         (showConfirmPassword && type === 'confirm-password')
  //           ? IconVisibilityOffOutlined
  //           : IconRemoveRedEyeOutlined
  //       }
  //       size="sm"
  //     />
  //   </Button>
  // );

  useEffect(() => {
    const verifyToken = async () => {
      if (!resetToken) {
        errorNotify('error', t('message:invalid-setup-token-password'));
        return navigate(PAGE_PATH_AUTH_SIGNIN);
      }
      try {
        await validateSetUpPassword(resetToken);
      } catch (error) {
        errorNotify(error, t('message:invalid-setup-token-password'));
        return navigate(PAGE_PATH_AUTH_SIGNIN);
      }
    };
    verifyToken();
  }, [resetToken]);
  return (
    <AuthWrapper>
      <div className="grid gap-10">
        <div>
          <h1 className="text-gray-900 typo-head-bold-huge mb-4">
            {t(`enter-password.title`)}
          </h1>
          <p className="text-gray-600 typo-para-medium">
            {t(`enter-password.description`)}
          </p>
        </div>

        <FormProvider {...form}>
          <Form
            onSubmit={form.handleSubmit(onSubmit)}
            className="flex flex-col w-full gap-y-5"
          >
            <Form.Field
              control={form.control}
              name="password"
              render={({ field }) => (
                <Form.Item className="py-0">
                  <Form.Label required>{t(`password`)}</Form.Label>
                  <Form.Control>
                    <InputGroup
                      className="w-full"
                      addonSlot="right"
                      addon={
                        <Button
                          type="button"
                          variant="grey"
                          className="text-gray-500 size-6"
                          onClick={() => setShowPassword(!showPassword)}
                        >
                          <Icon
                            icon={
                              showPassword
                                ? IconVisibilityOffOutlined
                                : IconRemoveRedEyeOutlined
                            }
                            size="sm"
                          />
                        </Button>
                      }
                    >
                      <Input
                        type={showPassword ? 'text' : 'password'}
                        placeholder={t('password-placeholder')}
                        autoComplete="current-password"
                        {...field}
                      />
                    </InputGroup>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <Form.Item className="py-0">
                  <Form.Label required>{t(`password-confirm`)}</Form.Label>
                  <Form.Control>
                    <InputGroup
                      className="w-full"
                      addonSlot="right"
                      addon={
                        <Button
                          type="button"
                          variant="grey"
                          className="text-gray-500 size-6"
                          onClick={() =>
                            setShowConfirmPassword(!showConfirmPassword)
                          }
                        >
                          <Icon
                            icon={
                              showConfirmPassword
                                ? IconVisibilityOffOutlined
                                : IconRemoveRedEyeOutlined
                            }
                            size="sm"
                          />
                        </Button>
                      }
                    >
                      <Input
                        type={showConfirmPassword ? 'text' : 'password'}
                        placeholder={t('password-confirm-placeholder')}
                        autoComplete="current-password"
                        {...field}
                      />
                    </InputGroup>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Button
              type="submit"
              loading={isSubmitting}
              // disabled={!isValid}
              className="mt-5 w-full"
            >
              {t('common:continue')}
            </Button>
          </Form>
        </FormProvider>
      </div>
    </AuthWrapper>
  );
};

export default SetupPassword;
