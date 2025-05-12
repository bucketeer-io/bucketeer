import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import Button from 'components/button';

const FlagOffDescription = () => {
  const { t } = useTranslation(['form']);
  const { watch, setValue } = useFormContext();
  const isShowRules = watch('isShowRules');
  const prerequisiteCount = watch('prerequisites')?.length;
  const individualRuleCount = watch('individualRules')?.length;
  const segmentRuleCount = watch('segmentRules')?.length;

  const hiddenRuleDesc = useMemo(() => {
    const getText = (count: number, key: string) =>
      `${count} ${t(count > 1 ? key : key.slice(0, -1))}`;
    let text = '';

    text +=
      prerequisiteCount > 0
        ? `${getText(prerequisiteCount, 'feature-flags.prerequisites')}, `
        : '';
    text +=
      individualRuleCount > 0
        ? `${getText(individualRuleCount, 'targeting.targets')}, `
        : '';
    text +=
      segmentRuleCount > 0 ? getText(segmentRuleCount, 'targeting.rules') : '';

    return text;
  }, [prerequisiteCount, individualRuleCount, segmentRuleCount]);

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
      <div className="flex-center flex-col gap-y-2">
        <Button
          variant="text"
          type="button"
          className="w-fit h-4 p-0 underline"
          onClick={() => setValue('isShowRules', !isShowRules)}
        >
          {t(
            `targeting.${isShowRules ? 'hide-rules' : 'view-targeting-rules'}`
          )}
        </Button>
        {!isShowRules && <p className="typo-para-small">({hiddenRuleDesc})</p>}
      </div>
    </div>
  );
};

export default FlagOffDescription;
