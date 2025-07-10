import { useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { yupResolver } from '@hookform/resolvers/yup';
import { PAGE_PATH_ROOT } from 'constants/routing';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { onGenerateSlug } from 'utils/converts';
import { IconBackspace } from '@icons';
import Button from 'components/button';
import Checkbox from 'components/checkbox';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import AuthWrapper from './elements/auth-wrapper';

interface AccessDemoForm {
  organizationName: string;
  organizationUrlCode: string;
  email: string;
  isAgree: boolean;
}

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    organizationName: yup.string().required(requiredMessage),
    organizationUrlCode: yup
      .string()
      .required(requiredMessage)
      .matches(
        /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
        translation('message:validation.id-rule', {
          name: translation('common:url-code')
        })
      ),
    email: yup.string().email().required(requiredMessage),
    isAgree: yup.boolean().required(requiredMessage)
  });

const AccessDemoPage = () => {
  const { t } = useTranslation(['auth', 'common', 'form']);
  const navigate = useNavigate();

  const [showAuthError, setShowAuthError] = useState(false);

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      organizationName: '',
      organizationUrlCode: '',
      email: '',
      isAgree: false
    },
    mode: 'onChange'
  });

  const onSubmit: SubmitHandler<AccessDemoForm> = async values => {
    try {
      console.log(values);
      navigate(PAGE_PATH_ROOT);
    } catch (error) {
      if (error) {
        setShowAuthError(true);
      }
    }
  };

  return (
    <AuthWrapper>
      <Button
        variant="secondary-2"
        onClick={() => navigate(PAGE_PATH_ROOT)}
        className="p-2 h-auto"
      >
        <Icon icon={IconBackspace} size="sm" />
      </Button>
      <h1 className="text-gray-900 typo-head-bold-huge mt-8">
        {t('access-demo-page')}
      </h1>
      <p className="text-gray-600 typo-para-medium mt-4">
        {t(`access-demo-page-desc`)}
      </p>
      {showAuthError && (
        <p className="text-accent-red-500 typo-para-medium mt-6">
          {t(`error-message.invalid-sign-in`)}
        </p>
      )}

      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)} className="mt-8">
          <Form.Field
            control={form.control}
            name="organizationName"
            render={({ field }) => (
              <Form.Item>
                <Form.Label required>{t('organization-name')}</Form.Label>
                <Form.Control>
                  <Input
                    placeholder={`${t('form:placeholder-name')}`}
                    {...field}
                    onChange={value => {
                      field.onChange(value);
                      const isUrlCodeDirty = form.getFieldState(
                        'organizationUrlCode'
                      ).isDirty;
                      const urlCode = form.getValues('organizationUrlCode');
                      form.setValue(
                        'organizationUrlCode',
                        isUrlCodeDirty ? urlCode : onGenerateSlug(value)
                      );
                    }}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="organizationUrlCode"
            render={({ field }) => (
              <Form.Item>
                <Form.Label required>{t('form:url-code')}</Form.Label>
                <Form.Control>
                  <Input
                    placeholder={`${t('form:placeholder-code')}`}
                    {...field}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="email"
            render={({ field }) => (
              <Form.Item>
                <Form.Label>{t('owner-email')}</Form.Label>
                <Form.Control>
                  <Input placeholder={t('email')} {...field} />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="isAgree"
            render={({ field }) => (
              <Form.Item>
                <Form.Control>
                  <Checkbox
                    title={t('agree-terms')}
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
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
            {t('common:submit')}
          </Button>
        </Form>
      </FormProvider>
    </AuthWrapper>
  );
};

export default AccessDemoPage;
