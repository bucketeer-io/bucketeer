import { FunctionComponent, useCallback, useState } from 'react';
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
import { invalidateTeams, useQueryTeams } from '@queries/teams';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, getEditorEnvironments, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import useOptions from 'hooks/use-options';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { Language, useTranslation } from 'i18n';
import uniqBy from 'lodash/uniqBy';
import * as yup from 'yup';
import { EnvironmentRoleType, OrganizationRole } from '@types';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { IconEnglishFlag, IconInfo, IconJapanFlag } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownOption
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { Tooltip } from 'components/tooltip';
import SelectMenu from 'elements/select-menu';
import EnvironmentRoles from './environment-roles';

interface AddMemberModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const defaultEnvironmentRole: EnvironmentRoleItem = {
  environmentId: '',
  role: 'Environment_UNASSIGNED'
};

interface LanguageItem {
  readonly label: string;
  readonly value: Language;
  readonly icon: FunctionComponent;
}
export const languageList: LanguageItem[] = [
  { label: '日本語', value: Language.JAPANESE, icon: IconJapanFlag },
  { label: 'English', value: Language.ENGLISH, icon: IconEnglishFlag }
];

export interface AddMemberForm {
  email: string;
  memberRole: OrganizationRole;
  environmentRoles: EnvironmentRoleItem[];
  teams: string[];
}

export const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    email: yup.string().email().required(requiredMessage),
    memberRole: yup.mixed<OrganizationRole>().required(requiredMessage),
    environmentRoles: yup
      .array()
      .when('memberRole', {
        is: (role: OrganizationRole) => role === 'Organization_MEMBER',
        then: schema => schema.required(requiredMessage)
      })
      .of(
        yup.object().shape({
          environmentId: yup.string().when('memberRole', {
            is: (role: OrganizationRole) => role === 'Organization_MEMBER',
            then: schema => schema.required(requiredMessage)
          }),
          role: yup
            .mixed<EnvironmentRoleType>()
            .when('memberRole', {
              is: (role: OrganizationRole) => role === 'Organization_MEMBER',
              then: schema => schema.required(requiredMessage)
            })
            .test('isUnassigned', (value, context) => {
              const isMemberRole =
                context?.from &&
                context?.from[1]?.value?.memberRole === 'Organization_MEMBER';
              if (value === 'Environment_UNASSIGNED' && isMemberRole)
                return context.createError({
                  message: requiredMessage,
                  path: context.path
                });
              return true;
            })
        })
      ),
    teams: yup.array().of(yup.string())
  });

const AddMemberModal = ({ isOpen, onClose }: AddMemberModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify, errorNotify } = useToast();
  const { organizationRoles } = useOptions();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { editorEnvironments } = getEditorEnvironments(consoleAccount!);
  const [teamOptions, setTeamOptions] = useState<DropdownOption[]>([]);

  const { data: teamCollection, isLoading: isLoadingTeams } = useQueryTeams({
    params: {
      cursor: String(0),
      organizationId: currentEnvironment.organizationId
    }
  });

  const form = useForm<AddMemberForm>({
    resolver: yupResolver(useFormSchema(formSchema)) as Resolver<AddMemberForm>,
    defaultValues: {
      email: '',
      memberRole: undefined,
      environmentRoles: [defaultEnvironmentRole],
      teams: []
    }
  });

  const {
    watch,
    formState: { isValid, isDirty, isSubmitting }
  } = form;
  const memberEnvironments = watch('environmentRoles');
  const roleWatch = watch('memberRole');
  const isAdminRole = roleWatch === 'Organization_ADMIN';
  const teams = teamCollection?.teams || [];

  const uniqueNameTeams = uniqBy(teams || [], 'name')?.map(item => ({
    label: item.name,
    value: item.name
  }));

  const teamDropdownOptions = uniqBy(
    [...teamOptions, ...uniqueNameTeams],
    'label'
  );

  const { formattedEnvironments } = onFormatEnvironments(editorEnvironments);

  const checkSubmitBtnDisabled = useCallback(() => {
    if (isValid && isAdminRole) return false;
    const checkEnvironments = memberEnvironments.every(
      item => item.environmentId && item.role !== 'Environment_UNASSIGNED'
    );
    if (!checkEnvironments || !isValid) return true;

    return false;
  }, [isValid, memberEnvironments, isAdminRole]);

  const onSubmit: SubmitHandler<AddMemberForm> = async values => {
    return accountCreator({
      organizationId: currentEnvironment.organizationId,
      email: values.email,
      organizationRole: values.memberRole,
      ...(isAdminRole
        ? {}
        : {
            environmentRoles: values.environmentRoles.map(item => ({
              ...item,
              environmentId: checkEnvironmentEmptyId(item.environmentId)
            }))
          }),
      teams: values.teams ?? []
    })
      .then(() => {
        notify({
          message: t('message:collection-action-success', {
            collection: t('member'),
            action: t('created')
          })
        });
        invalidateAccounts(queryClient);
        invalidateTeams(queryClient);
        onClose();
      })
      .catch(error => errorNotify(error));
  };

  useUnsavedLeavePage({ isShow: isDirty && !isSubmitting });
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
                      autoComplete="email"
                      {...field}
                      name="member-email"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="memberRole"
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
                            description={item.description}
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
              name={`teams`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label className="relative w-fit">
                    {t('teams')}
                    <Tooltip
                      align="start"
                      alignOffset={-30}
                      trigger={
                        <div className="flex-center absolute top-0 -right-6">
                          <Icon icon={IconInfo} size={'sm'} color="gray-600" />
                        </div>
                      }
                      content={t('form:teams-tooltip')}
                      className="!z-[100] max-w-[400px]"
                    />
                  </Form.Label>
                  <Form.Control>
                    <SelectMenu
                      options={teamDropdownOptions}
                      fieldValues={field.value}
                      disabled={isLoadingTeams}
                      inputPlaceholderKey="search-teams-placeholder"
                      dropdownPlaceholderKey="select-teams"
                      onChange={field.onChange}
                      onChangeOptions={setTeamOptions}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            {!isAdminRole && !!roleWatch && (
              <>
                <Divider className="mt-1 mb-3" />
                <EnvironmentRoles environments={formattedEnvironments} />
              </>
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
