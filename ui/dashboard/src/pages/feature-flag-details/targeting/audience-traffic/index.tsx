import { useTranslation } from 'i18n';
import { IconInfo, IconMember } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

const AudienceTraffic = () => {
  const { t } = useTranslation(['form']);
  return (
    <div className="flex-center w-full gap-x-2">
      <Icon icon={IconMember} size="sm" color="gray-500" />
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
        className="max-w-[450px]"
      />
    </div>
  );
};

export default AudienceTraffic;
