import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { IconToastWarning } from '@icons';
import Icon from 'components/icon';

const RolloutWarning = ({
  urlCode,
  experiments,
  hasRolloutRunning
}: {
  urlCode: string;
  experiments: Experiment[];
  hasRolloutRunning: boolean;
}) => {
  const { t } = useTranslation(['form']);
  return (
    <div className="flex w-full gap-x-3 p-4 rounded-md bg-accent-yellow-50 typo-para-small text-accent-yellow-700">
      <Icon icon={IconToastWarning} />
      <div className="flex flex-col flex-1">
        <p className="font-bold">{t('rollout-warning-title')}</p>
        <ul className="flex flex-col w-full gap-y-2 list-disc pl-5 mt-2">
          {experiments.length > 0 && (
            <li>
              <Trans
                i18nKey={'form:rollout-experiment-warning'}
                components={{
                  comp: (
                    <Link
                      to={`/${urlCode}${PAGE_PATH_EXPERIMENTS}/${experiments[0].id}/results`}
                      className="!text-primary-500 underline"
                    />
                  )
                }}
              />
            </li>
          )}
          {hasRolloutRunning && <li>{t('rollout-warning-desc')}</li>}
        </ul>

        <div className="mt-4">
          <Trans
            i18nKey={'form:more-info-see-document'}
            components={{
              comp: (
                <Link
                  to="https://docs.bucketeer.io/feature-flags/creating-feature-flags/auto-operation/progressive-rollout"
                  target="_blank"
                  className="!text-primary-500 underline"
                />
              )
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default RolloutWarning;
