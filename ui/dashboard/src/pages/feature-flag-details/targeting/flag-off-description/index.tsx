import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import Button from 'components/button';

const FlagOffDescription = () => {
  const { t } = useTranslation(['form']);
  const { watch, setValue } = useFormContext();
  const isShowRules = watch('isShowRules');
  return (
    <div
      className={cn(
        'flex-center w-full gap-x-2 py-2 typo-para-medium text-gray-600',
        {
          'flex-col gap-y-4': !isShowRules
        }
      )}
    >
      <p>{t('targeting.flag-off-desc')}</p>
      <Button
        variant="text"
        type="button"
        className="w-fit h-4 p-0 underline"
        onClick={() => setValue('isShowRules', !isShowRules)}
      >
        {t(`targeting.${isShowRules ? 'hide-rules' : 'view-targeting-rules'}`)}
      </Button>
    </div>
  );
};

export default FlagOffDescription;
