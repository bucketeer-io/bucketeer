import { useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { notificationCreator } from '@api/notification';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateNotifications } from '@queries/notifications';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { languageList } from 'constants/notification';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { NotificationLanguage, SourceType } from '@types';
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
  value: SourceType;
  label: string;
  description: string;
};

export interface AddNotificationForm {
  name: string;
  url: string;
  environment: string;
  language: NotificationLanguage;
  types: SourceType[];
}

export const formSchema = yup.object().shape({
  name: yup.string().required(),
  url: yup.string().required(),
  environment: yup.string().required(),
  language: yup.mixed<NotificationLanguage>().required(),
  types: yup.array().of(yup.string().required()).min(1, 'Required').required()
});

const AddNotificationModal = ({
  isOpen,
  onClose
}: AddNotificationModalProps) => {
  const { notify } = useToast();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const SOURCE_TYPE_ITEMS: NotificationOption[] = [
    {
      label: t(`source-type.account`),
      description: t(`source-type.account-description`),
      value: 'DOMAIN_EVENT_ACCOUNT'
    },
    {
      label: t(`source-type.api-key`),
      description: t(`source-type.api-key-description`),
      value: 'DOMAIN_EVENT_APIKEY'
    },
    {
      label: t(`source-type.auto-ops`),
      description: t(`source-type.auto-ops-description`),
      value: 'DOMAIN_EVENT_AUTOOPS_RULE'
    },
    {
      label: t(`source-type.experiment`),
      description: t(`source-type.experiment-description`),
      value: 'DOMAIN_EVENT_EXPERIMENT'
    },
    {
      label: t(`source-type.feature-flag`),
      description: t(`source-type.feature-flag-description`),
      value: 'DOMAIN_EVENT_FEATURE'
    },
    {
      label: t(`source-type.goal`),
      description: t(`source-type.goal-description`),
      value: 'DOMAIN_EVENT_GOAL'
    },
    {
      label: t(`source-type.mau-count`),
      description: t(`source-type.mau-count-description`),
      value: 'MAU_COUNT'
    },
    {
      label: t(`source-type.notification`),
      description: t(`source-type.notification-description`),
      value: 'DOMAIN_EVENT_SUBSCRIPTION'
    },
    {
      label: t(`source-type.push`),
      description: t(`source-type.push-description`),
      value: 'DOMAIN_EVENT_PUSH'
    },
    {
      label: t(`source-type.running-experiments`),
      description: t(`source-type.running-experiments-description`),
      value: 'EXPERIMENT_RUNNING'
    },
    {
      label: t(`source-type.segment`),
      description: t(`source-type.segment-description`),
      value: 'DOMAIN_EVENT_SEGMENT'
    },
    {
      label: t(`source-type.stale-feature-flag`),
      description: t(`source-type.stale-feature-flag-description`),
      value: 'FEATURE_STALE'
    }
  ];

  const [searchValue, setSearchValue] = useState('');
  const [filteredTypes, setSearchTypes] =
    useState<NotificationOption[]>(SOURCE_TYPE_ITEMS);

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
      types: []
    }
  });

  const {
    watch,
    getValues,
    formState: { isValid, isSubmitting }
  } = form;

  const checkedTypes = watch('types');

  const handleOnChange = (value: string, checked: boolean) => {
    if (checked) {
      checkedTypes.push(value);
      form.setValue('types', checkedTypes);
    } else {
      const checkedItems = checkedTypes.filter(item => item !== value);
      form.setValue('types', checkedItems);
    }
  };

  const onSearchTypes = (value: string) => {
    if (!value) {
      setSearchTypes(SOURCE_TYPE_ITEMS);
      setSearchValue('');
    } else {
      const regex = new RegExp(value, 'i');
      const searchTypes = SOURCE_TYPE_ITEMS.filter(
        item => regex.test(item.label) || regex.test(item.description)
      );
      setSearchTypes(searchTypes);
      setSearchValue(value);
    }
  };

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
    }).then(() => {
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{values.name}</b> {` has been successfully created!`}
          </span>
        )
      });
      invalidateNotifications(queryClient);
      onClose();
    });
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
              value={searchValue}
              onChange={onSearchTypes}
              placeholder={t(`form:search-notification-type`)}
            />

            <div className="mt-4 flex items-center justify-between">
              <div className="typo-para-tiny text-gray-500 uppercase">
                {t('all-types-selected', { count: checkedTypes.length })}
              </div>
              <Checkbox
                checked={checkedTypes.length === SOURCE_TYPE_ITEMS.length}
                onCheckedChange={checked => {
                  if (checked) {
                    form.setValue(
                      'types',
                      SOURCE_TYPE_ITEMS.map(item => item.value)
                    );
                  } else {
                    form.setValue('types', []);
                  }
                }}
              />
            </div>
            <Divider className="mt-3" />

            {filteredTypes.map(item => (
              <div key={item.value} className="flex items-center py-3 gap-x-5">
                <label htmlFor={item.value} className="flex-1 cursor-pointer">
                  <p className="typo-para-medium text-gray-800">{item.label}</p>
                  <p className="typo-para-small text-gray-600 mt-0.5">
                    {item.description}
                  </p>
                </label>
                <Checkbox
                  id={item.value}
                  checked={checkedTypes.includes(item.value)}
                  onCheckedChange={checked =>
                    handleOnChange(item.value, Boolean(checked))
                  }
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
