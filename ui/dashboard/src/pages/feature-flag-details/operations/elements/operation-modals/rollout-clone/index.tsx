import { ReactNode, useMemo } from 'react';
import { Trans } from 'react-i18next';
import { getLanguage, useTranslation } from 'i18n';
import { Rollout } from '@types';
import { RolloutTypeMap } from 'pages/feature-flag-details/operations/types';
import {
  getDateTimeDisplay,
  numberToOrdinalWord
} from 'pages/feature-flag-details/operations/utils';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import SlideModal from 'components/modal/slide';

const InfoItem = ({ title, desc }: { title: ReactNode; desc: string }) => (
  <div className="flex flex-col flex-1 typo-para-small gap-y-2">
    <p className="text-gray-600 capitalize">{title}</p>
    <p className="text-gray-700 capitalize">{desc}</p>
  </div>
);

const RolloutCloneModal = ({
  isOpen,
  selectedData,
  onClose
}: {
  isOpen: boolean;
  selectedData: Rollout;
  onClose: () => void;
}) => {
  const { t } = useTranslation(['form', 'common']);
  const isLanguageJapanese = getLanguage() === 'ja';

  const { clause, type } = selectedData;

  const currentSchedule = useMemo(
    () =>
      clause.schedules.find(
        (item, index) =>
          item.triggeredAt !== '0' &&
          (clause.schedules[index + 1]?.triggeredAt === '0' ||
            !clause.schedules[index + 1])
      ),
    [clause]
  );

  return (
    <SlideModal
      title={t(`common:source-type.progressive-rollout`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col gap-y-6 w-full p-5 pb-28">
        {type === RolloutTypeMap.TEMPLATE_SCHEDULE ? (
          <>
            <p className="typo-head-bold-small text-gray-800">
              {t('general-info')}
            </p>
            <div className="flex items-center w-full justify-between">
              <InfoItem
                title={t('feature-flags.start-date')}
                desc={getDateTimeDisplay(clause.schedules[0]?.executeAt)?.date}
              />
              <InfoItem
                title={t('increments')}
                desc={`${clause?.increments}%`}
              />
              <InfoItem
                title={t('frequency')}
                desc={clause?.interval?.toLowerCase() || ''}
              />
            </div>
          </>
        ) : (
          <div className="flex flex-col w-full divide-y divide-gray-200">
            <div className="flex items-center w-full justify-between pb-4">
              <InfoItem title={t('set')} desc={t('manual')} />
              <InfoItem
                title={t('current-progress')}
                desc={`${(currentSchedule?.weight || 0) / 1000}%`}
              />
            </div>
            {clause.schedules.map((item, index) => (
              <div
                key={index}
                className="flex items-center w-full justify-between py-4 first:pt-0 last:pb-0"
              >
                <InfoItem
                  title={
                    <Trans
                      i18nKey={'form:ordinal-increment'}
                      values={{
                        ordinal: isLanguageJapanese
                          ? index + 1
                          : numberToOrdinalWord(index + 1)
                      }}
                    />
                  }
                  desc={`${item.weight / 1000}%`}
                />
                <InfoItem
                  title={t('feature-flags.date')}
                  desc={getDateTimeDisplay(item.executeAt)?.date || ''}
                />
              </div>
            ))}
          </div>
        )}
      </div>
      <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
        <ButtonBar
          primaryButton={
            <Button variant="secondary" onClick={onClose}>
              {t(`common:cancel`)}
            </Button>
          }
          secondaryButton={
            <Button type="submit">{t(`clone-operation`)}</Button>
          }
        />
      </div>
    </SlideModal>
  );
};

export default RolloutCloneModal;
