import { useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { debuggerEvaluate } from '@api/debugger';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { Evaluation, Feature } from '@types';
import AddDebuggerForm from 'pages/debugger/add-debugger-form';
import {
  addDebuggerFormSchema,
  AddDebuggerFormType
} from 'pages/debugger/form-schema';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import SlideModal from 'components/modal/slide';

interface Props {
  isOpen: boolean;
  feature: Feature;
  evaluations: Evaluation[];
  onShowResults: () => void;
  setEvaluations: (value: Evaluation[]) => void;
  onClose: () => void;
}

const CreateDebuggerForm = ({
  isOpen,
  evaluations,
  feature,
  onShowResults,
  setEvaluations,
  onClose
}: Props) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { errorNotify } = useToast();

  const defaultValues = useMemo(
    () => ({
      flags: [''],
      userIds: [],
      attributes: [
        {
          key: '',
          value: ''
        }
      ]
    }),
    []
  );

  const form = useForm({
    resolver: yupResolver(addDebuggerFormSchema),
    defaultValues: {
      ...defaultValues
    },
    mode: 'onChange'
  });

  const {
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  const onSubmit = useCallback(
    async (values: AddDebuggerFormType) => {
      try {
        const dataMap = new Map();
        values?.attributes?.forEach(item => dataMap.set(item.key, item.value));

        const userData: { [key: string]: string } = {};
        dataMap?.forEach((value, key) => (userData[key] = value));

        const resp = await debuggerEvaluate({
          environmentId: currentEnvironment.id,
          featureIds: values.flags,
          users: values.userIds.map(item => ({
            id: item,
            data: userData
          }))
        });
        setEvaluations(resp.evaluations);
        onShowResults();
        onClose();
      } catch (error) {
        errorNotify(error);
      }
    },
    [currentEnvironment]
  );

  return (
    <SlideModal
      title={t(`navigation.debugger`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <AddDebuggerForm
            isOnTargeting
            isLoading={form.formState.isSubmitting}
            feature={feature}
            evaluations={evaluations}
            onCancel={onShowResults}
          />
          <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
            <ButtonBar
              primaryButton={
                <Button
                  variant={'secondary'}
                  className="w-fit"
                  onClick={onClose}
                >
                  {t('cancel')}
                </Button>
              }
              secondaryButton={
                <Button
                  className="w-fit"
                  loading={isSubmitting}
                  disabled={!isDirty || !isValid}
                >
                  {t('evaluate')}
                </Button>
              }
            />
          </div>
        </Form>
      </FormProvider>
    </SlideModal>
  );
};

export default CreateDebuggerForm;
