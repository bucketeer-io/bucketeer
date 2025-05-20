import { useCallback, useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { autoOpsCreator } from '@api/auto-ops';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures } from '@queries/features';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Feature } from '@types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { ReactDatePicker } from 'components/date-time-picker';
import Form from 'components/form';
import DialogModal from 'components/modal/dialog';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import ToastMessage from 'components/toast';

export type ConfirmationRequiredModalProps = {
  selectedFlag: Feature;
  isOpen: boolean;
  isEnabling: boolean;
  onClose: () => void;
};

export interface ScheduleFlagForm {
  id?: string;
  scheduleType: 'ENABLE' | 'DISABLE' | 'SCHEDULE';
  scheduleAt: string;
  comment?: string;
}

const ConfirmationRequiredModal = ({
  selectedFlag,
  isOpen,
  isEnabling,
  onClose
}: ConfirmationRequiredModalProps) => {
  const { t } = useTranslation(['common', 'form', 'table']);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const isRequireComment = currentEnvironment?.requireComment;

  const isEnabled = useMemo(() => selectedFlag.enabled, [selectedFlag]);

  const formSchema = yup.object().shape({
    id: yup.string(),
    comment: isRequireComment ? yup.string().required() : yup.string(),
    scheduleType: yup
      .string()
      .oneOf(['ENABLE', 'DISABLE', 'SCHEDULE'])
      .required(),
    scheduleAt: yup
      .string()
      .required()
      .test('test', function (value, context) {
        const scheduleType = context.from && context.from[0].value.scheduleType;
        if (scheduleType === 'SCHEDULE') {
          if (!value)
            return context.createError({
              message: `This field is required.`,
              path: context.path
            });
          if (+value * 1000 < new Date().getTime())
            return context.createError({
              message: `This must be later than the current time.`,
              path: context.path
            });
        }
        return true;
      })
  });

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: selectedFlag.id,
      comment: '',
      scheduleType: isEnabled ? 'DISABLE' : 'ENABLE',
      scheduleAt: String(Math.floor((new Date().getTime() + 3600000) / 1000))
    },
    mode: 'onChange'
  });

  const {
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<ScheduleFlagForm> = useCallback(
    async values => {
      try {
        const { scheduleType, id, comment, scheduleAt } = values;
        const isEnable = scheduleType === 'ENABLE';
        const isSchedule = scheduleType === 'SCHEDULE';
        let resp;
        if (['ENABLE', 'DISABLE'].includes(scheduleType)) {
          resp = await featureUpdater({
            environmentId: currentEnvironment.id,
            id,
            enabled: isEnable,
            comment
          });
        } else {
          resp = await autoOpsCreator({
            environmentId: currentEnvironment.id,
            featureId: selectedFlag.id,
            opsType: 'SCHEDULE',
            datetimeClauses: [
              {
                actionType: isEnabled ? 'DISABLE' : 'ENABLE',
                time: scheduleAt
              }
            ]
          });
        }
        if (resp) {
          notify({
            message: (
              <Trans
                i18nKey={
                  isSchedule
                    ? 'form:feature-flags.schedule-configured'
                    : 'form:feature-flags.flag-switch'
                }
                values={{
                  name: selectedFlag.name,
                  state: isEnable ? 'enabled' : 'disabled'
                }}
              />
            )
          });
          invalidateFeatures(queryClient);
          invalidateFeature(queryClient);
          onClose();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [isEnabling, selectedFlag, currentEnvironment, isEnabled]
  );

  return (
    <DialogModal
      className="w-[500px]"
      title={t('table:feature-flags.confirm-required')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full gap-y-5 items-start pt-5">
        <div className="typo-para-small text-gray-600 w-full px-5">
          {t('table:feature-flags.confirm-required-desc')}
        </div>
        <FormProvider {...form}>
          <Form className="w-full" onSubmit={form.handleSubmit(onSubmit)}>
            <div className="flex flex-col w-full px-5 pb-5">
              <Form.Field
                control={form.control}
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
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Field
                control={form.control}
                name="scheduleType"
                render={({ field }) => (
                  <Form.Item className="flex flex-col w-full py-0 gap-y-4 mt-5">
                    <Form.Control>
                      <RadioGroup
                        defaultValue={field.value}
                        className="flex flex-col w-full gap-y-4"
                        onValueChange={value => {
                          field.onChange(value);
                        }}
                      >
                        <div className="flex items-center gap-x-2">
                          <RadioGroupItem
                            id="active_now"
                            value={isEnabled ? 'DISABLE' : 'ENABLE'}
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
                              field.value ? new Date(+field.value * 1000) : null
                            }
                            onChange={date => {
                              if (date) {
                                field.onChange(String(date?.getTime() / 1000));
                                form.trigger('scheduleAt');
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
            </div>
            <ButtonBar
              secondaryButton={
                <Button
                  loading={isSubmitting}
                  disabled={(isRequireComment && !isDirty) || !isValid}
                >
                  {t(`submit`)}
                </Button>
              }
              primaryButton={
                <Button onClick={onClose} variant="secondary">
                  {t(`cancel`)}
                </Button>
              }
            />
          </Form>
        </FormProvider>
      </div>
    </DialogModal>
  );
};

export default ConfirmationRequiredModal;
