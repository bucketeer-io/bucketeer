import { yupResolver } from '@hookform/resolvers/yup';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch } from 'react-redux';
import { useHistory, useRouteMatch, useParams } from 'react-router-dom';

import { ConfirmDialog } from '../../../components/ConfirmDialog';
import { Overlay } from '../../../components/Overlay';
import { ProjectAddForm } from '../../../components/ProjectAddForm';
import { ProjectList } from '../../../components/ProjectList';
import { ProjectUpdateForm } from '../../../components/ProjectUpdateForm';
import { PROJECT_LIST_PAGE_SIZE } from '../../../constants/project';
import {
  ID_NEW,
  PAGE_PATH_ADMIN,
  PAGE_PATH_PROJECTS,
  PAGE_PATH_NEW,
} from '../../../constants/routing';
import { messages } from '../../../lang/messages';
import {
  convertProject,
  createProject,
  disableProject,
  enableProject,
  getProject,
  listProjects,
  updateProject,
  OrderBy,
  OrderDirection,
} from '../../../modules/projects';
import { Project } from '../../../proto/environment/project_pb';
import { ListProjectsRequest } from '../../../proto/environment/service_pb';
import { AppDispatch } from '../../../store';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_DESC,
  SORT_OPTIONS_NAME_ASC,
} from '../../../types/list';
import { ProjectSortOption, isProjectSortOption } from '../../../types/project';
import {
  useSearchParams,
  stringifySearchParams,
} from '../../../utils/search-params';

import {
  addFormSchema,
  updateFormSchema,
} from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: ProjectSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListProjectsRequest.OrderBy.CREATED_AT,
        orderDirection: ListProjectsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListProjectsRequest.OrderBy.CREATED_AT,
        orderDirection: ListProjectsRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListProjectsRequest.OrderBy.NAME,
        orderDirection: ListProjectsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_NAME_DESC:
      return {
        orderBy: ListProjectsRequest.OrderBy.NAME,
        orderDirection: ListProjectsRequest.OrderDirection.DESC,
      };
    default:
      return {
        orderBy: ListProjectsRequest.OrderBy.CREATED_AT,
        orderDirection: ListProjectsRequest.OrderDirection.DESC,
      };
  }
};

export const AdminProjectIndexPage: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const history = useHistory();
  const { url } = useRouteMatch();
  const { projectId } = useParams<{ projectId: string }>();
  const isNew = projectId == ID_NEW;
  const isUpdate = projectId ? projectId != ID_NEW : false;
  const [open, setOpen] = useState(isNew);
  const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);
  const [isConvertConfirmDialogOpen, setIsConvertConfirmDialogOpen] =
    useState(false);
  const updateProjectList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isProjectSortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * PROJECT_LIST_PAGE_SIZE;
      const disabled =
        options && options.enabled ? options.enabled === 'false' : null;
      dispatch(
        listProjects({
          pageSize: PROJECT_LIST_PAGE_SIZE,
          cursor: String(cursor),
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
          searchKeyword: options && (options.q as string),
          disabled,
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
      updateProjectList(options, 1);
    },
    [updateURL, updateProjectList]
  );

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateProjectList(searchOptions, page);
    },
    [updateURL, searchOptions, updateProjectList]
  );

  const handleOpenAdd = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ADMIN}${PAGE_PATH_PROJECTS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      name: '',
      urlCode: '',
      description: '',
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset: resetAdd } = addMethod;

  const handleOpenUpdate = useCallback(
    (p: Project.AsObject) => {
      setOpen(true);
      resetUpdate({
        id: p.id,
        name: p.name,
        urlCode: p.urlCode,
        description: p.description,
        creatorEmail: p.creatorEmail,
      });
      history.push({
        pathname: `${PAGE_PATH_ADMIN}${PAGE_PATH_PROJECTS}/${p.urlCode}`,
        search: location.search,
      });
    },
    [setOpen, history, location]
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

  const handleClose = useCallback(() => {
    resetAdd();
    resetUpdate();
    setOpen(false);
    history.replace({
      pathname: `${PAGE_PATH_ADMIN}${PAGE_PATH_PROJECTS}`,
      search: location.search,
    });
  }, [setOpen, history, location, resetAdd]);

  const handleAdd = useCallback(
    async (data) => {
      dispatch(
        createProject({
          name: data.name,
          urlCode: data.urlCode,
          description: data.description,
        })
      ).then(() => {
        resetAdd();
        setOpen(false);
        history.replace(`${PAGE_PATH_ADMIN}${PAGE_PATH_PROJECTS}`);
        updateProjectList(null, 1);
      });
    },
    [dispatch]
  );

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
        updateProject({
          id: data.id,
          name: name,
          description: description,
        })
      ).then(() => {
        dispatch(
          getProject({
            id: data.id,
          })
        );
        handleClose();
      });
    },
    [dispatch, dirtyFields]
  );

  const switchEnabledMethod = useForm({
    defaultValues: {
      projectId: '',
      enabled: false,
    },
    mode: 'onChange',
  });

  const {
    handleSubmit: switchEnableHandleSubmit,
    setValue: switchEnabledSetValue,
  } = switchEnabledMethod;

  const handleClickSwitchEnabled = useCallback(
    (project: Project.AsObject) => {
      switchEnabledSetValue('projectId', project.id);
      switchEnabledSetValue('enabled', project.disabled);
      setIsConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleSwitchEnabled = useCallback(
    async (data) => {
      dispatch(
        (() => {
          if (data.enabled) {
            return enableProject({
              id: data.projectId,
            });
          }
          return disableProject({
            id: data.projectId,
          });
        })()
      ).then(() => {
        setIsConfirmDialogOpen(false);
        dispatch(
          getProject({
            id: data.projectId,
          })
        );
      });
    },
    [dispatch, setIsConfirmDialogOpen]
  );

  const convertMethod = useForm({
    defaultValues: {
      project: null,
    },
    mode: 'onChange',
  });

  const {
    handleSubmit: convertHandleSubmit,
    setValue: convertSetValue,
    reset: convertReset,
  } = convertMethod;

  const handleClickConvert = useCallback(
    (project: Project.AsObject) => {
      convertSetValue('project', project);
      setIsConvertConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleConvert = useCallback(
    async (data) => {
      dispatch(
        convertProject({
          id: data.project.id,
        })
      ).then(() => {
        convertReset();
        setIsConvertConfirmDialogOpen(false);
        updateProjectList(null, 1);
      });
    },
    [dispatch, convertReset, setIsConvertConfirmDialogOpen]
  );

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
    updateProjectList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [updateProjectList]);

  return (
    <>
      <div className="flex items-stretch m-10 text-sm">
        <p className="text-gray-700">
          {f(messages.adminProject.list.header.description)}
        </p>
        <a
          className="link"
          target="_blank"
          href="https://bucketeer.io/docs/#/projects"
          rel="noreferrer"
        >
          {f(messages.readMore)}
        </a>
      </div>
      <div className="m-10">
        <ProjectList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onSwitchEnabled={handleClickSwitchEnabled}
          onChangeSearchOptions={handleSearchOptionsChange}
          onAdd={handleOpenAdd}
          onUpdate={handleOpenUpdate}
          onConvert={handleClickConvert}
        />
      </div>
      <Overlay open={open} onClose={handleClose}>
        {isNew && (
          <FormProvider {...addMethod}>
            <ProjectAddForm
              onSubmit={handleAddSubmit(handleAdd)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
        {isUpdate && (
          <FormProvider {...updateMethod}>
            <ProjectUpdateForm
              onSubmit={handleUpdateSubmit(handleUpdate)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
      </Overlay>
      <ConfirmDialog
        open={isConfirmDialogOpen}
        onConfirm={switchEnableHandleSubmit(handleSwitchEnabled)}
        onClose={() => setIsConfirmDialogOpen(false)}
        title={
          switchEnabledMethod.getValues().enabled
            ? f(messages.adminProject.confirm.enableTitle)
            : f(messages.adminProject.confirm.disableTitle)
        }
        description={f(
          switchEnabledMethod.getValues().enabled
            ? messages.adminProject.confirm.enableDescription
            : messages.adminProject.confirm.disableDescription,
          {
            projectId: switchEnabledMethod.getValues().projectId,
          }
        )}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
      <ConfirmDialog
        open={isConvertConfirmDialogOpen}
        onConfirm={convertHandleSubmit(handleConvert)}
        onClose={() => setIsConvertConfirmDialogOpen(false)}
        title={f(messages.adminProject.confirm.convertProjectTitle)}
        description={f(
          messages.adminProject.confirm.convertProjectDescription,
          {
            projectId:
              convertMethod.getValues().project &&
              convertMethod.getValues().project.id,
          }
        )}
        onCloseButton={f(messages.button.cancel)}
        onConfirmButton={f(messages.button.submit)}
      />
    </>
  );
});
