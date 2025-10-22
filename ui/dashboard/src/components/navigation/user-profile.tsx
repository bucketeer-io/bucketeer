import { useMemo, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { IconEditOutlined } from 'react-icons-material-design';
import { accountUpdater, AccountAvatar } from '@api/account/account-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import defaultAvatar from 'assets/avatars/default.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import { isNil } from 'lodash';
import * as yup from 'yup';
import { UserInfoForm } from '@types';
import { isNotEmptyObject } from 'utils/data-type';
import UpdatePassword from 'pages/signin/update-password';
import { AvatarImage } from 'components/avatar';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import DialogModal from 'components/modal/dialog';

export const userFormSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) =>
  yup.object().shape({
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
    avatar: yup.string()
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
  const [updatePassword, setUpdatePassword] = useState<boolean | null>(null);
  const { consoleAccount, onMeFetcher } = useAuth();
  const { notify, errorNotify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const form = useForm({
    resolver: yupResolver(useFormSchema(userFormSchema)),
    values: {
      firstName: consoleAccount?.firstName || '',
      lastName: consoleAccount?.lastName || ''
    },
    mode: 'onChange'
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

  const {
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<UserInfoForm> = async values => {
    try {
      if (consoleAccount) {
        const { firstName, lastName } = values;
        const environmentRoles = consoleAccount?.environmentRoles.map(item => ({
          environmentId: item.environment.id,
          role: item.role
        }));
        const resp = await accountUpdater({
          organizationId: currentEnvironment.organizationId,
          email: consoleAccount.email,
          firstName: firstName,
          lastName: lastName,
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
      className="w-[466px] overflow-x-hidden "
      title={updatePassword ? t('update-password') : t('edit-profile')}
      isOpen={isOpen}
      onClose={onClose}
    >
      {!updatePassword ? (
        <div
          className={`w-full ${!updatePassword && !isNil(updatePassword) ? 'animate-translate-right' : ''}`}
        >
          <FormProvider {...form}>
            <Form onSubmit={form.handleSubmit(onSubmit)}>
              <div className="p-5">
                <div className="flex items-center justify-center mb-2">
                  <div className="relative">
                    <AvatarImage
                      image={avatarSrc}
                      size="xl"
                      alt="user-avatar"
                    />
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
                          onChange={value => field.onChange(value)}
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
                          onChange={value => field.onChange(value)}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
                <div className="w-full flex justify-end">
                  <Button
                    type="button"
                    className="bg-transparent border-0 text-primary-500 hover:bg-transparent p-0"
                    onClick={() => setUpdatePassword(true)}
                  >
                    {t('common:change-password')}
                  </Button>
                </div>
              </div>
              <ButtonBar
                secondaryButton={
                  <Button
                    disabled={(!selectedAvatar && !isDirty) || !isValid}
                    loading={isSubmitting}
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
        </div>
      ) : (
        <UpdatePassword
          className="animate-translate-left"
          onBack={() => setUpdatePassword(false)}
          onCancel={onClose}
        />
      )}
    </DialogModal>
  );
};

export default UserProfileModal;
