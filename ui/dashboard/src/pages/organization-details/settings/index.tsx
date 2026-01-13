import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { organizationUpdater } from '@api/organization';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateAccounts, useQueryAccounts } from '@queries/accounts';
import { invalidateOrganizationDetails } from '@queries/organization-details';
import { useQueryClient } from '@tanstack/react-query';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Organization } from '@types';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import TextArea from 'components/textarea';

const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    name: yup.string().required(requiredMessage),
    urlCode: yup.string().required(requiredMessage),
    description: yup.string(),
    ownerEmail: yup.string().required(requiredMessage)
  });

export interface OrganizationSettingsForm {
  name: string;
  urlCode: string;
  description?: string;
  ownerEmail: string;
}

const OrganizationSettings = ({
  organization
}: {
  organization: Organization;
}) => {
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const params = useParams();
  const orgDetailsId = params.organizationId!;
  const { data: accounts } = useQueryAccounts({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      organizationId: orgDetailsId
    }
  });

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      name: organization.name,
      urlCode: organization.urlCode,
      description: organization.description,
      ownerEmail: organization.ownerEmail
    }
  });

  const onSubmit: SubmitHandler<OrganizationSettingsForm> = async values => {
    try {
      const resp = await organizationUpdater({
        id: orgDetailsId,
        name: values.name,
        ownerEmail: values.ownerEmail,
        description: values.description
      });
      if (resp) {
        invalidateOrganizationDetails(queryClient, { id: orgDetailsId });
        invalidateAccounts(queryClient);
        notify({
          message: t('message:collection-action-success', {
            collection: t('organization'),
            action: t('updated')
          })
        });
      }
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <div className="w-full px-3 sm:px-6">
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
                    <Dropdown
                      options={
                        accounts?.accounts.map(item => ({
                          label: item.email,
                          value: item.email
                        })) || []
                      }
                      value={field.value}
                      onChange={field.onChange}
                      placeholder={t('form:owner-email')}
                      className="w-full"
                      contentClassName="w-[400px]"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Button
              loading={form.formState.isSubmitting}
              disabled={!form.formState.isDirty}
              type="submit"
              className="w-fit mt-6"
            >
              {t(`save`)}
            </Button>
          </Form>
        </FormProvider>
      </div>
    </div>
  );
};

export default OrganizationSettings;
