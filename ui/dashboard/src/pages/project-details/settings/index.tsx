import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { IconLaunchOutlined } from 'react-icons-material-design';
import { useParams } from 'react-router-dom';
import { projectUpdater } from '@api/project';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateAccounts } from '@queries/accounts';
import { invalidateProjectDetails } from '@queries/project-details';
import { useQueryClient } from '@tanstack/react-query';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Project } from '@types';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import TextArea from 'components/textarea';

const formSchema = yup.object().shape({
  name: yup.string().required(),
  urlCode: yup.string().required(),
  description: yup.string()
});

export interface ProjectSettingsForm {
  name: string;
  urlCode: string;
  description?: string;
}

const ProjectSettings = ({ project }: { project: Project }) => {
  const { notify } = useToast();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const params = useParams();
  const projectDetailsId = params.projectId!;

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: project.name,
      urlCode: project.urlCode,
      description: project.description
    }
  });

  const onSubmit: SubmitHandler<ProjectSettingsForm> = async values => {
    try {
      const resp = await projectUpdater({
        id: projectDetailsId,
        description: values.description,
        name: values.name
      });
      if (resp) {
        invalidateProjectDetails(queryClient, {
          id: projectDetailsId,
          organizationId: project.organizationId
        });
        invalidateAccounts(queryClient);
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: (
            <span>
              <b>{values.name}</b> {`has been successfully updated!`}
            </span>
          )
        });
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
    <div className="flex flex-col w-full p-6">
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex lg:items-center justify-between flex-col lg:flex-row">
            <p className="typo-head-bold-big">{t(`settings`)}</p>
            <div className="flex items-center gap-4 mt-3 lg:mt-0">
              <Button
                variant="text"
                type="button"
                className="flex-1 lg:flex-none"
              >
                <Icon icon={IconLaunchOutlined} size="sm" />
                {t('documentation')}
              </Button>
              <Button
                loading={form.formState.isSubmitting}
                disabled={!form.formState.isDirty}
                type="submit"
                className="w-[120px]"
              >
                {t(`save`)}
              </Button>
            </div>
          </div>
          <div className="p-5 shadow-card rounded-lg bg-white mt-6">
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
              name="urlCode"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('form:url-code')}</Form.Label>
                  <Form.Control>
                    <Input
                      disabled
                      placeholder={`${t('form:placeholder-code')}`}
                      {...field}
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
                      placeholder={t('form:placeholder-desc')}
                      rows={4}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </div>
        </Form>
      </FormProvider>
    </div>
  );
};

export default ProjectSettings;
