import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Notification } from '@types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
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
import SearchInput from 'components/search-input';

interface EditNotificationModalProps {
  isOpen: boolean;
  onClose: () => void;
  notification: Notification;
}

type NotificationOption = {
  id: string;
  label: string;
  description: string;
};

export interface EditNotificationForm {
  name: string;
  url: string;
  language: string;
  role: string;
}

export const formSchema = yup.object().shape({
  name: yup.string().required(),
  url: yup.string().required(),
  language: yup.string().required(),
  role: yup.string().required()
});

const EditNotificationModal = ({
  isOpen,
  onClose,
  notification
}: EditNotificationModalProps) => {
  const { t } = useTranslation(['common', 'form']);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: notification.name,
      url: '',
      language: '',
      role: ''
    }
  });

  const languages = [
    {
      label: 'English',
      value: 'en'
    }
  ];

  const options: NotificationOption[] = [
    {
      id: 'project',
      label: t('project'),
      description: t('form:notification-type.project')
    },
    {
      id: 'environment',
      label: t('environment'),
      description: t('form:notification-type.environment')
    },
    {
      id: 'account',
      label: t('account'),
      description: t('form:notification-type.account')
    },
    {
      id: 'notification',
      label: t('notification'),
      description: t('form:notification-type.notification')
    }
  ];

  const {
    formState: { isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<EditNotificationForm> = () => {};

  return (
    <SlideModal
      title={t('edit-notification')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="w-full p-5 pb-28">
        <div className="typo-para-small text-gray-600 mb-3">
          {t('new-notification-subtitle')}
        </div>
        <p className="text-gray-800 typo-head-bold-small">
          {t('form:general-info')}
        </p>
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <Form.Field
              control={form.control}
              name="name"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('name')}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:placeholder-name')}`}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="url"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>
                    {t('slack-incoming-webhook')}
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:placeholder-url')}`}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            <Form.Field
              control={form.control}
              name={`language`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('language')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:select-language`)}
                        label=""
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {languages.map((item, index) => (
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

            <Divider className="my-3" />
            <p className="text-gray-800 typo-head-bold-small mb-4">
              {t('types')}
            </p>

            <SearchInput
              value=""
              onChange={() => {}}
              placeholder={t(`form:search-notification-type`)}
            />

            <div className="mt-4 flex items-center justify-between">
              <div className="typo-para-tiny text-gray-500 uppercase">
                {t('all-types-selected', { count: 2 })}
              </div>
              <Checkbox />
            </div>
            <Divider className="mt-3" />

            {options.map(({ id, label, description }) => (
              <div key={id} className="flex items-center py-3 gap-x-5">
                <label htmlFor={id} className="flex-1 cursor-pointer">
                  <p className="typo-para-medium text-gr ay-700">{label}</p>
                  <p className="typo-para-small text-gray-600">{description}</p>
                </label>
                <Form.Field
                  control={form.control}
                  name="role"
                  render={({ field }) => (
                    <Form.Item>
                      <Form.Control>
                        <Checkbox {...field} />
                      </Form.Control>
                    </Form.Item>
                  )}
                />
              </div>
            ))}

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
                    {t(`submit`)}
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

export default EditNotificationModal;
