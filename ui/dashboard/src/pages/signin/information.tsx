import { useCallback, useMemo, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { IconEditOutlined } from 'react-icons-material-design';
import { useNavigate } from 'react-router-dom';
import { AccountAvatar, accountUpdater } from '@api/account/account-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import defaultAvatar from 'assets/avatars/default.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { getLanguage, Language, useTranslation } from 'i18n';
import { clearIsLoginFirstTimeStorage } from 'storage/login';
import * as yup from 'yup';
import { isNotEmptyObject } from 'utils/data-type';
import { languageList } from 'pages/members/member-modal/add-member-modal';
import { AvatarImage } from 'components/avatar';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import EditPhotoProfileModal from 'components/navigation/edit-photo';
import UploadAvatarModal from 'components/navigation/upload-avatar';
import AuthWrapper from './elements/auth-wrapper';
import UpdatePassword from './update-password';

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
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
    avatar: yup.string(),
    language: yup.string().required(requiredMessage)
  });

interface UserInfoForm {
  firstName: string;
  lastName: string;
  avatar?: string;
  language: string;
}

const UserInformation = () => {
  const { t } = useTranslation(['auth', 'form', 'common', 'message']);
  const { consoleAccount, onMeFetcher } = useAuth();
  const [isOpenEditAvatarModal, onOpenEditAvatarModal, onCloseEditAvatarModal] =
    useToggleOpen(false);
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { notify, errorNotify } = useToast();
  const navigate = useNavigate();

  const [
    isOpenUploadAvatarModal,
    onOpenUploadAvatarModal,
    onCloseUploadAvatarModal
  ] = useToggleOpen(false);

  const [selectedAvatar, setSelectedAvatar] = useState<AccountAvatar | null>(
    null
  );

  const avatarType = useMemo(
    () => selectedAvatar?.avatarFileType || consoleAccount?.avatarFileType,
    [selectedAvatar, consoleAccount]
  );

  const avatarImage = useMemo(
    () => selectedAvatar?.avatarImage || consoleAccount?.avatarImage,
    [selectedAvatar, consoleAccount]
  );

  const isUserAvatar = useMemo(
    () => selectedAvatar || avatarImage,
    [selectedAvatar, avatarImage]
  );

  const avatarSrc = useMemo(
    () =>
      isUserAvatar
        ? `data:${avatarType};base64,${avatarImage}`
        : consoleAccount?.avatarUrl || defaultAvatar,
    [avatarImage, avatarType, defaultAvatar, isUserAvatar, consoleAccount]
  );

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      firstName: consoleAccount?.firstName || '',
      lastName: consoleAccount?.lastName || '',
      avatar: avatarSrc,
      language: getLanguage() || Language.ENGLISH
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isSubmitting },
    setValue
  } = form;

  const onSelectAvatar = useCallback(
    (avatar: AccountAvatar | null, cb?: () => void) => {
      setValue(
        'avatar',
        avatar
          ? `data:${avatar?.avatarFileType};base64,${avatar?.avatarImage}`
          : defaultAvatar,
        {
          shouldDirty: true
        }
      );
      setSelectedAvatar(avatar);
      if (cb) cb();
    },
    [defaultAvatar]
  );

  const onSubmit: SubmitHandler<UserInfoForm> = useCallback(
    async values => {
      try {
        if (consoleAccount) {
          const { firstName, lastName, language } = values;
          const resp = await accountUpdater({
            organizationId: currentEnvironment.organizationId,
            email: consoleAccount.email,
            firstName: firstName,
            lastName: lastName,
            ...(selectedAvatar && isNotEmptyObject(selectedAvatar)
              ? {
                  avatar: selectedAvatar
                }
              : {}),
            language
          });

          if (resp) {
            notify({
              message: t('message:collection-action-success', {
                collection: t('common:profile'),
                action: t('common:updated')
              })
            });
            onMeFetcher({ organizationId: currentEnvironment.organizationId });
            setSelectedAvatar(null);
            clearIsLoginFirstTimeStorage();
            navigate(PAGE_PATH_ROOT);
          }
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [consoleAccount, currentEnvironment, selectedAvatar]
  );

  return (
    <AuthWrapper>
      <div className="grid gap-10">
        <div>
          <h1 className="text-gray-900 typo-head-bold-huge mb-4">
            {t(`enter-information.title`)}
          </h1>
          <p className="text-gray-600 typo-para-medium">
            {t(`enter-information.description`)}
          </p>
        </div>

        <FormProvider {...form}>
          <Form
            onSubmit={form.handleSubmit(onSubmit)}
            className="flex flex-col w-full gap-y-5"
          >
            <Form.Field
              control={form.control}
              name="avatar"
              render={({ field }) => (
                <Form.Item className="py-0">
                  <div className="flex items-center justify-center">
                    <div className="relative">
                      <AvatarImage
                        image={field.value as string}
                        size="xl"
                        alt="user-avatar"
                      />
                      <Button
                        type="button"
                        className="absolute bottom-0 right-0 size-8 px-1"
                        onClick={onOpenUploadAvatarModal}
                      >
                        <Icon icon={IconEditOutlined} size="sm" />
                      </Button>
                    </div>
                  </div>
                </Form.Item>
              )}
            />

            <Form.Field
              control={form.control}
              name="firstName"
              render={({ field }) => (
                <Form.Item className="py-0">
                  <Form.Label required>{t(`first-name`)}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t(`first-name-placeholder`)}
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
                <Form.Item className="py-0">
                  <Form.Label required>{t(`last-name`)}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t(`last-name-placeholder`)}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="language"
              render={({ field }) => {
                const currentItem = languageList.find(
                  item => item.value === field.value
                );

                return (
                  <Form.Item className="py-0">
                    <Form.Label required>{t('common:language')}</Form.Label>
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
                          isExpand
                          className="w-[400px]"
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
            <Button
              type="submit"
              loading={isSubmitting}
              disabled={!isValid}
              className="mt-5 w-full"
            >
              {t('common:continue')}
            </Button>
          </Form>
        </FormProvider>
        <p>pass</p>
        <UpdatePassword className="animate-slide-left" />
      </div>
      {isOpenUploadAvatarModal && (
        <UploadAvatarModal
          isOpen={isOpenUploadAvatarModal}
          onClose={onCloseUploadAvatarModal}
          onUploadPhoto={() => {
            onCloseUploadAvatarModal();
            onOpenEditAvatarModal();
          }}
          onSelectAvatar={avatar =>
            onSelectAvatar(avatar, onCloseUploadAvatarModal)
          }
        />
      )}
      {isOpenEditAvatarModal && (
        <EditPhotoProfileModal
          onUpload={avatar => onSelectAvatar(avatar, onCloseEditAvatarModal)}
          isOpen={isOpenEditAvatarModal}
          onClose={onCloseEditAvatarModal}
        />
      )}
    </AuthWrapper>
  );
};

export default UserInformation;
