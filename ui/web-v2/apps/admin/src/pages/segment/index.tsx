import { yupResolver } from '@hookform/resolvers/yup';
import { SerializedError } from '@reduxjs/toolkit';
import React, { useCallback, FC, memo, useEffect, useState } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import {
  useHistory,
  useRouteMatch,
  useLocation,
  useParams,
} from 'react-router-dom';

import { Header } from '../../components/Header';
import { Overlay } from '../../components/Overlay';
import { SegmentAddForm } from '../../components/SegmentAddForm';
import { SegmentDeleteDialog } from '../../components/SegmentDeleteDialog';
import { SegmentList } from '../../components/SegmentList';
import { SegmentUpdateForm } from '../../components/SegmentUpdateForm';
import { SegmentUploadingDialog } from '../../components/SegmentUploadingDialog';
import {
  ID_NEW,
  PAGE_PATH_USER_SEGMENTS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { SEGMENT_LIST_PAGE_SIZE } from '../../constants/segment';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useCurrentEnvironment } from '../../modules/me';
import {
  selectById as selectSegmentById,
  listSegments,
  getSegment,
  bulkDownloadSegmentUsers,
  bulkUploadSegmentUsers,
  createSegment,
  updateSegment,
  deleteSegmentUser,
  OrderBy,
  OrderDirection,
} from '../../modules/segments';
import { Segment } from '../../proto/feature/segment_pb';
import { ListSegmentsRequest } from '../../proto/feature/service_pb';
import { AppDispatch } from '../../store';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
} from '../../types/list';
import { SegmentSortOption, isSegmentSortOption } from '../../types/segment';
import {
  stringifySearchParams,
  useSearchParams,
} from '../../utils/search-params';

import { addFormSchema, updateFormSchema } from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: SegmentSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListSegmentsRequest.OrderBy.CREATED_AT,
        orderDirection: ListSegmentsRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListSegmentsRequest.OrderBy.CREATED_AT,
        orderDirection: ListSegmentsRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListSegmentsRequest.OrderBy.NAME,
        orderDirection: ListSegmentsRequest.OrderDirection.ASC,
      };
    default:
      return {
        orderBy: ListSegmentsRequest.OrderBy.NAME,
        orderDirection: ListSegmentsRequest.OrderDirection.DESC,
      };
  }
};

export const SegmentIndexPage: FC = memo(() => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const history = useHistory();
  const location = useLocation();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const { url } = useRouteMatch();
  const { segmentId } = useParams<{ segmentId: string }>();
  const isNew = segmentId == ID_NEW;
  const isUpdate = segmentId ? segmentId != ID_NEW : false;
  const [open, setOpen] = useState(isNew);
  const [isUploadingDialogOpen, setIsUploadingDialogOpen] = useState(false);
  const [segment, getSegmentError] = useSelector<
    AppState,
    [Segment.AsObject | undefined, SerializedError | null]
  >(
    (state) => [
      selectSegmentById(state.segments, segmentId),
      state.segments.getSegmentError,
    ],
    shallowEqual
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

  const updateSegmentList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isSegmentSortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * SEGMENT_LIST_PAGE_SIZE;
      const inUse =
        options && options.inUse != null ? options.inUse === 'true' : null;
      dispatch(
        listSegments({
          environmentNamespace: currentEnvironment.id,
          pageSize: SEGMENT_LIST_PAGE_SIZE,
          cursor: String(cursor),
          searchKeyword: options && (options.q as string),
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
          inUse: inUse,
        })
      );
    },
    [dispatch]
  );

  const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

  const deleteMethod = useForm({
    defaultValues: {
      segment: null,
    },
    mode: 'onChange',
  });

  const { handleSubmit: deleteHandleSubmit, setValue: deleteSetValue } =
    deleteMethod;

  const handleOnClickDelete = useCallback(
    (s: Segment.AsObject) => {
      deleteSetValue('segment', s);
      if (s.isInUseStatus) {
        // TODO get the feature flag list using the segment
      }
      setIsConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleDeleteSegment = useCallback(
    (data) => {
      setIsConfirmDialogOpen(false);
      dispatch(
        deleteSegmentUser({
          environmentNamespace: currentEnvironment.id,
          id: data.segment.id,
        })
      );
    },
    [dispatch, setIsConfirmDialogOpen]
  );

  const handleOnClickDownload = useCallback(
    (segment: Segment.AsObject) => {
      dispatch(
        bulkDownloadSegmentUsers({
          environmentNamespace: currentEnvironment.id,
          segmentId: segment.id,
        })
      ).then((data) => {
        const url = window.URL.createObjectURL(
          new Blob([atob(String(data.payload))])
        );
        const link = window.document.createElement('a');
        link.href = url;
        link.setAttribute(
          'download',
          `${currentEnvironment.name}-${segment.name}.csv`
        );
        window.document.body.appendChild(link);
        link.click();
        if (link.parentNode) {
          link.parentNode.removeChild(link);
        }
      });
    },
    [dispatch]
  );

  const handleOnChangeSearchOptions = useCallback(
    (options) => {
      updateURL({ ...options, page: 1 });
      updateSegmentList(options, 1);
    },
    [updateURL, updateSegmentList]
  );

  const handleOnChangePage = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateSegmentList(searchOptions, page);
    },
    [updateURL, updateSegmentList, searchOptions]
  );

  const handleOnClickAdd = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const handleOnClickUpdate = useCallback(
    (s: Segment.AsObject) => {
      if (s.status === Segment.Status.UPLOADING) {
        setIsUploadingDialogOpen(true);
      } else {
        setOpen(true);
        resetUpdate({
          name: s.name,
          description: s.description,
          isInUseStatus: s.isInUseStatus,
          status: s.status,
          file: null,
          featureList: s.featuresList,
        });
        history.push({
          pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}/${s.id}`,
          search: location.search,
        });
      }
    },
    [setOpen, history, segment, location]
  );

  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      name: '',
      description: '',
      file: null,
      userIds: '',
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

  const handleOnClose = useCallback(() => {
    resetAdd();
    resetUpdate();
    setOpen(false);
    history.replace({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}`,
      search: location.search,
    });
  }, [setOpen, history, location, resetAdd, resetUpdate]);

  const add = useCallback(
    async (data) => {
      dispatch(
        createSegment({
          environmentNamespace: currentEnvironment.id,
          name: data.name,
          description: data.description,
        })
      ).then((response) => {
        let file: File;
        if (data.file && data.file.length > 0) {
          file = data.file[0];
        } else if (data.userIds) {
          // Convert string to file object
          file = new File([data.userIds], 'filename.txt', {
            type: 'text/plain',
          });
        }

        if (file) {
          convertFileToUint8Array(file, (uint8Array) => {
            dispatch(
              bulkUploadSegmentUsers({
                environmentNamespace: currentEnvironment.id,
                segmentId: response.payload as string,
                data: uint8Array,
              })
            ).then(addFinished);
          });
        } else {
          addFinished();
        }
      });
    },
    [dispatch]
  );

  const addFinished = () => {
    resetAdd();
    setOpen(false);
    history.replace(
      `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}`
    );
    updateSegmentList(null, 1);
  };

  const convertFileToUint8Array = (
    file: Blob,
    onLoad: (data: Uint8Array) => void
  ) => {
    const reader = new FileReader();
    reader.readAsArrayBuffer(file);
    reader.onload = () => {
      onLoad(new Uint8Array(reader.result as ArrayBuffer));
    };
  };

  const update = useCallback(
    async (data) => {
      let name: string;
      let description: String;
      let file: File;

      if (dirtyFields.name) {
        name = data.name;
      }
      if (dirtyFields.description) {
        description = data.description;
      }
      if (data.file && data.file.length > 0) {
        file = data.file[0];
      } else if (data.userIds) {
        // Convert string to file object
        file = new File([data.userIds], 'filename.txt', {
          type: 'text/plain',
        });
      }

      // File only
      if (!name && !description && file) {
        convertFileToUint8Array(file, (uint8Array) => {
          dispatch(
            bulkUploadSegmentUsers({
              environmentNamespace: currentEnvironment.id,
              segmentId: segmentId,
              data: uint8Array,
            })
          ).then(() => {
            dispatch(
              getSegment({
                environmentNamespace: currentEnvironment.id,
                id: segmentId,
              })
            ).then(handleOnClose);
          });
        });
        return;
      }
      // Name, description and file
      dispatch(
        updateSegment({
          environmentNamespace: currentEnvironment.id,
          id: segmentId,
          name: name,
          description: description,
        })
      ).then(() => {
        if (!file) {
          dispatch(
            getSegment({
              environmentNamespace: currentEnvironment.id,
              id: segmentId,
            })
          ).then(handleOnClose);
          return;
        }
        convertFileToUint8Array(file, (uint8Array) => {
          dispatch(
            bulkUploadSegmentUsers({
              environmentNamespace: currentEnvironment.id,
              segmentId: segmentId,
              data: uint8Array,
            })
          ).then(() => {
            dispatch(
              getSegment({
                environmentNamespace: currentEnvironment.id,
                id: segmentId,
              })
            ).then(handleOnClose);
          });
        });
      });
    },
    [dispatch, segmentId]
  );

  const handleUploadingClose = () => {
    setIsUploadingDialogOpen(false);
    updateSegmentList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  };

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
    updateSegmentList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [updateSegmentList]);

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.segment.list.header.title)}
          description={f(messages.segment.list.header.description)}
        />
      </div>
      <div className="m-10">
        <SegmentList
          searchOptions={searchOptions}
          onChangePage={handleOnChangePage}
          onDelete={handleOnClickDelete}
          onDownload={handleOnClickDownload}
          onAdd={handleOnClickAdd}
          onUpdate={handleOnClickUpdate}
          onChangeSearchOptions={handleOnChangeSearchOptions}
        />
        <Overlay open={open} onClose={handleOnClose}>
          {isNew && (
            <FormProvider {...addMethod}>
              <SegmentAddForm
                onSubmit={handleAddSubmit(add)}
                onCancel={handleOnClose}
              />
            </FormProvider>
          )}
          {isUpdate && (
            <FormProvider {...updateMethod}>
              <SegmentUpdateForm
                onSubmit={handleUpdateSubmit(update)}
                onCancel={handleOnClose}
              />
            </FormProvider>
          )}
        </Overlay>
      </div>
      <SegmentDeleteDialog
        open={isConfirmDialogOpen}
        segment={deleteMethod.getValues().segment}
        onConfirm={deleteHandleSubmit(handleDeleteSegment)}
        onClose={() => setIsConfirmDialogOpen(false)}
      />
      <SegmentUploadingDialog
        open={isUploadingDialogOpen}
        onClose={handleUploadingClose}
      />
    </>
  );
});
