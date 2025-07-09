import { useCallback } from 'react';
import { urls } from 'configs';
import { useToast } from 'hooks';
import { getLanguage, Language, useTranslation } from 'i18n';
import { copyToClipBoard } from 'utils/function';
import { IconCopy } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';

const SDKApiEndpoint = () => {
  const { t } = useTranslation(['form', 'message']);
  const isJapaneseLanguage = getLanguage() === Language.JAPANESE;
  const { notify } = useToast();
  const handleCopy = useCallback(() => {
    copyToClipBoard(urls.API_ENDPOINT || '');
    notify({
      message: t('message:copied')
    });
  }, [urls]);
  return (
    <div className="flex items-center gap-x-2 p-2 bg-gray-100 rounded">
      <p className="typo-para-small text-gray-600 whitespace-nowrap">
        {t('sdk-api-endpoint')}
        {isJapaneseLanguage ? '：' : ':'}
      </p>
      <div className="flex items-center gap-x-1">
        <p className="typo-para-small text-primary-500">{urls.API_ENDPOINT}</p>
        <Button variant={'grey'} className="size-fit p-0" onClick={handleCopy}>
          <Icon icon={IconCopy} size={'xs'} />
        </Button>
      </div>
    </div>
  );
};

export default SDKApiEndpoint;
