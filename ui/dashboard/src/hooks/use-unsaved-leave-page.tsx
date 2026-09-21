import {
  createContext,
  ReactNode,
  RefObject,
  useCallback,
  useContext,
  useEffect,
  useMemo,
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

interface RegisteredBlocker {
  isShow: boolean;
  title?: string;
  content?: string;
  titleLeave?: string;
  titleStay?: string;
  callBackCancel?: () => void;
}

interface ConfirmContextType {
  isShow: boolean;
  confirm: (options: ConfirmOptions) => void;
  options: ConfirmOptions | null;
  handleCancel: () => void;
  handleConfirm: () => void;
  registerBlocker: (id: number, state: RegisteredBlocker | null) => void;
  allowNavigation: (action: () => void) => void;
  resetAllBlockers: () => void;
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
  const { registerBlocker } = useConfirm();

  const instanceId = useRef(++nextInstanceId).current;

  // Register this instance's dirty state with the provider, which owns the
  // single router blocker and aggregates every registered instance into one
  // unsaved-state decision. Multiple instances can be mounted at once (e.g. a
  // dirty parent page with a dirty nested modal), so each is tracked
  // separately instead of each instance creating its own blocker.
  useEffect(() => {
    registerBlocker(instanceId, {
      isShow,
      title,
      content,
      titleLeave,
      titleStay,
      callBackCancel
    });
  }, [isShow, title, content, titleLeave, titleStay, callBackCancel]);

  useEffect(() => () => registerBlocker(instanceId, null), []);

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
  const [isShow, setIsShow] = useState(false);
  const bypassRef = useRef(false);
  const blockersRef = useRef<Map<number, RegisteredBlocker>>(new Map());

  const options = queue[0] ?? null;

  const confirm = useCallback((opts: ConfirmOptions) => {
    setQueue(prev => [...prev, opts]);
  }, []);

  const isAnyDirty = useCallback(
    () => [...blockersRef.current.values()].some(b => b.isShow),
    []
  );

  const registerBlocker = useCallback(
    (id: number, state: RegisteredBlocker | null) => {
      if (state) {
        blockersRef.current.set(id, state);
      } else {
        blockersRef.current.delete(id);
      }
      setIsShow(isAnyDirty());
    },
    [isAnyDirty]
  );

  // Leaving is a decision to discard every unsaved change, not just the one
  // that triggered the prompt. Clearing the whole registry here (rather than
  // relying on each instance to unmount and unregister itself) prevents a
  // still-mounted dirty instance from re-blocking a chained navigation that
  // fires before React has torn down the old route tree — e.g. switching
  // organizations issues a root navigation, which then triggers a follow-up
  // redirect to the env-specific URL from an effect in the newly mounted
  // route; without this, that second navigation re-opens the confirm dialog.
  const resetAllBlockers = useCallback(() => {
    blockersRef.current.clear();
    setIsShow(false);
  }, []);

  const blocker = useBlocker(
    useCallback<BlockerFunction>(() => {
      if (bypassRef.current) {
        bypassRef.current = false;
        return false;
      }
      return isAnyDirty();
    }, [])
  );

  // When blocker fires, show the confirmation dialog — unless the onboarding
  // walkthrough (driver.js) is running, in which case the navigation attempt
  // is silently suppressed (reset without prompting) so guided steps are
  // never interrupted.
  useEffect(() => {
    if (blocker.state !== 'blocked') return;

    if (isWalkthroughActive()) {
      blocker.reset();
      return;
    }

    // Use the most recently registered dirty instance to source the dialog's
    // copy, but proceeding must clear every dirty instance so none of them
    // re-blocks the very navigation the user just confirmed.
    const dirtyBlockers = [...blockersRef.current.values()].filter(
      b => b.isShow
    );
    const active = dirtyBlockers[dirtyBlockers.length - 1];

    confirm({
      title: active?.title,
      titleLeave: active?.titleLeave,
      titleStay: active?.titleStay,
      message: active?.content,
      onConfirm: () => {
        dirtyBlockers.forEach(b => b.callBackCancel?.());
        resetAllBlockers();
        blocker.proceed();
      },
      onCancel: () => {
        blocker.reset();
      }
    });
  }, [blocker.state]);

  const allowNavigation = useCallback((action: () => void) => {
    // One-shot bypass so the next blocked navigation isn't blocked.
    bypassRef.current = true;
    try {
      action();
    } finally {
      queueMicrotask(() => {
        bypassRef.current = false;
      });
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

  const value = useMemo<ConfirmContextType>(
    () => ({
      isShow,
      confirm,
      options,
      handleCancel,
      handleConfirm,
      registerBlocker,
      allowNavigation,
      resetAllBlockers,
      bypassRef
    }),
    [
      isShow,
      confirm,
      options,
      allowNavigation,
      registerBlocker,
      resetAllBlockers
    ]
  );

  return (
    <ConfirmContext.Provider value={value}>
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
