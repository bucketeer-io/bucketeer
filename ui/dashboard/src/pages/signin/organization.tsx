import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { switchOrganization } from '@api/auth';
import { yupResolver } from '@hookform/resolvers/yup';
import { useAuth } from 'auth';
import { requiredMessage } from 'constants/message';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useTranslation } from 'i18n';
import { jwtDecode } from 'jwt-decode';
import { setOrgIdStorage } from 'storage/organization';
import { getTokenStorage, setTokenStorage } from 'storage/token';
import * as yup from 'yup';
import { DecodedToken } from '@types';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import AuthWrapper from './elements/auth-wrapper';

const formSchema = yup.object().shape({
  organization: yup.string().required(requiredMessage)
});

const SelectOrganization = () => {
  const { t } = useTranslation(['auth', 'common']);
  const navigate = useNavigate();
  const { myOrganizations, onMeFetcher } = useAuth();

  const form = useForm({
    resolver: yupResolver(formSchema),
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

        const fetchUserData = () => {
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
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      label={
                        myOrganizations.find(
                          org => org.id === form.getValues('organization')
                        )?.name || t(`organization-placeholder`)
                      }
                      isExpand
                    />
                    <DropdownMenuContent className="w-[442px]">
                      {myOrganizations.map((org, index) => (
                        <DropdownMenuItem
                          {...field}
                          key={index}
                          label={org.name}
                          value={org.id}
                          onSelectOption={value => {
                            form.clearErrors();
                            form.setValue('organization', value as string);
                          }}
                        />
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
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
