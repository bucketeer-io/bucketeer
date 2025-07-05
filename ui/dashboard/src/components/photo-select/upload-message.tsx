import { Trans, useTranslation } from 'react-i18next';
import { formatBytes } from 'utils/files';
import UploadZone from './upload-zone';

const UploadMessage = ({ maxFileSize }: { maxFileSize: number }) => {
  const { t } = useTranslation(['form']);
  return (
    <UploadZone color="positive">
      <div className="flex-center flex-col gap-y-2 typo-para-small text-gray-600 mt-3">
        <div>
          <Trans
            i18nKey="form:upload-files"
            values={{ text: t('browse') }}
            components={{
              underline: <u className="text-primary-500 inline" />
            }}
          />
        </div>
        <Trans
          i18nKey="form:accept-file-types"
          values={{
            type: `JPG, JPGE, PNG (${t('max')}. ${formatBytes(maxFileSize)})`
          }}
        />
      </div>
    </UploadZone>
  );
};

export default UploadMessage;
