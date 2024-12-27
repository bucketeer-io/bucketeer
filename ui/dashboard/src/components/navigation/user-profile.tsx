import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { IconEditOutlined } from 'react-icons-material-design';
import { accountUpdater } from '@api/account/account-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import primaryAvatar from 'assets/avatars/primary.svg';
// import defaultAvatar from 'assets/avatars/default.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { UserInfoForm } from '@types';
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
    .required()
    .min(2, 'The first name you have provided must have at least 2 characters'),
  lastName: yup
    .string()
    .required()
    .min(2, 'The last name you have provided must have at least 2 characters'),
  language: yup.string().required()
});

export type FilterProps = {
  isOpen: boolean;
  onClose: () => void;
  onEditAvatar: () => void;
};

const UserProfileModal = ({ isOpen, onClose, onEditAvatar }: FilterProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { consoleAccount, onMeFetcher } = useAuth();
  const { notify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      firstName: consoleAccount?.firstName || '',
      lastName: consoleAccount?.lastName || '',
      language: consoleAccount?.language || ''
    }
  });

  const avatar = consoleAccount?.avatarImage
    ? consoleAccount.avatarImage
    : primaryAvatar;

  const onSubmit: SubmitHandler<UserInfoForm> = values => {
    return accountUpdater({
      organizationId: currentEnvironment.organizationId,
      email: consoleAccount!.email,
      changeFirstNameCommand: {
        firstName: values.firstName
      },
      changeLastNameCommand: {
        lastName: values.lastName
      },
      changeLanguageCommand: {
        language: values.language
      }
    }).then(() => {
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: `Profile has been successfully updated!`
      });
      onMeFetcher({ organizationId: currentEnvironment.organizationId });
      onClose();
    });
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
                <AvatarImage
                  image={`data:image/jpeg;base64,${avatar}`}
                  size="xl"
                  alt="user-avatar"
                />
                {/* <AvatarImage image={defaultAvatar} size="xl" alt="user-avatar" /> */}
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
                    <Input placeholder={t(`form:enter-last-name`)} {...field} />
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
                disabled={!form.formState.isDirty || !form.formState.isValid}
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
