import { useEffect } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { organizationUpdater } from '@api/organization';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryAccounts } from '@queries/accounts';
import {
  invalidateOrganizationDetails,
  useQueryOrganizationDetails
} from '@queries/organization-details';
import { invalidateOrganizations } from '@queries/organizations';
import { useQueryClient } from '@tanstack/react-query';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import Spinner from 'components/spinner';
import TextArea from 'components/textarea';

interface EditOrganizationModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface EditOrganizationForm {
  name: string;
  ownerEmail: string;
  description?: string;
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  description: yup.string(),
  ownerEmail: yup.string().email().required()
});

const EditOrganizationModal = ({
  isOpen,
  onClose
}: EditOrganizationModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();

  const { id: organizationId, errorToast } = useActionWithURL({
    idKey: '*'
  });

  const { data, isLoading, error } = useQueryOrganizationDetails({
    params: {
      id: organizationId as string
    }
  });

  const organization = data?.organization;

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      description: '',
      ownerEmail: ''
    }
  });

  const { data: accounts } = useQueryAccounts({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      organizationId: organizationId as string
    }
  });

  const onSubmit: SubmitHandler<EditOrganizationForm> = async values => {
    try {
      const resp = await organizationUpdater({
        id: organizationId as string,
        changeDescriptionCommand: {
          description: values.description
        },
        renameCommand: {
          name: values.name
        },
        changeOwnerEmailCommand: {
          ownerEmail: values.ownerEmail
        }
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
        invalidateOrganizations(queryClient);
        invalidateOrganizationDetails(queryClient, {
          id: organizationId as string
        });
        onClose();
      }
    } catch (error) {
      errorToast(error as Error);
    }
  };

  useEffect(() => {
    if (organization)
      form.reset({
        name: organization.name,
        description: organization.description,
        ownerEmail: organization.ownerEmail
      });
  }, [organization, form]);

  useEffect(() => {
    if (error) errorToast(error);
  }, [error]);

  return (
    <SlideModal title={t('update-org')} isOpen={isOpen} onClose={onClose}>
      {isLoading ? (
        <div className="flex-center py-10">
          <Spinner />
        </div>
      ) : (
        <div className="w-full p-5">
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

              <Form.Item>
                <Form.Label required>{t('form:url-code')}</Form.Label>
                <Form.Control>
                  <Input
                    value={organization?.urlCode || ''}
                    placeholder={`${t('form:placeholder-code')}`}
                    disabled
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>

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
                        />
                        <DropdownMenuContent
                          className="w-[500px]"
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

              <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button variant="secondary" onClick={onClose}>
                      {t(`common:cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <Button
                      type="submit"
                      disabled={!form.formState.isDirty}
                      loading={form.formState.isSubmitting}
                    >
                      {t(`update-org`)}
                    </Button>
                  }
                />
              </div>
            </Form>
          </FormProvider>
        </div>
      )}
    </SlideModal>
  );
};

export default EditOrganizationModal;
