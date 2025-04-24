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
import { useTranslation } from 'i18n';
import {
  AutoOpsRule,
  ClauseActionType,
  Feature,
  OpsEventRateClause
} from '@types';
import { IconInfo, IconPlus } from '@icons';
import {
  eventRateSchema,
  EventRateSchemaType
} from 'pages/feature-flag-details/operations/form-schema';
import { createEventRate } from 'pages/feature-flag-details/operations/utils';
import { OperationActionType } from 'pages/feature-flag-details/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import SlideModal from 'components/modal/slide';
import CreateGoalModal from 'elements/create-goal-modal';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';

export interface OperationModalProps {
  feature: Feature;
  environmentId: string;
  isOpen: boolean;
  actionType: OperationActionType;
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
  feature,
  environmentId,
  isOpen,
  actionType,
  selectedData,
  onClose,
  onSubmitOperationSuccess
}: OperationModalProps) => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { notify, errorNotify } = useToast();

  const isCreate = useMemo(() => actionType === 'NEW', [actionType]);

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
      feature.variations.map(item => ({
        label: item.name || item.value,
        value: item.id
      })),
    [feature]
  );
  const handleCreateDefaultValues = () => {
    const clause = selectedData?.clauses[0]?.clause as OpsEventRateClause;
    if (clause) {
      return {
        ...clause,
        minCount: +clause?.minCount || 0
      } as EventRateSchemaType;
    }
    return createEventRate(feature);
  };

  const form = useForm({
    resolver: yupResolver(eventRateSchema),
    defaultValues: {
      ...handleCreateDefaultValues()
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isSubmitting, errors }
  } = form;

  const onSubmit = useCallback(
    async (values: EventRateSchemaType) => {
      try {
        console.log(values);
        let resp: AutoOpsCreatorResponse | null = null;

        if (!isCreate && selectedData) {
          resp = await autoOpsUpdate({
            environmentId,
            updateDatetimeClauses: {
              id: selectedData.id,
              delete: false,
              clause: []
            }
          });
        } else {
          resp = await autoOpsCreator({
            featureId: feature.id,
            environmentId,
            opsType: 'EVENT_RATE',
            opsEventRateClauses: [
              {
                ...values,
                actionType: values.actionType as ClauseActionType,
                minCount: values.minCount.toString(),
                threadsholdRate: values.threadsholdRate / 100
              }
            ]
          });
        }

        if (resp) {
          onSubmitOperationSuccess();
          notify({
            message: `Schedule operation ${isCreate ? 'created' : 'updated'} successfully`
          });
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [isCreate, actionType, selectedData]
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
            <p className="typo-head-bold-small text-gray-800">
              {t('feature-flags.event-rate')}
            </p>
            <div className="flex items-center w-full">
              <div className="pr-4">
                <p className="flex-center w-[42px] h-[26px] rounded-[3px] bg-accent-pink-50 typo-para-small text-accent-pink-500">
                  {t('common:if')}
                </p>
              </div>
              <div className="flex flex-col flex-1 pl-4 gap-y-4 border-l border-primary-500">
                <Form.Field
                  control={form.control}
                  name={`variationId`}
                  render={({ field }) => (
                    <Form.Item className="py-0">
                      <Form.Label required className="relative w-fit">
                        {t('table:results.variation')}
                        <Icon
                          icon={IconInfo}
                          size="xs"
                          color="gray-500"
                          className="absolute -right-6"
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
                        <Icon
                          icon={IconInfo}
                          size="xs"
                          color="gray-500"
                          className="absolute -right-6"
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
                  <div className="flex flex-1 gap-x-4">
                    <Form.Field
                      control={form.control}
                      name={`operator`}
                      render={({ field }) => (
                        <Form.Item className="py-0 flex-1 h-full">
                          <Form.Label required className="relative w-fit">
                            {t('condition')}
                            <Icon
                              icon={IconInfo}
                              size="xs"
                              color="gray-500"
                              className="absolute -right-6"
                            />
                          </Form.Label>
                          <Form.Control>
                            <DropdownMenu>
                              <DropdownMenuTrigger
                                label={
                                  conditionOptions.find(
                                    item => item.value === field.value
                                  )?.label || ''
                                }
                                isExpand
                              />
                              <DropdownMenuContent
                                align="end"
                                className="min-w-[132px]"
                              >
                                {conditionOptions.map((item, index) => (
                                  <DropdownMenuItem
                                    key={index}
                                    label={item.label}
                                    value={item.value}
                                    onSelectOption={field.onChange}
                                  />
                                ))}
                              </DropdownMenuContent>
                            </DropdownMenu>
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
                            <Icon
                              icon={IconInfo}
                              size="xs"
                              color="gray-500"
                              className="absolute -right-6"
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
                            <Icon
                              icon={IconInfo}
                              size="xs"
                              color="gray-500"
                              className="absolute -right-6"
                            />
                          </Form.Label>
                          <Form.Control>
                            <Input
                              {...field}
                              value={field.value || ''}
                              type="number"
                              className="pr-4"
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
                <Button variant="secondary" onClick={onClose}>
                  {t(`common:cancel`)}
                </Button>
              }
              secondaryButton={
                <Button
                  type="submit"
                  loading={isSubmitting}
                  disabled={!isValid}
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
            form.setValue('goalId', goal.id);
          }}
        />
      )}
    </SlideModal>
  );
};

export default EventRateOperationModal;
