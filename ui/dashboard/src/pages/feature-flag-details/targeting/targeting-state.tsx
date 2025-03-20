import { useMemo, useState } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconPencil } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import Switch from 'components/switch';

const TargetingState = () => {
  const { t } = useTranslation(['common', 'table']);

  const [isTargetingOn, setIsTargetingOn] = useState(true);
  const [serveValue, setServeValue] = useState(false);

  const serveOptions = useMemo(
    () => [
      { label: t('false'), value: false },
      { label: t('true'), value: true }
    ],
    []
  );

  const currentServe = useMemo(
    () => serveOptions.find(item => item.value === serveValue),
    [serveValue]
  );

  return (
    <div className="flex w-full gap-x-4">
      <div className="flex flex-col flex-1 gap-y-2">
        <p className="typo-para-small py-[14px] text-gray-500">
          <Trans
            i18nKey={'table:feature-flags.when-targeting'}
            values={{
              state: isTargetingOn ? 'ON' : 'OFF'
            }}
          />
        </p>
        <div className="flex items-center gap-x-2">
          <p className="typo-para-small text-gray-600 uppercase mr-2">
            {t('table:feature-flags.serve')}
          </p>
          <div
            className={cn(
              'flex items-center gap-x-2 px-3 py-[14px] bg-gray-100 rounded-lg'
            )}
          >
            <FlagVariationPolygon index={serveValue ? 1 : 0} />
            <p className="typo-para-medium leading-5 text-gray-700">
              {currentServe?.label}
            </p>
          </div>
          <Popover
            trigger={
              <div className="flex-center size-12 border border-gray-200 rounded-lg">
                <Icon icon={IconPencil} color="gray-600" size={'sm'} />
              </div>
            }
            align="start"
            sideOffset={8}
            className="w-20"
            options={serveOptions}
            onClick={value => setServeValue(value as boolean)}
          />
        </div>
      </div>
      <div className="flex items-center h-fit gap-x-4">
        <Switch
          checked={isTargetingOn}
          onCheckedChange={value => setIsTargetingOn(value)}
        />
        <Button>{t('submit')}</Button>
      </div>
    </div>
  );
};

export default TargetingState;
