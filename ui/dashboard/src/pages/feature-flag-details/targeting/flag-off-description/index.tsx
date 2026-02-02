import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import Button from 'components/button';
import { TargetingSchema } from '../form-schema';

const FlagOffDescription = ({
  isShowRules,
  setIsShowRules
}: {
  isShowRules: boolean;
  setIsShowRules: (value: boolean) => void;
}) => {
  const { t } = useTranslation(['form']);
  const { watch } = useFormContext();

  const prerequisiteCount = watch('prerequisites')?.length;
  const individualRulesWatch: TargetingSchema['individualRules'] = [
    ...watch('individualRules')
  ];
  const segmentRuleCount: number = watch('segmentRules')?.length || 0;

  const individualRuleCount = useMemo(() => {
    const count = individualRulesWatch?.reduce((acc, curr) => {
      if (curr?.users?.length) acc++;
      return acc;
    }, 0);
    return count || 0;
  }, [individualRulesWatch]);

  const hiddenRuleDesc = useMemo(() => {
    const getText = (count: number, key: string) =>
      `${count} ${t(count > 1 ? key : key.slice(0, -1))}`;
    let text = '';

    text +=
      prerequisiteCount > 0
        ? `${getText(prerequisiteCount, 'feature-flags.prerequisites')}${
            individualRuleCount > 0 || segmentRuleCount > 0 ? ', ' : ''
          }`
        : '';
    text +=
      individualRuleCount > 0
        ? `${getText(individualRuleCount, 'targeting.targets')}${segmentRuleCount > 0 ? ', ' : ''}`
        : '';
    text +=
      segmentRuleCount > 0 ? getText(segmentRuleCount, 'targeting.rules') : '';

    return text;
  }, [prerequisiteCount, individualRuleCount, segmentRuleCount]);

  return (
    <div
      className={cn(
        'inline sm:flex-center w-full gap-x-2 py-2 typo-para-medium text-gray-600',
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
          onClick={() => setIsShowRules(!isShowRules)}
        >
          {t(
            `targeting.${isShowRules ? 'hide-rules' : 'view-targeting-rules'}`
          )}
        </Button>
        {!isShowRules && !!hiddenRuleDesc && (
          <p className="typo-para-small">({hiddenRuleDesc})</p>
        )}
      </div>
    </div>
  );
};

export default FlagOffDescription;
