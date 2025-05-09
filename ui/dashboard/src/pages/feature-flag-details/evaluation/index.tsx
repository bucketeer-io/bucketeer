import { useMemo, useRef, useState } from 'react';
import { useQueryEvaluation } from '@queries/evaluation';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { EvaluationTimeRange, Feature } from '@types';
import {
  ChartToggleLegendRef,
  DatasetReduceType
} from 'pages/experiment-details/collection-loader/results/goal-results/timeseries-area-line-chart';
import { Option } from 'components/creatable-select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Card from 'elements/card';
import PageLayout from 'elements/page-layout';
import { EvaluationChart } from './evaluation-chart';
import EvaluationTable from './evaluation-table';
import FilterBar from './filter-bar';

type EvaluationTab = 'EVENT_COUNT' | 'USER_COUNT';

export interface TimeRangeOption {
  label: string;
  value: EvaluationTimeRange;
}

const EvaluationPage = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const evaluationChartRef = useRef<ChartToggleLegendRef>(null);

  const [timeRange, setTimeRange] = useState<EvaluationTimeRange>(
    EvaluationTimeRange.THIRTY_DAYS
  );
  const [currentTab, setCurrentTab] = useState<EvaluationTab>('EVENT_COUNT');
  const [dataSets, setDataSets] = useState<DatasetReduceType[]>([]);

  const { data: evaluationCollection, isLoading } = useQueryEvaluation({
    params: {
      environmentId: currentEnvironment.id,
      featureId: feature.id,
      timeRange
    },
    gcTime: 0
  });

  const timeRangeOptions: TimeRangeOption[] = useMemo(
    () => [
      {
        label: t('table:evaluation.last-30-days'),
        value: EvaluationTimeRange.THIRTY_DAYS
      },
      {
        label: t('table:evaluation.last-14-days'),
        value: EvaluationTimeRange.FOURTEEN_DAYS
      },
      {
        label: t('table:evaluation.last-7-days'),
        value: EvaluationTimeRange.SEVEN_DAYS
      },
      {
        label: t('table:evaluation.last-24h'),
        value: EvaluationTimeRange.TWENTY_FOUR_HOURS
      }
    ],
    []
  );

  const timeRangeLabel = useMemo(
    () => timeRangeOptions.find(item => item.value === timeRange)?.label || '',
    [timeRangeOptions, timeRange]
  );

  const countData = useMemo(
    () =>
      (currentTab === 'EVENT_COUNT'
        ? evaluationCollection?.eventCounts
        : evaluationCollection?.userCounts) || [],
    [currentTab, evaluationCollection]
  );

  const variationValues: Option[] = useMemo(
    () =>
      countData?.map(item => ({
        value: item.variationId,
        label:
          feature.variations.find(v => v.id === item.variationId)?.value ||
          (item.variationId === 'default' ? 'default value' : '')
      })) || [],
    [countData, feature]
  );

  const timeseries = useMemo(
    () => countData[0]?.timeseries?.timestamps || [],
    [countData]
  );

  const data = useMemo(
    () =>
      countData.map(vt => {
        return vt.timeseries?.values?.map((v: number) => Math.round(v));
      }) || [],
    [countData]
  );

  return (
    <PageLayout.Content className="p-6 pt-0 gap-y-6 min-w-[900px]">
      <FilterBar
        isLoading={isLoading}
        timeRangeOptions={timeRangeOptions}
        timeRangeLabel={timeRangeLabel}
        onChangeTimeRange={range => setTimeRange(range)}
      />
      <Card className="h-full">
        <Tabs
          className="flex-1 flex h-full flex-col"
          value={currentTab}
          onValueChange={value => setCurrentTab(value as EvaluationTab)}
        >
          <TabsList>
            <TabsTrigger value="EVENT_COUNT">
              {t(`table:evaluation.event-count`)}
            </TabsTrigger>
            <TabsTrigger value="USER_COUNT">
              {t(`table:evaluation.user-count`)}
            </TabsTrigger>
          </TabsList>

          <TabsContent value={currentTab as EvaluationTab}>
            {isLoading ? (
              <PageLayout.LoadingState />
            ) : (
              <>
                <EvaluationChart
                  ref={evaluationChartRef}
                  data={data}
                  variationValues={variationValues}
                  timeseries={timeseries}
                  unit={
                    timeRange === EvaluationTimeRange.TWENTY_FOUR_HOURS
                      ? 'hour'
                      : 'day'
                  }
                  setDataSets={setDataSets}
                />
              </>
            )}
          </TabsContent>
        </Tabs>
      </Card>
      {!isLoading && (
        <EvaluationTable
          feature={feature}
          dataSets={dataSets}
          timeRangeLabel={timeRangeLabel}
          countData={countData}
          onToggleShowData={variationId =>
            evaluationChartRef.current?.toggleLegend(variationId)
          }
        />
      )}
    </PageLayout.Content>
  );
};

export default EvaluationPage;
