import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { goalCreator } from '@api/goal';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateGoals } from '@queries/goals';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { AxiosError } from 'axios';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { ConnectionType } from '@types';
import { IconInfo } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';

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

const formSchema = yup.object().shape({
  id: yup
    .string()
    .required()
    .matches(
      /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
      "Goal ID must start with a letter or number and only contain letters, numbers, or '-'"
    ),
  name: yup.string().required(),
  description: yup.string(),
  connectionType: yup.string()
});

const AddGoalModal = ({ isOpen, onClose }: AddGoalModalProps) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: '',
      name: '',
      connectionType: 'EXPERIMENT',
      description: ''
    }
  });

  const addSuccess = (name: string) => {
    {
      notify({
        message: (
          <span>
            <b>{name}</b> {`has been successfully created!`}
          </span>
        )
      });
      invalidateGoals(queryClient);
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
      if (resp.goal) addSuccess(values.name);
    } catch (error) {
      const { status, message } = error as AxiosError;
      notify({
        messageType: 'error',
        message:
          status === 409
            ? 'The Goal ID already exists.'
            : message || 'Something went wrong.'
      });
    }
  };

  const {
    formState: { isDirty, isSubmitting }
  } = form;

  return (
    <SlideModal title={t('new-goal')} isOpen={isOpen} onClose={onClose}>
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
            <Form.Field
              control={form.control}
              name="id"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required className="relative w-fit">
                    {t('form:goal-id')}
                    <Icon
                      icon={IconInfo}
                      size="xs"
                      color="gray-500"
                      className="absolute -right-6"
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:placeholder-goal-id')}`}
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
                          className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
                        >
                          {t('form:experiments')}
                        </label>
                      </div>

                      <div className="flex items-center gap-x-2">
                        <RadioGroupItem
                          id="operations-connection"
                          value="OPERATION"
                        />
                        <label
                          htmlFor="operations-connection"
                          className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
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
            <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
              <ButtonBar
                primaryButton={
                  <Button variant="secondary" onClick={onClose}>
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
