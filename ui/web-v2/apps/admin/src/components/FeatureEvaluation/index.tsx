import { COLORS } from '@/constants/colorPattern';
import { useCurrentEnvironment } from '@/modules/me';
import { AppDispatch } from '@/store';
import { SerializedError } from '@reduxjs/toolkit';
import { FC, memo, useState } from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { Select } from '../../components/Select';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { getEvaluationTimeseriesCount } from '../../modules/evaluationTimeseriesCount';
import { selectById as selectFeatureById } from '../../modules/features';
import { VariationTimeseries } from '../../proto/eventcounter/timeseries_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { Variation } from '../../proto/feature/variation_pb';
import { classNames } from '../../utils/css';
import { TimeseriesStackedLineChart } from '../TimeseriesStackedLineChart';

export enum TimeRange {
  TWENTY_FOUR_HOURS = 1,
  SEVEN_DAYS = 2,
  FOURTEEN_DAYS = 3,
  LAST_THIRTY_DAYS = 4,
}

enum CountsListType {
  USER_COUNT = 'userCount',
  Event_COUNT = 'eventCount',
}

export interface TimeRangeOption {
  value: string;
  label: string;
  data: string;
}

interface FeatureEvaluationProps {
  featureId: string;
  selectedTimeRange: TimeRangeOption;
  setSelectedTimeRange: React.Dispatch<React.SetStateAction<TimeRangeOption>>;
}

const countsListOptions = [
  { value: CountsListType.Event_COUNT, label: 'Event Count' },
  { value: CountsListType.USER_COUNT, label: 'User Count' },
];

export const timeRangeOptions: TimeRangeOption[] = [
  {
    value: TimeRange.LAST_THIRTY_DAYS.toString(),
    label: intl.formatMessage(messages.feature.evaluation.last30Days),
    data: 'day',
  },
  {
    value: TimeRange.FOURTEEN_DAYS.toString(),
    label: intl.formatMessage(messages.feature.evaluation.last14Days),
    data: 'day',
  },
  {
    value: TimeRange.SEVEN_DAYS.toString(),
    label: intl.formatMessage(messages.feature.evaluation.last7Days),
    data: 'day',
  },
  {
    value: TimeRange.TWENTY_FOUR_HOURS.toString(),
    label: intl.formatMessage(messages.feature.evaluation.last24Hours),
    data: 'hour',
  },
];

export const FeatureEvaluation: FC<FeatureEvaluationProps> = memo(
  ({ featureId, selectedTimeRange, setSelectedTimeRange }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();
    const [userCounts, eventCounts] = useSelector<
      AppState,
      [Array<VariationTimeseries.AsObject>, Array<VariationTimeseries.AsObject>]
    >(
      (state) => [
        state.evaluationTimeseriesCount.userCountsList,
        state.evaluationTimeseriesCount.eventCountsList,
      ],
      shallowEqual
    );
    const [feature, _] = useSelector<
      AppState,
      [Feature.AsObject | undefined, SerializedError | null]
    >(
      (state) => [
        selectFeatureById(state.features, featureId),
        state.features.getFeatureError,
      ],
      shallowEqual
    );
    const [selectedCountsListType, setSelectedCountsListType] = useState(
      countsListOptions[0]
    );

    const variationMap = new Map<string, Variation.AsObject>();
    feature.variationsList.forEach((v) => {
      variationMap.set(v.id, v);
    });

    // The default variation corresponds to the default value used on the client.
    // Because the SDK doesn't know what variation was used, we define it as "default".
    // So the user can see how many times the default value was used.
    const defaultVariation = new Variation();
    defaultVariation.setId('default');
    defaultVariation.setValue('default value');
    variationMap.set(defaultVariation.getId(), defaultVariation.toObject());

    const variationTSs =
      selectedCountsListType.value == CountsListType.USER_COUNT
        ? userCounts
        : eventCounts;

    const variationValues = variationTSs.map((vt) => {
      return variationMap.get(vt.variationId)?.value || '';
    });
    const timeseries = variationTSs[0]?.timeseries?.timestampsList;
    const data = variationTSs.map((vt, i) => {
      return vt.timeseries?.valuesList?.map((v: number) => Math.round(v));
    });

    const handleTimeRange = (o) => {
      setSelectedTimeRange(o);
      dispatch(
        getEvaluationTimeseriesCount({
          featureId: featureId,
          environmentNamespace: currentEnvironment.id,
          timeRange: o.value,
        })
      );
    };

    if (!timeseries) {
      return <p>No data</p>;
    }

    const dataLabels = variationTSs.map((vt, i) => {
      let variation = '';
      if (variationMap.get(vt.variationId)) {
        const { name, value } = variationMap.get(vt.variationId);
        variation = name ? `${name} - ${value}` : value;
      }

      return {
        variation:
          variation.length > 50
            ? `${variation.substring(0, 50)}...`
            : variation,
        backgroundColor: COLORS[i % COLORS.length],
        totalCounts: vt.timeseries.totalCounts,
      };
    });

    return (
      <div className="p-10 bg-gray-100">
        <div className="bg-white rounded-md p-3 border">
          <div className="flex justify-end space-x-4">
            <div className="flex">
              <div
                onClick={() => setSelectedCountsListType(countsListOptions[0])}
                className={`px-4 h-[42px] text-sm flex justify-center items-center rounded-l-lg cursor-pointer border ${
                  selectedCountsListType.value === countsListOptions[0].value
                    ? 'text-pink-400 border-pink-400 bg-pink-50 bg-opacity-50'
                    : 'text-gray-500'
                }`}
              >
                Event Count
              </div>
              <div
                onClick={() => setSelectedCountsListType(countsListOptions[1])}
                className={`px-4 h-[42px] text-sm flex justify-center items-center rounded-r-lg cursor-pointer text-gray-500 border ${
                  selectedCountsListType.value === countsListOptions[1].value
                    ? 'text-pink-400 border-pink-400 bg-pink-50 bg-opacity-50'
                    : 'text-gray-500'
                }`}
              >
                User Count
              </div>
            </div>
            <Select
              options={timeRangeOptions}
              className={classNames('flex-none w-[200px]')}
              value={selectedTimeRange}
              onChange={handleTimeRange}
            />
          </div>
          <TimeseriesStackedLineChart
            label={''}
            dataLabels={variationValues}
            timeseries={timeseries}
            data={data}
            unit={selectedTimeRange.data}
          />
          <div className="mt-8 flow-root">
            <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
              <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
                  <table className="min-w-full divide-y divide-gray-300">
                    <thead className="">
                      <tr>
                        <td className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-800 sm:pl-6">
                          Variation
                        </td>
                        <td className="px-3 py-3.5 text-left text-sm">
                          <span className="font-semibold text-gray-800">
                            Total evaluations
                          </span>
                          <span className="text-gray-600 ml-1">
                            ({selectedTimeRange.label})
                          </span>
                        </td>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200 bg-white">
                      {dataLabels.map(
                        (
                          { variation, backgroundColor, totalCounts },
                          index
                        ) => (
                          <tr key={index}>
                            <td className="p-4 text-sm text-gray-900 flex space-x-2">
                              <div
                                className="w-4 h-4"
                                style={{
                                  backgroundColor,
                                }}
                              />
                              <span className="">{variation}</span>
                            </td>
                            <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                              {Number(totalCounts).toLocaleString()}
                            </td>
                          </tr>
                        )
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
);
