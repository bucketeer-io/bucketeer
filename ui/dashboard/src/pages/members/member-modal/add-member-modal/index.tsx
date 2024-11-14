import { useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import {
  accountCreator,
  EnvironmentRoleItem
} from '@api/account/account-creator';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateAccounts } from '@queries/accounts';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { Language, useTranslation } from 'i18n';
import * as yup from 'yup';
import { OrganizationRole } from '@types';
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
import EnvironmentRoles from './environment-roles';

interface AddMemberModalProps {
  isOpen: boolean;
  onClose: () => void;
}

interface organizationRoleOption {
  value: OrganizationRole;
  label: string;
}

export const organizationRoles: organizationRoleOption[] = [
  {
    value: 'Organization_MEMBER',
    label: 'Member'
  },
  {
    value: 'Organization_ADMIN',
    label: 'Admin'
  }
];

interface LanguageItem {
  readonly label: string;
  readonly value: Language;
}

export const languageList: LanguageItem[] = [
  { label: '日本語', value: Language.JAPANESE },
  { label: 'English', value: Language.ENGLISH }
];

export interface AddMemberForm {
  email: string;
  role: string;
}

export const formSchema = yup.object().shape({
  email: yup.string().email().required(),
  role: yup.string().required()
});

const AddMemberModal = ({ isOpen, onClose }: AddMemberModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      email: '',
      role: undefined
    }
  });

  const [memberEnvironments, setMemberEnvironments] = useState<
    EnvironmentRoleItem[]
  >([]);

  const { data: collection } = useFetchEnvironments();
  const environments = collection?.environments || [];

  const isInvalidEnvironments = () => {
    const invalidEnv = memberEnvironments.find(
      item => !item.environmentId || item.role === 'Environment_UNASSIGNED'
    );
    return memberEnvironments.length > 0 && !!invalidEnv;
  };

  const onSubmit: SubmitHandler<AddMemberForm> = values => {
    return accountCreator({
      organizationId: currentEnvironment.organizationId,
      command: {
        email: values.email,
        organizationRole: values.role as OrganizationRole,
        environmentRoles:
          memberEnvironments.length > 0 ? memberEnvironments : undefined
      }
    }).then(() => {
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{values.email}</b> {` has been successfully created!`}{' '}
          </span>
        )
      });
      invalidateAccounts(queryClient);
      onClose();
    });
  };

  return (
    <SlideModal title={t('new-member')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full p-5 pb-28">
        <p className="text-gray-800 typo-head-bold-small">
          {t('form:general-info')}
        </p>
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <Form.Field
              control={form.control}
              name="email"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('email')}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t('form:placeholder-email')}
                      {...field}
                    />
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
            <EnvironmentRoles
              environments={environments}
              memberEnvironments={memberEnvironments}
              setMemberEnvironments={setMemberEnvironments}
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
                    disabled={
                      !form.formState.isValid || isInvalidEnvironments()
                    }
                    loading={form.formState.isSubmitting}
                  >
                    {t(`create-member`)}
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

export default AddMemberModal;
