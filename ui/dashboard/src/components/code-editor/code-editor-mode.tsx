import { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { EditorProps } from '@monaco-editor/react';
import { IconAlert } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import ReactCodeEditor from '.';

interface ReactCodeEditorModalProps extends EditorProps {
  isOpen: boolean;
  onClose: () => void;
  defaultLanguage?: string;
  error?: string;
  title: string | ReactNode;
  readOnly?: boolean;
}

export default function ReactCodeEditorModal({
  isOpen,
  onClose,
  title,
  value,
  defaultLanguage,
  onChange,
  error,
  readOnly = false,
  ...props
}: ReactCodeEditorModalProps) {
  const { t } = useTranslation('common');
  return (
    <DialogModal
      className="w-[60vw] h-[80%] overflow-hidden"
      overlayCls="bg-overlay-second"
      title={title}
      closeContent={
        <Button
          type="button"
          onClick={onClose}
          variant="secondary"
          className="!py-1 !px-2 h-7"
        >
          {t('common:close')}
        </Button>
      }
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="w-full h-full relative">
        <ReactCodeEditor
          defaultLanguage={defaultLanguage}
          value={value}
          onChange={onChange}
          isDefaulScroll
          isExpand={false}
          isResize={false}
          wrapCls="rounded-none"
          className="w-full h-full max-h-full"
          lastLine={10}
          readOnly={readOnly}
          {...props}
        />
        {error && (
          <div className="sticky left-0 bottom-0 pl-3 flex items-center gap-[10px] bg-white w-full h-9 typo-para-small text-accent-red-500">
            <Icon icon={IconAlert} size="xs" /> {error}
          </div>
        )}
      </div>
    </DialogModal>
  );
}
