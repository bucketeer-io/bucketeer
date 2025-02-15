import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import Divider from 'components/divider';
import { DefineAudienceField } from '..';

export interface DefineAudienceProps {
  field: DefineAudienceField;
}

const DefineAudience = ({ field }: DefineAudienceProps) => {
  const { t } = useTranslation(['form', 'common']);

  return (
    <div className="flex flex-col gap-y-5 mt-2.5">
      <p className="text-gray-800 typo-head-bold-small">
        {t('experiments.audience')}
      </p>
      <div className="flex flex-col w-full gap-y-4 typo-para-small leading-[14px] text-gray-600">
        <p className="text-gray-800 typo-head-bold-small leading-4">
          Default Rule
        </p>
        <p className="typo-para-medium leading-5 text-gray-500">
          {t('experiments.define-audience.any-traffic')}
        </p>
      </div>
      <div className="w-full p-4 bg-gray-100 rounded-lg">
        <div className="flex flex-col w-full gap-y-4 typo-para-small leading-[14px] text-gray-600">
          <p className="uppercase">{t('experiments.audience-included')}</p>
          <div className="w-full h-3 p-[1px] border border-gray-400 rounded-full bg-gray-100">
            <div
              className={cn('h-full bg-primary-500 rounded-l-full', {
                'rounded-r-full': field.value?.inExperiment === 100
              })}
              style={{
                width: `${field.value?.inExperiment}%`
              }}
            />
          </div>
          <div className="flex items-center w-full gap-x-4">
            <div className="flex items-center gap-x-2">
              <div className="flex-center size-5 m-0.5 rounded bg-primary-500" />
              <p>{`${t('experiments.define-audience.in-this-experiment')} - ${field.value?.inExperiment}%`}</p>
            </div>
            <div className="flex items-center gap-x-2">
              <div className="flex-center size-5 m-0.5 border border-gray-400 rounded bg-gray-100" />
              <p>{`${t('experiments.define-audience.not-in-experiment')} - ${field.value?.notInExperiment}%`}</p>
            </div>
          </div>
        </div>
        <Divider className="my-5" />
        <div className="flex items-center w-full gap-x-2 mt-4 typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
          <Trans
            i18nKey={'form:experiments.settings.in-this-experiment'}
            values={{
              percent: `${field.value?.inExperiment}%`
            }}
            components={{
              highlight: (
                <div className="flex-center size-fit p-3 rounded-lg typo-para-medium leading-5 text-gray-700 bg-gray-200" />
              )
            }}
          />
        </div>
        <div className="flex items-center w-full gap-x-1.5 mt-4 typo-para-medium leading-5 text-gray-600 whitespace-nowrap">
          <Trans
            i18nKey={'form:experiments.settings.not-in-experiment'}
            values={{
              percent: `${field.value?.notInExperiment}%`
            }}
            components={{
              highlight: (
                <div className="flex-center size-fit p-3 rounded-lg typo-para-medium leading-5 text-gray-700 bg-gray-200" />
              )
            }}
          />
          <div className="flex items-center gap-x-1.5">
            <div className="flex-center size-3 bg-accent-blue-500 rounded-sm rotate-45" />
            <p>test</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DefineAudience;
