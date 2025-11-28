import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { switchOrganization } from '@api/auth';
import { yupResolver } from '@hookform/resolvers/yup';
import { useAuth } from 'auth';
import { PAGE_PATH_ROOT } from 'constants/routing';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import { jwtDecode } from 'jwt-decode';
import { setOrgIdStorage } from 'storage/organization';
import { getTokenStorage, setTokenStorage } from 'storage/token';
import * as yup from 'yup';
import { DecodedToken } from '@types';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import AuthWrapper from './elements/auth-wrapper';

const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    organization: yup.string().required(requiredMessage)
  });

const SelectOrganization = () => {
  const { t } = useTranslation(['auth', 'common']);
  const navigate = useNavigate();
  const { myOrganizations, onMeFetcher } = useAuth();

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      organization: ''
    }
  });

  const onSubmit: SubmitHandler<{ organization: string }> = async values => {
    try {
      const organizationId = values.organization;
      const token = getTokenStorage();

      if (organizationId && token) {
        setOrgIdStorage(organizationId);
        const parsedToken: DecodedToken = jwtDecode(token?.accessToken);

        const fetchUserData = async () => {
          return onMeFetcher({ organizationId }).then(() => {
            navigate(PAGE_PATH_ROOT);
          });
        };

        if (parsedToken.organization_id === organizationId) {
          await fetchUserData();
        } else {
          const resp = await switchOrganization({
            organizationId,
            accessToken: token.accessToken
          });
          setTokenStorage(resp.token);
          await fetchUserData();
        }
      }
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <AuthWrapper>
      <h1 className="text-gray-900 typo-head-bold-huge">
        {t(`select-organization.title`)}
      </h1>
      <p className="text-gray-600 typo-para-medium mt-4">
        {t(`select-organization.description`)}
      </p>
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)} className="mt-10">
          <Form.Field
            control={form.control}
            name="organization"
            render={({ field }) => (
              <Form.Item>
                <Form.Label required>{t(`organization`)}</Form.Label>
                <Form.Control>
                  <Dropdown
                    options={myOrganizations.map(org => ({
                      label: org.name,
                      value: org.id
                    }))}
                    value={field.value}
                    onChange={value => {
                      form.clearErrors();
                      field.onChange(value);
                    }}
                    placeholder={t(`organization-placeholder`)}
                    contentClassName="w-[442px]"
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Button
            type="submit"
            loading={form.formState.isSubmitting}
            className="w-full mt-10"
          >
            {t(`common:continue`)}
          </Button>
        </Form>
      </FormProvider>
    </AuthWrapper>
  );
};

export default SelectOrganization;
