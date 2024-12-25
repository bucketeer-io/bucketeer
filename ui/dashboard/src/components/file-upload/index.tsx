import { createRef, useEffect, useRef, useState } from 'react';
import type { InputHTMLAttributes, ChangeEvent } from 'react';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import type { ServerError } from '@types';
import { formatBytes } from 'utils/files';
import DraggingOverlay from './dragging-overlay';

export type FileUploadFormat = 'pdf' | 'image';

export interface FileUploadProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'accept'> {
  uploading: boolean;
  error?: ServerError;
  onUpload: (file: File) => void;
  onReset: () => void;
  format: FileUploadFormat | FileUploadFormat[];
  maxFileSize?: number;
  Uploading: () => JSX.Element;
  Error: ({
    error,
    onClearError
  }: {
    error: ServerError | ValidationError;
    onClearError: () => void;
  }) => JSX.Element;
  children: JSX.Element;
}

const validates: Record<FileUploadFormat, (file: File) => boolean> = {
  pdf: (file: File) => file.type === 'application/pdf',
  image: (file: File) => file.type.startsWith('image/')
};
const accepts: Record<FileUploadFormat, string | string[]> = {
  pdf: '.pdf',
  image: 'image/*'
};

const FileUpload = ({
  multiple,
  format,
  uploading,
  error: ServerError,
  onUpload,
  onReset,
  maxFileSize,
  children,
  Uploading: UploadingComponent,
  Error: ErrorComponent,
  ...props
}: FileUploadProps) => {
  const { t } = useTranslation('common');
  const id = uuid();
  const dropzone = createRef<HTMLLabelElement>();
  const draggingCount = useRef(0);
  const formats = Array.isArray(format) ? format : [format];
  const accept = formats.map(item => accepts[item]).join(',');
  const validate = (file: File) => formats.every(item => validates[item](file));
  const [dragging, setDragging] = useState(false);
  const [validationError, setValidationError] =
    useState<ValidationError | null>(null);

  const handleValidate = async (_files?: FileList | null) => {
    // Reset previous upload state
    onReset();
    setValidationError(null);

    // This is required to convert FileList object to array
    const files = [...(_files || [])];

    // Validate single file upload
    if (!multiple && files.length > 1) {
      setValidationError({
        title: t('file-upload.error-multiple-file-title'),
        message: t('file-upload.error-multiple-file-message')
      });
      return;
    }

    // Ok, single file
    const file = files[0];
    if (!file) return;

    // Validate file format
    if (validate(file)) {
      setValidationError({
        title: t('file-upload.error-format-title'),
        message: t('file-upload.error-format-message', { format })
      });
      return;
    }

    // Validate file size
    if (maxFileSize && file.size > maxFileSize) {
      setValidationError({
        title: t('file-upload.error-file-size-title'),
        message: t('file-upload.error-file-size-message', {
          size: formatBytes(maxFileSize)
        })
      });
      return;
    }

    // Ok, upload
    onUpload(file);
  };

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    handleValidate(e.target?.files);
  };

  const handleClearError = () => {
    onReset();
    setValidationError(null);
  };

  const handleDragIn = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (uploading) return;
    draggingCount.current += 1;
    if (e.dataTransfer?.items?.length) {
      setDragging(true);
    }
  };

  const handleDragOut = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (uploading) return;
    draggingCount.current -= 1;
    if (draggingCount.current > 0) return;
    setDragging(false);
  };

  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (uploading) return;
    setDragging(false);
    draggingCount.current = 0;
    handleValidate(e.dataTransfer?.files);
  };

  const handleDrag = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  useEffect(() => {
    if (dropzone.current) {
      dropzone.current.addEventListener('dragenter', handleDragIn);
      dropzone.current.addEventListener('dragleave', handleDragOut);
      dropzone.current.addEventListener('dragover', handleDrag);
      dropzone.current.addEventListener('drop', handleDrop);
    }

    return () => {
      if (dropzone.current) {
        dropzone.current.removeEventListener('dragenter', handleDragIn);
        dropzone.current.removeEventListener('dragleave', handleDragOut);
        dropzone.current.removeEventListener('dragover', handleDrag);
        dropzone.current.removeEventListener('drop', handleDrop);
      }
    };
  }, []);

  return (
    <>
      <label ref={dropzone} htmlFor={id} className="relative overflow-hidden">
        {dragging && <DraggingOverlay />}
        {(() => {
          if (uploading) {
            return <UploadingComponent />;
          }

          if (ServerError) {
            return (
              <ErrorComponent
                error={ServerError}
                onClearError={handleClearError}
              />
            );
          }

          if (validationError) {
            return (
              <ErrorComponent
                error={validationError}
                onClearError={handleClearError}
              />
            );
          }

          return children;
        })()}
      </label>

      <input
        {...props}
        id={id}
        className="hidden"
        type="file"
        accept={accept}
        multiple={multiple}
        onChange={handleChange}
      />
    </>
  );
};

export default FileUpload;
