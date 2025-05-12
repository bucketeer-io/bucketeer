import { useState } from 'react';
import toast from 'react-hot-toast';
import { cva } from 'class-variance-authority';
import { NotifyProps } from 'hooks/use-toast';
import { cn } from 'utils/style';
import {
  IconChevronRight,
  IconClose,
  IconToastError,
  IconToastInfo,
  IconToastSuccess,
  IconToastWarning
} from '@icons';
import Icon from 'components/icon';

const otherToastCls = 'px-4 py-2 border-l-4 rounded';

const toastVariants = cva(
  'flex flex-col w-fit max-w-[90%] typo-para-medium select-none',
  {
    variants: {
      toastType: {
        toast: 'sm:max-w-[460px] p-4 rounded-xl border text-gray-700',
        'info-message': cn(otherToastCls, 'sm:max-w-[490px] gap-y-2'),
        'prerequisite-message': cn(otherToastCls, 'sm:max-w-[79.5%] gap-y-3')
      },
      messageType: {
        success: 'bg-accent-green-50 border-accent-green-500',
        error: 'bg-accent-red-50 border-accent-red-500',
        info: 'bg-accent-blue-50 border-accent-blue-500 text-accent-blue-500',
        warning:
          'bg-accent-yellow-50 border-accent-yellow-500 text-accent-yellow-500'
      }
    }
  }
);

const ToastMessage = ({
  toastType = 'toast',
  messageType = 'success',
  message,
  t,
  className,
  toastChildren
}: NotifyProps) => {
  const [isExpand, setIsExpand] = useState(false);
  const icon =
    messageType === 'success'
      ? IconToastSuccess
      : messageType === 'error'
        ? IconToastError
        : messageType === 'info'
          ? IconToastInfo
          : IconToastWarning;
  return (
    <div className={cn(toastVariants({ toastType, messageType }), className)}>
      <div
        className={cn('flex gap-x-2', {
          'items-center': ['info-message', 'prerequisite-message'].includes(
            toastType
          )
        })}
      >
        <div className="flex-center size-6">
          <Icon icon={icon} size={'fit'} />
        </div>
        <div>{message}</div>
        {['success', 'error'].includes(messageType) && (
          <div
            className="flex-center size-5 hover:cursor-pointer"
            onClick={() => toast.dismiss(t?.id)}
          >
            <Icon icon={IconClose} size={'fit'} />
          </div>
        )}
      </div>
      {toastType === 'info-message' && <>{toastChildren}</>}
      {toastType === 'prerequisite-message' && (
        <div className="flex flex-col w-full pl-6">
          <div
            className="flex items-center gap-x-2 text-gray-700 cursor-pointer"
            onClick={() => setIsExpand(!isExpand)}
          >
            <p className="typo-para-medium">See more</p>
            <div
              className={cn(
                'flex-center size-fit rotate-90 transition-transform duration-200',
                {
                  '-rotate-90': isExpand
                }
              )}
            >
              <Icon icon={IconChevronRight} size={'sm'} />
            </div>
          </div>
          <div
            className={cn('h-0 pt-0 opacity-0 transition-all duration-200', {
              'h-fit pt-3 opacity-100': isExpand
            })}
          >
            {toastChildren}
          </div>
        </div>
      )}
    </div>
  );
};

export default ToastMessage;
