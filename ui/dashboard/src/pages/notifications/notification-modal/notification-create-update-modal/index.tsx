import { useCallback, useMemo, useState } from 'react';
import {
  FormProvider,
  Resolver,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import {
  notificationCreator,
  NotificationResponse,
  notificationUpdater
} from '@api/notification';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateNotifications } from '@queries/notifications';
import { useQueryClient } from '@tanstack/react-query';
import { useAuth } from 'auth';
import { languageList } from 'constants/notification';
import { ID_NEW } from 'constants/routing';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import uniqBy from 'lodash/uniqBy';
import * as yup from 'yup';
import { Notification, NotificationLanguage, SourceType } from '@types';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { cn } from 'utils/style';
import { IconInfo, IconNoData } from '@icons';
import { useFetchTags } from 'pages/members/collection-loader';
import { SOURCE_TYPE_ITEMS } from 'pages/notifications/constants';
import { NotificationOption } from 'pages/notifications/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import { CreatableSelect } from 'components/creatable-select';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import SearchInput from 'components/search-input';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import EnvironmentEditorList from 'elements/environment-editor-list';
import FormLoading from 'elements/form-loading';

interface NotificationCreateUpdateModalProps {
  disabled?: boolean;
  notificationId?: string;
  isOpen: boolean;
  isLoadingNotification: boolean;
  notification?: Notification;
  onClose: () => void;
}

export interface NotificationCreateUpdateForm {
  name: string;
  url: string;
  environment: string;
  language: NotificationLanguage;
  types: SourceType[];
  tags: string[];
}

export const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    name: yup.string().required(requiredMessage),
    url: yup
      .string()
      .required(requiredMessage)
      .url(
        translation('message:validation.id-rule', {
          name: translation('common:url')
        })
      ),
    environment: yup.string().required(requiredMessage),
    language: yup.mixed<NotificationLanguage>().required(requiredMessage),
    types: yup.array().min(1).required(requiredMessage),
    tags: yup.array()
  });

const NotificationCreateUpdateModal = ({
  disabled,
  notificationId,
  isOpen,
  isLoadingNotification,
  notification,
  onClose
}: NotificationCreateUpdateModalProps) => {
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);

  const { consoleAccount } = useAuth();

  const [searchValue, setSearchValue] = useState('');
  const [filteredTypes, setSearchTypes] =
    useState<NotificationOption[]>(SOURCE_TYPE_ITEMS);

  const isEditNotification = useMemo(
    () =>
      notificationId !== ID_NEW || !!notification || !!isLoadingNotification,
    [notification, isLoadingNotification, notificationId]
  );

  const editorEnvironments = useMemo(
    () =>
      consoleAccount?.environmentRoles
        .filter(item => item.role === 'Environment_EDITOR')
        ?.map(item => item.environment) || [],
    [consoleAccount]
  );

  const { emptyEnvironmentId } = onFormatEnvironments(editorEnvironments);

  const form = useForm<NotificationCreateUpdateForm>({
    resolver: yupResolver(
      useFormSchema(formSchema)
    ) as Resolver<NotificationCreateUpdateForm>,
    values: {
      name: notification?.name || '',
      url: notification?.recipient.slackChannelRecipient.webhookUrl || '',
      environment: notification
        ? notification?.environmentId || emptyEnvironmentId
        : '',
      language: notification?.recipient.language || 'ENGLISH',
      types: notification?.sourceTypes.sort() || [],
      tags: notification?.featureFlagTags || []
    },
    mode: 'onChange'
  });

  const {
    watch,
    getValues,
    formState: { isValid, isSubmitting, isDirty }
  } = form;

  const checkedTypes = watch('types');
  const isSelectedEnv = !!watch('environment');

  const { data: tagCollection, isLoading: isLoadingTags } = useFetchTags({
    entityType: 'FEATURE_FLAG',
    environmentId: checkEnvironmentEmptyId(watch('environment') || ''),
    options: {
      enabled: isSelectedEnv
    }
  });
  const tagOptions = uniqBy(tagCollection?.tags || [], 'name')?.map(tag => ({
    label: tag.name,
    value: tag.name
  }));

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

  const onSubmit: SubmitHandler<NotificationCreateUpdateForm> = useCallback(
    async values => {
      try {
        let resp: NotificationResponse | null = null;
        const { name, environment, types, language, tags, url } = values;
        const environmentId = checkEnvironmentEmptyId(environment as string);

        if (isEditNotification) {
          resp = await notificationUpdater({
            id: notification!.id,
            environmentId,
            name,
            sourceTypes: types,
            language,
            featureFlagTags: types.includes('DOMAIN_EVENT_FEATURE') ? tags : []
          });
        } else {
          resp = await notificationCreator({
            environmentId,
            name,
            sourceTypes: types,
            recipient: {
              type: 'SlackChannel',
              slackChannelRecipient: { webhookUrl: url },
              language
            },
            featureFlagTags: tags
          });
        }
        if (resp) {
          notify({
            message: t('message:collection-action-success', {
              collection: t('notification'),
              action: t(isEditNotification ? 'updated' : 'created')
            })
          });
          invalidateNotifications(queryClient);
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [notification, isEditNotification]
  );

  return (
    <SlideModal
      title={t(isEditNotification ? 'update-notification' : 'new-notification')}
      isOpen={isOpen}
      onClose={onClose}
    >
      {isLoadingNotification ? (
        <FormLoading />
      ) : (
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
                        disabled={disabled}
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
                        disabled={disabled || isEditNotification}
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
                      <EnvironmentEditorList
                        value={field.value}
                        onSelectOption={field.onChange}
                        disabled={disabled || isEditNotification}
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
                          label={
                            languageList.find(
                              item => item.value === getValues('language')
                            )?.label
                          }
                          disabled={disabled || isEditNotification}
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

              {filteredTypes.length > 0 ? (
                <Form.Field
                  control={form.control}
                  name={`types`}
                  render={({ field }) => (
                    <>
                      <div className="mt-4 flex items-center justify-between">
                        <div className="typo-para-tiny text-gray-500 uppercase">
                          {t('all-types-selected', {
                            count: checkedTypes.length
                          })}
                        </div>
                        <Checkbox
                          checked={
                            checkedTypes.length === SOURCE_TYPE_ITEMS.length
                          }
                          disabled={disabled}
                          onCheckedChange={checked => {
                            if (checked) {
                              field.onChange(
                                SOURCE_TYPE_ITEMS.map(item => item.value)
                              );
                            } else {
                              field.onChange([]);
                            }
                          }}
                        />
                      </div>
                      <Divider className="mt-3" />

                      {filteredTypes.map(item => (
                        <div key={item.value}>
                          <div className="flex items-center py-3 gap-x-5">
                            <label
                              htmlFor={item.value}
                              className={cn(
                                'flex-1',
                                item.value === 'DOMAIN_EVENT_FEATURE' &&
                                  !isSelectedEnv
                                  ? 'opacity-50'
                                  : 'cursor-pointer'
                              )}
                            >
                              <p className="typo-para-medium text-gray-800">
                                {item.label}
                              </p>
                              <p className="typo-para-small text-gray-600 mt-0.5">
                                {item.description}
                              </p>
                            </label>
                            <Checkbox
                              id={item.value}
                              checked={checkedTypes.includes(item.value)}
                              disabled={
                                (item.value === 'DOMAIN_EVENT_FEATURE' &&
                                  !isSelectedEnv) ||
                                disabled
                              }
                              onCheckedChange={checked => {
                                if (checked) {
                                  checkedTypes.push(item.value);
                                  field.onChange(checkedTypes);
                                } else {
                                  const checkedItems = checkedTypes.filter(
                                    v => v !== item.value
                                  );
                                  field.onChange(checkedItems);
                                }
                              }}
                            />
                          </div>

                          {item.value === 'DOMAIN_EVENT_FEATURE' &&
                            checkedTypes.includes(item.value) && (
                              <Form.Field
                                control={form.control}
                                name={`tags`}
                                render={({ field }) => (
                                  <Form.Item className="-mt-2">
                                    <Form.Label className="relative w-fit">
                                      {t('tags')}
                                      <Tooltip
                                        align="start"
                                        alignOffset={-30}
                                        trigger={
                                          <div className="flex-center absolute top-0 -right-6">
                                            <Icon
                                              icon={IconInfo}
                                              size={'sm'}
                                              color="gray-600"
                                            />
                                          </div>
                                        }
                                        content={t(
                                          'form:tag-notifications-tooltip'
                                        )}
                                        className="!z-[100] max-w-[400px]"
                                      />
                                    </Form.Label>
                                    <Form.Control>
                                      <CreatableSelect
                                        disabled={
                                          isLoadingTags ||
                                          !tagOptions.length ||
                                          disabled
                                        }
                                        value={field.value?.map(tag => {
                                          const tagItem = tagOptions.find(
                                            item => item.value === tag
                                          );
                                          return {
                                            label: tagItem?.label || tag,
                                            value: tagItem?.value || tag
                                          };
                                        })}
                                        loading={isLoadingTags}
                                        placeholder={t(
                                          isSelectedEnv &&
                                            !tagOptions.length &&
                                            !isLoadingTags
                                            ? `form:no-tags-found`
                                            : `form:placeholder-tags`
                                        )}
                                        allowCreateWhileLoading={false}
                                        isValidNewOption={() => false}
                                        isClearable
                                        onKeyDown={e => {
                                          const { value } =
                                            e.target as HTMLInputElement;
                                          const isExists = tagOptions.find(
                                            item =>
                                              item.label
                                                .toLowerCase()
                                                .includes(
                                                  value.toLowerCase()
                                                ) &&
                                              !field.value?.includes(item.label)
                                          );
                                          if (
                                            e.key === 'Enter' &&
                                            (!isExists || !value)
                                          ) {
                                            e.preventDefault();
                                          }
                                        }}
                                        options={tagOptions}
                                        onChange={value =>
                                          field.onChange(
                                            value.map(tag => tag.value)
                                          )
                                        }
                                      />
                                    </Form.Control>
                                    <Form.Message />
                                  </Form.Item>
                                )}
                              />
                            )}
                        </div>
                      ))}
                    </>
                  )}
                />
              ) : (
                <div className="flex flex-col justify-center items-center gap-3 pt-16 pb-4">
                  <IconNoData />
                  <div className="typo-para-medium text-gray-500">
                    {t(`no-data`)}
                  </div>
                </div>
              )}

              <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button variant="secondary" onClick={onClose}>
                      {t(`cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <DisabledButtonTooltip
                      hidden={!disabled}
                      trigger={
                        <Button
                          type="submit"
                          disabled={!isValid || !isDirty || disabled}
                          loading={isSubmitting}
                        >
                          {t(`submit`)}
                        </Button>
                      }
                    />
                  }
                />
              </div>
            </Form>
          </FormProvider>
        </div>
      )}
    </SlideModal>
  );
};

export default NotificationCreateUpdateModal;
