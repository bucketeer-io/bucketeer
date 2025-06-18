import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { environmentCreator } from '@api/environment';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateEnvironments } from '@queries/environments';
import { useQueryProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { onGenerateSlug } from 'utils/converts';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';

interface AddEnvironmentModalProps {
  disabled?: boolean;
  isOpen: boolean;
  onClose: () => void;
}

export interface AddEnvironmentForm {
  projectId: string;
  name: string;
  urlCode: string;
  description?: string;
  requireComment: boolean;
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  urlCode: yup
    .string()
    .required()
    .matches(
      /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
      "urlCode must start with a letter or number and only contain letters, numbers, or '-'"
    ),
  description: yup.string(),
  projectId: yup.string().required(),
  requireComment: yup.boolean().required()
});

const AddEnvironmentModal = ({
  disabled,
  isOpen,
  onClose
}: AddEnvironmentModalProps) => {
  const queryClient = useQueryClient();
  const { projectId } = useParams();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: projectList } = useQueryProjects({
    params: {
      organizationId: currentEnvironment.organizationId,
      cursor: '0',
      pageSize: 9999
    }
  });

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      urlCode: '',
      projectId: projectId,
      description: '',
      requireComment: false
    }
  });

  const onSubmit: SubmitHandler<AddEnvironmentForm> = async values => {
    try {
      const resp = await environmentCreator({
        ...values
      });
      if (resp) {
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: (
            <span>
              <b>{values.name}</b> {`has been successfully created!`}
            </span>
          )
        });
        invalidateEnvironments(queryClient);
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
    <SlideModal title={t('new-env')} isOpen={isOpen} onClose={onClose}>
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
                      onChange={value => {
                        const isUrlCodeDirty =
                          form.getFieldState('urlCode').isDirty;
                        const urlCode = form.getValues('urlCode');
                        field.onChange(value);
                        form.setValue(
                          'urlCode',
                          isUrlCodeDirty ? urlCode : onGenerateSlug(value)
                        );
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
                  <Form.Label required>{t('form:url-code')}</Form.Label>
                  <Form.Control>
                    <Input
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
              name="projectId"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{`${t(`project`)}`}</Form.Label>
                  <Form.Control className="w-full">
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={`${t(`project`)}`}
                        label={
                          projectList?.projects.find(
                            item => item.id === field.value
                          )?.name
                        }
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[500px]"
                        align="start"
                        {...field}
                      >
                        {projectList?.projects?.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.id}
                            label={item.name}
                            onSelectOption={value => {
                              field.onChange(value);
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
                    disabled={!form.formState.isDirty || disabled}
                    loading={form.formState.isSubmitting}
                  >
                    {t(`create-env`)}
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

export default AddEnvironmentModal;
