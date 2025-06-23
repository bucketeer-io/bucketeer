import { ReactNode } from 'react';
import toast, { Toast } from 'react-hot-toast';
import { AxiosError } from 'axios';
import { useTranslation } from 'i18n';
import ToastMessage from 'components/toast';

type ToastType = 'toast' | 'info-message' | 'prerequisite-message';
type MessageType = 'success' | 'info' | 'warning' | 'error';

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
  const { t } = useTranslation(['message']);
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
    const { message: _message, status } = (error as AxiosError) || {};
    return notify({
      messageType: 'error',
      message:
        status === 409
          ? t('same-data-exists')
          : message || _message || t('something-went-wrong')
    });
  };
  return {
    notify,
    errorNotify
  };
};
