import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { UserInfoForm } from '@types';
import Button from 'components/button';
import Form from 'components/form';
import Input from 'components/input';
import AuthWrapper from './elements/auth-wrapper';

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    firstName: yup
      .string()
      .required(requiredMessage)
      .min(
        2,
        translation('message:validation.name-at-least-characters', {
          count: 2,
          name: translation('common:first-name').toLowerCase()
        })
      ),
    lastName: yup
      .string()
      .required(requiredMessage)
      .min(
        2,
        translation('message:validation.name-at-least-characters', {
          count: 2,
          name: translation('common:first-name').toLowerCase()
        })
      ),
    language: yup.string().required(requiredMessage)
  });

const UserInformation = () => {
  const { t } = useTranslation(['auth']);
  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      firstName: '',
      lastName: '',
      language: ''
    }
  });

  const onSubmit: SubmitHandler<UserInfoForm> = values => {
    console.log(values);
  };

  return (
    <AuthWrapper>
      <div className="grid gap-10">
        <div>
          <h1 className="text-gray-900 typo-head-bold-huge mb-4">
            {t(`enter-information.title`)}
          </h1>
          <p className="text-gray-600 typo-para-medium">
            {t(`enter-information.description`)}
          </p>
        </div>

        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)} className="mt-8">
            <Form.Field
              control={form.control}
              name="firstName"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t(`first-name`)}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t(`first-name-placeholder`)}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="lastName"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t(`last-name`)}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t(`last-name-placeholder`)}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            {/* <Select
              placeholder="Select your Language"
              options={[]}
              label="Language"
              className="mt-4"
              defaultValue={0}
              required
            /> */}
            <Button type="submit" className="mt-8 w-full">
              {`Sign In`}
            </Button>
          </Form>
        </FormProvider>
      </div>
    </AuthWrapper>
  );
};

export default UserInformation;
