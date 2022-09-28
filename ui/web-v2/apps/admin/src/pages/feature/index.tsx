import { yupResolver } from '@hookform/resolvers/yup';
import { unwrapResult } from '@reduxjs/toolkit';
import React, { useCallback, FC, memo, useEffect, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch } from 'react-redux';
import { useHistory, useRouteMatch, useParams } from 'react-router-dom';

import { FeatureAddForm } from '../../components/FeatureAddForm';
import { FeatureCloneForm } from '../../components/FeatureCloneForm';
import { FeatureConfirmDialog } from '../../components/FeatureConfirmDialog';
import { FeatureList } from '../../components/FeatureList';
import { Header } from '../../components/Header';
import { Overlay } from '../../components/Overlay';
import {
  FEATURE_ACCOUNT_PAGE_SIZE,
  FEATURE_LIST_PAGE_SIZE,
} from '../../constants/feature';
import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_CLONE,
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
} from '../../constants/routing';
import { messages } from '../../lang/messages';
import { listAccounts } from '../../modules/accounts';
import {
  createFeature,
  cloneFeature,
  disableFeature,
  enableFeature,
  archiveFeature,
  unarchiveFeature,
  getFeature,
  listFeatures,
  OrderBy,
  OrderDirection,
} from '../../modules/features';
import { useCurrentEnvironment, useEnvironments } from '../../modules/me';
import { Feature } from '../../proto/feature/feature_pb';
import { ListFeaturesRequest } from '../../proto/feature/service_pb';
import { AppDispatch } from '../../store';
import { isFeatureSortOption, FeatureSortOption } from '../../types/feature';
import {
  SORT_OPTIONS_CREATED_AT_ASC,
  SORT_OPTIONS_CREATED_AT_DESC,
  SORT_OPTIONS_NAME_ASC,
} from '../../types/list';
import {
  stringifySearchParams,
  useSearchParams,
} from '../../utils/search-params';

import {
  addFormSchema,
  cloneSchema,
  switchEnabledFormSchema,
  archiveFormSchema,
} from './formSchema';

interface Sort {
  orderBy: OrderBy;
  orderDirection: OrderDirection;
}

const createSort = (sortOption?: FeatureSortOption): Sort => {
  switch (sortOption) {
    case SORT_OPTIONS_CREATED_AT_ASC:
      return {
        orderBy: ListFeaturesRequest.OrderBy.CREATED_AT,
        orderDirection: ListFeaturesRequest.OrderDirection.ASC,
      };
    case SORT_OPTIONS_CREATED_AT_DESC:
      return {
        orderBy: ListFeaturesRequest.OrderBy.CREATED_AT,
        orderDirection: ListFeaturesRequest.OrderDirection.DESC,
      };
    case SORT_OPTIONS_NAME_ASC:
      return {
        orderBy: ListFeaturesRequest.OrderBy.NAME,
        orderDirection: ListFeaturesRequest.OrderDirection.ASC,
      };
    default:
      return {
        orderBy: ListFeaturesRequest.OrderBy.NAME,
        orderDirection: ListFeaturesRequest.OrderDirection.DESC,
      };
  }
};

export const FeatureIndexPage: FC = memo(() => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const environments = useEnvironments();
  const history = useHistory();
  const searchOptions = useSearchParams();
  searchOptions.sort = searchOptions.sort ? searchOptions.sort : '-createdAt';
  const { url } = useRouteMatch();
  const { featureId } = useParams<{ featureId: string }>();
  const isNew = `/${url.substring(url.lastIndexOf('/') + 1)}` == PAGE_PATH_NEW;
  const isClone =
    url.substring(url.indexOf(PAGE_PATH_FEATURE_CLONE)) ===
    PAGE_PATH_FEATURE_CLONE + '/' + featureId;
  const [open, setOpen] = useState(isNew || isClone);
  const addMethod = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      id: '',
      name: '',
      description: '',
      tags: [],
      variationType: Feature.VariationType.BOOLEAN.toString(),
      variations: [
        {
          value: 'true',
          name: '',
          description: '',
        },
        {
          value: 'false',
          name: '',
          description: '',
        },
      ],
      onVariation: { value: 0, label: `${f(messages.feature.variation)} 1` },
      offVariation: { value: 1, label: `${f(messages.feature.variation)} 2` },
    },
    mode: 'onChange',
  });
  const { handleSubmit: handleAddSubmit, reset } = addMethod;
  const switchEnabledMethod = useForm({
    resolver: yupResolver(switchEnabledFormSchema),
    defaultValues: {
      featureId: '',
      enabled: false,
      comment: '',
    },
    mode: 'onChange',
  });
  const {
    handleSubmit: switchEnableHandleSubmit,
    setValue: switchEnabledSetValue,
    reset: switchEnabledReset,
  } = switchEnabledMethod;
  const archiveMethod = useForm({
    resolver: yupResolver(archiveFormSchema),
    defaultValues: {
      feature: null,
      featureId: '',
      comment: '',
    },
    mode: 'onChange',
  });
  const {
    handleSubmit: archiveHandleSubmit,
    setValue: archiveSetValue,
    reset: archiveReset,
  } = archiveMethod;

  const cloneMethod = useForm({
    resolver: yupResolver(cloneSchema),
    defaultValues: {
      feature: null,
    },
    mode: 'onChange',
  });
  const {
    handleSubmit: cloneHandleSubmit,
    setValue: cloneSetValue,
    reset: cloneReset,
  } = cloneMethod;

  const [isSwitchEnableConfirmDialogOpen, setIsSwitchEnableConfirmDialogOpen] =
    useState(false);
  const [isArchiveConfirmDialogOpen, setIsArchiveConfirmDialogOpen] =
    useState(false);

  const updateFeatureList = useCallback(
    (options, page: number) => {
      const sort = createSort(
        isFeatureSortOption(options && options.sort)
          ? options.sort
          : '-createdAt'
      );
      const cursor = (page - 1) * FEATURE_LIST_PAGE_SIZE;
      const enabled =
        options && options.enabled ? options.enabled === 'true' : null;
      const archived =
        options && options.archived ? options.archived === 'true' : false;
      const hasExperiment =
        options && options.hasExperiment
          ? options.hasExperiment === 'true'
          : null;
      dispatch(
        listFeatures({
          environmentNamespace: currentEnvironment.namespace,
          pageSize: FEATURE_LIST_PAGE_SIZE,
          cursor: String(cursor),
          tags: [],
          orderBy: sort.orderBy,
          orderDirection: sort.orderDirection,
          searchKeyword: options && (options.q as string),
          enabled: enabled,
          archived: archived,
          hasExperiment: hasExperiment,
          maintainerId: options && (options.maintainerId as string),
        })
      );
    },
    [dispatch]
  );

  const clearURLParameters = useCallback(() => {
    history.replace(url);
  }, [history]);

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
      updateFeatureList(options, 1);
    },
    [updateURL, updateFeatureList]
  );

  const handleClearSearchOptionsChange = useCallback(() => {
    clearURLParameters();
    updateFeatureList(null, 1);
  }, [updateFeatureList]);

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...searchOptions, page });
      updateFeatureList(searchOptions, page);
    },
    [updateURL, searchOptions, updateFeatureList]
  );

  const handleOpen = useCallback(() => {
    setOpen(true);
    history.push({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_FEATURES}${PAGE_PATH_NEW}`,
      search: location.search,
    });
  }, [setOpen, history, location]);

  const handleClose = useCallback(() => {
    setOpen(false);
    reset();
    cloneReset();
    history.replace({
      pathname: `${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_FEATURES}`,
      search: location.search,
    });
  }, [setOpen, history, location, reset]);

  const handleAdd = useCallback(
    async (data) => {
      dispatch(
        createFeature({
          environmentNamespace: currentEnvironment.namespace,
          id: data.id,
          name: data.name,
          description: data.description,
          tagsList: data.tags,
          variationType: data.variationType,
          variations: data.variations.map((variation) => {
            return {
              value: variation.value,
              name: variation.name,
              description: variation.description,
            };
          }),
          defaultOnVariationIndex: data.onVariation.value,
          defaultOffVariationIndex: data.offVariation.value,
        })
      ).then(() => {
        setOpen(false);
        history.push(
          `${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_FEATURES}/${data.id}${PAGE_PATH_FEATURE_TARGETING}`
        );
      });
    },
    [dispatch]
  );

  const handleClone = useCallback(
    async (data) => {
      const destinationEnvironment = environments.find(
        (o) => o.id == data.destinationEnvironmentId
      );
      dispatch(
        cloneFeature({
          environmentNamespace: currentEnvironment.namespace,
          id: featureId,
          destinationEnvironmentNamespace: destinationEnvironment.namespace,
        })
      )
        .then(unwrapResult)
        .then(() => {
          cloneReset();
          history.replace(
            `${PAGE_PATH_ROOT}${destinationEnvironment.id}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_TARGETING}`
          );
        })
        .catch(() => {
          cloneReset();
        });
    },
    [dispatch, featureId]
  );

  const handleClickSwitchEnabled = useCallback(
    (featureId: string, enabled: boolean) => {
      switchEnabledSetValue('featureId', featureId);
      switchEnabledSetValue('enabled', enabled);
      setIsSwitchEnableConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleSwitchEnabled = useCallback(
    async (data) => {
      dispatch(
        (() => {
          if (data.enabled) {
            return enableFeature({
              environmentNamespace: currentEnvironment.namespace,
              id: data.featureId,
              comment: data.comment,
            });
          }
          return disableFeature({
            environmentNamespace: currentEnvironment.namespace,
            id: data.featureId,
            comment: data.comment,
          });
        })()
      ).then(() => {
        switchEnabledReset();
        setIsSwitchEnableConfirmDialogOpen(false);
        dispatch(
          getFeature({
            environmentNamespace: currentEnvironment.namespace,
            id: data.featureId,
          })
        );
      });
    },
    [dispatch, switchEnabledReset, setIsSwitchEnableConfirmDialogOpen]
  );

  const handleClickArchive = useCallback(
    (feature: Feature.AsObject) => {
      archiveSetValue('feature', feature);
      archiveSetValue('featureId', feature.id);
      setIsArchiveConfirmDialogOpen(true);
    },
    [dispatch]
  );

  const handleClickClone = useCallback(
    (feature: Feature.AsObject) => {
      setOpen(true);
      cloneSetValue('feature', feature);
      history.push({
        pathname: `${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_FEATURES}${PAGE_PATH_FEATURE_CLONE}/${feature.id}`,
        search: location.search,
      });
    },
    [dispatch, history, location]
  );

  const handleArchive = useCallback(
    async (data) => {
      dispatch(
        data.feature.archived
          ? unarchiveFeature({
              environmentNamespace: currentEnvironment.namespace,
              id: data.feature.id,
              comment: data.comment,
            })
          : archiveFeature({
              environmentNamespace: currentEnvironment.namespace,
              id: data.feature.id,
              comment: data.comment,
            })
      ).then(() => {
        archiveReset();
        setIsArchiveConfirmDialogOpen(false);
        history.replace(
          `${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_FEATURES}`
        );
        updateFeatureList(null, 1);
      });
    },
    [dispatch, archiveReset, setIsArchiveConfirmDialogOpen]
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
    if (isClone) {
      dispatch(
        getFeature({
          environmentNamespace: currentEnvironment.namespace,
          id: featureId,
        })
      ).then((e) => {
        const feature = e.payload as Feature.AsObject;
        cloneReset({
          feature: feature,
        });
      });
    }
    dispatch(
      listAccounts({
        environmentNamespace: currentEnvironment.namespace,
        pageSize: FEATURE_ACCOUNT_PAGE_SIZE,
        cursor: '',
      })
    );
    updateFeatureList(
      searchOptions,
      searchOptions.page ? Number(searchOptions.page) : 1
    );
  }, [dispatch, updateFeatureList]);

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.feature.list.header.title)}
          description={f(messages.feature.list.header.description)}
        />
      </div>
      <div className="m-10">
        <FeatureList
          searchOptions={searchOptions}
          onChangePage={handlePageChange}
          onSwitchEnabled={handleClickSwitchEnabled}
          onArchive={handleClickArchive}
          onClone={handleClickClone}
          onAdd={handleOpen}
          onChangeSearchOptions={handleSearchOptionsChange}
          onClearSearchOptions={handleClearSearchOptionsChange}
        />
      </div>
      <Overlay open={open} onClose={handleClose}>
        {isNew && (
          <FormProvider {...addMethod}>
            <FeatureAddForm
              onSubmit={handleAddSubmit(handleAdd)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
        {isClone && (
          <FormProvider {...cloneMethod}>
            <FeatureCloneForm
              onSubmit={cloneHandleSubmit(handleClone)}
              onCancel={handleClose}
            />
          </FormProvider>
        )}
      </Overlay>
      <FormProvider {...switchEnabledMethod}>
        <FeatureConfirmDialog
          open={isSwitchEnableConfirmDialogOpen}
          handleSubmit={switchEnableHandleSubmit(handleSwitchEnabled)}
          onClose={() => setIsSwitchEnableConfirmDialogOpen(false)}
          title={f(messages.feature.confirm.title)}
          description={f(messages.feature.confirm.description)}
        />
      </FormProvider>
      <FormProvider {...archiveMethod}>
        <FeatureConfirmDialog
          open={isArchiveConfirmDialogOpen}
          handleSubmit={archiveHandleSubmit(handleArchive)}
          onClose={() => setIsArchiveConfirmDialogOpen(false)}
          title={
            archiveMethod.getValues().feature &&
            archiveMethod.getValues().feature.archived
              ? f(messages.feature.confirm.unarchiveTitle)
              : f(messages.feature.confirm.archiveTitle)
          }
          description={
            archiveMethod.getValues().feature &&
            archiveMethod.getValues().feature.archived
              ? f(messages.feature.confirm.unarchiveDescription, {
                  featureId:
                    archiveMethod.getValues().feature &&
                    archiveMethod.getValues().feature.id,
                })
              : f(messages.feature.confirm.archiveDescription, {
                  featureId:
                    archiveMethod.getValues().feature &&
                    archiveMethod.getValues().feature.id,
                })
          }
        />
      </FormProvider>
    </>
  );
});
