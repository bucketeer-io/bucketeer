import { useCallback, useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { projectCreator, ProjectResponse, projectUpdater } from '@api/project';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { getAccountAccess, getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Project } from '@types';
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

interface ProjectCreateUpdateModalProps {
  isOpen: boolean;
  project?: Project;
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
  onClose,
  project
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

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      name: project?.name || '',
      urlCode: project?.urlCode || '',
      description: project?.description || ''
    },
    mode: 'onChange'
  });

  const onSubmit: SubmitHandler<ProjectCreateUpdateForm> = useCallback(
    async values => {
      try {
        let resp: ProjectResponse | null = null;
        if (project) {
          resp = await projectUpdater({
            id: project.id,
            description: values.description,
            name: values.name
          });
        } else {
          resp = await projectCreator({
            ...values,
            organizationId: currentEnvironment.organizationId
          });
        }

        if (resp) {
          invalidateProjects(queryClient);
          notify({
            message: t('message:collection-action-success', {
              collection: t('project'),
              action: t(project ? 'updated' : 'created')
            })
          });
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [project, currentEnvironment]
  );

  return (
    <SlideModal
      title={t(project ? 'update-project' : 'new-project')}
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

                        if (!project) {
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
                      content={t('form:project-url-tooltip')}
                      className="!z-[100] max-w-[400px]"
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input
                      value={field.value}
                      placeholder={`${t('form:placeholder-code')}`}
                      disabled={!!project || disabled}
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
                  <Button variant="secondary" onClick={onClose}>
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
                        {t(project ? `update-project` : 'create-project')}
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

export default ProjectCreateUpdateModal;
