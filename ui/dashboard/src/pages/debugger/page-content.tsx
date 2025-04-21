import { useCallback, useMemo, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { debuggerEvaluate } from '@api/debugger';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { Evaluation } from '@types';
import Form from 'components/form';
import PageLayout from 'elements/page-layout';
import AddDebuggerForm from './add-debugger-form';
import DebuggerResults from './debugger-results';
import { addDebuggerFormSchema, AddDebuggerFormType } from './form-schema';

const PageContent = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [isShowResults, setIsShowResults] = useState(false);
  const [evaluations, setEvaluations] = useState<Evaluation[]>([]);

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
        setIsShowResults(true);
      } catch (error) {
        errorNotify(error);
      }
    },
    [currentEnvironment]
  );

  return (
    <PageLayout.Content>
      {!isShowResults ? (
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <AddDebuggerForm isLoading={form.formState.isSubmitting} />
          </Form>
        </FormProvider>
      ) : (
        <DebuggerResults
          evaluations={evaluations}
          onEditFields={() => setIsShowResults(false)}
          onResetFields={() => {
            setIsShowResults(false);
            form.reset({
              ...defaultValues
            });
          }}
        />
      )}
    </PageLayout.Content>
  );
};

export default PageContent;
