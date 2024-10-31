import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { environmentUpdater } from '@api/environment';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateEnvironments } from '@queries/environments';
import { useQueryProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Environment } from '@types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import Divider from 'components/divider';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';

interface EditEnvironmentModalProps {
  isOpen: boolean;
  onClose: () => void;
  environment: Environment;
}

export interface EditEnvironmentForm {
  name: string;
  description?: string;
  requireComment: boolean;
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  description: yup.string(),
  requireComment: yup.boolean().required()
});

const EditEnvironmentModal = ({
  isOpen,
  onClose,
  environment
}: EditEnvironmentModalProps) => {
  const queryClient = useQueryClient();
  const { projectId } = useParams();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();
  const { data: collection } = useQueryProjects();

  const project = collection?.projects.find(item => item.id === projectId);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: environment.name,
      description: environment.description,
      requireComment: environment.requireComment
    }
  });

  const onSubmit: SubmitHandler<EditEnvironmentForm> = values => {
    return environmentUpdater({
      id: environment.id,
      renameCommand: {
        name: values.name
      },
      changeDescriptionCommand: {
        description: values.description
      },
      changeRequireCommentCommand: {
        requireComment: values.requireComment
      }
    }).then(() => {
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{values.name}</b> {`has been successfully updated!`}
          </span>
        )
      });
      invalidateEnvironments(queryClient);
      onClose();
    });
  };

  return (
    <SlideModal title={t('update-env')} isOpen={isOpen} onClose={onClose}>
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
                  value={environment.urlCode}
                  placeholder={`${t('form:placeholder-code')}`}
                  disabled
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>

            <Form.Item>
              <Form.Label required>{`${t(`project`)}`}</Form.Label>
              <Form.Control>
                <Input
                  value={project?.name || ''}
                  placeholder={`${t(`project`)}`}
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
                    {t(`common:cancel`)}
                  </Button>
                }
                secondaryButton={
                  <Button
                    type="submit"
                    disabled={!form.formState.isDirty}
                    loading={form.formState.isSubmitting}
                  >
                    {t(`update-env`)}
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

export default EditEnvironmentModal;
