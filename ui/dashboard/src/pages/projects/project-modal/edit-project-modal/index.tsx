import { useEffect, useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { projectUpdater } from '@api/project';
import { yupResolver } from '@hookform/resolvers/yup';
import {
  invalidateProjectDetails,
  useQueryProjectDetails
} from '@queries/project-details';
import { invalidateProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import Spinner from 'components/spinner';
import TextArea from 'components/textarea';

interface EditProjectModalProps {
  _projectId?: string;
  isOpen: boolean;
  onClose: () => void;
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
  _projectId,
  isOpen,
  onClose
}: EditProjectModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();
  const { id, errorToast } = useActionWithURL({
    idKey: '*'
  });
  const projectId = useMemo(() => _projectId || id, [id, _projectId]);

  const { data, isLoading, error } = useQueryProjectDetails({
    params: {
      id: projectId as string
    }
  });

  const project = data?.project;

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      description: ''
    }
  });

  const onSubmit: SubmitHandler<EditProjectForm> = async values => {
    try {
      const resp = await projectUpdater({
        id: projectId || '',
        changeDescriptionCommand: {
          description: values.description
        },
        renameCommand: {
          name: values.name
        }
      });
      if (resp) {
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: (
            <span>
              <b>{values.name}</b> {`has been successfully updated!`}
            </span>
          )
        });
        invalidateProjects(queryClient);
        invalidateProjectDetails(queryClient, {
          id: projectId as string
        });
        onClose();
      }
    } catch (error) {
      errorToast(error as Error);
    }
  };

  useEffect(() => {
    if (project)
      form.reset({
        name: project.name,
        description: project.description
      });
  }, [project, form]);

  useEffect(() => {
    if (error) errorToast(error);
  }, [error]);

  return (
    <SlideModal title={t('update-project')} isOpen={isOpen} onClose={onClose}>
      {isLoading ? (
        <div className="flex-center py-10">
          <Spinner />
        </div>
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
                    value={project?.urlCode || ''}
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
                    <Button
                      type="submit"
                      disabled={!form.formState.isDirty}
                      loading={form.formState.isSubmitting}
                    >
                      {t(`update-project`)}
                    </Button>
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

export default EditProjectModal;
