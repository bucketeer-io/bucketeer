import { useCallback } from 'react';
import {
  FormProvider,
  Resolver,
  SubmitHandler,
  useForm
} from 'react-hook-form';
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
import uniqBy from 'lodash/uniqBy';
import * as yup from 'yup';
import { EnvironmentRoleType, OrganizationRole } from '@types';
import { IconInfo } from '@icons';
import { useFetchTags } from 'pages/members/collection-loader';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { CreatableSelect } from 'components/creatable-select';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { Tooltip } from 'components/tooltip';
import EnvironmentRoles from './environment-roles';

interface AddMemberModalProps {
  isOpen: boolean;
  onClose: () => void;
}

interface OrganizationRoleOption {
  value: OrganizationRole;
  label: string;
}

export const defaultEnvironmentRole: EnvironmentRoleItem = {
  environmentId: '',
  role: 'Environment_UNASSIGNED'
};

export const organizationRoles: OrganizationRoleOption[] = [
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
  role: OrganizationRole;
  environmentRoles: EnvironmentRoleItem[];
  tags: string[];
}

export const formSchema = yup.object().shape({
  email: yup.string().email().required(),
  role: yup.mixed<OrganizationRole>().required(),
  environmentRoles: yup
    .array()
    .required()
    .of(
      yup.object().shape({
        environmentId: yup.string().required(),
        role: yup.mixed<EnvironmentRoleType>().required()
      })
    ),
  tags: yup.array().of(yup.string())
});

const AddMemberModal = ({ isOpen, onClose }: AddMemberModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: tagCollection, isLoading: isLoadingTags } = useFetchTags({
    organizationId: currentEnvironment.organizationId
  });

  const form = useForm<AddMemberForm>({
    resolver: yupResolver(formSchema) as Resolver<AddMemberForm>,
    defaultValues: {
      email: '',
      role: undefined,
      environmentRoles: [defaultEnvironmentRole],
      tags: []
    }
  });

  const {
    watch,
    formState: { dirtyFields, isValid, isSubmitting }
  } = form;
  const memberEnvironments = watch('environmentRoles');

  const tagOptions = uniqBy(tagCollection?.tags || [], 'name')?.filter(tag =>
    memberEnvironments.find(env => env.environmentId === tag.environmentId)
  );

  const { data: collection } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });
  const environments = collection?.environments || [];

  const checkSubmitBtnDisabled = useCallback(() => {
    const checkEnvironments = memberEnvironments.every(
      item => item.environmentId && item.role !== 'Environment_UNASSIGNED'
    );
    if (!checkEnvironments || !isValid) {
      return true;
    }
    return false;
  }, [dirtyFields, isValid, memberEnvironments]);

  const onSubmit: SubmitHandler<AddMemberForm> = values => {
    return accountCreator({
      organizationId: currentEnvironment.organizationId,
      command: {
        email: values.email,
        organizationRole: values.role,
        environmentRoles: values.environmentRoles,
        tags: values.tags ?? []
      }
    }).then(() => {
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{values.email}</b> {` has been successfully created!`}
          </span>
        )
      });
      invalidateAccounts(queryClient);
      onClose();
    });
  };

  return (
    <SlideModal title={t('invite-member')} isOpen={isOpen} onClose={onClose}>
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

            <Form.Field
              control={form.control}
              name={`tags`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label className="relative w-fit">
                    {t('tags')}
                    <Tooltip
                      align="start"
                      alignOffset={-130}
                      trigger={
                        <div className="flex-center absolute top-0 -right-6">
                          <Icon icon={IconInfo} size={'sm'} color="gray-600" />
                        </div>
                      }
                      content={t('form:member-tags-tooltip')}
                      className="!z-[100] max-w-[400px]"
                    />
                  </Form.Label>
                  <Form.Control>
                    <CreatableSelect
                      disabled={isLoadingTags}
                      loading={isLoadingTags}
                      placeholder={t(`form:placeholder-tags`)}
                      options={tagOptions?.map(tag => ({
                        label: tag.name,
                        value: tag.id
                      }))}
                      onChange={value =>
                        field.onChange(value.map(tag => tag.value))
                      }
                    />
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
                    <EnvironmentRoles
                      environments={environments}
                      memberEnvironments={memberEnvironments}
                      onChangeEnvironments={values => {
                        field.onChange(values);
                      }}
                    />
                  </Form.Control>
                  <Form.Message />
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
                    {t(`invite-member`)}
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
