import { useCallback } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { goalCreator } from '@api/goal';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateGoals } from '@queries/goals';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Goal } from '@types';
import { onGenerateSlug } from 'utils/converts';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import Input from 'components/input';
import DialogModal from 'components/modal/dialog';
import TextArea from 'components/textarea';

export type CreateGoalModalProps = {
  isOpen: boolean;
  onClose: () => void;
  onCompleted?: (goal: Goal) => void;
};

interface CreateGoalForm {
  id: string;
  name: string;
  description?: string;
}

const formSchema = yup.object().shape({
  id: yup.string().required('This field is required'),
  name: yup.string().required('This field is required'),
  description: yup.string()
});

const CreateGoalModal = ({
  isOpen,
  onClose,
  onCompleted
}: CreateGoalModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: '',
      name: '',
      description: ''
    }
  });

  const {
    control,
    formState: { isSubmitting },
    setValue
  } = form;

  const onSubmit: SubmitHandler<CreateGoalForm> = useCallback(async values => {
    try {
      const { name, id, description } = values;
      const resp = await goalCreator({
        connectionType: 'EXPERIMENT',
        environmentId: currentEnvironment.id,
        name,
        id,
        description
      });

      if (resp) {
        notify({
          message: (
            <span>
              <b>{name}</b> {`has been successfully created!`}
            </span>
          )
        });
        onCompleted?.(resp.goal);
        invalidateGoals(queryClient);
        onClose();
        form.reset();
      }
    } catch (error) {
      errorNotify(error);
    }
  }, []);

  return (
    <DialogModal
      className="w-[500px]"
      title={t('new-goal')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col w-full items-start p-5 gap-y-5">
            <p className="text-gray-800 typo-head-bold-small">
              {t('form:general-info')}
            </p>
            <Form.Field
              control={control}
              name="name"
              render={({ field }) => (
                <Form.Item className="w-full py-0">
                  <Form.Label required>{t('name')}</Form.Label>
                  <Form.Control>
                    <Input
                      {...field}
                      placeholder={t('form:placeholder-name')}
                      onChange={value => {
                        field.onChange(value);
                        const id = onGenerateSlug(value);
                        setValue('id', id);
                      }}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={control}
              name="id"
              render={({ field }) => (
                <Form.Item className="w-full py-0">
                  <Form.Label required>{t('form:goal-id')}</Form.Label>
                  <Form.Control>
                    <Input
                      {...field}
                      placeholder={t('form:placeholder-goal-id')}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={control}
              name="description"
              render={({ field }) => (
                <Form.Item className="w-full py-0">
                  <Form.Label optional>{t('form:description')}</Form.Label>
                  <Form.Control>
                    <TextArea
                      {...field}
                      placeholder={t('form:placeholder-desc')}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </div>
          <ButtonBar
            secondaryButton={
              <Button loading={isSubmitting}>{t(`create-goal`)}</Button>
            }
            primaryButton={
              <Button type="button" onClick={onClose} variant="secondary">
                {t(`cancel`)}
              </Button>
            }
          />
        </Form>
      </FormProvider>
    </DialogModal>
  );
};

export default CreateGoalModal;
