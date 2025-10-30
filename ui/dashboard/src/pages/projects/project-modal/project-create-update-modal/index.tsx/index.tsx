import { useCallback, useEffect, useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { projectCreator, ProjectResponse, projectUpdater } from '@api/project';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateOrganizations } from '@queries/organizations';
import {
  invalidateProjectDetails,
  useQueryProjectDetails
} from '@queries/project-details';
import { invalidateProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { getAccountAccess, getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_PROJECTS } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { onGenerateSlug } from 'utils/converts';
import { IconInfo } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import FormLoading from 'elements/form-loading';

interface ProjectCreateUpdateModalProps {
  isOpen: boolean;
  onClose: () => void;
}
export interface ProjectCreateUpdateForm {
  name: string;
  urlCode: string;
  description?: string;
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
    description: yup.string()
  });

const ProjectCreateUpdateModal = ({
  isOpen,
  onClose
}: ProjectCreateUpdateModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();
  const { consoleAccount } = useAuth();
  const { envEditable, isOrganizationAdmin } = getAccountAccess(
    consoleAccount!
  );
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const disabled = useMemo(
    () => !envEditable || !isOrganizationAdmin,
    [envEditable, isOrganizationAdmin]
  );

  const { isEdit, params } = useActionWithURL({
    closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}`
  });

  const projectId = useMemo(() => params?.projectId, [params]);

  const organnizationId = useMemo(
    () => params.organizationId || currentEnvironment.organizationId,
    [params, currentEnvironment]
  );

  const {
    data: projectCollections,
    isLoading: projectLoading,
    error: projectError
  } = useQueryProjectDetails({
    params: {
      id: projectId as string,
      organizationId: organnizationId
    },
    enabled: !!projectId
  });

  const projectDetail = useMemo(
    () => projectCollections?.project,
    [projectCollections]
  );

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      name: projectDetail?.name || '',
      urlCode: projectDetail?.urlCode || '',
      description: projectDetail?.description || ''
    },
    mode: 'onChange'
  });
  const { isDirty, isSubmitting } = form.formState;
  const onSubmit: SubmitHandler<ProjectCreateUpdateForm> = useCallback(
    async values => {
      try {
        let resp: ProjectResponse | null = null;
        if (projectDetail) {
          resp = await projectUpdater({
            id: projectDetail.id,
            description: values.description,
            name: values.name
          });
          invalidateProjectDetails(queryClient, {
            id: projectDetail?.id,
            organizationId: projectDetail?.organizationId
          });
        } else {
          resp = await projectCreator({
            ...values,
            organizationId: currentEnvironment.organizationId
          });
        }

        if (resp) {
          invalidateProjects(queryClient);
          invalidateOrganizations(queryClient);
          notify({
            message: t('message:collection-action-success', {
              collection: t('project'),
              action: t(projectDetail ? 'updated' : 'created')
            })
          });
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [projectDetail, currentEnvironment]
  );

  useEffect(() => {
    if (projectError) {
      errorNotify(projectError);
    }
  }, [projectError]);

  useUnsavedLeavePage({ isShow: isDirty && !isSubmitting });

  return (
    <SlideModal
      title={t(isEdit ? 'update-project' : 'new-project')}
      isOpen={isOpen}
      onClose={onClose}
    >
      {projectLoading ? (
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

                          if (!projectDetail) {
                            const isUrlCodeDirty =
                              form.getFieldState('urlCode').isDirty;
                            const urlCode = form.getValues('urlCode');
                            form.setValue(
                              'urlCode',
                              isUrlCodeDirty ? urlCode : onGenerateSlug(value)
                            );
                          }
                        }}
                        name="project-name"
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
                        content={t('form:project-url-tooltip')}
                        className="!z-[100] max-w-[400px]"
                      />
                    </Form.Label>
                    <Form.Control>
                      <Input
                        value={field.value}
                        placeholder={`${t('form:placeholder-code')}`}
                        disabled={!!projectDetail || disabled}
                        name="project-code"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Field
                control={form.control}
                name="description"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label optional>{t('form:description')}</Form.Label>
                    <Form.Control>
                      <TextArea
                        disabled={disabled}
                        placeholder={t('form:placeholder-desc')}
                        rows={4}
                        {...field}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button type="button" variant="secondary" onClick={onClose}>
                      {t(`cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <DisabledButtonTooltip
                      type={!envEditable ? 'editor' : 'admin'}
                      hidden={!disabled}
                      trigger={
                        <Button
                          type="submit"
                          disabled={!form.formState.isDirty || disabled}
                          loading={form.formState.isSubmitting}
                        >
                          {t(isEdit ? `update-project` : 'create-project')}
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

export default ProjectCreateUpdateModal;
