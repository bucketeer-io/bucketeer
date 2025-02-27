import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Experiment, GoalResult } from '@types';
import { getData, getTimeSeries } from 'utils/chart';
import { IconChevronDown } from '@icons';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import { ChartDataType, GoalResultState, GoalResultTab } from '..';
import ChartDataTypeDropdown from '../chart-data-type-dropdown';
import ConversionRateChart from './conversion-rate-chart';
import ConversionRateTable from './conversion-rate-table';
import EvaluationTable from './evaluation-table';
import TimeSeriesLineChart from './timeseries-line-chart';

const GoalResultItem = ({
  experiment,
  goalResult,
  goalResultState,
  onChangeResultState
}: {
  experiment: Experiment;
  goalResult: GoalResult;
  goalResultState: GoalResultState;
  onChangeResultState: (tab?: GoalResultTab, chartType?: ChartDataType) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);

  const variationValues = useMemo(
    () =>
      (goalResult?.variationResults?.map(vr => {
        return experiment.variations.find(item => vr.variationId === item.id)
          ?.value;
      }) as string[]) || [],
    [goalResult, experiment]
  );

  return (
    <div className="flex flex-col w-full min-w-fit gap-y-6 p-5 shadow-card rounded-lg bg-white">
      <div className="flex items-center justify-between w-full">
        <p className="text-gray-800 typo-head-bold-small">
          {experiment?.goals?.find(goal => goal.id === goalResult?.goalId)
            ?.name || ''}
        </p>
        <Icon
          icon={IconChevronDown}
          size={'md'}
          color="gray-500"
          className="rotate-180"
        />
      </div>
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
              <EvaluationTable
                goalResult={goalResult}
                experiment={experiment}
              />
              <ChartDataTypeDropdown
                tab={goalResultState?.tab}
                chartType={goalResultState?.chartType}
                onSelectOption={value =>
                  onChangeResultState(undefined, value as ChartDataType)
                }
              />
              <TimeSeriesLineChart
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
              />
            </div>
          )}
          {goalResultState?.tab === 'CONVERSION' && (
            <div className="flex flex-col gap-y-6">
              <ConversionRateTable
                goalResultState={goalResultState}
                goalResult={goalResult}
                experiment={experiment}
              />
              <ChartDataTypeDropdown
                tab={goalResultState?.tab}
                chartType={goalResultState?.chartType}
                onSelectOption={value =>
                  onChangeResultState(undefined, value as ChartDataType)
                }
              />
              <ConversionRateChart
                variationValues={variationValues}
                goalResult={goalResult}
                goalResultState={goalResultState}
              />
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default GoalResultItem;
