import { useCallback, useEffect, useMemo } from 'react';
import {
  ControllerRenderProps,
  FormProvider,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import {
  ExperimentCreateUpdateResponse,
  experimentCreator,
  experimentUpdater
} from '@api/experiment';
import { yupResolver } from '@hookform/resolvers/yup';
import {
  invalidateExperimentDetails,
  useQueryExperimentDetails
} from '@queries/experiment-details';
import { invalidateExperiments } from '@queries/experiments';
import { useQueryFeature } from '@queries/feature-details';
import { useQueryGoals } from '@queries/goals';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useFeatureFlagsLoader } from 'hooks/use-feature-loading-more';
import useFormSchema from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import { IconInfo, IconPlus } from '@icons';
import { createExperimentFormSchema } from 'pages/experiments/form-schema';
import CreateFlagForm from 'pages/feature-flags/flags-modal/add-flag-modal/create-flag-form';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { ReactDatePicker } from 'components/date-time-picker';
import Divider from 'components/divider';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import DialogModal from 'components/modal/dialog';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import CreateGoalModal from 'elements/create-goal-modal';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import FeatureFlagStatus from 'elements/feature-flag-status';
import FormLoading from 'elements/form-loading';
import VariationLabel from 'elements/variation-label';

interface ExperimentCreateUpdateModalProps {
  disabled: boolean;
  isOpen: boolean;
  onClose: () => void;
}

type StartType = 'manual' | 'schedule';

export interface ExperimentCreateUpdateForm {
  id?: string;
  baseVariationId: string;
  name: string;
  description?: string;
  startType: StartType;
  startAt: string;
  stopAt: string;
  audience?: {
    rule: string;
    inExperiment: number;
    notInExperiment: number;
    served: boolean;
    variationReassignment: boolean;
  };
  featureId: string;
  goalIds: string[];
}

export type DefineAudienceField = ControllerRenderProps<
  ExperimentCreateUpdateForm,
  'audience'
>;

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

const ExperimentCreateUpdateModal = ({
  disabled,
  isOpen,
  onClose
}: ExperimentCreateUpdateModalProps) => {
  const { t } = useTranslation(['form', 'common']);
  const { notify, errorNotify } = useToast();
  const formSchema = useFormSchema(createExperimentFormSchema);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);
  const queryClient = useQueryClient();

  const [
    isOpenCreateGoalModal,
    onOpenCreateGoalModal,
    onHiddenCreateGoalModal
  ] = useToggleOpen(false);

  const [
    isOpenCreateFlagModal,
    onOpenCreateFlagModal,
    onHiddenCreateFlagModal
  ] = useToggleOpen(false);

  const { isEdit, params } = useActionWithURL({
    closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}`
  });
  const experimentId = useMemo(() => params?.experimentId, [params]);

  const {
    data: experimentCollection,
    isLoading: experimentLoading,
    error: experimentError
  } = useQueryExperimentDetails({
    params: {
      id: experimentId as string,
      environmentId: currentEnvironment.id
    },
    enabled: !!experimentId
  });

  const experiment = useMemo(
    () => experimentCollection?.experiment,
    [experimentCollection]
  );

  const isEnabledEdit = useMemo(
    () =>
      ['WAITING', 'NOT_STARTED'].includes(experiment?.status as string) ||
      !experiment,
    [experiment]
  );

  // Fetch the experiment's feature when in edit mode to ensure we have the feature name
  const { data: experimentFeatureData } = useQueryFeature({
    params: {
      id: experiment?.featureId as string,
      environmentId: currentEnvironment.id
    },
    enabled: Boolean(experiment?.featureId && isEdit)
  });

  const experimentFeature = useMemo(
    () => experimentFeatureData?.feature,
    [experimentFeatureData]
  );

  const { data: goalCollection, isLoading: isLoadingGoals } = useQueryGoals({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
    }
  });

  const goalOptions = useMemo(() => {
    return (
      goalCollection?.goals
        ?.filter(item => !item.archived && item.connectionType !== 'OPERATION')
        .map(item => ({
          label: item.name,
          value: item.id
        })) || []
    );
  }, [goalCollection]);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: '',
      name: '',
      baseVariationId: '',
      description: '',
      startType: 'manual',
      startAt: '',
      stopAt: '',
      audience: {
        rule: '',
        inExperiment: 5,
        notInExperiment: 95,
        served: true,
        variationReassignment: false
      },
      featureId: '',
      goalIds: []
    },
    mode: 'onChange'
  });

  const {
    watch,
    formState: { isDirty, isSubmitting }
  } = form;

  useUnsavedLeavePage({
    isShow: isDirty && !isSubmitting
  });
  const featureId = watch('featureId');

  const {
    allAvailableFlags,
    remainingFlagOptions,
    isLoadingMore,
    isInitialLoading: isLoadingFeature,
    hasMore,
    onSearchChange,
    loadMore
  } = useFeatureFlagsLoader({
    environmentId: currentEnvironment.id,
    selectedFlagIds: featureId ? [featureId] : [],
    filterSelected: !isEdit
  });

  const featureFlagOptions = useMemo(
    () =>
      allAvailableFlags.map(feature => {
        return {
          value: feature.id,
          label: feature.name,
          enabled: feature.enabled,
          disabled: featureId === feature.id
        };
      }),
    [allAvailableFlags]
  );

  const variationOptions = useMemo(() => {
    // In edit mode, use variations from fetched experiment feature or experiment data
    const variations =
      isEdit && (experimentFeature || experiment)
        ? experimentFeature?.variations || experiment?.variations
        : allAvailableFlags?.find(item => item.id === featureId)?.variations;

    return (
      variations?.map((item, index) => ({
        label: <VariationLabel label={item.name || item.value} index={index} />,
        value: item.id
      })) || []
    );
  }, [isEdit, experimentFeature, experiment, allAvailableFlags, featureId]);

  // const startOptions = [
  //   {
  //     label: 'Manual Start',
  //     value: 'manual'
  //   },
  //   {
  //     label: 'Schedule',
  //     value: 'schedule'
  //   }
  // ];

  const onSubmit: SubmitHandler<ExperimentCreateUpdateForm> = useCallback(
    async values => {
      try {
        const {
          id,
          baseVariationId,
          featureId,
          goalIds,
          name,
          startAt,
          stopAt,
          description
        } = values;
        let resp: ExperimentCreateUpdateResponse | null = null;
        const formatStartAt = Math.floor(Number(startAt)).toString();
        const formatStopAt = Math.floor(Number(stopAt)).toString();
        if (isEdit) {
          resp = await experimentUpdater({
            id,
            name,
            description,
            startAt: formatStartAt,
            stopAt: formatStopAt,
            environmentId: currentEnvironment.id
          });
        } else {
          resp = await experimentCreator({
            baseVariationId,
            featureId,
            goalIds,
            name,
            startAt: formatStartAt,
            stopAt: formatStopAt,
            description,
            environmentId: currentEnvironment.id
          });
        }
        if (resp) {
          notify({
            message: t('message:collection-action-success', {
              collection: t('common:source-type.experiment'),
              action: t(isEdit ? 'common:updated' : 'common:created')
            })
          });
          invalidateExperiments(queryClient);
          invalidateExperimentDetails(queryClient, {
            id: experimentId as string,
            environmentId: currentEnvironment.id
          });
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [currentEnvironment, experimentId, isEdit]
  );

  useEffect(() => {
    if (experiment) {
      const {
        id,
        baseVariationId,
        name,
        description,
        startAt,
        stopAt,
        featureId,
        goalIds
      } = experiment;
      form.reset({
        id,
        baseVariationId,
        name,
        description,
        startAt,
        stopAt,
        startType: startAt && stopAt ? 'schedule' : 'manual',
        audience: {
          rule: '',
          inExperiment: 5,
          notInExperiment: 95,
          served: true,
          variationReassignment: false
        },
        featureId,
        goalIds
      });
    }
  }, [experiment, form]);

  useEffect(() => {
    if (experimentError) {
      errorNotify(experimentError);
    }
  }, [experimentError]);

  return (
    <SlideModal
      title={t(`common:${isEdit ? 'edit' : 'new'}-experiment`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      {experimentLoading ? (
        <FormLoading />
      ) : (
        <div className="p-5 pb-28">
          <p className="text-gray-800 typo-head-bold-small">
            {t('general-info')}
          </p>
          <FormProvider {...form}>
            <Form onSubmit={form.handleSubmit(onSubmit)}>
              <Form.Field
                control={form.control}
                name="name"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label required>{t('common:name')}</Form.Label>
                    <Form.Control>
                      <Input
                        disabled={disabled}
                        placeholder={`${t('placeholder-name')}`}
                        {...field}
                        name="experiment-name"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Field
                control={form.control}
                name="description"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label optional>{t('description')}</Form.Label>
                    <Form.Control>
                      <TextArea
                        disabled={disabled}
                        placeholder={t('placeholder-desc')}
                        rows={4}
                        {...field}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              {/* <RadioGroup
                defaultValue={startType}
                value={startType}
                onValueChange={value =>
                  form.setValue('startType', value as StartType)
                }
                className="flex flex-col gap-y-[18px]"
              >
                {startOptions.map(({ label, value }) => (
                  <Form.Field
                    key={value}
                    control={form.control}
                    name="startType"
                    render={() => (
                      <Form.Item className="py-0 last:pb-2">
                        <Form.Control>
                          <div className="flex items-center gap-x-2">
                            <RadioGroupItem value={value} id={value} />
                            <label
                              htmlFor={value}
                              className="typo-para-medium leading-4 text-gray-600 cursor-pointer"
                            >
                              {label}
                            </label>
                          </div>
                        </Form.Control>
                      </Form.Item>
                    )}
                  />
                ))}
              </RadioGroup> */}
              <div className="flex items-center w-full gap-x-4">
                <Form.Field
                  control={form.control}
                  name="startAt"
                  render={({ field }) => (
                    <Form.Item className="flex flex-col flex-1 h-full self-stretch">
                      <Form.Label required>{t('start-at')}</Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          disabled={!isEnabledEdit || disabled}
                          dateFormat={'yyyy/MM/dd'}
                          showTimeSelect={false}
                          selected={
                            field.value ? new Date(+field.value * 1000) : null
                          }
                          onChange={date => {
                            if (date) {
                              const timestamp = new Date(date)?.getTime();
                              field.onChange(timestamp / 1000);
                              form.trigger('startAt');
                            }
                          }}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
                <Form.Field
                  control={form.control}
                  name="startAt"
                  render={({ field }) => (
                    <Form.Item className="flex flex-col flex-1 h-full self-stretch">
                      <Form.Label required>{t('experiments.time')}</Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          disabled={!isEnabledEdit || disabled}
                          dateFormat={'HH:mm'}
                          showTimeSelect
                          showTimeSelectOnly={true}
                          selected={
                            field.value ? new Date(+field.value * 1000) : null
                          }
                          onChange={date => {
                            if (date) {
                              const timestamp = new Date(date)?.getTime();
                              field.onChange(timestamp / 1000);
                              form.trigger('startAt');
                            }
                          }}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              </div>
              <div className="flex items-center w-full gap-x-4">
                <Form.Field
                  control={form.control}
                  name="stopAt"
                  render={({ field }) => (
                    <Form.Item className="flex flex-col flex-1 h-full self-stretch">
                      <Form.Label required>{t('end-at')}</Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          disabled={!isEnabledEdit || disabled}
                          dateFormat={'yyyy/MM/dd'}
                          showTimeSelect={false}
                          selected={
                            field.value ? new Date(+field.value * 1000) : null
                          }
                          onChange={date => {
                            if (date) {
                              const timestamp = new Date(date)?.getTime();
                              field.onChange(timestamp / 1000);
                              form.trigger('stopAt');
                            }
                          }}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
                <Form.Field
                  control={form.control}
                  name="stopAt"
                  render={({ field }) => (
                    <Form.Item className="flex flex-col flex-1 h-full self-stretch">
                      <Form.Label required>{t('experiments.time')}</Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          disabled={!isEnabledEdit || disabled}
                          dateFormat={'HH:mm'}
                          showTimeSelect
                          showTimeSelectOnly={true}
                          selected={
                            field.value ? new Date(+field.value * 1000) : null
                          }
                          onChange={date => {
                            if (date) {
                              const timestamp = new Date(date)?.getTime();
                              field.onChange(timestamp / 1000);
                              form.trigger('stopAt');
                            }
                          }}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              </div>
              <Divider className="mt-3 mb-4" />
              <p className="text-gray-800 typo-head-bold-small mb-1">
                {t('common:flag')}
              </p>
              <Form.Field
                control={form.control}
                name={`featureId`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col w-full">
                    <Form.Label required>{t('common:flag')}</Form.Label>
                    <Form.Control>
                      <DropdownMenuWithSearch
                        disabled={!!isEdit || disabled}
                        hidden={isOpenCreateFlagModal}
                        isLoading={isLoadingFeature}
                        isLoadingMore={isLoadingMore}
                        isHasMore={hasMore || isLoadingMore}
                        onSearchChange={onSearchChange}
                        onHasMoreOptions={loadMore}
                        placeholder={t(`experiments.select-flag`)}
                        label={
                          (isEdit && experimentFeature
                            ? experimentFeature.name
                            : featureFlagOptions.find(
                                item => item.value === field.value
                              )?.label) || ''
                        }
                        options={remainingFlagOptions}
                        selectedOptions={[field.value]}
                        additionalElement={item => (
                          <FeatureFlagStatus
                            status={t(
                              item.enabled
                                ? 'experiments.on'
                                : 'experiments.off'
                            )}
                            enabled={item.enabled as boolean}
                          />
                        )}
                        createNewOption={
                          disabled ? undefined : (
                            <CreateNewOptionButton
                              text={t('common:create-a-new-flag')}
                              onClick={onOpenCreateFlagModal}
                            />
                          )
                        }
                        onSelectOption={field.onChange}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              {featureId && (
                <Form.Field
                  control={form.control}
                  name={`baseVariationId`}
                  render={({ field }) => (
                    <Form.Item className="flex flex-col w-full overflow-hidden">
                      <Form.Label required>
                        {t('experiments.base-variation')}
                      </Form.Label>
                      <Form.Control>
                        <Dropdown
                          disabled={!!isEdit || disabled}
                          placeholder={t(`experiments.select-variation`)}
                          className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                          contentClassName="min-w-[502px]"
                          options={variationOptions}
                          value={field.value}
                          onChange={field.onChange}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              )}
              <Divider className="mt-3 mb-4" />
              <p className="text-gray-800 typo-head-bold-small mb-1">
                {t('common:goals')}
              </p>
              <Form.Field
                control={form.control}
                name={`goalIds`}
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label required className="relative w-fit">
                      {t('common:goals')}
                      <Tooltip
                        align="start"
                        alignOffset={-52}
                        content={t('experiments.goals-tooltip')}
                        trigger={
                          <div className="flex-center absolute top-0 -right-6">
                            <Icon icon={IconInfo} size="xs" color="gray-500" />
                          </div>
                        }
                        className="max-w-[300px]"
                      />
                    </Form.Label>
                    <Form.Control>
                      <DropdownMenuWithSearch
                        isMultiselect
                        disabled={!!isEdit || disabled}
                        hidden={isOpenCreateGoalModal}
                        isLoading={isLoadingGoals}
                        placeholder={t(`experiments.select-goal`)}
                        label={
                          field.value
                            ?.map(
                              item =>
                                goalOptions.find(goal => goal.value === item)
                                  ?.label
                            )
                            ?.join(', ') || ''
                        }
                        options={goalOptions}
                        selectedOptions={field.value as string[]}
                        createNewOption={
                          disabled ? undefined : (
                            <CreateNewOptionButton
                              text={t('common:create-a-new-goal')}
                              onClick={onOpenCreateGoalModal}
                            />
                          )
                        }
                        onSelectOption={value => {
                          const isExisted = field.value?.find(
                            item => item === value
                          );
                          field.onChange(
                            isExisted
                              ? field.value?.filter(item => item !== value)
                              : [...field.value, value]
                          );
                        }}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button type="button" variant="secondary" onClick={onClose}>
                      {t(`common:cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <DisabledButtonTooltip
                      hidden={editable}
                      trigger={
                        <Button
                          type="submit"
                          disabled={!isDirty || !editable}
                          loading={isSubmitting}
                        >
                          {t(`common:submit`)}
                        </Button>
                      }
                    />
                  }
                />
              </div>
            </Form>
          </FormProvider>
        </div>
      )}
      {isOpenCreateGoalModal && (
        <CreateGoalModal
          isOpen={isOpenCreateGoalModal}
          onClose={onHiddenCreateGoalModal}
          onCompleted={goal => {
            form.setValue('goalIds', [
              ...(form.getValues('goalIds') || []),
              goal.id
            ]);
          }}
        />
      )}
      {isOpenCreateFlagModal && (
        <DialogModal
          className="max-w-[850px] w-full h-full max-h-[90vh] overflow-hidden"
          title={t('common:new-flag')}
          isOpen={isOpenCreateFlagModal}
          onClose={onHiddenCreateFlagModal}
        >
          <CreateFlagForm
            className={
              'flex flex-col flex-1 h-full overflow-auto small-scroll max-h-[90vh] pb-[170px]'
            }
            onClose={onHiddenCreateFlagModal}
            onCompleted={flag => {
              form.setValue('featureId', flag.id);
            }}
          />
        </DialogModal>
      )}
    </SlideModal>
  );
};

export default ExperimentCreateUpdateModal;
