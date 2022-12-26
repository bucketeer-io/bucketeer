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

interface FeatureEvaluationProps {
  featureId: string;
}

interface Option {
  readonly value: Type;
  readonly label: string;
}

const typeOptions: Array<Option> = [
  { value: 'userCount', label: 'User Count' },
  { value: 'eventCount', label: 'Event Count' },
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
    const variationMap = new Map<string, Variation.AsObject>();
    feature.variationsList.forEach((v) => {
      variationMap.set(v.id, v);
    });
    const variationTSs = type == 'userCount' ? userCounts : eventCounts;
    const variationValues = variationTSs.map((vt) => {
      return variationMap.get(vt.variationId)?.value;
    });
    const timeseries = variationTSs[0]?.timeseries?.timestampsList;
    const data = variationTSs.map((vt) => {
      return vt.timeseries?.valuesList?.map((v: number) => Math.round(v));
    });

    const handleChange = (o: Option) => {
      setType(o.value);
    };

    if (!timeseries) {
      return <p>No data</p>;
    }
    return (
      <div className="p-10 bg-gray-100">
        <div className="bg-white rounded-md p-3 border ">
          <Select
            options={typeOptions}
            className={classNames('flex-none w-[200px]')}
            value={typeOptions.find((o) => o.value === type)}
            onChange={handleChange}
          />
          <TimeseriesStackedLineChart
            label={''}
            dataLabels={variationValues}
            timeseries={timeseries}
            data={data}
          />
        </div>
      </div>
    );
  }
);
