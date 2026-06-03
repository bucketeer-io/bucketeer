import { useRef } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { goalCreator } from '@api/goal';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { ConnectionType } from '@types';
import { onGenerateSlug } from 'utils/converts';
import { checkFieldDirty } from 'utils/function';
import { IconInfo } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';

interface AddGoalModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface AddGoalForm {
  id: string;
  name: string;
  connectionType?: string;
  description?: string;
}

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    id: yup
      .string()
      .required(requiredMessage)
      .matches(
        /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
        translation('message:validation.id-rule', {
          name: translation('common:source-type.feature-flag')
        })
      ),
    name: yup.string().required(requiredMessage),
    description: yup.string(),
    connectionType: yup.string()
  });

const AddGoalModal = ({ isOpen, onClose }: AddGoalModalProps) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const isIdEdited = useRef(false);

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      id: '',
      name: '',
      connectionType: 'EXPERIMENT',
      description: ''
    }
  });

  const addSuccess = () => {
    {
      notify({
        message: t('message:collection-action-success', {
          collection: t('source-type.goal'),
          action: t('created')
        })
      });
      onClose();
    }
  };

  const onSubmit: SubmitHandler<AddGoalForm> = async values => {
    try {
      const resp = await goalCreator({
        ...values,
        connectionType: values.connectionType as ConnectionType,
        environmentId: currentEnvironment.id
      });
      if (resp.goal) addSuccess();
    } catch (error) {
      errorNotify(error);
    }
  };

  const {
    formState: { isDirty, isSubmitting, dirtyFields }
  } = form;

  useUnsavedLeavePage({
    isShow: checkFieldDirty(dirtyFields) && !isSubmitting
  });

  return (
    <SlideModal title={t('new-goal')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full p-5">
        <p className="text-gray-800 dark:text-dark-gray-400 typo-head-bold-small">
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
                        field.onChange(value);
                        if (!isIdEdited.current) {
                          form.setValue('id', onGenerateSlug(value), {
                            shouldDirty: false,
                            shouldValidate: true
                          });
                        }
                      }}
                      name="goal-name"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="id"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required className="relative w-fit">
                    {t('form:goal-id')}
                    <Tooltip
                      align="start"
                      alignOffset={-62}
                      trigger={
                        <div className="flex-center absolute top-0 -right-6">
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </div>
                      }
                      content={t('form:goal-id-tooltip')}
                      className="!z-[100] max-w-[400px]"
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:placeholder-goal-id')}`}
                      {...field}
                      name="goal-id"
                      autoComplete="off"
                      onChange={value => {
                        isIdEdited.current = value !== '';
                        field.onChange(value);
                      }}
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
            <Form.Field
              control={form.control}
              name="connectionType"
              render={({ field }) => (
                <Form.Item className="flex flex-col w-full py-0 gap-y-4">
                  <Form.Label>{t('form:connections')}</Form.Label>
                  <Form.Control>
                    <RadioGroup
                      defaultValue={field.value}
                      className="flex flex-col w-full gap-y-5"
                      onValueChange={field.onChange}
                    >
                      <div className="flex items-center gap-x-2">
                        <RadioGroupItem
                          id="experiments-connection"
                          value="EXPERIMENT"
                        />
                        <label
                          htmlFor="experiments-connection"
                          className="typo-para-medium leading-4 text-gray-700 dark:text-dark-gray-300 cursor-pointer"
                        >
                          {t('experiments')}
                        </label>
                      </div>

                      <div className="flex items-center gap-x-2">
                        <RadioGroupItem
                          id="operations-connection"
                          value="OPERATION"
                        />
                        <label
                          htmlFor="operations-connection"
                          className="typo-para-medium leading-4 text-gray-700 dark:text-dark-gray-300 cursor-pointer"
                        >
                          {t('form:operations')}
                        </label>
                      </div>
                    </RadioGroup>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <div className="absolute left-0 bottom-0 bg-gray-50 dark:bg-dark-black-800 w-full rounded-b-lg">
              <ButtonBar
                primaryButton={
                  <Button type="button" variant="secondary" onClick={onClose}>
                    {t(`cancel`)}
                  </Button>
                }
                secondaryButton={
                  <Button
                    type="submit"
                    disabled={!isDirty || isSubmitting}
                    loading={isSubmitting}
                  >
                    {t(`submit`)}
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

export default AddGoalModal;
