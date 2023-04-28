import { COLORS } from '@/constants/colorPattern';
import { SerializedError } from '@reduxjs/toolkit';
import { FC, memo, useState } from 'react';
import { shallowEqual, useSelector } from 'react-redux';

import { Select } from '../../components/Select';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { VariationTimeseries } from '../../proto/eventcounter/timeseries_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { Variation } from '../../proto/feature/variation_pb';
import { classNames } from '../../utils/css';
import { TimeseriesStackedLineChart } from '../TimeseriesStackedLineChart';

type Type = 'userCount' | 'eventCount';
type FilterType = 'monthly' | 'last14Days' | 'last7Days' | 'last24Hours';

interface FeatureEvaluationProps {
  featureId: string;
}

interface Option {
  readonly value: Type;
  readonly label: string;
}

interface FilterOption {
  readonly value: FilterType;
  readonly label: string;
}

const typeOptions: Array<Option> = [
  { value: 'userCount', label: 'User Count' },
  { value: 'eventCount', label: 'Event Count' },
];
const filterOptions: Array<FilterOption> = [
  { value: 'monthly', label: 'Monthly' },
  { value: 'last14Days', label: 'Last 14 days' },
  { value: 'last7Days', label: 'Last 7 days' },
  { value: 'last24Hours', label: 'Last 24 hours' },
];

export const FeatureEvaluation: FC<FeatureEvaluationProps> = memo(
  ({ featureId }) => {
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
    const [type, setType] = useState<Type>('userCount');
    const [filterType, setFilterType] = useState<FilterType>('monthly');
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

    const variationTSs = type == 'userCount' ? userCounts : eventCounts;
    console.log('variationTSs', variationTSs);
    const variationValues = variationTSs.map((vt) => {
      return variationMap.get(vt.variationId)?.value;
    });
    const timeseries = variationTSs[0]?.timeseries?.timestampsList;
    const data = variationTSs.map((vt, i) => {
      // return vt.timeseries?.valuesList?.map((v: number, i2) =>
      //   Math.round(i2 * 5)
      // );
      return vt.timeseries?.valuesList?.map((v: number) => Math.round(v));
    });

    const handleChange = (o: Option) => {
      setType(o.value);
    };

    const handleFilterChange = (o: FilterOption) => {
      setFilterType(o.value);
    };

    if (!timeseries) {
      return <p>No data</p>;
    }

    const dataLabels = variationTSs.map((vt, i) => {
      const { name, value } = variationMap.get(vt.variationId);
      const variation = name ? `${name} - ${value}` : value;
      return {
        variation:
          variation.length > 50
            ? `${variation.substring(0, 50)}...`
            : variation,
        backgroundColor: COLORS[i % COLORS.length],
      };
    });

    return (
      <div className="p-10 bg-gray-100">
        <div className="bg-white rounded-md p-3 border">
          <div className="flex justify-end">
            <Select
              options={filterOptions}
              className={classNames('flex-none w-[200px]')}
              value={filterOptions.find((o) => o.value === filterType)}
              onChange={handleFilterChange}
            />
          </div>
          <TimeseriesStackedLineChart
            label={''}
            dataLabels={variationValues}
            timeseries={timeseries}
            data={data}
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
                        <td className="px-3 py-3.5 text-left text-sm font-semibold text-gray-800">
                          Total evaluations
                        </td>
                        <td className="]">
                          <Select
                            options={typeOptions}
                            className={classNames('mt-1 w-[260px]')}
                            value={typeOptions.find((o) => o.value === type)}
                            onChange={handleChange}
                          />
                        </td>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200 bg-white">
                      {dataLabels.map(
                        ({ variation, backgroundColor }, index) => (
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
                            <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 w-1/4">
                              -
                            </td>
                            <td className="w-1/4" />
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
