import { useCallback } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import { FlagFormSchema } from 'pages/create-flag/form-schema';
import { FlagSwitchVariationType } from 'pages/create-flag/types';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

const buttonCls =
  'typo-para-medium !text-gray-600 !shadow-none border border-gray-200 hover:border-gray-400 disabled:pointer-events-none';
const buttonActiveCls =
  '!text-accent-pink-500 border-accent-pink-500 hover:!text-accent-pink-500 hover:border-accent-pink-500 disabled:pointer-events-none';

const VariationsSwitch = () => {
  const { t } = useTranslation(['form', 'common']);

  const { watch, setValue, resetField } = useFormContext<FlagFormSchema>();

  const currentSwitchVariation = watch('switchVariationType');

  const handleSwitchVariation = useCallback(
    (value: FlagSwitchVariationType) => {
      const onVariationId = uuid();
      const offVariation = uuid();
      resetField('variations');
      setValue('switchVariationType', value);
      setValue('defaultOnVariation', onVariationId);
      setValue('defaultOffVariation', offVariation);
      if (value === FlagSwitchVariationType.EXPERIMENT) {
        setValue('variationType', 'STRING');

        return setValue('variations', [
          {
            id: onVariationId,
            name: t('control'),
            value: '',
            description: ''
          },
          {
            id: offVariation,
            name: `${t('treatment')} 1`,
            value: '',
            description: ''
          },
          {
            id: uuid(),
            name: `${t('treatment')} 2`,
            value: '',
            description: ''
          }
        ]);
      }
      const isRelease = value === FlagSwitchVariationType.RELEASE;
      const isKillSwitch = value === FlagSwitchVariationType.KILL_SWITCH;
      setValue('variationType', 'BOOLEAN');
      return setValue('variations', [
        {
          id: onVariationId,
          name: isRelease ? t('available') : isKillSwitch ? t('enabled') : '',
          value: 'true',
          description: ''
        },
        {
          id: offVariation,
          name: isRelease
            ? t('unavailable')
            : isKillSwitch
              ? t('disabled')
              : '',
          value: 'false',
          description: ''
        }
      ]);
    },
    []
  );

  return (
    <div className="flex items-center w-full justify-between">
      <p className="typo-para-medium text-gray-700">
        {t('feature-flags.flag-variations')}
      </p>
      <Form.Field
        name="switchVariationType"
        render={() => (
          <Form.Item className="py-0">
            <Form.Control>
              <div className="flex items-center">
                <Button
                  variant={'secondary-2'}
                  type="button"
                  className={cn(
                    'rounded-r-none',
                    buttonCls,
                    currentSwitchVariation === FlagSwitchVariationType.CUSTOM &&
                      buttonActiveCls
                  )}
                  onClick={() =>
                    handleSwitchVariation(FlagSwitchVariationType.CUSTOM)
                  }
                >
                  <div className="flex items-center gap-x-1">
                    {t(`custom`)}
                    <Tooltip
                      align="start"
                      trigger={
                        <div className="flex-center">
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </div>
                      }
                      content={t('template-tooltip.custom')}
                      className="!z-[100] max-w-[300px]"
                    />
                  </div>
                </Button>
                <Button
                  variant={'secondary-2'}
                  type="button"
                  className={cn(
                    'rounded-none',
                    buttonCls,
                    currentSwitchVariation ===
                      FlagSwitchVariationType.RELEASE && buttonActiveCls
                  )}
                  onClick={() =>
                    handleSwitchVariation(FlagSwitchVariationType.RELEASE)
                  }
                >
                  <div className="flex items-center gap-x-1">
                    {t(`release`)}
                    <Tooltip
                      align="start"
                      trigger={
                        <div className="flex-center">
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </div>
                      }
                      content={t('template-tooltip.release')}
                      className="!z-[100] max-w-[300px]"
                    />
                  </div>
                </Button>
                <Button
                  variant={'secondary-2'}
                  type="button"
                  className={cn(
                    'rounded-none',
                    buttonCls,
                    currentSwitchVariation ===
                      FlagSwitchVariationType.KILL_SWITCH && buttonActiveCls
                  )}
                  onClick={() =>
                    handleSwitchVariation(FlagSwitchVariationType.KILL_SWITCH)
                  }
                >
                  <div className="flex items-center gap-x-1">
                    {t(`kill-switch`)}
                    <Tooltip
                      align="start"
                      trigger={
                        <div className="flex-center">
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </div>
                      }
                      content={t('template-tooltip.kill-switch')}
                      className="!z-[100] max-w-[300px]"
                    />
                  </div>
                </Button>
                <Button
                  variant={'secondary-2'}
                  type="button"
                  className={cn(
                    'rounded-l-none',
                    buttonCls,
                    currentSwitchVariation ===
                      FlagSwitchVariationType.EXPERIMENT && buttonActiveCls
                  )}
                  onClick={() =>
                    handleSwitchVariation(FlagSwitchVariationType.EXPERIMENT)
                  }
                >
                  <div className="flex items-center gap-x-1">
                    {t(`common:source-type.experiment`)}
                    <Tooltip
                      align="start"
                      trigger={
                        <div className="flex-center">
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </div>
                      }
                      content={t('template-tooltip.experiment')}
                      className="!z-[100] max-w-[300px]"
                    />
                  </div>
                </Button>
              </div>
            </Form.Control>
          </Form.Item>
        )}
      />
    </div>
  );
};

export default VariationsSwitch;
