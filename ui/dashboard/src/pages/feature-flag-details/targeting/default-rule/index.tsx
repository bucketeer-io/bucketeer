import { useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { IconUndoOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import {
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature, Rollout, StrategyType } from '@types';
import { IconInfo, IconInfoFilled, IconPercentage } from '@icons';
import Card from 'pages/feature-flag-details/elements/card';
import Form from 'components/form';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import VariationLabel from 'elements/variation-label';
import { DefaultRuleSchema, TargetingSchema } from '../form-schema';
import Strategy from '../segment-rule/strategy';
import { VariationOption } from '../segment-rule/variation';
import { DiscardChangesType } from '../types';

const DefaultRule = ({
  editable,
  urlCode,
  feature,
  waitingRunningRollouts,
  handleDiscardChanges
}: {
  editable: boolean;
  urlCode: string;
  feature: Feature;
  waitingRunningRollouts: Rollout[];
  handleDiscardChanges: (type: DiscardChangesType) => void;
}) => {
  const { t } = useTranslation(['form']);

  const { control, watch } = useFormContext<TargetingSchema>();

  const commonName = 'defaultRule';
  const defaultRule = watch(commonName);
  const rolloutStrategy = watch('defaultRule.rolloutStrategy');

  const variationOptions: VariationOption[] = useMemo(() => {
    const variations = feature.variations.map((item, index) => ({
      label: <VariationLabel label={item.name || item.value} index={index} />,
      value: item.id,
      type: StrategyType.FIXED,
      variationValue: item.value
    }));
    return [
      ...variations,
      {
        label: t('form:manual-percentage'),
        value: StrategyType.MANUAL,
        type: StrategyType.MANUAL,
        icon: IconPercentage
      }
    ];
  }, [feature]);

  const percentageValueCount = useMemo(
    () =>
      rolloutStrategy?.variations?.filter(item => item.weight > 0)?.length || 0,
    [rolloutStrategy]
  );

  const handleSelectStrategy = useCallback(
    (item: VariationOption, onChange: (item: DefaultRuleSchema) => void) => {
      const { type, value } = item;
      const isFixed = type === StrategyType.FIXED;
      onChange({
        ...defaultRule,
        type: type as StrategyType,
        currentOption: value,
        fixedStrategy: {
          variation: isFixed ? value : ''
        }
      });
    },
    [feature, variationOptions, commonName]
  );

  return (
    <Card>
      {waitingRunningRollouts.length > 0 && (
        <div className="flex items-center gap-x-3 p-4 rounded bg-accent-blue-50 border-l-4 border-accent-blue-500 text-accent-blue-500 typo-para-medium">
          <Icon icon={IconInfoFilled} color="accent-blue-500" size="sm" />
          <div className="flex items-center [&>a]:ml-1">
            <Trans
              i18nKey={'form:targeting.rollout-running-message'}
              components={{
                comp: (
                  <Link
                    className="text-primary-500 underline"
                    to={`/${urlCode}${PAGE_PATH_FEATURES}/${feature.id}${PAGE_PATH_FEATURE_AUTOOPS}`}
                  />
                )
              }}
            />
          </div>
        </div>
      )}
      <div className="flex justify-between">
        <div className="flex flex-col w-full gap-y-3">
          <div className="flex items-center gap-x-2">
            <p className="typo-para-medium text-gray-700">
              {t('targeting.default-rule')}
            </p>
            <Tooltip
              align="start"
              alignOffset={-68}
              content={t('form:targeting.tooltip.default-rule')}
              trigger={
                <div className="flex-center size-fit">
                  <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
                </div>
              }
              className="max-w-[400px]"
            />
          </div>
          <p className="typo-para-small text-gray-500">
            {t('targeting.default-rule-desc')}
          </p>
        </div>
        <div
          className="flex-center h-8 w-8 px-2 rounded-md cursor-pointer group border border-gray-300 hover:border-gray-800"
          onClick={() => {
            handleDiscardChanges(DiscardChangesType.DEFAULT);
          }}
        >
          <Icon
            icon={IconUndoOutlined}
            size={'sm'}
            className="flex-center text-gray-500 group-hover:text-gray-700"
          />
        </div>
      </div>
      <Form.Field
        control={control}
        name="defaultRule"
        render={({ field }) => {
          return (
            <Form.Item className="flex flex-col py-0 gap-y-6">
              <Strategy
                rootName={commonName}
                strategyName={'rolloutStrategy'}
                variationOptions={variationOptions}
                percentageValueCount={percentageValueCount}
                handleSelectStrategy={item =>
                  handleSelectStrategy(item, field.onChange)
                }
                isDisabled={waitingRunningRollouts.length > 0 || !editable}
              />
            </Form.Item>
          );
        }}
      />
    </Card>
  );
};

export default DefaultRule;
