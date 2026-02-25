import {
  ReactNode,
  useState,
  createContext,
  useContext,
  useEffect,
  Dispatch,
  SetStateAction
} from 'react';
import { useTranslation } from 'react-i18next';
import { UNSAFE_NavigationContext as NavigationContext } from 'react-router-dom';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

interface ConfirmOptions {
  title?: string;
  message?: string;
  titleLeave?: string;
  titleStay?: string;
  onConfirm: () => void;
  onCancel?: () => void;
}

interface ConfirmContextType {
  isShow: boolean;
  setIsShow: Dispatch<SetStateAction<boolean>>;
  confirm: (options: ConfirmOptions) => void;
  options: ConfirmOptions | null;
  handleCancel: () => void;
  handleConfirm: () => void;
}
interface Props {
  title?: string;
  titleStay?: string;
  titleLeave?: string;
  message?: string;
  isOpen: boolean;
  onClose?: () => void;
  onConfirm: () => void;
}

const ConfirmContext = createContext<ConfirmContextType | null>(null);

export function useUnsavedLeavePage({
  isShow,
  title = 'message:leave-page-unsaved-changes',
  content = 'message:leave-page-unsaved-changes-content',
  callBackCancel
}: {
  isShow: boolean;
  title?: string;
  content?: string;
  titleLeave?: string;
  titleStay?: string;
  callBackCancel?: () => void;
}) {
  const { confirm } = useConfirm();
  const { t } = useTranslation(['form', 'common', 'table', 'message']);
  const navigator = useContext(NavigationContext).navigator;

  useEffect(() => {
    if (!isShow) return;

    const push = navigator.push;
    const replace = navigator.replace;

    navigator.push = (...args: Parameters<typeof push>) => {
      confirm({
        title: title,
        message: content,
        onConfirm: () => {
          if (callBackCancel) {
            callBackCancel();
          }
          return push(...args);
        }
      });
    };

    navigator.replace = (...args: Parameters<typeof replace>) => {
      confirm({
        title: title,
        message: content,
        onConfirm: () => {
          if (callBackCancel) {
            callBackCancel();
          }
          return replace(...args);
        }
      });
    };

    return () => {
      navigator.push = push;
      navigator.replace = replace;
    };
  }, [isShow, title, content, navigator]);

  useEffect(() => {
    if (!isShow) return;
    history.pushState(null, '', window.location.href);

    const handlePopState = () => {
      confirm({
        title: t(title),
        message: t(content),
        onConfirm: () => {
          if (callBackCancel) {
            callBackCancel();
          }
          history.back();
        },
        onCancel: () => {
          history.pushState(null, '', window.location.href);
        }
      });
    };

    window.addEventListener('popstate', handlePopState);

    return () => {
      window.removeEventListener('popstate', handlePopState);
    };
  }, [isShow, title, content]);

  useEffect(() => {
    if (!isShow) return;
    const handler = (e: BeforeUnloadEvent) => {
      e.preventDefault();
      e.returnValue = '';
    };
    window.addEventListener('beforeunload', handler);
    return () => window.removeEventListener('beforeunload', handler);
  }, [isShow]);
  return { isShow };
}

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const [options, setOptions] = useState<ConfirmOptions | null>(null);
  const [isShow, setIsShow] = useState<boolean>(false);
  const confirm = (opts: ConfirmOptions) => setOptions(opts);

  const handleConfirm = () => {
    options?.onConfirm();
    setOptions(null);
  };

  const handleCancel = () => {
    options?.onCancel?.();
    setOptions(null);
  };

  return (
    <ConfirmContext.Provider
      value={{
        confirm,
        options,
        isShow,
        setIsShow,
        handleCancel,
        handleConfirm
      }}
    >
      {children}
      {options && (
        <PopupGlobal
          title={options.title}
          titleLeave={options.titleLeave}
          titleStay={options.titleStay}
          message={options.message}
          isOpen={true}
          onClose={handleCancel}
          onConfirm={handleConfirm}
        />
      )}
    </ConfirmContext.Provider>
  );
}

export function useConfirm() {
  const context = useContext(ConfirmContext);
  const { t } = useTranslation(['message']);

  if (!context) {
    throw new Error(t('auth-context-error'));
  }
  return context;
}

export function PopupGlobal({
  titleLeave,
  titleStay,
  title = 'message:leave-page-unsaved-changes',
  message = 'message:leave-page-unsaved-changes-content',
  isOpen,
  onClose,
  onConfirm
}: Props) {
  const { t } = useTranslation(['message', 'form']);
  return (
    <DialogModal
      className="max-w-[500px]"
      title={t(title)}
      isOpen={isOpen}
      onClose={() => onClose?.()}
    >
      <div className="p-5">{t(message)}</div>

      <ButtonBar
        primaryButton={
          <Button
            type="button"
            variant="secondary"
            className="p-2 h-9 font-bold text-sm rounded-md"
            onClick={onClose}
          >
            {titleStay ? t(titleStay) : t(`common:continue-editing`)}
          </Button>
        }
        secondaryButton={
          <Button
            type="button"
            variant="negative"
            className="p-2 h-9 font-bold text-sm rounded-md"
            onClick={onConfirm}
          >
            {titleLeave ? t(titleLeave) : t(`common:leave-page`)}
          </Button>
        }
      />
    </DialogModal>
  );
}
