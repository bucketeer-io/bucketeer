import {
  forwardRef,
  useCallback,
  useEffect,
  useId,
  useMemo,
  useRef,
  useState
} from 'react';
import type { InputHTMLAttributes, ChangeEvent } from 'react';
import { useTranslation } from 'i18n';
import { formatBytes } from 'utils/files';
import DraggingOverlay from './dragging-overlay';
import ErrorMessage from './error-message';
import UploadMessage from './upload-message';

export type PhotoSelectRef = HTMLInputElement;

export interface PhotoSelectProps extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  'accept' | 'onChange'
> {
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
    const id = useId();
    const [inputKey, setInputKey] = useState(Date.now());
    const dropzone = useRef<HTMLLabelElement>(null);
    const draggingCount = useRef(0);
    const formats = useMemo(
      () => (Array.isArray(format) ? format : [format]),
      [format]
    );
    const accept = formats.map(item => accepts[item]).join(',');
    const validate = useCallback(
      (file: File) => formats.every(item => validates[item](file)),
      [formats]
    );
    const [dragging, setDragging] = useState(false);
    const [validationError, setValidationError] =
      useState<ValidationError | null>(null);

    const handleValidate = useCallback(
      async (_files?: FileList | null) => {
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
      },
      [t, format, validate, maxFileSize, onChange]
    );

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
      handleValidate(e.target?.files);
      setInputKey(Date.now());
    };

    const handleClearError = () => {
      setValidationError(null);
    };

    const handleDragIn = useCallback((e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      draggingCount.current += 1;
      if (e.dataTransfer?.items?.length) {
        setDragging(true);
      }
    }, []);

    const handleDragOut = useCallback((e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      draggingCount.current -= 1;
      if (draggingCount.current > 0) return;
      setDragging(false);
    }, []);

    const handleDrop = useCallback(
      (e: DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setDragging(false);
        draggingCount.current = 0;
        handleValidate(e.dataTransfer?.files);
      },
      [handleValidate]
    );

    const handleDrag = useCallback((e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
    }, []);

    useEffect(() => {
      const el = dropzone.current;
      if (!el) return;
      el.addEventListener('dragenter', handleDragIn);
      el.addEventListener('dragleave', handleDragOut);
      el.addEventListener('dragover', handleDrag);
      el.addEventListener('drop', handleDrop);

      return () => {
        el.removeEventListener('dragenter', handleDragIn);
        el.removeEventListener('dragleave', handleDragOut);
        el.removeEventListener('dragover', handleDrag);
        el.removeEventListener('drop', handleDrop);
      };
    }, [handleDragIn, handleDragOut, handleDrag, handleDrop]);

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
