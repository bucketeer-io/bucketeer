import { useCallback } from 'react';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { copyToClipBoard } from 'utils/function';
import { IconCopy } from '@icons';
import Icon from 'components/icon';

const HeaderDetailsID = ({ id, message }: { id: string; message?: string }) => {
  const { t } = useTranslation(['message']);
  const { notify } = useToast();

  const handleCopyId = useCallback(() => {
    copyToClipBoard(id);
    notify({
      message: t(message || 'copied')
    });
  }, [id, message]);
  return (
    <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 select-none pl-0 sm:pl-10 mt-2">
      {id}
      <div onClick={handleCopyId}>
        <Icon
          icon={IconCopy}
          size={'sm'}
          className="opacity-100 cursor-pointer"
        />
      </div>
    </div>
  );
};

export default HeaderDetailsID;
