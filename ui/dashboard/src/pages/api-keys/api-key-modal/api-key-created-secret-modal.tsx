import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { copyToClipBoard } from 'utils/function';
import { IconCopy } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

type APIKeyCreatedSecretModalProps = {
  apiKeySecret: string | null;
  onClose: () => void;
};

const APIKeyCreatedSecretModal = ({
  apiKeySecret,
  onClose
}: APIKeyCreatedSecretModalProps) => {
  const { t } = useTranslation(['message']);
  const { notify } = useToast();

  if (!apiKeySecret) {
    return null;
  }

  const handleCopy = () => {
    copyToClipBoard(apiKeySecret);
    notify({ message: t('message:copied') });
  };

  return (
    <DialogModal
      title={t('message:api-key-created-title')}
      isOpen={!!apiKeySecret}
      onClose={onClose}
      closeOnClickOutside={false}
      className="w-full max-w-lg"
    >
      <div className="flex flex-col gap-4 p-5 pb-6">
        <p
          className="typo-para-medium rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-amber-900"
          role="alert"
        >
          {t('message:api-key-secret-warning')}
        </p>
        <div className="rounded-lg border border-gray-200 bg-gray-100 p-3">
          <p className="typo-para-small mb-2 text-gray-600">
            {t('message:api-key-secret-label')}
          </p>
          <code className="typo-para-small block break-all font-mono text-gray-900">
            {apiKeySecret}
          </code>
        </div>
        <div className="flex flex-wrap items-center gap-3">
          <Button
            type="button"
            variant="secondary"
            onClick={handleCopy}
            className="gap-2"
          >
            <Icon icon={IconCopy} size="sm" />
            {t('message:copy-api-key')}
          </Button>
          <Button type="button" variant="primary" onClick={onClose}>
            {t('message:api-key-secret-close')}
          </Button>
        </div>
      </div>
    </DialogModal>
  );
};

export default APIKeyCreatedSecretModal;
