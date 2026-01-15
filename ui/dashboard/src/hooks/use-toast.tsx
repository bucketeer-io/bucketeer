import { ReactNode } from 'react';
import toast, { Toast } from 'react-hot-toast';
import { Trans } from 'react-i18next';
import { AxiosError } from 'axios';
import { useTranslation } from 'i18n';
import ToastMessage from 'components/toast';

type ToastType = 'toast' | 'info-message' | 'prerequisite-message';
type MessageType = 'success' | 'info' | 'warning' | 'error';
type BackendErrorMetadata = {
  metadata: Record<string, string | undefined>;
};
export type ServerErrorType = AxiosError<{
  message?: string;
  details?: BackendErrorMetadata[];
}>;

export type NotifyProps = {
  toastType?: ToastType;
  messageType?: MessageType;
  message: ReactNode;
  t?: Toast;
  duration?: number;
  toastChildren?: ReactNode;
  className?: string;
};

export const useToast = () => {
  const { t } = useTranslation(['message', 'backend']);
  const notify = ({
    toastType = 'toast',
    messageType = 'success',
    message,
    duration = 5000,
    toastChildren,
    className
  }: NotifyProps) =>
    toast.custom(
      t => (
        <ToastMessage
          toastType={toastType}
          messageType={messageType}
          message={message}
          t={t}
          toastChildren={toastChildren}
          className={className}
        />
      ),
      {
        duration
      }
    );

  const errorNotify = (error?: unknown, message?: string) => {
    const response = (error as ServerErrorType | undefined)?.response;
    const details = response?.data?.details;
    if (!details?.length) {
      return notify({
        messageType: 'error',
        message: message || t('something-went-wrong')
      });
    }
    details.forEach(({ metadata }) => {
      const errorMessage = (
        <Trans
          i18nKey={
            metadata?.messageKey
              ? `backend:errors.${metadata.messageKey}`
              : 'message:something-went-wrong'
          }
          values={Object.fromEntries(
            Object.entries(metadata || {})
              .filter(([key]) => key !== 'messageKey')
              .map(([key, value]) => [
                key,
                value
                  ? t(`backend:nouns.${value}`, { defaultValue: value })
                  : value
              ])
          )}
        />
      );
      return notify({
        messageType: 'error',
        message: errorMessage
      });
    });
  };
  return {
    notify,
    errorNotify
  };
};
