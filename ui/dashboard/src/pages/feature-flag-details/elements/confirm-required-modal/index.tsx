import { ReactNode, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { isNil } from 'lodash';
import { Feature, FeatureRuleStrategy } from '@types';
import { IconInfo, IconToastWarning, IconWatch } from '@icons';
import { TargetingSchema } from 'pages/feature-flag-details/targeting/form-schema';
import { DiscardChangesStateData } from 'pages/feature-flag-details/targeting/types';
import {
  checkDefaultRuleDiscardChanges,
  handleCheckIndividualDiscardChanges,
  handleCheckPrerequisiteDiscardChanges
} from 'pages/feature-flag-details/targeting/utils';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import { ReactDatePicker } from 'components/date-time-picker';
import Form from 'components/form';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import DiscardChangeItems from '../discard-change-items';
import {
  CustomRuleDiscardItem,
  IndividualDiscardItem,
  PrerequisiteDiscardItem
} from '../discard-changes-modal';
import {
  formSchema,
  SCHEDULE_TYPE_SCHEDULE,
  SCHEDULE_TYPE_UPDATE_NOW
} from './form-schema';

export type ConfirmationRequiredModalProps = {
  feature: Feature;
  targetingRule?: TargetingSchema;
  activeFeatures?: Feature[];
  isOpen: boolean;
  isShowScheduleSelect?: boolean;
  isShowRolloutWarning?: boolean;
  onSegmentRuleDeleted?: () => DiscardChangesStateData[];
  onSegmentRuleChannge?: (
    index: number,
    isAction: boolean
  ) => DiscardChangesStateData[];
  onClose: () => void;
  onSubmit: (values: ConfirmRequiredValues) => Promise<void>;
};

export interface ConfirmRequiredValues {
  resetSampling?: boolean;
  comment?: string;
  requireComment?: boolean;
  scheduleType?: string;
  scheduleAt?: string;
}

interface DiscardItemRendererProps<T> {
  title: string;
  items: T[];
  Renderer: React.FC<T>;
}

const ConfirmationRequiredModal = ({
  feature,
  targetingRule,
  activeFeatures,
  isOpen,
  isShowScheduleSelect,
  isShowRolloutWarning,
  onSegmentRuleDeleted,
  onSegmentRuleChannge,
  onClose,
  onSubmit
}: ConfirmationRequiredModalProps) => {
  const { t } = useTranslation(['common', 'form', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const navigate = useNavigate();

  const { data: rolloutCollection } = useQueryRollouts({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      featureIds: [feature.id]
    },
    refetchOnMount: !!feature && isShowRolloutWarning ? 'always' : false,
    enabled: !!feature && isShowRolloutWarning
  });

  const hasRolloutRunning = !!rolloutCollection?.progressiveRollouts?.find(
    item => ['WAITING', 'RUNNING'].includes(item.status)
  );

  const individualChange = useMemo(
    () =>
      targetingRule?.individualRules
        ? handleCheckIndividualDiscardChanges(
            feature,
            targetingRule.individualRules
          )
        : [],
    [targetingRule]
  );

  const prerequisiteChanges = useMemo(
    () =>
      targetingRule?.prerequisites && activeFeatures
        ? handleCheckPrerequisiteDiscardChanges(
            targetingRule.prerequisites,
            feature,
            activeFeatures
          )
        : [],
    [targetingRule]
  );

  const segmentRuleDeletedChanges = useMemo(() => {
    if (onSegmentRuleDeleted) {
      return onSegmentRuleDeleted();
    }
    return [];
  }, [onSegmentRuleDeleted]);

  const segmentRulesChange = useMemo(() => {
    const change: {
      rule: number;
      changes: DiscardChangesStateData[];
      action?: 'new-rule' | 'edit-rule';
    }[] = [];

    if (!targetingRule) return [];
    for (let i = 0; i < targetingRule.segmentRules!.length; i++) {
      if (!onSegmentRuleChannge) continue;

      const changes = onSegmentRuleChannge(i, false);
      if (changes.length > 0) {
        // Detect action type from the first change item's changeType
        const action =
          changes[0]?.changeType === 'new-rule' ? 'new-rule' : 'edit-rule';
        change.push({
          rule: i + 1,
          changes: changes,
          action: action
        });
      }
    }

    return change;
  }, [targetingRule]);

  const defaultRulesChange = useMemo(
    () =>
      targetingRule?.defaultRule
        ? checkDefaultRuleDiscardChanges(
            feature.defaultStrategy,
            targetingRule.defaultRule as FeatureRuleStrategy,
            feature.variations
          )
        : [],
    [targetingRule]
  );

  const changeBreakdown = useMemo(() => {
    let adds = 0;
    let updates = 0;
    let deletes = 0;

    // Count prerequisites
    prerequisiteChanges?.forEach(item => {
      if (item.labelType === 'ADD') adds++;
      else if (item.labelType === 'UPDATE') updates++;
      else if (item.labelType === 'REMOVE') deletes++;
    });

    // Count individual targets
    individualChange?.forEach(item => {
      if (item.labelType === 'ADD') adds++;
      else if (item.labelType === 'REMOVE') deletes++;
    });

    // Count segment rule changes
    segmentRulesChange?.forEach(({ changes, action }) => {
      if (action === 'new-rule') {
        adds++;
      } else {
        changes.forEach(change => {
          if (change.labelType === 'ADD') adds++;
          else if (change.labelType === 'UPDATE') updates++;
          else if (change.labelType === 'REMOVE') deletes++;
        });
      }
    });

    // Count deleted rules
    deletes += segmentRuleDeletedChanges?.length ?? 0;

    // Count default rule changes
    defaultRulesChange?.forEach(item => {
      if (item.labelType === 'ADD') adds++;
      else if (item.labelType === 'UPDATE') updates++;
    });

    const total = adds + updates + deletes;
    return { adds, updates, deletes, total };
  }, [
    segmentRulesChange,
    defaultRulesChange,
    individualChange,
    prerequisiteChanges,
    segmentRuleDeletedChanges
  ]);

  const isShowChange = changeBreakdown.total;

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      comment: '',
      resetSampling: false,
      requireComment: currentEnvironment.requireComment,
      scheduleType: SCHEDULE_TYPE_UPDATE_NOW,
      scheduleAt: String(Math.floor((new Date().getTime() + 3600000) / 1000))
    },
    mode: 'onChange'
  });

  const {
    control,
    formState: { isDirty, isValid, isSubmitting },
    watch,
    setValue
  } = form;
  const isRequireComment = watch('requireComment');
  const isShowSchedule = watch('scheduleType') === SCHEDULE_TYPE_SCHEDULE;

  const handleOnSubmit = async (values: ConfirmRequiredValues) => {
    await onSubmit(values);
  };

  const renderDiscardSection = <T,>({
    title,
    items,
    Renderer
  }: DiscardItemRendererProps<T>): ReactNode => {
    if (!items || !items?.length) return null;
    return (
      <DiscardChangeItems title={title}>
        <div className="flex flex-col gap-2 pl-4">
          {items.map((item, idx) => (
            <Renderer key={idx} {...item} />
          ))}
        </div>
      </DiscardChangeItems>
    );
  };

  const renderSegmentRuleChanges = (): ReactNode => {
    const showCustomRuleChange =
      !!segmentRulesChange.length || !!segmentRuleDeletedChanges.length;
    if (!showCustomRuleChange) return null;
    return (
      <DiscardChangeItems title={t('common:custom-rule')}>
        <div className="flex flex-col gap-2 pl-4">
          {segmentRulesChange.map(({ rule, changes, action }) => (
            <div key={rule}>
              <div className="flex pb-2 gap-1 items-center typo-para-medium leading-[1px] my-2 text-gray-700">
                <Trans
                  i18nKey={
                    action === 'new-rule'
                      ? 'common:add-rule'
                      : 'common:edit-rule'
                  }
                  values={{ rule }}
                  components={{
                    b: <strong />
                  }}
                />
              </div>
              {changes.map((item, idx) => {
                if (!isNil(item)) {
                  return <CustomRuleDiscardItem key={idx} {...item} />;
                }
              })}
            </div>
          ))}
          {!!segmentRuleDeletedChanges.length && (
            <>
              {segmentRuleDeletedChanges.map((item, index) => (
                <div key={index}>
                  {item.ruleIndex && (
                    <div className="flex pb-2 gap-1 items-center typo-para-medium leading-[1px] my-2 text-gray-700 text-accent-red-500">
                      <Trans
                        i18nKey="common:delete-rule"
                        values={{ rule: item.ruleIndex }}
                        components={{
                          b: <strong />
                        }}
                      />
                    </div>
                  )}
                  <CustomRuleDiscardItem {...item} />
                </div>
              ))}
            </>
          )}
        </div>
      </DiscardChangeItems>
    );
  };

  const changeKey = isShowChange > 1 ? 'changes' : 'change';
  const changeText = t(`common:${changeKey}`)?.toLowerCase();
  const versionFeature = feature.version + 1;

  const transComponents = useMemo(() => ({ b: <strong /> }), []);

  return (
    <DialogModal
      className="w-full max-w-[640px]"
      title={t('table:feature-flags.confirm-required')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(handleOnSubmit)}>
          <div className="flex flex-col w-full max-h-[80vh] gap-y-5 items-start pt-5">
            <div className="relative overflow-auto w-full h-full small-scroll">
              {!!isShowChange && (
                <>
                  <div className="sticky top-0 z-20 bg-white typo-para-small text-gray-600 w-full px-5">
                    <p className="typo-para-medium leading-4 text-gray-700 pb-5">
                      <Trans
                        i18nKey="common:change-count-breakdown"
                        values={{
                          count: isShowChange,
                          changeText,
                          versionFeature,
                          adds: changeBreakdown.adds,
                          updates: changeBreakdown.updates,
                          deletes: changeBreakdown.deletes
                        }}
                        components={transComponents}
                      />
                    </p>
                  </div>

                  <div className="w-full flex flex-col px-5 pb-5 gap-6 ">
                    {renderDiscardSection({
                      title: t('form:feature-flags.prerequisites'),
                      items: prerequisiteChanges ? prerequisiteChanges : [],
                      Renderer: PrerequisiteDiscardItem
                    })}

                    {renderDiscardSection({
                      title: t('form:targeting.individual-target'),
                      items: individualChange ? individualChange : [],
                      Renderer: IndividualDiscardItem
                    })}

                    {renderSegmentRuleChanges()}

                    {renderDiscardSection({
                      title: t('form:targeting.default-rule'),
                      items: defaultRulesChange ? defaultRulesChange : [],
                      Renderer: CustomRuleDiscardItem
                    })}
                  </div>
                </>
              )}
              <div className="flex flex-col w-full px-5 pb-5">
                <Form.Field
                  control={control}
                  name="comment"
                  render={({ field }) => (
                    <Form.Item className="py-0">
                      <Form.Label
                        required={isRequireComment && !isShowSchedule}
                      >
                        {t('form:comment-for-update')}
                      </Form.Label>
                      <Form.Control>
                        <TextArea
                          placeholder={`${t('form:placeholder-comment')}`}
                          rows={3}
                          {...field}
                          onChange={value => {
                            field.onChange(value);
                          }}
                          name="comment"
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
                <Form.Field
                  control={control}
                  name="resetSampling"
                  render={({ field }) => (
                    <Form.Item className="flex flex-col w-full py-0 gap-y-4 mt-5">
                      <div className="flex items-center gap-x-2">
                        <Form.Control>
                          <Checkbox
                            ref={field.ref}
                            checked={field.value}
                            onCheckedChange={checked => field.onChange(checked)}
                            title={t('form:reset-sampling')}
                          />
                        </Form.Control>
                        <Tooltip
                          align="start"
                          content={t('form:reset-sampling-tooltip')}
                          trigger={
                            <div className="flex-center size-fit">
                              <Icon
                                icon={IconInfo}
                                size="xs"
                                color="gray-500"
                              />
                            </div>
                          }
                          className="max-w-[400px]"
                        />
                      </div>
                      <Form.Message />
                    </Form.Item>
                  )}
                />

                {isShowScheduleSelect && (
                  <>
                    <Form.Field
                      control={form.control}
                      name="scheduleType"
                      render={({ field }) => (
                        <Form.Item className="flex flex-col w-full py-0 gap-y-4 mt-5">
                          <Form.Control>
                            <RadioGroup
                              defaultValue={field.value}
                              className="flex flex-col w-full gap-y-4"
                              onValueChange={value => {
                                field.onChange(value);
                                setValue(
                                  'requireComment',
                                  value !== SCHEDULE_TYPE_SCHEDULE &&
                                    currentEnvironment.requireComment
                                );
                              }}
                            >
                              <div className="flex items-center gap-x-2">
                                <RadioGroupItem
                                  id="active_now"
                                  value={SCHEDULE_TYPE_UPDATE_NOW}
                                />
                                <label
                                  htmlFor="active_now"
                                  className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
                                >
                                  {t('update-now')}
                                </label>
                              </div>

                              <div className="flex items-center gap-x-2">
                                <RadioGroupItem
                                  id="schedule"
                                  value={SCHEDULE_TYPE_SCHEDULE}
                                />
                                <label
                                  htmlFor="schedule"
                                  className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
                                >
                                  {t('form:feature-flags.schedule-the-updates')}
                                </label>
                                <span className="px-2 py-1.5 rounded-[3px] bg-accent-blue-50 text-accent-blue-500 typo-para-small leading-[14px] whitespace-nowrap uppercase">
                                  New
                                </span>
                                <Tooltip
                                  align="start"
                                  content={t(
                                    'form:feature-flags.schedule-the-updates-tooltip'
                                  )}
                                  trigger={
                                    <div className="flex-center size-fit">
                                      <Icon
                                        icon={IconInfo}
                                        size="xs"
                                        color="gray-500"
                                      />
                                    </div>
                                  }
                                  className="max-w-[400px]"
                                />
                              </div>
                            </RadioGroup>
                          </Form.Control>
                          <Form.Message />
                        </Form.Item>
                      )}
                    />
                    {isShowSchedule && (
                      <Form.Field
                        control={form.control}
                        name="scheduleAt"
                        render={({ field }) => {
                          const scheduleDate = field.value
                            ? new Date(+field.value * 1000)
                            : null;

                          return (
                            <Form.Item className="py-0 mt-5">
                              <Form.Control>
                                <div className="flex gap-x-4">
                                  <div>
                                    <Form.Label required>
                                      {t('form:feature-flags.update-date')}
                                    </Form.Label>
                                    <ReactDatePicker
                                      dateFormat="yyyy/MM/dd"
                                      minDate={new Date()}
                                      selected={scheduleDate}
                                      showTimeSelect={false}
                                      className="w-[186px]"
                                      onChange={date => {
                                        if (date) {
                                          if (scheduleDate) {
                                            date.setHours(
                                              scheduleDate.getHours(),
                                              scheduleDate.getMinutes(),
                                              0,
                                              0
                                            );
                                          }
                                          field.onChange(
                                            String(
                                              Math.floor(date.getTime() / 1000)
                                            )
                                          );
                                        }
                                      }}
                                    />
                                  </div>
                                  <div>
                                    <Form.Label required>
                                      {t('form:feature-flags.update-time')}
                                    </Form.Label>
                                    <ReactDatePicker
                                      dateFormat="HH:mm"
                                      timeFormat="HH:mm"
                                      selected={scheduleDate}
                                      showTimeSelectOnly={true}
                                      className="w-[124px]"
                                      onChange={date => {
                                        if (date) {
                                          field.onChange(
                                            String(
                                              Math.floor(date.getTime() / 1000)
                                            )
                                          );
                                        }
                                      }}
                                      icon={
                                        <Icon
                                          icon={IconWatch}
                                          className="flex-center"
                                        />
                                      }
                                    />
                                  </div>
                                </div>
                              </Form.Control>
                              <Form.Message />
                            </Form.Item>
                          );
                        }}
                      />
                    )}
                  </>
                )}
                {isShowRolloutWarning && hasRolloutRunning && (
                  <div className="flex w-full gap-x-3 p-4 mt-5 rounded-md bg-accent-yellow-50 typo-para-small">
                    <Icon icon={IconToastWarning} />
                    <p className="w-full typo-para-medium text-accent-yellow-700">
                      <Trans
                        i18nKey={'form:has-rollout-running'}
                        components={{
                          comp: (
                            <span
                              onClick={() =>
                                navigate(
                                  `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${feature.id}${PAGE_PATH_FEATURE_AUTOOPS}`
                                )
                              }
                              className="inline-flex text-primary-500 underline whitespace-nowrap cursor-pointer"
                            />
                          )
                        }}
                      />
                    </p>
                  </div>
                )}
              </div>
            </div>
            <ButtonBar
              secondaryButton={
                <Button
                  type="submit"
                  loading={isSubmitting}
                  disabled={(isRequireComment && !isDirty) || !isValid}
                >
                  {t(`submit`)}
                </Button>
              }
              primaryButton={
                <Button type="button" onClick={onClose} variant="secondary">
                  {t(`cancel`)}
                </Button>
              }
            />
          </div>
        </Form>
      </FormProvider>
    </DialogModal>
  );
};

export default ConfirmationRequiredModal;
