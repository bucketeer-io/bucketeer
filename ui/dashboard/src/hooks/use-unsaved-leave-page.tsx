import {
  ReactNode,
  useState,
  createContext,
  useContext,
  useEffect
} from 'react';
import { useTranslation } from 'react-i18next';
import { UNSAFE_NavigationContext as NavigationContext } from 'react-router-dom';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export function useNavigationPage({
  isShow,
  title,
  content
}: {
  isShow: boolean;
  title: string;
  content: string;
}) {
  const { confirm } = useConfirm();
  const navigator = useContext(NavigationContext).navigator;

  useEffect(() => {
    if (!isShow) return;

    const push = navigator.push;
    const replace = navigator.replace;

    navigator.push = (...args: Parameters<typeof push>) => {
      confirm({
        title,
        message: content,
        onConfirm: () => push(...args)
      });
    };

    navigator.replace = (...args: Parameters<typeof replace>) => {
      confirm({
        title,
        message: content,
        onConfirm: () => replace(...args)
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
        title,
        message: content,
        onConfirm: () => {
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
}

type ConfirmOptions = {
  title: string;
  message: string;
  onConfirm: () => void;
  onCancel?: () => void;
};

type ConfirmContextType = {
  confirm: (options: ConfirmOptions) => void;
};

const ConfirmContext = createContext<ConfirmContextType | null>(null);

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const { t } = useTranslation(['message', 'form']);
  const [options, setOptions] = useState<ConfirmOptions | null>(null);
  const confirm = (opts: ConfirmOptions) => setOptions(opts);
  const handleConfirm = () => {
    options?.onConfirm();
    setOptions(null);
  };
  const handleCancel = () => {
    options?.onCancel?.();
    setOptions(null);
  };
  console.log('check optionn :::: ', options);
  return (
    <ConfirmContext.Provider value={{ confirm }}>
      {children}
      {options && (
        <DialogModal
          className="w-[500px]"
          title={options.title}
          isOpen={!!options}
          onClose={handleCancel}
        >
          <div className="p-5">{options.message}</div>

          <ButtonBar
            primaryButton={
              <Button
                type="button"
                variant="secondary"
                className="p-2 h-9 font-bold text-sm rounded-md"
                onClick={handleCancel}
              >
                {t(`common:stay-on`)}
              </Button>
            }
            secondaryButton={
              <Button
                type="button"
                variant="negative"
                className="p-2 h-9 font-bold text-sm rounded-md"
                onClick={handleConfirm}
              >
                {t(`common:leave-page`)}
              </Button>
            }
          />
        </DialogModal>
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
