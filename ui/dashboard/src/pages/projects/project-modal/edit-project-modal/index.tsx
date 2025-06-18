import { useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { projectUpdater } from '@api/project';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { getAccountAccess, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Project } from '@types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';

interface EditProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
  project: Project;
}

export interface EditProjectForm {
  name: string;
  description?: string;
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  description: yup.string()
});

const EditProjectModal = ({
  isOpen,
  onClose,
  project
}: EditProjectModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();
  const { consoleAccount } = useAuth();
  const { envEditable, isOrganizationAdmin } = getAccountAccess(
    consoleAccount!
  );

  const disabled = useMemo(
    () => !envEditable || !isOrganizationAdmin,
    [envEditable, isOrganizationAdmin]
  );

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: project.name,
      description: project.description
    }
  });

  const onSubmit: SubmitHandler<EditProjectForm> = async values => {
    try {
      const resp = await projectUpdater({
        id: project.id,
        description: values.description,
        name: values.name
      });
      if (resp) {
        invalidateProjects(queryClient);
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: (
            <span>
              <b>{values.name}</b> {`has been successfully updated!`}
            </span>
          )
        });
        onClose();
      }
    } catch (error) {
      notify({
        toastType: 'toast',
        messageType: 'error',
        message: (error as Error)?.message
      });
    }
  };

  return (
    <SlideModal title={t('update-project')} isOpen={isOpen} onClose={onClose}>
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
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            <Form.Item>
              <Form.Label required>{t('form:url-code')}</Form.Label>
              <Form.Control>
                <Input
                  value={project.urlCode}
                  placeholder={`${t('form:placeholder-code')}`}
                  disabled
                />
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
                    {t(`common:cancel`)}
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
                        {t(`update-project`)}
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

export default EditProjectModal;
