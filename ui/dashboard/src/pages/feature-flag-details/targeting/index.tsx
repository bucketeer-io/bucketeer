import { useCallback, useMemo, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryFeatures } from '@queries/features';
import { getCurrentEnvironment, useAuth } from 'auth';
import { v4 as uuid } from 'uuid';
import { Feature, FeatureRuleClauseOperator, StrategyType } from '@types';
import Form from 'components/form';
import AddRuleButton from './add-rule-button';
import AddRuleDropdown from './add-rule-dropdown';
import { initialPrerequisitesRule } from './constants';
import DefaultRule from './default-rule';
import { formSchema, RuleClauseType, RuleSchema } from './form-schema';
import IndividualRule from './individual-rule';
import PrerequisiteRule from './prerequisite-rule';
import TargetSegmentRule from './target-segment-rule';
import TargetingState from './targeting-state';
import {
  IndividualRuleItem,
  RuleCategory,
  TargetingForm,
  TargetPrerequisiteItem
} from './types';

const Targeting = ({ feature }: { feature: Feature }) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [targetSegmentRules, setTargetSegmentRules] = useState<RuleSchema[]>(
    []
  );
  const [targetIndividualRules, setTargetIndividualRules] = useState<
    IndividualRuleItem[]
  >([]);
  const [prerequisitesRules, setPrerequisitesRules] = useState<
    TargetPrerequisiteItem[]
  >([]);

  const { data: collection } = useQueryFeatures({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
    },
    enabled: !!currentEnvironment?.id
  });

  const features = useMemo(() => collection?.features || [], [collection]);

  const form = useForm({
    resolver: yupResolver(formSchema),
    values: {
      prerequisitesRules,
      rules: targetSegmentRules,
      targetIndividualRules
    }
  });

  const defaultRolloutStrategy = useMemo(
    () =>
      feature.variations?.map(val => ({
        variation: val.id,
        weight: 0
      })),
    [feature]
  );

  const defaultRule = useMemo(() => {
    return {
      id: uuid(),
      strategy: {
        currentOption: feature?.variations[0]?.id || '',
        fixedStrategy: {
          variation: feature?.variations[0]?.id || ''
        },
        rolloutStrategy: defaultRolloutStrategy,
        type: StrategyType.FIXED
      },
      clauses: [
        {
          id: uuid(),
          type: RuleClauseType.COMPARE,
          attribute: '',
          operator: FeatureRuleClauseOperator.EQUALS,
          values: []
        }
      ]
    };
  }, [feature, defaultRolloutStrategy]);

  const onAddRule = useCallback(
    (type: RuleCategory) => {
      if (type === 'target-segments') {
        const _rules = [...targetSegmentRules, defaultRule];
        form.setValue('rules', _rules);
        return setTargetSegmentRules(_rules);
      }
      if (type === 'target-individuals') {
        const data = feature?.variations?.map(({ name, id }) => ({
          variationId: id,
          name,
          users: []
        }));
        return setTargetIndividualRules(data);
      }
      setPrerequisitesRules([
        ...prerequisitesRules,
        {
          index: prerequisitesRules.length + 1,
          rules: [initialPrerequisitesRule]
        }
      ]);
    },
    [
      targetSegmentRules,
      targetIndividualRules,
      prerequisitesRules,
      feature,
      defaultRule,
      form
    ]
  );

  const onSubmit: SubmitHandler<TargetingForm> = async values => {
    console.log(values);
  };

  return (
    <>
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col size-full gap-y-6 overflow-visible">
            <TargetingState />
            <AddRuleDropdown onAddRule={onAddRule} />
            {prerequisitesRules.length > 0 && (
              <Form.Field
                control={form.control}
                name={'prerequisitesRules'}
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Control>
                      <PrerequisiteRule
                        features={features}
                        feature={feature}
                        prerequisitesRules={prerequisitesRules}
                        onChangePrerequisitesRules={rules => {
                          field.onChange(rules);
                          setPrerequisitesRules(rules);
                        }}
                      />
                    </Form.Control>
                  </Form.Item>
                )}
              />
            )}
            {targetIndividualRules.length > 0 && (
              <Form.Field
                control={form.control}
                name={'targetIndividualRules'}
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Control>
                      <IndividualRule
                        individualRules={targetIndividualRules}
                        onChangeIndividualRules={rules => {
                          field.onChange(rules);
                          setTargetIndividualRules(rules);
                        }}
                      />
                    </Form.Control>
                  </Form.Item>
                )}
              />
            )}
            {targetSegmentRules.length > 0 && (
              <>
                <Form.Field
                  control={form.control}
                  name={'rules'}
                  render={({ field }) => (
                    <Form.Item className="py-0">
                      <Form.Control>
                        <TargetSegmentRule
                          feature={feature}
                          defaultRolloutStrategy={defaultRolloutStrategy}
                          features={features}
                          targetSegmentRules={targetSegmentRules}
                          onChangeTargetSegmentRules={rules => {
                            field.onChange(rules);
                            setTargetSegmentRules(rules);
                          }}
                          onAddRule={() => onAddRule('target-segments')}
                        />
                      </Form.Control>
                    </Form.Item>
                  )}
                />
                <AddRuleButton
                  isCenter
                  onAddRule={() => onAddRule('target-segments')}
                />
              </>
            )}
            <DefaultRule />
          </div>
        </Form>
      </FormProvider>
    </>
  );
};

export default Targeting;
