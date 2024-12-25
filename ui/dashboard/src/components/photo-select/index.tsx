import { createRef, forwardRef, useEffect, useRef, useState } from 'react';
import type { InputHTMLAttributes, ChangeEvent } from 'react';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { formatBytes } from 'utils/files';
import DraggingOverlay from './dragging-overlay';
import ErrorMessage from './error-message';
import UploadMessage from './upload-message';

export type PhotoSelectRef = HTMLInputElement;

export interface PhotoSelectProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, 'accept' | 'onChange'> {
  onChange: (file: File) => void;
  format?: PhotoSelectFormat | PhotoSelectFormat[];
  maxFileSize: number;
}

export interface ValidationError {
  title: string;
  message: string;
}

export type PhotoSelectFormat = 'image';

const validates: Record<PhotoSelectFormat, (file: File) => boolean> = {
  image: (file: File) => file.type.startsWith('image/')
};
const accepts: Record<PhotoSelectFormat, string | string[]> = {
  image: 'image/jpeg, image/jpg, image/png'
};

const PhotoSelect = forwardRef<PhotoSelectRef, PhotoSelectProps>(
  ({ format = 'image', onChange, maxFileSize, ...props }, ref) => {
    const { t } = useTranslation('common');
    const id = uuid();
    const [inputKey, setInputKey] = useState(Date.now());
    const dropzone = createRef<HTMLLabelElement>();
    const draggingCount = useRef(0);
    const formats = Array.isArray(format) ? format : [format];
    const accept = formats.map(item => accepts[item]).join(',');
    const validate = (file: File) =>
      formats.every(item => validates[item](file));
    const [dragging, setDragging] = useState(false);
    const [validationError, setValidationError] =
      useState<ValidationError | null>(null);

    const handleValidate = async (_files?: FileList | null) => {
      // Clear validation error
      setValidationError(null);

      // This is required to convert FileList object to array
      const files = [...(_files || [])];

      // Validate single file upload
      if (files.length > 1) {
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
      if (!validate(file)) {
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
      onChange(file);
    };

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
      handleValidate(e.target?.files);
      setInputKey(Date.now());
    };

    const handleClearError = () => {
      setValidationError(null);
    };

    const handleDragIn = (e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      draggingCount.current += 1;
      if (e.dataTransfer?.items?.length) {
        setDragging(true);
      }
    };

    const handleDragOut = (e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      draggingCount.current -= 1;
      if (draggingCount.current > 0) return;
      setDragging(false);
    };

    const handleDrop = (e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
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

    let content;
    if (validationError) {
      content = (
        <ErrorMessage error={validationError} onClearError={handleClearError} />
      );
    } else {
      content = <UploadMessage maxFileSize={maxFileSize} />;
    }

    return (
      <>
        <label ref={dropzone} htmlFor={id} className="relative overflow-hidden">
          {dragging && <DraggingOverlay />}
          {content}
        </label>

        <input
          {...props}
          ref={ref}
          id={id}
          key={inputKey}
          className="hidden"
          type="file"
          accept={accept}
          onChange={handleChange}
        />
      </>
    );
  }
);

export default PhotoSelect;
