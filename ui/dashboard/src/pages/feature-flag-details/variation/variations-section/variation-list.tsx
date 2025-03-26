import { useCallback, useMemo } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconAddOutlined } from 'react-icons-material-design';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { OperationStatus, StrategyType } from '@types';
import { cn } from 'utils/style';
import { IconInfo, IconTrash } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import TextArea from 'components/textarea';
import { VariationProps } from '..';
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

const VariationList = ({ feature }: VariationProps) => {
  const { t } = useTranslation(['common', 'form']);

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: rolloutCollection } = useQueryRollouts({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment?.id,
      featureIds: [feature?.id]
    },
    enabled: !!currentEnvironment?.id && !!feature?.id
  });

  const rollouts = rolloutCollection?.progressiveRollouts || [];

  const { control, watch, trigger } = useFormContext<VariationForm>();
  const { append, remove } = useFieldArray({
    control,
    name: 'variations',
    keyName: 'key'
  });

  const offVariation = watch('offVariation');
  const variations = watch('variations');

  const offVariationValue = useMemo(() => {
    const variation = variations.find(item => item.id === offVariation);
    return variation?.value || variation?.name || '';
  }, [offVariation, variations]);

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
    (variationId: string) =>
      isBoolean ||
      [
        ...new Set([offVariation, ...onVariationIds, ...ruleVariationIds])
      ].includes(variationId),
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

  const onAddVariation = useCallback(() => {
    append({
      id: uuid(),
      value: '',
      name: '',
      description: ''
    });
  }, [variations]);

  return (
    <>
      <Form.Field
        control={control}
        name="variations"
        render={() => (
          <Form.Item className="flex flex-col w-full py-0">
            <Form.Control>
              <>
                {variations.map((v, index) => (
                  <div key={index} className="flex items-end w-full gap-x-2">
                    <Form.Field
                      name={`variations.${index}.value`}
                      control={control}
                      render={({ field }) => (
                        <Form.Item
                          className={cn(formItemClassName, 'w-[20%]', {
                            'w-[10%]': isBoolean,
                            'flex-1 w-full': isJSON
                          })}
                        >
                          {isBoolean && (
                            <VariationLabel index={index} className="mb-3" />
                          )}
                          <Form.Label required={isBoolean}>
                            {isBoolean ? (
                              t(field.value ? 'true' : 'false')
                            ) : (
                              <VariationLabel index={index} />
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
                      )}
                    />
                    <Form.Field
                      name={`variations.${index}.name`}
                      control={control}
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
                      name={`variations.${index}.description`}
                      control={control}
                      render={({ field }) => (
                        <Form.Item
                          className={cn(formItemClassName, 'flex-1', {
                            'flex-1 w-full': isJSON
                          })}
                        >
                          <Form.Label>{t('form:description')}</Form.Label>
                          <Form.Control>
                            <Input
                              {...field}
                              placeholder={t('form:description')}
                            />
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
                      disabled={isDisableRemoveBtn(v.id)}
                      onClick={() => {
                        remove(index);
                        trigger('variations');
                      }}
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
            </Form.Control>
          </Form.Item>
        )}
      />

      <Form.Field
        control={control}
        name={'offVariation'}
        render={({ field }) => (
          <Form.Item className="pt-6 pb-0">
            <Form.Label required className="relative w-fit mb-6">
              {t('form:off-variation')}
              <Icon
                icon={IconInfo}
                size="xs"
                color="gray-500"
                className="absolute -right-6"
              />
            </Form.Label>
            <Form.Control>
              <DropdownMenu>
                <DropdownMenuTrigger
                  label={offVariationValue}
                  isExpand
                  className={isBoolean ? 'capitalize' : ''}
                />
                <DropdownMenuContent align="start">
                  {variations?.map((item, index) => (
                    <DropdownMenuItem
                      {...field}
                      key={index}
                      label={item.value || item.name}
                      value={item.id}
                      className={isBoolean ? 'capitalize' : ''}
                      onSelectOption={value => field.onChange(value)}
                    />
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </Form.Control>
          </Form.Item>
        )}
      />
    </>
  );
};

export default VariationList;
