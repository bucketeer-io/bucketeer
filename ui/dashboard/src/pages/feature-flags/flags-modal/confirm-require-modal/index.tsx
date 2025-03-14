import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import DialogModal from 'components/modal/dialog';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';

export type ConfirmationRequiredModalProps = {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: () => void;
};

export const formSchema = yup.object().shape({
  comment: yup.string().required(),
  publishType: yup.string().required(),
  startAt: yup.string().test('isRequired', function (value, context) {
    const publishType = context.from && context.from[0].value.publishType;
    if (publishType === 'SCHEDULE' && !value) {
      return context.createError({
        message: `This field is required.`,
        path: context.path
      });
    }
    return true;
  })
});

const ConfirmationRequiredModal = ({
  isOpen,
  onClose,
  onSubmit
}: ConfirmationRequiredModalProps) => {
  const { t } = useTranslation(['common', 'form', 'table']);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      comment: '',
      publishType: 'ACTIVE_NOW',
      startAt: ''
    }
  });

  const {
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  return (
    <DialogModal
      className="w-[500px]"
      title={t('table:feature-flags.confirm-required')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full gap-y-5 items-start p-5">
        <div className="typo-para-small text-gray-600 w-full">
          {t('table:feature-flags.confirm-required-desc')}
        </div>
        <FormProvider {...form}>
          <Form className="w-full" onSubmit={form.handleSubmit(onSubmit)}>
            <Form.Field
              control={form.control}
              name="comment"
              render={({ field }) => (
                <Form.Item className="py-0">
                  <Form.Label>{t('form:comment-for-update')}</Form.Label>
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
              name="publishType"
              render={({ field }) => (
                <Form.Item className="flex flex-col w-full py-0 gap-y-4 mt-5">
                  <Form.Control>
                    <RadioGroup
                      defaultValue={field.value}
                      className="flex flex-col w-full gap-y-4"
                      onValueChange={field.onChange}
                    >
                      <div className="flex items-center gap-x-2">
                        <RadioGroupItem id="active_now" value="ACTIVE_NOW" />
                        <label
                          htmlFor="active_now"
                          className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
                        >
                          {t('form:feature-flags.active-now')}
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
            <Form.Field
              control={form.control}
              name="startAt"
              render={({ field }) => (
                <Form.Item className="py-0 mt-5">
                  <Form.Label required>
                    {t('form:feature-flags:start-at')}
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
          </Form>
        </FormProvider>
      </div>

      <ButtonBar
        secondaryButton={
          <Button
            loading={isSubmitting}
            disabled={!isDirty || isValid}
            onClick={onSubmit}
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
    </DialogModal>
  );
};

export default ConfirmationRequiredModal;
