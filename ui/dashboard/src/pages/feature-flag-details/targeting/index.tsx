import { useCallback, useMemo } from 'react';
import { FormProvider, useFieldArray, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryFeatures } from '@queries/features';
import { getCurrentEnvironment, useAuth } from 'auth';
import { cloneDeep } from 'lodash';
import { Feature } from '@types';
import Divider from 'components/divider';
import Form from 'components/form';
import PageLayout from 'elements/page-layout';
import AddRule from './add-rule';
import AudienceTraffic from './audience-traffic';
import { initialIndividualRule, initialPrerequisite } from './constants';
import FlagOffDescription from './flag-off-description';
import FlagSwitch from './flag-switch';
import { formSchema, TargetingSchema } from './form-schema';
import IndividualRule from './individual-rule';
import PrerequisiteRule from './prerequisite-rule';
import { RuleCategory } from './types';

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
      individualRules: [],
      defaultStrategy: {},
      enabled: feature.enabled,
      isShowRules: feature.enabled
    }
  });

  const { control, watch, setValue } = form;

  const {
    fields: prerequisites,
    append: prerequisiteAppend,
    remove: prerequisiteRemove
  } = useFieldArray({
    control,
    name: 'prerequisites'
  });

  const {
    fields: individualRules,
    append: individualAppend,
    remove: individualRemove
  } = useFieldArray({
    control,
    name: 'individualRules'
  });

  const onAddRule = useCallback(
    (rule: RuleCategory) => {
      if (rule === RuleCategory.PREREQUISITE) {
        prerequisiteAppend(cloneDeep(initialPrerequisite));
      }
      if (rule === RuleCategory.INDIVIDUAL) {
        individualAppend(cloneDeep(initialIndividualRule));
      }
    },
    [prerequisites]
  );

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

          {prerequisites?.length > 0 && (
            <>
              <PrerequisiteRule
                features={features}
                feature={feature}
                prerequisites={prerequisites}
                onRemovePrerequisite={prerequisiteRemove}
                onAddPrerequisite={() => onAddRule(RuleCategory.PREREQUISITE)}
              />
              <TargetingDivider />
            </>
          )}
          <AddRule onAddRule={onAddRule} />
          {individualRules?.length > 0 && (
            <>
              <IndividualRule individualRules={individualRules} />
              <TargetingDivider />
            </>
          )}
        </Form>
      </FormProvider>
    </PageLayout.Content>
  );
};

export default TargetingPage;
