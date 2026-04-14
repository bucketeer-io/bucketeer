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
    (
      value: FlagSwitchVariationType,
      isInitial = false,
      shouldRegenerateIds = true
    ) => {
      const currentOnVariationId = watch('defaultOnVariation');
      const currentOffVariationId = watch('defaultOffVariation');
      const onVariationId = shouldRegenerateIds ? uuid() : currentOnVariationId;
      const offVariation = shouldRegenerateIds ? uuid() : currentOffVariationId;
      const setValueOptions = isInitial ? { shouldDirty: false } : {};

      const previousSwitchVariation = watch('switchVariationType');

      resetField('variations');
      setValue('switchVariationType', value, setValueOptions);
      if (shouldRegenerateIds) {
        setValue('defaultOnVariation', onVariationId, setValueOptions);
        setValue('defaultOffVariation', offVariation, setValueOptions);
      }

      // Handle EXPERIMENT template
      if (value === FlagSwitchVariationType.EXPERIMENT) {
        const currentVariations = watch('variations');
        const thirdVariationId =
          shouldRegenerateIds || !currentVariations?.[2]?.id
            ? uuid()
            : currentVariations[2].id;

        // Set variationType to STRING when switching TO EXPERIMENT from another template
        if (
          previousSwitchVariation !== FlagSwitchVariationType.EXPERIMENT ||
          isInitial
        ) {
          setValue('variationType', 'STRING', setValueOptions);
        }

        const variationTypeToUse = watch('variationType');

        switch (variationTypeToUse) {
          case 'BOOLEAN':
            // For BOOLEAN, only show 2 variations
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('control'),
                  value: 'true'
                },
                {
                  id: offVariation,
                  name: `${t('treatment')} 1`,
                  value: 'false'
                }
              ],
              setValueOptions
            );

          case 'STRING':
            return setValue(
              'variations',
              [
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
                  id: thirdVariationId,
                  name: `${t('treatment')} 2`,
                  value: 'treatment-2'
                }
              ],
              setValueOptions
            );

          case 'NUMBER':
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('control'),
                  value: '0'
                },
                {
                  id: offVariation,
                  name: `${t('treatment')} 1`,
                  value: '1'
                },
                {
                  id: thirdVariationId,
                  name: `${t('treatment')} 2`,
                  value: '2'
                }
              ],
              setValueOptions
            );

          case 'JSON':
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('control'),
                  value: '{"group": "control"}'
                },
                {
                  id: offVariation,
                  name: `${t('treatment')} 1`,
                  value: '{"group": "treatment-1"}'
                },
                {
                  id: thirdVariationId,
                  name: `${t('treatment')} 2`,
                  value: '{"group": "treatment-2"}'
                }
              ],
              setValueOptions
            );

          case 'YAML':
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('control'),
                  value: 'group: control'
                },
                {
                  id: offVariation,
                  name: `${t('treatment')} 1`,
                  value: 'group: treatment-1'
                },
                {
                  id: thirdVariationId,
                  name: `${t('treatment')} 2`,
                  value: 'group: treatment-2'
                }
              ],
              setValueOptions
            );

          default:
            // Fallback to string
            return setValue(
              'variations',
              [
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
                  id: thirdVariationId,
                  name: `${t('treatment')} 2`,
                  value: 'treatment-2'
                }
              ],
              setValueOptions
            );
        }
      }

      // Handle CUSTOM template - set defaults based on current variationType
      if (value === FlagSwitchVariationType.CUSTOM) {
        const currentVariationType = watch('variationType');

        setValue('defaultOnVariation', onVariationId, setValueOptions);
        setValue('defaultOffVariation', offVariation, setValueOptions);

        switch (currentVariationType) {
          case 'BOOLEAN':
            return setValue(
              'variations',
              [
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
              ],
              setValueOptions
            );

          case 'STRING':
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('variation-n', { number: 1 }),
                  value: 'variation-1'
                },
                {
                  id: offVariation,
                  name: t('variation-n', { number: 2 }),
                  value: 'variation-2'
                }
              ],
              setValueOptions
            );

          case 'NUMBER':
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('variation-n', { number: 1 }),
                  value: '1'
                },
                {
                  id: offVariation,
                  name: t('variation-n', { number: 2 }),
                  value: '2'
                }
              ],
              setValueOptions
            );

          case 'JSON':
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('variation-n', { number: 1 }),
                  value: '{"variation": "variation-1"}'
                },
                {
                  id: offVariation,
                  name: t('variation-n', { number: 2 }),
                  value: '{"variation": "variation-2"}'
                }
              ],
              setValueOptions
            );

          case 'YAML':
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('variation-n', { number: 1 }),
                  value: 'variation: variation-1'
                },
                {
                  id: offVariation,
                  name: t('variation-n', { number: 2 }),
                  value: 'variation: variation-2'
                }
              ],
              setValueOptions
            );

          default:
            // Fallback to string if variationType is not set
            return setValue(
              'variations',
              [
                {
                  id: onVariationId,
                  name: t('variation-n', { number: 1 }),
                  value: 'variation-1'
                },
                {
                  id: offVariation,
                  name: t('variation-n', { number: 2 }),
                  value: 'variation-2'
                }
              ],
              setValueOptions
            );
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
            return setValue(
              'variations',
              [
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
              ],
              setValueOptions
            );

          case 'STRING':
            return setValue(
              'variations',
              [
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
              ],
              setValueOptions
            );

          case 'NUMBER':
            return setValue(
              'variations',
              [
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
              ],
              setValueOptions
            );

          case 'JSON':
            return setValue(
              'variations',
              [
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
              ],
              setValueOptions
            );

          case 'YAML':
            return setValue(
              'variations',
              [
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
              ],
              setValueOptions
            );

          default:
            // Fallback to boolean if variationType is not set
            return setValue(
              'variations',
              [
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
              ],
              setValueOptions
            );
        }
      }
    },
    [watch, setValue, resetField, t]
  );

  // Trigger CUSTOM template on initial mount
  useEffect(() => {
    if (isInitialMount.current) {
      isInitialMount.current = false;
      handleSwitchVariation(FlagSwitchVariationType.CUSTOM, true);
    }
  }, [handleSwitchVariation]);

  // Watch for variationType changes and update variations when CUSTOM, RELEASE, KILL_SWITCH, or EXPERIMENT template is active
  useEffect(() => {
    if (
      !isInitialMount.current &&
      (currentSwitchVariation === FlagSwitchVariationType.CUSTOM ||
        currentSwitchVariation === FlagSwitchVariationType.RELEASE ||
        currentSwitchVariation === FlagSwitchVariationType.KILL_SWITCH ||
        currentSwitchVariation === FlagSwitchVariationType.EXPERIMENT)
    ) {
      // Don't regenerate IDs when only updating variation values
      handleSwitchVariation(currentSwitchVariation, false, false);
    }
  }, [currentVariationType, handleSwitchVariation, currentSwitchVariation]);

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
                        <button
                          type="button"
                          className="flex-center cursor-pointer"
                          onClick={e => {
                            e.preventDefault();
                            e.stopPropagation();
                          }}
                          aria-label={t('template-info-aria-label')}
                        >
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </button>
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
                        <button
                          type="button"
                          className="flex-center cursor-pointer"
                          onClick={e => {
                            e.preventDefault();
                            e.stopPropagation();
                          }}
                          aria-label={t('template-info-aria-label')}
                        >
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </button>
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
                        <button
                          type="button"
                          className="flex-center cursor-pointer"
                          onClick={e => {
                            e.preventDefault();
                            e.stopPropagation();
                          }}
                          aria-label={t('template-info-aria-label')}
                        >
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </button>
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
                        <button
                          type="button"
                          className="flex-center cursor-pointer"
                          onClick={e => {
                            e.preventDefault();
                            e.stopPropagation();
                          }}
                          aria-label={t('template-info-aria-label')}
                        >
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </button>
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
