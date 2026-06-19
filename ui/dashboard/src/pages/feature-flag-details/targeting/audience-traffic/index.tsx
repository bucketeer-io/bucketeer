import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

const AudienceTraffic = () => {
  const { t } = useTranslation(['form']);
  return (
    // Left-aligned to match the visual rhythm of the cards below: the title
    // sits flush with the card column, anchored to the start node on the spine.
    <div className="flex items-center w-full gap-x-2">
      <p className="typo-para-medium text-gray-700">
        {t('targeting.all-audience-traffic')}
      </p>
      <Tooltip
        content={t('targeting.tooltip.audience')}
        trigger={
          <div className="flex-center size-fit">
            <Icon icon={IconInfo} size="xxs" color="gray-500" />
          </div>
        }
        className="max-w-[450px] whitespace-pre-line"
      />
    </div>
  );
};

export default AudienceTraffic;
