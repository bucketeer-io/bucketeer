import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import DialogModal from 'components/modal/dialog';
import TextArea from 'components/textarea';

export type SaveWithCommentModalProps = {
  isOpen: boolean;
  isRequired: boolean;
  onSubmit: () => void;
  onClose: () => void;
};

const SaveWithCommentModal = ({
  isOpen,
  isRequired,
  onClose,
  onSubmit
}: SaveWithCommentModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const {
    control,
    formState: { isValid, isSubmitting },
    handleSubmit
  } = useFormContext();

  return (
    <DialogModal
      className="w-[500px]"
      title={t('update-flag')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start px-5 py-8">
        <Form.Field
          name="comment"
          control={control}
          render={({ field }) => (
            <Form.Item className="py-0 w-full">
              <Form.Label required={isRequired} optional={!isRequired}>
                {t('form:comment-for-update')}
              </Form.Label>
              <Form.Control>
                <TextArea
                  {...field}
                  placeholder={t('form:placeholder-comment')}
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
            disabled={!isValid}
            loading={isSubmitting}
            onClick={handleSubmit(onSubmit)}
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

export default SaveWithCommentModal;
