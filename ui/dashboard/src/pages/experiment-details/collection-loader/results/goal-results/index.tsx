import { useCallback, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { featureUpdater } from '@api/features';
import { useToast, useToggleOpen } from 'hooks';
import { Experiment, Feature, GoalResult, StrategyType } from '@types';
import { getData, getTimeSeries } from 'utils/chart';
import { cn } from 'utils/style';
import { IconChevronDown, IconInfo, IconOutperformed } from '@icons';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import { Tooltip } from 'components/tooltip';
import { ChartDataType, GoalResultState, GoalResultTab } from '..';
import ChartDataTypeDropdown from '../chart-data-type-dropdown';
import ConfidenceVariants from './confidence-variants';
import ConversionRateChart from './conversion-rate-chart';
import ConversionRateTable from './conversion-rate-table';
import EvaluationTable from './evaluation-table';
import RolloutVariantModal, { RolloutVariant } from './rollout-variant-modal';
import {
  ChartToggleLegendRef,
  DatasetReduceType
} from './timeseries-area-line-chart';
import TimeSeriesLineChart from './timeseries-line-chart';

const GoalResultItem = ({
  isNarrow,
  isPrimary,
  experiment,
  feature,
  goalResult,
  goalResultState,
  environmentId,
  isRequireComment,
  onChangeResultState,
  handleNarrowGoalResult
}: {
  isNarrow: boolean;
  isPrimary: boolean;
  experiment: Experiment;
  feature?: Feature;
  goalResult: GoalResult;
  goalResultState: GoalResultState;
  environmentId: string;
  isRequireComment: boolean;
  onChangeResultState: (tab?: GoalResultTab, chartType?: ChartDataType) => void;
  handleNarrowGoalResult: (goalId: string) => void;
}) => {
  const { t } = useTranslation(['common', 'form', 'message', 'table']);

  const conversionRateChartRef = useRef<ChartToggleLegendRef>(null);
  const evaluationChartRef = useRef<ChartToggleLegendRef>(null);
  const { notify, errorNotify } = useToast();

  const [conversionRateDataSets, setConversionRateDataSets] = useState<
    DatasetReduceType[]
  >([]);
  const [evaluationDataSets, setEvaluationDataSets] = useState<
    DatasetReduceType[]
  >([]);

  const [isOpenRolloutVariant, onOpenRolloutVariant, onCloseRolloutVariant] =
    useToggleOpen(false);

  // Pick the best-variations list and the safe-to-stop flag that match the
  // currently selected chart type (driven by chartType, not the tab — value
  // charts are reachable from both tabs). The value charts (`value-user`,
  // `value-total`) use the per-user value posterior; all other chart types
  // use the conversion-rate list.
  const isValueChart = useMemo(
    () =>
      goalResultState?.chartType === 'value-user' ||
      goalResultState?.chartType === 'value-total',
    [goalResultState]
  );

  const activeBestVariations = useMemo(
    () =>
      (isValueChart
        ? goalResult?.summary?.bestVariationsValue
        : goalResult?.summary?.bestVariations) ?? [],
    [goalResult, isValueChart]
  );

  // safeToStop reflects the sequential Bayes Factor verdict for the active
  // metric. When false the evidence threshold has not been met yet and the
  // rollout CTA should be disabled with an explanatory tooltip.
  const safeToStop = useMemo(
    () =>
      isValueChart
        ? (goalResult?.summary?.valueSafeToStop ?? false)
        : (goalResult?.summary?.cvrSafeToStop ?? false),
    [goalResult, isValueChart]
  );

  const variationValues = useMemo(
    () =>
      goalResult?.variationResults?.map(vr => {
        const variation = experiment.variations.find(
          item => vr.variationId === item.id
        );
        const { name, value, id } = variation || {};
        return {
          label: name || value || '',
          value: id || '',
          variationType: feature?.variationType
        };
      }) || [],
    [goalResult, experiment, feature]
  );

  const onSubmitRolloutVariation = useCallback(
    async (values: RolloutVariant) => {
      try {
        if (values.variation && feature) {
          const resp = await featureUpdater({
            id: feature.id,
            environmentId,
            defaultStrategy: {
              ...feature.defaultStrategy,
              type: StrategyType.FIXED,
              fixedStrategy: {
                variation: values.variation
              }
            },
            comment: values?.comment
          });
          if (resp) {
            notify({
              message: t('message:collection-action-success', {
                collection: t('rollout-variant'),
                action: t('updated')
              })
            });
            onCloseRolloutVariant();
          }
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [feature, environmentId]
  );

  return (
    <div
      className={cn(
        'flex flex-col w-full min-w-fit gap-y-6 p-5 shadow-card rounded-lg bg-white',
        {
          'h-16 overflow-hidden': isNarrow
        }
      )}
    >
      <div className="flex items-center justify-between w-full">
        <div className="flex items-center gap-x-3 min-w-0">
          <p className="text-gray-800 typo-head-bold-small truncate">
            {experiment?.goals?.find(goal => goal.id === goalResult?.goalId)
              ?.name || ''}
          </p>
          {isPrimary && (
            <Tooltip
              content={t('table:results.primary-badge-tooltip')}
              trigger={
                <span className="flex items-center gap-x-1 shrink-0 px-2 py-0.5 rounded bg-primary-100/30 typo-para-small font-medium text-primary-500 whitespace-nowrap">
                  <Icon
                    icon={IconOutperformed}
                    size="xxs"
                    color="primary-500"
                  />
                  {t('table:results.primary-badge')}
                  <Icon icon={IconInfo} size="xxs" color="gray-500" />
                </span>
              }
              className="max-w-xs"
            />
          )}
        </div>
        <div
          className={cn(
            'flex-center cursor-pointer rounded hover:bg-gray-200 transition-all duration-300',
            {
              'rotate-180': isNarrow
            }
          )}
          onClick={() => handleNarrowGoalResult(goalResult.goalId)}
        >
          <Icon
            icon={IconChevronDown}
            size={'md'}
            color="gray-500"
            className="rotate-180"
          />
        </div>
      </div>
      {isPrimary
        ? activeBestVariations.length > 0 && (
            <ConfidenceVariants
              bestVariations={activeBestVariations}
              variations={experiment.variations}
              safeToStop={safeToStop}
              onOpenRolloutVariant={onOpenRolloutVariant}
            />
          )
        : !isNarrow && (
            <div className="flex items-center w-full gap-x-2 rounded-lg border-l-4 border-gray-300 bg-gray-100 px-4 py-2">
              <Icon icon={IconInfo} size="xxs" color="gray-600" />
              <p className="typo-para-small text-gray-600">
                {t('table:results.secondary-note')}
              </p>
            </div>
          )}
      {isPrimary && isOpenRolloutVariant && (
        <RolloutVariantModal
          isOpen={isOpenRolloutVariant}
          variations={experiment.variations}
          defaultStrategy={feature?.defaultStrategy}
          bestVariations={activeBestVariations}
          isRequireComment={isRequireComment}
          onClose={onCloseRolloutVariant}
          onSubmit={onSubmitRolloutVariation}
        />
      )}
      <Tabs
        className="flex-1 flex h-full flex-col"
        value={goalResultState?.tab}
        onValueChange={tab =>
          onChangeResultState(
            tab as GoalResultTab,
            (tab as GoalResultTab) === 'EVALUATION'
              ? 'evaluation-user'
              : 'conversion-rate'
          )
        }
      >
        <TabsList>
          <TabsTrigger value="CONVERSION">{t(`conversion-rate`)}</TabsTrigger>
          <TabsTrigger value="EVALUATION">{t(`evaluation`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={goalResultState?.tab} className="mt-6">
          {goalResultState?.tab === 'EVALUATION' && (
            <div className="flex flex-col gap-y-6">
              <ChartDataTypeDropdown
                tab={goalResultState?.tab}
                chartType={goalResultState?.chartType}
                onSelectOption={value =>
                  onChangeResultState(undefined, value as ChartDataType)
                }
              />
              <TimeSeriesLineChart
                ref={evaluationChartRef}
                timeseries={getTimeSeries(
                  goalResult?.variationResults,
                  goalResultState?.chartType,
                  goalResultState?.tab
                )}
                dataLabels={variationValues}
                data={getData(
                  goalResult?.variationResults,
                  goalResultState?.chartType
                )}
                setDataSets={setEvaluationDataSets}
              />
              <EvaluationTable
                variationType={feature?.variationType}
                goalResult={goalResult}
                experiment={experiment}
                evaluationDataSets={evaluationDataSets}
                onToggleShowData={variationId =>
                  evaluationChartRef.current?.toggleLegend(variationId)
                }
              />
            </div>
          )}
          {goalResultState?.tab === 'CONVERSION' && (
            <div className="flex flex-col gap-y-6">
              <ChartDataTypeDropdown
                tab={goalResultState?.tab}
                chartType={goalResultState?.chartType}
                onSelectOption={value =>
                  onChangeResultState(undefined, value as ChartDataType)
                }
              />
              <ConversionRateChart
                ref={conversionRateChartRef}
                variationValues={variationValues}
                goalResult={goalResult}
                goalResultState={goalResultState}
                setConversionRateDataSets={setConversionRateDataSets}
              />
              <ConversionRateTable
                variationType={feature?.variationType}
                conversionRateDataSets={conversionRateDataSets}
                goalResultState={goalResultState}
                goalResult={goalResult}
                experiment={experiment}
                onToggleShowData={variationId =>
                  conversionRateChartRef.current?.toggleLegend(variationId)
                }
              />
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default GoalResultItem;
