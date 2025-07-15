import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { onGenerateSlug } from 'utils/converts';
import Button from 'components/button';
import Checkbox from 'components/checkbox';
import Form from 'components/form';
import Input from 'components/input';

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
    isAgree: yup
      .boolean()
      .isTrue(translation('message:required-agreement-terms'))
      .required(requiredMessage)
  });

const CreateDemoOrganizationForm = ({
  isDemoSiteEnabled
}: {
  isDemoSiteEnabled?: boolean;
}) => {
  const { t } = useTranslation(['auth', 'common', 'form', 'message']);

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      organizationName: '',
      organizationUrlCode: '',
      email: '',
      isAgree: undefined
    },
    mode: 'onChange'
  });

  const {
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<AccessDemoForm> = async values => {
    try {
      console.log(values);
    } catch (error) {
      if (error) {
        console.log(error);
      }
    }
  };

  return (
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
              <Form.Label required>{t('owner-email')}</Form.Label>
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

export default CreateDemoOrganizationForm;
