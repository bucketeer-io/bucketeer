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
import { type BlockerFunction, useBlocker } from 'react-router';
import { LEAVE_PAGE_CANCELLED_EVENT } from 'constants/walkthrough';
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
  registerProceed: (id: number, fn: (() => void) | null) => void;
  allowNavigation: (action: () => void) => void;
  bypassRef: RefObject<boolean>;
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

// While the onboarding walkthrough (driver.js) runs, navigation attempts are
// ignored instead of prompting, so its guided steps are never interrupted.
const isWalkthroughActive = () =>
  document.body.classList.contains('driver-active');

let nextInstanceId = 0;

export function useUnsavedLeavePage({
  isShow,
  title = 'message:leave-page-unsaved-changes',
  content = 'message:leave-page-unsaved-changes-content',
  titleLeave,
  titleStay,
  callBackCancel
}: {
  isShow: boolean;
  title?: string;
  content?: string;
  titleLeave?: string;
  titleStay?: string;
  callBackCancel?: () => void;
}) {
  const {
    confirm,
    setIsShow: setIsShowGlobal,
    registerProceed,
    bypassRef
  } = useConfirm();

  const instanceId = useRef(++nextInstanceId).current;

  const blocker = useBlocker(
    useCallback<BlockerFunction>(() => {
      if (bypassRef.current) {
        bypassRef.current = false;
        return false;
      }
      if (isWalkthroughActive()) return false;
      return isShow;
    }, [isShow])
  );

  useEffect(() => {
    setIsShowGlobal(isShow);
  }, [isShow]);

  // When blocker fires, show the confirmation dialog
  useEffect(() => {
    if (blocker.state !== 'blocked') return;

    confirm({
      title,
      titleLeave,
      titleStay,
      message: content,
      onConfirm: () => {
        callBackCancel?.();
        setIsShowGlobal(false);
        blocker.proceed();
      },
      onCancel: () => {
        blocker.reset();
      }
    });
  }, [blocker.state]);

  // Register blocker.proceed only when the blocker is actually blocked,
  // and unregister on unmount so a stale proceed is never invoked later.
  useEffect(() => {
    registerProceed(
      instanceId,
      blocker.state === 'blocked' ? () => blocker.proceed() : null
    );
    return () => registerProceed(instanceId, null);
  }, [blocker.state]);

  // Browser tab close / reload guard
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
  const [queue, setQueue] = useState<ConfirmOptions[]>([]);
  const [isShow, setIsShow] = useState<boolean>(false);
  const proceedFnsRef = useRef<Map<number, () => void>>(new Map());
  const bypassRef = useRef(false);

  const options = queue[0] ?? null;

  const confirm = useCallback((opts: ConfirmOptions) => {
    setQueue(prev => [...prev, opts]);
  }, []);

  const registerProceed = useCallback((id: number, fn: (() => void) | null) => {
    if (fn) {
      proceedFnsRef.current.set(id, fn);
    } else {
      proceedFnsRef.current.delete(id);
    }
  }, []);

  const allowNavigation = useCallback((action: () => void) => {
    if (proceedFnsRef.current.size > 0) {
      // One or more blockers are active — proceed through all of them, then run the action
      proceedFnsRef.current.forEach(proceed => proceed());
      proceedFnsRef.current.clear();
      action();
    } else {
      // No blocker is active — set a one-shot bypass so the next navigation isn't blocked.
      bypassRef.current = true;
      try {
        action();
      } finally {
        queueMicrotask(() => {
          bypassRef.current = false;
        });
      }
    }
  }, []);

  const handleConfirm = () => {
    options?.onConfirm();
    setQueue(prev => prev.slice(1));
  };

  const handleCancel = () => {
    options?.onCancel?.();
    setQueue(prev => prev.slice(1));
    document.dispatchEvent(new CustomEvent(LEAVE_PAGE_CANCELLED_EVENT));
  };

  return (
    <ConfirmContext.Provider
      value={{
        confirm,
        options,
        isShow,
        setIsShow,
        handleCancel,
        handleConfirm,
        registerProceed,
        allowNavigation,
        bypassRef
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
