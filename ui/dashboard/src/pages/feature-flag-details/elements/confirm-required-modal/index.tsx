import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { cn } from 'utils/style';
import { IconInfo, IconToastWarning } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import { ReactDatePicker } from 'components/date-time-picker';
import Form from 'components/form';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import ToastMessage from 'components/toast';
import { Tooltip } from 'components/tooltip';
import { formSchema } from './form-schema';

export type ConfirmationRequiredModalProps = {
  feature: Feature;
  isOpen: boolean;
  isShowScheduleSelect?: boolean;
  isShowRolloutWarning?: boolean;
  onClose: () => void;
  onSubmit: (args: { comment?: string; scheduleAt?: string }) => Promise<void>;
};

export interface ConfirmRequiredValues {
  resetSampling?: boolean;
  comment?: string;
  requireComment?: boolean;
  scheduleType?: string;
  scheduleAt?: string;
}

const ConfirmationRequiredModal = ({
  feature,
  isOpen,
  isShowScheduleSelect,
  isShowRolloutWarning,
  onClose,
  onSubmit
}: ConfirmationRequiredModalProps) => {
  const { t } = useTranslation(['common', 'form', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const navigate = useNavigate();

  const { data: rolloutCollection } = useQueryRollouts({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      featureIds: [feature.id]
    },
    refetchOnMount: !!feature && isShowRolloutWarning ? 'always' : false,
    enabled: !!feature && isShowRolloutWarning
  });

  const hasRolloutRunning = !!rolloutCollection?.progressiveRollouts?.find(
    item => ['WAITING', 'RUNNING'].includes(item.status)
  );

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      comment: '',
      resetSampling: false,
      requireComment: currentEnvironment.requireComment,
      scheduleType: feature.enabled ? 'DISABLE' : 'ENABLE',
      scheduleAt: String(Math.floor((new Date().getTime() + 3600000) / 1000))
    },
    mode: 'onChange'
  });

  const {
    control,
    formState: { isDirty, isValid, isSubmitting },
    watch,
    setValue
  } = form;
  const isRequireComment = watch('requireComment');
  const isShowSchedule = watch('scheduleType') === 'SCHEDULE';

  const handleOnSubmit = async (values: ConfirmRequiredValues) => {
    await onSubmit(values);
  };

  return (
    <DialogModal
      className="w-full max-w-[350px] sm:max-w-[500px]"
      title={t('table:feature-flags.confirm-required')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(handleOnSubmit)}>
          <div className="flex flex-col w-full gap-y-5 items-start pt-5">
            <div className="typo-para-small text-gray-600 w-full px-5">
              {t('table:feature-flags.confirm-required-desc')}
            </div>

            <div className="flex flex-col w-full px-5 pb-5">
              {!isShowSchedule && (
                <>
                  <Form.Field
                    control={control}
                    name="comment"
                    render={({ field }) => (
                      <Form.Item className="py-0">
                        <Form.Label required={isRequireComment}>
                          {t('form:comment-for-update')}
                        </Form.Label>
                        <Form.Control>
                          <TextArea
                            placeholder={`${t('form:placeholder-comment')}`}
                            rows={3}
                            {...field}
                            onChange={value => {
                              field.onChange(value);
                            }}
                            name="comment"
                          />
                        </Form.Control>
                        <Form.Message />
                      </Form.Item>
                    )}
                  />
                  <Form.Field
                    control={control}
                    name="resetSampling"
                    render={({ field }) => (
                      <Form.Item className="flex flex-col w-full py-0 gap-y-4 mt-5">
                        <div className="flex items-center gap-x-2">
                          <Form.Control>
                            <Checkbox
                              ref={field.ref}
                              checked={field.value}
                              onCheckedChange={checked =>
                                field.onChange(checked)
                              }
                              title={t('form:reset-sampling')}
                            />
                          </Form.Control>
                          <Tooltip
                            align="start"
                            content={t('form:reset-sampling-tooltip')}
                            trigger={
                              <div className="flex-center size-fit">
                                <Icon
                                  icon={IconInfo}
                                  size="xs"
                                  color="gray-500"
                                />
                              </div>
                            }
                            className="max-w-[400px]"
                          />
                        </div>
                        <Form.Message />
                      </Form.Item>
                    )}
                  />
                </>
              )}

              {isShowScheduleSelect && (
                <>
                  <Form.Field
                    control={form.control}
                    name="scheduleType"
                    render={({ field }) => (
                      <Form.Item
                        className={cn(
                          'flex flex-col w-full py-0 gap-y-4 mt-5',
                          {
                            'mt-0': isShowSchedule
                          }
                        )}
                      >
                        <Form.Control>
                          <RadioGroup
                            defaultValue={field.value}
                            className="flex flex-col w-full gap-y-4"
                            onValueChange={value => {
                              field.onChange(value);
                              setValue('requireComment', value !== 'SCHEDULE');
                            }}
                          >
                            <div className="flex items-center gap-x-2">
                              <RadioGroupItem
                                id="active_now"
                                value={feature.enabled ? 'DISABLE' : 'ENABLE'}
                              />
                              <label
                                htmlFor="active_now"
                                className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
                              >
                                {t('update-now')}
                              </label>
                            </div>

                            <div className="flex items-center gap-x-2">
                              <RadioGroupItem id="schedule" value="SCHEDULE" />
                              <label
                                htmlFor="schedule"
                                className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
                              >
                                {t('form:feature-flags.schedule')}
                              </label>
                            </div>
                          </RadioGroup>
                        </Form.Control>
                        <Form.Message />
                      </Form.Item>
                    )}
                  />
                  {form.watch('scheduleType') === 'SCHEDULE' && (
                    <div className="flex flex-col w-full gap-y-5 mt-5">
                      <ToastMessage
                        toastType="info-message"
                        messageType="info"
                        message={t('form:feature-flags.schedule-info')}
                      />
                      <Form.Field
                        control={form.control}
                        name="scheduleAt"
                        render={({ field }) => (
                          <Form.Item className="py-0">
                            <Form.Label required>
                              {t('form:feature-flags:start-at')}
                            </Form.Label>
                            <Form.Control>
                              <ReactDatePicker
                                minDate={new Date()}
                                selected={
                                  field.value
                                    ? new Date(+field.value * 1000)
                                    : null
                                }
                                onChange={date => {
                                  if (date) {
                                    field.onChange(
                                      String(date?.getTime() / 1000)
                                    );
                                  }
                                }}
                              />
                            </Form.Control>
                            <Form.Message />
                          </Form.Item>
                        )}
                      />
                    </div>
                  )}
                </>
              )}
              {isShowRolloutWarning && hasRolloutRunning && (
                <div className="flex w-full gap-x-3 p-4 mt-5 rounded-md bg-accent-yellow-50 typo-para-small">
                  <Icon icon={IconToastWarning} />
                  <p className="w-full typo-para-medium text-accent-yellow-700">
                    <Trans
                      i18nKey={'form:has-rollout-running'}
                      components={{
                        comp: (
                          <span
                            onClick={() =>
                              navigate(
                                `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${feature.id}${PAGE_PATH_FEATURE_AUTOOPS}`
                              )
                            }
                            className="inline-flex text-primary-500 underline whitespace-nowrap cursor-pointer"
                          />
                        )
                      }}
                    />
                  </p>
                </div>
              )}
            </div>
            <ButtonBar
              secondaryButton={
                <Button
                  type="submit"
                  loading={isSubmitting}
                  disabled={(isRequireComment && !isDirty) || !isValid}
                >
                  {t(`submit`)}
                </Button>
              }
              primaryButton={
                <Button type="button" onClick={onClose} variant="secondary">
                  {t(`cancel`)}
                </Button>
              }
            />
          </div>
        </Form>
      </FormProvider>
    </DialogModal>
  );
};

export default ConfirmationRequiredModal;
