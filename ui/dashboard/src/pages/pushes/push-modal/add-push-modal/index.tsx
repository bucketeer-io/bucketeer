import { useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { pushCreator } from '@api/push';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidatePushes } from '@queries/pushes';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import uniqBy from 'lodash/uniqBy';
import * as yup from 'yup';
import { covertFileToUint8ToBase64 } from 'utils/converts';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { IconInfo } from '@icons';
import { UserMessage } from 'pages/feature-flag-details/targeting/individual-rule';
import { useFetchTags } from 'pages/members/collection-loader';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { CreatableSelect } from 'components/creatable-select';
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
import UploadFiles from 'components/upload-files';

interface AddPushModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface AddPushForm {
  name: string;
  fcmServiceAccount: Uint8Array | string;
  tags?: string[];
  environmentId: string;
}

export const formSchema = yup.object().shape({
  name: yup.string().required(),
  fcmServiceAccount: yup.string().required(),
  tags: yup.array(),
  environmentId: yup.string().required()
});

const AddPushModal = ({ isOpen, onClose }: AddPushModalProps) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify, errorNotify } = useToast();

  const [files, setFiles] = useState<File[]>([]);

  const { data: collection, isLoading: isLoadingEnvs } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });
  const environments = collection?.environments || [];
  const { formattedEnvironments } = onFormatEnvironments(environments);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      fcmServiceAccount: '',
      tags: undefined,
      environmentId: ''
    }
  });

  const {
    watch,
    getValues,
    formState: { isValid, isSubmitting }
  } = form;

  const isEnabledTags = !!watch('environmentId');

  const { data: tagCollection, isLoading: isLoadingTags } = useFetchTags({
    entityType: 'FEATURE_FLAG',
    environmentId: watch('environmentId'),
    options: {
      enabled: isEnabledTags
    }
  });

  const tagOptions = (uniqBy(tagCollection?.tags || [], 'name') || [])?.map(
    tag => ({
      label: tag.name,
      value: tag.name
    })
  );

  const onSubmit: SubmitHandler<AddPushForm> = async values => {
    try {
      const base64String: string = await new Promise(rs =>
        covertFileToUint8ToBase64(files[0], data => rs(data))
      );
      const { environmentId } = values;
      const resp = await pushCreator({
        ...values,
        environmentId: checkEnvironmentEmptyId(environmentId),
        fcmServiceAccount: base64String
      });
      if (resp) {
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: (
            <span>
              <b>{values.name}</b> {` has been successfully created!`}
            </span>
          )
        });
        invalidatePushes(queryClient);
        onClose();
      }
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <SlideModal title={t('new-push')} isOpen={isOpen} onClose={onClose}>
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
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="fcmServiceAccount"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required className="relative w-fit">
                    {t('fcm-api-key')}
                    <Icon
                      icon={IconInfo}
                      className="absolute -right-8"
                      size={'sm'}
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

            <Form.Field
              control={form.control}
              name={`environmentId`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('environment')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:select-environment`)}
                        label={
                          formattedEnvironments.find(
                            item => item.id === getValues('environmentId')
                          )?.name || ''
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
                        {formattedEnvironments.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.id}
                            label={item.name}
                            onSelectOption={value => {
                              field.onChange(value);
                              form.setValue('tags', []);
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
              name={`tags`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label>{t('form:feature-flag-tags')}</Form.Label>
                  <Form.Control>
                    <CreatableSelect
                      disabled={
                        isLoadingTags || !isEnabledTags || !tagOptions.length
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
                      placeholder={t(
                        isEnabledTags && !tagOptions.length && !isLoadingTags
                          ? `form:no-tags-found`
                          : `form:placeholder-tags`
                      )}
                      options={tagOptions}
                      onChange={value =>
                        field.onChange(value.map(tag => tag.value))
                      }
                      noOptionsMessage={() => (
                        <UserMessage message={t('no-options-found')} />
                      )}
                      onCreateOption={() => {}}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

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

export default AddPushModal;
