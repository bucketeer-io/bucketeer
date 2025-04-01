import { useMemo } from 'react';
import {
  ControllerRenderProps,
  FormProvider,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import { experimentCreator } from '@api/experiment';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateExperiments } from '@queries/experiments';
import { useQueryFeatures } from '@queries/features';
import { useQueryGoals } from '@queries/goals';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconInfo, IconPlus } from '@icons';
import { experimentFormSchema } from 'pages/experiments/form-schema';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { CreatableSelect } from 'components/creatable-select';
import { ReactDatePicker } from 'components/date-time-picker';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';
import CreateGoalModal from 'elements/create-goal-modal';

interface AddExperimentModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export type StartType = 'manual' | 'schedule';

export interface AddExperimentForm {
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
  AddExperimentForm,
  'audience'
>;

const FeatureFlagStatus = ({
  status,
  enabled
}: {
  status: string;
  enabled: boolean;
}) => {
  return (
    <div
      className={cn(
        'flex-center py-0.5 px-2 rounded-lg typo-para-small text-white bg-primary-500 border border-gray-300',
        {
          'text-gray-700 bg-gray-200': !enabled
        }
      )}
    >
      {status}
    </div>
  );
};

const AddExperimentModal = ({ isOpen, onClose }: AddExperimentModalProps) => {
  const { t } = useTranslation(['form', 'common']);
  const { notify } = useToast();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();

  const [
    isOpenCreateGoalModal,
    onOpenCreateGoalModal,
    onHiddenCreateGoalModal
  ] = useToggleOpen(false);

  const { data: goalCollection, isLoading: isLoadingGoals } = useQueryGoals({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      archived: false,
      connectionType: 'EXPERIMENT'
    }
  });

  const { data: featureCollection, isLoading: isLoadingFeature } =
    useQueryFeatures({
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
      baseVariationId: '',
      name: '',
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
    }
  });

  const {
    watch,
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  const featureId = watch('featureId');

  const variationOptions =
    featureFlagOptions?.find(item => item.value === featureId)?.variations ||
    [];

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

  const onSubmit: SubmitHandler<AddExperimentForm> = async values => {
    try {
      const {
        baseVariationId,
        featureId,
        goalIds,
        name,
        startAt,
        stopAt,
        description
      } = values;
      const resp = await experimentCreator({
        baseVariationId,
        featureId,
        goalIds,
        name,
        startAt,
        stopAt,
        description,
        environmentId: currentEnvironment.id
      });
      if (resp) {
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: 'Experiment created successfully.'
        });
        invalidateExperiments(queryClient);
        onClose();
      }
    } catch (error) {
      const errorMessage = (error as Error)?.message;
      notify({
        toastType: 'toast',
        messageType: 'error',
        message: errorMessage || 'Something went wrong.'
      });
    }
  };

  return (
    <SlideModal
      title={t('common:new-experiment')}
      isOpen={isOpen}
      onClose={onClose}
    >
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
                      placeholder={`${t('placeholder-name')}`}
                      {...field}
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
                <Form.Item className="flex flex-col w-full overflow-hidden">
                  <Form.Label required className="relative w-fit">
                    {t('common:flag')}
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
                        placeholder={t(`experiments.select-flag`)}
                        label={
                          featureFlagOptions.find(
                            item => item.value === field.value
                          )?.label || ''
                        }
                        disabled={isLoadingFeature}
                        variant="secondary"
                        className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {featureFlagOptions.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.value}
                            label={item.label}
                            className="justify-between gap-x-4"
                            additionalElement={
                              <FeatureFlagStatus
                                status={t(
                                  item.enabled
                                    ? 'experiments.on'
                                    : 'experiments.off'
                                )}
                                enabled={item.enabled}
                              />
                            }
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
                  <Form.Item className="flex flex-col w-full overflow-hidden">
                    <Form.Label required>
                      {t('experiments.base-variation')}
                    </Form.Label>
                    <Form.Control>
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          placeholder={t(`experiments.select-variation`)}
                          label={
                            variationOptions?.find(
                              item => item.id === field.value
                            )?.name || ''
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
                              value={item.id}
                              label={item.name}
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
                    <Icon
                      icon={IconInfo}
                      size="xs"
                      color="gray-500"
                      className="absolute -right-6"
                    />
                  </Form.Label>
                  <Form.Control>
                    <CreatableSelect
                      disabled={isLoadingGoals}
                      loading={isLoadingGoals}
                      placeholder={t(`experiments.select-goal`)}
                      options={goalOptions?.map(goal => ({
                        label: goal.label,
                        value: goal.value
                      }))}
                      // const formatOptionLabel = ({ label, enabled }) => {
                      //   return (
                      //     <div className="flex justify-between space-x-4 pr-2">
                      //       <span className="flex-1 truncate">{label}</span>
                      //       <span
                      //         className={`border rounded-lg text-sm w-11 flex justify-center ${
                      //           enabled
                      //             ? 'bg-primary border-primary text-white'
                      //             : 'bg-gray-100 border-gray-300'
                      //         }`}
                      //       >
                      //         {enabled ? 'On' : 'Off'}
                      //       </span>
                      //     </div>
                      //   );
                      // };
                      onChange={value =>
                        field.onChange(value.map(goal => goal.value))
                      }
                      onCreateOption={onOpenCreateGoalModal}
                      formatCreateLabel={() => (
                        <Button
                          variant="text"
                          className="h-6 self-center w-full"
                        >
                          <Icon
                            icon={IconPlus}
                            color="primary-500"
                            size={'xs'}
                          />
                          {t('common:create-a-new-goal')}
                        </Button>
                      )}
                      noOptionsMessage={() => (
                        <p>{t('common:no-options-found')}</p>
                      )}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            {/* <Divider className="mt-4 mb-5" />
            <div className="mb-3">
              <p className="text-gray-800 typo-head-bold-small mb-3">
                {t('experiments.define-audience.title')}
              </p>
              <Form.Field
                control={form.control}
                name={`audience`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col w-full py-2 gap-y-5">
                    <DefineAudience field={field as DefineAudienceField} />
                  </Form.Item>
                )}
              />
            </div> */}

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
                    disabled={!isValid || !isDirty}
                    loading={isSubmitting}
                  >
                    {t(`common:submit`)}
                  </Button>
                }
              />
            </div>
          </Form>
        </FormProvider>
      </div>
      {isOpenCreateGoalModal && (
        <CreateGoalModal
          isOpen={isOpenCreateGoalModal}
          onClose={onHiddenCreateGoalModal}
        />
      )}
    </SlideModal>
  );
};

export default AddExperimentModal;
