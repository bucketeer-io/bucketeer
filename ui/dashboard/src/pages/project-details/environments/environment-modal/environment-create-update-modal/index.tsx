import { useCallback, useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import {
  environmentCreator,
  EnvironmentResponse,
  environmentUpdater
} from '@api/environment';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateEnvironments } from '@queries/environments';
import {
  invalidateEnvironmentDetails,
  useQueryEnvironmentDetails
} from '@queries/environments-details';
import { invalidateOrganizations } from '@queries/organizations';
import { useQueryProjectDetails } from '@queries/project-details';
import { invalidateProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { getAccountAccess, getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Environment } from '@types';
import { onGenerateSlug } from 'utils/converts';
import { IconInfo } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import Divider from 'components/divider';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import SlideModal from 'components/modal/slide';
import Spinner from 'components/spinner';
import Switch from 'components/switch';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import FormLoading from 'elements/form-loading';

interface EnvironmentCreateUpdateModalProps {
  organizationId: string;
  isOpen: boolean;
  environment?: Environment;
  onClose: () => void;
}

export interface EnvironmentCreateUpdateForm {
  projectId: string;
  name: string;
  urlCode: string;
  description?: string;
  requireComment: boolean;
  autoArchiveEnabled: boolean;
  autoArchiveUnusedDays?: number;
  autoArchiveCheckCodeRefs: boolean;
}

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    name: yup.string().required(requiredMessage),
    urlCode: yup
      .string()
      .required(requiredMessage)
      .matches(
        /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
        translation('message:validation.id-rule', {
          name: translation('common:url-code')
        })
      ),
    description: yup.string(),
    projectId: yup.string().required(requiredMessage),
    requireComment: yup.boolean().required(requiredMessage),
    autoArchiveEnabled: yup.boolean().required(),
    autoArchiveUnusedDays: yup.number().when('autoArchiveEnabled', {
      is: true,
      then: schema => schema.min(1).required(requiredMessage),
      otherwise: schema => schema.nullable()
    }),
    autoArchiveCheckCodeRefs: yup.boolean().required()
  });

const EnvironmentCreateUpdateModal = ({
  organizationId,
  isOpen,
  onClose
}: EnvironmentCreateUpdateModalProps) => {
  const queryClient = useQueryClient();
  const { projectId, environmentId } = useParams();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const { consoleAccount, onMeFetcher } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { envEditable, isOrganizationAdmin } = getAccountAccess(
    consoleAccount!
  );

  const disabled = useMemo(
    () => !envEditable || !isOrganizationAdmin,
    [envEditable, isOrganizationAdmin]
  );

  const { data: collection, isLoading: isLoadingProject } =
    useQueryProjectDetails({
      params: {
        id: projectId!,
        organizationId
      },
      enabled: !!projectId && !!organizationId
    });

  const { data: envCollections, isLoading: isLoadingEnv } =
    useQueryEnvironmentDetails({
      params: {
        id: environmentId as string
      },
      enabled: !!environmentId
    });

  const environmentDetail = useMemo(
    () => envCollections?.environment,
    [envCollections]
  );

  const project = collection?.project;

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      name: environmentDetail?.name || '',
      description: environmentDetail?.description,
      requireComment: environmentDetail?.requireComment || false,
      projectId: projectId || '',
      urlCode: environmentDetail?.urlCode || '',
      autoArchiveEnabled: environmentDetail?.autoArchiveEnabled || false,
      autoArchiveUnusedDays: environmentDetail?.autoArchiveUnusedDays || 90,
      autoArchiveCheckCodeRefs:
        environmentDetail?.autoArchiveCheckCodeRefs || false
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isDirty, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<EnvironmentCreateUpdateForm> = useCallback(
    async values => {
      try {
        let resp: EnvironmentResponse | null = null;
        if (environmentDetail) {
          resp = await environmentUpdater({
            id: environmentDetail!.id,
            name: values.name,
            description: values.description,
            requireComment: values.requireComment,
            autoArchiveEnabled: values.autoArchiveEnabled,
            autoArchiveUnusedDays: values.autoArchiveUnusedDays,
            autoArchiveCheckCodeRefs: values.autoArchiveCheckCodeRefs
          });
        } else {
          resp = await environmentCreator({
            ...values
          });
        }

        if (resp) {
          notify({
            message: t('message:collection-action-success', {
              collection: t('source-type.environment'),
              action: t(environmentDetail ? 'updated' : 'created')
            })
          });
          invalidateOrganizations(queryClient);
          invalidateProjects(queryClient);
          invalidateEnvironments(queryClient);
          invalidateEnvironmentDetails(queryClient);
          onMeFetcher({ organizationId: currentEnvironment.organizationId });
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [environmentDetail, currentEnvironment]
  );
  useUnsavedLeavePage({ isShow: isDirty && !isSubmitting });
  return (
    <SlideModal
      title={t(environmentDetail ? 'update-env' : 'new-env')}
      isOpen={isOpen}
      onClose={onClose}
    >
      {isLoadingEnv ? (
        <FormLoading />
      ) : (
        <div className="w-full p-5">
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
                        disabled={disabled}
                        placeholder={`${t('form:placeholder-name')}`}
                        {...field}
                        onChange={value => {
                          field.onChange(value);
                          if (!environmentDetail) {
                            const isUrlCodeDirty =
                              form.getFieldState('urlCode').isDirty;
                            const urlCode = form.getValues('urlCode');
                            form.setValue(
                              'urlCode',
                              isUrlCodeDirty ? urlCode : onGenerateSlug(value)
                            );
                          }
                        }}
                        name="environment-name"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <Form.Field
                control={form.control}
                name="urlCode"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label required className="relative w-fit">
                      {t('form:url-code')}
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
                        content={t('form:env-url-tooltip')}
                        className="!z-[100] max-w-[400px]"
                      />
                    </Form.Label>
                    <Form.Control>
                      <Input
                        {...field}
                        value={field.value}
                        placeholder={`${t('form:placeholder-code')}`}
                        disabled={disabled || !!environmentDetail}
                        name="environment-code"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <Form.Item>
                <Form.Label required>{`${t(`project`)}`}</Form.Label>
                <Form.Control>
                  <InputGroup
                    addon={isLoadingProject ? <Spinner size="sm" /> : null}
                    addonSize="sm"
                    addonSlot="right"
                    className="w-full"
                  >
                    <Input
                      value={project?.name || ''}
                      placeholder={`${t(`project`)}`}
                      disabled
                    />
                  </InputGroup>
                </Form.Control>
                <Form.Message />
              </Form.Item>

              <Form.Field
                control={form.control}
                name="description"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label optional>{t('form:description')}</Form.Label>
                    <Form.Control>
                      <TextArea
                        placeholder={t('form:placeholder-desc')}
                        rows={4}
                        disabled={disabled}
                        {...field}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <Divider className="mb-5" />
              <h3 className="typo-head-bold-small text-gray-900">
                {t(`form:env-settings`)}
              </h3>
              <Form.Field
                control={form.control}
                name="requireComment"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Control>
                      <Checkbox
                        disabled={disabled}
                        onCheckedChange={checked => field.onChange(checked)}
                        checked={field.value}
                        title={`${t(`form:require-comments-flag`)}`}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <Divider className="my-5" />
              <h3 className="typo-head-bold-small text-gray-900 mb-4">
                {t('form:auto-archive-settings')}
              </h3>

              <Form.Field
                control={form.control}
                name="autoArchiveEnabled"
                render={({ field }) => (
                  <Form.Item className="mb-4">
                    <div className="flex items-center gap-x-3">
                      <Form.Control>
                        <Switch
                          disabled={disabled}
                          checked={field.value}
                          onCheckedChange={checked => field.onChange(checked)}
                        />
                      </Form.Control>
                      <Form.Label className="relative w-fit cursor-pointer mb-0">
                        {t('form:auto-archive-enable')}
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
                          content={t('form:auto-archive-enable-tooltip')}
                          className="!z-[100] max-w-[400px]"
                        />
                      </Form.Label>
                    </div>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              {form.watch('autoArchiveEnabled') && (
                <>
                  <Form.Field
                    control={form.control}
                    name="autoArchiveUnusedDays"
                    render={({ field }) => (
                      <Form.Item className="mb-4">
                        <Form.Label required>
                          {t('form:auto-archive-unused-days')}
                        </Form.Label>
                        <Form.Control>
                          <Input
                            type="number"
                            min={1}
                            disabled={disabled}
                            placeholder={t(
                              'form:auto-archive-unused-days-placeholder'
                            )}
                            {...field}
                            value={field.value ?? ''}
                            onChange={value => {
                              const numValue = value ? parseInt(value, 10) : '';
                              field.onChange(numValue);
                            }}
                          />
                        </Form.Control>
                        <Form.Message />
                      </Form.Item>
                    )}
                  />

                  <Form.Field
                    control={form.control}
                    name="autoArchiveCheckCodeRefs"
                    render={({ field }) => (
                      <Form.Item className="mb-4">
                        <div className="flex items-center gap-x-3">
                          <Form.Control>
                            <Checkbox
                              disabled={disabled}
                              onCheckedChange={checked => field.onChange(checked)}
                              checked={field.value}
                            />
                          </Form.Control>
                          <Form.Label className="relative w-fit cursor-pointer mb-0">
                            {t('form:auto-archive-check-code-refs')}
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
                              content={t(
                                'form:auto-archive-check-code-refs-tooltip'
                              )}
                              className="!z-[100] max-w-[400px]"
                            />
                          </Form.Label>
                        </div>
                        <Form.Message />
                      </Form.Item>
                    )}
                  />
                </>
              )}

              {/* Spacer for fixed ButtonBar */}
              <div className="h-20" />

              <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button type="button" variant="secondary" onClick={onClose}>
                      {t(`cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <DisabledButtonTooltip
                      type={!isOrganizationAdmin ? 'admin' : 'editor'}
                      hidden={!disabled}
                      trigger={
                        <Button
                          type="submit"
                          disabled={!isDirty || !isValid || disabled}
                          loading={isSubmitting}
                        >
                          {t(environmentDetail ? `update-env` : 'create-env')}
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

export default EnvironmentCreateUpdateModal;
