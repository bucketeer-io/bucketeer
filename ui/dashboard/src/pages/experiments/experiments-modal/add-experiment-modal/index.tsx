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
import { useQueryGoals } from '@queries/goals';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { IconInfo } from '@icons';
import {
  EXPERIMENT_DESCRIPTION_MAX_LENGTH,
  EXPERIMENT_GOAL_MAX_LENGTH,
  EXPERIMENT_GOAL_MIN_LENGTH,
  EXPERIMENT_NAME_MAX_LENGTH,
  EXPERIMENT_START_AT_OLDEST_DAYS
} from 'pages/experiments/constants';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
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
import DefineAudience from './define-audience';

interface AddExperimentModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface AddExperimentForm {
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
  AddExperimentForm,
  'audience'
>;

export const experimentFormSchema = yup.object().shape({
  id: yup.string().max(EXPERIMENT_NAME_MAX_LENGTH),
  name: yup.string().required(),
  baseVariationId: yup.string().required(),
  startAt: yup
    .string()
    .required()
    .test(
      'laterThanStartAt',
      `This must be later than or equal to ${EXPERIMENT_START_AT_OLDEST_DAYS} days ago.`,
      function (value) {
        const startDate = new Date(+value * 1000);
        const d = new Date();
        d.setDate(d.getDate() - EXPERIMENT_START_AT_OLDEST_DAYS);
        return startDate >= d;
      }
    ),
  stopAt: yup
    .string()
    .required()
    .test('laterThanStartAt', (value, context) => {
      const endDate = new Date(+value * 1000);
      const startAtValue = context?.from && context?.from[0]?.value?.startAt;
      const startDate = new Date(+startAtValue * 1000);
      const startTime = startDate.getTime();
      const endTime = endDate.getTime();
      if (startTime && endTime && endTime < startTime) {
        return context.createError({
          message: 'Stop at must be later than the start at.',
          path: context.path
        });
      }
      return true;
    })
    .test('lessThanOrEquals30Days', (value, context) => {
      const maxPeriodSeconds = 60 * 60 * 24 * 30;
      const startAtValue = context?.from && context?.from[0]?.value?.startAt;
      const startDate = new Date(+startAtValue * 1000);
      const endDate = new Date(+value * 1000);
      const startTime = startDate.getTime();
      const endTime = endDate.getTime();
      if (endTime / 1000 - startTime / 1000 <= maxPeriodSeconds) {
        return context.createError({
          message: `The period must be less than or equals to ${EXPERIMENT_START_AT_OLDEST_DAYS} days.`,
          path: context.path
        });
      }

      return true;
    }),
  description: yup.string().max(EXPERIMENT_DESCRIPTION_MAX_LENGTH),
  audience: yup.mixed(),
  featureId: yup.string().required(),
  goalIds: yup
    .array()
    .min(EXPERIMENT_GOAL_MIN_LENGTH)
    .max(EXPERIMENT_GOAL_MAX_LENGTH)
    .required()
});

const AddExperimentModal = ({ isOpen, onClose }: AddExperimentModalProps) => {
  const { t } = useTranslation(['form', 'common']);
  const { notify } = useToast();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();

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

  const flagOptions = [
    {
      label: 'Flag 1',
      value: 'flag-1'
    },
    {
      label: 'Flag 2',
      value: 'flag-2'
    }
  ];

  const form = useForm({
    resolver: yupResolver(experimentFormSchema),
    defaultValues: {
      baseVariationId: '',
      name: '',
      description: '',
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
    formState: { isDirty, isValid, isSubmitting }
  } = form;

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
      <div className="p-5 pb-28 relative">
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
              name="baseVariationId"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required className="relative w-fit">
                    {t('experiments.experiment-id')}
                    <Icon
                      icon={IconInfo}
                      className="absolute -right-8"
                      size={'sm'}
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('experiments.placeholder-id')}`}
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
            <div className="flex items-center w-full gap-x-4 mb-3">
              <Form.Field
                control={form.control}
                name="startAt"
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 h-full self-stretch">
                    <Form.Label required>{t('start-at')}</Form.Label>
                    <Form.Control>
                      <ReactDatePicker
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
                  <Form.Item className="flex flex-col flex-1 h-full self-stretch">
                    <Form.Label required>{t('end-at')}</Form.Label>
                    <Form.Control>
                      <ReactDatePicker
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
            <Divider className="mb-5" />
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
            </div>
            <Divider className="mb-5" />
            <p className="text-gray-800 typo-head-bold-small mb-3">
              {t('link')}
            </p>
            <Form.Field
              control={form.control}
              name={`featureId`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('experiments.link-flag')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
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
            <Form.Field
              control={form.control}
              name={`goalIds`}
              render={({ field }) => (
                <Form.Item className="py-2">
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
    </SlideModal>
  );
};

export default AddExperimentModal;
