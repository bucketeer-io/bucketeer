import { Dialog } from '@headlessui/react';
import { FC, memo, useCallback, useEffect, useState } from 'react';
import { Controller, useFormContext, useController } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectAll as selectAllFeatures,
  listFeatures,
} from '../../modules/features';
import { listGoals, selectAll as selectAllGoals } from '../../modules/goals';
import { useCurrentEnvironment } from '../../modules/me';
import { Goal } from '../../proto/experiment/goal_pb';
import { ListGoalsRequest } from '../../proto/experiment/service_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { ListFeaturesRequest } from '../../proto/feature/service_pb';
import { AppDispatch } from '../../store';
import { DatetimePicker } from '../DatetimePicker';
import { DetailSkeleton } from '../DetailSkeleton';
import { Option, Select } from '../Select';
import { OptionFeatureFlag, SelectFeatureFlag } from '../SelectFeatureFlag';

export interface ExperimentAddFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const ExperimentAddForm: FC<ExperimentAddFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();
    const features = useSelector<AppState, Feature.AsObject[]>(
      (state) => selectAllFeatures(state.features),
      shallowEqual
    );
    const isListFeatureLoading = useSelector<AppState, boolean>(
      (state) => state.features.listFeaturesLoading,
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
    const isLoading = isListFeatureLoading || isGoalLoading;
    const featureFlagOptions = features.map((feature) => {
      return {
        value: feature.id,
        label: `${feature.id}(${feature.name})`,
        enabled: feature.enabled,
      };
    });
    const goalOptions = goals.map((goal) => {
      return {
        value: goal.id,
        label: `${goal.id}(${goal.name})`,
      };
    });
    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors, isSubmitting, isDirty, isValid },
      getValues,
      reset,
      watch,
    } = methods;
    const watchBaselineVariation = watch('baselineVariation', null);
    const featureId = getValues('featureId');
    const [baselineVariationOptions, setBaselineVariationOptions] = useState<
      Option[] | null
    >(featureId && createBaselineVariationOptions(features, featureId));

    const handleOnChangeFeature = useCallback(
      (featureId: string) => {
        reset(
          { ...getValues(), baselineVariation: null },
          {
            keepDirty: true,
            keepErrors: true,
            keepIsValid: true,
            keepTouched: true,
          }
        );
      },
      [features, setBaselineVariationOptions]
    );

    useEffect(() => {
      dispatch(
        listGoals({
          environmentNamespace: currentEnvironment.id,
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
          environmentNamespace: currentEnvironment.id,
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
    }, [dispatch]);

    useEffect(() => {
      setBaselineVariationOptions(
        createBaselineVariationOptions(features, featureId)
      );
    }, [features, featureId]);

    return isLoading ? (
      <div className="p-9 bg-gray-100">
        <DetailSkeleton />
      </div>
    ) : (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.experiment.add.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.experiment.add.header.description)}
                </p>
              </div>
            </div>
            <div
              className="
                flex-1
                flex flex-col
                justify-between
              "
            >
              <div
                className="
                  space-y-6 px-5 pt-6 pb-5
                  flex flex-col
                "
              >
                <div className="">
                  <label htmlFor="name">
                    <span className="input-label">{f({ id: 'name' })}</span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('name')}
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="description" className="block">
                    <span className="input-label">
                      {f(messages.description)}
                    </span>
                    <span className="input-label-optional">
                      {' '}
                      {f(messages.input.optional)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <textarea
                      {...register('description')}
                      id="description"
                      name="description"
                      rows={4}
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.description && (
                        <span role="alert">{errors.description.message}</span>
                      )}
                    </p>
                  </div>
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
                        <SelectFeatureFlag
                          options={featureFlagOptions}
                          className="w-full"
                          onChange={(e: OptionFeatureFlag) => {
                            field.onChange(e.value);
                            handleOnChangeFeature(e.value);
                          }}
                          value={featureFlagOptions.find(
                            (o) => o.value === field.value
                          )}
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
                {baselineVariationOptions ? (
                  <div>
                    <label className="input-label">
                      {f(messages.experiment.baselineVariation)}
                    </label>
                    <Controller
                      name="baselineVariation"
                      control={control}
                      render={({ field }) => {
                        return (
                          <Select
                            options={baselineVariationOptions}
                            className="w-full"
                            onChange={(e) => {
                              field.onChange(e.value);
                            }}
                            value={
                              watchBaselineVariation
                                ? baselineVariationOptions?.find(
                                    (o) => o.value == watchBaselineVariation
                                  )
                                : { value: '', label: '' }
                            }
                          />
                        );
                      }}
                    />
                    <p className="input-error">
                      {errors.baselineVariation?.message && (
                        <span role="alert">
                          {errors.baselineVariation?.message}
                        </span>
                      )}
                    </p>
                  </div>
                ) : null}
                <div>
                  <label className="input-label">
                    {f(messages.experiment.goalIds)}
                  </label>
                  <Controller
                    name="goalIds"
                    control={control}
                    render={({ field }) => {
                      return (
                        <Select
                          isMulti
                          options={goalOptions}
                          className="w-full"
                          onChange={(e) => {
                            field.onChange(e.map((o) => o.value));
                          }}
                        />
                      );
                    }}
                  />
                  <p className="input-error">
                    {errors.goalIds?.message && (
                      <span role="alert">{errors.goalIds?.message}</span>
                    )}
                  </p>
                </div>
                <div>
                  <label className="input-label">
                    {f(messages.experiment.startAt)}
                  </label>
                  <DatetimePicker name={`startAt`} />
                  <p className="input-error">
                    {errors.startAt?.message && (
                      <span role="alert">{errors.startAt?.message}</span>
                    )}
                  </p>
                </div>
                <div>
                  <label className="input-label">
                    {f(messages.experiment.stopAt)}
                  </label>
                  <DatetimePicker name={`stopAt`} />
                  <p className="input-error">
                    {errors.stopAt?.message && (
                      <span role="alert">{errors.stopAt?.message}</span>
                    )}
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div className="flex-shrink-0 px-4 py-4 flex justify-end">
            <div className="mr-3">
              <button
                type="button"
                className="btn-cancel"
                disabled={false}
                onClick={onCancel}
              >
                {f(messages.button.cancel)}
              </button>
            </div>
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
      </div>
    );
  }
);

function createBaselineVariationOptions(
  features: Feature.AsObject[],
  featureId: string
): Option[] {
  return features
    .find((f) => f.id === featureId)
    ?.variationsList.map((v) => {
      return {
        value: v.id,
        label: v.value,
      };
    });
}
