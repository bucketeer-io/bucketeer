import { useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { IconEditOutlined } from 'react-icons-material-design';
import { accountUpdater, AccountAvatar } from '@api/account/account-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import defaultAvatar from 'assets/avatars/default.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { requiredMessage, translation } from 'constants/message';
import { useToast } from 'hooks';
import { Language, useTranslation } from 'i18n';
import * as yup from 'yup';
import { UserInfoForm } from '@types';
import { isNotEmptyObject } from 'utils/data-type';
import { onChangeFontWithLocalized } from 'utils/function';
import { languageList } from 'pages/members/member-modal/add-member-modal';
import { AvatarImage } from 'components/avatar';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import DialogModal from 'components/modal/dialog';

const formSchema = yup.object().shape({
  firstName: yup
    .string()
    .required(requiredMessage)
    .min(
      2,
      translation('message:validation.name-at-least-characters', {
        count: 2,
        name: translation('common:first-name').toLowerCase()
      })
    ),
  lastName: yup
    .string()
    .required(requiredMessage)
    .min(
      2,
      translation('message:validation.name-at-least-characters', {
        count: 2,
        name: translation('common:last-name').toLowerCase()
      })
    ),
  language: yup.string().required(requiredMessage)
});
export type FilterProps = {
  selectedAvatar: AccountAvatar | null;
  isOpen: boolean;
  onClose: () => void;
  onEditAvatar: () => void;
};

const UserProfileModal = ({
  selectedAvatar,
  isOpen,
  onClose,
  onEditAvatar
}: FilterProps) => {
  const { t } = useTranslation(['common', 'form', 'message']);
  const { consoleAccount, onMeFetcher } = useAuth();
  const { notify, errorNotify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      firstName: consoleAccount?.firstName || '',
      lastName: consoleAccount?.lastName || '',
      language: consoleAccount?.language || ''
    }
  });

  const avatarType =
    selectedAvatar?.avatarFileType || consoleAccount?.avatarFileType;
  const avatar = selectedAvatar?.avatarImage || consoleAccount?.avatarImage;
  const isUserAvatar = selectedAvatar || avatar;

  const avatarSrc = useMemo(
    () =>
      isUserAvatar ? `data:${avatarType};base64,${avatar}` : defaultAvatar,
    [avatar, selectedAvatar, defaultAvatar]
  );

  const { trigger } = form;

  const onSubmit: SubmitHandler<UserInfoForm> = async values => {
    try {
      if (consoleAccount) {
        const { firstName, lastName, language } = values;
        const environmentRoles = consoleAccount?.environmentRoles.map(item => ({
          environmentId: item.environment.id,
          role: item.role
        }));
        const resp = await accountUpdater({
          organizationId: currentEnvironment.organizationId,
          email: consoleAccount.email,
          firstName: firstName,
          lastName: lastName,
          language: language,
          environmentRoles,
          ...(selectedAvatar && isNotEmptyObject(selectedAvatar)
            ? {
                avatar: selectedAvatar
              }
            : {})
        });

        if (resp) {
          notify({
            message: t('message:collection-action-success', {
              collection: t('profile'),
              action: t('updated')
            })
          });
          onChangeFontWithLocalized(language === Language.JAPANESE);
          onMeFetcher({ organizationId: currentEnvironment.organizationId });
          onClose();
        }
      }
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <DialogModal
      className="w-[466px]"
      title={t('edit-profile')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="p-5">
            <div className="flex items-center justify-center mb-2">
              <div className="relative">
                <AvatarImage image={avatarSrc} size="xl" alt="user-avatar" />
                <Button
                  type="button"
                  className="absolute bottom-0 right-0 size-8 px-1"
                  onClick={onEditAvatar}
                >
                  <Icon icon={IconEditOutlined} size="sm" />
                </Button>
              </div>
            </div>

            <Form.Field
              control={form.control}
              name="firstName"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t(`first-name`)}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t(`form:enter-first-name`)}
                      {...field}
                      onChange={value => {
                        field.onChange(value);
                        trigger('firstName');
                      }}
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
                  <Form.Label required>{t(`last-name`)}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t(`form:enter-last-name`)}
                      {...field}
                      onChange={value => {
                        field.onChange(value);
                        trigger('lastName');
                      }}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
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
                        className="w-[426px]"
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
                              field.onBlur();
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
          </div>
          <ButtonBar
            secondaryButton={
              <Button
                disabled={
                  (!selectedAvatar && !form.formState.isDirty) ||
                  !form.formState.isValid
                }
                loading={form.formState.isSubmitting}
              >
                {t(`save`)}
              </Button>
            }
            primaryButton={
              <Button onClick={onClose} variant="secondary">
                {t(`cancel`)}
              </Button>
            }
          />
        </Form>
      </FormProvider>
    </DialogModal>
  );
};

export default UserProfileModal;
