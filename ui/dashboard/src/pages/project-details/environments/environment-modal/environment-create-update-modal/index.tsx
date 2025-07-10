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
import { invalidateOrganizations } from '@queries/organizations';
import { useQueryProjectDetails } from '@queries/project-details';
import { invalidateProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { getAccountAccess, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
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
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';

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
    requireComment: yup.boolean().required(requiredMessage)
  });

const EnvironmentCreateUpdateModal = ({
  organizationId,
  isOpen,
  onClose,
  environment
}: EnvironmentCreateUpdateModalProps) => {
  const queryClient = useQueryClient();
  const { projectId } = useParams();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const { consoleAccount } = useAuth();

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

  const project = collection?.project;

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      name: environment?.name || '',
      description: environment?.description,
      requireComment: environment?.requireComment || false,
      projectId: projectId || '',
      urlCode: environment?.urlCode || ''
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
        if (environment) {
          resp = await environmentUpdater({
            id: environment!.id,
            name: values.name,
            description: values.description,
            requireComment: values.requireComment
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
              action: t(environment ? 'updated' : 'created')
            })
          });
          invalidateOrganizations(queryClient);
          invalidateProjects(queryClient);
          invalidateEnvironments(queryClient);
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [environment]
  );

  return (
    <SlideModal
      title={t(environment ? 'update-env' : 'new-env')}
      isOpen={isOpen}
      onClose={onClose}
    >
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
                        if (!environment) {
                          const isUrlCodeDirty =
                            form.getFieldState('urlCode').isDirty;
                          const urlCode = form.getValues('urlCode');
                          form.setValue(
                            'urlCode',
                            isUrlCodeDirty ? urlCode : onGenerateSlug(value)
                          );
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
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </div>
                      }
                      content={t('form:env-url-tooltip')}
                      className="!z-[100] max-w-[400px]"
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input
                      value={field.value}
                      placeholder={`${t('form:placeholder-code')}`}
                      disabled={disabled || !!environment}
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

            <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
              <ButtonBar
                primaryButton={
                  <Button variant="secondary" onClick={onClose}>
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
                        {t(environment ? `update-env` : 'create-env')}
                      </Button>
                    }
                  />
                }
              />
            </div>
          </Form>
        </FormProvider>
      </div>
    </SlideModal>
  );
};

export default EnvironmentCreateUpdateModal;
