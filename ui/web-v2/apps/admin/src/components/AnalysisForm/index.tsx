import { FC, memo, useCallback, useEffect } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { ANALYSIS_USER_METADATA_REGEX } from '../../constants/analysis';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectAll as selectAllFeatures,
  listFeatures,
} from '../../modules/features';
import { listGoals, selectAll as selectAllGoals } from '../../modules/goals';
import { useCurrentEnvironment } from '../../modules/me';
import {
  listUserMetadata,
  selectAll as selectAllUserMetadata,
} from '../../modules/userMetadata';
import { Goal } from '../../proto/experiment/goal_pb';
import { ListGoalsRequest } from '../../proto/experiment/service_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { Reason } from '../../proto/feature/reason_pb';
import { ListFeaturesRequest } from '../../proto/feature/service_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';
import { DatetimePicker } from '../DatetimePicker';
import { DetailSkeleton } from '../DetailSkeleton';
import { Option, Select } from '../Select';

function getReasonKey(val: number): string {
  const keys = [...Object.keys(Reason.Type)];
  const values = [...Object.values(Reason.Type)];
  return keys[values.findIndex((v) => v == val)];
}

export const reasonOptions: Option[] = [
  {
    value: getReasonKey(Reason.Type.TARGET),
    label: intl.formatMessage(messages.reason.target),
  },
  {
    value: getReasonKey(Reason.Type.RULE),
    label: intl.formatMessage(messages.reason.rule),
  },
  {
    value: getReasonKey(Reason.Type.OFF_VARIATION),
    label: intl.formatMessage(messages.reason.offVariation),
  },
  {
    value: getReasonKey(Reason.Type.CLIENT),
    label: intl.formatMessage(messages.reason.client),
  },
];

export interface AnalysisFormProps {
  onSubmit: () => void;
}

export const AnalysisForm: FC<AnalysisFormProps> = memo(({ onSubmit }) => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const features = useSelector<AppState, Feature.AsObject[]>(
    (state) => selectAllFeatures(state.features),
    shallowEqual
  );
  const isFeatureLoading = useSelector<AppState, boolean>(
    (state) => state.features.loading,
    shallowEqual
  );
  const goals = useSelector<AppState, Goal.AsObject[]>(
    (state) => selectAllGoals(state.goals),
    shallowEqual
  );
  const isGoalLoading = useSelector<AppState, boolean>(
    (state) => state.goals.loading,
    shallowEqual
  );
  const userMetadata = useSelector<AppState, string[]>(
    (state) => selectAllUserMetadata(state.userMetadata),
    shallowEqual
  );
  const isUserMetadataLoading = useSelector<AppState, boolean>(
    (state) => state.userMetadata.loading,
    shallowEqual
  );
  const isLoading = isFeatureLoading || isGoalLoading || isUserMetadataLoading;
  const featureOptions = features.map((feature) => {
    return {
      value: feature.id,
      label: feature.name,
    };
  });
  const goalOptions = goals.map((goal) => {
    return {
      value: goal.id,
      label: goal.name,
    };
  });
  const segmentOptions = userMetadata.map((data) => {
    return {
      value: data,
      label: getSegmentLabel(data),
    };
  });
  segmentOptions.push({
    value: 'tag',
    label: f(messages.tags),
  });
  const methods = useFormContext();
  const {
    control,
    formState: { errors, isSubmitting, isDirty, isValid, dirtyFields },
    setValue,
    watch,
  } = methods;
  const featureId = watch('featureId');

  const handleOnChangeFeature = useCallback(
    (featureId: string) => {
      const feature = features?.find((f) => f.id === featureId);
      setValue('featureVersion', feature ? feature.version : 0);
      !feature && setValue('reason', '');
    },
    [features, setValue]
  );

  useEffect(() => {
    dispatch(
      listGoals({
        environmentNamespace: currentEnvironment.namespace,
        pageSize: 99999,
        cursor: '',
        searchKeyword: null,
        status: null,
        orderBy: ListGoalsRequest.OrderBy.DEFAULT,
        orderDirection: ListGoalsRequest.OrderDirection.ASC,
      })
    );
    dispatch(
      listFeatures({
        environmentNamespace: currentEnvironment.namespace,
        pageSize: 99999,
        cursor: '',
        tags: [],
        searchKeyword: null,
        enabled: null,
        hasExperiment: null,
        maintainerId: null,
        orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
        orderDirection: ListFeaturesRequest.OrderDirection.ASC,
      })
    );
    dispatch(
      listUserMetadata({
        environmentNamespace: currentEnvironment.namespace,
      })
    );
  }, [dispatch]);

  return isLoading ? (
    <div className="p-9 bg-gray-100">
      <DetailSkeleton />
    </div>
  ) : (
    <form className="flex flex-col">
      <div className={classNames('space-y-6 pt-6 flex flex-col mx-5')}>
        <div>
          <div className="flex flex-row space-x-5">
            <div>
              <label className="input-label">
                {f(messages.experiment.startAt)}
              </label>
              <DatetimePicker name={`startAt`} />
            </div>
            <div>
              <label className="input-label">
                {f(messages.experiment.stopAt)}
              </label>
              <DatetimePicker name={`endAt`} />
            </div>
          </div>
          <p className="input-error">
            {errors.startAt?.message && (
              <span role="alert">{errors.startAt?.message}</span>
            )}
          </p>
          <p className="input-error">
            {errors.endAt?.message && (
              <span role="alert">{errors.endAt?.message}</span>
            )}
          </p>
        </div>
        <div>
          <label className="input-label">
            {f(messages.experiment.goalIds)}
          </label>
          <Controller
            name="goalId"
            control={control}
            render={({ field }) => {
              return (
                <Select
                  options={goalOptions}
                  onChange={(o: Option) => {
                    field.onChange(o.value);
                  }}
                />
              );
            }}
          />
          <p className="input-error">
            {errors.goalId?.message && (
              <span role="alert">{errors.goalId?.message}</span>
            )}
          </p>
        </div>
        <div>
          <label className="input-label">
            {f(messages.experiment.feature)}
          </label>
          <Controller
            name="featureId"
            control={control}
            render={({ field }) => {
              return (
                <Select
                  options={featureOptions}
                  onChange={(o: Option) => {
                    field.onChange(o ? o.value : '');
                    handleOnChangeFeature(o?.value.toString());
                  }}
                  value={featureOptions.find((o) => o.value === field.value)}
                  clearable={true}
                />
              );
            }}
          />
          <p className="input-error">
            {errors.featureid?.message && (
              <span role="alert">{errors.featureid?.message}</span>
            )}
          </p>
        </div>
        {featureId && (
          <div>
            <label className="input-label">{f(messages.reason.reason)}</label>
            <Controller
              name="reason"
              control={control}
              render={({ field }) => {
                return (
                  <Select
                    options={reasonOptions}
                    onChange={(o: Option) => {
                      field.onChange(o?.value);
                    }}
                    clearable={true}
                  />
                );
              }}
            />
          </div>
        )}
        <div>
          <label className="input-label">{f(messages.analysis.segment)}</label>
          <Controller
            name="segments"
            control={control}
            render={({ field }) => {
              return (
                <Select
                  isMulti={true}
                  options={segmentOptions}
                  onChange={(ops: Option[]) => {
                    field.onChange(ops.map((o) => o.value));
                  }}
                  clearable={true}
                />
              );
            }}
          />
          <p className="input-error">
            {errors.userMetadata?.message && (
              <span role="alert">{errors.userMetadata?.message}</span>
            )}
          </p>
        </div>
      </div>
      <div className="flex-shrink-0 px-4 py-4 flex justify-end">
        <button
          type="button"
          className="btn-submit"
          disabled={!isDirty || !isValid || isSubmitting}
          onClick={onSubmit}
        >
          {f(messages.button.submit)}
        </button>
      </div>
    </form>
  );
});

const getSegmentLabel = (segment: string): string => {
  const match = segment.match(ANALYSIS_USER_METADATA_REGEX);
  if (match && match?.length > 1) {
    return `${match[1]} (${intl.formatMessage(messages.analysis.clientData)})`;
  }
  return segment;
};
