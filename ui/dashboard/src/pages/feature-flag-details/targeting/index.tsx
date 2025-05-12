import { useCallback, useMemo } from 'react';
import { FormProvider, useFieldArray, useForm } from 'react-hook-form';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures, useQueryFeatures } from '@queries/features';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { Feature } from '@types';
import { IconDebugger } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import Icon from 'components/icon';
import PageLayout from 'elements/page-layout';
import AddRule from './add-rule';
import AudienceTraffic from './audience-traffic';
import { initialPrerequisite } from './constants';
import DefaultRule from './default-rule';
import FlagOffDescription from './flag-off-description';
import FlagSwitch from './flag-switch';
import { formSchema, TargetingSchema } from './form-schema';
import IndividualRule from './individual-rule';
import PrerequisiteRule from './prerequisite-rule';
import TargetSegmentRule from './segment-rule';
import { RuleCategory } from './types';
import {
  getDefaultRule,
  handleCheckIndividualRules,
  handleCheckPrerequisites,
  handleCheckSegmentRules,
  handleCreateDefaultValues
} from './utils';

const TargetingDivider = () => (
  <Divider vertical className="!h-6 w-px self-center my-4" />
);

const TargetingPage = ({ feature }: { feature: Feature }) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { t } = useTranslation(['common']);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const { data: collection } = useQueryFeatures({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
    },
    enabled: !!currentEnvironment?.id
  });

  const features = useMemo(() => collection?.features || [], [collection]);
  console.log({ feature });

  const form = useForm<TargetingSchema>({
    resolver: yupResolver(formSchema),
    defaultValues: handleCreateDefaultValues(feature)
  });

  const {
    control,
    formState: { isDirty, isValid },
    watch
  } = form;

  const isShowRules = watch('isShowRules');

  const {
    fields: prerequisites,
    append: prerequisiteAppend,
    remove: prerequisiteRemove
  } = useFieldArray({
    control,
    name: 'prerequisites'
  });

  const { fields: individualRules, append: individualAppend } = useFieldArray({
    control,
    name: 'individualRules'
  });

  const {
    fields: segmentRules,
    append: segmentRulesAppend,
    remove: segmentRulesRemove,
    swap: segmentRulesSwap
  } = useFieldArray({
    control,
    name: 'segmentRules',
    keyName: 'segmentId'
  });

  const onAddRule = useCallback(
    (rule: RuleCategory) => {
      if (rule === RuleCategory.PREREQUISITE) {
        return prerequisiteAppend(cloneDeep(initialPrerequisite));
      }
      if (rule === RuleCategory.INDIVIDUAL) {
        return individualAppend(
          feature.variations.map(item => ({
            variationId: item.id,
            name: item.name,
            users: []
          }))
        );
      }
      segmentRulesAppend(getDefaultRule(feature));
    },
    [prerequisites, feature]
  );

  const onSubmit = useCallback(
    async (values: TargetingSchema) => {
      try {
        const { enabled, individualRules, segmentRules, prerequisites } =
          values;
        const {
          id,
          rules,
          targets,
          prerequisites: featurePrerequisites
        } = feature;
        const resp = await featureUpdater({
          id,
          environmentId: currentEnvironment.id,
          enabled,
          ruleChanges: handleCheckSegmentRules(rules, segmentRules),
          targetChanges: handleCheckIndividualRules(targets, individualRules),
          prerequisiteChanges: handleCheckPrerequisites(
            featurePrerequisites,
            prerequisites
          ),
          comment: 'test'
        });
        if (resp) {
          notify({
            message: t('message:flag-updated')
          });
          invalidateFeature(queryClient);
          invalidateFeatures(queryClient);
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [feature, currentEnvironment]
  );

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
          {!feature.enabled && <FlagOffDescription />}
          {isShowRules && (
            <>
              {prerequisites?.length > 0 && (
                <>
                  <PrerequisiteRule
                    features={features}
                    feature={feature}
                    prerequisites={prerequisites}
                    onRemovePrerequisite={prerequisiteRemove}
                    onAddPrerequisite={() =>
                      onAddRule(RuleCategory.PREREQUISITE)
                    }
                  />
                  <TargetingDivider />
                </>
              )}
              <AddRule
                individualRules={individualRules}
                onAddRule={onAddRule}
              />
              <TargetingDivider />
              {individualRules?.length > 0 && (
                <>
                  <IndividualRule individualRules={individualRules} />
                  <TargetingDivider />
                  <AddRule
                    individualRules={individualRules}
                    onAddRule={onAddRule}
                  />
                </>
              )}
              {segmentRules.length > 0 && (
                <>
                  <TargetSegmentRule
                    feature={feature}
                    features={features}
                    segmentRules={segmentRules}
                    segmentRulesRemove={segmentRulesRemove}
                    segmentRulesSwap={segmentRulesSwap}
                  />
                  <TargetingDivider />
                  <AddRule
                    individualRules={individualRules}
                    onAddRule={onAddRule}
                  />
                </>
              )}
            </>
          )}
          <TargetingDivider />
          <DefaultRule feature={feature} />
          <ButtonBar
            primaryButton={
              <Button type="button" variant={'secondary-2'} className="size-12">
                <Icon icon={IconDebugger} color="gray-500" />
              </Button>
            }
            secondaryButton={
              <Button disabled={!isDirty || !isValid}>
                {t('save-with-comment')}
              </Button>
            }
          />
        </Form>
      </FormProvider>
    </PageLayout.Content>
  );
};

export default TargetingPage;
