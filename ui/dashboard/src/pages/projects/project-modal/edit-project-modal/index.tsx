import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { projectUpdater } from '@api/project';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { requiredMessage } from 'constants/message';
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
  name: yup.string().required(requiredMessage),
  description: yup.string()
});

const EditProjectModal = ({
  isOpen,
  onClose,
  project
}: EditProjectModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

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
          message: t('message:collection-action-success', {
            collection: t('project'),
            action: t('updated')
          })
        });
        onClose();
      }
    } catch (error) {
      errorNotify(error);
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
    </SlideModal>
  );
};

export default EditProjectModal;
