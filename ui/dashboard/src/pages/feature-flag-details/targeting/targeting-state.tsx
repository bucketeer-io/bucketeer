import { useState } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import Switch from 'components/switch';
import ServeDropdown from './serve-dropdown';

const TargetingState = () => {
  const [isTargetingOn, setIsTargetingOn] = useState(true);
  const [serveValue, setServeValue] = useState(0);

  const { t } = useTranslation(['common']);
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
        <ServeDropdown serveValue={serveValue} onChangeServe={setServeValue} />
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
