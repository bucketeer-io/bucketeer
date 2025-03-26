import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import Form from 'components/form';
import DialogModal from 'components/modal/dialog';
import TextArea from 'components/textarea';

export type ConfirmationRequiredModalProps = {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: () => void;
};

export interface ConfirmRequiredForm {
  resetSamplingSeed?: boolean;
  comment?: string;
}

const ConfirmationRequiredModal = ({
  isOpen,
  onClose,
  onSubmit
}: ConfirmationRequiredModalProps) => {
  const { t } = useTranslation(['common', 'form', 'table']);

  const {
    control,
    formState: { isDirty, isValid, isSubmitting },
    watch
  } = useFormContext();

  const isRequireComment = watch('requireComment');

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

        <div className="flex flex-col w-full px-5 pb-5">
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
                <Form.Control>
                  <Checkbox
                    ref={field.ref}
                    checked={field.value}
                    onCheckedChange={checked => field.onChange(checked)}
                    title={t('form:reset-sampling')}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
        </div>
        <ButtonBar
          secondaryButton={
            <Button
              type="submit"
              loading={isSubmitting}
              disabled={(isRequireComment && !isDirty) || !isValid}
              onClick={onSubmit}
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
    </DialogModal>
  );
};

export default ConfirmationRequiredModal;
