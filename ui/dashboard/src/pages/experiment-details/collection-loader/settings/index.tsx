import { useMemo } from 'react';
import {
  ControllerRenderProps,
  FormProvider,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import { experimentUpdater, ExperimentUpdaterParams } from '@api/experiment';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateExperimentDetails } from '@queries/experiment-details';
import { invalidateExperiments } from '@queries/experiments';
import { useQueryFeatures } from '@queries/features';
import { useQueryGoals } from '@queries/goals';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { IconInfo } from '@icons';
import { createExperimentFormSchema } from 'pages/experiments/form-schema';
import Button from 'components/button';
import { CreatableSelect } from 'components/creatable-select';
import { ReactDatePicker } from 'components/date-time-picker';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import VariationLabel from 'elements/variation-label';

export interface ExperimentSettingsForm {
  id?: string;
  baseVariationId: string;
  name: string;
  description?: string;
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
  ExperimentSettingsForm,
  'audience'
>;

const ExperimentSettings = ({ experiment }: { experiment: Experiment }) => {
  const { t } = useTranslation(['form', 'common', 'table', 'message']);
  const { notify, errorNotify } = useToast();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  const queryClient = useQueryClient();

  const isEnabledEdit = useMemo(
    () => ['WAITING', 'NOT_STARTED'].includes(experiment?.status as string),
    [experiment]
  );

  const { data: goalCollection, isLoading: isLoadingGoals } = useQueryGoals({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
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

  const { data: featureCollection } = useQueryFeatures({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
    }
  });

  const featureFlagOptions = (featureCollection?.features || []).map(
    feature => {
      return {
        value: feature.id,
        label: feature.name,
        enabled: feature.enabled,
        variations: feature.variations
      };
    }
  );

  const form = useForm({
    resolver: yupResolver(useFormSchema(createExperimentFormSchema)),
    defaultValues: {
      id: experiment.id,
      name: experiment.name,
      baseVariationId: experiment.baseVariationId,
      description: experiment.description,
      startAt: experiment.startAt,
      stopAt: experiment.stopAt,
      startType: 'manual',
      audience: {
        rule: '',
        inExperiment: 5,
        notInExperiment: 95,
        served: true,
        variationReassignment: false
      },
      featureId: experiment.featureId,
      goalIds: experiment.goalIds
    },
    mode: 'onChange'
  });

  const {
    watch,
    formState: { isDirty, isValid }
  } = form;

  const featureId = watch('featureId');

  const variationOptions =
    featureFlagOptions
      ?.find(item => item.value === featureId)
      ?.variations.map((variation, index) => ({
        label: (
          <VariationLabel
            label={variation.name || variation.value}
            index={index}
          />
        ),
        value: variation.id
      })) || [];

  const onSubmit: SubmitHandler<ExperimentSettingsForm> = async values => {
    const { id, name, description, startAt, stopAt } = values;
    return onUpdateExperiment({
      id,
      name,
      description,
      startAt,
      stopAt,
      environmentId: currentEnvironment.id
    });
  };

  const mutationState = useMutation({
    mutationFn: async (params: ExperimentUpdaterParams) => {
      return experimentUpdater(params);
    },
    onSuccess: data => {
      invalidateExperimentDetails(queryClient, {
        id: data.experiment.id,
        environmentId: currentEnvironment.id
      });
      invalidateExperiments(queryClient);
      notify({
        message: t('message:collection-action-success', {
          collection: t('common:source-type.experiment'),
          action: t('common:updated')
        })
      });
      mutationState.reset();

      form.reset({
        ...form.getValues(),
        name: data?.experiment?.name,
        description: data?.experiment?.description
      });
    },
    onError: error => errorNotify(error)
  });

  const onUpdateExperiment = async (payload: ExperimentUpdaterParams) =>
    mutationState.mutate(payload);

  return (
    <div className="flex flex-col w-full gap-y-6">
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col w-full gap-y-6">
            <div className="flex items-center w-full justify-between">
              <p className="text-gray-800 typo-head-bold-small">
                {t('common:settings')}
              </p>
              <DisabledButtonTooltip
                hidden={editable}
                trigger={
                  <Button
                    type="submit"
                    disabled={!isDirty || !isValid || !editable}
                    loading={mutationState.isPending}
                  >
                    {t('common:save')}
                  </Button>
                }
              />
            </div>

            <div className="flex flex-col w-full gap-y-5 p-5 shadow-card rounded-lg bg-white">
              <p className="text-gray-800 typo-head-bold-small">
                {t('general-info')}
              </p>
              <Form.Field
                control={form.control}
                name="name"
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Label required>{t('common:name')}</Form.Label>
                    <Form.Control>
                      <Input
                        placeholder={`${t('placeholder-name')}`}
                        disabled={!editable}
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
                  <Form.Item className="py-0">
                    <Form.Label optional>{t('description')}</Form.Label>
                    <Form.Control>
                      <TextArea
                        placeholder={t('placeholder-desc')}
                        rows={4}
                        disabled={!editable}
                        {...field}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <div className="flex items-center w-full gap-x-4">
                <Form.Field
                  control={form.control}
                  name="startAt"
                  render={({ field }) => (
                    <Form.Item className="flex flex-col flex-1 h-full self-stretch py-0">
                      <Form.Label required>{t('start-at')}</Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          disabled={!isEnabledEdit || !editable}
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
                    <Form.Item className="flex flex-col flex-1 h-full self-stretch py-0">
                      <Form.Label required>{t('experiments.time')}</Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          disabled={!isEnabledEdit || !editable}
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
                    <Form.Item className="flex flex-col flex-1 h-full self-stretch py-0">
                      <Form.Label required>{t('end-at')}</Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          disabled={!isEnabledEdit || !editable}
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
                    <Form.Item className="flex flex-col flex-1 h-full self-stretch py-0">
                      <Form.Label required>{t('experiments.time')}</Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          disabled={!isEnabledEdit || !editable}
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
            </div>
            <div className="flex flex-col w-full gap-y-5 p-5 shadow-card rounded-lg bg-white">
              <p className="text-gray-800 typo-head-bold-small">
                {t('common:flag')}
              </p>
              <Form.Field
                control={form.control}
                name={`featureId`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 overflow-hidden py-0">
                    <Form.Label required className="relative w-fit">
                      {t('common:flag')}
                    </Form.Label>
                    <Form.Control>
                      <Dropdown
                        disabled
                        placeholder={t(`experiments.select-flag`)}
                        contentClassName="min-w-[502px]"
                        className="w-full"
                        options={featureFlagOptions as DropdownOption[]}
                        value={field.value}
                        onChange={field.onChange}
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
                    <Form.Item className="flex flex-col flex-1 overflow-hidden py-0">
                      <Form.Label required>
                        {t('experiments.base-variation')}
                      </Form.Label>
                      <Form.Control>
                        <Dropdown
                          disabled
                          placeholder={t(`experiments.select-flag`)}
                          options={variationOptions as DropdownOption[]}
                          value={field.value}
                          onChange={field.onChange}
                          contentClassName="min-w-[502px]"
                          className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              )}
            </div>
            <div className="flex flex-col w-full gap-y-5 p-5 shadow-card rounded-lg bg-white">
              <p className="text-gray-800 typo-head-bold-small">
                {t('common:goals')}
              </p>
              <Form.Field
                control={form.control}
                name={`goalIds`}
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Label required className="relative w-fit">
                      {t('common:goals')}
                      <Tooltip
                        align="start"
                        alignOffset={-50}
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
                      <Form.Control>
                        <CreatableSelect
                          disabled
                          loading={isLoadingGoals}
                          value={goalOptions.filter(item =>
                            field.value.includes(item.value)
                          )}
                          placeholder={t(`experiments.select-goal`)}
                          options={goalOptions?.map(goal => ({
                            label: goal.label,
                            value: goal.value
                          }))}
                          onChange={value =>
                            field.onChange(value.map(goal => goal.value))
                          }
                          onCreateOption={() => {}}
                        />
                      </Form.Control>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            </div>
          </div>
        </Form>
      </FormProvider>
    </div>
  );
};

export default ExperimentSettings;
