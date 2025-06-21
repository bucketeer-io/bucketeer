import { useCallback, useMemo } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import {
  AutoOpsRule,
  Feature,
  OperationStatus,
  OpsEventRateClause,
  Rollout,
  StrategyType
} from '@types';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import ReactCodeEditor from 'components/code-editor';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import { Tooltip } from 'components/tooltip';
import { VariationForm } from '../form-schema';

const VariationLabel = ({
  index,
  className,
  specificColor
}: {
  index: number;
  className?: string;
  specificColor?: string;
}) => (
  <div className={cn('flex items-center gap-x-2 text-gray-600', className)}>
    <FlagVariationPolygon index={index} specificColor={specificColor} />
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
  rollouts,
  isRunningExperiment,
  eventRateOperations,
  editable
}: {
  feature: Feature;
  rollouts: Rollout[];
  isRunningExperiment?: boolean;
  eventRateOperations: AutoOpsRule[];
  editable: boolean;
}) => {
  const { t } = useTranslation(['common', 'form', 'table']);

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
    () => 'flex flex-col flex-1 py-0 h-full self-stretch',
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

  const targetVariationIds = useMemo(() => {
    if (feature?.targets?.length) {
      const ids = feature.targets
        .filter(target => target.users?.length > 0)
        ?.map(item => item.variation);
      return [...new Set(ids)];
    }
    return [];
  }, [feature]);

  const prerequisiteVariationIds = useMemo(() => {
    if (feature?.prerequisites?.length) {
      const ids = feature.prerequisites.map(pre => pre.variationId);
      return [...new Set(ids)];
    }
    return [];
  }, [feature]);

  const rolloutVariationIds = useMemo(
    () => rollouts?.map(item => item.clause.variationId),
    [rollouts]
  );

  const eventRateVariationIds = useMemo(
    () =>
      eventRateOperations
        ?.flatMap(item => item.clauses)
        ?.map(item => (item?.clause as OpsEventRateClause)?.variationId),
    [eventRateOperations]
  );

  const isDisableRemoveBtn = useCallback(
    (variationId: string) => {
      return (
        !editable ||
        isBoolean ||
        isRunningExperiment ||
        fields.length <= 2 ||
        [
          ...new Set([
            offVariation,
            ...onVariationIds,
            ...ruleVariationIds,
            ...targetVariationIds,
            ...prerequisiteVariationIds,
            ...rolloutVariationIds,
            ...eventRateVariationIds
          ])
        ].includes(variationId)
      );
    },
    [
      isBoolean,
      onVariationIds,
      ruleVariationIds,
      offVariation,
      fields,
      targetVariationIds,
      prerequisiteVariationIds,
      isRunningExperiment,
      rolloutVariationIds,
      eventRateVariationIds,
      editable
    ]
  );
  const isProgressiveRolloutsRunningWaiting = (status: OperationStatus) =>
    ['RUNNING', 'WAITING'].includes(status);

  const isDisableAddBtn = useCallback(
    () =>
      !editable ||
      isBoolean ||
      isRunningExperiment ||
      rollouts.filter(item => isProgressiveRolloutsRunningWaiting(item.status))
        ?.length > 0,
    [isBoolean, rollouts, isRunningExperiment, editable]
  );

  const onAddVariation = () => {
    append({
      id: uuid(),
      value: isJSON ? '{}' : '',
      name: '',
      description: ''
    });
  };

  const getTooltipContent = useCallback(
    (variationId: string) => {
      if (onVariationIds.includes(variationId))
        return t('table:feature-flags.default-variation-disabled-delete');
      if (offVariation === variationId)
        return t('table:feature-flags.off-variation-disabled-delete');
      if (
        [
          ...ruleVariationIds,
          ...targetVariationIds,
          ...prerequisiteVariationIds
        ].includes(variationId)
      )
        return t('table:feature-flags.in-used-variation-disabled-delete');
      return '';
    },
    [
      offVariation,
      onVariationIds,
      ruleVariationIds,
      targetVariationIds,
      prerequisiteVariationIds
    ]
  );

  return (
    <div className="flex flex-col w-full gap-y-6">
      {fields.map((variation, variationIndex) => (
        <div key={variation.variationField} className="flex w-full gap-x-2">
          <div className="flex flex-col w-full gap-y-3">
            <VariationLabel index={variationIndex} />
            <div className="flex flex-col w-full gap-y-5">
              <div className="flex items-end w-full gap-x-2">
                {!isJSON && (
                  <Form.Field
                    control={control}
                    name={`variations.${variationIndex}.value`}
                    render={({ field }) => {
                      return (
                        <Form.Item className={cn(formItemClassName)}>
                          <Form.Label required>
                            {t('form:feature-flags.value')}
                          </Form.Label>
                          <Form.Control>
                            <Input
                              {...field}
                              disabled={
                                isBoolean || isRunningExperiment || !editable
                              }
                              placeholder={t('form:feature-flags.value')}
                              className="px-3"
                            />
                          </Form.Control>
                          <Form.Message />
                        </Form.Item>
                      );
                    }}
                  />
                )}

                <Form.Field
                  control={control}
                  name={`variations.${variationIndex}.name`}
                  render={({ field }) => (
                    <Form.Item className={cn(formItemClassName)}>
                      <Form.Label required>{t('name')}</Form.Label>
                      <Form.Control>
                        <Input
                          {...field}
                          placeholder={t('name')}
                          disabled={isRunningExperiment || !editable}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
                <Form.Field
                  control={control}
                  name={`variations.${variationIndex}.description`}
                  render={({ field }) => (
                    <Form.Item className={cn(formItemClassName)}>
                      <Form.Label>{t('form:description')}</Form.Label>
                      <Form.Control>
                        <Input
                          {...field}
                          placeholder={t('form:description')}
                          disabled={isRunningExperiment || !editable}
                        />
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  )}
                />
              </div>
              {isJSON && (
                <Form.Field
                  control={control}
                  name={`variations.${variationIndex}.value`}
                  render={({ field }) => {
                    return (
                      <Form.Item className={cn(formItemClassName)}>
                        <Form.Label required>
                          {t('form:feature-flags.value')}
                        </Form.Label>
                        <Form.Control>
                          <ReactCodeEditor
                            readOnly={isRunningExperiment || !editable}
                            value={field.value}
                            onChange={field.onChange}
                          />
                        </Form.Control>
                        <Form.Message />
                      </Form.Item>
                    );
                  }}
                />
              )}
            </div>
          </div>
          <Tooltip
            align="end"
            alignOffset={-20}
            content={getTooltipContent(variation.id)}
            trigger={
              <Button
                variant="grey"
                size="icon"
                type="button"
                className="p-0 size-5 self-end mb-4"
                disabled={isDisableRemoveBtn(variation.id)}
                onClick={() => remove(variationIndex)}
              >
                <Icon icon={IconTrash} size="sm" />
              </Button>
            }
            className="max-w-[350px]"
          />
        </div>
      ))}
      <Button
        onClick={onAddVariation}
        variant="text"
        type="button"
        disabled={isDisableAddBtn()}
        className="h-6 px-0 self-start"
      >
        <Icon icon={IconAddOutlined} />
        {t(`form:feature-flags.add-variation`)}
      </Button>
    </div>
  );
};

export default Variations;
