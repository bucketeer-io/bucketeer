import { useEffect, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import {
  IconRemoveRedEyeOutlined,
  IconVisibilityOffOutlined
} from 'react-icons-material-design';
import { useNavigate } from 'react-router-dom';
import { signIn } from '@api/auth';
import { yupResolver } from '@hookform/resolvers/yup';
import { useAuth } from 'auth';
import {
  DEMO_SIGN_IN_EMAIL,
  DEMO_SIGN_IN_ENABLED,
  DEMO_SIGN_IN_PASSWORD
} from 'configs';
import { PAGE_PATH_ROOT } from 'constants/routing';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { SignInForm } from '@types';
import { IconBackspace } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import AuthWrapper from './elements/auth-wrapper';

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    email: yup.string().email().required(requiredMessage),
    password: yup
      .string()
      .required(requiredMessage)
      .min(
        4,
        translation('message:validation.name-at-least-characters', {
          name: translation('auth:password').toLowerCase(),
          count: 4
        })
      )
  });

const SignInWithEmail = () => {
  const { t } = useTranslation(['auth']);
  const { syncSignIn, setIsInitialLoading } = useAuth();
  const navigate = useNavigate();

  const [showPassword, setShowPassword] = useState(false);
  const [showAuthError, setShowAuthError] = useState(false);

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      email: '',
      password: ''
    }
  });

  const onSubmit: SubmitHandler<SignInForm> = async values => {
    try {
      setShowAuthError(false);
      const response = await signIn(values);
      await syncSignIn(response?.token);
      setIsInitialLoading(true);
      navigate(PAGE_PATH_ROOT);
    } catch (error) {
      if (error) {
        setShowAuthError(true);
      }
    }
  };

  useEffect(() => {
    if (!DEMO_SIGN_IN_ENABLED) {
      navigate(PAGE_PATH_ROOT);
    }
  }, []);

  const PasswordAddonAction = () => (
    <Button
      type="button"
      variant="grey"
      className="text-gray-500 size-6"
      onClick={() => setShowPassword(!showPassword)}
    >
      <Icon
        icon={
          showPassword ? IconVisibilityOffOutlined : IconRemoveRedEyeOutlined
        }
        size="sm"
      />
    </Button>
  );

  return (
    <AuthWrapper>
      <Button
        variant="secondary-2"
        onClick={() => navigate(PAGE_PATH_ROOT)}
        className="p-2 h-auto"
      >
        <Icon icon={IconBackspace} size="sm" />
      </Button>
      <h1 className="text-gray-900 typo-head-bold-huge mt-8">{`Sign in`}</h1>
      <p className="text-gray-600 typo-para-medium mt-4">
        {t(`sign-in.description`)}
      </p>
      <div className="text-gray-600 typo-para-medium mt-6">
        <p>{`${t('email')}: ${DEMO_SIGN_IN_EMAIL}`}</p>
        <p>{`${t('password')}: ${DEMO_SIGN_IN_PASSWORD}`}</p>
      </div>

      {showAuthError && (
        <p className="text-accent-red-500 typo-para-medium mt-6">
          {t(`error-message.invalid-sign-in`)}
        </p>
      )}

      <FormProvider {...form}>
        <Form
          onSubmit={form.handleSubmit(onSubmit)}
          onChange={() => setShowAuthError(false)}
          className="mt-8"
        >
          <Form.Field
            control={form.control}
            name="email"
            render={({ field }) => (
              <Form.Item>
                <Form.Label>{t('email')}</Form.Label>
                <Form.Control>
                  <Input placeholder={t('email')} {...field} />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="password"
            render={({ field }) => (
              <Form.Item>
                <Form.Label>{t('password')}</Form.Label>
                <Form.Control>
                  <InputGroup
                    className="w-full"
                    addonSlot="right"
                    addon={<PasswordAddonAction />}
                  >
                    <Input
                      type={showPassword ? 'text' : 'password'}
                      placeholder={t('password')}
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
            loading={form.formState.isSubmitting}
            className="mt-8 w-full"
          >
            {`Sign In`}
          </Button>
        </Form>
      </FormProvider>
    </AuthWrapper>
  );
};

export default SignInWithEmail;
