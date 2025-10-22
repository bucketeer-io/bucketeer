import { useEffect } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { forgotPassword } from '@api/auth/setup-password';
import { yupResolver } from '@hookform/resolvers/yup';
import { DEMO_SIGN_IN_ENABLED } from 'configs';
import { PAGE_PATH_AUTH_SIGNIN, PAGE_PATH_ROOT } from 'constants/routing';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { ForgotPasswordForm } from '@types';
import { IconBackspace } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import AuthWrapper from './elements/auth-wrapper';

const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    email: yup.string().email().required(requiredMessage)
  });

const ForgotPassword = () => {
  const { t } = useTranslation(['auth', 'common']);
  const { notify, errorNotify } = useToast();
  const navigate = useNavigate();

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      email: ''
    }
  });

  const onSubmit: SubmitHandler<ForgotPasswordForm> = async values => {
    try {
      const response = await forgotPassword(values);
      notify({ toastType: 'info-message', message: response.message });
    } catch (error) {
      errorNotify(error);
    }
  };

  useEffect(() => {
    if (!DEMO_SIGN_IN_ENABLED) {
      navigate(PAGE_PATH_ROOT);
    }
  }, []);

  return (
    <AuthWrapper>
      <div className="animate-translate-left w-full h-full">
        <Button
          variant="secondary-2"
          onClick={() => navigate(PAGE_PATH_AUTH_SIGNIN)}
          className="p-2 h-auto"
        >
          <Icon icon={IconBackspace} size="sm" />
        </Button>
        <div className="pt-5">
          <h1 className="text-gray-900 typo-head-bold-huge mb-4">
            {t(`forgot-password.title`)}
          </h1>
          <p className="text-gray-600 typo-para-medium">
            {t(`forgot-password.description`)}
          </p>
        </div>
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)} className="mt-6">
            <Form.Field
              control={form.control}
              name="email"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label>{t('email')}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t('email-placeholder')}
                      autoComplete="email"
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            <Button
              type="submit"
              loading={form.formState.isSubmitting}
              disabled={!form.formState.isValid}
              className="mt-8 w-full"
            >
              {t('common:continue')}
            </Button>
          </Form>
        </FormProvider>
      </div>
    </AuthWrapper>
  );
};

export default ForgotPassword;
