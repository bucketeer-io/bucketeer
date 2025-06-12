import { useEffect, useState } from 'react';
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
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import uniqBy from 'lodash/uniqBy';
import * as yup from 'yup';
import { Account, EnvironmentRoleType, OrganizationRole } from '@types';
import { joinName } from 'utils/name';
import { IconInfo } from '@icons';
import { useFetchTags } from 'pages/members/collection-loader';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
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
import TagsSelectMenu from 'elements/tags-select-menu';
import { languageList, organizationRoles } from '../add-member-modal';
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
  role: string;
  environmentRoles: EnvironmentRoleItem[];
  tags: string[];
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
        environmentId: yup
          .string()
          .required(`Environment is a required field.`),
        role: yup
          .mixed<EnvironmentRoleType>()
          .required()
          .test('isUnassigned', (value, context) => {
            if (value === 'Environment_UNASSIGNED')
              return context.createError({
                message: 'Role is a required field.',
                path: context.path
              });
            return true;
          })
      })
    ),
  tags: yup.array().of(yup.string())
});

const EditMemberModal = ({ isOpen, onClose, member }: EditMemberModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const [tagOptions, setTagOptions] = useState<DropdownOption[]>([]);

  const { data: tagCollection, isLoading: isLoadingTags } = useFetchTags({
    organizationId: currentEnvironment.organizationId,
    entityType: 'ACCOUNT'
  });

  const tags = tagCollection?.tags || [];
  const form = useForm<EditMemberForm>({
    resolver: yupResolver(formSchema) as Resolver<EditMemberForm>,
    defaultValues: {
      firstName: member.firstName,
      lastName: member.lastName,
      language: member.language,
      role: member.organizationRole,
      environmentRoles: member.environmentRoles,
      tags: member.tags
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isSubmitting }
  } = form;

  const { data: collection } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });
  const environments = collection?.environments || [];

  const onSubmit: SubmitHandler<EditMemberForm> = values => {
    return accountUpdater({
      organizationId: currentEnvironment.organizationId,
      email: member.email,
      organizationRole: {
        role: values.role as OrganizationRole
      },
      environmentRoles: values.environmentRoles,
      firstName: values.firstName,
      lastName: values.lastName,
      language: values.language,
      tags: {
        values: values.tags
      }
    }).then(() => {
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
      onClose();
    });
  };

  useEffect(() => {
    if (tags.length > 0) {
      const uniqueTags = uniqBy(tagCollection?.tags || [], 'name')?.map(
        item => ({
          label: item.name,
          value: item.id
        })
      );
      setTagOptions(uniqueTags);
    }
  }, [tags]);

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
                  value={member.email}
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
                          languageList.find(item => item.value === field.value)
                            ?.label
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
            <Form.Field
              control={form.control}
              name={`tags`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label className="relative w-fit">
                    {t('tags')}
                    <Tooltip
                      align="start"
                      alignOffset={-30}
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
                    <TagsSelectMenu
                      tagOptions={tagOptions}
                      fieldValues={field.value}
                      onChange={field.onChange}
                      disabled={isLoadingTags}
                      onChangeTagOptions={setTagOptions}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Divider className="mt-1 mb-3" />

            <EnvironmentRoles environments={environments} />
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
