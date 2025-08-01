import { useCallback, useMemo } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import flatmap from 'lodash/flatmap';
import uniqBy from 'lodash/uniqBy';
import { v4 as uuid } from 'uuid';
import { Feature, OperationStatus, Rollout, StrategyType } from '@types';
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
  editable,
  features
}: {
  feature: Feature;
  rollouts: Rollout[];
  isRunningExperiment?: boolean;
  editable: boolean;
  features: Feature[];
}) => {
  const { t } = useTranslation(['common', 'form', 'table']);

  const { control } = useFormContext<VariationForm>();

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'variations',
    keyName: 'variationField'
  });

  const isBoolean = useMemo(
    () => feature.variationType === 'BOOLEAN',
    [feature]
  );
  const isJSON = useMemo(() => feature.variationType === 'JSON', [feature]);

  const formItemClassName = useMemo(
    () => 'flex flex-col flex-1 py-0 h-full self-stretch',
    []
  );

  const prerequisiteVariationIds = useMemo(() => {
    return uniqBy(
      flatmap(features.map(item => item.prerequisites)),
      'variationId'
    ).map(item => item.variationId);
  }, [features]);

  const ruleVariationIds = useMemo(() => {
    const featureRules = flatmap(features.map(item => item.rules));
    const variationIds: string[] = [];

    featureRules.forEach(rule => {
      const { strategy } = rule;
      if (strategy.type === StrategyType.FIXED) {
        variationIds.push(strategy.fixedStrategy.variation);
      } else if (strategy.type === StrategyType.ROLLOUT) {
        strategy.rolloutStrategy.variations.filter(item => {
          if (item.weight > 0) variationIds.push(item.variation);
        });
      }
    });
    return [...new Set(variationIds)];
  }, [feature]);

  const isDisableRemoveBtn = useCallback(
    (variationId: string) => {
      return (
        !editable ||
        isBoolean ||
        isRunningExperiment ||
        fields.length <= 2 ||
        [...prerequisiteVariationIds, ...ruleVariationIds].includes(variationId)
      );
    },
    [
      editable,
      isBoolean,
      fields,
      isRunningExperiment,
      ruleVariationIds,
      prerequisiteVariationIds
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
            content={t('table:feature-flags.default-variation-disabled-delete')}
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
