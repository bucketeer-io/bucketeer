import { useCallback, useEffect } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { EnvironmentRoleItem } from '@api/account/account-creator';
import { accountUpdater } from '@api/account/account-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import {
  invalidateAccountOrganizationDetails,
  useQueryAccountOrganizationDetails
} from '@queries/account-by-organization-details';
import { invalidateAccounts } from '@queries/accounts';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { EnvironmentRoleType, OrganizationRole } from '@types';
import { joinName } from 'utils/name';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import FormLoading from 'elements/form-loading';
import {
  defaultEnvironmentRole,
  languageList,
  organizationRoles
} from '../add-member-modal';
import EnvironmentRoles from '../add-member-modal/environment-roles';

interface EditMemberModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface EditMemberForm {
  firstName: string;
  lastName: string;
  language: string;
  role: string;
  environmentRoles: EnvironmentRoleItem[];
}

export const formSchema = yup.object().shape({
  firstName: yup.string().required(),
  lastName: yup.string().required(),
  language: yup.string().required(),
  role: yup.string().required(),
  environmentRoles: yup
    .array()
    .required()
    .of(
      yup.object().shape({
        environmentId: yup.string().required(),
        role: yup.mixed<EnvironmentRoleType>().required()
      })
    )
});

const EditMemberModal = ({ isOpen, onClose }: EditMemberModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { id: memberEmail, errorToast } = useActionWithURL({});

  const {
    data: memberCollection,
    isLoading,
    error
  } = useQueryAccountOrganizationDetails({
    params: {
      email: memberEmail as string,
      organizationId: currentEnvironment.organizationId
    }
  });

  const member = memberCollection?.account;

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      firstName: '',
      lastName: '',
      language: '',
      role: '',
      environmentRoles: []
    }
  });

  const {
    formState: { isDirty, isSubmitting, isValid },
    watch
  } = form;

  const { data: collection } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });
  const environments = collection?.environments || [];
  const memberEnvironments = watch('environmentRoles');

  const checkSubmitBtnDisabled = useCallback(() => {
    const checkEnvironments = memberEnvironments.every(
      item => item.environmentId && item.role !== 'Environment_UNASSIGNED'
    );
    if (!checkEnvironments || !isValid || !isDirty) {
      return true;
    }
    return false;
  }, [isDirty, isValid, memberEnvironments]);

  const onSubmit: SubmitHandler<EditMemberForm> = async values => {
    try {
      const resp = await accountUpdater({
        organizationId: currentEnvironment.organizationId,
        email: memberEmail || '',
        changeOrganizationRoleCommand: {
          role: values.role as OrganizationRole
        },
        changeEnvironmentRolesCommand: {
          roles: values.environmentRoles,
          writeType: 'WriteType_OVERRIDE'
        },
        changeFirstNameCommand: {
          firstName: values.firstName
        },
        changeLastNameCommand: {
          lastName: values.lastName
        },
        changeLanguageCommand: {
          language: values.language
        }
      });
      if (resp) {
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: (
            <span>
              <b>{joinName(values.firstName, values.lastName)}</b>
              {` has been successfully updated!`}
            </span>
          )
        });
        invalidateAccounts(queryClient);
        invalidateAccountOrganizationDetails(queryClient, {
          email: memberEmail!,
          organizationId: currentEnvironment.organizationId
        });
        onClose();
      }
    } catch (error) {
      errorToast(error as Error);
    }
  };

  useEffect(() => {
    if (member) {
      form.reset({
        firstName: member?.firstName,
        lastName: member?.lastName,
        language: member?.language,
        role: member?.organizationRole,
        environmentRoles: member?.environmentRoles?.length
          ? member?.environmentRoles
          : [defaultEnvironmentRole]
      });
    }
  }, [member]);

  useEffect(() => {
    if (error) errorToast(error);
  }, [error]);

  return (
    <SlideModal title={t('update-member')} isOpen={isOpen} onClose={onClose}>
      {isLoading ? (
        <FormLoading />
      ) : (
        <div className="w-full p-5 pb-28">
          <p className="text-gray-800 typo-head-bold-small">
            {t('form:general-info')}
          </p>
          <FormProvider {...form}>
            <Form onSubmit={form.handleSubmit(onSubmit)}>
              <Form.Field
                control={form.control}
                name="firstName"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label required>{t('first-name')}</Form.Label>
                    <Form.Control>
                      <Input
                        placeholder={`${t('form:enter-first-name')}`}
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
                    <Form.Label required>{t('last-name')}</Form.Label>
                    <Form.Control>
                      <Input
                        placeholder={`${t('form:enter-first-name')}`}
                        {...field}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Item>
                <Form.Label required>{t('email')}</Form.Label>
                <Form.Control>
                  <Input
                    disabled
                    value={member?.email || ''}
                    placeholder={t('form:placeholder-email')}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
              <Form.Field
                control={form.control}
                name="language"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label required>{t('language')}</Form.Label>
                    <Form.Control className="w-full">
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          placeholder={t('form:select-language')}
                          label={
                            languageList.find(
                              item => item.value === field.value
                            )?.label
                          }
                          variant="secondary"
                          className="w-full"
                        />
                        <DropdownMenuContent
                          className="w-[500px]"
                          align="start"
                          {...field}
                        >
                          {languageList.map((item, index) => (
                            <DropdownMenuItem
                              {...field}
                              key={index}
                              value={item.value}
                              label={item.label}
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
              <Form.Field
                control={form.control}
                name="role"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label required>{t('role')}</Form.Label>
                    <Form.Control className="w-full">
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          placeholder={t('form:select-role')}
                          label={
                            organizationRoles.find(
                              item => item.value === field.value
                            )?.label
                          }
                          variant="secondary"
                          className="w-full"
                        />
                        <DropdownMenuContent
                          className="w-[500px]"
                          align="start"
                          {...field}
                        >
                          {organizationRoles.map((item, index) => (
                            <DropdownMenuItem
                              {...field}
                              key={index}
                              value={item.value}
                              label={item.label}
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
              <Divider className="mt-1 mb-3" />
              <Form.Field
                control={form.control}
                name="environmentRoles"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Control>
                      {memberEnvironments && (
                        <EnvironmentRoles
                          environments={environments}
                          memberEnvironments={memberEnvironments}
                          onChangeEnvironments={values =>
                            field.onChange(values)
                          }
                        />
                      )}
                    </Form.Control>
                  </Form.Item>
                )}
              />
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
                      disabled={checkSubmitBtnDisabled()}
                      loading={isSubmitting}
                    >
                      {t(`save`)}
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

export default EditMemberModal;
