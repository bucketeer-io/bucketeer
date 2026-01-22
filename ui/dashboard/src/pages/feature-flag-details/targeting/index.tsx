import { useCallback, useMemo, useState } from 'react';
import { FormProvider, useFieldArray, useForm } from 'react-hook-form';
import { autoOpsCreator } from '@api/auto-ops';
import { FeatureResponse, featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures, useQueryFeatures } from '@queries/features';
import { useQueryRollouts } from '@queries/rollouts';
import { invalidateUserSegments } from '@queries/user-segments';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import cloneDeep from 'lodash/cloneDeep';
import { v4 as uuid } from 'uuid';
import { Evaluation, Feature } from '@types';
import { IconDebugger } from '@icons';
import { AddDebuggerFormType } from 'pages/debugger/form-schema';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import PageLayout from 'elements/page-layout';
import ConfirmationRequiredModal, {
  ConfirmRequiredValues
} from '../elements/confirm-required-modal';
import AddRule from './add-rule';
import AudienceTraffic from './audience-traffic';
import { initialPrerequisite } from './constants';
import CreateDebuggerForm from './debugger/create-form';
import TargetingDebuggerResults from './debugger/results';
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

export const TargetingDivider = () => (
  <Divider vertical className="!h-6 w-px self-center my-4 !border-gray-400" />
);

const TargetingPage = ({
  feature,
  editable
}: {
  feature: Feature;
  editable: boolean;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const [isOpenConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const [isOpenDebuggerModal, onOpenDebuggerModal, onCloseDebuggerModal] =
    useToggleOpen(false);

  const [evaluations, setEvaluations] = useState<Evaluation[]>([]);
  const [isOpenResults, onOpenResultsModal, onCloseResultsModal] =
    useToggleOpen(false);
  const [debuggerForm, setDebuggerForm] = useState<AddDebuggerFormType | null>(
    null
  );
  const [isShowRules, setIsShowRules] = useState<boolean>(feature.enabled);

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
    enabled: !!currentEnvironment
  });
  const waitingRunningRollouts =
    rolloutCollection?.progressiveRollouts?.filter(item =>
      ['WAITING', 'RUNNING'].includes(item.status)
    ) || [];
  const features = useMemo(() => collection?.features || [], [collection]);
  const activeFeatures = useMemo(
    () => features.filter(item => !item.archived) || [],
    [features]
  );

  const form = useForm<TargetingSchema>({
    resolver: yupResolver(formSchema),
    defaultValues: handleCreateDefaultValues(feature),
    mode: 'onChange'
  });

  const {
    control,
    formState: { isDirty, isValid, dirtyFields, isSubmitting },
    watch,
    reset
  } = form;

  useUnsavedLeavePage({
    isShow: isDirty && !isSubmitting
  });

  const enabledWatch = watch('enabled');
  const prerequisitesWatch = [...(watch('prerequisites') || [])];
  const segmentRulesWatch = [...(watch('segmentRules') || [])];

  const hasPrerequisiteFlags = activeFeatures.filter(item =>
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
    const filterFeatures = activeFeatures.filter(f => f.id !== feature.id);
    const featuresSelected = prerequisitesWatch.map(
      (item: PrerequisiteSchema) => item.featureId
    );

    return (
      !filterFeatures.filter(f => !featuresSelected.includes(f.id))?.length ||
      filterFeatures?.length === prerequisitesWatch?.length
    );
  }, [prerequisitesWatch, activeFeatures, feature]);

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
    update: segmentRulesUpdate,
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
    [feature]
  );

  const handleSwapSegmentRule = useCallback(
    (indexA: number, indexB: number) => {
      segmentRulesUpdate(indexA, {
        ...segmentRulesWatch[indexA],
        id: uuid()
      });
      segmentRulesUpdate(indexB, {
        ...segmentRulesWatch[indexB],
        id: uuid()
      });
      segmentRulesSwap(indexA, indexB);
    },
    [segmentRulesWatch]
  );

  const onSubmit = useCallback(
    async (
      values: TargetingSchema,
      additionalValues?: ConfirmRequiredValues
    ) => {
      if (editable) {
        try {
          const {
            enabled,
            individualRules,
            segmentRules,
            prerequisites,
            defaultRule,
            offVariation
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
              targetChanges: handleCheckIndividualRules(
                targets,
                individualRules
              ),
              prerequisiteChanges: handleCheckPrerequisites(
                featurePrerequisites,
                prerequisites
              ),
              comment,
              resetSamplingSeed: resetSampling,
              offVariation
            });
          }
          if (resp) {
            notify({
              message: t('message:collection-action-success', {
                collection: t('source-type.feature-flag'),
                action: t('updated')
              })
            });
            invalidateFeature(queryClient);
            invalidateFeatures(queryClient);
            invalidateUserSegments(queryClient);
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
      }
    },
    [feature, currentEnvironment, isShowUpdateSchedule, editable]
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
          <FlagSwitch
            feature={feature}
            setIsShowRules={setIsShowRules}
            editable={editable}
          />
          {(!feature.enabled || !enabledWatch) && (
            <>
              <TargetingDivider />
              <FlagOffDescription
                isShowRules={isShowRules}
                setIsShowRules={setIsShowRules}
              />
            </>
          )}
          {isShowRules && (
            <>
              {(prerequisites?.length > 0 ||
                hasPrerequisiteFlags?.length > 0) && (
                <>
                  <TargetingDivider />
                  <PrerequisiteRule
                    currentEnvironment={currentEnvironment}
                    isDisableAddPrerequisite={isDisableAddPrerequisite}
                    features={activeFeatures}
                    feature={feature}
                    prerequisites={prerequisites}
                    hasPrerequisiteFlags={hasPrerequisiteFlags}
                    onRemovePrerequisite={prerequisiteRemove}
                    onAddPrerequisite={() =>
                      onAddRule(RuleCategory.PREREQUISITE)
                    }
                  />
                </>
              )}
              {(!prerequisitesWatch?.length || !individualRules?.length) && (
                <>
                  <TargetingDivider />
                  <AddRule
                    isDisableAddPrerequisite={prerequisitesWatch?.length > 0}
                    isDisableAddIndividualRules={individualRules?.length > 0}
                    onAddRule={onAddRule}
                  />
                </>
              )}
              {individualRules?.length > 0 && (
                <>
                  <TargetingDivider />
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
                  <TargetingDivider />
                  <TargetSegmentRule
                    feature={feature}
                    features={activeFeatures}
                    segmentRules={segmentRules}
                    isDisableAddPrerequisite={prerequisitesWatch?.length > 0}
                    isDisableAddIndividualRules={individualRules?.length > 0}
                    onAddRule={onAddRule}
                    segmentRulesRemove={segmentRulesRemove}
                    segmentRulesSwap={handleSwapSegmentRule}
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
            editable={editable}
            urlCode={currentEnvironment.urlCode}
            feature={feature}
            waitingRunningRollouts={waitingRunningRollouts}
          />
          <ButtonBar
            primaryButton={
              <Tooltip
                side="top"
                className="max-w-[320px]"
                content={t('form:targeting.tooltip.debugger')}
                trigger={
                  <Button
                    type="button"
                    variant={'secondary-2'}
                    className="size-12"
                    onClick={onOpenDebuggerModal}
                  >
                    <Icon icon={IconDebugger} color="gray-500" />
                  </Button>
                }
              />
            }
            secondaryButton={
              <DisabledButtonTooltip
                hidden={editable}
                trigger={
                  <Button
                    type="button"
                    disabled={!isDirty || !isValid || !editable}
                    onClick={onOpenConfirmModal}
                  >
                    {t('save-with-comment')}
                  </Button>
                }
              />
            }
          />
        </Form>
      </FormProvider>
      {isOpenConfirmModal && (
        <ConfirmationRequiredModal
          feature={feature}
          isOpen={isOpenConfirmModal}
          isShowScheduleSelect={isShowUpdateSchedule}
          onClose={onCloseConfirmModal}
          onSubmit={additionalValues =>
            form.handleSubmit(values => onSubmit(values, additionalValues))()
          }
        />
      )}
      {isOpenDebuggerModal && (
        <CreateDebuggerForm
          isOpen={isOpenDebuggerModal}
          feature={feature}
          evaluations={evaluations}
          debuggerForm={debuggerForm}
          onClose={onCloseDebuggerModal}
          setEvaluations={setEvaluations}
          setDebuggerForm={setDebuggerForm}
          onShowResults={onOpenResultsModal}
        />
      )}
      {isOpenResults && (
        <TargetingDebuggerResults
          isOpen={isOpenResults}
          features={activeFeatures}
          evaluations={evaluations}
          onClose={() => {
            setEvaluations([]);
            onCloseResultsModal();
            setDebuggerForm(null);
          }}
          onEditFields={() => {
            onCloseResultsModal();
            setEvaluations([]);
            onOpenDebuggerModal();
          }}
        />
      )}
    </PageLayout.Content>
  );
};

export default TargetingPage;
