import {
  createContext,
  Dispatch,
  ReactNode,
  RefObject,
  SetStateAction,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState
} from 'react';
import { useTranslation } from 'react-i18next';
import { useBlocker } from 'react-router';
import { LEAVE_PAGE_CANCELLED_EVENT } from 'constants/walkthrough';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import { createUnsavedForms, UnsavedForm } from './unsaved-forms';

interface ConfirmOptions {
  title?: string;
  message?: string;
  titleLeave?: string;
  titleStay?: string;
  onConfirm: () => void;
  onCancel?: () => void;
}

interface ConfirmContextType {
  register: (form: RefObject<UnsavedForm>) => () => void;
  syncDirtyState: () => void;
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

// While the onboarding walkthrough (driver.js) runs, navigation attempts are
// ignored instead of prompting, so its guided steps are never interrupted.
const isWalkthroughActive = () =>
  document.body.classList.contains('driver-active');

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
  const { register, syncDirtyState } = useConfirm();
  const form = useRef({ isShow, title, content, callBackCancel });
  form.current = { isShow, title, content, callBackCancel };

  useEffect(() => register(form), [register]);
  useEffect(() => {
    syncDirtyState();
  }, [isShow, syncDirtyState]);
  return { isShow };
}

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const [options, setOptions] = useState<ConfirmOptions | null>(null);
  const [isShow, setIsShow] = useState<boolean>(false);
  const [forms] = useState(createUnsavedForms);
  const getDirtyForms = forms.getDirtyForms;
  const syncDirtyState = useCallback(() => {
    setIsShow(getDirtyForms().length > 0);
  }, [getDirtyForms]);
  const register = useCallback(
    (form: RefObject<UnsavedForm>) => {
      const unregister = forms.register(form);
      syncDirtyState();
      return () => {
        unregister();
        syncDirtyState();
      };
    },
    [forms, syncDirtyState]
  );
  const confirm = useCallback((opts: ConfirmOptions) => setOptions(opts), []);
  const shouldBlock = useCallback(() => {
    if (bypassNavigation) {
      bypassNavigation = false;
      return false;
    }
    return getDirtyForms().length > 0;
  }, [getDirtyForms]);
  const blocker = useBlocker(shouldBlock);

  useEffect(() => {
    if (blocker.state !== 'blocked') return;
    if (isWalkthroughActive()) {
      blocker.reset();
      return;
    }
    const dirtyForms = getDirtyForms();
    if (!dirtyForms.length) {
      blocker.proceed();
      return;
    }
    // One confirmation covers all dirty forms. The first supplies the copy.
    confirm({
      title: dirtyForms[0].title,
      message: dirtyForms[0].content,
      onConfirm: () => {
        getDirtyForms().forEach(form => form.callBackCancel?.());
        setIsShow(false);
        blocker.proceed();
      },
      onCancel: () => blocker.reset()
    });
  }, [blocker, confirm, getDirtyForms]);

  useEffect(() => {
    const handler = (event: BeforeUnloadEvent) => {
      if (!getDirtyForms().length) return;
      event.preventDefault();
      event.returnValue = '';
    };
    window.addEventListener('beforeunload', handler);
    return () => window.removeEventListener('beforeunload', handler);
  }, [getDirtyForms]);

  const handleConfirm = () => {
    options?.onConfirm();
    setOptions(null);
  };

  const handleCancel = () => {
    options?.onCancel?.();
    setOptions(null);
    document.dispatchEvent(new CustomEvent(LEAVE_PAGE_CANCELLED_EVENT));
  };

  return (
    <ConfirmContext.Provider
      value={{
        register,
        syncDirtyState,
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
      <div className="p-5 dark:text-dark-gray-300">{t(message)}</div>

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
