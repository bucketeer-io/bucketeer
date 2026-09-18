import { useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import {
  autoOpsCreator,
  AutoOpsCreatorResponse,
  autoOpsUpdate
} from '@api/auto-ops';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryGoals } from '@queries/goals';
import { useToast, useToggleOpen } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import { AutoOpsRule, Feature, OpsEventRateClause } from '@types';
import { IconInfo, IconPlus } from '@icons';
import {
  eventRateSchema,
  EventRateSchemaType
} from 'pages/feature-flag-details/operations/form-schema';
import { OperationActionType } from 'pages/feature-flag-details/operations/types';
import { createEventRate } from 'pages/feature-flag-details/operations/utils';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import SlideModal from 'components/modal/slide';
import { Tooltip } from 'components/tooltip';
import CreateGoalModal from 'elements/create-goal-modal';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';

export interface OperationModalProps {
  editable: boolean;
  feature: Feature;
  environmentId: string;
  isOpen: boolean;
  actionType: OperationActionType;
  isFinishedTab: boolean;
  selectedData?: AutoOpsRule;
  onClose: () => void;
  onSubmitOperationSuccess: () => void;
}

const CreateNewOptionButton = ({
  text,
  onClick
}: {
  text: string;
  onClick: () => void;
}) => (
  <Button
    type="button"
    variant="text"
    className="h-10 self-center w-full bg-white hover:bg-gray-100 sticky left-0 right-0 bottom-0 border-t border-gray-200"
    onClick={onClick}
  >
    <Icon icon={IconPlus} color="primary-500" size={'xs'} />
    {text}
  </Button>
);

const conditionOptions = [
  {
    label: '>=',
    value: 'GREATER_OR_EQUAL'
  },
  {
    label: '<=',
    value: 'LESS_OR_EQUAL'
  }
];

const EventRateOperationModal = ({
  editable,
  feature,
  environmentId,
  isOpen,
  actionType,
  selectedData,
  isFinishedTab,
  onClose,
  onSubmitOperationSuccess
}: OperationModalProps) => {
  const { t } = useTranslation(['form', 'common', 'table', 'message']);
  const { notify, errorNotify } = useToast();

  const isCreate = useMemo(() => actionType === 'NEW', [actionType]);

  const isDisabled = useMemo(
    () => (!isCreate && isFinishedTab) || !editable,
    [isCreate, isFinishedTab, editable]
  );

  const [
    isOpenCreateGoalModal,
    onOpenCreateGoalModal,
    onHiddenCreateGoalModal
  ] = useToggleOpen(false);

  const { data: goalCollection, isLoading: isLoadingGoals } = useQueryGoals({
    params: {
      cursor: String(0),
      environmentId,
      archived: false,
      connectionType: 'OPERATION'
    }
  });

  const goalOptions = useMemo(() => {
    return (
      goalCollection?.goals?.map(item => ({
        label: item.name,
        value: item.id
      })) || []
    );
  }, [goalCollection]);

  const variationOptions = useMemo(
    () =>
      feature.variations.map((item, index) => ({
        label: (
          <div className="flex items-center gap-x-2 pl-0.5">
            <FlagVariationPolygon index={index} />
            <p className="-mt-0.5 truncate">{item.name || item.value}</p>
          </div>
        ),
        labelText: item.name || item.value,
        value: item.id
      })),
    [feature]
  );
  const handleCreateDefaultValues = () => {
    const clause = selectedData?.clauses[0]?.clause as OpsEventRateClause;
    if (clause) {
      return {
        ...clause,
        minCount: +clause?.minCount || 0,
        threadsholdRate: clause.threadsholdRate * 100 || 0
      } as EventRateSchemaType;
    }
    return createEventRate(feature);
  };

  const form = useForm({
    resolver: yupResolver(useFormSchema(eventRateSchema)),
    defaultValues: {
      ...handleCreateDefaultValues()
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isDirty, isSubmitting, errors }
  } = form;
  useUnsavedLeavePage({ isShow: isDirty && !isSubmitting });
  const onSubmit = useCallback(
    async (values: EventRateSchemaType) => {
      if (editable) {
        try {
          let resp: AutoOpsCreatorResponse | null = null;
          if (!isCreate && selectedData) {
            resp = await autoOpsUpdate({
              environmentId,
              id: selectedData.id,
              opsEventRateClauseChanges: [
                {
                  id: selectedData.clauses[0].id,
                  changeType: 'UPDATE',
                  clause: {
                    ...values,
                    minCount: values.minCount.toString(),
                    threadsholdRate: values.threadsholdRate / 100
                  }
                }
              ]
            });
          } else {
            resp = await autoOpsCreator({
              featureId: feature.id,
              environmentId,
              opsType: 'EVENT_RATE',
              opsEventRateClauses: [
                {
                  ...values,
                  actionType: 'DISABLE',
                  minCount: values.minCount.toString(),
                  threadsholdRate: values.threadsholdRate / 100
                }
              ]
            });
          }

          if (resp) {
            onSubmitOperationSuccess();
            notify({
              message: t(`message:collection-action-success`, {
                collection: t('common:operation'),
                action: t(`common:${isCreate ? 'created' : 'updated'}`)
              })
            });
          }
        } catch (error) {
          errorNotify(error);
        }
      }
    },
    [isCreate, actionType, selectedData, editable]
  );

  return (
    <SlideModal
      title={t(`common:${isCreate ? 'new' : 'update'}-operation`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col gap-y-5 w-full p-5 pb-28">
            <div className="flex flex-col gap-y-3">
              <p className="typo-head-bold-small text-gray-800">
                {t('feature-flags.kill-switch')}
              </p>
              <p className="typo-para-small text-gray-500">
                {t('table:feature-flags.event-rate-create-desc')}
              </p>
            </div>
            <div className="flex items-center w-full">
              <div className="pr-4">
                <p className="flex-center w-[42px] h-[26px] rounded-[3px] bg-accent-pink-50 typo-para-small text-accent-pink-500">
                  {t('common:if')}
                </p>
              </div>
              <div className="flex flex-col flex-1 max-w-[430px] pl-4 gap-y-4 border-l border-primary-500">
                <Form.Field
                  control={form.control}
                  name={`variationId`}
                  render={({ field }) => (
                    <Form.Item className="py-0">
                      <Form.Label required className="relative w-fit">
                        {t('form:flag-variation')}
                        <Tooltip
                          content={t('event-rate-tooltip.variation')}
                          trigger={
                            <div className="flex-center absolute top-0 -right-6">
                              <Icon
                                icon={IconInfo}
                                size="xs"
                                color="gray-500"
                              />
                            </div>
                          }
                          className="max-w-[300px]"
                        />
                      </Form.Label>
                      <Form.Control>
                        <DropdownMenuWithSearch
                          align="end"
                          label={
                            variationOptions.find(
                              item => item.value === field.value
                            )?.label || ''
                          }
                          itemSelected={field.value}
                          disabled={isDisabled}
                          contentClassName="max-w-[427px] [&>div.wrapper-menu-items>div]:px-4"
                          options={variationOptions}
                          onSelectOption={field.onChange}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
                <Form.Field
                  control={form.control}
                  name={`goalId`}
                  render={({ field }) => (
                    <Form.Item className="py-0">
                      <Form.Label required className="relative w-fit">
                        {t('common:goal')}
                        <Tooltip
                          align="start"
                          alignOffset={-50}
                          content={t('event-rate-tooltip.goal')}
                          trigger={
                            <div className="flex-center absolute top-0 -right-6">
                              <Icon
                                icon={IconInfo}
                                size="xs"
                                color="gray-500"
                              />
                            </div>
                          }
                          className="max-w-[300px]"
                        />
                      </Form.Label>
                      <Form.Control>
                        <DropdownMenuWithSearch
                          align="end"
                          hidden={isOpenCreateGoalModal}
                          isLoading={isLoadingGoals}
                          placeholder={t(`experiments.select-goal`)}
                          label={
                            goalOptions.find(item => item.value === field.value)
                              ?.label || ''
                          }
                          itemSelected={field.value}
                          disabled={isDisabled}
                          contentClassName="max-w-[427px] [&>div.wrapper-menu-items>div]:px-4"
                          options={goalOptions}
                          createNewOption={
                            <CreateNewOptionButton
                              text={t('common:create-a-new-goal')}
                              onClick={onOpenCreateGoalModal}
                            />
                          }
                          onSelectOption={field.onChange}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
                <div className="flex flex-col w-full gap-y-1">
                  <div className="flex flex-col sm:flex-row flex-1 gap-4">
                    <Form.Field
                      control={form.control}
                      name={`operator`}
                      render={({ field }) => (
                        <Form.Item className="py-0 flex-1 h-full">
                          <Form.Label required className="relative w-fit">
                            {t('condition')}
                            <Tooltip
                              content={t('event-rate-tooltip.condition')}
                              trigger={
                                <div className="flex-center absolute top-0 -right-6">
                                  <Icon
                                    icon={IconInfo}
                                    size="xs"
                                    color="gray-500"
                                  />
                                </div>
                              }
                              className="max-w-[300px]"
                            />
                          </Form.Label>
                          <Form.Control>
                            <Dropdown
                              options={conditionOptions}
                              value={field.value}
                              onChange={field.onChange}
                              placeholder={t('form:select-condition')}
                              disabled={isDisabled}
                              className="w-full"
                              contentClassName="min-w-[132px]"
                            />
                          </Form.Control>
                        </Form.Item>
                      )}
                    />
                    <Form.Field
                      control={form.control}
                      name={`threadsholdRate`}
                      render={({ field }) => (
                        <Form.Item className="py-0 flex-1 h-full">
                          <Form.Label required className="relative w-fit">
                            {t('threshold')}
                            <Tooltip
                              align="start"
                              alignOffset={-250}
                              content={
                                <p
                                  dangerouslySetInnerHTML={{
                                    __html: t('event-rate-tooltip.threshold')
                                  }}
                                />
                              }
                              trigger={
                                <div className="flex-center absolute top-0 -right-6">
                                  <Icon
                                    icon={IconInfo}
                                    size="xs"
                                    color="gray-500"
                                  />
                                </div>
                              }
                              className="max-w-[450px]"
                            />
                          </Form.Label>
                          <Form.Control>
                            <InputGroup
                              className="w-full"
                              addonSlot="right"
                              addonSize="md"
                              addon={'%'}
                            >
                              <Input
                                {...field}
                                value={field.value || ''}
                                type="number"
                                className="pr-8"
                                disabled={isDisabled}
                              />
                            </InputGroup>
                          </Form.Control>
                        </Form.Item>
                      )}
                    />
                    <Form.Field
                      control={form.control}
                      name={`minCount`}
                      render={({ field }) => (
                        <Form.Item className="py-0 flex-1 h-full">
                          <Form.Label required className="relative w-fit">
                            {t('minimum-count')}
                            <Tooltip
                              align="end"
                              alignOffset={-10}
                              content={t('event-rate-tooltip.min-count')}
                              trigger={
                                <div className="flex-center absolute top-0 -right-6">
                                  <Icon
                                    icon={IconInfo}
                                    size="xs"
                                    color="gray-500"
                                  />
                                </div>
                              }
                              className="max-w-[300px]"
                            />
                          </Form.Label>
                          <Form.Control>
                            <Input
                              {...field}
                              value={field.value || ''}
                              type="number"
                              className="pr-4"
                              disabled={isDisabled}
                            />
                          </Form.Control>
                        </Form.Item>
                      )}
                    />
                  </div>
                  <div className="flex flex-col gap-y-0.5">
                    {errors?.operator?.message && (
                      <Form.Message>{errors.operator.message}</Form.Message>
                    )}
                    {errors?.threadsholdRate?.message && (
                      <Form.Message>
                        {errors.threadsholdRate.message}
                      </Form.Message>
                    )}
                    {errors?.minCount?.message && (
                      <Form.Message>{errors.minCount.message}</Form.Message>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
            <ButtonBar
              primaryButton={
                <Button type="button" variant="secondary" onClick={onClose}>
                  {t(`common:cancel`)}
                </Button>
              }
              secondaryButton={
                <Button
                  type="submit"
                  loading={isSubmitting}
                  disabled={!isValid || !isDirty || isDisabled}
                >
                  {t(
                    isCreate
                      ? `feature-flags.create-operation`
                      : 'common:update-operation'
                  )}
                </Button>
              }
            />
          </div>
        </Form>
      </FormProvider>
      {isOpenCreateGoalModal && (
        <CreateGoalModal
          isOpen={isOpenCreateGoalModal}
          connectionType="OPERATION"
          onClose={onHiddenCreateGoalModal}
          onCompleted={goal => {
            form.setValue('goalId', goal.id, {
              shouldDirty: true,
              shouldValidate: true
            });
          }}
        />
      )}
    </SlideModal>
  );
};

export default EventRateOperationModal;
