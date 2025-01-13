import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Goal } from '@types';
import Button from 'components/button';
import Form from 'components/form';
import Input from 'components/input';
import TextArea from 'components/textarea';

export interface GoalDetailsForm {
  name: string;
  id: string;
  description?: string;
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  id: yup.string().required(),
  description: yup.string()
});

const GoalUpdateForm = ({ goal }: { goal: Goal }) => {
  const { t } = useTranslation(['common', 'form']);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: goal.name,
      id: goal.id,
      description: goal.description
    }
  });

  const onSubmit: SubmitHandler<GoalDetailsForm> = async values => {
    console.log(values);
  };

  return (
    <div className="p-5 shadow-card rounded-lg bg-white">
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
                <Form.Label required>{t('form:goal-id')}</Form.Label>
                <Form.Control>
                  <Input
                    disabled
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

          <Button
            loading={form.formState.isSubmitting}
            disabled={!form.formState.isDirty}
            type="submit"
            className="w-fit"
            variant={'secondary'}
          >
            {t(`save-with-comment`)}
          </Button>
        </Form>
      </FormProvider>
    </div>
  );
};

export default GoalUpdateForm;
