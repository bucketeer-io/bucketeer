import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useEffect,
  useState
} from 'react';
import { useTranslation } from 'react-i18next';
import { UNSAFE_NavigationContext as NavigationContext } from 'react-router';
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
  setOptions: Dispatch<SetStateAction<ConfirmOptions | null>>;
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

let bypassNavigation = false;

export function allowNavigation(action?: () => void) {
  bypassNavigation = true;
  if (action) {
    try {
      action();
    } finally {
      bypassNavigation = false;
    }
  } else {
    // Fallback for existing call sites that do not pass a callback:
    // ensure the bypass is short-lived and does not leak indefinitely.
    setTimeout(() => {
      bypassNavigation = false;
    }, 0);
  }
}

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
  const { confirm, setIsShow: setIsShowGlobal, isShow: global } = useConfirm();
  const navigator = useContext(NavigationContext).navigator;

  useEffect(() => {
    setIsShowGlobal(isShow);
  }, [isShow]);

  useEffect(() => {
    if (!global) return;

    const push = navigator.push;
    const replace = navigator.replace;

    navigator.push = (...args: Parameters<typeof push>) => {
      if (bypassNavigation) {
        bypassNavigation = false;
        return push(...args);
      }
      confirm({
        title: title,
        message: content,
        onConfirm: () => {
          if (callBackCancel) {
            callBackCancel();
          }
          setIsShowGlobal(false);
          return push(...args);
        }
      });
    };

    navigator.replace = (...args: Parameters<typeof replace>) => {
      if (bypassNavigation) {
        bypassNavigation = false;
        return replace(...args);
      }
      confirm({
        title: title,
        message: content,
        onConfirm: () => {
          if (callBackCancel) {
            callBackCancel();
          }
          setIsShowGlobal(false);
          return replace(...args);
        }
      });
    };

    return () => {
      navigator.push = push;
      navigator.replace = replace;
    };
  }, [global, title, content, navigator]);

  useEffect(() => {
    if (!global) return;
    history.pushState(null, '', window.location.href);

    const handlePopState = () => {
      if (bypassNavigation) {
        bypassNavigation = false;
        return;
      }
      confirm({
        title: title,
        message: content,
        onConfirm: () => {
          if (callBackCancel) {
            callBackCancel();
          }
          setIsShowGlobal(false);
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
  }, [global, title, content]);

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
        setOptions,
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
      className="w-[500px]"
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
