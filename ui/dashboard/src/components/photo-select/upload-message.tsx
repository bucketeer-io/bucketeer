import { formatBytes } from 'utils/files';
import UploadZone from './upload-zone';

const UploadMessage = ({ maxFileSize }: { maxFileSize: number }) => {
  return (
    <UploadZone color="positive">
      <div className="typo-para-small text-gray-600 mt-3">
        {`Drag and drop your picture here`} or
        <span className="underline text-primary-600 ml-1">{`browse`}</span>
      </div>

      <div className="typo-para-tiny text-gray-500 mt-2">
        {`JPG, JPGE, PNG (max. ${formatBytes(maxFileSize)})`}
      </div>
    </UploadZone>
  );
};

export default UploadMessage;
