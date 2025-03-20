import { useCallback } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { scheduleFlagCreator } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, useAuth } from 'auth';
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

export type ConfirmationRequiredModalProps = {
  selectedFlag: Feature;
  isOpen: boolean;
  isEnabling: boolean;
  onClose: () => void;
};

export interface ScheduleFlagForm {
  id?: string;
  scheduleType: 'ACTIVE_NOW' | 'SCHEDULE';
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

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  // const { data: collection } = useQueryScheduleFlags({
  //   params: {
  //     environmentId: currentEnvironment?.id,
  //     featureId: selectedFlag.id
  //   }
  // });

  const formSchema = yup.object().shape({
    id: yup.string(),
    comment: currentEnvironment?.requireComment
      ? yup.string().required()
      : yup.string(),
    scheduleType: yup.string().oneOf(['ACTIVE_NOW', 'SCHEDULE']).required(),
    scheduleAt: yup
      .string()
      .required()
      .test('test', function (value, context) {
        const scheduleType = context.from && context.from[0].value.scheduleType;
        if (scheduleType === 'SCHEDULE' && !value) {
          return context.createError({
            message: `This field is required.`,
            path: context.path
          });
        }
        return true;
      })
  });

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: '',
      comment: '',
      scheduleType: 'ACTIVE_NOW',
      scheduleAt: String(Math.floor(new Date().getTime() / 1000))
    }
  });

  const {
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<ScheduleFlagForm> = useCallback(
    async values => {
      try {
        console.log(values);
        const resp = await scheduleFlagCreator({
          environmentId: currentEnvironment.id,
          featureId: selectedFlag.id,
          scheduledAt: values.scheduleAt,
          scheduledChanges: []
        });
        console.log({ resp });
      } catch (error) {
        console.log(error);
      }
    },
    [isEnabling, selectedFlag, currentEnvironment]
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
                    <Form.Label required={currentEnvironment?.requireComment}>
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
                          form.setValue(
                            'scheduleAt',
                            String(Math.floor(new Date().getTime() / 1000))
                          );
                        }}
                      >
                        <div className="flex items-center gap-x-2">
                          <RadioGroupItem id="active_now" value="ACTIVE_NOW" />
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
                <Form.Field
                  control={form.control}
                  name="scheduleAt"
                  render={({ field }) => (
                    <Form.Item className="py-0 mt-5">
                      <Form.Label required>
                        {t('form:feature-flags:start-at')}
                      </Form.Label>
                      <Form.Control>
                        <ReactDatePicker
                          selected={
                            field.value ? new Date(+field.value * 1000) : null
                          }
                          onChange={date => field.onChange(date)}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              )}
            </div>
            <ButtonBar
              secondaryButton={
                <Button loading={isSubmitting} disabled={!isDirty || !isValid}>
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
