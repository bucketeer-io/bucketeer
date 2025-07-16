import { useCallback, useEffect, useState } from 'react';
import {
  FormProvider,
  Resolver,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import { EnvironmentRoleItem } from '@api/account/account-creator';
import { accountUpdater } from '@api/account/account-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateAccounts } from '@queries/accounts';
import { invalidateTeams, useQueryTeams } from '@queries/teams';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, getEditorEnvironments, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import useOptions from 'hooks/use-options';
import { Language, setLanguage, useTranslation } from 'i18n';
import uniqBy from 'lodash/uniqBy';
import * as yup from 'yup';
import {
  Account,
  EnvironmentRoleType,
  OrganizationRole,
  TeamChange
} from '@types';
import {
  checkEnvironmentEmptyId,
  onChangeFontWithLocalized,
  onFormatEnvironments
} from 'utils/function';
import { IconInfo } from '@icons';
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
import { defaultEnvironmentRole, languageList } from '../add-member-modal';
import EnvironmentRoles from '../add-member-modal/environment-roles';

interface EditMemberModalProps {
  isOpen: boolean;
  onClose: () => void;
  member: Account;
}

export interface EditMemberForm {
  firstName: string;
  lastName: string;
  language: string;
  memberRole: string;
  environmentRoles: EnvironmentRoleItem[];
  teams: string[];
}

export const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    firstName: yup.string().required(requiredMessage),
    lastName: yup.string().required(requiredMessage),
    language: yup.string().required(requiredMessage),
    memberRole: yup.string().required(requiredMessage),
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

const EditMemberModal = ({ isOpen, onClose, member }: EditMemberModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify } = useToast();
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

  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(editorEnvironments);

  const teams = teamCollection?.teams || [];
  const form = useForm<EditMemberForm>({
    resolver: yupResolver(
      useFormSchema(formSchema)
    ) as Resolver<EditMemberForm>,
    values: {
      firstName: member.firstName,
      lastName: member.lastName,
      language: member.language,
      memberRole: member.organizationRole,
      environmentRoles: member.environmentRoles?.length
        ? member.environmentRoles.map(item => ({
            ...item,
            environmentId: item.environmentId || emptyEnvironmentId
          }))
        : [defaultEnvironmentRole],
      teams: member.teams
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isSubmitting },
    watch
  } = form;
  const roleWatch = watch('memberRole');
  const isAdminRole = roleWatch === 'Organization_ADMIN';

  const handleCheckTags = useCallback(
    (teamValues: string[]) => {
      const teamChanges: TeamChange[] = [];
      const { teams } = member;
      teams?.forEach(item => {
        if (!teamValues.find(tag => tag === item)) {
          teamChanges.push({
            changeType: 'DELETE',
            team: item
          });
        }
      });
      teamValues.forEach(item => {
        const currentTeam = teams.find(team => team === item);
        if (!currentTeam) {
          teamChanges.push({
            changeType: 'CREATE',
            team: item
          });
        }
      });

      return teamChanges;
    },
    [member]
  );

  const onSubmit: SubmitHandler<EditMemberForm> = async values => {
    return accountUpdater({
      organizationId: currentEnvironment.organizationId,
      email: member.email,
      organizationRole: {
        role: values.memberRole as OrganizationRole
      },
      ...(isAdminRole
        ? {}
        : {
            environmentRoles: values.environmentRoles.map(item => ({
              ...item,
              environmentId: checkEnvironmentEmptyId(item.environmentId)
            }))
          }),
      firstName: values.firstName,
      lastName: values.lastName,
      language: values.language,
      teamChanges: handleCheckTags(values.teams)
    }).then(() => {
      const { email } = consoleAccount!;
      if (email === member.email && values.language !== member.language) {
        setLanguage(values.language as Language);
        onChangeFontWithLocalized(values.language === Language.JAPANESE);
      }
      notify({
        message: t('message:collection-action-success', {
          collection: t('member'),
          action: t('updated')
        })
      });
      invalidateAccounts(queryClient);
      invalidateTeams(queryClient);
      onClose();
    });
  };

  useEffect(() => {
    if (teams.length > 0) {
      const uniqueTeams = uniqBy(teams || [], 'name')?.map(item => ({
        label: item.name,
        value: item.name
      }));
      setTeamOptions(uniqueTeams);
    }
  }, [teams]);

  return (
    <SlideModal title={t('update-member')} isOpen={isOpen} onClose={onClose}>
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
                      autoComplete="given-name"
                      {...field}
                      name="first-name"
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
                      placeholder={`${t('form:enter-last-name')}`}
                      autoComplete="family-name"
                      {...field}
                      name="last-name"
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
                  value={member.email}
                  placeholder={t('form:placeholder-email')}
                  autoComplete="email"
                  name="member-email"
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
            <Form.Field
              control={form.control}
              name="language"
              render={({ field }) => {
                const currentItem = languageList.find(
                  item => item.value === field.value
                );

                return (
                  <Form.Item>
                    <Form.Label required>{t('language')}</Form.Label>
                    <Form.Control className="w-full">
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          placeholder={t('form:select-language')}
                          trigger={
                            <div className="flex items-center gap-x-2">
                              {currentItem?.icon && (
                                <div className="flex-center size-fit mt-0.5">
                                  <Icon icon={currentItem?.icon} size={'sm'} />
                                </div>
                              )}
                              {currentItem?.label}
                            </div>
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
                              iconElement={
                                <div className="flex-center size-fit mt-0.5">
                                  <Icon icon={item.icon} size={'sm'} />
                                </div>
                              }
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
                );
              }}
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
                      options={teamOptions}
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
                    disabled={!isValid}
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
    </SlideModal>
  );
};

export default EditMemberModal;
