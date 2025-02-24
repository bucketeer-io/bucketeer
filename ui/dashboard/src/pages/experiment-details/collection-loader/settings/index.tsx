import { useMemo } from 'react';
import {
  ControllerRenderProps,
  FormProvider,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import { Trans } from 'react-i18next';
import { experimentUpdater, ExperimentUpdaterParams } from '@api/experiment';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateExperimentDetails } from '@queries/experiment-details';
import { invalidateExperiments } from '@queries/experiments';
import { useQueryGoals } from '@queries/goals';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import {
  booleanVariations,
  flagOptions,
  jsonVariations,
  numberVariations,
  stringVariations
} from 'pages/experiments/experiments-modal/mocks';
import { experimentFormSchema } from 'pages/experiments/form-schema';
import GoalActions from 'pages/goal-details/elements/goal-actions';
import Button from 'components/button';
import { ReactDatePicker } from 'components/date-time-picker';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import InfoMessage from 'components/info-message';
import Input from 'components/input';
import TextArea from 'components/textarea';
import ConfirmModal from 'elements/confirm-modal';

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
  const { t } = useTranslation(['form', 'common', 'table']);
  const { notify } = useToast();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const { data: goalCollection } = useQueryGoals({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
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

  const form = useForm({
    resolver: yupResolver(experimentFormSchema),
    defaultValues: {
      id: experiment.id,
      name: experiment.name,
      baseVariationId: experiment.baseVariationId,
      description: experiment.description,
      startAt: experiment.startAt,
      stopAt: experiment.stopAt,
      audience: {
        rule: '',
        inExperiment: 5,
        notInExperiment: 95,
        served: true,
        variationReassignment: false
      },
      featureId: experiment.featureId,
      goalIds: experiment.goalIds
    }
  });

  const {
    watch,
    formState: { isDirty, isSubmitting }
  } = form;

  const featureId = watch('featureId');
  const isStringVariation = featureId.includes('string');
  const isNumberVariation = featureId.includes('number');
  const isBooleanVariation = featureId.includes('boolean');

  const variationOptions = isStringVariation
    ? stringVariations
    : isNumberVariation
      ? numberVariations
      : isBooleanVariation
        ? booleanVariations
        : jsonVariations;

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
      onCloseConfirmModal();
      invalidateExperimentDetails(queryClient, {
        id: data.experiment.id,
        environmentId: currentEnvironment.id
      });
      invalidateExperiments(queryClient);
      notify({
        message: (
          <span>
            <b>{data?.experiment?.name}</b> {`has been successfully updated!`}
          </span>
        )
      });
      mutationState.reset();
    },
    onError: error =>
      notify({
        messageType: 'error',
        message: error?.message || 'Something went wrong.'
      })
  });

  const onUpdateExperiment = async (payload: ExperimentUpdaterParams) =>
    mutationState.mutate(payload);

  return (
    <div className="p-5">
      <p className="text-gray-800 typo-head-bold-small">{t('general-info')}</p>
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <Form.Field
            control={form.control}
            name="name"
            render={({ field }) => (
              <Form.Item className="py-2.5">
                <Form.Label required>{t('common:name')}</Form.Label>
                <Form.Control>
                  <Input placeholder={`${t('placeholder-name')}`} {...field} />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="description"
            render={({ field }) => (
              <Form.Item className="py-2.5">
                <Form.Label optional>{t('description')}</Form.Label>
                <Form.Control>
                  <TextArea
                    placeholder={t('placeholder-desc')}
                    rows={4}
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
                <Form.Item className="flex flex-col flex-1 py-2.5 h-full self-stretch">
                  <Form.Label required>{t('start-at')}</Form.Label>
                  <Form.Control>
                    <ReactDatePicker
                      disabled
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
              name="stopAt"
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 py-2.5 h-full self-stretch">
                  <Form.Label required>{t('end-at')}</Form.Label>
                  <Form.Control>
                    <ReactDatePicker
                      disabled
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
          <p className="text-gray-800 typo-head-bold-small mt-5 mb-2.5">
            {t('link')}
          </p>
          <div className="flex items-center w-full gap-x-4">
            <Form.Field
              control={form.control}
              name={`featureId`}
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 overflow-hidden py-2.5">
                  <Form.Label required>{t('common:flag')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        disabled
                        placeholder={t(`experiments.select-flag`)}
                        label={
                          flagOptions.find(item => item.value === field.value)
                            ?.label || ''
                        }
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {flagOptions.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.value}
                            label={item.label}
                            onSelectOption={value => {
                              field.onChange(value);
                            }}
                          />
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
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
                  <Form.Item className="flex flex-col flex-1 overflow-hidden py-2.5">
                    <Form.Label required>
                      {t('experiments.base-variation')}
                    </Form.Label>
                    <Form.Control>
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          disabled
                          placeholder={t(`experiments.select-flag`)}
                          label={
                            variationOptions.find(
                              item => item.value === field.value
                            )?.label || ''
                          }
                          variant="secondary"
                          className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                        />
                        <DropdownMenuContent
                          className="w-[502px]"
                          align="start"
                          {...field}
                        >
                          {variationOptions.map((item, index) => (
                            <DropdownMenuItem
                              {...field}
                              key={index}
                              value={item.value}
                              label={item.label}
                              onSelectOption={value => {
                                field.onChange(value);
                              }}
                            />
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
          </div>

          <Form.Field
            control={form.control}
            name={`goalIds`}
            render={({ field }) => (
              <Form.Item className="py-2.5">
                <Form.Label required>{t('experiments.link-goal')}</Form.Label>
                <Form.Control>
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      placeholder={t(`experiments.select-goal`)}
                      label={field.value
                        .map(
                          item =>
                            goalOptions.find(opt => opt.value === item)?.label
                        )
                        .join(', ')}
                      variant="secondary"
                      className="w-full"
                    />
                    <DropdownMenuContent
                      className="w-[502px]"
                      align="start"
                      {...field}
                    >
                      {goalOptions.map((item, index) => (
                        <DropdownMenuItem
                          {...field}
                          isMultiselect
                          isSelected={field.value.includes(item.value)}
                          key={index}
                          value={item.value}
                          label={item.label}
                          onSelectOption={value => {
                            const newValue = field.value.includes(value)
                              ? field.value.filter(item => item !== value)
                              : [...field.value, value];
                            field.onChange(newValue);
                          }}
                        />
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Button
            className="mt-2.5 mb-5"
            variant="secondary"
            disabled={!isDirty}
            loading={isSubmitting}
          >
            {t('common:save')}
          </Button>
          <GoalActions
            title={
              experiment.archived
                ? t(`table:popover.unarchive-experiment`)
                : t(`table:popover.archive-experiment`)
            }
            description={
              experiment?.archived ? '' : t('form:experiments.archive-desc')
            }
            btnText={
              experiment.archived
                ? t(`table:popover.unarchive-experiment`)
                : t(`table:popover.archive-experiment`)
            }
            onClick={onOpenConfirmModal}
            disabled={experiment.status === 'RUNNING'}
          >
            {experiment.status === 'RUNNING' && (
              <InfoMessage
                description={t('form:experiments.archive-warning-desc')}
              />
            )}
          </GoalActions>
          {/* <Form.Field
            control={form.control}
            name={`audience`}
            render={({ field }) => (
              <Form.Item className="flex flex-col w-full py-0">
                <DefineAudience field={field as DefineAudienceField} />
              </Form.Item>
            )}
          /> */}
        </Form>
      </FormProvider>
      {openConfirmModal && (
        <ConfirmModal
          isOpen={openConfirmModal}
          loading={mutationState.isPending}
          title={
            experiment.archived
              ? t(`table:popover.unarchive-experiment`)
              : t(`table:popover.archive-experiment`)
          }
          description={
            <Trans
              i18nKey={
                experiment.archived
                  ? 'table:experiment.confirm-unarchive-desc'
                  : 'table:experiment.confirm-archive-desc'
              }
              values={{ name: experiment?.name }}
              components={{ bold: <strong /> }}
            />
          }
          onClose={onCloseConfirmModal}
          onSubmit={() =>
            onUpdateExperiment({
              id: experiment.id,
              environmentId: currentEnvironment.id,
              archived: !experiment.archived
            })
          }
        />
      )}
    </div>
  );
};

export default ExperimentSettings;
