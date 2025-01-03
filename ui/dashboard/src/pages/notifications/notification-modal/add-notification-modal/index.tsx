import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { notificationCreator } from '@api/notification';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, useAuth } from 'auth';
import { languageList } from 'constants/notification';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { NotificationLanguage } from '@types';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
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

interface AddNotificationModalProps {
  isOpen: boolean;
  onClose: () => void;
}

type NotificationOption = {
  id: string;
  label: string;
  description: string;
};

export interface AddNotificationForm {
  name: string;
  url: string;
  environment: string;
  language: NotificationLanguage;
  type: string;
}

export const formSchema = yup.object().shape({
  name: yup.string().required(),
  url: yup.string().required(),
  environment: yup.string().required(),
  language: yup.mixed<NotificationLanguage>().required(),
  type: yup.string().required()
});

const AddNotificationModal = ({
  isOpen,
  onClose
}: AddNotificationModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: collection, isLoading: isLoadingEnvs } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });
  const environments = (collection?.environments || []).filter(item => item.id);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      url: '',
      environment: '',
      language: undefined,
      type: ''
    }
  });

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
    getValues,
    formState: { isValid, isSubmitting }
  } = form;
  // 'https://dev.bucketeer.jp/hookauth=CiQAQFReLhnIle3NdlT3KBlNsZInL46XvTqeFrEf_yYlZdbJoIISgwEAemffGZYq1vkzNUV4CPfYEgIJt1y9enp1B36b_XGNds58ELMAOWXP5q84peCShNIXjareVnaThwO73_RJP5STk-gbdhxF_TWDDejo_6y1zI9iOqlqLetAxM7GTnfBGd9DnpsLaLucKnKvGyGkgwVX06l6Mw2ovP30ZaMU6HIQbFLl9A'

  const onSubmit: SubmitHandler<AddNotificationForm> = values => {
    return notificationCreator({
      environmentId: values.environment,
      name: values.name,
      sourceTypes: ['DOMAIN_EVENT_FEATURE'],
      recipient: {
        type: 'SlackChannel',
        slackChannelRecipient: { webhookUrl: values.url },
        language: values.language
      }
    }).then(() => {});
  };

  return (
    <SlideModal title={t('new-notification')} isOpen={isOpen} onClose={onClose}>
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
              name={`environment`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('environment')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:select-environment`)}
                        label={
                          environments.find(
                            item => item.id === getValues('environment')
                          )?.name
                        }
                        disabled={isLoadingEnvs}
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {environments.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.id}
                            label={item.name}
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
              name={`language`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('language')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:select-language`)}
                        label={
                          languageList.find(
                            item => item.value === getValues('language')
                          )?.label
                        }
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
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

            <Divider className="my-3" />
            <p className="text-gray-800 typo-head-bold-small mb-4">
              {t('types')}
            </p>

            <SearchInput
              value={''}
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
                  name="type"
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

export default AddNotificationModal;
