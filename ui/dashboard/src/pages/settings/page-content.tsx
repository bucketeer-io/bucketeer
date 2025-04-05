import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { organizationUpdater } from '@api/organization';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryAccounts } from '@queries/accounts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast } from 'hooks';
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
import PageLayout from 'elements/page-layout';

const formSchema = yup.object().shape({
  name: yup.string().required(),
  urlCode: yup.string().required(),
  description: yup.string(),
  ownerEmail: yup.string().required()
});

export interface PageContentForm {
  name: string;
  urlCode: string;
  description?: string;
  ownerEmail: string;
}

const PageContent = ({ organization }: { organization: Organization }) => {
  const { notify } = useToast();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { t } = useTranslation(['common', 'form']);
  const { data: accounts, isLoading: isLoadingAccounts } = useQueryAccounts({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      organizationId: currentEnvironment.organizationId
    }
  });

  const form = useForm({
    resolver: yupResolver(formSchema),
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
          toastType: 'toast',
          messageType: 'success',
          message: (
            <span>
              <b>{values.name}</b> {`has been successfully updated!`}
            </span>
          )
        });
      }
    } catch (error) {
      notify({
        toastType: 'toast',
        messageType: 'error',
        message: (error as Error)?.message || 'Something went wrong.'
      });
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
                      placeholder={`${t('form:placeholder-name')}`}
                      {...field}
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
                        disabled={isLoadingAccounts}
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
    </PageLayout.Content>
  );
};

export default PageContent;
