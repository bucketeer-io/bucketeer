import { useCallback, useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import {
  organizationCreator,
  OrganizationResponse,
  organizationUpdater
} from '@api/organization';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryAccounts } from '@queries/accounts';
import { invalidateOrganizations } from '@queries/organizations';
import { useQueryClient } from '@tanstack/react-query';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Organization } from '@types';
import { onGenerateSlug } from 'utils/converts';
import { cn } from 'utils/style';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownOption
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';
import DropdownList from 'elements/dropdown-list';

interface OrganizationCreateUpdateModalProps {
  isOpen: boolean;
  organization?: Organization;
  onClose: () => void;
}

export interface OrganizationCreateUpdateForm {
  name: string;
  urlCode: string;
  ownerEmail: string;
  isTrial?: boolean;
  description?: string;
}

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    name: yup.string().required(requiredMessage),
    urlCode: yup
      .string()
      .required(requiredMessage)
      .matches(
        /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
        translation('message:validation.id-rule', {
          name: translation('common:url-code')
        })
      ),
    description: yup.string(),
    ownerEmail: yup.string().email().required(requiredMessage),
    isTrial: yup.bool()
  });

const OrganizationCreateUpdateModal = ({
  isOpen,
  onClose,
  organization
}: OrganizationCreateUpdateModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const form = useForm<OrganizationCreateUpdateForm>({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      name: organization?.name || '',
      urlCode: organization?.urlCode || '',
      description: organization?.description,
      ownerEmail: organization?.ownerEmail || '',
      isTrial: organization?.trial || true
    },
    mode: 'onChange'
  });

  const { data: accounts } = useQueryAccounts({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      organizationId: organization?.id
    },
    enabled: !!organization
  });

  const accountOptions = useMemo(
    () =>
      accounts?.accounts.map(item => ({
        label: item.email,
        value: item.email
      })),
    [accounts]
  );

  const onSubmit: SubmitHandler<OrganizationCreateUpdateForm> = useCallback(
    async values => {
      try {
        let resp: OrganizationResponse | null = null;
        if (organization) {
          resp = await organizationUpdater({
            id: organization!.id,
            ...values
          });
        } else {
          resp = await organizationCreator({
            ...values,
            isSystemAdmin: false
          });
        }
        if (resp) {
          notify({
            message: t('message:collection-action-success', {
              collection: t('organization'),
              action: t(organization ? 'updated' : 'created')
            })
          });
          invalidateOrganizations(queryClient);
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [organization]
  );

  return (
    <SlideModal
      title={t(organization ? 'update-org' : 'new-org')}
      isOpen={isOpen}
      onClose={onClose}
    >
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
                      onChange={value => {
                        field.onChange(value);
                        if (!organization) {
                          const isUrlCodeDirty =
                            form.getFieldState('urlCode').isDirty;
                          const urlCode = form.getValues('urlCode');
                          form.setValue(
                            'urlCode',
                            isUrlCodeDirty ? urlCode : onGenerateSlug(value)
                          );
                        }
                      }}
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
                      placeholder={`${t('form:placeholder-code')}`}
                      disabled={!!organization}
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
                    {organization ? (
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          placeholder={t('form:owner-email')}
                          label={
                            accountOptions?.find(
                              item => item.value === field.value
                            )?.label
                          }
                          variant="secondary"
                          className="w-full"
                        />
                        <DropdownMenuContent
                          className={cn('w-[500px]', {
                            'hidden-scroll':
                              accountOptions && accountOptions?.length > 15
                          })}
                          align="start"
                          {...field}
                        >
                          <DropdownList
                            options={accountOptions as DropdownOption[]}
                            onSelectOption={field.onChange}
                          />
                        </DropdownMenuContent>
                      </DropdownMenu>
                    ) : (
                      <Input
                        placeholder={`${t('form:placeholder-email')}`}
                        {...field}
                      />
                    )}
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            {!organization && (
              <Form.Field
                control={form.control}
                name="isTrial"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Control>
                      <Checkbox
                        onCheckedChange={checked => field.onChange(checked)}
                        checked={field.value}
                        title={`${t(`form:trial`)}`}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}

            <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
              <ButtonBar
                primaryButton={
                  <Button variant="secondary" onClick={onClose}>
                    {t(`cancel`)}
                  </Button>
                }
                secondaryButton={
                  <Button
                    type="submit"
                    disabled={!form.formState.isDirty}
                    loading={form.formState.isSubmitting}
                  >
                    {t(organization ? `update-org` : 'create-org')}
                  </Button>
                }
              />
            </div>
          </Form>
        </FormProvider>
      </div>
    </SlideModal>
  );
};

export default OrganizationCreateUpdateModal;
