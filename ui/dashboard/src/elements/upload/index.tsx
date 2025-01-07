import React, {
  ChangeEvent,
  ReactNode,
  useCallback,
  useMemo,
  useRef
} from 'react';
import { Trans } from 'react-i18next';
import { useToast } from 'hooks';
import { formatFileSize } from 'utils/converts';
import { cn } from 'utils/style';
import { IconFolder, IconUpload } from '@icons';
import Icon from 'components/icon';

type Props = {
  files?: File[];
  accept?: string[];
  multiple?: boolean;
  children?: ReactNode;
  className?: string;
  onChange: (files: File[]) => void;
};

const UploadFiles = ({
  files,
  accept = ['.csv', '.txt'],
  multiple = false,
  children,
  className,
  onChange
}: Props) => {
  const uploadRef = useRef<HTMLInputElement>(null);
  const { notify } = useToast();

  const formatAccept = useMemo(() => {
    if (typeof accept === 'string') return accept;
    const _accept = accept.join(', ');
    return _accept;
  }, [accept]);

  const validateFileType = useCallback(
    (files: File[]) => {
      const isInvalidField = files.find(file => {
        const fileExtension = file?.name
          .substring(file.name?.lastIndexOf('.'))
          ?.toLowerCase();
        return !accept.includes(fileExtension);
      });
      if (isInvalidField) {
        notify({
          message:
            'Some files were not uploaded because they are of an unsupported type. Please try again with the allowed formats.',
          messageType: 'error',
          toastType: 'toast'
        });
        return false;
      }
      return true;
    },
    [accept]
  );

  const handleDrop = useCallback(
    (event: React.DragEvent<HTMLDivElement>) => {
      event.preventDefault();
      event.stopPropagation();
      if (event.dataTransfer.files) {
        const filesArray = Array.from(event.dataTransfer.files);
        const isValidate = validateFileType(filesArray);
        if (isValidate) onChange(filesArray);
      }
    },
    [accept]
  );

  const handleChange = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      if (e.target.files) {
        const filesArray = Array.from(e.target.files);
        const isValidate = validateFileType(filesArray);
        if (isValidate) onChange(filesArray);
      }
    },
    [accept]
  );

  const preventAndStopFnc = useCallback(
    (event: React.DragEvent<HTMLDivElement>) => {
      event.preventDefault();
      event.stopPropagation();
    },
    []
  );

  return (
    <div className={cn('flex flex-col w-full gap-y-2 h-fit', className)}>
      <div className="flex h-[200px] min-h-[200px] w-full border border-dashed border-gray-400 rounded-lg bg-gray-100">
        <input
          ref={uploadRef}
          type="file"
          multiple={multiple}
          accept={formatAccept}
          onChange={handleChange}
          className="hidden"
        />
        <div
          className="flex-center flex-col size-full gap-y-4 cursor-pointer"
          onDrop={handleDrop}
          onDragOver={preventAndStopFnc}
          onDragEnter={preventAndStopFnc}
          onDragLeave={preventAndStopFnc}
          onClick={() => uploadRef?.current?.click()}
        >
          <Icon icon={IconUpload} size={'xl'} color="primary-500" />
          {children || (
            <div className="flex-center flex-col gap-y-2 typo-para-medium text-gray-600">
              <div>
                <Trans
                  i18nKey="form:upload-files"
                  values={{ text: 'browse' }}
                  components={{
                    underline: <u className="text-primary-500 inline" />
                  }}
                />
              </div>
              <Trans i18nKey="form:any-text" values={{ text: 'CSV, TXT' }} />
            </div>
          )}
        </div>
      </div>
      {files && files?.length > 0 && (
        <div className="flex flex-col w-full gap-y-2">
          {files.map(file => (
            <div className="flex items-center w-full gap-x-4 p-1 border-2 border-dashed border-gray-300 rounded">
              <Icon icon={IconFolder} size={'xl'} color="gray-500" />
              <div className="flex flex-col flex-1 gap-y-1 truncate">
                <p className="typo-para-small text-gray-800">{file.name}</p>
                <p className="typo-para-tiny text-gray-600">
                  {formatFileSize(file.size)}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default UploadFiles;
