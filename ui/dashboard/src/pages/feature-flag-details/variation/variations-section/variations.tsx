import { useCallback, useMemo } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { Feature, OperationStatus, Rollout, StrategyType } from '@types';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import TextArea from 'components/textarea';
import { VariationForm } from '../form-schema';

const VariationLabel = ({
  index,
  className
}: {
  index: number;
  className?: string;
}) => (
  <div className={cn('flex items-center gap-x-2 text-gray-600', className)}>
    <FlagVariationPolygon index={index} />
    <Trans
      i18nKey={'form:feature-flags.variation'}
      values={{
        index: `${index + 1}`
      }}
    />
  </div>
);

const Variations = ({
  feature,
  rollouts
}: {
  feature: Feature;
  rollouts: Rollout[];
}) => {
  const { t } = useTranslation(['common', 'form']);

  const { control, watch } = useFormContext<VariationForm>();

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'variations',
    keyName: 'variationField'
  });

  const offVariation = watch('offVariation');

  const isBoolean = useMemo(
    () => feature.variationType === 'BOOLEAN',
    [feature]
  );
  const isJSON = useMemo(() => feature.variationType === 'JSON', [feature]);

  const formItemClassName = useMemo(
    () => 'flex flex-col pt-6 pb-0 h-full self-stretch',
    []
  );

  const onVariationIds = useMemo(() => {
    if (feature?.defaultStrategy) {
      const { fixedStrategy, rolloutStrategy, type } = feature.defaultStrategy;
      if (type === StrategyType.FIXED) return [fixedStrategy.variation];
      if (type === StrategyType.ROLLOUT)
        return rolloutStrategy.variations
          .filter(v => v.weight > 0)
          ?.map(item => item?.variation);
    }
    return [];
  }, [feature]);

  const ruleVariationIds = useMemo(() => {
    if (feature?.rules?.length) {
      const arr: string[] = [];
      feature.rules.forEach(rule => {
        const { strategy } = rule;
        if (strategy.type === StrategyType.FIXED) {
          arr.push(strategy.fixedStrategy.variation);
        } else if (strategy.type === StrategyType.ROLLOUT) {
          strategy.rolloutStrategy.variations.filter(item => {
            if (item.weight > 0) arr.push(item.variation);
          });
        }
      });
      return [...new Set(arr)];
    }
    return [];
  }, [feature]);

  const isDisableRemoveBtn = useCallback(
    (variationId: string) => {
      return (
        isBoolean ||
        [
          ...new Set([offVariation, ...onVariationIds, ...ruleVariationIds])
        ].includes(variationId)
      );
    },
    [isBoolean, onVariationIds, ruleVariationIds, offVariation]
  );
  const isProgressiveRolloutsRunningWaiting = (status: OperationStatus) =>
    ['RUNNING', 'WAITING'].includes(status);

  const isDisableAddBtn = useCallback(
    () =>
      isBoolean ||
      rollouts.filter(item => isProgressiveRolloutsRunningWaiting(item.status))
        ?.length > 0,
    [isBoolean, rollouts]
  );

  const onAddVariation = () => {
    append({
      id: uuid(),
      value: '',
      name: '',
      description: ''
    });
  };

  return (
    <>
      {fields.map((variation, variationIndex) => (
        <div
          key={variation.variationField}
          className="flex items-end w-full gap-x-2"
        >
          <Form.Field
            control={control}
            name={`variations.${variationIndex}.value`}
            render={({ field }) => {
              return (
                <Form.Item
                  className={cn(formItemClassName, 'w-[20%]', {
                    'w-[10%]': isBoolean,
                    'flex-1 w-full': isJSON
                  })}
                >
                  {isBoolean && (
                    <VariationLabel index={variationIndex} className="mb-3" />
                  )}
                  <Form.Label required={isBoolean}>
                    {isBoolean ? (
                      t(field.value ? 'true' : 'false')
                    ) : (
                      <VariationLabel index={variationIndex} />
                    )}
                  </Form.Label>
                  <Form.Control>
                    {isJSON ? (
                      <TextArea
                        {...field}
                        rows={3}
                        placeholder={t('form:feature-flags.value')}
                      />
                    ) : (
                      <Input
                        {...field}
                        disabled={isBoolean}
                        placeholder={t('form:feature-flags.value')}
                        className={cn('px-3', {
                          capitalize: isBoolean
                        })}
                      />
                    )}
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              );
            }}
          />
          <Form.Field
            control={control}
            name={`variations.${variationIndex}.name`}
            render={({ field }) => (
              <Form.Item
                className={cn(formItemClassName, 'w-[30%]', {
                  'flex-1 w-full': isJSON
                })}
              >
                <Form.Label required>{t('name')}</Form.Label>
                <Form.Control>
                  <Input {...field} placeholder={t('name')} />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={control}
            name={`variations.${variationIndex}.description`}
            render={({ field }) => (
              <Form.Item
                className={cn(formItemClassName, 'flex-1', {
                  'flex-1 w-full': isJSON
                })}
              >
                <Form.Label>{t('form:description')}</Form.Label>
                <Form.Control>
                  <Input {...field} placeholder={t('form:description')} />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Button
            variant="grey"
            size="icon"
            type="button"
            className="p-0 size-5 mb-4"
            disabled={isDisableRemoveBtn(variation.id)}
            onClick={() => remove(variationIndex)}
          >
            <Icon icon={IconTrash} size="sm" />
          </Button>
        </div>
      ))}
      <Button
        onClick={onAddVariation}
        variant="text"
        type="button"
        disabled={isDisableAddBtn()}
        className="h-6 mt-6 px-0 self-start"
      >
        <Icon icon={IconAddOutlined} />
        {t(`form:feature-flags.add-variation`)}
      </Button>
    </>
  );
};

export default Variations;
