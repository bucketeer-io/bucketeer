import { useMemo } from 'react';
import {
  ControllerRenderProps,
  FormProvider,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import { experimentUpdater } from '@api/experiment';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateExperimentDetails } from '@queries/experiment-details';
import { invalidateExperiments } from '@queries/experiments';
import { useQueryGoals } from '@queries/goals';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { IconInfo } from '@icons';
import { experimentFormSchema } from 'pages/experiments/form-schema';
import Button from 'components/button';
import { ReactDatePicker } from 'components/date-time-picker';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import TextArea from 'components/textarea';
import DefineAudience from './define-audience';

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
    formState: { isDirty, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<ExperimentSettingsForm> = async values => {
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
          id: id!,
          environmentId: currentEnvironment.id
        });
      }
    } catch (error) {
      notify({
        toastType: 'toast',
        messageType: 'error',
        message: (error as Error)?.message || 'Something went wrong.'
      });
    }
  };

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
            name="baseVariationId"
            render={({ field }) => (
              <Form.Item className="py-2.5">
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
                    disabled={true}
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
          <Form.Field
            control={form.control}
            name={`audience`}
            render={({ field }) => (
              <Form.Item className="flex flex-col w-full py-0">
                <DefineAudience field={field as DefineAudienceField} />
              </Form.Item>
            )}
          />
          <p className="text-gray-800 typo-head-bold-small mt-5 mb-2.5">
            {t('link')}
          </p>
          <Form.Field
            control={form.control}
            name={`featureId`}
            render={({ field }) => (
              <Form.Item className="py-2.5">
                <Form.Label required>{t('experiments.link-flag')}</Form.Label>
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
          <Form.Field
            control={form.control}
            name={`goalIds`}
            render={({ field }) => (
              <Form.Item className="py-2.5">
                <Form.Label required>{t('experiments.link-goal')}</Form.Label>
                <Form.Control>
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      disabled
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
            className="mt-2.5"
            variant="secondary"
            disabled={!isDirty}
            loading={isSubmitting}
          >
            {t('common:save')}
          </Button>
        </Form>
      </FormProvider>
    </div>
  );
};

export default ExperimentSettings;
