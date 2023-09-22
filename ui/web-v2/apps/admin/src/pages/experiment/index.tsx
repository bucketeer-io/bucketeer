import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { useCallback, FC, memo, useEffect, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import {
  useHistory,
  useRouteMatch,
  useParams,
  useLocation,
} from 'react-router-dom';

import { ConfirmDialog } from '../../components/ConfirmDialog';
import { ExperimentAddForm } from '../../components/ExperimentAddForm';
import { ExperimentList } from '../../components/ExperimentList';
import { ExperimentUpdateForm } from '../../components/ExperimentUpdateForm';
import { Header } from '../../components/Header';
import { Overlay } from '../../components/Overlay';
import { EXPERIMENT_LIST_PAGE_SIZE } from '../../constants/experiment';
import {
  ID_NEW,
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  archiveExperiment,
  selectById as selectExperimentById,
  createExperiment,
  getExperiment,
  listExperiments,
  OrderBy,
  OrderDirection,
  updateExperiment,
} from '../../modules/experiments';
import { useCurrentEnvironment } from '../../modules/me';
import {
  ChangeExperimentDescriptionCommand,
  ChangeExperimentNameCommand,
} from '../../proto/experiment/command_pb';
import { Experiment } from '../../proto/experiment/experiment_pb';
import { ListExperimentsRequest } from '../../proto/experiment/service_pb';
import { AppDispatch } from '../../store';
import {
  ExperimentSortOption,
  isExperimentSortOption,
} from '../../types/experiment';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
} from '../../types/list';
import {
  SearchParams,
  stringifySearchParams,
  useSearchParams,
} from '../../utils/search-params';

import { addFormSchema, updateFormSchema } from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: ExperimentSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListExperimentsRequest.OrderBy.CREATED_AT,
        orderDirection: ListExperimentsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListExperimentsRequest.OrderBy.CREATED_AT,
        orderDirection: ListExperimentsRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListExperimentsRequest.OrderBy.NAME,
        orderDirection: ListExperimentsRequest.OrderDirection.ASC,
      };
    default:
      return {
        orderBy: ListExperimentsRequest.OrderBy.NAME,
        orderDirection: ListExperimentsRequest.OrderDirection.DESC,
      };
  }
};

export const ExperimentIndexPage: FC = memo(() => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const history = useHistory();
  const location = useLocation();
  const searchParams = useSearchParams();
  const searchOptions: SearchParams = {
    ...searchParams,
    sort: searchParams.sort || '-createdAt',
  };
  const { url } = useRouteMatch();
  const { experimentId } = useParams<{ experimentId: string }>();
  const isNew = experimentId == ID_NEW;
  const isUpdate = experimentId ? experimentId != ID_NEW : false;
  const [open, setOpen] = useState(isNew || isUpdate);
  const [isArchiveConfirmDialogOpen, setIsArchiveConfirmDialogOpen] =
    useState(false);
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
  const defaultValues = {
    name: '',
    description: '',
    featureId: null,
    featureVersion: null,
    baselineVariation: null,
    goalIds: null,
    startAt: null,
    stopAt: null,
  };
  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      ...defaultValues,
      ...{ featureId: isNew && searchOptions.fid ? searchOptions.fid : null },
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const updateMethod = useForm({
    resolver: yupResolver(updateFormSchema),
    mode: 'onChange',
  });
  const {
    handleSubmit: handleUpdateSubmit,
    formState: { dirtyFields },
    reset: resetUpdate,
  } = updateMethod;

  const updateExperimentList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isExperimentSortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * EXPERIMENT_LIST_PAGE_SIZE;
      const status = options && options.status ? Number(options.status) : null;
      const archived =
        options && options.archived ? options.archived === 'true' : false;
      dispatch(
        listExperiments({
          environmentNamespace: currentEnvironment.id,
          pageSize: EXPERIMENT_LIST_PAGE_SIZE,
          cursor: String(cursor),
          searchKeyword: options && (options.q as string),
          status: status,
          maintainer: options && (options.maintainerId as string),
          archived,
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
        })
      );
    },
    [dispatch]
  );

  const updateURL = useCallback(
    (options: Record<string, string | number | boolean | undefined>) => {
      history.replace(
        `${url}?${stringifySearchParams({
          ...options,
        })}`
      );
    },
    [history]
  );

  const handleSearchOptionsChange = useCallback(
    (options) => {
      updateURL({ ...options, page: 1 });
      updateExperimentList(options, 1);
    },
    [updateURL]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateExperimentList(searchOptions, page);
    },
    [updateURL, searchOptions]
  );

  const handleOpenAdd = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const handleOpenUpdate = useCallback(
    (e: Experiment.AsObject) => {
      setOpen(true);
      resetUpdate({
        name: e.name,
        description: e.description,
        maintainer: e.maintainer,
      });
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${e.id}`,
        search: location.search,
      });
    },
    [setOpen, resetUpdate, history, experiment, location]
  );

  const handleClose = useCallback(() => {
    setOpen(false);
    resetAdd(defaultValues);
    resetUpdate();
    const { fid, ...opts } = searchParams;
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}`,
      search: stringifySearchParams(opts),
    });
  }, [searchParams, setOpen, history, resetAdd, resetUpdate]);

  const handleAdd = useCallback(
    async (data) => {
      dispatch(
        createExperiment({
          environmentNamespace: currentEnvironment.id,
          name: data.name,
          description: data.description,
          featureId: data.featureId,
          baseVariationId: data.baselineVariation,
          goalIdsList: data.goalIds,
          startAt: data.startAt.getTime() / 1000,
          stopAt: data.stopAt.getTime() / 1000,
        })
      ).then(() => {
        handleClose();
        updateExperimentList(null, 1);
      });
    },
    [dispatch]
  );

  const handleUpdate = useCallback(
    async (data) => {
      let changeExperimentNameCommand: ChangeExperimentNameCommand;
      let changeExperimentDescriptionCommand: ChangeExperimentDescriptionCommand;
      if (dirtyFields.name) {
        changeExperimentNameCommand = new ChangeExperimentNameCommand();
        changeExperimentNameCommand.setName(data.name);
      }
      if (dirtyFields.description) {
        changeExperimentDescriptionCommand =
          new ChangeExperimentDescriptionCommand();
        changeExperimentDescriptionCommand.setDescription(data.description);
      }
      dispatch(
        updateExperiment({
          environmentNamespace: currentEnvironment.id,
          id: experimentId,
          changeNameCommand: changeExperimentNameCommand,
          changeDescriptionCommand: changeExperimentDescriptionCommand,
        })
      ).then(() => {
        dispatch(
          getExperiment({
            environmentNamespace: currentEnvironment.id,
            id: experimentId,
          })
        );
        handleClose();
      });
    },
    [dispatch, dirtyFields, experimentId]
  );

  const archiveMethod = useForm({
    defaultValues: {
      experiment: null,
    },
    mode: 'onChange',
  });
  const {
    handleSubmit: archiveHandleSubmit,
    setValue: archiveSetValue,
    reset: archiveReset,
  } = archiveMethod;

  const handleClickArchive = useCallback(
    (experiment: Experiment.AsObject) => {
      archiveSetValue('experiment', experiment);
      setIsArchiveConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleArchive = useCallback(
    async (data) => {
      dispatch(
        archiveExperiment({
          environmentNamespace: currentEnvironment.id,
          id: data.experiment.id,
        })
      ).then(() => {
        archiveReset();
        setIsArchiveConfirmDialogOpen(false);
        updateExperimentList(null, 1);
      });
    },
    [dispatch, archiveReset, setIsArchiveConfirmDialogOpen]
  );

  useEffect(() => {
    if (isUpdate) {
      dispatch(
        getExperiment({
          environmentNamespace: currentEnvironment.id,
          id: experimentId,
        })
      ).then((e) => {
        const experiment = e.payload as Experiment.AsObject;
        resetUpdate({
          name: experiment.name,
          description: experiment.description,
          maintainer: experiment.maintainer,
        });
      });
    }
    updateExperimentList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [dispatch, updateExperimentList]);

  useEffect(() => {
    history.listen(() => {
      // Handle browser's back button
      if (history.action === 'POP') {
        if (open) {
          setOpen(false);
        }
      }
    });
  });

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.experiment.list.header.title)}
          description={f(messages.experiment.list.header.description)}
        />
      </div>
      <div className="m-10">
        <ExperimentList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onAdd={handleOpenAdd}
          onUpdate={handleOpenUpdate}
          onChangeSearchOptions={handleSearchOptionsChange}
          onArchive={handleClickArchive}
        />
      </div>
      <Overlay open={open} onClose={handleClose}>
        {isNew && (
          <FormProvider {...addMethod}>
            <ExperimentAddForm
              onSubmit={handleAddSubmit(handleAdd)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
        {isUpdate && (
          <FormProvider {...updateMethod}>
            <ExperimentUpdateForm
              onSubmit={handleUpdateSubmit(handleUpdate)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
      </Overlay>
      <ConfirmDialog
        open={isArchiveConfirmDialogOpen}
        onConfirm={archiveHandleSubmit(handleArchive)}
        onClose={() => setIsArchiveConfirmDialogOpen(false)}
        title={f(messages.experiment.confirm.archiveTitle)}
        description={f(messages.experiment.confirm.archiveDescription, {
          experimentName:
            archiveMethod.getValues().experiment &&
            archiveMethod.getValues().experiment.name,
        })}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
