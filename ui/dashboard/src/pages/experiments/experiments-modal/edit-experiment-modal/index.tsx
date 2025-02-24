import { useEffect, useMemo } from 'react';
import {
  ControllerRenderProps,
  FormProvider,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import { experimentUpdater } from '@api/experiment';
import { yupResolver } from '@hookform/resolvers/yup';
import {
  invalidateExperimentDetails,
  useQueryExperimentDetails
} from '@queries/experiment-details';
import { invalidateExperiments } from '@queries/experiments';
import { useQueryGoals } from '@queries/goals';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useToast } from 'hooks';
import useActionWithURL from 'hooks/use-action-with-url';
import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
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
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import FormLoading from 'elements/form-loading';
import { StartType } from '../add-experiment-modal';
import {
  booleanVariations,
  flagOptions,
  jsonVariations,
  numberVariations,
  stringVariations
} from '../mocks';

interface EditExperimentModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface EditExperimentForm {
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
  EditExperimentForm,
  'audience'
>;

const EditExperimentModal = ({ isOpen, onClose }: EditExperimentModalProps) => {
  const { t } = useTranslation(['form', 'common']);
  const { notify } = useToast();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();

  const { id: experimentId, errorToast } = useActionWithURL({
    closeModalPath: `/${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}`
  });

  const {
    data: experimentCollection,
    isLoading: experimentLoading,
    error: experimentError
  } = useQueryExperimentDetails({
    params: {
      id: experimentId as string,
      environmentId: currentEnvironment.id
    }
  });

  const experiment = useMemo(
    () => experimentCollection?.experiment,
    [experimentCollection]
  );

  const { data: goalCollection, isLoading: isLoadingGoals } = useQueryGoals({
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
    }
  });

  const {
    watch,
    formState: { isDirty, isSubmitting }
  } = form;
  const startType = watch('startType');
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

  const startOptions = [
    {
      label: 'Manual Start',
      value: 'manual'
    },
    {
      label: 'Schedule',
      value: 'schedule'
    }
  ];

  const onSubmit: SubmitHandler<EditExperimentForm> = async values => {
    try {
      const { id, name, description, startAt, stopAt } = values;
      const resp = await experimentUpdater({
        id,
        name,
        description,
        startAt,
        stopAt,
        environmentId: currentEnvironment.id
      });
      if (resp) {
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: 'Experiment updated successfully.'
        });
        invalidateExperiments(queryClient);
        invalidateExperimentDetails(queryClient, {
          id: experimentId as string,
          environmentId: currentEnvironment.id
        });
        onClose();
      }
    } catch (error) {
      errorToast(error as Error);
    }
  };

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
      errorToast(experimentError);
    }
  }, [experimentError]);

  return (
    <SlideModal
      title={t('common:edit-experiment')}
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
              <RadioGroup
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
              </RadioGroup>
              {startType === 'schedule' && (
                <>
                  <div className="flex items-center w-full gap-x-4">
                    <Form.Field
                      control={form.control}
                      name="startAt"
                      render={({ field }) => (
                        <Form.Item className="flex flex-col flex-1 h-full self-stretch">
                          <Form.Label required>{t('start-at')}</Form.Label>
                          <Form.Control>
                            <ReactDatePicker
                              disabled={!!experiment?.startAt}
                              dateFormat={'yyyy/MM/dd'}
                              showTimeSelect={false}
                              selected={
                                field.value
                                  ? new Date(+field.value * 1000)
                                  : null
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
                          <Form.Label required>
                            {t('experiments.time')}
                          </Form.Label>
                          <Form.Control>
                            <ReactDatePicker
                              disabled={!!experiment?.startAt}
                              dateFormat={'HH:mm'}
                              showTimeSelect
                              showTimeSelectOnly={true}
                              selected={
                                field.value
                                  ? new Date(+field.value * 1000)
                                  : null
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
                              disabled={!!experiment?.stopAt}
                              dateFormat={'yyyy/MM/dd'}
                              showTimeSelect={false}
                              selected={
                                field.value
                                  ? new Date(+field.value * 1000)
                                  : null
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
                          <Form.Label required>
                            {t('experiments.time')}
                          </Form.Label>
                          <Form.Control>
                            <ReactDatePicker
                              disabled={!!experiment?.stopAt}
                              dateFormat={'HH:mm'}
                              showTimeSelect
                              showTimeSelectOnly={true}
                              selected={
                                field.value
                                  ? new Date(+field.value * 1000)
                                  : null
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
                </>
              )}
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
                          disabled
                          placeholder={t(`experiments.select-flag`)}
                          label={
                            flagOptions.find(item => item.value === field.value)
                              ?.label || ''
                          }
                          variant="secondary"
                          className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
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
                    <Form.Item className="flex flex-col w-full overflow-hidden">
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
                    <Form.Message />
                  </Form.Item>
                )}
              />

              {/* <Divider className="mt-4 mb-5" />
              <div>
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
                      disabled={!isDirty}
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
      )}
    </SlideModal>
  );
};

export default EditExperimentModal;
