import { useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { organizationUpdater } from '@api/organization';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryAccounts } from '@queries/accounts';
import { invalidateOrganizationDetails } from '@queries/organization-details';
import { invalidateOrganizations } from '@queries/organizations';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth, useAuthAccess } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Organization } from '@types';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import TextArea from 'components/textarea';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import PageLayout from 'elements/page-layout';

const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    name: yup.string().required(requiredMessage),
    urlCode: yup.string().required(requiredMessage),
    description: yup.string(),
    ownerEmail: yup.string().required(requiredMessage)
  });

export interface PageContentForm {
  name: string;
  urlCode: string;
  description?: string;
  ownerEmail: string;
}

const PageContent = ({ organization }: { organization: Organization }) => {
  const { notify, errorNotify } = useToast();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { envEditable, isOrganizationAdmin } = useAuthAccess();
  const queryClient = useQueryClient();

  const { t } = useTranslation(['common', 'form', 'message']);
  const { data: accounts, isLoading: isLoadingAccounts } = useQueryAccounts({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      organizationId: currentEnvironment.organizationId
    }
  });
  const disabled = useMemo(
    () => !envEditable || !isOrganizationAdmin,
    [envEditable, isOrganizationAdmin]
  );

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      name: organization.name,
      urlCode: organization.urlCode,
      description: organization.description,
      ownerEmail: organization.ownerEmail
    }
  });

  const onSubmit: SubmitHandler<PageContentForm> = async values => {
    try {
      const resp = await organizationUpdater({
        id: currentEnvironment.organizationId,
        description: values.description,
        name: values.name,
        ownerEmail: values.ownerEmail
      });
      if (resp) {
        notify({
          message: t('message:collection-action-success', {
            collection: t('organization'),
            action: t('updated')
          })
        });
        invalidateOrganizationDetails(queryClient, {
          id: currentEnvironment.organizationId
        });
        invalidateOrganizations(queryClient);
      }
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <PageLayout.Content className="p-6">
      <div className="p-5 shadow-card rounded-lg bg-white">
        <p className="text-gray-800 typo-head-bold-small">
          {t('form:general-info')}
        </p>
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <Form.Field
              control={form.control}
              name="name"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('name')}</Form.Label>
                  <Form.Control>
                    <Input
                      disabled={disabled}
                      placeholder={`${t('form:placeholder-name')}`}
                      {...field}
                      name="organization-name"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="urlCode"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('form:url-code')}</Form.Label>
                  <Form.Control>
                    <Input
                      disabled
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
              name="description"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label optional>{t('form:description')}</Form.Label>
                  <Form.Control>
                    <TextArea
                      disabled={disabled}
                      placeholder={t('form:placeholder-desc')}
                      rows={4}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="ownerEmail"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('form:owner-email')}</Form.Label>
                  <Form.Control className="w-full">
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t('form:owner-email')}
                        label={
                          accounts?.accounts.find(
                            item => item.email === field.value
                          )?.email
                        }
                        variant="secondary"
                        className="w-full"
                        disabled={isLoadingAccounts || disabled}
                      />
                      <DropdownMenuContent
                        className="w-[400px]"
                        align="start"
                        {...field}
                      >
                        {accounts?.accounts?.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.email}
                            label={item.email}
                            onSelectOption={value => {
                              field.onChange(value);
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
            <DisabledButtonTooltip
              align="start"
              type={!isOrganizationAdmin ? 'admin' : 'editor'}
              hidden={!disabled}
              trigger={
                <Button
                  loading={form.formState.isSubmitting}
                  disabled={!form.formState.isDirty || disabled}
                  type="submit"
                  className="w-fit mt-6"
                >
                  {t(`save`)}
                </Button>
              }
            />
          </Form>
        </FormProvider>
      </div>
    </PageLayout.Content>
  );
};

export default PageContent;
