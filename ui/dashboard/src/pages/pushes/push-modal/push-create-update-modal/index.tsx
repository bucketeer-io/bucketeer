import { useCallback, useMemo, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { pushCreator, PushResponse, pushUpdater, TagChange } from '@api/push';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidatePushes } from '@queries/pushes';
import { useQueryClient } from '@tanstack/react-query';
import { useAuth } from 'auth';
import { ID_NEW } from 'constants/routing';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import uniqBy from 'lodash/uniqBy';
import * as yup from 'yup';
import { Push } from '@types';
import { covertFileToUint8ToBase64 } from 'utils/converts';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { IconInfo } from '@icons';
import { UserMessage } from 'pages/feature-flag-details/targeting/individual-rule';
import { useFetchTags } from 'pages/members/collection-loader';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { CreatableSelect } from 'components/creatable-select';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { Tooltip } from 'components/tooltip';
import UploadFiles from 'components/upload-files';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import EnvironmentEditorList from 'elements/environment-editor-list';
import FormLoading from 'elements/form-loading';

interface PushCreateUpdateModalProps {
  disabled?: boolean;
  isOpen: boolean;
  pushId?: string;
  isLoadingPush: boolean;
  push?: Push;
  resetPush: () => void;
  onClose: (isRefresh?: boolean) => void;
}

export interface PushCreateUpdateForm {
  isEditPush?: boolean;
  name: string;
  fcmServiceAccount?: Uint8Array | string;
  tags?: string[];
  environmentId: string;
}

const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    isEditPush: yup.boolean(),
    name: yup.string().required(requiredMessage),
    fcmServiceAccount: yup.string().when('isEditPush', {
      is: (isEditPush: boolean) => !isEditPush,
      then: schema => schema.required(requiredMessage)
    }),
    tags: yup.array(),
    environmentId: yup.string().required(requiredMessage)
  });

const PushCreateUpdateModal = ({
  disabled,
  isOpen,
  pushId,
  isLoadingPush,
  resetPush,
  push,
  onClose
}: PushCreateUpdateModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();
  const [files, setFiles] = useState<File[]>([]);

  const isEditPush = useMemo(() => pushId !== ID_NEW || !!push, [push, pushId]);

  const editorEnvironments = useMemo(
    () =>
      consoleAccount?.environmentRoles
        .filter(item => item.role === 'Environment_EDITOR')
        ?.map(item => item.environment) || [],
    [consoleAccount]
  );

  const { emptyEnvironmentId } = onFormatEnvironments(editorEnvironments);

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      isEditPush,
      name: push?.name || '',
      tags: push?.tags || [],
      environmentId: isEditPush
        ? push?.environmentId || emptyEnvironmentId
        : '',
      fcmServiceAccount: ''
    }
  });

  const { watch } = form;

  const isEnabledTags = !!watch('environmentId');
  const { data: tagCollection, isLoading: isLoadingTags } = useFetchTags({
    environmentId: checkEnvironmentEmptyId(watch('environmentId')),
    entityType: 'FEATURE_FLAG',
    options: {
      enabled: isEnabledTags,
      gcTime: 0
    }
  });

  const tagOptions = (uniqBy(tagCollection?.tags || [], 'name') || [])?.map(
    tag => ({
      label: tag.name,
      value: tag.name
    })
  );

  const {
    formState: { isValid, isSubmitting, isDirty }
  } = form;

  const handleCheckTags = useCallback(
    (tagValues: string[]) => {
      const tagChanges: TagChange[] = [];
      const tags = push?.tags || [];
      tags?.forEach(item => {
        if (!tagValues.find(tag => tag === item)) {
          tagChanges.push({
            changeType: 'DELETE',
            tag: item
          });
        }
      });
      tagValues.forEach(item => {
        const currentTag = tags.find(tag => tag === item);
        if (!currentTag) {
          tagChanges.push({
            changeType: 'CREATE',
            tag: item
          });
        }
      });

      return tagChanges;
    },
    [push]
  );

  const onSubmit: SubmitHandler<PushCreateUpdateForm> = useCallback(
    async values => {
      try {
        let resp: PushResponse | null = null;
        const { name, tags, environmentId } = values;

        if (isEditPush) {
          resp = await pushUpdater({
            name,
            tagChanges: handleCheckTags(tags || []),
            id: push!.id,
            environmentId: checkEnvironmentEmptyId(environmentId)
          });
        } else {
          const base64String: string = await new Promise(rs =>
            covertFileToUint8ToBase64(files[0], data => rs(data))
          );

          resp = await pushCreator({
            name,
            tags,
            environmentId: checkEnvironmentEmptyId(environmentId),
            fcmServiceAccount: base64String
          });
        }

        if (resp) {
          notify({
            message: t('message:collection-action-success', {
              collection: t('push-notification'),
              action: t('updated')
            })
          });
          invalidatePushes(queryClient);
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [isEditPush, push, files]
  );

  useUnsavedLeavePage({
    isShow: isDirty && !isSubmitting,
    callBackCancel: resetPush
  });

  return (
    <SlideModal
      title={t(isEditPush ? 'edit-push' : 'new-push')}
      isOpen={isOpen}
      onClose={() => onClose(false)}
    >
      {isLoadingPush ? (
        <FormLoading />
      ) : (
        <div className="w-full p-5 pb-28">
          <div className="typo-para-small text-gray-600 mb-3">
            {t('new-push-subtitle')}
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
                        name="push-name"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              {!isEditPush && (
                <Form.Field
                  control={form.control}
                  name="fcmServiceAccount"
                  render={({ field }) => (
                    <Form.Item>
                      <Form.Label required className="relative w-fit">
                        {t('fcm-api-key')}
                        <Tooltip
                          align="start"
                          alignOffset={-76}
                          trigger={
                            <div className="flex-center absolute top-0 -right-6">
                              <Icon
                                icon={IconInfo}
                                size={'sm'}
                                color="gray-500"
                              />
                            </div>
                          }
                          content={t('form:firebase-service-account-tooltip')}
                          className="!z-[100] max-w-[400px]"
                        />
                      </Form.Label>
                      <Form.Control>
                        <UploadFiles
                          files={files}
                          accept={['.json']}
                          acceptTypeText="JSON"
                          onChange={files => {
                            if (files?.length) {
                              field.onChange(files[0]);
                              setFiles(files);
                            } else {
                              field.onChange('');
                              setFiles([]);
                            }
                          }}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              )}
              <Form.Field
                control={form.control}
                name={`environmentId`}
                render={({ field }) => (
                  <Form.Item className="py-2">
                    <Form.Label required>{t('environment')}</Form.Label>
                    <Form.Control>
                      <EnvironmentEditorList
                        value={field.value}
                        disabled={disabled || isEditPush}
                        onSelectOption={field.onChange}
                      />
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
                      {t('form:feature-flag-tags')}
                      <Tooltip
                        align="start"
                        alignOffset={-76}
                        trigger={
                          <div className="flex-center absolute top-0 -right-6">
                            <Icon
                              icon={IconInfo}
                              size={'sm'}
                              color="gray-500"
                            />
                          </div>
                        }
                        content={t('form:feature-flag-tags-tooltip')}
                        className="!z-[100] max-w-[400px]"
                      />
                    </Form.Label>
                    <Form.Control>
                      <CreatableSelect
                        value={field.value?.map(tag => {
                          const tagItem = tagOptions.find(
                            item => item.value === tag
                          );
                          return {
                            label: tagItem?.label || tag,
                            value: tagItem?.value || tag
                          };
                        })}
                        disabled={
                          isLoadingTags || !tagOptions.length || disabled
                        }
                        loading={isLoadingTags}
                        allowCreateWhileLoading={false}
                        isValidNewOption={() => false}
                        isClearable
                        onKeyDown={e => {
                          const { value } = e.target as HTMLInputElement;
                          const isExists = tagOptions.find(
                            item =>
                              item.label
                                .toLowerCase()
                                .includes(value.toLowerCase()) &&
                              !field.value?.includes(item.label)
                          );
                          if (e.key === 'Enter' && (!isExists || !value)) {
                            e.preventDefault();
                          }
                        }}
                        placeholder={t(`form:placeholder-tags`)}
                        options={tagOptions}
                        onChange={value =>
                          field.onChange(value.map(tag => tag.value))
                        }
                        noOptionsMessage={() => (
                          <UserMessage message={t('no-options-found')} />
                        )}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={() => onClose(false)}
                    >
                      {t(`cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <DisabledButtonTooltip
                      align="center"
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

export default PushCreateUpdateModal;
