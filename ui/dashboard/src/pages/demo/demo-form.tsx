import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { organizationDemoCreator } from '@api/organization';
import { yupResolver } from '@hookform/resolvers/yup';
import { AxiosError } from 'axios';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import { clearDemoTokenStorage } from 'storage/demo-token';
import * as yup from 'yup';
import { onGenerateSlug } from 'utils/converts';
import Button from 'components/button';
import Checkbox from 'components/checkbox';
import Form from 'components/form';
import Input from 'components/input';

interface AccessDemoForm {
  organizationName: string;
  organizationUrlCode: string;
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
    isAgree: yup
      .boolean()
      .isTrue(translation('message:required-agreement-terms'))
      .required(requiredMessage)
  });

const DemoForm = ({
  isDemoSiteEnabled,
  onDemoAuthenticated
}: {
  isDemoSiteEnabled?: boolean;
  onDemoAuthenticated: (v: boolean) => void;
}) => {
  const { errorNotify, notify } = useToast();
  const navigate = useNavigate();

  const { t } = useTranslation(['common', 'form', 'auth', 'message']);

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      organizationName: '',
      organizationUrlCode: '',
      isAgree: undefined
    },
    mode: 'onChange'
  });

  const {
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<AccessDemoForm> = async values => {
    try {
      const response = await organizationDemoCreator({
        name: values.organizationName,
        urlCode: values.organizationUrlCode
      });

      if (response?.organization) {
        notify({
          message: t('message:collection-action-success', {
            collection: t('organization'),
            action: t('created')
          })
        });
        clearDemoTokenStorage();
        navigate(PAGE_PATH_ROOT);
      }
    } catch (error) {
      if ((error as AxiosError).status === 401) {
        onDemoAuthenticated(false);
        clearDemoTokenStorage();
      }
      errorNotify(error);
    }
  };

  return (
    <FormProvider {...form}>
      <Form onSubmit={form.handleSubmit(onSubmit)} className="mt-5">
        <Form.Field
          control={form.control}
          name="organizationName"
          render={({ field }) => (
            <Form.Item>
              <Form.Label required>{t('auth:organization-name')}</Form.Label>
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
          name="isAgree"
          render={({ field }) => (
            <Form.Item>
              <Form.Control>
                <Checkbox
                  title={t('auth:demo-agree-terms')}
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
          disabled={!isValid || !isDirty || !isDemoSiteEnabled}
          loading={isSubmitting}
          className="mt-8 w-full"
        >
          {t('common:submit')}
        </Button>
      </Form>
    </FormProvider>
  );
};

export default DemoForm;
