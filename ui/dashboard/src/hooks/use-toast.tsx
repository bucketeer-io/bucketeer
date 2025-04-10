import { ReactNode } from 'react';
import toast, { Toast } from 'react-hot-toast';
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
};

export const useToast = () => {
  const notify = ({
    toastType = 'toast',
    messageType = 'success',
    message,
    duration = 5000,
    toastChildren
  }: NotifyProps) =>
    toast.custom(
      t => (
        <ToastMessage
          toastType={toastType}
          messageType={messageType}
          message={message}
          t={t}
          toastChildren={toastChildren}
        />
      ),
      {
        duration
      }
    );

  const errorNotify = (error?: unknown, message?: string) =>
    notify({
      messageType: 'error',
      message: message || (error as Error)?.message || 'Something went wrong.'
    });

  return {
    notify,
    errorNotify
  };
};
