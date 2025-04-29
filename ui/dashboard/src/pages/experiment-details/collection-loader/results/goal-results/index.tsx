import { useCallback, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { featureUpdater } from '@api/features';
import { invalidateExperimentDetails } from '@queries/experiment-details';
import { invalidateExperimentResultDetails } from '@queries/experiment-result';
import { invalidateFeature } from '@queries/feature-details';
import { useQueryClient } from '@tanstack/react-query';
import { useToast, useToggleOpen } from 'hooks';
import { Experiment, Feature, GoalResult } from '@types';
import { getData, getTimeSeries } from 'utils/chart';
import { cn } from 'utils/style';
import { IconChevronDown } from '@icons';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
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
  experiment: Experiment;
  feature?: Feature;
  goalResult: GoalResult;
  goalResultState: GoalResultState;
  environmentId: string;
  isRequireComment: boolean;
  onChangeResultState: (tab?: GoalResultTab, chartType?: ChartDataType) => void;
  handleNarrowGoalResult: (goalId: string) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);

  const conversionRateChartRef = useRef<ChartToggleLegendRef>(null);
  const evaluationChartRef = useRef<ChartToggleLegendRef>(null);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const [conversionRateDataSets, setConversionRateDataSets] = useState<
    DatasetReduceType[]
  >([]);
  const [evaluationDataSets, setEvaluationDataSets] = useState<
    DatasetReduceType[]
  >([]);

  const [isOpenRolloutVariant, onOpenRolloutVariant, onCloseRolloutVariant] =
    useToggleOpen(false);

  const variationValues = useMemo(
    () =>
      (goalResult?.variationResults?.map(vr => {
        const variation = experiment.variations.find(
          item => vr.variationId === item.id
        );
        const { name, value } = variation || {};
        return name || value || '';
      }) as string[]) || [],
    [goalResult, experiment]
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
              type: 'FIXED',
              fixedStrategy: {
                variation: values.variation
              }
            },
            comment: values?.comment
          });
          if (resp) {
            notify({
              message: 'Rollout variant updated successfully.'
            });
            invalidateFeature(queryClient);
            invalidateExperimentDetails(queryClient, {
              environmentId,
              id: experiment.id
            });
            invalidateExperimentResultDetails(queryClient, {
              environmentId,
              experimentId: experiment.id
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
        <p className="text-gray-800 typo-head-bold-small">
          {experiment?.goals?.find(goal => goal.id === goalResult?.goalId)
            ?.name || ''}
        </p>
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
      {goalResult?.summary?.bestVariations?.length > 0 &&
        experiment?.status === 'RUNNING' && (
          <ConfidenceVariants
            bestVariations={goalResult.summary.bestVariations}
            variations={experiment.variations}
            onOpenRolloutVariant={onOpenRolloutVariant}
          />
        )}
      {isOpenRolloutVariant && (
        <RolloutVariantModal
          isOpen={isOpenRolloutVariant}
          variations={experiment.variations}
          defaultStrategy={feature?.defaultStrategy}
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
          <TabsTrigger value="EVALUATION">{t(`evaluation`)}</TabsTrigger>
          <TabsTrigger value="CONVERSION">{t(`conversion-rate`)}</TabsTrigger>
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
                goalResult={goalResult}
                experiment={experiment}
                evaluationDataSets={evaluationDataSets}
                onToggleShowData={label =>
                  evaluationChartRef.current?.toggleLegend(label)
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
                conversionRateDataSets={conversionRateDataSets}
                goalResultState={goalResultState}
                goalResult={goalResult}
                experiment={experiment}
                onToggleShowData={label =>
                  conversionRateChartRef.current?.toggleLegend(label)
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
