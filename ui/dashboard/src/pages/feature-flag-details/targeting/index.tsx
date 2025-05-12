import { useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryFeatures } from '@queries/features';
import { getCurrentEnvironment, useAuth } from 'auth';
import { Feature } from '@types';
import Divider from 'components/divider';
import Form from 'components/form';
import PageLayout from 'elements/page-layout';
import AddRule from './add-rule';
import AudienceTraffic from './audience-traffic';
import FlagOffDescription from './flag-off-description';
import FlagSwitch from './flag-switch';
import { formSchema, TargetingSchema } from './form-schema';

const TargetingDivider = () => (
  <Divider vertical className="!h-6 w-px self-center my-4" />
);

const TargetingPage = ({ feature }: { feature: Feature }) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { data: collection } = useQueryFeatures({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
    },
    enabled: !!currentEnvironment?.id
  });

  const features = useMemo(() => collection?.features || [], [collection]);

  const form = useForm<TargetingSchema>({
    resolver: yupResolver(formSchema),
    defaultValues: {
      prerequisites: [],
      rules: [],
      targetIndividualRules: [],
      defaultStrategy: {},
      enabled: feature.enabled,
      isShowRules: feature.enabled
    }
  });
  console.log(features);
  const onSubmit = useCallback(async (values: TargetingSchema) => {
    console.log(values);
  }, []);

  return (
    <PageLayout.Content className="p-6 pt-0 gap-y-6 min-w-[900px]">
      <FormProvider {...form}>
        <Form
          onSubmit={form.handleSubmit(onSubmit)}
          className="flex flex-col w-full items-center"
        >
          <AudienceTraffic />
          <TargetingDivider />
          <FlagSwitch />
          <TargetingDivider />
          {!feature.enabled && (
            <>
              <FlagOffDescription />
              <TargetingDivider />
            </>
          )}
          <AddRule />
        </Form>
      </FormProvider>
    </PageLayout.Content>
  );
};

export default TargetingPage;
