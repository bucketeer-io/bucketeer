import { isExperimentStatusWaitingRunnning } from '../../components/ExperimentList';
import { formatDate } from '../../utils/date';
import { PlusIcon, SelectorIcon } from '@heroicons/react/solid';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';

import { ConfirmDialog } from '../../components/ConfirmDialog';
import { ExperimentResultDetail } from '../../components/ExperimentResultDetail';
import { statusOptions } from '../../components/ExperimentSearch';
import { HoverPopover } from '../../components/HoverPopover';
import {
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  getExperiment,
  listExperiments,
  selectAll as selectAllExperiment,
  stopExperiment
} from '../../modules/experiments';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { Experiment } from '../../proto/experiment/experiment_pb';
import { ListExperimentsRequest } from '../../proto/experiment/service_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';
import { useSearchParams } from '../../utils/search-params';
import { Select } from '../../components/Select';

interface FeatureExperimentsPageProps {
  featureId: string;
}

export const FeatureExperimentsPage: FC<FeatureExperimentsPageProps> = memo(
  ({ featureId }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const dispatch = useDispatch<AppDispatch>();
    const searchOptions = useSearchParams();
    const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);
    const experiments = useSelector<AppState, Experiment.AsObject[]>(
      (state) => {
        const exprs = selectAllExperiment(state.experiments);
        exprs?.sort((a, b) => b.createdAt - a.createdAt);
        return exprs;
      },
      shallowEqual
    );
    const isLoading = useSelector<AppState, boolean>(
      (state) => state.experiments.loading,
      shallowEqual
    );
    const [experimentId, setExperimentId] = useState<string>('');
    const history = useHistory();
    const currentEnvironment = useCurrentEnvironment();
    const experimentOptions = experiments.map((e) => {
      return {
        value: e.id,
        label: e.name
      };
    });
    const experiment = experiments?.find((e) => e.id === experimentId);

    const handleStop = useCallback(async () => {
      dispatch(
        stopExperiment({
          environmentNamespace: currentEnvironment.id,
          experimentId: experimentId
        })
      ).then(() => {
        dispatch(
          getExperiment({
            environmentNamespace: currentEnvironment.id,
            id: experimentId
          })
        );
        setIsConfirmDialogOpen(false);
      });
    }, [dispatch, experimentId]);

    useEffect(() => {
      dispatch(
        listExperiments({
          featureId: featureId,
          environmentNamespace: currentEnvironment.id,
          searchKeyword: '',
          pageSize: 1000,
          cursor: '',
          orderBy: ListExperimentsRequest.OrderBy.CREATED_AT,
          orderDirection: ListExperimentsRequest.OrderDirection.DESC
        })
      );
    }, [dispatch, featureId, currentEnvironment, searchOptions]);

    useEffect(() => {
      const id = searchOptions.experimentId
        ? (searchOptions.experimentId as string)
        : experiments[0]?.id;
      setExperimentId(id);
    }, [experiments]);

    if (isLoading) {
      return <div>loading</div>;
    }
    return (
      <>
        <div className="p-10 bg-gray-100 space-y-4">
          <div className="flex items-center">
            {experimentOptions.length > 0 && (
              <Select
                className={classNames('text-sm w-[300px]')}
                onChange={(o) => {
                  setExperimentId(o.value.toString());
                }}
                options={experimentOptions}
                value={
                  experimentId
                    ? experimentOptions.find((o) => o.value === experimentId)
                    : experimentOptions[0]
                }
              />
            )}
            <div className="flex-grow" />
            {experiment && editable && (
              <button
                type="button"
                className="btn-cancel mx-2"
                disabled={false}
                onClick={() =>
                  history.push(
                    `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${experiment.id}`
                  )
                }
              >
                {f(messages.button.edit)}
              </button>
            )}
            {editable && (
              <button
                type="button"
                className="btn-submit mx-2"
                disabled={false}
                onClick={() =>
                  history.push(
                    `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}${PAGE_PATH_NEW}?fid=${featureId}`
                  )
                }
              >
                <PlusIcon className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
                {f(messages.button.add)}
              </button>
            )}
          </div>
          {experimentOptions.length == 0 && (
            <div className="my-10 flex justify-center">
              <div className="w-[600px] text-gray-700 text-center">
                <h1 className="text-lg">
                  {f(messages.noData.title, {
                    title: f(messages.experiment.list.header.title)
                  })}
                </h1>
                <p className="mt-5">
                  {f(messages.experiment.list.noData.description)}
                </p>
                <a
                  href="https://docs.bucketeer.io/feature-flags/testing-with-flags"
                  target="_blank"
                  rel="noreferrer"
                  className="link"
                >
                  {f(messages.readMore)}
                </a>
              </div>
            </div>
          )}
          {experiment ? (
            <div className="space-y-5">
              <div className="bg-white rounded-md border border-gray-300 p-5">
                <ExperimentDetail
                  experiment={experiment}
                  onStopExperiment={() => setIsConfirmDialogOpen(true)}
                />
              </div>
              <div className="bg-white rounded-md border border-gray-300">
                <ExperimentResultDetail experimentId={experimentId} />
              </div>
            </div>
          ) : null}
        </div>
        <ConfirmDialog
          open={isConfirmDialogOpen}
          title={f(messages.experiment.stop.dialog.title)}
          description={f(messages.experiment.stop.dialog.description)}
          onConfirm={handleStop}
          onClose={() => setIsConfirmDialogOpen(false)}
          onCloseButton={f(messages.button.cancel)}
          onConfirmButton={f(messages.button.submit)}
        />
      </>
    );
  }
);

interface ExperimentDetailProps {
  experiment: Experiment.AsObject;
  onStopExperiment: () => void;
}

export const ExperimentDetail: FC<ExperimentDetailProps> = memo(
  ({ experiment, onStopExperiment }) => {
    const { formatMessage: f } = useIntl();
    const startAt = new Date(experiment.startAt * 1000);
    const endAt =
      experiment.status === Experiment.Status.FORCE_STOPPED
        ? new Date(Number(experiment.stoppedAt) * 1000)
        : new Date(experiment.stopAt * 1000);

    return (
      <div className="space-y-3">
        <div className="flex flex-row">
          <StatusButton status={experiment.status} onClick={onStopExperiment} />
          <div className="flex-grow" />
          <div className="inline-flex justify-center text-sm">
            {`${f(messages.experiment.period)}
            ${formatDate({ date: startAt })} - ${formatDate({
              date: endAt
            })}`}
          </div>
        </div>
        {experiment.description && (
          <div className="text-sm">
            <p>{f(messages.description)}</p>
            <p>{experiment.description}</p>
          </div>
        )}
      </div>
    );
  }
);

export interface StatusButtonProps {
  status: Experiment.StatusMap[keyof Experiment.StatusMap];
  onClick: () => void;
}

export const StatusButton: FC<StatusButtonProps> = memo(
  ({ status, onClick }) => {
    const editable = useIsEditable();
    const { formatMessage: f } = useIntl();
    const isStoppable = isExperimentStatusWaitingRunnning(status);

    return isStoppable ? (
      <HoverPopover
        disabled={!editable}
        render={() => {
          return (
            <button
              type="button"
              className={classNames(
                'bg-white border text-gray-600',
                'p-1.5 text-xs rounded',
                'hover:filter hover:brightness-75'
              )}
              onClick={onClick}
            >
              {f(messages.experiment.stop.button)}
            </button>
          );
        }}
      >
        <span
          className={classNames(
            'inline-flex justify-center',
            'py-2 pl-4 pr-3 border border-gray-300',
            'text-sm font-medium rounded-md',
            'text-white bg-primary',
            editable && 'hover:filter hover:brightness-75'
          )}
        >
          {
            statusOptions.find((option) => option.value == status.toString())
              .label
          }
          {editable && (
            <SelectorIcon className="w-5 h-5 text-white" aria-hidden="true" />
          )}
        </span>
      </HoverPopover>
    ) : (
      <span
        className={classNames(
          'inline-flex justify-center',
          'py-2 px-4 border border-gray-300',
          'text-sm font-medium rounded-md',
          'text-white bg-primary'
        )}
      >
        {
          statusOptions.find((option) => option.value == status.toString())
            .label
        }
      </span>
    );
  }
);
