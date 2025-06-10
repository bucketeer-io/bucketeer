import { useEffect, useMemo, useRef, useState } from 'react';
import { useQueryEvaluation } from '@queries/evaluation';
import { useLoaderData, useSearch } from '@tanstack/react-router';
import { getCurrentEnvironment, useAuth } from 'auth';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import {
  Route as featureDetailsLayoutRoute,
  FeatureDetailsLoaderData
} from 'routes/_default-layout/$env/features/$featureId/_feature-details-layout';
import {
  EvaluationSearch,
  Route as EvaluationsRoute
} from 'routes/_default-layout/$env/features/$featureId/_feature-details-layout/evaluations/';
import { EvaluationTimeRange } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
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
import { EvaluationTab, TimeRangeOption } from './types';

const EvaluationPage = () => {
  const { t } = useTranslation(['common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const evaluationChartRef = useRef<ChartToggleLegendRef>(null);
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const [dataSets, setDataSets] = useState<DatasetReduceType[]>([]);

  const loaderData: FeatureDetailsLoaderData = useLoaderData({
    from: featureDetailsLayoutRoute.id
  });

  const feature = loaderData?.feature;

  const searchFilters: EvaluationSearch = useSearch({
    from: EvaluationsRoute.id
  });

  const defaultFilters = {
    tab: EvaluationTab.EVENT_COUNT,
    period: EvaluationTimeRange.THIRTY_DAYS,
    ...searchFilters
  } as EvaluationSearch;

  const [filters, setFilters] =
    usePartialState<EvaluationSearch>(defaultFilters);

  const { data: evaluationCollection, isLoading } = useQueryEvaluation({
    params: {
      environmentId: currentEnvironment.id,
      featureId: feature?.id,
      timeRange: filters.period!
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
    () =>
      timeRangeOptions.find(item => item.value === filters.period)?.label || '',
    [timeRangeOptions, filters]
  );

  const countData = useMemo(
    () =>
      (filters.tab === EvaluationTab.EVENT_COUNT
        ? evaluationCollection?.eventCounts
        : evaluationCollection?.userCounts) || [],
    [filters, evaluationCollection]
  );

  const variationValues: Option[] = useMemo(
    () =>
      countData?.map(item => ({
        value: item.variationId,
        label:
          feature.variations?.find(v => v.id === item.variationId)?.value ||
          (item.variationId === 'default' ? 'default value' : ''),
        variationType: feature.variationType
      })) || [],
    [countData, feature]
  );

  const timeseries = useMemo(
    () => countData[0]?.timeseries?.timestamps || [],
    [countData]
  );

  const chartData = useMemo(
    () =>
      countData.map(vt => {
        return vt.timeseries?.values?.map((v: number) => Math.round(v));
      }) || [],
    [countData]
  );

  const onChangeFilters = (values: Partial<EvaluationSearch>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      onChangeFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content className="p-6 pt-0 gap-y-6 min-w-[900px]">
      <FilterBar
        isLoading={isLoading}
        timeRangeOptions={timeRangeOptions}
        timeRangeLabel={timeRangeLabel}
        onChangeTimeRange={range =>
          onChangeFilters({
            period: range
          })
        }
      />
      <Card className="h-full">
        <Tabs
          className="flex-1 flex h-full flex-col"
          value={filters.tab}
          onValueChange={value =>
            onChangeFilters({
              tab: value as EvaluationTab
            })
          }
        >
          <TabsList>
            <TabsTrigger value={EvaluationTab.EVENT_COUNT}>
              {t(`table:evaluation.event-count`)}
            </TabsTrigger>
            <TabsTrigger value={EvaluationTab.USER_COUNT}>
              {t(`table:evaluation.user-count`)}
            </TabsTrigger>
          </TabsList>

          <TabsContent value={filters.tab || ''}>
            {isLoading ? (
              <PageLayout.LoadingState />
            ) : (
              <>
                <EvaluationChart
                  ref={evaluationChartRef}
                  data={chartData}
                  variationValues={variationValues}
                  timeseries={timeseries}
                  unit={
                    filters.period === EvaluationTimeRange.TWENTY_FOUR_HOURS
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
