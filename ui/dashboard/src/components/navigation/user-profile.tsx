import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { IconEditOutlined } from 'react-icons-material-design';
import { yupResolver } from '@hookform/resolvers/yup';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useAuth } from 'auth';
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
  first_name: yup
    .string()
    .required()
    .min(2, 'The first name you have provided must have at least 2 characters'),
  last_name: yup
    .string()
    .required()
    .min(2, 'The last name you have provided must have at least 2 characters'),
  language: yup.string().required()
});

export type FilterProps = {
  isOpen: boolean;
  onClose: () => void;
};

const UserProfileModal = ({ isOpen, onClose }: FilterProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { consoleAccount } = useAuth();

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      first_name: consoleAccount?.firstName || '',
      last_name: consoleAccount?.lastName || '',
      language: consoleAccount?.language || ''
    }
  });

  const avatar = consoleAccount?.avatarUrl
    ? consoleAccount.avatarUrl
    : primaryAvatar;

  const onSubmit: SubmitHandler<UserInfoForm> = values => {
    console.log(values);
  };

  return (
    <DialogModal
      className="w-[466px]"
      title={t('edit-profile')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="p-5">
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="flex items-center justify-center mb-2">
              <div className="relative">
                <AvatarImage image={avatar} size="xl" alt="user-avatar" />
                <Button
                  type="button"
                  className="absolute bottom-0 right-0 size-8 px-1"
                >
                  <Icon icon={IconEditOutlined} size="sm" />
                </Button>
              </div>
            </div>

            <Form.Field
              control={form.control}
              name="first_name"
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
              name="last_name"
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
          </Form>
        </FormProvider>
      </div>

      <ButtonBar
        secondaryButton={<Button>{t(`save`)}</Button>}
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {t(`cancel`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default UserProfileModal;
