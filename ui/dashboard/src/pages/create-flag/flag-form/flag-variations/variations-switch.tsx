import { useCallback, useEffect, useRef } from 'react';
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
  const currentVariationType = watch('variationType');
  const isInitialMount = useRef(true);

  const handleSwitchVariation = useCallback(
    (value: FlagSwitchVariationType) => {
      const onVariationId = uuid();
      const offVariation = uuid();
      resetField('variations');
      setValue('switchVariationType', value);
      setValue('defaultOnVariation', onVariationId);
      setValue('defaultOffVariation', offVariation);

      // Handle EXPERIMENT template
      if (value === FlagSwitchVariationType.EXPERIMENT) {
        setValue('variationType', 'STRING');

        return setValue('variations', [
          {
            id: onVariationId,
            name: t('control'),
            value: 'control'
          },
          {
            id: offVariation,
            name: `${t('treatment')} 1`,
            value: 'treatment-1'
          },
          {
            id: uuid(),
            name: `${t('treatment')} 2`,
            value: 'treatment-2'
          }
        ]);
      }

      // Handle CUSTOM template - set defaults based on current variationType
      if (value === FlagSwitchVariationType.CUSTOM) {
        const currentVariationType = watch('variationType');

        switch (currentVariationType) {
          case 'BOOLEAN':
            return setValue('variations', [
              {
                id: onVariationId,
                name: 'true',
                value: 'true'
              },
              {
                id: offVariation,
                name: 'false',
                value: 'false'
              }
            ]);

          case 'STRING':
            return setValue('variations', [
              {
                id: onVariationId,
                name: 'Variation 1',
                value: 'variation-1'
              },
              {
                id: offVariation,
                name: 'Variation 2',
                value: 'variation-2'
              }
            ]);

          case 'NUMBER':
            return setValue('variations', [
              {
                id: onVariationId,
                name: 'Variation 1',
                value: '1'
              },
              {
                id: offVariation,
                name: 'Variation 2',
                value: '2'
              }
            ]);

          case 'JSON':
            return setValue('variations', [
              {
                id: onVariationId,
                name: 'Variation 1',
                value: '{"variation": "variation-1"}'
              },
              {
                id: offVariation,
                name: 'Variation 2',
                value: '{"variation": "variation-2"}'
              }
            ]);

          case 'YAML':
            return setValue('variations', [
              {
                id: onVariationId,
                name: 'Variation 1',
                value: 'variation: variation-1'
              },
              {
                id: offVariation,
                name: 'Variation 2',
                value: 'variation: variation-2'
              }
            ]);

          default:
            // Fallback to string if variationType is not set
            return setValue('variations', [
              {
                id: onVariationId,
                name: 'Variation 1',
                value: 'variation-1'
              },
              {
                id: offVariation,
                name: 'Variation 2',
                value: 'variation-2'
              }
            ]);
        }
      }

      // Handle RELEASE and KILL_SWITCH templates
      const isRelease = value === FlagSwitchVariationType.RELEASE;
      const isKillSwitch = value === FlagSwitchVariationType.KILL_SWITCH;

      if (isRelease || isKillSwitch) {
        const currentVariationType = watch('variationType');
        const onName = isRelease ? t('available') : t('enabled');
        const offName = isRelease ? t('unavailable') : t('disabled');

        switch (currentVariationType) {
          case 'BOOLEAN':
            return setValue('variations', [
              {
                id: onVariationId,
                name: onName,
                value: 'true'
              },
              {
                id: offVariation,
                name: offName,
                value: 'false'
              }
            ]);

          case 'STRING':
            return setValue('variations', [
              {
                id: onVariationId,
                name: onName,
                value: 'true'
              },
              {
                id: offVariation,
                name: offName,
                value: 'false'
              }
            ]);

          case 'NUMBER':
            return setValue('variations', [
              {
                id: onVariationId,
                name: onName,
                value: '1'
              },
              {
                id: offVariation,
                name: offName,
                value: '0'
              }
            ]);

          case 'JSON':
            return setValue('variations', [
              {
                id: onVariationId,
                name: onName,
                value: '{"status": true}'
              },
              {
                id: offVariation,
                name: offName,
                value: '{"status": false}'
              }
            ]);

          case 'YAML':
            return setValue('variations', [
              {
                id: onVariationId,
                name: onName,
                value: 'status: true'
              },
              {
                id: offVariation,
                name: offName,
                value: 'status: false'
              }
            ]);

          default:
            // Fallback to boolean if variationType is not set
            return setValue('variations', [
              {
                id: onVariationId,
                name: onName,
                value: 'true'
              },
              {
                id: offVariation,
                name: offName,
                value: 'false'
              }
            ]);
        }
      }
    },
    [watch, setValue, resetField, t]
  );

  // Trigger CUSTOM template on initial mount
  useEffect(() => {
    if (isInitialMount.current) {
      isInitialMount.current = false;
      handleSwitchVariation(FlagSwitchVariationType.CUSTOM);
    }
  }, [handleSwitchVariation]);

  // Watch for variationType changes and update variations when CUSTOM, RELEASE, or KILL_SWITCH template is active
  useEffect(() => {
    if (
      !isInitialMount.current &&
      (currentSwitchVariation === FlagSwitchVariationType.CUSTOM ||
        currentSwitchVariation === FlagSwitchVariationType.RELEASE ||
        currentSwitchVariation === FlagSwitchVariationType.KILL_SWITCH)
    ) {
      handleSwitchVariation(currentSwitchVariation);
    }
  }, [currentVariationType, currentSwitchVariation, handleSwitchVariation]);

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
