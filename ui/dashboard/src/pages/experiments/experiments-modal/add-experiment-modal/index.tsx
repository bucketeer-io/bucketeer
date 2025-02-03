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

export const formSchema = yup.object().shape({
  baseVariationId: yup.string().required(),
  name: yup.string().required(),
  startAt: yup.string().required(),
  stopAt: yup.string().required(),
  description: yup.string(),
  audience: yup.mixed(),
  featureId: yup.string().required(),
  goalIds: yup.array().min(1).required()
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
    resolver: yupResolver(formSchema),
    defaultValues: {
      baseVariationId: '',
      name: '',
      description: '',
      startAt: '1602829513',
      stopAt: '1737508939',
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
    formState: { isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<AddExperimentForm> = async values => {
    try {
      const { audience, ...rest } = values;
      console.log({ audience });
      const resp = await experimentCreator({
        ...rest,
        featureId: 'test',
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
      console.log(error)
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
                  <Form.Item className="flex flex-col flex-1">
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
                  <Form.Item className="flex flex-col flex-1">
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
                        label={
                          goalOptions.find(
                            item => item.value === field.value[0]
                          )?.label
                        }
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
                            key={index}
                            value={item.value}
                            label={item.label}
                            onSelectOption={value => {
                              field.onChange([value]);
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
                    disabled={!isValid}
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
