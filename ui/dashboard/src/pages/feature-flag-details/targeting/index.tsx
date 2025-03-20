import { useCallback, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { Feature } from '@types';
import Form from 'components/form';
import AddRuleButton from './add-rule-button';
import AddRuleDropdown from './add-rule-dropdown';
import {
  initialIndividualRule,
  initialPrerequisitesRule,
  initialSegmentCondition
} from './constants';
import DefaultRule from './default-rule';
import { formSchema } from './form-schema';
import IndividualRule from './individual-rule';
import PrerequisiteRule from './prerequisite-rule';
import TargetSegmentRule from './target-segment-rule';
import TargetingState from './targeting-state';
import {
  IndividualRuleItem,
  RuleCategory,
  TargetingForm,
  TargetPrerequisiteItem,
  TargetSegmentItem
} from './types';

const Targeting = ({ feature }: { feature: Feature }) => {
  const [targetSegmentRules, setTargetSegmentRules] = useState<
    TargetSegmentItem[]
  >([]);
  const [targetIndividualRules, setTargetIndividualRules] = useState<
    IndividualRuleItem[]
  >([]);
  const [prerequisitesRules, setPrerequisitesRules] = useState<
    TargetPrerequisiteItem[]
  >([]);

  const form = useForm({
    resolver: yupResolver(formSchema),
    values: {
      prerequisitesRules,
      targetSegmentRules,
      targetIndividualRules
    }
  });

  const onAddRule = useCallback(
    (type: RuleCategory) => {
      if (type === 'target-segments') {
        return setTargetSegmentRules([
          ...targetSegmentRules,
          {
            index: targetSegmentRules.length + 1,
            rules: [
              {
                variation: true,
                conditions: [initialSegmentCondition]
              }
            ]
          }
        ]);
      }
      if (type === 'target-individuals') {
        return setTargetIndividualRules([
          ...targetIndividualRules,
          initialIndividualRule
        ]);
      }
      setPrerequisitesRules([
        ...prerequisitesRules,
        {
          index: prerequisitesRules.length + 1,
          rules: [initialPrerequisitesRule]
        }
      ]);
    },
    [targetSegmentRules, targetIndividualRules, prerequisitesRules]
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
              <>
                <Form.Field
                  control={form.control}
                  name={'prerequisitesRules'}
                  render={({ field }) => (
                    <Form.Item className="py-0">
                      <Form.Control>
                        <PrerequisiteRule
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

                <AddRuleButton
                  isCenter
                  onAddRule={() => onAddRule('set-prerequisites')}
                />
              </>
            )}
            {targetIndividualRules.length > 0 && (
              <>
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
                <AddRuleButton
                  isCenter
                  onAddRule={() => onAddRule('target-individuals')}
                />
              </>
            )}
            {targetSegmentRules.length > 0 && (
              <>
                <Form.Field
                  control={form.control}
                  name={'targetSegmentRules'}
                  render={({ field }) => (
                    <Form.Item className="py-0">
                      <Form.Control>
                        <TargetSegmentRule
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
