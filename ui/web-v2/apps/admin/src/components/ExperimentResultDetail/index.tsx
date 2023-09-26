import { css, jsx } from '@emotion/react';
import { makeStyles, Theme, useTheme } from '@material-ui/core/styles';
import MUExitToAppIcon from '@material-ui/icons/ExitToApp';
import { SerializedError } from '@reduxjs/toolkit';
import dayjs from 'dayjs';
import { unwrapUndefinable } from 'option-t/lib/Undefinable/unwrap';
import { FC, useEffect, useState } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useSelector, shallowEqual, useDispatch } from 'react-redux';

import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  getExperimentResult,
  selectById as selectExperimentResultById,
} from '../../modules/experimentResult';
import { selectById as selectExperimentById } from '../../modules/experiments';
import { useCurrentEnvironment } from '../../modules/me';
import { ExperimentResult } from '../../proto/eventcounter/experiment_result_pb';
import { Experiment } from '../../proto/experiment/experiment_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';
import { DetailSkeleton } from '../DetailSkeleton';
import { GoalResultDetail } from '../GoalResultDetail';
import { Option, Select } from '../Select';

interface ExperimentResultDetailProps {
  experimentId: string;
}

export const ExperimentResultDetail: FC<ExperimentResultDetailProps> = ({
  experimentId,
}) => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const [experiment, getExperimentError] = useSelector<
    AppState,
    [Experiment.AsObject | undefined, SerializedError | null]
  >(
    (state) => [
      selectExperimentById(state.experiments, experimentId),
      state.experiments.getExperimentError,
    ],
    shallowEqual
  );
  const [experimentResult, getExperimentResultError] = useSelector<
    AppState,
    [ExperimentResult.AsObject | undefined, SerializedError | null]
  >(
    (state) => [
      selectExperimentResultById(state.experimentResults, experimentId),
      state.experimentResults.getExperimentResultError,
    ],
    shallowEqual
  );
  const isExperimentLoading = useSelector<AppState, boolean>(
    (state) => state.experiments.loading,
    shallowEqual
  );
  const isExperimentResultLoading = useSelector<AppState, boolean>(
    (state) => state.experimentResults.loading,
    shallowEqual
  );
  const isLoading = isExperimentLoading || isExperimentResultLoading;

  const [goalId, setGoalId] = useState<string>('');
  const goalOptions = experimentResult?.goalResultsList.map((gr) => {
    return {
      value: gr.goalId,
      label: gr.goalId,
    };
  });

  useEffect(() => {
    if (experimentId && experiment.startAt < Number(Date.now() / 1000)) {
      dispatch(
        getExperimentResult({
          environmentNamespace: currentEnvironment.id,
          experimentId: experimentId,
        })
      );
    }
  }, [experimentId, dispatch, currentEnvironment]);

  useEffect(() => {
    setGoalId(experimentResult?.goalResultsList[0]?.goalId);
  }, [experimentResult]);

  if (isLoading) {
    return (
      <div className="p-9 bg-gray-100">
        <DetailSkeleton />
      </div>
    );
  }
  return experimentResult ? (
    <div>
      {goalId && (
        <div className={classNames('bg-white m-5')}>
          <Select
            options={goalOptions}
            className={classNames('text-sm w-[300px] mb-5')}
            value={goalOptions.find((o) => o.value === goalId)}
            isSearchable={false}
            onChange={(e) => {
              setGoalId(e.value);
            }}
          />
          <GoalResultDetail experimentId={experimentId} goalId={goalId} />
        </div>
      )}
    </div>
  ) : (
    <div className="my-10 flex justify-center">
      <div className="w-[600px] text-gray-700 text-center">
        <h1 className="text-lg">
          {f(messages.noData.title, {
            title: f(messages.experiment.result.noData.experimentResult),
          })}
        </h1>
        <p className="mt-5">
          {f(messages.experiment.result.noData.description)}
        </p>
      </div>
    </div>
  );
};
