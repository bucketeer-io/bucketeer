import { useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { projectCreator } from '@api/project';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateProjects } from '@queries/projects';
import { useQueryClient } from '@tanstack/react-query';
import { getAccountAccess, getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
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

interface AddProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface AddProjectForm {
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
    description: yup.string(),
    id: yup.string()
  });

const AddProjectModal = ({ isOpen, onClose }: AddProjectModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { envEditable, isOrganizationAdmin } = getAccountAccess(
    consoleAccount!
  );

  const disabled = useMemo(
    () => !envEditable || !isOrganizationAdmin,
    [envEditable, isOrganizationAdmin]
  );

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      name: '',
      urlCode: '',
      description: ''
    }
  });

  const onSubmit: SubmitHandler<AddProjectForm> = async values => {
    try {
      const resp = await projectCreator({
        ...values,
        organizationId: currentEnvironment.organizationId
      });

      if (resp) {
        notify({
          message: t('message:collection-action-success', {
            collection: t('project'),
            action: t('created')
          })
        });
        invalidateProjects(queryClient);
        onClose();
      }
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <SlideModal title={t('new-project')} isOpen={isOpen} onClose={onClose}>
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
                      disabled={disabled}
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
                      disabled={disabled}
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
                      disabled={disabled}
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
                        {t(`create-project`)}
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

export default AddProjectModal;
