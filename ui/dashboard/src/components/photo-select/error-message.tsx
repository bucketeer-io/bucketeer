import { IconErrorOutlineRound } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import type { ServerError } from '@types';
import Button from 'components/button';
import { ValidationError } from 'components/photo-select';
import UploadZone from './upload-zone';

const ErrorMessage = ({
  error,
  onClearError
}: {
  error: ServerError | ValidationError;
  onClearError: () => void;
}) => {
  const { t } = useTranslation('common');

  return (
    <UploadZone color="negative" icon={IconErrorOutlineRound}>
      {(error as ValidationError).title && (
        <div className="typo-header-light-tiny mb-1 text-accent-red-500">
          {(error as ValidationError).title}
        </div>
      )}
      {error.message && (
        <div className="typo-body-medium mb-1 max-w-[214px]">
          {error.message}
        </div>
      )}
      <Button variant="text" onClick={onClearError}>
        {t('file-upload.try-again')}
      </Button>
    </UploadZone>
  );
};

export default ErrorMessage;
