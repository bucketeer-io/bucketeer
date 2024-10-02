import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { yupResolver } from '@hookform/resolvers/yup';
import { useAuth } from 'auth';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useTranslation } from 'i18n';
import { setOrgIdStorage } from 'storage/organization';
import * as yup from 'yup';
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
  organization: yup.string().required()
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

  const onSubmit: SubmitHandler<{ organization: string }> = values => {
    if (values.organization) {
      setOrgIdStorage(values.organization);
      return onMeFetcher({ organizationId: values.organization }).then(() => {
        navigate(PAGE_PATH_ROOT);
      });
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
            {t(`continue`)}
          </Button>
        </Form>
      </FormProvider>
    </AuthWrapper>
  );
};

export default SelectOrganization;
