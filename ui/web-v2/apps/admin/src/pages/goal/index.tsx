import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { useCallback, FC, memo, useEffect, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory, useRouteMatch, useParams } from 'react-router-dom';

import { ConfirmDialog } from '../../components/ConfirmDialog';
import { GoalAddForm } from '../../components/GoalAddForm';
import { GoalList } from '../../components/GoalList';
import { GoalUpdateForm } from '../../components/GoalUpdateForm';
import { Header } from '../../components/Header';
import { Overlay } from '../../components/Overlay';
import { GOAL_LIST_PAGE_SIZE } from '../../constants/goal';
import {
  ID_NEW,
  PAGE_PATH_GOALS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  archiveGoal,
  selectById as selectGoalById,
  createGoal,
  listGoals,
  getGoal,
  updateGoal,
  OrderBy,
  OrderDirection,
} from '../../modules/goals';
import { useCurrentEnvironment } from '../../modules/me';
import { Goal } from '../../proto/experiment/goal_pb';
import { ListGoalsRequest } from '../../proto/experiment/service_pb';
import { AppDispatch } from '../../store';
import { GoalSortOption, isGoalSortOption } from '../../types/goal';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
} from '../../types/list';
import {
  stringifySearchParams,
  useSearchParams,
} from '../../utils/search-params';

import { addFormSchema, updateFormSchema } from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: GoalSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListGoalsRequest.OrderBy.CREATED_AT,
        orderDirection: ListGoalsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListGoalsRequest.OrderBy.CREATED_AT,
        orderDirection: ListGoalsRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListGoalsRequest.OrderBy.NAME,
        orderDirection: ListGoalsRequest.OrderDirection.ASC,
      };
    default:
      return {
        orderBy: ListGoalsRequest.OrderBy.NAME,
        orderDirection: ListGoalsRequest.OrderDirection.DESC,
      };
  }
};

export const GoalIndexPage: FC = memo(() => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const history = useHistory();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const { url } = useRouteMatch();
  const { goalId } = useParams<{ goalId: string }>();
  const isNew = goalId == ID_NEW;
  const isUpdate = goalId ? goalId != ID_NEW : false;
  const [open, setOpen] = useState(isNew);
  const [isArchiveConfirmDialogOpen, setIsArchiveConfirmDialogOpen] =
    useState(false);
  const [goal] = useSelector<
    AppState,
    [Goal.AsObject | undefined, SerializedError | null]
  >(
    (state) => [selectGoalById(state.goals, goalId), state.goals.getGoalError],
    shallowEqual
  );
  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      id: '',
      name: '',
      description: '',
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const handleAdd = useCallback(
    async (data) => {
      dispatch(
        createGoal({
          environmentNamespace: currentEnvironment.id,
          id: data.id,
          name: data.name,
          description: data.description,
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        history.replace(
          `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_GOALS}`
        );
        updateGoalList(null, 1);
      });
    },
    [dispatch]
  );

  const updateMethod = useForm({
    resolver: yupResolver(updateFormSchema),
    mode: 'onChange',
  });

  const {
    handleSubmit: handleUpdateSubmit,
    formState: { dirtyFields },
    reset: resetUpdate,
  } = updateMethod;

  const handleUpdate = useCallback(
    async (data) => {
      let name: string;
      let description: string;
      if (dirtyFields.name) {
        name = data.name;
      }
      if (dirtyFields.description) {
        description = data.description;
      }
      dispatch(
        updateGoal({
          environmentNamespace: currentEnvironment.id,
          id: goalId,
          name: name,
          description: description,
        })
      ).then(() => {
        dispatch(
          getGoal({
            environmentNamespace: currentEnvironment.id,
            id: goalId,
          })
        );
        handleClose();
      });
    },
    [dispatch, goalId, dirtyFields]
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

  const updateGoalList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isGoalSortOption(options && options.sort) ? options.sort : '-createdAt'
      );
      const cursor = (page - 1) * GOAL_LIST_PAGE_SIZE;
      const status =
        options && options.status != null ? options.status === 'true' : null;
      const archived =
        options && options.archived != null
          ? options.archived === 'true'
          : false;
      dispatch(
        listGoals({
          environmentNamespace: currentEnvironment.id,
          pageSize: GOAL_LIST_PAGE_SIZE,
          cursor: String(cursor),
          searchKeyword: options && (options.q as string),
          status: status,
          archived,
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
        })
      );
    },
    [dispatch]
  );

  const handleSearchOptionsChange = useCallback(
    (options) => {
      updateURL({ ...options, page: 1 });
      updateGoalList(options, 1);
    },
    [updateURL, updateGoalList]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateGoalList(searchOptions, page);
    },
    [updateURL, updateGoalList, searchOptions]
  );

  const handleOpen = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_GOALS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const handleOpenUpdate = useCallback(
    (g: Goal.AsObject) => {
      setOpen(true);
      resetUpdate({
        id: g.id,
        name: g.name,
        description: g.description,
      });
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_GOALS}/${g.id}`,
        search: location.search,
      });
    },
    [setOpen, history, goal, location]
  );

  const archiveMethod = useForm({
    defaultValues: {
      goal: null,
    },
    mode: 'onChange',
  });
  const {
    handleSubmit: archiveHandleSubmit,
    setValue: archiveSetValue,
    reset: archiveReset,
  } = archiveMethod;

  const handleClickArchive = useCallback(
    (goal: Goal.AsObject) => {
      archiveSetValue('goal', goal);
      setIsArchiveConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleArchive = useCallback(
    async (data) => {
      dispatch(
        archiveGoal({
          environmentNamespace: currentEnvironment.id,
          id: data.goal.id,
        })
      ).then(() => {
        archiveReset();
        setIsArchiveConfirmDialogOpen(false);
        updateGoalList(null, 1);
      });
    },
    [dispatch, archiveReset, setIsArchiveConfirmDialogOpen]
  );

  const handleClose = useCallback(() => {
    setOpen(false);
    resetAdd();
    resetUpdate();
    history.replace({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_GOALS}`,
      search: location.search,
    });
  }, [setOpen, history, location, resetAdd, resetUpdate]);

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

  useEffect(() => {
    updateGoalList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [updateGoalList]);

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.goal.list.header.title)}
          description={f(messages.goal.list.header.description)}
        />
      </div>
      <div className="m-10">
        <GoalList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onAdd={handleOpen}
          onChangeSearchOptions={handleSearchOptionsChange}
          onUpdate={handleOpenUpdate}
          onArchive={handleClickArchive}
        />
      </div>
      <Overlay open={open} onClose={handleClose}>
        {isNew && (
          <FormProvider {...addMethod}>
            <GoalAddForm
              onSubmit={handleAddSubmit(handleAdd)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
        {isUpdate && (
          <FormProvider {...updateMethod}>
            <GoalUpdateForm
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
        title={f(messages.goal.confirm.archiveTitle)}
        description={f(messages.goal.confirm.archiveDescription, {
          goalId:
            archiveMethod.getValues().goal && archiveMethod.getValues().goal.id,
        })}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
