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
  registerProceed: (fn: (() => void) | null) => void;
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
  const {
    confirm,
    setIsShow: setIsShowGlobal,
    registerProceed,
    bypassRef
  } = useConfirm();

  const blocker = useBlocker(
    useCallback<BlockerFunction>(
      ({ currentLocation, nextLocation }) => {
        if (bypassRef.current) {
          bypassRef.current = false;
          return false;
        }
        return isShow && currentLocation.pathname !== nextLocation.pathname;
      },
      [isShow]
    )
  );

  useEffect(() => {
    setIsShowGlobal(isShow);
  }, [isShow]);

  // When blocker fires, show the confirmation dialog
  useEffect(() => {
    if (blocker.state !== 'blocked') return;

    confirm({
      title,
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

  // Register blocker.proceed only when the blocker is actually blocked
  useEffect(() => {
    registerProceed(
      blocker.state === 'blocked' ? () => blocker.proceed() : null
    );
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
  const [options, setOptions] = useState<ConfirmOptions | null>(null);
  const [isShow, setIsShow] = useState<boolean>(false);
  const proceedRef = useRef<(() => void) | null>(null);
  const bypassRef = useRef(false);

  const confirm = (opts: ConfirmOptions) => setOptions(opts);

  const registerProceed = useCallback((fn: (() => void) | null) => {
    proceedRef.current = fn;
  }, []);

  const allowNavigation = useCallback((action: () => void) => {
    if (proceedRef.current) {
      // Blocker is active — proceed through it, then run the action
      proceedRef.current();
      action();
    } else {
      // Blocker is idle — set a one-shot bypass so the next navigation isn't blocked.
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
