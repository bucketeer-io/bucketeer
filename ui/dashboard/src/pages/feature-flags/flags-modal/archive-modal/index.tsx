import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { cn } from 'utils/style';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import DialogModal from 'components/modal/dialog';
import TextArea from 'components/textarea';
import ArchiveWarning from './archive-warning';

export type ArchiveModalProps = {
  isArchiving: boolean;
  isShowWarning: boolean;
  isOpen: boolean;
  title: string;
  description: React.ReactElement | string;
  className?: string;
  isLoading?: boolean;
  onClose: () => void;
  onSubmit: ({ comment }: { comment?: string }) => Promise<void>;
};

const ArchiveModal = ({
  isArchiving,
  isShowWarning,
  isOpen,
  title,
  description,
  className,
  isLoading,
  onClose,
  onSubmit
}: ArchiveModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  const formSchema = yup.object().shape({
    comment: currentEnvironment?.requireComment
      ? yup.string().required()
      : yup.string()
  });

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      comment: ''
    }
  });

  const {
    formState: { isValid, isSubmitting }
  } = form;

  return (
    <DialogModal
      className="w-[500px]"
      title={title}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div
        className={cn(
          'flex flex-col w-full gap-y-5 items-start p-5',
          className
        )}
      >
        <div className="typo-para-small text-gray-600 w-full">
          {description}
        </div>
        {isShowWarning && <ArchiveWarning />}
        <FormProvider {...form}>
          <Form className="w-full" onSubmit={form.handleSubmit(onSubmit)}>
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
                      disabled={!editable}
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
            loading={isSubmitting || isLoading}
            disabled={!isValid || !editable}
            onClick={form.handleSubmit(onSubmit)}
          >
            {t(isArchiving ? `archive-flag` : 'unarchive-flag')}
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

export default ArchiveModal;
