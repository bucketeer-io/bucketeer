import {
  KeyboardEvent,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState
} from 'react';
import {
  IconAddOutlined,
  IconArchiveOutlined,
  IconEditOutlined
} from 'react-icons-material-design';
import {
  environmentCreator,
  EnvironmentCreatorCommand,
  EnvironmentsFetcherParams
} from '@api/environment';
import { environmentArchive } from '@api/environment/environment-archive';
import { environmentUpdate } from '@api/environment/environment-update';
import { ProjectFetcherParams } from '@api/project';
import { useQueryEnvironments } from '@queries/environments';
import { useQueryProjects } from '@queries/projects';
import { LIST_PAGE_SIZE } from 'constants/app';
import { FormFieldProps } from 'containers/common-form';
import CommonSlider from 'containers/common-slider';
import Filter from 'containers/filter';
import { onSubmitSuccess, SubmitType } from 'containers/pages/organizations';
import { SortingType } from 'containers/pages/projects';
import TableContent from 'containers/table-content';
import { useToast } from 'hooks';
import { ColumnType } from 'hooks/use-table';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import * as yup from 'yup';
import { Environment, EnvironmentCollection, OrderBy } from '@types';
import { sortingFn } from 'utils/sort';
import { IconInfo } from '@icons';
import { Button } from 'components/button';
import Checkbox from 'components/checkbox';
import Divider from 'components/divider';
import Icon from 'components/icon';
import { PopoverValue } from 'components/popover';

export const ProjectEnvironments = ({ projectId }: { projectId?: string }) => {
  const { notify } = useToast();
  const { t } = useTranslation(['common', 'form', 'table']);

  const initLoadedRef = useRef(true);

  const defaultParams: EnvironmentsFetcherParams = useMemo(
    () => ({
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      orderBy: 'DEFAULT',
      orderDirection: 'ASC',
      searchKeyword: '',
      disabled: false,
      archived: false,
      projectId
    }),
    [projectId]
  );

  const [isOpenSlider, setIsOpenSlider] = useState(false);
  const [isLoadingSubmit, setIsLoadingSubmit] = useState(false);
  const [submitType, setSubmitType] = useState<SubmitType>('created');
  const [envSelected, setEnvSelected] = useState<Environment | undefined>({
    projectId: String(projectId)
  } as Environment);
  const [sortingState, setSortingState] = useState<SortingType>({
    id: 'default',
    orderBy: 'DEFAULT',
    orderDirection: 'ASC'
  });
  const [environmentData, setEnvironmentData] =
    useState<EnvironmentCollection>();
  const [cursor, setCursor] = useState(0);
  const [searchValue, setSearchValue] = useState('');
  const [environmentParams, setEnvironmentParams] =
    useState<EnvironmentsFetcherParams>(defaultParams);

  const projectParams: ProjectFetcherParams = useMemo(
    () => ({
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      orderBy: 'DEFAULT',
      orderDirection: 'ASC',
      searchKeyword: '',
      disabled: false,
      organizationIds: []
    }),
    []
  );

  const { data, isLoading, refetch } = useQueryEnvironments({
    params: environmentParams,
    enabled: !!projectId
  });

  const { data: projectList } = useQueryProjects({
    params: projectParams,
    enabled: !!isOpenSlider
  });

  const columns = useMemo<ColumnType<Environment>[]>(
    () => [
      {
        accessorKey: 'name',
        id: 'name',
        header: `${t('name')}`,
        size: '35%',
        minSize: 160,
        sorting: true,
        cellType: 'title'
      },
      {
        accessorKey: 'featureFlagCount',
        id: 'featureFlagCount',
        header: `${t('table:flags')}`,
        size: '30%',
        minSize: 160,
        headerCellType: 'title',
        expandable: true,
        sorting: true
      },
      {
        accessorKey: 'createdAt',
        id: 'createdAt',
        header: `${t('table:created-at')}`,
        size: '30%',
        minSize: 160,
        sorting: true,
        headerCellType: 'title'
      },

      {
        accessorKey: 'actions',
        id: 'actions',
        header: 'Actions',
        size: '5%',
        minSize: 52,
        headerCellType: 'empty',
        cellType: 'icon',
        options: [
          {
            label: `${t('table:popover.edit-env')}`,
            icon: IconEditOutlined,
            value: 'edit-env'
          },
          {
            label: `${t('table:popover.archive-env', {
              ns: 'table'
            })}`,
            icon: IconArchiveOutlined,
            value: 'archive-env'
          }
        ],
        onClickPopover: (value, row) => handleClickPopover(value, row)
      }
    ],
    []
  );

  const formFields: FormFieldProps[] = useMemo(
    () => [
      {
        name: 'name',
        label: `${t('name')}`,
        placeholder: `${t('form:placeholder-name')}`,
        isRequired: true,
        isExpand: true,
        fieldType: 'input'
      },
      {
        name: 'urlCode',
        label: `${t('form:url-code')}`,
        placeholder: `${t('form:placeholder-code')}`,
        isRequired: true,
        isExpand: true,
        fieldType: 'input',
        isDisabled: submitType === 'updated' ? true : false,
        labelIcon: (
          <div className="inline-flex items-center ml-2.5">
            <Icon icon={IconInfo} size={'fit'} />
          </div>
        )
      },
      {
        name: 'projectId',
        label: `${t(`form:project-id`)}`,
        placeholder: `${t('form:placeholder-prj-id')}`,
        isRequired: true,
        isExpand: true,
        isDisabled: submitType === 'created' ? false : true,
        fieldType: 'dropdown',
        dropdownOptions: projectList?.projects?.map(item => ({
          label: item.name,
          value: item.id
        }))
      },
      {
        name: 'description',
        label: `${t(`form:description`)}`,
        placeholder: `${t('form:placeholder-desc')}`,
        isOptional: true,
        isExpand: true,
        fieldType: 'textarea'
      },
      {
        name: 'requireComment',
        isExpand: true,
        fieldType: 'additional',
        render: field => (
          <div className="flex flex-col w-full gap-y-5">
            <Divider />
            <h3 className="typo-head-bold-small text-gray-900">
              {t(`form:env-settings`)}
            </h3>
            <Checkbox
              onCheckedChange={checked => field.onChange(checked)}
              checked={field.value}
              title={`${t(`form:require-comments-flag`)}`}
              {...field}
            />
          </div>
        )
      }
    ],
    [submitType]
  );

  const formSchema = useMemo(
    () =>
      yup.object().shape({
        name: yup.string().required(),
        urlCode: yup.string().required(),
        projectId: yup.string().required()
      }),
    []
  );

  const handleChangeSearchValue = useCallback(
    debounce((value: string) => {
      initLoadedRef.current = false;
      setSearchValue(value);
    }),
    []
  );

  const handleKeyDown = useCallback((e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') handleChangeSearchValue.cancel();
  }, []);

  const handleOnSubmit = useCallback(
    (formValues: EnvironmentCreatorCommand) => {
      setIsLoadingSubmit(true);
      return (
        submitType === 'created'
          ? environmentCreator({
              command: {
                ...formValues
              }
            })
          : environmentUpdate({
              id: formValues.id,
              renameCommand: {
                name: formValues?.name
              },
              changeDescriptionCommand: {
                description: formValues?.description
              },
              changeRequireCommentCommand: {
                requireComment: formValues?.requireComment
              }
            })
      )
        .then(() =>
          onSubmitSuccess({
            name: formValues.name,
            submitType,
            notify,
            cb: () => {
              refetch();
              setIsOpenSlider(false);
              setEnvSelected(undefined);
            }
          })
        )
        .finally(() => {
          setIsLoadingSubmit(false);
        });
    },
    [submitType]
  );

  const handleClickPopover = useCallback(
    (value: PopoverValue, row?: Environment) => {
      if (value === 'edit-env') {
        if (submitType !== 'updated') setSubmitType('updated');
        setIsOpenSlider(true);
        return setEnvSelected(row);
      }
      handleArchiveEnv(row);
    },
    [submitType]
  );

  const handleArchiveEnv = useCallback((row?: Environment) => {
    environmentArchive({
      id: row?.id || '',
      command: {}
    }).then(() =>
      onSubmitSuccess({
        name: row?.name || '',
        submitType: 'archived',
        notify,
        cb: () => {
          refetch();
        }
      })
    );
  }, []);

  const onSortingTable = useCallback(
    (accessorKey: string, sortingKey?: OrderBy) => {
      initLoadedRef.current = false;
      sortingFn({
        accessorKey,
        sortingKey,
        sortingState,
        setSortingState
      });
    },
    [data, sortingState]
  );

  useEffect(() => {
    if (data) {
      setEnvironmentData(data);
    }
  }, [data]);

  useEffect(() => {
    if ((cursor >= 0 || sortingState) && !initLoadedRef.current) {
      setEnvironmentParams(prev => ({
        ...prev,
        searchKeyword: searchValue.trim().toLowerCase(),
        orderBy: sortingState.orderBy,
        orderDirection: sortingState.orderDirection,
        cursor: String(cursor)
      }));
    }
  }, [cursor, sortingState, searchValue]);

  return (
    <div>
      <Filter
        searchValue={searchValue}
        onChangeSearchValue={handleChangeSearchValue}
        onKeyDown={handleKeyDown}
        additionalActions={
          <Button
            className="flex-1 lg:flex-none"
            onClick={() => {
              setIsOpenSlider(true);
              if (submitType !== 'created') setSubmitType('created');
            }}
          >
            <Icon icon={IconAddOutlined} size="sm" />
            {t('new-env')}
          </Button>
        }
      />
      <TableContent
        isLoading={isLoading && initLoadedRef.current}
        columns={columns}
        data={
          environmentData?.environments?.length
            ? environmentData.environments
            : data?.environments
        }
        emptyTitle={t(`table:empty.env-title`)}
        emptyDescription={t(`table:empty.env-desc`)}
        paginationProps={{
          cursor,
          pageSize: LIST_PAGE_SIZE,
          totalCount: environmentData?.totalCount
            ? Number(environmentData?.totalCount)
            : 0,
          setCursor,
          cb: () => (initLoadedRef.current = false)
        }}
        sortingState={sortingState}
        onSortingTable={onSortingTable}
      />
      <CommonSlider
        title={t(submitType === 'created' ? 'new-project' : 'update-project')}
        isLoading={isLoadingSubmit}
        isOpen={isOpenSlider}
        formFields={formFields}
        submitTextBtn={t(
          submitType === 'created' ? 'create-project' : 'update-project'
        )}
        formSchema={formSchema}
        formData={envSelected}
        onSubmit={formValues =>
          handleOnSubmit(formValues as EnvironmentCreatorCommand)
        }
        onClose={() => {
          setEnvSelected(undefined);
          setIsOpenSlider(false);
        }}
      />
    </div>
  );
};
