import { useCallback, useMemo } from 'react';
import { FormProvider, useFieldArray, useForm } from 'react-hook-form';
import { autoOpsCreator } from '@api/auto-ops';
import { FeatureResponse, featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures, useQueryFeatures } from '@queries/features';
import { useQueryRollouts } from '@queries/rollouts';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
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
import ConfirmationRequiredModal, {
  ConfirmRequiredValues
} from '../elements/confirm-required-modal';
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
import { PrerequisiteSchema, RuleCategory } from './types';
import {
  getDefaultRule,
  handleCheckIndividualRules,
  handleCheckPrerequisites,
  handleCheckSegmentRules,
  handleCreateDefaultValues,
  handleGetDefaultRuleStrategy
} from './utils';

const TargetingDivider = () => (
  <Divider vertical className="!h-6 w-px self-center my-4" />
);

const TargetingPage = ({ feature }: { feature: Feature }) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { t } = useTranslation(['common', 'form']);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const [isOpenConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const { data: rolloutCollection } = useQueryRollouts({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      featureIds: [feature.id]
    }
  });

  const { data: collection } = useQueryFeatures({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      archived: false
    },
    enabled: !!currentEnvironment?.id
  });
  const waitingRunningRollouts =
    rolloutCollection?.progressiveRollouts?.filter(item =>
      ['WAITING', 'RUNNING'].includes(item.status)
    ) || [];
  const features = useMemo(() => collection?.features || [], [collection]);

  const form = useForm<TargetingSchema>({
    resolver: yupResolver(formSchema),
    defaultValues: handleCreateDefaultValues(feature)
  });

  const {
    control,
    formState: { isDirty, isValid, dirtyFields },
    watch,
    reset
  } = form;

  const isShowRules = watch('isShowRules');

  const prerequisitesWatch = [...(watch('prerequisites') || [])];

  const hasPrerequisiteFlags = features.filter(item =>
    item.prerequisites.find(p => p.featureId === feature.id)
  );

  const isShowUpdateSchedule =
    dirtyFields?.enabled &&
    Object.keys(dirtyFields).filter(
      key =>
        !['requireComment', 'comment', 'scheduleType', 'scheduleAt'].includes(
          key
        )
    ).length <= 1;

  const isDisableAddPrerequisite = useMemo(() => {
    if (!features?.length) return true;
    const filterFeatures = features.filter(f => f.id !== feature.id);
    const featuresSelected = prerequisitesWatch.map(
      (item: PrerequisiteSchema) => item.featureId
    );

    return (
      !filterFeatures.filter(f => !featuresSelected.includes(f.id))?.length ||
      filterFeatures?.length === prerequisitesWatch?.length
    );
  }, [prerequisitesWatch, features, feature]);

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
    async (
      values: TargetingSchema,
      additionalValues?: ConfirmRequiredValues
    ) => {
      try {
        const {
          enabled,
          individualRules,
          segmentRules,
          prerequisites,
          defaultRule
        } = values;

        const { comment, resetSampling, scheduleType, scheduleAt } =
          additionalValues || {};

        const {
          id,
          rules,
          targets,
          prerequisites: featurePrerequisites
        } = feature;
        let resp;
        const isScheduleUpdate =
          isShowUpdateSchedule &&
          !['ENABLE', 'DISABLE'].includes(scheduleType as string);
        if (isScheduleUpdate) {
          resp = await autoOpsCreator({
            environmentId: currentEnvironment.id,
            featureId: feature.id,
            opsType: 'SCHEDULE',
            datetimeClauses: [
              {
                actionType: enabled ? 'ENABLE' : 'DISABLE',
                time: scheduleAt as string
              }
            ]
          });
        } else {
          resp = await featureUpdater({
            id,
            environmentId: currentEnvironment.id,
            enabled: isScheduleUpdate ? false : enabled,
            defaultStrategy: handleGetDefaultRuleStrategy(defaultRule),
            ruleChanges: handleCheckSegmentRules(rules, segmentRules),
            targetChanges: handleCheckIndividualRules(targets, individualRules),
            prerequisiteChanges: handleCheckPrerequisites(
              featurePrerequisites,
              prerequisites
            ),
            comment,
            resetSamplingSeed: resetSampling
          });
        }
        if (resp) {
          notify({
            message: t('message:flag-updated')
          });
          invalidateFeature(queryClient);
          invalidateFeatures(queryClient);
          reset(
            handleCreateDefaultValues(
              isScheduleUpdate ? feature : (resp as FeatureResponse)?.feature
            )
          );
          onCloseConfirmModal();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [feature, currentEnvironment, isShowUpdateSchedule]
  );

  return (
    <PageLayout.Content className="p-6 pt-0 gap-y-6 min-w-[900px]">
      <FormProvider {...form}>
        <Form
          onSubmit={form.handleSubmit(values => onSubmit(values))}
          className="flex flex-col w-full items-center"
        >
          <AudienceTraffic />
          <TargetingDivider />
          <FlagSwitch />
          <TargetingDivider />
          {!feature.enabled && <FlagOffDescription />}
          {isShowRules && (
            <>
              {(prerequisites?.length > 0 ||
                hasPrerequisiteFlags?.length > 0) && (
                <>
                  <PrerequisiteRule
                    isDisableAddPrerequisite={isDisableAddPrerequisite}
                    features={features}
                    feature={feature}
                    prerequisites={prerequisites}
                    hasPrerequisiteFlags={hasPrerequisiteFlags}
                    onRemovePrerequisite={prerequisiteRemove}
                    onAddPrerequisite={() =>
                      onAddRule(RuleCategory.PREREQUISITE)
                    }
                  />
                  <TargetingDivider />
                </>
              )}
              <AddRule
                isDisableAddPrerequisite={prerequisitesWatch?.length > 0}
                isDisableAddIndividualRules={individualRules?.length > 0}
                onAddRule={onAddRule}
              />
              <TargetingDivider />
              {individualRules?.length > 0 && (
                <>
                  <IndividualRule individualRules={individualRules} />
                  <TargetingDivider />
                  <AddRule
                    isDisableAddPrerequisite={prerequisitesWatch?.length > 0}
                    isDisableAddIndividualRules={individualRules?.length > 0}
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
                    isDisableAddPrerequisite={prerequisitesWatch?.length > 0}
                    isDisableAddIndividualRules={individualRules?.length > 0}
                    onAddRule={onAddRule}
                  />
                </>
              )}
            </>
          )}
          <TargetingDivider />
          <DefaultRule
            urlCode={currentEnvironment.urlCode}
            feature={feature}
            waitingRunningRollouts={waitingRunningRollouts}
          />
          <ButtonBar
            primaryButton={
              <Button type="button" variant={'secondary-2'} className="size-12">
                <Icon icon={IconDebugger} color="gray-500" />
              </Button>
            }
            secondaryButton={
              <Button
                type="button"
                disabled={!isDirty || !isValid}
                onClick={onOpenConfirmModal}
              >
                {t('save-with-comment')}
              </Button>
            }
          />
        </Form>
      </FormProvider>
      {isOpenConfirmModal && (
        <ConfirmationRequiredModal
          feature={feature}
          isOpen={isOpenConfirmModal}
          isShowScheduleSelect={isShowUpdateSchedule}
          isShowRolloutWarning={true}
          onClose={onCloseConfirmModal}
          onSubmit={additionalValues =>
            form.handleSubmit(values => onSubmit(values, additionalValues))()
          }
        />
      )}
    </PageLayout.Content>
  );
};

export default TargetingPage;
