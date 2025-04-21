import { useCallback, useMemo, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import Form from 'components/form';
import PageLayout from 'elements/page-layout';
import AddDebuggerForm from './add-debugger-form';
import DebuggerResults from './debugger-results';
import { addDebuggerFormSchema, AddDebuggerFormType } from './form-schema';

const PageContent = () => {
  const [isShowResults, setIsShowResults] = useState(false);

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

  const onSubmit = useCallback((values: AddDebuggerFormType) => {
    console.log(values);
    setIsShowResults(true);
  }, []);

  return (
    <PageLayout.Content className="p-6">
      {!isShowResults ? (
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <AddDebuggerForm />
          </Form>
        </FormProvider>
      ) : (
        <DebuggerResults
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
